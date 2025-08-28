package notify

import (
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/notify/providers"
	"github.com/certimate-go/certimate/pkg/core"
)

type notifierProviderOptions struct {
	Provider              domain.NotificationProviderType
	ProviderAccessConfig  map[string]any
	ProviderServiceConfig map[string]any
}

func createNotifierProvider(options *notifierProviderOptions) (core.Notifier, error) {
	provider, err := providers.Registries.Get(options.Provider)
	if err != nil {
		return nil, err
	}

	return provider(&providers.ProviderFactoryOptions{
		AccessConfig:   options.ProviderAccessConfig,
		ProviderConfig: options.ProviderServiceConfig,
	})
}
