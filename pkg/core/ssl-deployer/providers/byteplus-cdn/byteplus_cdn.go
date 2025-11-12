package bytepluscdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	bpcdn "github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
	bp "github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/byteplus-cdn"
	xcerthostname "github.com/certimate-go/certimate/pkg/utils/cert/hostname"
)

type SSLDeployerProviderConfig struct {
	// BytePlus AccessKey。
	AccessKey string `json:"accessKey"`
	// BytePlus SecretKey。
	SecretKey string `json:"secretKey"`
	// 域名匹配模式。
	// 零值时默认值 [DOMAIN_MATCH_PATTERN_EXACT]。
	DomainMatchPattern string `json:"domainMatchPattern,omitempty"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *bpcdn.CDN
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client := bpcdn.NewInstance()
	client.Client.SetAccessKey(config.AccessKey)
	client.Client.SetSecretKey(config.SecretKey)

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &SSLDeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
		sslManager: sslmgr,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sslManager.SetLogger(logger)
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取待部署的域名列表
	domains := make([]string, 0)
	switch d.config.DomainMatchPattern {
	case "", DOMAIN_MATCH_PATTERN_EXACT:
		{
			if d.config.Domain == "" {
				return nil, errors.New("config `domain` is required")
			}

			domains = []string{d.config.Domain}
		}

	case DOMAIN_MATCH_PATTERN_WILDCARD:
		{
			if d.config.Domain == "" {
				return nil, errors.New("config `domain` is required")
			}

			if strings.HasPrefix(d.config.Domain, "*.") {
				domainCandidates, err := d.getMatchedDomainsByWildcard(ctx, d.config.Domain)
				if err != nil {
					return nil, err
				}

				domains = domainCandidates
			} else {
				domains = []string{d.config.Domain}
			}
		}

	case DOMAIN_MATCH_PATTERN_CERTSAN:
		{
			domainCandidates, err := d.getMatchedDomainsByCertId(ctx, upres.CertId)
			if err != nil {
				return nil, err
			}

			domains = domainCandidates
		}

	default:
		return nil, fmt.Errorf("unsupported domain match pattern: '%s'", d.config.DomainMatchPattern)
	}

	// 遍历绑定证书
	if len(domains) == 0 {
		d.logger.Info("no cdn domains to deploy")
	} else {
		d.logger.Info("found cdn domains to deploy", slog.Any("domains", domains))
		var errs []error

		for _, domain := range domains {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if err := d.updateDomainCertificate(ctx, domain, upres.CertId); err != nil {
					errs = append(errs, err)
				}
			}
		}

		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}

	return &core.SSLDeployResult{}, nil
}

func (d *SSLDeployerProvider) getMatchedDomainsByWildcard(ctx context.Context, wildcardDomain string) ([]string, error) {
	domains := make([]string, 0)

	// 遍历获取加速域名列表，获取匹配的域名
	// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/ListCdnDomains_en-us
	listCdnDomainsPageNum := int64(1)
	listCdnDomainsPageSize := int64(100)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCdnDomainsReq := &bpcdn.ListCdnDomainsRequest{
			Domain:   bp.String(strings.TrimPrefix(wildcardDomain, "*.")),
			Status:   bp.String("online"),
			PageNum:  bp.Int64(listCdnDomainsPageNum),
			PageSize: bp.Int64(listCdnDomainsPageSize),
		}
		listCdnDomainsResp, err := d.sdkClient.ListCdnDomains(listCdnDomainsReq)
		d.logger.Debug("sdk request 'cdn.ListCdnDomains'", slog.Any("request", listCdnDomainsReq), slog.Any("response", listCdnDomainsResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.ListCdnDomains': %w", err)
		}

		for _, domainInfo := range listCdnDomainsResp.Result.Data {
			if xcerthostname.IsMatch(wildcardDomain, domainInfo.Domain) {
				domains = append(domains, domainInfo.Domain)
			}
		}

		if len(listCdnDomainsResp.Result.Data) < int(listCdnDomainsPageSize) {
			break
		} else {
			listCdnDomainsPageSize++
		}
	}

	return domains, nil
}

func (d *SSLDeployerProvider) getMatchedDomainsByCertId(ctx context.Context, cloudCertId string) ([]string, error) {
	domains := make([]string, 0)

	// 获取指定证书可关联的域名
	// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-describecertconfig-9ea17
	describeCertConfigReq := &bpcdn.DescribeCertConfigRequest{
		CertId: cloudCertId,
	}
	describeCertConfigResp, err := d.sdkClient.DescribeCertConfig(describeCertConfigReq)
	d.logger.Debug("sdk request 'cdn.DescribeCertConfig'", slog.Any("request", describeCertConfigReq), slog.Any("response", describeCertConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.DescribeCertConfig': %w", err)
	}

	if describeCertConfigResp.Result.CertNotConfig != nil {
		for i := range describeCertConfigResp.Result.CertNotConfig {
			domains = append(domains, describeCertConfigResp.Result.CertNotConfig[i].Domain)
		}
	}

	if describeCertConfigResp.Result.OtherCertConfig != nil {
		for i := range describeCertConfigResp.Result.OtherCertConfig {
			domains = append(domains, describeCertConfigResp.Result.OtherCertConfig[i].Domain)
		}
	}

	if len(domains) == 0 {
		if len(describeCertConfigResp.Result.SpecifiedCertConfig) == 0 {
			return nil, errors.New("no domains matched by certificate")
		}
	}

	return domains, nil
}

func (d *SSLDeployerProvider) updateDomainCertificate(ctx context.Context, domain string, cloudCertId string) error {
	// 关联证书与加速域名
	// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-batchdeploycert
	batchDeployCertReq := &bpcdn.BatchDeployCertRequest{
		CertId: cloudCertId,
		Domain: domain,
	}
	batchDeployCertResp, err := d.sdkClient.BatchDeployCert(batchDeployCertReq)
	d.logger.Debug("sdk request 'cdn.BatchDeployCert'", slog.Any("request", batchDeployCertReq), slog.Any("response", batchDeployCertResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'cdn.BatchDeployCert': %w", err)
	}

	return nil
}
