package providers

import (
	"fmt"
	"net/http"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	webhook "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/webhook"
	xhttp "github.com/certimate-go/certimate/pkg/utils/http"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	if err := Registries.Register(domain.DeploymentProviderTypeWebhook, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		access := domain.AccessConfigForWebhook{}
		if err := xmaps.Populate(options.AccessConfig, &access); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		mergedHeaders := make(map[string]string)
		if defaultHeadersString := access.HeadersString; defaultHeadersString != "" {
			h, err := xhttp.ParseHeaders(defaultHeadersString)
			if err != nil {
				return nil, fmt.Errorf("failed to parse webhook headers: %w", err)
			}
			for key := range h {
				mergedHeaders[http.CanonicalHeaderKey(key)] = h.Get(key)
			}
		}
		if extendedHeadersString := xmaps.GetString(options.ProviderConfig, "headers"); extendedHeadersString != "" {
			h, err := xhttp.ParseHeaders(extendedHeadersString)
			if err != nil {
				return nil, fmt.Errorf("failed to parse webhook headers: %w", err)
			}
			for key := range h {
				mergedHeaders[http.CanonicalHeaderKey(key)] = h.Get(key)
			}
		}

		provider, err := webhook.NewSSLDeployerProvider(&webhook.SSLDeployerProviderConfig{
			WebhookUrl:               access.Url,
			WebhookData:              xmaps.GetOrDefaultString(options.ProviderConfig, "webhookData", access.DataString),
			Method:                   access.Method,
			Headers:                  mergedHeaders,
			AllowInsecureConnections: access.AllowInsecureConnections,
		})
		return provider, err
	}); err != nil {
		panic(err)
	}
}
