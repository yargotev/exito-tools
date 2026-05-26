package orders_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/domain/orders"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

func TestHTTPGetterPostsRequestAndMapsProviderResponse(t *testing.T) {
	t.Parallel()

	var gotPath string
	var gotAuth string
	var gotRequestID string
	var gotCorrelationID string
	var gotBody struct {
		ID string `json:"id"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotRequestID = r.Header.Get(httpclient.HeaderRequestID)
		gotCorrelationID = r.Header.Get(httpclient.HeaderCorrelationID)
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode(request body) error = %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"order":{
				"id":"A123",
				"status":"created",
				"createdAt":"2026-05-26T00:00:00Z"
			}
		}`))
	}))
	defer server.Close()

	ctx := httpclient.ContextWithRequestMetadata(context.Background(), httpclient.RequestMetadata{
		RequestID:     "req_orders",
		CorrelationID: "corr-orders",
	})
	getter := orders.NewHTTPGetter(orders.HTTPGetterConfig{BaseURL: server.URL, Token: "token-123"}, server.Client())
	got, err := getter.Get(ctx, orders.GetInput{ID: "A123"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if gotPath != "/orders/get" {
		t.Fatalf("request path = %q, want /orders/get", gotPath)
	}
	if gotAuth != "Bearer token-123" {
		t.Fatalf("Authorization = %q, want bearer token", gotAuth)
	}
	if gotRequestID != "req_orders" {
		t.Fatalf("%s = %q, want req_orders", httpclient.HeaderRequestID, gotRequestID)
	}
	if gotCorrelationID != "corr-orders" {
		t.Fatalf("%s = %q, want corr-orders", httpclient.HeaderCorrelationID, gotCorrelationID)
	}
	if gotBody.ID != "A123" {
		t.Fatalf("request body = %#v, want id A123", gotBody)
	}
	if got.ID != "A123" || got.Status != "created" || got.CreatedAt != "2026-05-26T00:00:00Z" {
		t.Fatalf("mapped order = %#v, want provider fields mapped", got)
	}
}

func TestHTTPGetterReturnsStructuredErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		getter   orders.HTTPGetter
		wantCode string
	}{
		{
			name:     "missing config",
			getter:   orders.NewHTTPGetter(orders.HTTPGetterConfig{}, nil),
			wantCode: orders.ErrorOrdersNotConfigured,
		},
		{
			name:     "invalid base URL",
			getter:   orders.NewHTTPGetter(orders.HTTPGetterConfig{BaseURL: "://bad", Token: "token"}, nil),
			wantCode: orders.ErrorOrdersNotConfigured,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := tc.getter.Get(context.Background(), orders.GetInput{ID: "A123"})
			assertStructuredCode(t, err, tc.wantCode)
		})
	}
}

func TestHTTPGetterHandlesProviderFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantCode string
	}{
		{
			name: "not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "missing", http.StatusNotFound)
			},
			wantCode: orders.ErrorOrderNotFound,
		},
		{
			name: "non success status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "provider down", http.StatusBadGateway)
			},
			wantCode: orders.ErrorOrdersProviderUnavailable,
		},
		{
			name: "invalid JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`not-json`))
			},
			wantCode: orders.ErrorOrdersProviderInvalidResponse,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tc.handler)
			defer server.Close()

			getter := orders.NewHTTPGetter(orders.HTTPGetterConfig{BaseURL: server.URL, Token: "token"}, server.Client())
			_, err := getter.Get(context.Background(), orders.GetInput{ID: "A123"})
			assertStructuredCode(t, err, tc.wantCode)
		})
	}
}

func assertStructuredCode(t *testing.T, err error, want string) {
	t.Helper()

	var structured capability.StructuredError
	if !errors.As(err, &structured) {
		t.Fatalf("error = %T, want StructuredError", err)
	}
	if structured.Code != want {
		t.Fatalf("StructuredError.Code = %q, want %q", structured.Code, want)
	}
}
