package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	baotawafsite "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/baotawaf-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeBaotaWAFSite, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForBaotaWAF{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := baotawafsite.NewSSLDeployerProvider(&baotawafsite.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiKey:                   access.ApiKey,
			AllowInsecureConnections: access.AllowInsecureConnections,
			SiteName:                 xmaps.GetString(options.ProviderConfig, "siteName"),
			SitePort:                 xmaps.GetOrDefaultInt32(options.ProviderConfig, "sitePort", 443),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
