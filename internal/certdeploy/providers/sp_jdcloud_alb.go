package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	jdcloudalb "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/jdcloud-alb"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeJDCloudALB, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForJDCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := jdcloudalb.NewSSLDeployerProvider(&jdcloudalb.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			RegionId:        xmaps.GetString(options.ProviderConfig, "regionId"),
			ResourceType:    jdcloudalb.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			LoadbalancerId:  xmaps.GetString(options.ProviderConfig, "loadbalancerId"),
			ListenerId:      xmaps.GetString(options.ProviderConfig, "listenerId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
