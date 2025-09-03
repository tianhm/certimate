package notify

import (
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/notify/notifiers"
	"github.com/certimate-go/certimate/pkg/core"
)

type notifierProviderOptions struct {
	Provider               domain.NotificationProviderType
	ProviderAccessConfig   map[string]any
	ProviderExtendedConfig map[string]any
}

func createNotifierProvider(options *notifierProviderOptions) (core.Notifier, error) {
	provider, err := notifiers.Registries.Get(options.Provider)
	if err != nil {
		return nil, err
	}

	return provider(&notifiers.ProviderFactoryOptions{
		ProviderAccessConfig:   options.ProviderAccessConfig,
		ProviderExtendedConfig: options.ProviderExtendedConfig,
	})
}
