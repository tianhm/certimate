package providers

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	aliyuncasdeploy "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aliyun-cas-deploy"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAliyunCASDeploy, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAliyun{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := aliyuncasdeploy.NewSSLDeployerProvider(&aliyuncasdeploy.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			ResourceGroupId: access.ResourceGroupId,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
			ResourceIds:     lo.Filter(strings.Split(xmaps.GetString(options.ProviderConfig, "resourceIds"), ";"), func(s string, _ int) bool { return s != "" }),
			ContactIds:      lo.Filter(strings.Split(xmaps.GetString(options.ProviderConfig, "contactIds"), ";"), func(s string, _ int) bool { return s != "" }),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
