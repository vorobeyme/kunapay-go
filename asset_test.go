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

	ctx := context.Background()
	assets, _, err := client.Asset.GetBalance(ctx, []string{"USDT ", "  ", "btc"}...)
	if err != nil {
		t.Errorf("Asset.GetBalance returned error: %v", err)
	}

	want := []*Asset{assetMock(), assetMock()}
	if !reflect.DeepEqual(assets, want) {
		t.Errorf("Asset.GetBalance returned %+v, want %+v", assets, want)
	}

	const method = "Asset.GetBalance"
	testBadPathParams(t, method, func() error {
		_, _, err = client.Transaction.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, method, client, func() (*Response, error) {
		got, resp, err := client.Asset.GetBalance(ctx, []string{"usdt", "btc"}...)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", method, got)
		}
		return resp, err
	})
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
