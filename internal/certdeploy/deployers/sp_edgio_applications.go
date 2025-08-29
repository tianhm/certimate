package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	edgioapps "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/edgio-applications"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeEdgioApplications, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForEdgio{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := edgioapps.NewSSLDeployerProvider(&edgioapps.SSLDeployerProviderConfig{
			ClientId:      access.ClientId,
			ClientSecret:  access.ClientSecret,
			EnvironmentId: xmaps.GetString(options.ProviderConfig, "environmentId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
