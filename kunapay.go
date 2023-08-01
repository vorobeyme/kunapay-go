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
	"runtime"
	"time"
)

const (
	// libVersion is the current version of the library.
	libVersion = "0.0.1"

	// apiURL is the default base URL for the KunaPay API.
	apiURL = "https://api-kunapayapp.kuna.io/v1/"
)

var (
	// userAgent is the default user agent string to use when sending requests
	// to the KunaPay API.
	userAgent = fmt.Sprintf("kunapay-go/%s (%s %s) go/%s", libVersion, runtime.GOOS, runtime.GOARCH, runtime.Version())
)

// Client manages communication with the KunaPay API.
type Client struct {
	// base URL for API requests.
	baseURL *url.URL

	// API key used to make authenticated API calls.
	// apiKey string

	// Private key used to make authenticated API calls.
	privateKey string

	// Public key used to make authenticated API calls.
	publicKey string

	// User agent used when communicating with the API.
	userAgent string

	// HTTP client used to communicate with the API.
	client *http.Client

	// Services used for talking to different parts of the KunaPay API.
	Asset       *AssetService
	Invoice     *InvoiceService
	Transaction *TransactionService
	Withdraw    *WithdrawService
}

// NewClient returns a new KunaPay API client.
// If a nil httpClient is provided, http.DefaultClient will be used.
func New(publicKey, privateKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(apiURL)

	client := &Client{
		client:     httpClient,
		baseURL:    baseURL,
		privateKey: privateKey,
		publicKey:  publicKey,
		userAgent:  userAgent,
	}

	client.Asset = &AssetService{client: client}
	client.Invoice = &InvoiceService{client: client}
	client.Transaction = &TransactionService{client: client}
	client.Withdraw = &WithdrawService{client: client}

	return client
}

// RequestOption specifies the optional parameters to the Client.NewRequest method
// that can modify a http.Request.
type RequestOption func(req *http.Request)

// NewRequest creates an API request.
// A relative URL can be provided in the path, it will be resolved in relation to the Client's baseURL.
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	rel, err := url.Parse(path)
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

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("nonce", fmt.Sprint(time.Now().UnixNano()))
	req.Header.Set("public-key", c.publicKey)

	sign, err := c.Sign(req.Header.Get("nonce"), u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("signature", sign)

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

// Do sends an API request and returns the API response.
// The JSON response from the API is decoded and saved in the pointed value v.
// If there is an API error, an error response is returned instead.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

// CheckResponse checks the API response for errors and returns
// them if they are found.
func CheckResponse(r *http.Response) error {
	if code := r.StatusCode; code >= http.StatusOK && code <= 299 {
		return nil
	}

	errorResp := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, errorResp)
	}

	return errorResp
}

// ErrorResponse represents an error response from the KunaPay API.
type ErrorResponse struct {
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
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Errors,
	)
}

// Sign calculates the signature for the request using HMAC-SHA384 algorithm.
func (c *Client) Sign(nonce, url string, body interface{}) (string, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("%s%s%s", url, nonce, string(bodyBytes))

	hash := hmac.New(sha512.New384, []byte(c.privateKey))
	hash.Write([]byte(message))
	signature := hex.EncodeToString(hash.Sum(nil))

	return signature, err
}
