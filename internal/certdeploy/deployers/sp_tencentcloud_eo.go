package deployers

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	tencentcloudeo "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-eo"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeTencentCloudEO, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForTencentCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := tencentcloudeo.NewSSLDeployerProvider(&tencentcloudeo.SSLDeployerProviderConfig{
			SecretId:  access.SecretId,
			SecretKey: access.SecretKey,
			Endpoint:  xmaps.GetString(options.ProviderConfig, "endpoint"),
			ZoneId:    xmaps.GetString(options.ProviderConfig, "zoneId"),
			Domains:   lo.Filter(strings.Split(xmaps.GetString(options.ProviderConfig, "domains"), ";"), func(s string, _ int) bool { return s != "" }),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
