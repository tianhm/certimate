package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/cdnfly"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeCdnfly, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForCdnfly{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		deployer, err := cdnfly.NewSSLDeployerProvider(&cdnfly.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiKey:                   access.ApiKey,
			ApiSecret:                access.ApiSecret,
			AllowInsecureConnections: access.AllowInsecureConnections,
			ResourceType:             cdnfly.ResourceType(xmaps.GetOrDefaultString(options.ProviderConfig, "resourceType", string(cdnfly.RESOURCE_TYPE_SITE))),
			SiteId:                   xmaps.GetString(options.ProviderConfig, "siteId"),
			CertificateId:            xmaps.GetString(options.ProviderConfig, "certificateId"),
		})
		return deployer, err
	}); err != nil {
		panic(err)
	}
}
