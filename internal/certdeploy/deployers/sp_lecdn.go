package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/lecdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeLeCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForLeCDN{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := lecdn.NewSSLDeployerProvider(&lecdn.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiVersion:               access.ApiVersion,
			ApiRole:                  access.ApiRole,
			Username:                 access.Username,
			Password:                 access.Password,
			AllowInsecureConnections: access.AllowInsecureConnections,
			ResourceType:             lecdn.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			CertificateId:            xmaps.GetInt64(options.ProviderConfig, "certificateId"),
			ClientId:                 xmaps.GetInt64(options.ProviderConfig, "clientId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
