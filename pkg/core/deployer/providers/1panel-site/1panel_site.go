package onepanelsite

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/certimate-go/certimate/pkg/core/certmgr"
	mcertmgr "github.com/certimate-go/certimate/pkg/core/certmgr/providers/1panel"
	"github.com/certimate-go/certimate/pkg/core/deployer"
	onepanelsdk "github.com/certimate-go/certimate/pkg/sdk3rd/1panel"
	onepanelsdk2 "github.com/certimate-go/certimate/pkg/sdk3rd/1panel/v2"
)

type DeployerConfig struct {
	// 1Panel 服务地址。
	ServerUrl string `json:"serverUrl"`
	// 1Panel 版本。
	// 可取值 "v1"、"v2"。
	ApiVersion string `json:"apiVersion"`
	// 1Panel 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 子节点名称。
	// 选填。
	NodeName string `json:"nodeName,omitempty"`
	// 部署资源类型。
	ResourceType string `json:"resourceType"`
	// 网站 ID。
	// 部署资源类型为 [RESOURCE_TYPE_WEBSITE] 时必填。
	WebsiteId int64 `json:"websiteId,omitempty"`
	// 证书 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时必填。
	CertificateId int64 `json:"certificateId,omitempty"`
}

type Deployer struct {
	config     *DeployerConfig
	logger     *slog.Logger
	sdkClient  any
	sdkCertmgr certmgr.Provider
}

var _ deployer.Provider = (*Deployer)(nil)

func NewDeployer(config *DeployerConfig) (*Deployer, error) {
	if config == nil {
		return nil, errors.New("the configuration of the deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiVersion, config.ApiKey, config.AllowInsecureConnections, config.NodeName)
	if err != nil {
		return nil, fmt.Errorf("could not create client: %w", err)
	}

	pcertmgr, err := mcertmgr.NewCertmgr(&mcertmgr.CertmgrConfig{
		ServerUrl:                config.ServerUrl,
		ApiVersion:               config.ApiVersion,
		ApiKey:                   config.ApiKey,
		AllowInsecureConnections: config.AllowInsecureConnections,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create certmgr: %w", err)
	}

	return &Deployer{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
		sdkCertmgr: pcertmgr,
	}, nil
}

func (d *Deployer) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sdkCertmgr.SetLogger(logger)
}

func (d *Deployer) Deploy(ctx context.Context, certPEM, privkeyPEM string) (*deployer.DeployResult, error) {
	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_WEBSITE:
		if err := d.deployToWebsite(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	case RESOURCE_TYPE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type '%s'", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *Deployer) deployToWebsite(ctx context.Context, certPEM, privkeyPEM string) error {
	if d.config.WebsiteId == 0 {
		return errors.New("config `websiteId` is required")
	}

	// 上传证书
	upres, err := d.sdkCertmgr.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	switch sdkClient := d.sdkClient.(type) {
	case *onepanelsdk.Client:
		{
			// 获取网站 HTTPS 配置
			websiteHttpsGetResp, err := sdkClient.WebsiteHttpsGet(d.config.WebsiteId)
			d.logger.Debug("sdk request '1panel.WebsiteHttpsGet'", slog.Int64("websiteId", d.config.WebsiteId), slog.Any("response", websiteHttpsGetResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.WebsiteHttpsGet': %w", err)
			}

			// 修改网站 HTTPS 配置
			sslId, _ := strconv.ParseInt(upres.CertId, 10, 64)
			websiteHttpsPostReq := &onepanelsdk.WebsiteHttpsPostRequest{
				WebsiteID:    d.config.WebsiteId,
				Type:         "existed",
				WebsiteSSLID: sslId,
				Enable:       websiteHttpsGetResp.Data.Enable,
				HttpConfig:   websiteHttpsGetResp.Data.HttpConfig,
				SSLProtocol:  websiteHttpsGetResp.Data.SSLProtocol,
				Algorithm:    websiteHttpsGetResp.Data.Algorithm,
				Hsts:         websiteHttpsGetResp.Data.Hsts,
			}
			websiteHttpsPostResp, err := sdkClient.WebsiteHttpsPost(d.config.WebsiteId, websiteHttpsPostReq)
			d.logger.Debug("sdk request '1panel.WebsiteHttpsPost'", slog.Int64("websiteId", d.config.WebsiteId), slog.Any("request", websiteHttpsPostReq), slog.Any("response", websiteHttpsPostResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.WebsiteHttpsPost': %w", err)
			}
		}

	case *onepanelsdk2.Client:
		{
			// 获取网站 HTTPS 配置
			websiteHttpsGetResp, err := sdkClient.WebsiteHttpsGet(d.config.WebsiteId)
			d.logger.Debug("sdk request '1panel.WebsiteHttpsGet'", slog.Int64("websiteId", d.config.WebsiteId), slog.Any("response", websiteHttpsGetResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.WebsiteHttpsGet': %w", err)
			}

			// 修改网站 HTTPS 配置
			sslId, _ := strconv.ParseInt(upres.CertId, 10, 64)
			websiteHttpsPostReq := &onepanelsdk2.WebsiteHttpsPostRequest{
				WebsiteID:    d.config.WebsiteId,
				Type:         "existed",
				WebsiteSSLID: sslId,
				Enable:       websiteHttpsGetResp.Data.Enable,
				HttpConfig:   websiteHttpsGetResp.Data.HttpConfig,
				SSLProtocol:  websiteHttpsGetResp.Data.SSLProtocol,
				Algorithm:    websiteHttpsGetResp.Data.Algorithm,
				Hsts:         websiteHttpsGetResp.Data.Hsts,
			}
			websiteHttpsPostResp, err := sdkClient.WebsiteHttpsPost(d.config.WebsiteId, websiteHttpsPostReq)
			d.logger.Debug("sdk request '1panel.WebsiteHttpsPost'", slog.Int64("websiteId", d.config.WebsiteId), slog.Any("request", websiteHttpsPostReq), slog.Any("response", websiteHttpsPostResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.WebsiteHttpsPost': %w", err)
			}
		}

	default:
		panic("unreachable")
	}

	return nil
}

func (d *Deployer) deployToCertificate(ctx context.Context, certPEM, privkeyPEM string) error {
	if d.config.CertificateId == 0 {
		return errors.New("config `certificateId` is required")
	}

	// 替换证书
	opres, err := d.sdkCertmgr.Replace(ctx, fmt.Sprintf("%d", d.config.CertificateId), certPEM, privkeyPEM)
	if err != nil {
		return fmt.Errorf("failed to replace certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate replaced", slog.Any("result", opres))
	}

	return nil
}

const (
	sdkVersionV1 = "v1"
	sdkVersionV2 = "v2"
)

func createSDKClient(serverUrl, apiVersion, apiKey string, skipTlsVerify bool, nodeName string) (any, error) {
	if apiVersion == sdkVersionV1 {
		client, err := onepanelsdk.NewClient(serverUrl, apiKey)
		if err != nil {
			return nil, err
		}

		if skipTlsVerify {
			client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
		}

		return client, nil
	} else if apiVersion == sdkVersionV2 {
		var client *onepanelsdk2.Client
		if nodeName == "" {
			temp, err := onepanelsdk2.NewClient(serverUrl, apiKey)
			if err != nil {
				return nil, err
			}
			client = temp
		} else {
			temp, err := onepanelsdk2.NewClientWithNode(serverUrl, apiKey, nodeName)
			if err != nil {
				return nil, err
			}
			client = temp
		}

		if skipTlsVerify {
			client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
		}

		return client, nil
	}

	return nil, fmt.Errorf("invalid 1panel api version")
}
