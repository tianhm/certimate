package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	rainyunrcdn "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/rainyun-rcdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeRainYunRCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForRainYun{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := rainyunrcdn.NewSSLDeployerProvider(&rainyunrcdn.SSLDeployerProviderConfig{
			ApiKey:     access.ApiKey,
			InstanceId: xmaps.GetInt32(options.ProviderConfig, "instanceId"),
			Domain:     xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
