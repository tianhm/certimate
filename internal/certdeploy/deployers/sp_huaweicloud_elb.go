package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	huaweicloudelb "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/huaweicloud-elb"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeHuaweiCloudELB, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForHuaweiCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := huaweicloudelb.NewSSLDeployerProvider(&huaweicloudelb.SSLDeployerProviderConfig{
			AccessKeyId:         access.AccessKeyId,
			SecretAccessKey:     access.SecretAccessKey,
			EnterpriseProjectId: access.EnterpriseProjectId,
			Region:              xmaps.GetString(options.ProviderConfig, "region"),
			ResourceType:        huaweicloudelb.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			CertificateId:       xmaps.GetString(options.ProviderConfig, "certificateId"),
			LoadbalancerId:      xmaps.GetString(options.ProviderConfig, "loadbalancerId"),
			ListenerId:          xmaps.GetString(options.ProviderConfig, "listenerId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
