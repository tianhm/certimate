package qiniupili

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/qiniu/go-sdk/v7/pili"
	"github.com/samber/lo"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/qiniu-sslcert"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type SSLDeployerProviderConfig struct {
	// 七牛云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 七牛云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 直播空间名。
	Hub string `json:"hub"`
	// 域名匹配模式。
	// 零值时默认值 [DOMAIN_MATCH_PATTERN_EXACT]。
	DomainMatchPattern string `json:"domainMatchPattern,omitempty"`
	// 直播流域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *pili.Manager
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	manager := pili.NewManager(pili.ManagerConfig{AccessKey: config.AccessKey, SecretKey: config.SecretKey})

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
		sdkClient:  manager,
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
	if d.config.Domain == "" {
		return nil, fmt.Errorf("config `domain` is required")
	}

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

	case DOMAIN_MATCH_PATTERN_CERTSAN:
		{
			certX509, err := xcert.ParseCertificateFromPEM(certPEM)
			if err != nil {
				return nil, err
			}

			domainCandidates, err := d.getAllDomainsByHub(ctx, d.config.Hub)
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
				if err := d.updateDomainCertificate(ctx, d.config.Hub, domain, upres.CertName); err != nil {
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

func (d *SSLDeployerProvider) getAllDomainsByHub(ctx context.Context, hub string) ([]string, error) {
	domains := make([]string, 0)

	// 查询域名列表
	// REF: https://developer.qiniu.com/pili/9910/pili-service-sdk#6
	getDomainListReq := pili.GetDomainsListRequest{
		Hub: hub,
	}
	getDomainListResp, err := d.sdkClient.GetDomainsList(ctx, getDomainListReq)
	d.logger.Debug("sdk request 'pili.GetDomainsList'", slog.Any("request", getDomainListReq), slog.Any("response", getDomainListResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'pili.GetDomainsList': %w", err)
	}

	for _, domainItem := range getDomainListResp.Domains {
		domains = append(domains, domainItem.Domain)
	}

	return domains, nil
}

func (d *SSLDeployerProvider) updateDomainCertificate(ctx context.Context, hub string, domain string, cloudCertName string) error {
	// 修改域名证书配置
	// REF: https://developer.qiniu.com/pili/9910/pili-service-sdk#6
	setDomainCertReq := pili.SetDomainCertRequest{
		Hub:      hub,
		Domain:   domain,
		CertName: cloudCertName,
	}
	err := d.sdkClient.SetDomainCert(ctx, setDomainCertReq)
	d.logger.Debug("sdk request 'pili.SetDomainCert'", slog.Any("request", setDomainCertReq))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'pili.SetDomainCert': %w", err)
	}

	return nil
}
