package ksyuncdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/KscSDK/ksc-sdk-go/ksc"
	ksccdnv1 "github.com/KscSDK/ksc-sdk-go/service/cdnv1"
	"github.com/go-viper/mapstructure/v2"

	"github.com/certimate-go/certimate/pkg/core/deployer"
)

type DeployerConfig struct {
	// 金山云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 金山云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 域名匹配模式。暂时只支持精确匹配。
	// 零值时默认值 [DOMAIN_MATCH_PATTERN_EXACT]。
	DomainMatchPattern string `json:"domainMatchPattern,omitempty"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
	// 证书 ID。
	// 选填。零值时表示新建证书；否则表示更新证书。
	CertificateId string `json:"certificateId,omitempty"`
}

type Deployer struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *ksccdnv1.Cdnv1
}

var _ deployer.Provider = (*Deployer)(nil)

func NewDeployer(config *DeployerConfig) (*Deployer, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &Deployer{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *Deployer) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

func (d *Deployer) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 如果原证书 ID 为空，则创建证书；否则更新证书。
	if d.config.CertificateId == "" {
		if err := d.deployToDomain(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}
	} else {
		if err := d.deployToCertificate(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}
	}

	return &deployer.DeployResult{}, nil
}

func (d *Deployer) deployToDomain(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 获取域名 ID
	domainId, err := d.findDomainIdByDomain(ctx, d.config.Domain)
	if err != nil {
		return err
	}

	if err := d.updateDomainCertificate(ctx, domainId, certPEM, privkeyPEM); err != nil {
		return err
	}

	return nil
}

func (d *Deployer) deployToCertificate(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.CertificateId == "" {
		return errors.New("config `certificateId` is required")
	}

	// 更新证书
	// https://docs.ksyun.com/documents/259
	setCertificateInput := map[string]any{
		"CertificateId":     d.config.CertificateId,
		"CertificateName":   fmt.Sprintf("certimate_%d", time.Now().UnixMilli()),
		"ServerCertificate": certPEM,
		"PrivateKey":        privkeyPEM,
	}
	setCertificateReq, setCertificateOutput := d.sdkClient.SetCertificatePostRequest(&setCertificateInput)
	setCertificateErr := setCertificateReq.Send()
	d.logger.Debug("sdk request 'cdn.SetCertificate'", slog.Any("request", setCertificateInput), slog.Any("response", setCertificateOutput))
	if setCertificateErr != nil {
		return fmt.Errorf("failed to execute sdk request 'cdn.SetCertificate': %w", setCertificateErr)
	}

	return nil
}

func (d *Deployer) findDomainIdByDomain(ctx context.Context, domain string) (string, error) {
	// 查询域名列表
	// https://docs.ksyun.com/documents/198
	getCdnDomainsPageNumber := 1
	getCdnDomainsPageSize := 100
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		getCdnDomainsInput := map[string]any{
			"PageNumber": getCdnDomainsPageNumber,
			"PageSize":   getCdnDomainsPageSize,
			"DomainName": domain,
			"FuzzyMatch": "off",
		}
		getCdnDomainsReq, getCdnDomainsOutput := d.sdkClient.GetCdnDomainsPostRequest(&getCdnDomainsInput)
		getCdnDomainsErr := getCdnDomainsReq.Send()
		d.logger.Debug("sdk request 'cdn.GetCdnDomains'", slog.Any("request", getCdnDomainsInput), slog.Any("response", getCdnDomainsOutput))
		if getCdnDomainsErr != nil {
			return "", fmt.Errorf("failed to execute sdk request 'cdn.GetCdnDomains': %w", getCdnDomainsErr)
		}

		type GetCdnDomainsResponse struct {
			PageNumber int32 `json:"PageNumber"`
			PageSize   int32 `json:"PageSize"`
			TotalCount int32 `json:"TotalCount"`
			Domains    []*struct {
				DomainId     string `json:"DomainId"`
				DomainName   string `json:"DomainName"`
				Cname        string `json:"Cname"`
				CdnType      string `json:"CdnType"`
				CreatedTime  string `json:"CreatedTime"`
				ModifiedTime string `json:"ModifiedTime"`
				Region       string `json:"Region"`
			} `json:"Domains"`
		}
		var getCdnDomainsResp *GetCdnDomainsResponse
		mapstructure.Decode(getCdnDomainsOutput, &getCdnDomainsResp)
		if getCdnDomainsResp == nil {
			break
		}

		for _, domainItem := range getCdnDomainsResp.Domains {
			if strings.EqualFold(domainItem.DomainName, domain) {
				return domainItem.DomainId, nil
			}
		}

		if len(getCdnDomainsResp.Domains) < getCdnDomainsPageSize {
			break
		}

		getCdnDomainsPageNumber++
	}

	return "", fmt.Errorf("could not find domain '%s'", domain)
}

func (d *Deployer) updateDomainCertificate(ctx context.Context, domainId string, certPEM string, privkeyPEM string) error {
	// 为加速域名配置证书接口
	// https://docs.ksyun.com/documents/261
	configCertificateInput := map[string]any{
		"Enable":            "on",
		"DomainIds":         domainId,
		"CertificateName":   fmt.Sprintf("certimate_%d", time.Now().UnixMilli()),
		"ServerCertificate": certPEM,
		"PrivateKey":        privkeyPEM,
	}
	configCertificateReq, configCertificateOutput := d.sdkClient.ConfigCertificatePostRequest(&configCertificateInput)
	configCertificateErr := configCertificateReq.Send()
	d.logger.Debug("sdk request 'cdn.ConfigCertificate'", slog.Any("request", configCertificateInput), slog.Any("response", configCertificateOutput))
	if configCertificateErr != nil {
		return fmt.Errorf("failed to execute sdk request 'cdn.ConfigCertificate': %w", configCertificateErr)
	}

	return nil
}

func createSDKClient(accessKeyId, secretAccessKey string) (*ksccdnv1.Cdnv1, error) {
	region := "cn-beijing-6"
	client := ksccdnv1.SdkNew(ksc.NewClient(accessKeyId, secretAccessKey), &ksc.Config{Region: &region})
	return client, nil
}
