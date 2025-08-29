package applicators

import (
	"fmt"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/certimate-go/certimate/internal/domain"
	teo "github.com/certimate-go/certimate/pkg/core/ssl-applicator/acme-dns01/providers/tencentcloud-eo"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := ACMEDns01Registries.Register(domain.ACMEDns01ProviderTypeTencentCloudEO, func(options *ProviderFactoryOptions) (challenge.Provider, error) {
		access := domain.AccessConfigForTencentCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := teo.NewChallengeProvider(&teo.ChallengeProviderConfig{
			SecretId:              access.SecretId,
			SecretKey:             access.SecretKey,
			ZoneId:                xmaps.GetString(options.ProviderConfig, "zoneId"),
			DnsPropagationTimeout: options.DnsPropagationTimeout,
			DnsTTL:                options.DnsTTL,
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
