package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/notifier/providers/email"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.NotificationProviderTypeEmail, func(options *ProviderFactoryOptions) (core.Notifier, error) {
		access := domain.AccessConfigForEmail{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := email.NewNotifierProvider(&email.NotifierProviderConfig{
			SmtpHost:        access.SmtpHost,
			SmtpPort:        access.SmtpPort,
			SmtpTls:         access.SmtpTls,
			Username:        access.Username,
			Password:        access.Password,
			SenderAddress:   access.SenderAddress,
			SenderName:      access.SenderName,
			ReceiverAddress: xmaps.GetOrDefaultString(options.ProviderConfig, "receiverAddress", access.ReceiverAddress),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
