package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/notifier/providers/larkbot"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.NotificationProviderTypeLarkBot, func(options *ProviderFactoryOptions) (core.Notifier, error) {
		access := domain.AccessConfigForLarkBot{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := larkbot.NewNotifierProvider(&larkbot.NotifierProviderConfig{
			WebhookUrl: access.WebhookUrl,
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
