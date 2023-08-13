package kunapay

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestTransactionService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Transaction{}, "{}")
	m := transactionMock()
	want := `{
		"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
		"address": "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
		"amount": "100.011",
		"asset": "ETH",
		"fee": "0.0000001",
		"processedAmount": "100.011",
		"status": "Processed",
		"paymentCode": "ETH",
		"type": "Deposit",
		"createdAt": "2023-07-30T00:00:00.000Z"
	}`
	testJSONMarshal(t, m, want)
}

func TestTransactionService_Get(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/transaction/c94c0c95-e735-45ea-982e-a95f7f52ca49", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"data": {
				"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
				"address": "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
				"amount": "100.011",
				"asset": "ETH",
				"fee": "0.0000001",
				"processedAmount": "100.011",
				"status": "Processed",
				"paymentCode": "ETH",
				"type": "Deposit",
				"createdAt": "2023-07-30T00:00:00.000Z"
			}
		}`)
	})

	transaction, _, err := client.Transaction.Get(context.Background(), "c94c0c95-e735-45ea-982e-a95f7f52ca49")
	if err != nil {
		t.Errorf("Transaction.Get returned error: %v", err)
	}

	want := transactionMock()

	if !reflect.DeepEqual(transaction, want) {
		t.Errorf("Transaction.Get returned %+v, want %+v", transaction, want)
	}

	testBadPathParams(t, "Transaction.Get", func() error {
		_, _, err = client.Transaction.Get(context.Background(), "")
		return err
	})
}

func TestTransactionService_List(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/v1/transaction", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testURL(t, r, "/v1/transaction?asset=ETH&createdFrom=2023-07-30T15%3A10%3A08%2B03%3A00&createdTo=2023-07-31T15%3A10%3A08%2B03%3A00&orderBy=createdAt&skip=10&take=10")
		fmt.Fprint(w, `{
			"data": [
				{
					"id": "c94c0c95-e735-45ea-982e-a95f7f52ca49",
					"address": "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
					"amount": "100.011",
					"asset": "ETH",
					"fee": "0.0000001",
					"processedAmount": "100.011",
					"status": "Processed",
					"paymentCode": "ETH",
					"type": "Deposit",
					"createdAt": "2023-07-30T00:00:00.000Z"
				}
			]
		}`)
	})

	createdFrom, _ := time.Parse(time.RFC3339, "2023-07-30T15:10:08+03:00")
	createdTo, _ := time.Parse(time.RFC3339, "2023-07-31T15:10:08+03:00")
	transactionListOptions := &TransactionListOpts{
		Take:        10,
		Skip:        10,
		Asset:       "ETH",
		CreatedFrom: &createdFrom,
		CreatedTo:   &createdTo,
		OrderBy:     "createdAt",
	}
	transactions, _, err := client.Transaction.List(context.Background(), transactionListOptions)
	if err != nil {
		t.Errorf("Transaction.List returned error: %v", err)
	}

	want := []*Transaction{transactionMock()}
	if !reflect.DeepEqual(transactions, want) {
		t.Errorf("Transaction.List returned %+v, want %+v", transactions, want)
	}
}

func transactionMock() *Transaction {
	return &Transaction{
		ID:              "c94c0c95-e735-45ea-982e-a95f7f52ca49",
		Address:         "tb1q0xrgwsd7e0uad3sy98klppjwjq26023mcx224d",
		Amount:          "100.011",
		Asset:           "ETH",
		Fee:             "0.0000001",
		ProcessedAmount: "100.011",
		Status:          "Processed",
		PaymentCode:     "ETH",
		Type:            "Deposit",
		CreatedAt:       "2023-07-30T00:00:00.000Z",
	}
}
