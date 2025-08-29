package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	bunnycdn "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/bunny-cdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeBunnyCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForBunny{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := bunnycdn.NewSSLDeployerProvider(&bunnycdn.SSLDeployerProviderConfig{
			ApiKey:     access.ApiKey,
			PullZoneId: xmaps.GetString(options.ProviderConfig, "pullZoneId"),
			Hostname:   xmaps.GetString(options.ProviderConfig, "hostname"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
