package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	aliyunesa "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aliyun-esa"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAliyunESA, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAliyun{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := aliyunesa.NewSSLDeployerProvider(&aliyunesa.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
			SiteId:          xmaps.GetInt64(options.ProviderConfig, "siteId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
