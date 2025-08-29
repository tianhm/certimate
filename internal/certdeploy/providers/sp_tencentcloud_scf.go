package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	tencentcloudscf "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-scf"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeTencentCloudSCF, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForTencentCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := tencentcloudscf.NewSSLDeployerProvider(&tencentcloudscf.SSLDeployerProviderConfig{
			SecretId:  access.SecretId,
			SecretKey: access.SecretKey,
			Endpoint:  xmaps.GetString(options.ProviderConfig, "endpoint"),
			Region:    xmaps.GetString(options.ProviderConfig, "region"),
			Domain:    xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
