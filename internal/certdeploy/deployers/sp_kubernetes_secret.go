package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	k8ssecret "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/k8s-secret"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeKubernetesSecret, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForKubernetes{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := k8ssecret.NewSSLDeployerProvider(&k8ssecret.SSLDeployerProviderConfig{
			KubeConfig:          access.KubeConfig,
			Namespace:           xmaps.GetOrDefaultString(options.ProviderConfig, "namespace", "default"),
			SecretName:          xmaps.GetString(options.ProviderConfig, "secretName"),
			SecretType:          xmaps.GetOrDefaultString(options.ProviderConfig, "secretType", "kubernetes.io/tls"),
			SecretDataKeyForCrt: xmaps.GetOrDefaultString(options.ProviderConfig, "secretDataKeyForCrt", "tls.crt"),
			SecretDataKeyForKey: xmaps.GetOrDefaultString(options.ProviderConfig, "secretDataKeyForKey", "tls.key"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
