package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	ratpanelsite "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ratpanel-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeRatPanelSite, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForRatPanel{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := ratpanelsite.NewSSLDeployerProvider(&ratpanelsite.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			AccessTokenId:            access.AccessTokenId,
			AccessToken:              access.AccessToken,
			AllowInsecureConnections: access.AllowInsecureConnections,
			SiteName:                 xmaps.GetString(options.ProviderConfig, "siteName"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
