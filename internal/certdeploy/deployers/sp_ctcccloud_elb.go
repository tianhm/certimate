package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	ctcccloudelb "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ctcccloud-elb"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeCTCCCloudELB, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForCTCCCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := ctcccloudelb.NewSSLDeployerProvider(&ctcccloudelb.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			SecretAccessKey: access.SecretAccessKey,
			RegionId:        xmaps.GetString(options.ProviderConfig, "regionId"),
			ResourceType:    ctcccloudelb.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			LoadbalancerId:  xmaps.GetString(options.ProviderConfig, "loadbalancerId"),
			ListenerId:      xmaps.GetString(options.ProviderConfig, "listenerId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
