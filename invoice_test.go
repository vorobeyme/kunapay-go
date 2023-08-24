package kunapay

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestInvoiceService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Invoice{}, "{}")
	m := invoiceMock()
	want := `{
		"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
		"status": "CREATED",
		"addressId": "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
		"externalOrderId": "c94c0c95-e735-45ea-982e-111111111111",
		"paymentAmount": "100.011",
		"invoiceAmount": "1000.011",
		"invoiceAssetCode": "UAH",
		"paymentAssetCode": "ETH",
		"expireAt": "2023-07-31T00:00:00.000Z",
		"completedAt": "2023-07-30T00:00:00.000Z",
		"createdAt": "2023-07-29T00:00:00.000Z"
	}`
	testJSONMarshal(t, m, want)
}

func TestInvoiceService_CreateWithRequiredParams(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/invoice", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"amount":"100.11","asset":"USDT"}`+"\n")
		fmt.Fprint(w, `{
			"data": {
				"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
				"paymentLink": "https://example.com/invoice/c94c0c95-e735-45ea-982e-a95f7f52ca49"
			}
		}`)
	})

	ctx := context.Background()
	invoice, _, err := client.Invoice.Create(ctx, &CreateInvoiceRequest{Amount: "100.11", Asset: "USDT"})
	if err != nil {
		t.Errorf("Invoice.Create returned error: %v", err)
	}

	want := &CreateInvoiceResponse{
		ID:          "c94c0c95-e735-45ea-982e-a95f7f52ca49",
		PaymentLink: "https://example.com/invoice/c94c0c95-e735-45ea-982e-a95f7f52ca49",
	}

	if !reflect.DeepEqual(invoice, want) {
		t.Errorf("Invoice.Create returned %+v, want %+v", invoice, want)
	}

	const method = "Invoice.Create"
	testNewRequestAndDoFailure(t, method, client, func() (*Response, error) {
		got, resp, err := client.Invoice.Create(ctx, &CreateInvoiceRequest{Amount: "100.11", Asset: "USDT"})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", method, got)
		}
		return resp, err
	})
}

func TestInvoiceService_CreateWithAllParams(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/invoice", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"amount":"100.011","asset":"USDT","externalOrderId":"c94c0c95-e735-45ea-982e-111111111111","productDescription":"Product description","productCategory":"Product category","callbackUrl":"https://example.com/callback"}`+"\n")
		fmt.Fprint(w, `{
			"data": {
				"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
				"paymentLink": "https://example.com/invoice/c94c0c95-e735-45ea-982e-a95f7f52ca49"
			}
		}`)
	})

	externalOrderID := "c94c0c95-e735-45ea-982e-111111111111"
	callbackURL := "https://example.com/callback"

	invoice, _, err := client.Invoice.Create(context.Background(), &CreateInvoiceRequest{
		Amount:             "100.011",
		Asset:              "USDT",
		ExternalOrderID:    externalOrderID,
		ProductDescription: "Product description",
		ProductCategory:    "Product category",
		CallbackURL:        callbackURL,
	})
	if err != nil {
		t.Errorf("Invoice.Create returned error: %v", err)
	}

	want := &CreateInvoiceResponse{
		ID:          "c94c0c95-e735-45ea-982e-a95f7f52ca49",
		PaymentLink: "https://example.com/invoice/c94c0c95-e735-45ea-982e-a95f7f52ca49",
	}

	if !reflect.DeepEqual(invoice, want) {
		t.Errorf("Invoice.Create returned %+v, want %+v", invoice, want)
	}
}

func TestInvoiceService_CreateWithRequestValidationErr(t *testing.T) {
	client, _, teardown := setupClient()
	defer teardown()

	ctx := context.Background()
	_, _, amountErr := client.Invoice.Create(ctx, &CreateInvoiceRequest{})
	if amountErr != nil && amountErr.Error() != "amount is required" {
		t.Errorf("Invoice.Create returned error: %v", amountErr)
	}
	_, _, assetErr := client.Invoice.Create(ctx, &CreateInvoiceRequest{Amount: "100.00"})
	if amountErr != nil && assetErr.Error() != "asset code is required" {
		t.Errorf("Invoice.Create returned error: %v", assetErr)
	}
}

func TestInvoiceService_List(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/invoice", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testURL(t, r, "/v1/invoice")
		fmt.Fprint(w, `{
			"data": [
				{
					"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
					"status": "CREATED",
					"addressId": "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
					"externalOrderId": "c94c0c95-e735-45ea-982e-111111111111",
					"paymentAmount": "100.011",
					"invoiceAmount": "1000.011",
					"invoiceAssetCode": "UAH",
					"paymentAssetCode": "ETH",
					"expireAt": "2023-07-31T00:00:00.000Z",
					"completedAt": "2023-07-30T00:00:00.000Z",
					"createdAt": "2023-07-29T00:00:00.000Z"
				}
			]
		}`)
	})

	ctx := context.Background()
	invoice, _, err := client.Invoice.List(ctx, nil)
	if err != nil {
		t.Errorf("Invoice.List returned error: %v", err)
	}

	want := []*Invoice{invoiceMock()}
	if !reflect.DeepEqual(invoice, want) {
		t.Errorf("Invoice.List returned %+v, want %+v", invoice, want)
	}

	const method = "Invoice.List"
	testNewRequestAndDoFailure(t, method, client, func() (*Response, error) {
		got, resp, err := client.Invoice.List(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", method, got)
		}
		return resp, err
	})
}

