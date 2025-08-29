package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ssh"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeSSH, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForSSH{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		jumpServers := make([]ssh.JumpServerConfig, len(access.JumpServers))
		for i, jumpServer := range access.JumpServers {
			jumpServers[i] = ssh.JumpServerConfig{
				SshHost:          jumpServer.Host,
				SshPort:          jumpServer.Port,
				SshAuthMethod:    jumpServer.AuthMethod,
				SshUsername:      jumpServer.Username,
				SshPassword:      jumpServer.Password,
				SshKey:           jumpServer.Key,
				SshKeyPassphrase: jumpServer.KeyPassphrase,
			}
		}

		provider, err := ssh.NewSSLDeployerProvider(&ssh.SSLDeployerProviderConfig{
			SshHost:                  access.Host,
			SshPort:                  access.Port,
			SshAuthMethod:            access.AuthMethod,
			SshUsername:              access.Username,
			SshPassword:              access.Password,
			SshKey:                   access.Key,
			SshKeyPassphrase:         access.KeyPassphrase,
			JumpServers:              jumpServers,
			UseSCP:                   xmaps.GetBool(options.ProviderConfig, "useSCP"),
			PreCommand:               xmaps.GetString(options.ProviderConfig, "preCommand"),
			PostCommand:              xmaps.GetString(options.ProviderConfig, "postCommand"),
			OutputFormat:             ssh.OutputFormatType(xmaps.GetOrDefaultString(options.ProviderConfig, "format", string(ssh.OUTPUT_FORMAT_PEM))),
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
