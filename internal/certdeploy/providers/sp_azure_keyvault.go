package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	azurekeyvault "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/azure-keyvault"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeAzureKeyVault, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForAzure{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := azurekeyvault.NewSSLDeployerProvider(&azurekeyvault.SSLDeployerProviderConfig{
			TenantId:        access.TenantId,
			ClientId:        access.ClientId,
			ClientSecret:    access.ClientSecret,
			CloudName:       access.CloudName,
			KeyVaultName:    xmaps.GetString(options.ProviderConfig, "keyvaultName"),
			CertificateName: xmaps.GetString(options.ProviderConfig, "certificateName"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
