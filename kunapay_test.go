package kunapay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func setupClient() (client *Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	client, _ = New("public_key", "private_key")
	url, _ := url.Parse(server.URL)
	client.baseURL = url

	return client, mux, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %s, want %s", got, want)
	}
}

func testURL(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.RequestURI; got != want {
		t.Errorf("Request URL: %s, want %s", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		t.Fatalf("Failed to Read Body: %v", err)
	}
	if got := buf.String(); got != want {
		t.Errorf("Request body: %s, want %s", got, want)
	}
}

func testBadPathParams(t *testing.T, method string, fn func() error) {
	t.Helper()
	if err := fn(); err == nil {
		t.Errorf("%v bad params, err = nil, want error", method)
	}
}

func testNewRequestAndDoFailure(t *testing.T, method string, client *Client, f func() (*Response, error)) {
	t.Helper()

	client.baseURL.Scheme = ""
	resp, err := f()
	if resp != nil {
		t.Errorf("client.baseURL.Path='' %v resp = %#v, want nil", method, resp)
	}
	if err == nil {
		t.Errorf("client.baseURL.Path='' %v err = nil, want error", method)
	}
}

func testJSONMarshal(t *testing.T, v interface{}, want string) {
	t.Helper()

	u := reflect.New(reflect.TypeOf(v)).Interface()
	if err := json.Unmarshal([]byte(want), &u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v: %v", want, err)
	}
	w, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", u)
	}

	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", v)
	}

	if string(w) != string(j) {
		t.Errorf("json.Marshal(%q) \nreturned %s,\nwant %s", v, j, w)
	}
}

func TestNew(t *testing.T) {
	pubKey := "public_key"
	privKey := "private_key"
	c, _ := New(pubKey, privKey)
	if c.publicKey != pubKey {
		t.Errorf("Client publicKey is %v, want %v", c.publicKey, pubKey)
	}
	if string(c.privateKey) != privKey {
		t.Errorf("Client privateKey is %v, want %v", string(c.privateKey), privKey)
	}
	if c.userAgent != userAgent {
		t.Errorf("Client userAgent is %v, want %v", c.userAgent, userAgent)
	}
	if c.baseURL.String() != apiURL {
		t.Errorf("Client baseURL is %v, want %v", c.baseURL.String(), apiURL)
	}
}

func TestNew_emptyKeys(t *testing.T) {
	if _, err := New("", ""); err == nil {
		t.Errorf("New() with empty keys returned nil, want error")
	}
}

func TestNew_setBadBaseURL(t *testing.T) {
	apiURL := "bad\nURL"
	if _, err := New("public_key", "private_key", SetBaseURL(apiURL)); err == nil {
		t.Errorf("New() returned nil, want error")
	}
}

func TestNewWithAPIKey(t *testing.T) {
	apiKey := "api_key"
	apiURL := "http://127.0.0.1"
	ua := "test/0.0.1"
	c, _ := NewWithAPIKey(apiKey, WithHTTPClient(&http.Client{}), SetUserAgent(ua), SetBaseURL(apiURL))
	if c.apiKey != apiKey {
		t.Errorf("Client API key is %v, want %v", c.apiKey, apiKey)
	}
	if c.userAgent != ua {
		t.Errorf("Client userAgent is %v, want %v", c.userAgent, ua)
	}
	if c.baseURL.String() != apiURL {
		t.Errorf("Client baseURL is %v, want %v", c.baseURL.String(), apiURL)
	}
}

func TestNewWithAPIKey_emptyAPIKey(t *testing.T) {
	if _, err := NewWithAPIKey(""); err == nil {
		t.Errorf("NewWithAPIKey with empty API key returned nil, want error")
	}
}

func TestNewWithAPIKey_setBadBaseURL(t *testing.T) {
	apiURL := "bad\nURL"
	if _, err := NewWithAPIKey("api_key", SetBaseURL(apiURL)); err == nil {
		t.Errorf("NewWithAPIKey() returned nil, want error")
	}
}

func TestNewRequest(t *testing.T) {
	c, _ := NewWithAPIKey("api_key")

	var (
		inURL  = "withdraw"
		outURL = apiURL + apiVersion + "/withdraw"

		inBody = &CreateWithdrawRequest{
			Amount:        "100",
			Asset:         "USDT",
			PaymentMethod: "USDT",
		}
		outBody = `{"amount":"100","asset":"USDT","paymentMethod":"USDT"}` + "\n"

		ctx = context.Background()
	)

	req, _ := c.NewRequest(ctx, "GET", inURL, inBody)

	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	body, _ := io.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%v) Body is %v, want %v", inBody, got, want)
	}

	userAgent := req.Header.Get("User-Agent")
	if got, want := userAgent, c.userAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}

	if !strings.Contains(userAgent, libVersion) {
		t.Errorf("NewRequest() User-Agent should contain %v, found %v", libVersion, userAgent)
	}
}

func TestNewRequest_invalidBody(t *testing.T) {
	c, _ := NewWithAPIKey("api_key")

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.NewRequest(context.Background(), "GET", "", &T{})

	if err == nil {
		t.Fatal("NewRequest returned nil; expected error")
	}

	var unsupportedTypeError *json.UnsupportedTypeError
	if !errors.As(err, &unsupportedTypeError) {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c, _ := NewWithAPIKey("api_key")
	_, err := c.NewRequest(context.Background(), "GET", "\n", nil)
	if err == nil {
		t.Fatal("NewRequest returned nil; expected error")
	}

	var urlParseError *url.Error
	if !errors.As(err, &urlParseError) || urlParseError.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestNewRequest_badMethod(t *testing.T) {
	c, _ := NewWithAPIKey("api_key")
	if _, err := c.NewRequest(context.Background(), "\nMETHOD", "/", nil); err == nil {
		t.Fatal("NewRequest returned nil; expected error")
	}
}

func TestDo(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	type test struct {
		Foo string `json:"foo"`
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"foo":"bar"}`)
	})

	req, _ := client.NewRequest(context.Background(), "GET", ".", nil)
	body := &test{}
	_, _ = client.Do(req, body)

	want := &test{"bar"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpBadRequest(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest(context.Background(), "GET", "/", nil)
	resp, err := client.Do(req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, no error returned.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

func TestCheckResponse(t *testing.T) {
	tests := []struct {
		title    string
		input    *http.Response
		expected *ResponseError
	}{
		{
			title: "400 Bad Request",
			input: &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`{
					"errors": [{
						"code": "BAD_REQUEST",
						"message": "Field must be a valid type"
					}]
				}`)),
			},
			expected: &ResponseError{
				Response: &http.Response{},
				Errors: []Error{
					{
						Code:    "BAD_REQUEST",
						Message: "Field must be a valid type",
					},
				},
			},
		},
		{
			title: "no body",
			input: &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			expected: &ResponseError{},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := handleErrorResponse(test.input)
			if err == nil {
				t.Errorf("Expected error response.")
			}
			test.expected.Response = test.input

			if !errors.As(err, &test.expected) {
				t.Errorf("Error = %#v, want %#v", err, test.expected)
			}
		})
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{
		Request: &http.Request{},
	}
	err := ResponseError{
		Response: res,
		Errors: []Error{
			{
				Code:    "BAD_REQUEST",
				Message: "Field must be a valid type",
			},
		},
	}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}
