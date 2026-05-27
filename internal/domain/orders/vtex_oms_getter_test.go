package orders_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yargotev/exito-tools/internal/domain/orders"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

func TestVTEXOMSGetterSendsServerSideCredentialsAndMapsOrder(t *testing.T) {
	t.Parallel()

	var gotPath string
	var gotAppKey string
	var gotAppToken string
	var gotRequestID string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAppKey = r.Header.Get("X-VTEX-API-AppKey")
		gotAppToken = r.Header.Get("X-VTEX-API-AppToken")
		gotRequestID = r.Header.Get(httpclient.HeaderRequestID)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"orderId":           "1611511090420-01",
			"sequence":          "12345",
			"status":            "ready-for-handling",
			"statusDescription": "Ready for handling",
			"creationDate":      "2026-05-27T15:00:00Z",
			"value":             123456,
			"clientProfileData": map[string]any{"firstName": "Ada", "lastName": "Lovelace", "email": "ada@example.test"},
		})
	}))
	defer server.Close()

	ctx := httpclient.ContextWithRequestMetadata(context.Background(), httpclient.RequestMetadata{RequestID: "req-vtex"})
	getter := orders.NewVTEXOMSGetter(orders.VTEXOMSGetterConfig{BaseURL: server.URL, AppKey: "app-key", AppToken: "app-token"}, server.Client())
	got, err := getter.GetVTEXOMS(ctx, orders.GetVTEXOMSInput{ID: "1611511090420-01", Brand: "exito"})
	if err != nil {
		t.Fatalf("GetVTEXOMS() error = %v", err)
	}

	if gotPath != "/api/oms/pvt/orders/1611511090420-01" {
		t.Fatalf("request path = %q, want VTEX OMS order detail path", gotPath)
	}
	if gotAppKey != "app-key" || gotAppToken != "app-token" {
		t.Fatalf("VTEX credential headers = (%q,%q), want configured app key/token", gotAppKey, gotAppToken)
	}
	if gotRequestID != "req-vtex" {
		t.Fatalf("%s = %q, want req-vtex", httpclient.HeaderRequestID, gotRequestID)
	}
	if got.ID != "1611511090420-01" || got.Status != "ready-for-handling" || got.ClientName != "Ada Lovelace" || got.Email != "ada@example.test" || got.Brand != "exito" {
		t.Fatalf("mapped order = %#v", got)
	}
	if got.Details == nil || got.Details["orderId"] != "1611511090420-01" {
		t.Fatalf("details should preserve provider payload for diagnostics: %#v", got.Details)
	}
}
