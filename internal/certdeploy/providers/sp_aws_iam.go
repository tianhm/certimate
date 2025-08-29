package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	awsiam "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aws-iam"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAWSIAM, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAWS{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := awsiam.NewSSLDeployerProvider(&awsiam.SSLDeployerProviderConfig{
			AccessKeyId:     access.AccessKeyId,
			SecretAccessKey: access.SecretAccessKey,
			Region:          xmaps.GetString(options.ProviderConfig, "region"),
			CertificatePath: xmaps.GetOrDefaultString(options.ProviderConfig, "certificatePath", "/"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
