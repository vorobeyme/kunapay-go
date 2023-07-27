package kunapay

import "net/http"

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

// Balance returns the balance of the asset.
// https://docs-pay.kuna.io/reference/assetcontroller_getbalances
func (s *AssetService) Balance() ([]*Asset, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "assets/balance", nil)
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
