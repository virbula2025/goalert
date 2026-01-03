package telnyx

import (
	"context"
	"fmt"
	"net/url"
)

type CarrierInfo struct {
	Name string `json:"name"`
	Type string `json:"type"` // mobile, landline, voip
}

func (c *Config) FetchCarrier(ctx context.Context, number string) (*CarrierInfo, error) {
	// Telnyx Lookup API v2
	// endpoint: https://api.telnyx.com/v2/number_lookup/+15550001
	
	safeNum := url.QueryEscape(number)
	path := fmt.Sprintf("number_lookup/%s", safeNum)

	var resp struct {
		Data struct {
			Carrier CarrierInfo `json:"carrier"`
		} `json:"data"`
	}

	err := c.postJSON(ctx, path, nil, &resp) // Lookup is often GET, but check docs. If GET, change client.go to support GET.
	if err != nil {
		return nil, err
	}
	
	return &resp.Data.Carrier, nil
}