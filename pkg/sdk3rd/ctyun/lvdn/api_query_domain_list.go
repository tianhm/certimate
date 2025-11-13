package lvdn

import (
	"context"
	"net/http"
	"strconv"
)

type QueryDomainListRequest struct {
	Page        *int32  `json:"page,omitempty"`
	PageSize    *int32  `json:"page_size,omitempty"`
	Domain      *string `json:"domain,omitempty"`
	ProductCode *string `json:"product_code,omitempty"`
	Status      *int32  `json:"status,omitempty"`
	AreaScope   *int32  `json:"area_scope,omitempty"`
}

type QueryDomainListResponse struct {
	apiResponseBase

	ReturnObj *struct {
		Results   []*DomainRecord `json:"result,omitempty"`
		Page      int32           `json:"page,omitempty"`
		PageSize  int32           `json:"page_size,omitempty"`
		PageCount int32           `json:"page_count,omitempty"`
		Total     int32           `json:"total,omitempty"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) QueryDomainList(req *QueryDomainListRequest) (*QueryDomainListResponse, error) {
	return c.QueryDomainListWithContext(context.Background(), req)
}

func (c *Client) QueryDomainListWithContext(ctx context.Context, req *QueryDomainListRequest) (*QueryDomainListResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/v1/domain/query-domain-list")
	if err != nil {
		return nil, err
	} else {
		if req.Page != nil {
			httpreq.SetQueryParam("page", strconv.Itoa(int(*req.Page)))
		}
		if req.PageSize != nil {
			httpreq.SetQueryParam("page_size", strconv.Itoa(int(*req.PageSize)))
		}
		if req.Domain != nil {
			httpreq.SetQueryParam("domain", *req.Domain)
		}
		if req.ProductCode != nil {
			httpreq.SetQueryParam("product_code", *req.ProductCode)
		}
		if req.Status != nil {
			httpreq.SetQueryParam("status", strconv.Itoa(int(*req.Status)))
		}
		if req.AreaScope != nil {
			httpreq.SetQueryParam("area_scope", strconv.Itoa(int(*req.AreaScope)))
		}

		httpreq.SetContext(ctx)
	}

	result := &QueryDomainListResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
