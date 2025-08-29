package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/flexcdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeFlexCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForFlexCDN{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := flexcdn.NewSSLDeployerProvider(&flexcdn.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiRole:                  access.ApiRole,
			AccessKeyId:              access.AccessKeyId,
			AccessKey:                access.AccessKey,
			AllowInsecureConnections: access.AllowInsecureConnections,
			ResourceType:             flexcdn.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			CertificateId:            xmaps.GetInt64(options.ProviderConfig, "certificateId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
