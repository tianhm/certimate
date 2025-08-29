package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	opconsole "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/1panel-console"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderType1PanelConsole, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigFor1Panel{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := opconsole.NewSSLDeployerProvider(&opconsole.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiVersion:               access.ApiVersion,
			ApiKey:                   access.ApiKey,
			AllowInsecureConnections: access.AllowInsecureConnections,
			AutoRestart:              xmaps.GetBool(options.ProviderConfig, "autoRestart"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
