package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	volcenginecertcenter "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/volcengine-certcenter"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeVolcEngineCertCenter, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForVolcEngine{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := volcenginecertcenter.NewSSLDeployerProvider(&volcenginecertcenter.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.SecretAccessKey,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
