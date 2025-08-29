package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	aliyunnlb "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aliyun-nlb"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAliyunNLB, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAliyun{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := aliyunnlb.NewSSLDeployerProvider(&aliyunnlb.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			ResourceGroupId: access.ResourceGroupId,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
			ResourceType:    aliyunnlb.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			LoadbalancerId:  xmaps.GetString(options.ProviderConfig, "loadbalancerId"),
			ListenerId:      xmaps.GetString(options.ProviderConfig, "listenerId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
