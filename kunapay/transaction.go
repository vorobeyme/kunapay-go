package kunapay

import "net/http"

// TransactionService handles communication with the transaction related.
type TransactionService struct {
	client *Client
}

// Transaction represents a KunaPay transaction.
type Transaction struct {
	ID              string `json:"id"`
	Address         string `json:"address"`
	Amount          string `json:"amount"`
	Asset           string `json:"asset"`
	Fee             string `json:"fee"`
	ProcessedAmount string `json:"processedAmount"`
	Status          string `json:"status"`
	PaymentCode     string `json:"paymentCode"`
	Type            string `json:"type"`
	CreatedAt       string `json:"createdAt"`
	InvoiceID       string `json:"invoiceId"`
}

// List returns information on all invoices and withdrawal operations.
// https://docs-pay.kuna.io/reference/transactioncontroller_gettransactions
func (s *TransactionService) List() ([]*Transaction, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "transaction", nil)
	if err != nil {
		return nil, nil, err
	}

	var transactions []*Transaction
	resp, err := s.client.Do(req, &transactions)
	if err != nil {
		return nil, resp, err
	}

	return transactions, resp, err
}

// Get returns detailed information on a single transaction.
// The transaction identifier is passed in the ID parameter.
// https://docs-pay.kuna.io/reference/transactioncontroller_gettransactionbyid
func (s *TransactionService) Get(ID string) (*Transaction, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "transaction/"+ID, nil)
	if err != nil {
		return nil, nil, err
	}

	var transaction *Transaction
	resp, err := s.client.Do(req, &transaction)
	if err != nil {
		return nil, resp, err
	}

	return transaction, resp, err
}
