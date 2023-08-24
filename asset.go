package kunapay

import (
	"context"
	"net/http"
	"strings"
)

// AssetService handles communication with the assets related.
type AssetService struct {
	client *Client
}

// Asset represents a KunaPay asset.
type Asset struct {
	Balance     string `json:"balance"`
	LockBalance string `json:"lockBalance"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Icons       struct {
		SVG string `json:"svg"`
		PNG string `json:"png"`
	} `json:"icons"`
}

// GetBalance returns the balance of the assets.
// If assets is empty, returns the balance of all assets,
// otherwise returns the balance of the specified assets.
//
// API docs: https://docs-pay.kuna.io/reference/assetcontroller_getbalances
func (s *AssetService) GetBalance(ctx context.Context, assets ...string) ([]*Asset, *Response, error) {
	u := "asset/balance"
	if len(assets) > 0 {
		var assetCodes []string
		for _, asset := range assets {
			if asset = strings.TrimSpace(asset); asset != "" {
				assetCodes = append(assetCodes, strings.ToUpper(asset))
			}
		}
		u += "?assetCodes=" + strings.ToUpper(strings.Join(assetCodes, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	var root struct {
		Data []*Asset `json:"data"`
	}

	resp, err := s.client.Do(req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, err
}
