package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/goedge"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeGoEdge, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForGoEdge{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := goedge.NewSSLDeployerProvider(&goedge.SSLDeployerProviderConfig{
			ServerUrl:                access.ServerUrl,
			ApiRole:                  access.ApiRole,
			AccessKeyId:              access.AccessKeyId,
			AccessKey:                access.AccessKey,
			AllowInsecureConnections: access.AllowInsecureConnections,
			ResourceType:             goedge.ResourceType(xmaps.GetString(options.ProviderConfig, "resourceType")),
			CertificateId:            xmaps.GetInt64(options.ProviderConfig, "certificateId"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
