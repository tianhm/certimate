package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/safeline"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeSafeLine, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForSafeLine{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := safeline.NewSSLDeployerProvider(&safeline.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiToken:                 access.ApiToken,
			AllowInsecureConnections: access.AllowInsecureConnections,
			ResourceType:             safeline.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			CertificateId:            xmaps.GetInt32(options.ProviderConfig, "certificateId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
