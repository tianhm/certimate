package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	tencentcloudgaap "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-gaap"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeTencentCloudGAAP, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForTencentCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := tencentcloudgaap.NewSSLDeployerProvider(&tencentcloudgaap.SSLDeployerProviderConfig{
			SecretId:     access.SecretId,
			SecretKey:    access.SecretKey,
			Endpoint:     xmaps.GetString(options.ProviderConfig, "endpoint"),
			ResourceType: tencentcloudgaap.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			ProxyId:      xmaps.GetString(options.ProviderConfig, "proxyId"),
			ListenerId:   xmaps.GetString(options.ProviderConfig, "listenerId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
