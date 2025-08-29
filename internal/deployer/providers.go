package deployer

import (
	"github.com/certimate-go/certimate/internal/certdeploy/providers"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
)

type deployerProviderOptions struct {
	Provider              domain.DeploymentProviderType
	ProviderAccessConfig  map[string]any
	ProviderServiceConfig map[string]any
}

func createSSLDeployerProvider(options *deployerProviderOptions) (core.SSLDeployer, error) {
	provider, err := providers.Registries.Get(options.Provider)
	if err != nil {
		return nil, err
	}

	return provider(&providers.ProviderFactoryOptions{
		AccessConfig:   options.ProviderAccessConfig,
		ProviderConfig: options.ProviderServiceConfig,
	})
}
