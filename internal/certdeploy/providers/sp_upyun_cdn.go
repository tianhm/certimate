package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	upyuncdn "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/upyun-cdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeUpyunCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForUpyun{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := upyuncdn.NewSSLDeployerProvider(&upyuncdn.SSLDeployerProviderConfig{
			Username: access.Username,
			Password: access.Password,
			Domain:   xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
