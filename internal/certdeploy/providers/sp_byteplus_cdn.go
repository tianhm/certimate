package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	bytepluscdn "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/byteplus-cdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeBytePlusCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForBytePlus{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := bytepluscdn.NewSSLDeployerProvider(&bytepluscdn.SSLDeployerProviderConfig{
			AccessKey: access.AccessKey,
			SecretKey: access.SecretKey,
			Domain:    xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
