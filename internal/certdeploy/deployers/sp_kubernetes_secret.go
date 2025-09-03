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
		credentials := domain.AccessConfigForKubernetes{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := k8ssecret.NewSSLDeployerProvider(&k8ssecret.SSLDeployerProviderConfig{
			KubeConfig:          credentials.KubeConfig,
			Namespace:           xmaps.GetOrDefaultString(options.ProviderExtendedConfig, "namespace", "default"),
			SecretName:          xmaps.GetString(options.ProviderExtendedConfig, "secretName"),
			SecretType:          xmaps.GetOrDefaultString(options.ProviderExtendedConfig, "secretType", "kubernetes.io/tls"),
			SecretDataKeyForCrt: xmaps.GetOrDefaultString(options.ProviderExtendedConfig, "secretDataKeyForCrt", "tls.crt"),
			SecretDataKeyForKey: xmaps.GetOrDefaultString(options.ProviderExtendedConfig, "secretDataKeyForKey", "tls.key"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
