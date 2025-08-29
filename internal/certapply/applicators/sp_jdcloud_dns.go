package applicators

import (
	"fmt"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core/ssl-applicator/acme-dns01/providers/jdcloud"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := ACMEDns01Registries.Register(domain.ACMEDns01ProviderTypeJDCloudDNS, func(options *ProviderFactoryOptions) (challenge.Provider, error) {
		access := domain.AccessConfigForJDCloud{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := jdcloud.NewChallengeProvider(&jdcloud.ChallengeProviderConfig{
			AccessKeyId:           access.AccessKeyId,
			AccessKeySecret:       access.AccessKeySecret,
			RegionId:              xmaps.GetString(options.ProviderConfig, "regionId"),
			DnsPropagationTimeout: options.DnsPropagationTimeout,
			DnsTTL:                options.DnsTTL,
		})
		return provider, err
	}); err != nil {
		panic(err)
	}

	if err := ACMEDns01Registries.RegisterAlias(domain.ACMEDns01ProviderTypeJDCloud, domain.ACMEDns01ProviderTypeJDCloudDNS); err != nil {
		panic(err)
	}
}
