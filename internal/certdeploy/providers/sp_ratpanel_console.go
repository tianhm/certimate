package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	ratpanelconsole "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ratpanel-console"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeRatPanelConsole, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForRatPanel{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := ratpanelconsole.NewSSLDeployerProvider(&ratpanelconsole.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			AccessTokenId:            access.AccessTokenId,
			AccessToken:              access.AccessToken,
			AllowInsecureConnections: access.AllowInsecureConnections,
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
