package providers

import (
	"fmt"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core/ssl-applicator/acme-dns01/providers/acmehttpreq"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := ACMEDns01Registries.Register(domain.ACMEDns01ProviderTypeACMEHttpReq, func(options *ProviderFactoryOptions) (challenge.Provider, error) {
		access := domain.AccessConfigForACMEHttpReq{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := acmehttpreq.NewChallengeProvider(&acmehttpreq.ChallengeProviderConfig{
			Endpoint:              access.Endpoint,
			Mode:                  access.Mode,
			Username:              access.Username,
			Password:              access.Password,
			DnsPropagationTimeout: options.DnsPropagationTimeout,
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
