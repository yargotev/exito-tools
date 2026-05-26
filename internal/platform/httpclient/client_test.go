package httpclient_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

func TestClientNewJSONRequestAppliesSharedHeaders(t *testing.T) {
	t.Parallel()

	ctx := httpclient.ContextWithRequestMetadata(context.Background(), httpclient.RequestMetadata{
		RequestID:     "req_shared",
		CorrelationID: "corr-shared",
	})
	client := httpclient.New(httpclient.Config{BaseURL: "https://provider.test/api/", Token: "token-123"})

	request, err := client.NewJSONRequest(ctx, http.MethodPost, "/orders/get", map[string]string{"id": "A123"})
	if err != nil {
		t.Fatalf("NewJSONRequest() error = %v", err)
	}

	if request.URL.String() != "https://provider.test/api/orders/get" {
		t.Fatalf("URL = %q, want joined endpoint", request.URL.String())
	}
	if request.Header.Get("Accept") != "application/json" {
		t.Fatalf("Accept = %q, want application/json", request.Header.Get("Accept"))
	}
	if request.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", request.Header.Get("Content-Type"))
	}
	if request.Header.Get("Authorization") != "Bearer token-123" {
		t.Fatalf("Authorization = %q, want bearer token", request.Header.Get("Authorization"))
	}
	if request.Header.Get(httpclient.HeaderRequestID) != "req_shared" {
		t.Fatalf("%s = %q, want req_shared", httpclient.HeaderRequestID, request.Header.Get(httpclient.HeaderRequestID))
	}
	if request.Header.Get(httpclient.HeaderCorrelationID) != "corr-shared" {
		t.Fatalf("%s = %q, want corr-shared", httpclient.HeaderCorrelationID, request.Header.Get(httpclient.HeaderCorrelationID))
	}

	var body map[string]string
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		t.Fatalf("Decode(request body) error = %v", err)
	}
	if body["id"] != "A123" {
		t.Fatalf("body = %#v, want id A123", body)
	}
}

func TestClientEndpointRequiresAbsoluteBaseURL(t *testing.T) {
	t.Parallel()

	client := httpclient.New(httpclient.Config{BaseURL: "provider.test", Token: "token"})
	if _, err := client.Endpoint("/orders/get"); err == nil {
		t.Fatalf("Endpoint() error = nil, want invalid base URL error")
	}
}

func TestClientConfiguredRequiresBaseURLAndToken(t *testing.T) {
	t.Parallel()

	if !httpclient.New(httpclient.Config{BaseURL: "https://provider.test", Token: "token"}).Configured() {
		t.Fatalf("Configured() = false, want true")
	}
	if httpclient.New(httpclient.Config{BaseURL: "https://provider.test"}).Configured() {
		t.Fatalf("Configured() = true without token, want false")
	}
	if httpclient.New(httpclient.Config{Token: "token"}).Configured() {
		t.Fatalf("Configured() = true without base URL, want false")
	}
}

func TestClientDoUsesConfiguredTimeout(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := httpclient.New(httpclient.Config{BaseURL: server.URL, Token: "token", Timeout: time.Millisecond})
	request, err := client.NewRequest(context.Background(), http.MethodGet, "/slow", nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}

	response, err := client.Do(request)
	if response != nil {
		_ = response.Body.Close()
	}
	if err == nil {
		t.Fatalf("Do() error = nil, want timeout")
	}
}

func TestResponseHelpers(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	response, err := server.Client().Get(server.URL)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer func() { _ = response.Body.Close() }()

	if !httpclient.Successful(response) {
		t.Fatalf("Successful() = false, want true")
	}
	var body struct {
		OK bool `json:"ok"`
	}
	if err := httpclient.DecodeJSONResponse(response, &body); err != nil {
		t.Fatalf("DecodeJSONResponse() error = %v", err)
	}
	if !body.OK {
		t.Fatalf("decoded body = %#v, want ok true", body)
	}
}

func TestSuccessfulRejectsNilAndNon2xxResponses(t *testing.T) {
	t.Parallel()

	if httpclient.Successful(nil) {
		t.Fatalf("Successful(nil) = true, want false")
	}
	if httpclient.Successful(&http.Response{StatusCode: http.StatusBadGateway}) {
		t.Fatalf("Successful(502) = true, want false")
	}
}

func TestDecodeJSONResponseReturnsDecodeError(t *testing.T) {
	t.Parallel()

	response := &http.Response{Body: errReaderCloser{}}
	var target map[string]string
	if err := httpclient.DecodeJSONResponse(response, &target); err == nil {
		t.Fatalf("DecodeJSONResponse() error = nil, want decode error")
	}
}

type errReaderCloser struct{}

func (errReaderCloser) Read([]byte) (int, error) { return 0, errors.New("read failed") }
func (errReaderCloser) Close() error             { return nil }
