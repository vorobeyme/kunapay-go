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

	mux.HandleFunc("/v1/asset/balance", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testURL(t, r, "/v1/asset/balance?assetCodes=USDT,BTC")
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

	assets, _, err := client.Asset.GetBalance(context.Background(), []string{"USDT", "BTC"}...)
	if err != nil {
		t.Errorf("Asset.GetBalance returned error: %v", err)
	}

	want := []*Asset{assetMock(), assetMock()}
	if !reflect.DeepEqual(assets, want) {
		t.Errorf("Asset.GetBalance returned %+v, expected %+v", assets, want)
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
