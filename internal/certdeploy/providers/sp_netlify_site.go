package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	netlifysite "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/netlify-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeNetlifySite, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForNetlify{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := netlifysite.NewSSLDeployerProvider(&netlifysite.SSLDeployerProviderConfig{
			ApiToken: access.ApiToken,
			SiteId:   xmaps.GetString(options.ProviderConfig, "siteId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
