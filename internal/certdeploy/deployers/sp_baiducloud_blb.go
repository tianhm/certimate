package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	baiducloudblb "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/baiducloud-blb"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeBaiduCloudBLB, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForBaiduCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := baiducloudblb.NewSSLDeployerProvider(&baiducloudblb.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			SecretAccessKey: access.SecretAccessKey,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
			ResourceType:    baiducloudblb.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			LoadbalancerId:  xmaps.GetString(options.ProviderConfig, "loadbalancerId"),
			ListenerPort:    xmaps.GetInt32(options.ProviderConfig, "listenerPort"),
			Domain:          xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
