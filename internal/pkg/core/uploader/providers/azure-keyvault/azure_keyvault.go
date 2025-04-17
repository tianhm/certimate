﻿package azurekeyvault

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
	azcommon "github.com/usual2970/certimate/internal/pkg/vendors/azure-sdk/common"
)

type UploaderConfig struct {
	// Azure TenantId。
	TenantId string `json:"tenantId"`
	// Azure ClientId。
	ClientId string `json:"clientId"`
	// Azure ClientSecret。
	ClientSecret string `json:"clientSecret"`
	// Azure 主权云环境。
	CloudName string `json:"cloudName,omitempty"`
	// Key Vault 名称。
	KeyVaultName string `json:"keyvaultName"`
	// Certificate 名称。
	CertificateName string `json:"certificateName,omitempty"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *azcertificates.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.TenantId, config.ClientId, config.ClientSecret, config.CloudName, config.KeyVaultName)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &UploaderProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) WithLogger(logger *slog.Logger) uploader.Uploader {
	if logger == nil {
		u.logger = slog.Default()
	} else {
		u.logger = logger
	}
	return u
}

func (u *UploaderProvider) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 生成 Azure 业务参数
	const TAG_CERTCN = "certimate/cert-cn"
	const TAG_CERTSN = "certimate/cert-sn"
	certCN := certX509.Subject.CommonName
	certSN := certX509.SerialNumber.Text(16)

	// 获取证书列表，避免重复上传
	// REF: https://learn.microsoft.com/en-us/rest/api/keyvault/certificates/get-certificates/get-certificates
	listCertificatesPager := u.sdkClient.NewListCertificatesPager(nil)
	for listCertificatesPager.More() {
		page, err := listCertificatesPager.NextPage(context.TODO())
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'keyvault.GetCertificates'")
		}

		for _, certItem := range page.Value {
			// 如果已经指定了证书名称，则跳过证书名称不匹配的证书
			if u.config.CertificateName != "" && certItem.ID.Name() != u.config.CertificateName {
				continue
			}

			// 先对比证书有效期
			if certItem.Attributes == nil {
				continue
			}
			if certItem.Attributes.NotBefore == nil || !certItem.Attributes.NotBefore.Equal(certX509.NotBefore) {
				continue
			}
			if certItem.Attributes.Expires == nil || !certItem.Attributes.Expires.Equal(certX509.NotAfter) {
				continue
			}

			// 再对比 Tag 中的通用名称
			if v, ok := certItem.Tags[TAG_CERTCN]; !ok || v == nil {
				continue
			} else if *v != certCN {
				continue
			}

			// 再对比 Tag 中的序列号
			if v, ok := certItem.Tags[TAG_CERTSN]; !ok || v == nil {
				continue
			} else if *v != certSN {
				continue
			}

			// 最后对比证书内容
			getCertificateResp, err := u.sdkClient.GetCertificate(context.TODO(), certItem.ID.Name(), certItem.ID.Version(), nil)
			u.logger.Debug("sdk request 'keyvault.GetCertificate'", slog.String("request.certificateName", certItem.ID.Name()), slog.String("request.certificateVersion", certItem.ID.Version()), slog.Any("response", getCertificateResp))
			if err != nil {
				return nil, xerrors.Wrap(err, "failed to execute sdk request 'keyvault.GetCertificate'")
			} else {
				oldCertX509, err := x509.ParseCertificate(getCertificateResp.CER)
				if err != nil {
					continue
				}

				if !certutil.EqualCertificate(certX509, oldCertX509) {
					continue
				}
			}

			// 如果以上信息都一致，则视为已存在相同证书，直接返回
			u.logger.Info("ssl certificate already exists")
			return &uploader.UploadResult{
				CertId:   string(*certItem.ID),
				CertName: certItem.ID.Name(),
			}, nil
		}
	}

	certName := u.config.CertificateName
	if certName == "" {
		// 未指定证书名称时，生成包含timestamp的新证书名（需符合 Azure 命名规则）
		certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())
	}

	// Azure Key Vault 不支持导入带有Certificiate Chain的PEM证书。
	// Issue Link: https://github.com/Azure/azure-cli/issues/19017
	// 暂时的解决方法是，将 PEM 证书转换成 PFX 格式，然后再导入。
	pfxCert, err := certutil.TransformCertificateFromPEMToPFX(certPem, privkeyPem, "")
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to transform certificate from PEM to PFX")
	}

	// 导入证书
	// REF: https://learn.microsoft.com/en-us/rest/api/keyvault/certificates/import-certificate/import-certificate
	importCertificateParams := azcertificates.ImportCertificateParameters{
		Base64EncodedCertificate: to.Ptr(base64.StdEncoding.EncodeToString(pfxCert)),
		CertificatePolicy: &azcertificates.CertificatePolicy{
			SecretProperties: &azcertificates.SecretProperties{
				ContentType: to.Ptr("application/x-pkcs12"),
			},
		},
		Tags: map[string]*string{
			TAG_CERTCN: to.Ptr(certCN),
			TAG_CERTSN: to.Ptr(certSN),
		},
	}
	importCertificateResp, err := u.sdkClient.ImportCertificate(context.TODO(), certName, importCertificateParams, nil)
	u.logger.Debug("sdk request 'keyvault.ImportCertificate'", slog.String("request.certificateName", certName), slog.Any("request.parameters", importCertificateParams), slog.Any("response", importCertificateResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'keyvault.ImportCertificate'")
	}

	return &uploader.UploadResult{
		CertId:   string(*importCertificateResp.ID),
		CertName: certName,
	}, nil
}

func createSdkClient(tenantId, clientId, clientSecret, cloudName, keyvaultName string) (*azcertificates.Client, error) {
	env, err := azcommon.GetCloudEnvironmentConfiguration(cloudName)
	if err != nil {
		return nil, err
	}
	clientOptions := azcore.ClientOptions{Cloud: env}

	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret,
		&azidentity.ClientSecretCredentialOptions{ClientOptions: clientOptions})
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("https://%s.vault.azure.net", keyvaultName)
	if azcommon.IsEnvironmentGovernment(cloudName) {
		endpoint = fmt.Sprintf("https://%s.vault.usgovcloudapi.net", keyvaultName)
	} else if azcommon.IsEnvironmentChina(cloudName) {
		endpoint = fmt.Sprintf("https://%s.vault.azure.cn", keyvaultName)
	}

	client, err := azcertificates.NewClient(endpoint, credential, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
