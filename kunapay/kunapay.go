package kunapay

import (
	"bytes"
	"context"
	"crypto/sha512"
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

// NewRequest creates an API request.
// A relative URL can be provided in the path, in which case it is resolved relative
// to the baseURL of the Client.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
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
	req.Header.Set("signature", c.Sign(req.Header.Get("nonce"), u.String(), buf.(*bytes.Buffer).String()))
	req.Header.Set("public-key", c.publicKey)

	return req, nil
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v,
// or returned as an error if an API error has occurred.
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

	return resp, err
}

// Sign calculates the signature for the request.
// TODO: implement
func (c *Client) Sign(nonce, url string, body string) string {
	return ""
}
