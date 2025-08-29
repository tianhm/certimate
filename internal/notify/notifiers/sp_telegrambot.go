package notifiers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/notifier/providers/telegrambot"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.NotificationProviderTypeTelegramBot, func(options *ProviderFactoryOptions) (core.Notifier, error) {
		access := domain.AccessConfigForTelegramBot{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := telegrambot.NewNotifierProvider(&telegrambot.NotifierProviderConfig{
			BotToken: access.BotToken,
			ChatId:   xmaps.GetOrDefaultInt64(options.ProviderConfig, "chatId", access.ChatId),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
