package kunapay

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// WithdrawService handles communication with the withdraw related.
type WithdrawService struct {
	client *Client
}

// Withdraw represents a KunaPay withdraw.
type Withdraw struct {
	Code        string          `json:"code"`
	Asset       string          `json:"asset"`
	Network     string          `json:"network"`
	Position    int64           `json:"position"`
	Name        string          `json:"name"`
	Icon        string          `json:"icon"`
	Description string          `json:"description"`
	CustomTitle string          `json:"customTitle"`
	Fields      []WithdrawField `json:"fields"`
}

// Field represents a KunaPay withdraw fields that should be used
// with withdraw request.
type WithdrawField struct {
	Name          string `json:"name"`
	Label         string `json:"label"`
	Description   string `json:"description"`
	Position      int64  `json:"position"`
	Type          string `json:"type"`
	IsRequired    bool   `json:"isRequired"`
	IsMasked      bool   `json:"isMasked"`
	IsResultField bool   `json:"isResultField"`
}

// CreateWithdrawRequest represents a KunaPay create withdraw request.
type CreateWithdrawRequest struct {
	Amount        string            `json:"amount"`
	Asset         string            `json:"asset"`
	PaymentMethod string            `json:"paymentMethod"`
	Fields        map[string]string `json:"fields,omitempty"`
	Comment       string            `json:"comment,omitempty"`
	WithdrawAll   bool              `json:"withdrawAll,omitempty"`
	CallbackURL   string            `json:"callbackUrl,omitempty"`
}

// validate checks if request values are valid.
func (r *CreateWithdrawRequest) validate() error {
	if r.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	if r.Asset == "" {
		return fmt.Errorf("asset code is required")
	}
	if r.PaymentMethod == "" {
		return fmt.Errorf("payment method is required")
	}

	return nil
}

// CreateWithdrawResponse represents a KunaPay create withdraw response.
type CreateWithdrawResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

// Create create withdraw in crypto to any specified address.
//
// API docs: https://docs-pay.kuna.io/reference/withdrawcontroller_makewithdraw
func (s *WithdrawService) Create(ctx context.Context, request *CreateWithdrawRequest) (*CreateWithdrawResponse, *http.Response, error) {
	if err := request.validate(); err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, "withdraw", request)
	if err != nil {
		return nil, nil, err
	}

	var root struct {
		Data *CreateWithdrawResponse `json:"data"`
	}

	resp, err := s.client.Do(req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, err
}

// GetMethods returns information on available withdraw methods.
//
// API docs: https://docs-pay.kuna.io/reference/withdrawcontroller_prerequestwithdraw
func (s *WithdrawService) GetMethods(ctx context.Context, asset string) ([]*Withdraw, *http.Response, error) {
	if asset == "" {
		return nil, nil, fmt.Errorf("asset code is required")
	}
	u := fmt.Sprintf("withdraw/pre-request?asset=%s", strings.ToUpper(asset))
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	var root struct {
		Data []*Withdraw `json:"data"`
	}

	resp, err := s.client.Do(req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, err
}
