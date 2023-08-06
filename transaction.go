package kunapay

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TransactionService handles communication with the transaction related.
type TransactionService struct {
	client *Client
}

// Transaction statuses.
const (
	TransactionStatusCreated            = "Created"
	TransactionStatusCanceled           = "Canceled"
	TransactionStatusProcessing         = "Processing"
	TransactionStatusProcessed          = "Processed"
	TransactionStatusPartiallyProcessed = "PartiallyProcessed"
)

// Transaction types.
const (
	TransactionTypeDeposit  = "Deposit"
	TransactionTypeWithdraw = "Withdraw"
	TransactionTypeRefund   = "Refund"
)

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
	InvoiceID       string `json:"invoiceId,omitempty"`
}

// TransactionListOpts specifies the optional parameters to the
// TransactionService.List method.
type TransactionListOpts struct {
	Take        int64
	Skip        int64
	Asset       string
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	OrderBy     string
}

// values converts TransactionListOpts to url.Values to be used in query string.
func (o *TransactionListOpts) values() url.Values {
	v := url.Values{}
	if o.Take > 0 {
		v.Add("take", fmt.Sprintf("%d", o.Take))
	}
	if o.Skip > 0 {
		v.Add("skip", fmt.Sprintf("%d", o.Skip))
	}
	if o.Asset != "" {
		v.Add("asset", string(o.Asset))
	}
	if o.CreatedFrom != nil {
		v.Add("createdFrom", o.CreatedFrom.Format(time.RFC3339))
	}
	if o.CreatedTo != nil {
		v.Add("createdTo", o.CreatedTo.Format(time.RFC3339))
	}
	if o.OrderBy != "" {
		v.Add("orderBy", o.OrderBy)
	}

	return v
}

// List returns information on all invoices and withdrawal operations.
// 
// API docs: https://docs-pay.kuna.io/reference/transactioncontroller_gettransactions
func (s *TransactionService) List(ctx context.Context, opts *TransactionListOpts) ([]*Transaction, *http.Response, error) {
	u := "transaction"
	if opts != nil {
		u += "?" + opts.values().Encode()
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	var root struct {
		Data []*Transaction `json:"data"`
	}

	resp, err := s.client.Do(req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, err
}

// Get returns detailed information on a single transaction.
// The transaction identifier is passed in the ID parameter.
//
// API docs: https://docs-pay.kuna.io/reference/transactioncontroller_gettransactionbyid
func (s *TransactionService) Get(ctx context.Context, ID string) (*Transaction, *http.Response, error) {
	if ID == "" {
		return nil, nil, fmt.Errorf("transaction ID is required")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, "transaction/"+ID, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	var root struct {
		Data *Transaction `json:"data"`
	}

	resp, err := s.client.Do(req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, err
}
