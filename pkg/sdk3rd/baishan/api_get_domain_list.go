package baishan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetDomainListRequest struct {
	PageNumber   *int32  `json:"page_number,omitempty"`
	PageSize     *int32  `json:"page_size,omitempty"`
	DomainStatus *string `json:"domain_status,omitempty"`
}

type GetDomainListResponse struct {
	apiResponseBase

	Data []*struct {
		List        []*DomainRecord `json:"list"`
		PageNumber  json.Number     `json:"page_number"`
		PageSize    json.Number     `json:"page_size"`
		TotalNumber json.Number     `json:"total_number"`
	} `json:"data,omitempty"`
}

func (c *Client) GetDomainList(req *GetDomainListRequest) (*GetDomainListResponse, error) {
	return c.GetDomainListWithContext(context.Background(), req)
}

func (c *Client) GetDomainListWithContext(ctx context.Context, req *GetDomainListRequest) (*GetDomainListResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/v2/domain/list")
	if err != nil {
		return nil, err
	} else {
		if req.PageNumber != nil {
			httpreq.SetQueryParam("page_number", fmt.Sprintf("%d", *req.PageNumber))
		}
		if req.PageSize != nil {
			httpreq.SetQueryParam("page_number", fmt.Sprintf("%d", *req.PageSize))
		}
		if req.DomainStatus != nil {
			httpreq.SetQueryParam("domain_status", *req.DomainStatus)
		}

		httpreq.SetContext(ctx)
	}

	result := &GetDomainListResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
