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

	var gotPaths []string
	var gotAuth string
	var gotRequestID string
	var gotCorrelationID string
	var gotBody struct {
		Filters struct {
			OrderNumber string `json:"orderNumber"`
			OrderType   string `json:"orderType"`
		} `json:"filters"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		gotAuth = r.Header.Get("Authorization")
		gotRequestID = r.Header.Get(httpclient.HeaderRequestID)
		gotCorrelationID = r.Header.Get(httpclient.HeaderCorrelationID)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/findOrders":
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				t.Fatalf("Decode(request body) error = %v", err)
			}
			_, _ = w.Write([]byte(`{
				"data":[{
					"orderNumber":"A123",
					"createdDate":"2026-05-26T00:00:00Z",
					"customerName":"Customer One",
					"email":"customer@example.test",
					"orderTotal":980584,
					"statusOrderMax":"7500",
					"statusOrderMin":"7000"
				}],
				"total":1
			}`))
		case "/getOrder":
			_, _ = w.Write([]byte(`{"data":{"infClient":{"email":"customer@example.test"},"paymentInformation":{"orderTotal":980584}}}`))
		case "/findItemsByOrder":
			var itemBody struct {
				NotFood bool `json:"notFood"`
			}
			if err := json.NewDecoder(r.Body).Decode(&itemBody); err != nil {
				t.Fatalf("Decode(item body) error = %v", err)
			}
			if itemBody.NotFood {
				_, _ = w.Write([]byte(`{"data":[{"orderLineId":"1","itemId":"6450"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"data":[]}`))
			}
		default:
			http.NotFound(w, r)
		}
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

	if len(gotPaths) != 4 || gotPaths[0] != "/findOrders" || gotPaths[1] != "/getOrder" || gotPaths[2] != "/findItemsByOrder" || gotPaths[3] != "/findItemsByOrder" {
		t.Fatalf("request paths = %#v, want GEOMS summary, detail, and items calls", gotPaths)
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
	if gotBody.Filters.OrderNumber != "A123" || gotBody.Filters.OrderType != "ExitoEcomm" {
		t.Fatalf("request body = %#v, want GEOMS order filter", gotBody)
	}
	if got.ID != "A123" || got.Status != "7500" || got.CreatedAt != "2026-05-26T00:00:00Z" || got.CustomerName != "Customer One" || got.Email != "customer@example.test" || got.OrderTotal != 980584 {
		t.Fatalf("mapped order = %#v, want provider summary fields mapped", got)
	}
	if got.Details == nil || got.Items == nil || len(got.Items.NotFood) != 1 || len(got.Items.Food) != 0 {
		t.Fatalf("enriched order = %#v, want details and grouped items", got)
	}
}

func TestHTTPGetterFetchesGEOMSTokenAndReusesUntilExpiry(t *testing.T) {
	t.Parallel()

	var tokenRequests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			tokenRequests++
			if err := r.ParseForm(); err != nil {
				t.Fatalf("ParseForm() error = %v", err)
			}
			if r.Form.Get("scope") != "api://scope-value/.default" || r.Form.Get("grant_type") != "client_credentials" {
				t.Fatalf("token form = %#v, want Azure client credentials", r.Form)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"dynamic-token","token_type":"Bearer","expires_in":3599}`))
		case "/findOrders":
			if got := r.Header.Get("Authorization"); got != "Bearer dynamic-token" {
				t.Fatalf("Authorization = %q, want dynamic token", got)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[{"orderNumber":"A123","statusOrderMax":"7500","createdDate":"2026-05-26T00:00:00Z"}]}`))
		case "/getOrder":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{}}`))
		case "/findItemsByOrder":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	getter := orders.NewHTTPGetter(orders.HTTPGetterConfig{
		BaseURL:      server.URL,
		TokenURL:     server.URL + "/token",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scope:        "scope-value",
	}, server.Client())

	for range 2 {
		if _, err := getter.Get(context.Background(), orders.GetInput{ID: "A123"}); err != nil {
			t.Fatalf("Get() error = %v", err)
		}
	}
	if tokenRequests != 1 {
		t.Fatalf("token requests = %d, want cached single request", tokenRequests)
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
