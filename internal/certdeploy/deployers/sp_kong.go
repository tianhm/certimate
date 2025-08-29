package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/kong"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeKong, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForKong{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := kong.NewSSLDeployerProvider(&kong.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiToken:                 access.ApiToken,
			AllowInsecureConnections: access.AllowInsecureConnections,
			ResourceType:             kong.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			Workspace:                xmaps.GetString(options.ProviderConfig, "workspace"),
			CertificateId:            xmaps.GetString(options.ProviderConfig, "certificateId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
