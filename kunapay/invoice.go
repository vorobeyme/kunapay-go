package kunapay

import (
	"fmt"
	"net/http"
)

// InvoiceService handles communication with the invoice related
type InvoiceService struct {
	client *Client
}

// Invoice represents a KunaPay invoice.
type Invoice struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	AddressID        string `json:"addressId"`
	ExternalOrderID  string `json:"externalOrderId"`
	PaymentAmount    string `json:"paymentAmount"`
	InvoiceAmount    string `json:"invoiceAmount"`
	InvoiceAssetCode string `json:"invoiceAssetCode"`
	PaymentAssetCode string `json:"paymentAssetCode"`
	ExpireAt         string `json:"expireAt"`
	CompletedAt      string `json:"completedAt"`
	CreatedAt        string `json:"createdAt"`
}

// InvoiceDetail represents a KunaPay invoice details.
type InvoiceDetail struct {
	Invoice
	CreatorID          string       `json:"creatorId"`
	Transactions       Transactions `json:"transactions"`
	ProductCategory    string       `json:"productCategory"`
	ProductDescription string       `json:"productDescription"`
	IsCreatedByAPI     bool         `json:"isCreatedByApi"`
	UpdateAt           string       `json:"updatedAt"`
}

// Transactions represents a KunaPay transactions associated with the invoice.
type Transactions struct {
	ID              string   `json:"id"`
	Address         string   `json:"address"`
	Amount          string   `json:"amount"`
	Asset           string   `json:"asset"`
	CreatorComment  string   `json:"creatorComment"`
	Fee             string   `json:"fee"`
	ProcessedAmount string   `json:"processedAmount"`
	Reason          []string `json:"reason"`
	Status          string   `json:"status"`
	Type            string   `json:"type"`
	CreateAt        string   `json:"createdAt"`
	UpdateAt        string   `json:"updatedAt"`
	PaymentCode     string   `json:"paymentCode"`
}

// InvoiceCurrency represents a KunaPay invoice currencies.
type InvoiceCurrency struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Position    int64  `json:"position"`
	Precision   int64  `json:"precision"`
	Type        string `json:"type"`
	Icons       struct {
		SVG string `json:"svg"`
		PNG string `json:"png"`
	} `json:"icons"`
}

type CreateInvoiceRequest struct {
	Amount             string `json:"amount"`
	Asset              string `json:"asset"`
	ExternalOrderID    string `json:"externalOrderId"`
	ProductDescription string `json:"productDescription"`
	ProductCategory    string `json:"productCategory"`
	CallbackUrl        string `json:"callbackUrl"`
}

type CreateInvoiceResponse struct {
	ID          string `json:"id"`
	PaymentLink string `json:"paymentLink"`
}

// Create creates invoice for a client for a specified amount.
// https://docs-pay.kuna.io/reference/invoicecontroller_createinvoice
func (s *InvoiceService) Create(request CreateInvoiceRequest) (*CreateInvoiceResponse, *http.Response, error) {
	req, err := s.client.NewRequest("POST", "invoice", request)
	if err != nil {
		return nil, nil, err
	}

	var createResp *CreateInvoiceResponse
	resp, err := s.client.Do(req, &createResp)
	if err != nil {
		return nil, resp, err
	}

	return createResp, resp, err
}

// List returns crypto invoices list.
// https://docs-pay.kuna.io/reference/invoicecontroller_getinvoices
func (s *InvoiceService) List() ([]*Invoice, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "invoice", nil)
	if err != nil {
		return nil, nil, err
	}

	var invoices []*Invoice
	resp, err := s.client.Do(req, &invoices)
	if err != nil {
		return nil, resp, err
	}

	return invoices, resp, err
}

// Get returns detailed information on a single crypto invoice.
// The invoice identifier is passed in the ID parameter.
// https://docs-pay.kuna.io/reference/invoicecontroller_getinvoicebyid
func (s *InvoiceService) Get(ID string) (*InvoiceDetail, *http.Response, error) {
	u := fmt.Sprintf("invoice/%s", ID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var invoice *InvoiceDetail
	resp, err := s.client.Do(req, &invoice)
	if err != nil {
		return nil, resp, err
	}

	return invoice, resp, err
}

// Currencies returns information on available crypto currencies for invoice creation.
// https://docs-pay.kuna.io/reference/invoicecontroller_getinvoiceassets
func (s *InvoiceService) Currencies() ([]*InvoiceCurrency, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "invoice/assets", nil)
	if err != nil {
		return nil, nil, err
	}

	var currencies []*InvoiceCurrency
	resp, err := s.client.Do(req, &currencies)
	if err != nil {
		return nil, resp, err
	}

	return currencies, resp, err
}
