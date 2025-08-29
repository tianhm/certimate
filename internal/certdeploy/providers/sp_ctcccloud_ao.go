package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	ctcccloudao "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ctcccloud-ao"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeCTCCCloudAO, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForCTCCCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := ctcccloudao.NewSSLDeployerProvider(&ctcccloudao.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			SecretAccessKey: access.SecretAccessKey,
			Domain:          xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
