package deployers

import (
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/local"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeLocal, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		provider, err := local.NewSSLDeployerProvider(&local.SSLDeployerProviderConfig{
			ShellEnv:                 local.ShellEnvType(xmaps.GetString(options.ProviderConfig, "shellEnv")),
			PreCommand:               xmaps.GetString(options.ProviderConfig, "preCommand"),
			PostCommand:              xmaps.GetString(options.ProviderConfig, "postCommand"),
			OutputFormat:             local.OutputFormatType(xmaps.GetOrDefaultString(options.ProviderConfig, "format", string(local.OUTPUT_FORMAT_PEM))),
			OutputCertPath:           xmaps.GetString(options.ProviderConfig, "certPath"),
			OutputServerCertPath:     xmaps.GetString(options.ProviderConfig, "certPathForServerOnly"),
			OutputIntermediaCertPath: xmaps.GetString(options.ProviderConfig, "certPathForIntermediaOnly"),
			OutputKeyPath:            xmaps.GetString(options.ProviderConfig, "keyPath"),
			PfxPassword:              xmaps.GetString(options.ProviderConfig, "pfxPassword"),
			JksAlias:                 xmaps.GetString(options.ProviderConfig, "jksAlias"),
			JksKeypass:               xmaps.GetString(options.ProviderConfig, "jksKeypass"),
			JksStorepass:             xmaps.GetString(options.ProviderConfig, "jksStorepass"),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
