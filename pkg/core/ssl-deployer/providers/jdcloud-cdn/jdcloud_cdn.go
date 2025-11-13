package jdcloudcdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdcdn "github.com/jdcloud-api/jdcloud-sdk-go/services/cdn/apis"
	"github.com/samber/lo"

	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/jdcloud-cdn/internal"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/jdcloud-ssl"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
	xcerthostname "github.com/certimate-go/certimate/pkg/utils/cert/hostname"
)

type SSLDeployerProviderConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 域名匹配模式。
	// 零值时默认值 [DOMAIN_MATCH_PATTERN_EXACT]。
	DomainMatchPattern string `json:"domainMatchPattern,omitempty"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *internal.CdnClient
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
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
	var domains []string
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
				domainCandidates, err := d.getAllDomains(ctx)
				if err != nil {
					return nil, err
				}

				domains = lo.Filter(domainCandidates, func(domain string, _ int) bool {
					return xcerthostname.IsMatch(d.config.Domain, domain)
				})
				if len(domains) == 0 {
					return nil, errors.New("could not find any domains matched by wildcard")
				}
			} else {
				domains = []string{d.config.Domain}
			}
		}

	case DOMAIN_MATCH_PATTERN_CERTSAN:
		{
			certX509, err := xcert.ParseCertificateFromPEM(certPEM)
			if err != nil {
				return nil, err
			}

			domainCandidates, err := d.getAllDomains(ctx)
			if err != nil {
				return nil, err
			}

			domains = lo.Filter(domainCandidates, func(domain string, _ int) bool {
				return certX509.VerifyHostname(domain) == nil
			})
			if len(domains) == 0 {
				return nil, errors.New("could not find any domains matched by certificate")
			}
		}

	default:
		return nil, fmt.Errorf("unsupported domain match pattern: '%s'", d.config.DomainMatchPattern)
	}

	// 遍历更新域名证书
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

func (d *SSLDeployerProvider) getAllDomains(ctx context.Context) ([]string, error) {
	domains := make([]string, 0)

	// 查询域名列表
	// REF: https://docs.jdcloud.com/cn/cdn/api/getdomainlist
	getDomainListPageNumber := 1
	getDomainListPageSize := 50
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		getDomainListReq := jdcdn.NewGetDomainListRequestWithoutParam()
		getDomainListReq.SetPageNumber(getDomainListPageNumber)
		getDomainListReq.SetPageSize(getDomainListPageSize)
		getDomainListResp, err := d.sdkClient.GetDomainList(getDomainListReq)
		d.logger.Debug("sdk request 'cdn.GetDomainList'", slog.Any("request", getDomainListReq), slog.Any("response", getDomainListResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.GetDomainList': %w", err)
		}

		ignoredStatuses := []string{"offline"}
		for _, domainInfo := range getDomainListResp.Result.Domains {
			if lo.Contains(ignoredStatuses, domainInfo.Status) {
				continue
			}

			domains = append(domains, domainInfo.Domain)
		}

		if len(getDomainListResp.Result.Domains) < getDomainListPageSize {
			break
		} else {
			getDomainListPageNumber++
		}
	}

	return domains, nil
}

func (d *SSLDeployerProvider) updateDomainCertificate(ctx context.Context, domain string, cloudCertId string) error {
	// 查询域名配置信息
	// REF: https://docs.jdcloud.com/cn/cdn/api/querydomainconfig
	queryDomainConfigReq := jdcdn.NewQueryDomainConfigRequestWithoutParam()
	queryDomainConfigReq.SetDomain(d.config.Domain)
	queryDomainConfigResp, err := d.sdkClient.QueryDomainConfig(queryDomainConfigReq)
	d.logger.Debug("sdk request 'cdn.QueryDomainConfig'", slog.Any("request", queryDomainConfigReq), slog.Any("response", queryDomainConfigResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'cdn.QueryDomainConfig': %w", err)
	}

	// 设置通讯协议
	// REF: https://docs.jdcloud.com/cn/cdn/api/sethttptype
	setHttpTypeReq := jdcdn.NewSetHttpTypeRequestWithoutParam()
	setHttpTypeReq.SetDomain(d.config.Domain)
	setHttpTypeReq.SetHttpType("https")
	setHttpTypeReq.SetCertFrom("ssl")
	setHttpTypeReq.SetSslCertId(cloudCertId)
	setHttpTypeReq.SetJumpType(queryDomainConfigResp.Result.HttpsJumpType)
	setHttpTypeResp, err := d.sdkClient.SetHttpType(setHttpTypeReq)
	d.logger.Debug("sdk request 'cdn.QueryDomainConfig'", slog.Any("request", setHttpTypeReq), slog.Any("response", setHttpTypeResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'cdn.SetHttpType': %w", err)
	}

	return nil
}

func createSDKClient(accessKeyId, accessKeySecret string) (*internal.CdnClient, error) {
	clientCredentials := jdcore.NewCredentials(accessKeyId, accessKeySecret)
	client := internal.NewCdnClient(clientCredentials)
	return client, nil
}
