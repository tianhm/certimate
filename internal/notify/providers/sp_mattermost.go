package providers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/notifier/providers/mattermost"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.NotificationProviderTypeMattermost, func(options *ProviderFactoryOptions) (core.Notifier, error) {
		access := domain.AccessConfigForMattermost{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := mattermost.NewNotifierProvider(&mattermost.NotifierProviderConfig{
			ServerUrl: access.ServerUrl,
			Username:  access.Username,
			Password:  access.Password,
			ChannelId: xmaps.GetOrDefaultString(options.ProviderConfig, "channelId", access.ChannelId),
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
