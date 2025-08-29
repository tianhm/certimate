package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	uclouducdn "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ucloud-ucdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeUCloudUCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForUCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := uclouducdn.NewSSLDeployerProvider(&uclouducdn.SSLDeployerProviderConfig{
			PrivateKey: access.PrivateKey,
			PublicKey:  access.PublicKey,
			ProjectId:  access.ProjectId,
			DomainId:   xmaps.GetString(options.ProviderConfig, "domainId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
