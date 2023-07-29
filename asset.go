package kunapay

import (
	"context"
	"net/http"
	"strings"
)

// AssetService handles communication with the asset related
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

// BalanceListOpts specifies the optional parameters to the AssetService.Balance method.
type BalanceListOpts struct {
	AssetCodes []string
}

// values returns a string of values in the format: "BTC,ETH,..."
func (o *BalanceListOpts) values() string {
	assets := strings.Join(o.AssetCodes[:], ",")
	return strings.ToUpper(assets)
}

// GetBalance returns the balance of the asset.
// https://docs-pay.kuna.io/reference/assetcontroller_getbalances
func (s *AssetService) GetBalance(ctx context.Context, opts *BalanceListOpts) ([]*Asset, *http.Response, error) {
	u := "assets/balance"
	if opts != nil {
		u += "?assetCodes=" + opts.values()
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	var assets []*Asset
	resp, err := s.client.Do(req, &assets)
	if err != nil {
		return nil, resp, err
	}

	return assets, resp, err
}
