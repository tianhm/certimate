package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/proxmoxve"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeProxmoxVE, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForProxmoxVE{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := proxmoxve.NewSSLDeployerProvider(&proxmoxve.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiToken:                 access.ApiToken,
			ApiTokenSecret:           access.ApiTokenSecret,
			AllowInsecureConnections: access.AllowInsecureConnections,
			NodeName:                 xmaps.GetString(options.ProviderConfig, "nodeName"),
			AutoRestart:              xmaps.GetBool(options.ProviderConfig, "autoRestart"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
