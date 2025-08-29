package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	aliyunga "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aliyun-ga"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAliyunGA, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAliyun{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := aliyunga.NewSSLDeployerProvider(&aliyunga.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			ResourceGroupId: access.ResourceGroupId,
			ResourceType:    aliyunga.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			AcceleratorId:   xmaps.GetString(options.ProviderConfig, "acceleratorId"),
			ListenerId:      xmaps.GetString(options.ProviderConfig, "listenerId"),
			Domain:          xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