func TestInvoiceService_ListWithQueryParams(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/invoice", func(w http.ResponseWriter, r *http.Request) {
		expectedURI := "/v1/invoice?" +
			"completedFrom=2023-07-30T15%3A10%3A08%2B03%3A00&" +
			"completedTo=2023-07-31T15%3A10%3A08%2B03%3A00&" +
			"createdFrom=2023-07-30T15%3A10%3A08%2B03%3A00&" +
			"createdTo=2023-07-31T15%3A10%3A08%2B03%3A00&" +
			"externalOrderId=c94c0c95-e735-45ea-982e-111111111111&" +
			"invoiceAssetCode=UAH&" +
			"orderBy=createdAt&" +
			"paymentAssetCode=USDT&" +
			"skip=10&" +
			"take=10"

		testMethod(t, r, "GET")
		testURL(t, r, expectedURI)
		fmt.Fprint(w, `{
			"data": [
				{
					"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
					"status": "CREATED",
					"addressId": "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
					"externalOrderId": "c94c0c95-e735-45ea-982e-111111111111",
					"paymentAmount": "100.011",
					"invoiceAmount": "1000.011",
					"invoiceAssetCode": "UAH",
					"paymentAssetCode": "ETH",
					"expireAt": "2023-07-31T00:00:00.000Z",
					"completedAt": "2023-07-30T00:00:00.000Z",
					"createdAt": "2023-07-29T00:00:00.000Z"
				}
			]
		}`)
	})

	createdFrom, _ := time.Parse(time.RFC3339, "2023-07-30T15:10:08+03:00")
	createdTo, _ := time.Parse(time.RFC3339, "2023-07-31T15:10:08+03:00")
	completedFrom, _ := time.Parse(time.RFC3339, "2023-07-30T15:10:08+03:00")
	completedTo, _ := time.Parse(time.RFC3339, "2023-07-31T15:10:08+03:00")
	invoiceListOpts := &InvoiceListOpts{
		Take:             10,
		Skip:             10,
		CreatedFrom:      &createdFrom,
		CreatedTo:        &createdTo,
		CompletedFrom:    &completedFrom,
		CompletedTo:      &completedTo,
		ExternalOrderID:  "c94c0c95-e735-45ea-982e-111111111111",
		InvoiceAssetCode: "UAH",
		PaymentAssetCode: "USDT",
		OrderBy:          InvoiceOrderByCreatedAt,
	}
	invoice, _, err := client.Invoice.List(context.Background(), invoiceListOpts)
	if err != nil {
		t.Errorf("Invoice.List returned error: %v", err)
	}

	want := []*Invoice{invoiceMock()}
	if !reflect.DeepEqual(invoice, want) {
		t.Errorf("Invoice.List returned %+v, want %+v", invoice, want)
	}
}

