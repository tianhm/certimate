package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	wangsucdnpro "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/wangsu-cdnpro"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeWangsuCDNPro, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForWangsu{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := wangsucdnpro.NewSSLDeployerProvider(&wangsucdnpro.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			ApiKey:          access.ApiKey,
			Environment:     xmaps.GetOrDefaultString(options.ProviderConfig, "environment", "production"),
			Domain:          xmaps.GetString(options.ProviderConfig, "domain"),
			CertificateId:   xmaps.GetString(options.ProviderConfig, "certificateId"),
			WebhookId:       xmaps.GetString(options.ProviderConfig, "webhookId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
