package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	aliyunclb "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aliyun-clb"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAliyunCLB, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAliyun{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := aliyunclb.NewSSLDeployerProvider(&aliyunclb.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			ResourceGroupId: access.ResourceGroupId,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
			ResourceType:    aliyunclb.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			LoadbalancerId:  xmaps.GetString(options.ProviderConfig, "loadbalancerId"),
			ListenerPort:    xmaps.GetOrDefaultInt32(options.ProviderConfig, "listenerPort", 443),
			Domain:          xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
