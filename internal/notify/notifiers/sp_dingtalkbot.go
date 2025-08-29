package notifiers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/notifier/providers/dingtalkbot"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.NotificationProviderTypeDingTalkBot, func(options *ProviderFactoryOptions) (core.Notifier, error) {
		access := domain.AccessConfigForDingTalkBot{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := dingtalkbot.NewNotifierProvider(&dingtalkbot.NotifierProviderConfig{
			WebhookUrl: access.WebhookUrl,
			Secret:     access.Secret,
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
