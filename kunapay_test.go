package kunapay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func setupClient() (client *Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client = New("public_key", "private_key", nil)
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

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
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

func TestNewClient(t *testing.T) {
	pubKey := "public_key"
	privKey := "private_key"

	c := New(pubKey, privKey, nil)
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
