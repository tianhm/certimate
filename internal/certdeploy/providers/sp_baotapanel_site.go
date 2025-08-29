package providers

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	baotapanelsite "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/baotapanel-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeBaotaPanelSite, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForBaotaPanel{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := baotapanelsite.NewSSLDeployerProvider(&baotapanelsite.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiKey:                   access.ApiKey,
			AllowInsecureConnections: access.AllowInsecureConnections,
			SiteType:                 xmaps.GetOrDefaultString(options.ProviderConfig, "siteType", "other"),
			SiteName:                 xmaps.GetString(options.ProviderConfig, "siteName"),
			SiteNames:                lo.Filter(strings.Split(xmaps.GetString(options.ProviderConfig, "siteNames"), ";"), func(s string, _ int) bool { return s != "" }),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
