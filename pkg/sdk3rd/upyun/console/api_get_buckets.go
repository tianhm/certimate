package console

import (
	"context"
	"fmt"
	"net/http"
)

type GetBucketsRequest struct {
	BucketName    string `json:"status"`
	BusinessType  string `json:"business_type"`
	Type          string `json:"type"`
	Status        string `json:"state"`
	Tag           string `json:"tag"`
	IsSecurityCDN bool   `json:"security_cdn"`
	WithDomains   bool   `json:"with_domains"`
	Page          int32  `json:"page"`
	PerPage       int32  `json:"perPage"`
}

type GetBucketsResponse struct {
	apiResponseBase
	Data *struct {
		apiResponseBaseData
		Buckets []*BucketInfo `json:"buckets"`
		Pager   BucketPager   `json:"pager"`
	} `json:"data,omitempty"`
}

type BucketInfo struct {
	BucketName    string          `json:"bucket_name"`
	BusinessType  string          `json:"business_type"`
	Type          string          `json:"type"`
	Status        string          `json:"status"`
	Tag           string          `json:"tag"`
	IsFusionCDN   bool            `json:"fusion_cdn"`
	IsSecurityCDN bool            `json:"security_cdn"`
	Domains       []*BucketDomain `json:"domains"`
	Visible       bool            `json:"visible"`
	CreatedAt     string          `json:"created_at"`
}

type BucketDomain struct {
	Domain string `json:"domain"`
	Status string `json:"status"`
}

type BucketPager struct {
	Page  int32 `json:"page"`
	Pages int64 `json:"pages"`
	Total int64 `json:"total"`
}

func (c *Client) GetBuckets(req *GetBucketsRequest) (*GetBucketsResponse, error) {
	return c.GetBucketsWithContext(context.Background(), req)
}

func (c *Client) GetBucketsWithContext(ctx context.Context, req *GetBucketsRequest) (*GetBucketsResponse, error) {
	if err := c.ensureCookieExists(); err != nil {
		return nil, err
	}

	httpreq, err := c.newRequest(http.MethodGet, "/api/v2/buckets")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetQueryParam("bucket_name", req.BucketName)
		httpreq.SetQueryParam("business_type", req.BusinessType)
		httpreq.SetQueryParam("type", req.Type)
		httpreq.SetQueryParam("state", req.Status)
		httpreq.SetQueryParam("tag", req.Tag)
		httpreq.SetQueryParam("security_cdn", fmt.Sprintf("%v", req.IsSecurityCDN))
		httpreq.SetQueryParam("with_domains", fmt.Sprintf("%v", req.WithDomains))
		httpreq.SetQueryParam("page", fmt.Sprintf("%d", req.Page))
		httpreq.SetQueryParam("perPage", fmt.Sprintf("%d", req.PerPage))
		httpreq.SetContext(ctx)
	}

	result := &GetBucketsResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
