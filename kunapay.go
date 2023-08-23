package kunapay

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	libVersion = "0.1.0"

	apiURL     = "https://api-kunapayapp.kuna.io/"
	apiVersion = "v1"

	userAgent = "kunapay-go/" + libVersion

	headerNonce     = "nonce"
	headerSignature = "signature"
	headerPublicKey = "public-key"
	headerAPIKey    = "api-key"
)

// Client manages communication with the KunaPay API.
type Client struct {
	// base URL for API requests.
	baseURL *url.URL

	// API key used to make authenticated API calls.
	apiKey string

	// Public key used to make authenticated API calls.
	publicKey string

	// Private key used to make authenticated API calls.
	privateKey []byte

	// User agent used when communicating with the API.
	userAgent string

	// HTTP client used to communicate with the API.
	httpClient *http.Client

	// Services used for talking to different parts of the KunaPay API.
	Asset       *AssetService
	Invoice     *InvoiceService
	Transaction *TransactionService
	Withdraw    *WithdrawService
}

// New returns a new KunaPay API client that uses signature authentication
// with the provided public and private keys.
func New(publicKey, privateKey string, opts ...ClientOptions) (*Client, error) {
	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("public and private keys are required")
	}

	client, err := newClient(opts...)
	if err != nil {
		return nil, err
	}

	client.publicKey = publicKey
	client.privateKey = []byte(privateKey)

	return client, err
}

// NewWithAPIKey returns a new KunaPay API client using the provided API key.
func NewWithAPIKey(apiKey string, opts ...ClientOptions) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	client, err := newClient(opts...)
	if err != nil {
		return nil, err
	}

	client.apiKey = apiKey

	return client, nil
}

type ClientOptions func(*Client) error

// newClient returns a new KunaPay API client instance.
func newClient(opts ...ClientOptions) (*Client, error) {
	baseURL, _ := url.Parse(apiURL)

	client := &Client{
		baseURL:   baseURL,
		userAgent: userAgent,
	}

	client.Asset = &AssetService{client: client}
	client.Invoice = &InvoiceService{client: client}
	client.Transaction = &TransactionService{client: client}
	client.Withdraw = &WithdrawService{client: client}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	if client.httpClient == nil {
		client.httpClient = http.DefaultClient
	}

	return client, nil
}

// WithHTTPClient sets the owning http.Client to use for API requests.
func WithHTTPClient(client *http.Client) ClientOptions {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

// SetUserAgent sets the custom user agent string to use when sending requests.
func SetUserAgent(userAgent string) ClientOptions {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in the path,
// it will be resolved in relation to the Client's baseURL. If specified,
// the value pointed to by body will be JSON encoded and included as the request body.
func (c *Client) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	rel, err := url.Parse(apiVersion + "/" + path)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if err = c.setAuth(req, body); err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) setAuth(req *http.Request, body any) error {
	if c.publicKey != "" && c.privateKey != nil {
		ts := fmt.Sprintf("%d", time.Now().UnixMilli())
		sign, err := c.sign(ts, req.URL.RequestURI(), body)
		if err != nil {
			return err
		}
		req.Header.Set(headerNonce, ts)
		req.Header.Set(headerSignature, sign)
		req.Header.Set(headerPublicKey, c.publicKey)
	} else if c.apiKey != "" {
		req.Header.Set(headerAPIKey, c.apiKey)
	}

	return nil
}

// Do sends an API request and returns the API response.
// The JSON response from the API is decoded and saved in the pointed value v.
// If there is an API error, an error response is returned instead.
func (c *Client) Do(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if err = handleErrorResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

// handleErrorResponse checks the API response for errors and returns
// them if they are found.
func handleErrorResponse(r *http.Response) error {
	if code := r.StatusCode; code >= http.StatusOK && code <= 299 {
		return nil
	}

	errResp := &ResponseError{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, errResp)
		if err != nil {
			errResp.Errors = append(errResp.Errors, Error{
				Code:    "",
				Message: string(data),
			})
		}
	}

	return errResp
}

// ResponseError represents an error response from the KunaPay API.
type ResponseError struct {
	// HTTP response that caused this error.
	Response *http.Response `json:"-"`

	Errors []Error `json:"errors"`
}

// Error represents a KunaPay API error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error returns the string representation of the error.
func (r *ResponseError) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Errors,
	)
}

// sign calculates the signature for the request using HMAC-SHA384 algorithm.
func (c *Client) sign(nonce, url string, body any) (string, error) {
	var reqBody = "{}"
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("sign calculation: %w", err)
		}
		reqBody = string(b)
	}

	hash := hmac.New(sha512.New384, c.privateKey)
	data := []byte(url + nonce + reqBody)
	hash.Write(data)

	return hex.EncodeToString(hash.Sum(nil)), nil
}
