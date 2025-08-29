package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	tencentcloudclb "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-clb"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeTencentCloudCLB, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForTencentCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := tencentcloudclb.NewSSLDeployerProvider(&tencentcloudclb.SSLDeployerProviderConfig{
			SecretId:       access.SecretId,
			SecretKey:      access.SecretKey,
			Endpoint:       xmaps.GetString(options.ProviderConfig, "endpoint"),
			Region:         xmaps.GetString(options.ProviderConfig, "region"),
			ResourceType:   tencentcloudclb.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			LoadbalancerId: xmaps.GetString(options.ProviderConfig, "loadbalancerId"),
			ListenerId:     xmaps.GetString(options.ProviderConfig, "listenerId"),
			Domain:         xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
