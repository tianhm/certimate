package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	baotawafconsole "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/baotawaf-console"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeBaotaWAFConsole, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForBaotaWAF{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := baotawafconsole.NewSSLDeployerProvider(&baotawafconsole.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiKey:                   access.ApiKey,
			AllowInsecureConnections: access.AllowInsecureConnections,
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
