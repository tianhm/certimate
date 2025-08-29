package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	aliyunwaf "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aliyun-waf"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAliyunWAF, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAliyun{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := aliyunwaf.NewSSLDeployerProvider(&aliyunwaf.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			ResourceGroupId: access.ResourceGroupId,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
			ServiceVersion:  xmaps.GetOrDefaultString(options.ProviderConfig, "serviceVersion", "3.0"),
			InstanceId:      xmaps.GetString(options.ProviderConfig, "instanceId"),
			Domain:          xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
