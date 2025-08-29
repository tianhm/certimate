package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	qiniupili "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/qiniu-pili"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeQiniuPili, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForQiniu{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := qiniupili.NewSSLDeployerProvider(&qiniupili.SSLDeployerProviderConfig{
			AccessKey: access.AccessKey,
			SecretKey: access.SecretKey,
			Hub:       xmaps.GetString(options.ProviderConfig, "hub"),
			Domain:    xmaps.GetString(options.ProviderConfig, "domain"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
