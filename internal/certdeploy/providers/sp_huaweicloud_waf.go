package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	huaweicloudwaf "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/huaweicloud-waf"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeHuaweiCloudWAF, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForHuaweiCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := huaweicloudwaf.NewSSLDeployerProvider(&huaweicloudwaf.SSLDeployerProviderConfig{
			AccessKeyId:         access.AccessKeyId,
			SecretAccessKey:     access.SecretAccessKey,
			EnterpriseProjectId: access.EnterpriseProjectId,
			Region:              xmaps.GetString(options.ProviderConfig, "region"),
			ResourceType:        huaweicloudwaf.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			CertificateId:       xmaps.GetString(options.ProviderConfig, "certificateId"),
			Domain:              xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
