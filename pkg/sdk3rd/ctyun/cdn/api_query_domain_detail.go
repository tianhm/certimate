package cdn

import (
	"context"
	"net/http"
)

type QueryDomainDetailRequest struct {
	Domain        *string `json:"domain,omitempty"`
	ProductCode   *string `json:"product_code,omitempty"`
	FunctionNames *string `json:"function_names,omitempty"`
}

type QueryDomainDetailResponse struct {
	apiResponseBase

	ReturnObj *DomainDetail `json:"returnObj,omitempty"`
}

func (c *Client) QueryDomainDetail(req *QueryDomainDetailRequest) (*QueryDomainDetailResponse, error) {
	return c.QueryDomainDetailWithContext(context.Background(), req)
}

func (c *Client) QueryDomainDetailWithContext(ctx context.Context, req *QueryDomainDetailRequest) (*QueryDomainDetailResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/v1/domain/query-domain-detail")
	if err != nil {
		return nil, err
	} else {
		if req.Domain != nil {
			httpreq.SetQueryParam("domain", *req.Domain)
		}
		if req.ProductCode != nil {
			httpreq.SetQueryParam("product_code", *req.ProductCode)
		}
		if req.FunctionNames != nil {
			httpreq.SetQueryParam("function_names", *req.FunctionNames)
		}

		httpreq.SetContext(ctx)
	}

	result := &QueryDomainDetailResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
