package applicant

import (
	"github.com/go-acme/lego/v4/challenge"

	"github.com/certimate-go/certimate/internal/certapply/providers"
	"github.com/certimate-go/certimate/internal/domain"
)

type applicantProviderOptions struct {
	Domains                 []string
	ContactEmail            string
	Provider                domain.ACMEDns01ProviderType
	ProviderAccessConfig    map[string]any
	ProviderServiceConfig   map[string]any
	CAProvider              domain.CAProviderType
	CAProviderAccessId      string
	CAProviderAccessConfig  map[string]any
	CAProviderServiceConfig map[string]any
	KeyAlgorithm            string
	Nameservers             []string
	DnsPropagationWait      int32
	DnsPropagationTimeout   int32
	DnsTTL                  int32
	ACMEProfile             string
	DisableFollowCNAME      bool
	ARIReplaceAcct          string
	ARIReplaceCert          string
}

func createApplicantProvider(options *applicantProviderOptions) (challenge.Provider, error) {
	provider, err := providers.ACMEDns01Registries.Get(options.Provider)
	if err != nil {
		return nil, err
	}

	return provider(&providers.ProviderFactoryOptions{
		AccessConfig:          options.ProviderAccessConfig,
		ProviderConfig:        options.ProviderServiceConfig,
		DnsPropagationWait:    options.DnsPropagationWait,
		DnsPropagationTimeout: options.DnsPropagationTimeout,
		DnsTTL:                options.DnsTTL,
	})
}
