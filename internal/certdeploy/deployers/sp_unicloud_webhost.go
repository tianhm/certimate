package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	unicloudwebhost "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/unicloud-webhost"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeUniCloudWebHost, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForUniCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := unicloudwebhost.NewSSLDeployerProvider(&unicloudwebhost.SSLDeployerProviderConfig{
			Username:      access.Username,
			Password:      access.Password,
			SpaceProvider: xmaps.GetString(options.ProviderConfig, "spaceProvider"),
			SpaceId:       xmaps.GetString(options.ProviderConfig, "spaceId"),
			Domain:        xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
