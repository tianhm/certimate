package applicators

import (
	"fmt"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/certimate-go/certimate/internal/domain"
	xinnet "github.com/certimate-go/certimate/pkg/core/ssl-applicator/acme-dns01/providers/xinnet"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	ACMEDns01Registries.MustRegister(domain.ACMEDns01ProviderTypeXinnet, func(options *ProviderFactoryOptions) (challenge.Provider, error) {
		credentials := domain.AccessConfigForXinnet{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := xinnet.NewChallengeProvider(&xinnet.ChallengeProviderConfig{
			AgentId:               credentials.AgentId,
			ApiPassword:           credentials.ApiPassword,
			DnsPropagationTimeout: options.DnsPropagationTimeout,
			DnsTTL:                options.DnsTTL,
		})
		return provider, err
	})
}
