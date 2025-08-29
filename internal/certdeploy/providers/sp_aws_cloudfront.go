package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	awscloudfront "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aws-cloudfront"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAWSCloudFront, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAWS{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := awscloudfront.NewSSLDeployerProvider(&awscloudfront.SSLDeployerProviderConfig{
			AccessKeyId:       access.AccessKeyId,
			SecretAccessKey:   access.SecretAccessKey,
			Region:            xmaps.GetString(options.ProviderConfig, "region"),
			DistributionId:    xmaps.GetString(options.ProviderConfig, "distributionId"),
			CertificateSource: xmaps.GetOrDefaultString(options.ProviderConfig, "certificateSource", "ACM"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