func TestInvoiceService_Get(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/invoice/c94c0c95-e735-45ea-982e-a95f7f52ca49", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testURL(t, r, "/v1/invoice/c94c0c95-e735-45ea-982e-a95f7f52ca49")
		fmt.Fprint(w, `{
			"data": {
				"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
				"status": "CREATED",
				"externalOrderId": "c94c0c95-e735-45ea-982e-111111111111",
				"addressId": "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
				"creatorId": "c94c0c95-e735-45ea-982e-111111111111",
				"invoiceAmount": "1000.011",
				"paymentAmount": "100.011",
				"invoiceAssetCode": "UAH",
				"paymentAssetCode": "ETH",
				"transactions": [
					{
						"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
						"address": "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
						"amount": "0.011",
						"asset": "ETH",
						"creatorComment": "Creator comment",
						"fee": "0.0000001",
						"processedAmount": "0.011",
						"reason": [
							"reason1"
						],
						"status": "Processed",
						"type": "Deposit",
						"createdAt": "2023-07-30T00:00:00.000Z",
						"updatedAt": "2023-07-30T00:00:00.000Z",
						"paymentCode": "ETH"
					}
				],
				"productCategory": "Product category",
				"productDescription": "Product description",
				"isCreatedByApi": true,
				"expireAt": "2023-07-31T00:00:00.000Z",
				"completedAt": "2023-07-30T00:00:00.000Z",
				"createdAt": "2023-07-29T00:00:00.000Z",
				"updatedAt": "2023-07-29T00:00:00.000Z"
			}
		}`)
	})

	ctx := context.Background()
	invoice, _, err := client.Invoice.Get(ctx, "c94c0c95-e735-45ea-982e-a95f7f52ca49")
	if err != nil {
		t.Errorf("Invoice.Get returned error: %v", err)
	}

	want := &InvoiceDetail{
		ID:               "c94c0c95-e735-45ea-982e-a95f7f52ca49",
		Status:           "CREATED",
		ExternalOrderID:  "c94c0c95-e735-45ea-982e-111111111111",
		AddressID:        "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
		CreatorID:        "c94c0c95-e735-45ea-982e-111111111111",
		InvoiceAmount:    "1000.011",
		PaymentAmount:    "100.011",
		InvoiceAssetCode: "UAH",
		PaymentAssetCode: "ETH",
		Transactions: []InvoiceTransaction{
			{
				Address:         "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
				Amount:          "0.011",
				Asset:           "ETH",
				CreatorComment:  "Creator comment",
				Fee:             "0.0000001",
				ID:              "c94c0c95-e735-45ea-982e-a95f7f52ca49",
				ProcessedAmount: "0.011",
				Reason:          []string{"reason1"},
				Status:          "Processed",
				Type:            "Deposit",
				CreatedAt:       "2023-07-30T00:00:00.000Z",
				UpdatedAt:       "2023-07-30T00:00:00.000Z",
				PaymentCode:     "ETH",
			},
		},
		ProductCategory:    "Product category",
		ProductDescription: "Product description",
		IsCreatedByAPI:     true,
		ExpireAt:           "2023-07-31T00:00:00.000Z",
		CompletedAt:        "2023-07-30T00:00:00.000Z",
		CreatedAt:          "2023-07-29T00:00:00.000Z",
		UpdateAt:           "2023-07-29T00:00:00.000Z",
	}

	if !reflect.DeepEqual(invoice, want) {
		t.Errorf("Invoice.Get returned %+v, want %+v", invoice, want)
	}

	const method = "Invoice.Get"
	testBadPathParams(t, method, func() error {
		_, _, err = client.Invoice.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, method, client, func() (*Response, error) {
		got, resp, err := client.Invoice.Get(ctx, "c94c0c95-e735-45ea-982e-a95f7f52ca49")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", method, got)
		}
		return resp, err
	})
}

func TestInvoiceService_GetCurrencies(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/invoice/assets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testURL(t, r, "/v1/invoice/assets?skip=10&take=10")
		fmt.Fprint(w, `{
			"data": [
				{
					"code": "USDT",
					"name": "Tether",
					"description": "Description",
					"position": 1,
					"precision": 8,
					"type": "crypto",
					"icons": {
						"svg": "https://example.com/assets/currencies/svg/usdt.svg",
						"png": "https://example.com/assets/currencies/png/usdt.png"
					}
				}
			]
		}`)
	})

	ctx := context.Background()
	invoiceCurrencyListOpts := &InvoiceCurrencyListOpts{
		Take: 10,
		Skip: 10,
	}
	assets, _, err := client.Invoice.GetCurrencies(ctx, invoiceCurrencyListOpts)
	if err != nil {
		t.Errorf("Invoice.GetCurrencies returned error: %v", err)
	}

	want := []*InvoiceCurrency{
		{
			Code:        "USDT",
			Name:        "Tether",
			Description: "Description",
			Position:    1,
			Precision:   8,
			Type:        "crypto",
			Icons: struct {
				SVG string "json:\"svg\""
				PNG string "json:\"png\""
			}{
				SVG: "https://example.com/assets/currencies/svg/usdt.svg",
				PNG: "https://example.com/assets/currencies/png/usdt.png",
			},
		},
	}
	if !reflect.DeepEqual(assets, want) {
		t.Errorf("Invoice.GetCurrencies returned %+v, want %+v", assets, want)
	}

	const method = "Invoice.GetCurrencies"
	testNewRequestAndDoFailure(t, method, client, func() (*Response, error) {
		got, resp, err := client.Invoice.GetCurrencies(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", method, got)
		}
		return resp, err
	})
}

func invoiceMock() *Invoice {
	return &Invoice{
		ID:               "c94c0c95-e735-45ea-982e-a95f7f52ca49",
		Status:           "CREATED",
		AddressID:        "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
		ExternalOrderID:  "c94c0c95-e735-45ea-982e-111111111111",
		PaymentAmount:    "100.011",
		InvoiceAmount:    "1000.011",
		InvoiceAssetCode: "UAH",
		PaymentAssetCode: "ETH",
		ExpireAt:         "2023-07-31T00:00:00.000Z",
		CompletedAt:      "2023-07-30T00:00:00.000Z",
		CreatedAt:        "2023-07-29T00:00:00.000Z",
	}
}
