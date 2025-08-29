package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	pDogeCDN "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/dogecloud-cdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeDogeCloudCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForDogeCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := pDogeCDN.NewSSLDeployerProvider(&pDogeCDN.SSLDeployerProviderConfig{
			AccessKey: access.AccessKey,
			SecretKey: access.SecretKey,
			Domain:    xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
