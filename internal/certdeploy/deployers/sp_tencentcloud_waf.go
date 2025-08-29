package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	tencentcloudwaf "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-waf"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeTencentCloudWAF, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForTencentCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := tencentcloudwaf.NewSSLDeployerProvider(&tencentcloudwaf.SSLDeployerProviderConfig{
			SecretId:   access.SecretId,
			SecretKey:  access.SecretKey,
			Endpoint:   xmaps.GetString(options.ProviderConfig, "endpoint"),
			Domain:     xmaps.GetString(options.ProviderConfig, "domain"),
			DomainId:   xmaps.GetString(options.ProviderConfig, "domainId"),
			InstanceId: xmaps.GetString(options.ProviderConfig, "instanceId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
