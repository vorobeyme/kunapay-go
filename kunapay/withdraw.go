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
	Code        string   `json:"code"`
	Asset       string   `json:"asset"`
	Network     string   `json:"network"`
	Position    int64    `json:"position"`
	Name        string   `json:"name"`
	Icon        string   `json:"icon"`
	Description string   `json:"description"`
	CustomTitle string   `json:"customTitle"`
	Fields      []Fields `json:"fields"`
}

// Fields represents a KunaPay withdraw fields.
type Fields struct {
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
	CallbackUrl   string            `json:"callbackUrl,omitempty"`
}

// CreateWithdrawResponse represents a KunaPay create withdraw response.
type CreateWithdrawResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

// Create create withdraw in crypto to any specified address.
// https://docs-pay.kuna.io/reference/withdrawcontroller_makewithdraw
func (s *WithdrawService) Create(ctx context.Context, request CreateWithdrawRequest) (*CreateWithdrawResponse, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "withdraw", request)
	if err != nil {
		return nil, nil, err
	}

	var createResp *CreateWithdrawResponse
	resp, err := s.client.Do(req, &createResp)
	if err != nil {
		return nil, resp, err
	}

	return createResp, resp, err
}

// Methods returns information on available withdraw methods.
// https://docs-pay.kuna.io/reference/withdrawcontroller_prerequestwithdraw
func (s *WithdrawService) Methods(ctx context.Context, asset string) ([]*Withdraw, *http.Response, error) {
	u := fmt.Sprintf("withdraw/pre-request?asset=%s", strings.ToUpper(asset))
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	var withdraws []*Withdraw
	resp, err := s.client.Do(req, &withdraws)
	if err != nil {
		return nil, resp, err
	}

	return withdraws, resp, err
}
