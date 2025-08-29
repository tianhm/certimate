package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	opsite "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/1panel-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderType1PanelSite, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigFor1Panel{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := opsite.NewSSLDeployerProvider(&opsite.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiVersion:               access.ApiVersion,
			ApiKey:                   access.ApiKey,
			AllowInsecureConnections: access.AllowInsecureConnections,
			NodeName:                 xmaps.GetString(options.ProviderConfig, "nodeName"),
			ResourceType:             opsite.ResourceType(xmaps.GetOrDefaultString(options.ProviderConfig, "resourceType", string(opsite.RESOURCE_TYPE_WEBSITE))),
			WebsiteId:                xmaps.GetInt64(options.ProviderConfig, "websiteId"),
			CertificateId:            xmaps.GetInt64(options.ProviderConfig, "certificateId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
