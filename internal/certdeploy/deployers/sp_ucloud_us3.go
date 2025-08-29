package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	ucloudus3 "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ucloud-us3"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeUCloudUS3, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForUCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := ucloudus3.NewSSLDeployerProvider(&ucloudus3.SSLDeployerProviderConfig{
			PrivateKey: access.PrivateKey,
			PublicKey:  access.PublicKey,
			ProjectId:  access.ProjectId,
			Region:     xmaps.GetString(options.ProviderConfig, "region"),
			Bucket:     xmaps.GetString(options.ProviderConfig, "bucket"),
			Domain:     xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
