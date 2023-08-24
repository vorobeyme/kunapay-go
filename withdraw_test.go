package kunapay

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestWithdrawService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Withdraw{}, "{}")
	m := withdrawMock()
	want := `{
		"code": "USDT",
		"asset": "USDT",
		"network": "ERC20",
		"position": 0,
		"name": "Tether",
		"icon": "https://example.com/assets/currencies/png/usdt.png",
		"description": "Description",
		"customTitle": "Title",
		"fields": [
			{
				"name": "address",
				"label": "Address",
				"description": "Address description",
				"position": 1,
				"type": "text",
				"isRequired": true,
				"isMasked": false,
				"isResultField": true
			}
		]
	}`
	testJSONMarshal(t, m, want)
}

func TestWithdrawService_GetMethods(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/withdraw/pre-request", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testURL(t, r, "/v1/withdraw/pre-request?asset=USDT")
		fmt.Fprint(w, `{
			"data": [
				{
					"code": "USDT",
					"asset": "USDT",
					"network": "ERC20",
					"position": 0,
					"name": "Tether",
					"icon": "https://example.com/assets/currencies/png/usdt.png",
					"description": "Description",
					"customTitle": "Title",
					"fields": [
						{
							"name": "address",
							"label": "Address",
							"description": "Address description",
							"position": 1,
							"type": "text",
							"isRequired": true,
							"isMasked": false,
							"isResultField": true
						}
					]
				}
			]
		}`)
	})

	ctx := context.Background()
	withdraw, _, err := client.Withdraw.GetMethods(ctx, "USDT")
	if err != nil {
		t.Errorf("Withdraw.GetMethods returned error: %v", err)
	}

	want := []*Withdraw{withdrawMock()}
	if !reflect.DeepEqual(withdraw, want) {
		t.Errorf("Withdraw.GetMethods returned %+v, want %+v", withdraw, want)
	}

	const method = "Withdraw.GetMethods"
	testBadPathParams(t, method, func() error {
		_, _, err = client.Withdraw.GetMethods(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, method, client, func() (*Response, error) {
		got, resp, err := client.Withdraw.GetMethods(ctx, "btc")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", method, got)
		}
		return resp, err
	})
}

func TestWithdrawService_Create(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/withdraw", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"data": {
				"id": "da24ed52-6698-49f6-b6b9-a3f5bf79818d",
				"success": true
			}
		}`)
	})

	ctx := context.Background()
	createReq := &CreateWithdrawRequest{
		Amount:        "100.00",
		Asset:         "USDT",
		PaymentMethod: "USDT",
	}
	withdraw, _, err := client.Withdraw.Create(ctx, createReq)
	if err != nil {
		t.Errorf("Withdraw.Create returned error: %v", err)
	}

	want := &CreateWithdrawResponse{
		ID:      "da24ed52-6698-49f6-b6b9-a3f5bf79818d",
		Success: true,
	}
	if !reflect.DeepEqual(withdraw, want) {
		t.Errorf("Withdraw.Create returned %+v, want %+v", withdraw, want)
	}

	const method = "Withdraw.Create"
	testNewRequestAndDoFailure(t, method, client, func() (*Response, error) {
		got, resp, err := client.Withdraw.Create(ctx, &CreateWithdrawRequest{
			Amount:        "100.00",
			Asset:         "USDT",
			PaymentMethod: "USDT",
		})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", method, got)
		}
		return resp, err
	})
}

func TestWithdrawService_CreateWithRequestValidationErr(t *testing.T) {
	client, _, teardown := setupClient()
	defer teardown()

	_, _, amountErr := client.Withdraw.Create(context.Background(), &CreateWithdrawRequest{})
	if amountErr != nil && amountErr.Error() != "amount is required" {
		t.Errorf("Withdraw.Create returned error: %v", amountErr)
	}
	_, _, assetErr := client.Withdraw.Create(context.Background(), &CreateWithdrawRequest{Amount: "100.00"})
	if assetErr != nil && assetErr.Error() != "asset code is required" {
		t.Errorf("Withdraw.Create returned error: %v", assetErr)
	}
	_, _, methodErr := client.Withdraw.Create(context.Background(), &CreateWithdrawRequest{Amount: "100.00", Asset: "USDT"})
	if assetErr != nil && methodErr.Error() != "payment method is required" {
		t.Errorf("Withdraw.Create returned error: %v", methodErr)
	}
}

func withdrawMock() *Withdraw {
	return &Withdraw{
		Code:        "USDT",
		Asset:       "USDT",
		Network:     "ERC20",
		Position:    0,
		Name:        "Tether",
		Icon:        "https://example.com/assets/currencies/png/usdt.png",
		Description: "Description",
		CustomTitle: "Title",
		Fields: []WithdrawField{
			{
				Name:          "address",
				Label:         "Address",
				Description:   "Address description",
				Position:      1,
				Type:          "text",
				IsRequired:    true,
				IsMasked:      false,
				IsResultField: true,
			},
		},
	}
}
