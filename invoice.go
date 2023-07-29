package kunapay

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// InvoiceService handles communication with the invoice related
type InvoiceService struct {
	client *Client
}

// The statuses of the invoice.
const (
	InvoiceStatusCreated              = "CREATED"
	InvoiceStatusPaymentAwaiting      = "PAYMENT_AWAITING"
	InvoiceStatusConfirmetionAwaiting = "CONFIRMATION_AWAITING"
	InvoiceStatusLimitsOutOfRange     = "LIMITS_OUT_OF_RANGE"
	InvoiceStatusPaid                 = "PAID"
	InvoiceStatusPartiallyPaid        = "PARTIALLY_PAID"
	InvoiceStatusTimeout              = "TIMEOUT"
	InvoiceStatusDeactivated          = "DEACTIVATED"
	InvoiceStatusDeclined             = "DECLINED"
)

// Invoice represents a KunaPay invoice response.
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

// InvoiceDetail represents a KunaPay invoice details response.
type InvoiceDetail struct {
	Invoice
	CreatorID          string               `json:"creatorId"`
	Transactions       []InvoiceTransaction `json:"transactions"`
	ProductCategory    string               `json:"productCategory"`
	ProductDescription string               `json:"productDescription"`
	IsCreatedByAPI     bool                 `json:"isCreatedByApi"`
	UpdateAt           string               `json:"updatedAt"`
}

// Transactions represents a KunaPay transactions associated with the invoice.
type InvoiceTransaction struct {
	Transaction
	CreatorComment string   `json:"creatorComment"`
	Reason         []string `json:"reason"`
	UpdateAt       string   `json:"updatedAt"`
}

// InvoiceCurrency represents a KunaPay invoice currencies response.
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
	ExternalOrderID    string `json:"externalOrderId,omitempty"`
	ProductDescription string `json:"productDescription,omitempty"`
	ProductCategory    string `json:"productCategory,omitempty"`
	CallbackUrl        string `json:"callbackUrl,omitempty"`
}

func (r *CreateInvoiceRequest) validate() error {
	if r.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	if r.Asset == "" {
		return fmt.Errorf("asset code is required")
	}
	return nil
}

type CreateInvoiceResponse struct {
	ID          string `json:"id"`
	PaymentLink string `json:"paymentLink"`
}

// Create creates invoice for a client for a specified amount.
// https://docs-pay.kuna.io/reference/invoicecontroller_createinvoice
func (s *InvoiceService) Create(ctx context.Context, request CreateInvoiceRequest) (*CreateInvoiceResponse, *http.Response, error) {
	if err := request.validate(); err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, "invoice", request)
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

type InvoiceOrderBy string

const (
	OrderByCreatedAt        InvoiceOrderBy = "createdAt"
	OrderByCompletedAt      InvoiceOrderBy = "completedAt"
	OrderByInvoiceAssetCode InvoiceOrderBy = "invoiceAssetCode"
	OrderByPaymentAssetCode InvoiceOrderBy = "paymentAssetCode"
	OrderByStatus           InvoiceOrderBy = "status"
)

// InvoiceListOpts specifies the optional parameters to the InvoiceService.List method.
type InvoiceListOpts struct {
	Take             int64
	Skip             int64
	CreatedFrom      *time.Time
	CreatedTo        *time.Time
	CompletedFrom    *time.Time
	CompletedTo      *time.Time
	ExternalOrderID  string
	InvoiceAssetCode string
	PaymentAssetCode string
	OrderBy          InvoiceOrderBy
}

func (o *InvoiceListOpts) values() url.Values {
	v := url.Values{}
	if o.Take > 0 {
		v.Add("take", fmt.Sprintf("%d", o.Take))
	}
	if o.Skip > 0 {
		v.Add("skip", fmt.Sprintf("%d", o.Skip))
	}
	if o.CreatedFrom != nil {
		v.Add("createdFrom", o.CreatedFrom.Format(time.RFC3339))
	}
	if o.CreatedTo != nil {
		v.Add("createdTo", o.CreatedTo.Format(time.RFC3339))
	}
	if o.CompletedFrom != nil {
		v.Add("completedFrom", o.CompletedFrom.Format(time.RFC3339))
	}
	if o.CompletedTo != nil {
		v.Add("completedTo", o.CompletedTo.Format(time.RFC3339))
	}
	if o.ExternalOrderID != "" {
		v.Add("externalOrderId", o.ExternalOrderID)
	}
	if o.InvoiceAssetCode != "" {
		v.Add("invoiceAssetCode", o.InvoiceAssetCode)
	}
	if o.PaymentAssetCode != "" {
		v.Add("paymentAssetCode", o.PaymentAssetCode)
	}
	if o.OrderBy != "" {
		v.Add("orderBy", string(o.OrderBy))
	}

	return v
}

// List returns crypto invoices list.
// https://docs-pay.kuna.io/reference/invoicecontroller_getinvoices
func (s *InvoiceService) List(ctx context.Context, opts *InvoiceListOpts) ([]*Invoice, *http.Response, error) {
	u := "invoice"
	if opts != nil {
		u += "?" + opts.values().Encode()
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, http.NoBody)
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
func (s *InvoiceService) Get(ctx context.Context, ID string) (*InvoiceDetail, *http.Response, error) {
	u := fmt.Sprintf("invoice/%s", ID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, http.NoBody)
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

// InvoiceUpdateOpts specifies the optional parameters to the InvoiceService.Currencies method.
type InvoiceCurrencyListOpts struct {
	Take int64
	Skip int64
}

func (o *InvoiceCurrencyListOpts) values() url.Values {
	v := url.Values{}
	if o.Take > 0 {
		v.Add("take", fmt.Sprintf("%d", o.Take))
	}
	if o.Skip > 0 {
		v.Add("skip", fmt.Sprintf("%d", o.Skip))
	}

	return v
}

// GetCurrencies returns information on available crypto currencies for invoice creation.
// https://docs-pay.kuna.io/reference/invoicecontroller_getinvoiceassets
func (s *InvoiceService) GetCurrencies(ctx context.Context, opts *InvoiceCurrencyListOpts) ([]*InvoiceCurrency, *http.Response, error) {
	u := "invoice/assets"
	if opts != nil {
		u += "?" + opts.values().Encode()
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, http.NoBody)
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
