package kunapay

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAssetService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Asset{}, "{}")
	m := assetMock()
	want := `{
		"balance": "123.99",
		"lockBalance": "23.00",
		"code": "USDT",
		"name": "Tether",
		"icons": {
			"svg": "https://example.com/assets/currencies/svg/usdt.svg",
			"png": "https://example.com/assets/currencies/png/usdt.png"
		}
	}`
	testJSONMarshal(t, m, want)
}

func TestAssetService_GetBalance(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	expectedAssets := []*Asset{assetMock(), assetMock()}

	mux.HandleFunc("/assets/balance", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testURL(t, r, "/assets/balance?assetCodes=USDT,BTC")
		fmt.Fprint(w, `{
			"data": [
				{
					"balance": "123.99",
					"lockBalance": "23.00",
					"code": "USDT",
					"name": "Tether",
					"icons": {
						"svg": "https://example.com/assets/currencies/svg/usdt.svg",
						"png": "https://example.com/assets/currencies/png/usdt.png"
					}
				},
				{
					"balance": "123.99",
					"lockBalance": "23.00",
					"code": "USDT",
					"name": "Tether",
					"icons": {
						"svg": "https://example.com/assets/currencies/svg/usdt.svg",
						"png": "https://example.com/assets/currencies/png/usdt.png"
					}
				}
			]
		}`)
	})

	assets, _, err := client.Asset.GetBalance(context.Background(), &BalanceListOpts{AssetCodes: []string{"USDT", "BTC"}})
	if err != nil {
		t.Errorf("Asset.GetBalance returnted error: %v", err)
	}

	if !reflect.DeepEqual(assets, expectedAssets) {
		t.Errorf("Asset.GetBalance returned %+v, expected %+v", assets, expectedAssets)
	}
}

func assetMock() *Asset {
	return &Asset{
		Balance:     "123.99",
		LockBalance: "23.00",
		Code:        "USDT",
		Name:        "Tether",
		Icons: struct {
			SVG string `json:"svg"`
			PNG string `json:"png"`
		}{
			SVG: "https://example.com/assets/currencies/svg/usdt.svg",
			PNG: "https://example.com/assets/currencies/png/usdt.png",
		},
	}
}