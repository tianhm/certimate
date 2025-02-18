﻿package tencentcloudwaf

import (
	"context"
	"errors"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcWaf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type TencentCloudWAFDeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 腾讯云地域。
	Region string `json:"region"`
	// 防护域名（不支持泛域名）。
	Domain string `json:"domain"`
	// 防护域名 ID。
	DomainId string `json:"domainId"`
	// 防护域名所属实例 ID。
	InstanceId string `json:"instanceId"`
}

type TencentCloudWAFDeployer struct {
	config      *TencentCloudWAFDeployerConfig
	logger      logger.Logger
	sdkClient   *tcWaf.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*TencentCloudWAFDeployer)(nil)

func New(config *TencentCloudWAFDeployerConfig) (*TencentCloudWAFDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *TencentCloudWAFDeployerConfig, logger logger.Logger) (*TencentCloudWAFDeployer, error) {
	if config == nil {
		panic("config is nil")
	}

	if logger == nil {
		panic("logger is nil")
	}

	client, err := createSdkClient(config.SecretId, config.SecretKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := uploadersp.New(&uploadersp.TencentCloudSSLUploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &TencentCloudWAFDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCloudWAFDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}
	if d.config.DomainId == "" {
		return nil, errors.New("config `domainId` is required")
	}
	if d.config.InstanceId == "" {
		return nil, errors.New("config `instanceId` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Logt("certificate file uploaded", upres)
	}

	// 查询单个 SaaS 型 WAF 域名详情
	// REF: https://cloud.tencent.com/document/api/627/82938
	describeDomainDetailsSaasReq := tcWaf.NewDescribeDomainDetailsSaasRequest()
	describeDomainDetailsSaasReq.Domain = common.StringPtr(d.config.Domain)
	describeDomainDetailsSaasReq.DomainId = common.StringPtr(d.config.DomainId)
	describeDomainDetailsSaasReq.InstanceId = common.StringPtr(d.config.InstanceId)
	describeDomainDetailsSaasResp, err := d.sdkClient.DescribeDomainDetailsSaas(describeDomainDetailsSaasReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'waf.DescribeDomainDetailsSaas'")
	} else {
		d.logger.Logt("已查询到 SaaS 型 WAF 域名详情", describeDomainDetailsSaasResp.Response)
	}

	// 编辑 SaaS 型 WAF 域名
	// REF: https://cloud.tencent.com/document/api/627/94309
	modifySpartaProtectionReq := tcWaf.NewModifySpartaProtectionRequest()
	modifySpartaProtectionReq.Domain = common.StringPtr(d.config.Domain)
	modifySpartaProtectionReq.DomainId = common.StringPtr(d.config.DomainId)
	modifySpartaProtectionReq.InstanceID = common.StringPtr(d.config.InstanceId)
	modifySpartaProtectionReq.CertType = common.Int64Ptr(2)
	modifySpartaProtectionReq.SSLId = common.StringPtr(upres.CertId)
	modifySpartaProtectionResp, err := d.sdkClient.ModifySpartaProtection(modifySpartaProtectionReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'waf.ModifySpartaProtection'")
	} else {
		d.logger.Logt("已编辑 SaaS 型 WAF 域名", modifySpartaProtectionResp.Response)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(secretId, secretKey, region string) (*tcWaf.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcWaf.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
