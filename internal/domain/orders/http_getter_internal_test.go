package orders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
)

func TestGEOMSOrderMappingHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     geomsFindOrdersResponse
		wantID  string
		wantOK  bool
		wantSum float64
	}{
		{
			name: "data array",
			raw: geomsFindOrdersResponse{
				"data": []any{map[string]any{
					"orderNumber": "A123",
					"orderTotal":  json.Number("42.5"),
				}},
			},
			wantID:  "A123",
			wantOK:  true,
			wantSum: 42.5,
		},
		{
			name: "nested results",
			raw: geomsFindOrdersResponse{
				"data": map[string]any{"results": []any{map[string]any{"orderID": float64(12345)}}},
			},
			wantID: "12345",
			wantOK: true,
		},
		{
			name:    "top-level order",
			raw:     geomsFindOrdersResponse{"order": map[string]any{"id": "B456", "orderTotal": 12}},
			wantID:  "B456",
			wantOK:  true,
			wantSum: 12,
		},
		{
			name:   "empty data",
			raw:    geomsFindOrdersResponse{"data": []any{}},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, ok := tc.raw.firstOrder()
			if ok != tc.wantOK {
				t.Fatalf("firstOrder() ok = %v, want %v", ok, tc.wantOK)
			}
			if got.ID != tc.wantID || got.OrderTotal != tc.wantSum {
				t.Fatalf("firstOrder() = %#v, want ID %q total %v", got, tc.wantID, tc.wantSum)
			}
		})
	}
}

func TestGEOMSRequestHelpersNormalizeOrderValues(t *testing.T) {
	t.Parallel()

	findOrders := geomsFindOrdersRequest(GetInput{ID: "  1611511090420-01  ", OrderType: "  CarullaEcomm  "})
	if findOrders.Filters.OrderNumber == nil || *findOrders.Filters.OrderNumber != "1611511090420" {
		t.Fatalf("findOrders order number = %#v, want clean order number", findOrders.Filters.OrderNumber)
	}
	if findOrders.Filters.OrderType != "CarullaEcomm" {
		t.Fatalf("findOrders order type = %q, want CarullaEcomm", findOrders.Filters.OrderType)
	}

	getOrder := geomsGetOrderRequest("  1611511090420-01  ")
	if getOrder.Order != "1611511090420" {
		t.Fatalf("getOrder order = %q, want clean order number", getOrder.Order)
	}

	items := geomsFindItemsByOrderRequest("  1611511090420-01  ", true)
	if items.Order != "1611511090420" || !items.NotFood || items.PerPageItem != 999 || items.PageNumberItem != 1 {
		t.Fatalf("items request = %#v, want cleaned not-food request", items)
	}
}

func TestStructuredCode(t *testing.T) {
	t.Parallel()

	structured := capability.StructuredError{Code: ErrorOrderNotFound, Message: "Order not found."}
	if got := structuredCode(structured); got != ErrorOrderNotFound {
		t.Fatalf("structuredCode() = %q, want %s", got, ErrorOrderNotFound)
	}
	if got := structuredCode(errors.New("plain")); got != "" {
		t.Fatalf("structuredCode(plain) = %q, want empty", got)
	}
}

func TestHTTPGetterTreatsMissingDetailsAndItemsAsEmpty(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/findOrders":
			_, _ = w.Write([]byte(`{"data":[{"orderNumber":"A123","statusOrderMax":"7500","createdDate":"2026-05-26T00:00:00Z"}]}`))
		case "/getOrder":
			_, _ = w.Write([]byte(`{"data":null}`))
		case "/findItemsByOrder":
			http.NotFound(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	getter := NewHTTPGetter(HTTPGetterConfig{BaseURL: server.URL, Token: "token-123"}, server.Client())
	got, err := getter.Get(context.Background(), GetInput{ID: "A123"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got.Details) != 0 || got.Items == nil || len(got.Items.Food) != 0 || len(got.Items.NotFood) != 0 {
		t.Fatalf("Get() = %#v, want empty details and item groups", got)
	}
}

func TestHTTPGetterPropagatesDetailAndItemErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		detailCode  int
		itemCode    int
		wantMessage string
	}{
		{name: "detail error", detailCode: http.StatusBadGateway, itemCode: http.StatusOK, wantMessage: "Orders provider returned an unsuccessful response."},
		{name: "item error", detailCode: http.StatusOK, itemCode: http.StatusBadGateway, wantMessage: "Orders provider returned an unsuccessful response."},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/findOrders":
					_, _ = w.Write([]byte(`{"data":[{"orderNumber":"A123","statusOrderMax":"7500","createdDate":"2026-05-26T00:00:00Z"}]}`))
				case "/getOrder":
					w.WriteHeader(tc.detailCode)
					if tc.detailCode == http.StatusOK {
						_, _ = w.Write([]byte(`{"data":{}}`))
					}
				case "/findItemsByOrder":
					w.WriteHeader(tc.itemCode)
					if tc.itemCode == http.StatusOK {
						_, _ = w.Write([]byte(`{"data":[]}`))
					}
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()

			getter := NewHTTPGetter(HTTPGetterConfig{BaseURL: server.URL, Token: "token-123"}, server.Client())
			_, err := getter.Get(context.Background(), GetInput{ID: "A123"})

			var structured capability.StructuredError
			if !errors.As(err, &structured) {
				t.Fatalf("Get() error = %T, want StructuredError", err)
			}
			if structured.Code != ErrorOrdersProviderUnavailable || structured.Message != tc.wantMessage {
				t.Fatalf("StructuredError = %#v, want provider unavailable", structured)
			}
		})
	}
}

func TestGEOMSTokenSourceErrorPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tokenURL string
		handler  http.HandlerFunc
		wantCode string
	}{
		{
			name:     "invalid token URL",
			tokenURL: "://bad",
			handler:  nil,
			wantCode: ErrorOrdersNotConfigured,
		},
		{
			name:     "unsuccessful token endpoint",
			handler:  func(w http.ResponseWriter, r *http.Request) { http.Error(w, "down", http.StatusBadGateway) },
			wantCode: ErrorOrdersProviderUnavailable,
		},
		{
			name:     "invalid token response",
			handler:  func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`{"access_token":""}`)) },
			wantCode: ErrorOrdersProviderInvalidResponse,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := http.DefaultClient
			tokenURL := tc.tokenURL
			if tc.handler != nil {
				server := httptest.NewServer(tc.handler)
				defer server.Close()
				client = server.Client()
				tokenURL = server.URL
			}

			source := newGEOMSTokenSource(geomsTokenConfig{
				TokenURL:     tokenURL,
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				Scope:        "scope-value",
			}, client)
			_, err := source.token(context.Background())

			var structured capability.StructuredError
			if !errors.As(err, &structured) {
				t.Fatalf("token() error = %T, want StructuredError", err)
			}
			if structured.Code != tc.wantCode {
				t.Fatalf("StructuredError.Code = %q, want %s", structured.Code, tc.wantCode)
			}
		})
	}
}
