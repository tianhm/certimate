package ctcccloudcdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/samber/lo"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/ctcccloud-cdn"
	ctyuncdn "github.com/certimate-go/certimate/pkg/sdk3rd/ctyun/cdn"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
	xcerthostname "github.com/certimate-go/certimate/pkg/utils/cert/hostname"
)

type SSLDeployerProviderConfig struct {
	// 天翼云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 天翼云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 域名匹配模式。
	// 零值时默认值 [DOMAIN_MATCH_PATTERN_EXACT]。
	DomainMatchPattern string `json:"domainMatchPattern,omitempty"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *ctyuncdn.Client
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
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
				if err := d.updateDomainCertificate(ctx, domain, upres.CertId, upres.CertName); err != nil {
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

	// 遍历查询域名列表
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13816&data=174&isNormal=1&vid=167
	queryDomainListPage := int32(1)
	queryDomainListPageSize := int32(100)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		queryDomainListReq := &ctyuncdn.QueryDomainListRequest{
			Page:        lo.ToPtr(queryDomainListPage),
			PageSize:    lo.ToPtr(queryDomainListPageSize),
			ProductCode: lo.ToPtr("020"),
		}
		queryDomainListResp, err := d.sdkClient.QueryDomainList(queryDomainListReq)
		d.logger.Debug("sdk request 'cdn.QueryDomainList'", slog.Any("request", queryDomainListReq), slog.Any("response", queryDomainListResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.QueryDomainList': %w", err)
		}

		if queryDomainListResp.ReturnObj != nil {
			filteredProductCodes := []string{"001", "003", "004", "008"}
			ignoredStatuses := []int32{1, 5, 6, 7, 8, 9, 11, 12}
			for _, domainInfo := range queryDomainListResp.ReturnObj.Results {
				if !lo.Contains(filteredProductCodes, domainInfo.ProductCode) {
					continue
				}
				if lo.Contains(ignoredStatuses, domainInfo.Status) {
					continue
				}

				domains = append(domains, domainInfo.Domain)
			}
		}

		if queryDomainListResp.ReturnObj == nil || len(queryDomainListResp.ReturnObj.Results) < int(queryDomainListPageSize) {
			break
		} else {
			queryDomainListPage++
		}
	}

	return domains, nil
}

func (d *SSLDeployerProvider) updateDomainCertificate(ctx context.Context, domain string, cloudCertId, cloudCertName string) error {
	// 查询域名配置信息
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11304&data=161&isNormal=1&vid=154
	queryDomainDetailReq := &ctyuncdn.QueryDomainDetailRequest{
		Domain: lo.ToPtr(domain),
	}
	queryDomainDetailResp, err := d.sdkClient.QueryDomainDetail(queryDomainDetailReq)
	d.logger.Debug("sdk request 'cdn.QueryDomainDetail'", slog.Any("request", queryDomainDetailReq), slog.Any("response", queryDomainDetailResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'cdn.QueryDomainDetail': %w", err)
	}

	// 修改域名配置
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11308&data=161&isNormal=1&vid=154
	updateDomainReq := &ctyuncdn.UpdateDomainRequest{
		Domain:      lo.ToPtr(domain),
		HttpsStatus: lo.ToPtr("on"),
		CertName:    lo.ToPtr(cloudCertName),
	}
	updateDomainResp, err := d.sdkClient.UpdateDomain(updateDomainReq)
	d.logger.Debug("sdk request 'cdn.UpdateDomain'", slog.Any("request", updateDomainReq), slog.Any("response", updateDomainResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'cdn.UpdateDomain': %w", err)
	}

	return nil
}

func createSDKClient(accessKeyId, secretAccessKey string) (*ctyuncdn.Client, error) {
	return ctyuncdn.NewClient(accessKeyId, secretAccessKey)
}
