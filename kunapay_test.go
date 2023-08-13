package kunapay

import (
	"bytes"
	"encoding/json"
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
	client, _ = NewWithAPIKey("API_key")
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

// func testHeader(t *testing.T, r *http.Request, header string, want string) {
// 	t.Helper()
// 	if got := r.Header.Get(header); got != want {
// 		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
// 	}
// }

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

func TestNewClientWithSignture(t *testing.T) {
	pubKey := "public_key"
	privKey := "private_key"

	c, _ := New(pubKey, privKey)
	if c.publicKey != pubKey {
		t.Errorf("Client publicKey is %v, want %v", c.publicKey, pubKey)
	}
	if c.privateKey != privKey {
		t.Errorf("Client privateKey is %v, want %v", c.privateKey, privKey)
	}
	if c.userAgent != userAgent {
		t.Errorf("Client userAgent is %v, want %v", c.userAgent, userAgent)
	}
	if c.baseURL.String() != apiURL {
		t.Errorf("Client baseURL is %v, want %v", c.baseURL.String(), apiURL)
	}
}

func TestNewRequest(t *testing.T) {
	t.Skip("TODO")
}

func TestDo(t *testing.T) {
	t.Skip("TODO")
}

func TestDo_httpBadRequest(t *testing.T) {
	t.Skip("TODO")
}

func TestDo_redirectLoop(t *testing.T) {
	t.Skip("TODO")
}

func TestCheckResponse(t *testing.T) {

	tests := []struct {
		title    string
		input    *http.Response
		expected *ErrorResponse
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
			expected: &ErrorResponse{
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
			expected: &ErrorResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := handleErrorResponse(test.input).(*ErrorResponse)
			if err == nil {
				t.Errorf("Expected error response.")
			}
			test.expected.Response = test.input

			if !reflect.DeepEqual(err, test.expected) {
				t.Errorf("Error = %#v, want %#v", err, test.expected)
			}
		})
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{
		Request: &http.Request{},
	}
	err := ErrorResponse{
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
