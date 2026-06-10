package cloudflare

import (
	"fmt"
	"strings"
)

type sdkResponse interface {
	GetErrors() error
	GetSuccess() bool
}

type sdkResponseBase struct {
	Errors   APIErrors    `json:"errors,omitempty"`
	Messages []APIMessage `json:"messages,omitempty"`
	Success  bool         `json:"success,omitempty"`
}

type APIMessage struct {
	Code             int             `json:"code"`
	Message          string          `json:"message"`
	DocumentationURL string          `json:"documentation_url"`
	ErrorChain       []APIErrorChain `json:"error_chain"`
	Source           *APISource      `json:"source"`
}

type APIErrors []APIMessage

type APIErrorChain struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type APISource struct {
	Pointer string `json:"pointer"`
}

func (e APIErrors) Error() string {
	builder := &strings.Builder{}

	for _, item := range e {
		fmt.Fprintf(builder, "%d: %s", item.Code, item.Message)

		for _, link := range item.ErrorChain {
			fmt.Fprintf(builder, "; %d: %s", link.Code, link.Message)
		}
	}

	return builder.String()
}

func (r *sdkResponseBase) GetErrors() error {
	if len(r.Errors) > 0 {
		return r.Errors
	}
	return nil
}

func (r *sdkResponseBase) GetSuccess() bool {
	return r.Success
}

var _ sdkResponse = (*sdkResponseBase)(nil)

type GeoRestriction struct {
	Label string `json:"label"`
}

type CustomCertificate struct {
	ID                 string           `json:"id"`
	ZoneID             string           `json:"zone_id"`
	BundleMethod       string           `json:"bundle_method"`
	CustomCsrID        string           `json:"custom_csr_id"`
	GeoRestrictions    []GeoRestriction `json:"geo_restrictions"`
	Hosts              []string         `json:"hosts"`
	Issuer             string           `json:"issuer"`
	PolicyRestrictions string           `json:"policy_restrictions"`
	Priority           float64          `json:"priority"`
	Signature          string           `json:"signature"`
	Status             string           `json:"status"`
	ExpiresOn          string           `json:"expires_on"`
	UploadedOn         string           `json:"uploaded_on"`
	ModifiedOn         string           `json:"modified_on"`
}
