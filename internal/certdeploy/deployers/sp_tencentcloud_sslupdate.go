package deployers

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	tencentcloudsslupdate "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-ssl-update"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeTencentCloudSSLUpdate, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForTencentCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := tencentcloudsslupdate.NewSSLDeployerProvider(&tencentcloudsslupdate.SSLDeployerProviderConfig{
			SecretId:        access.SecretId,
			SecretKey:       access.SecretKey,
			Endpoint:        xmaps.GetString(options.ProviderConfig, "endpoint"),
			CertificateId:   xmaps.GetString(options.ProviderConfig, "certificateId"),
			IsReplaced:      xmaps.GetBool(options.ProviderConfig, "isReplaced"),
			ResourceTypes:   lo.Filter(strings.Split(xmaps.GetString(options.ProviderConfig, "resourceTypes"), ";"), func(s string, _ int) bool { return s != "" }),
			ResourceRegions: lo.Filter(strings.Split(xmaps.GetString(options.ProviderConfig, "resourceRegions"), ";"), func(s string, _ int) bool { return s != "" }),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
