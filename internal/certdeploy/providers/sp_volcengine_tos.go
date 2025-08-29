package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	volcenginetos "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/volcengine-tos"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeVolcEngineTOS, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForVolcEngine{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := volcenginetos.NewSSLDeployerProvider(&volcenginetos.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.SecretAccessKey,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
			Bucket:          xmaps.GetString(options.ProviderConfig, "bucket"),
			Domain:          xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
