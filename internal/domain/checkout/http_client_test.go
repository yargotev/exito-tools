package checkout_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yargotev/exito-tools/internal/domain/checkout"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

func TestHTTPClientGetOrderFormMapsRedactedSummary(t *testing.T) {
	t.Parallel()

	var gotRequestID string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = r.Header.Get(httpclient.HeaderRequestID)
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/api/checkout/pub/orderForm/of-123" {
			t.Fatalf("path = %q, want orderForm path", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"orderFormId":"of-123",
			"salesChannel":"1",
			"value":12345,
			"clientProfileData":{"email":"secret@example.com","document":"123"},
			"shippingData":{"address":{"receiverName":"Jane"}},
			"totalizers":[{"id":"Items","name":"Items Total","value":10000}],
			"items":[{"id":"sku-1","productId":"prod-1","name":"Safe Item Name","quantity":2,"seller":"1","price":5000,"sellingPrice":4500,"availability":"available"}]
		}`))
	}))
	defer server.Close()

	client := checkout.NewHTTPClient(checkout.HTTPClientConfig{BaseURL: server.URL}, server.Client())
	ctx := httpclient.ContextWithRequestMetadata(context.Background(), httpclient.RequestMetadata{RequestID: "req-checkout"})
	got, err := client.GetOrderForm(ctx, checkout.GetOrderFormInput{Brand: "exito", OrderFormID: "of-123"})
	if err != nil {
		t.Fatalf("GetOrderForm() error = %v", err)
	}
	if gotRequestID != "req-checkout" {
		t.Fatalf("request ID header = %q, want req-checkout", gotRequestID)
	}
	if got.ID != "of-123" || got.Brand != "exito" || got.Value != 12345 || got.ItemCount != 1 {
		t.Fatalf("summary = %#v, want mapped orderForm", got)
	}
	if !got.ClientProfileDataSet || !got.ShippingDataSet {
		t.Fatalf("presence flags = client:%v shipping:%v, want true", got.ClientProfileDataSet, got.ShippingDataSet)
	}
	if got.Items[0].ID != "sku-1" || got.Items[0].Quantity != 2 || got.Items[0].Seller != "1" {
		t.Fatalf("item summary = %#v, want mapped safe item fields", got.Items[0])
	}
	if got.Diagnostics.RequestPath != "/api/checkout/pub/orderForm/of-123" || got.Diagnostics.ProviderStatus != http.StatusOK {
		t.Fatalf("diagnostics = %#v, want safe request path/status", got.Diagnostics)
	}
}

func TestHTTPClientCreateOrderFormUsesForceNewCartAndSalesChannel(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/checkout/pub/orderForm" {
			t.Fatalf("path = %q, want create path", r.URL.Path)
		}
		if r.URL.Query().Get("forceNewCart") != "true" || r.URL.Query().Get("sc") != "2" {
			t.Fatalf("query = %q, want forceNewCart=true&sc=2", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"orderFormId":"new-of","salesChannel":"2","value":0,"items":[],"totalizers":[]}`))
	}))
	defer server.Close()

	client := checkout.NewHTTPClient(checkout.HTTPClientConfig{BaseURL: server.URL}, server.Client())
	got, err := client.CreateOrderForm(context.Background(), checkout.CreateOrderFormInput{Brand: "carulla", SalesChannel: "2"})
	if err != nil {
		t.Fatalf("CreateOrderForm() error = %v", err)
	}
	if got.ID != "new-of" || got.Brand != "carulla" || got.SalesChannel != "2" {
		t.Fatalf("summary = %#v, want created orderForm", got)
	}
}

func TestHTTPClientAddItemsPostsOrderItems(t *testing.T) {
	t.Parallel()

	var gotRequestID string
	var gotBody struct {
		OrderItems []struct {
			ID       string `json:"id"`
			Quantity int    `json:"quantity"`
			Seller   string `json:"seller"`
			Index    int    `json:"index"`
		} `json:"orderItems"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = r.Header.Get(httpclient.HeaderRequestID)
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/api/checkout/pub/orderForm/of-123/items" {
			t.Fatalf("path = %q, want add-items path", r.URL.Path)
		}
		if r.URL.Query().Get("allowedOutdatedData") != "false" {
			t.Fatalf("query = %q, want allowedOutdatedData=false", r.URL.RawQuery)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode request body error = %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"orderFormId":"of-123",
			"salesChannel":"1",
			"value":9990,
			"items":[{"id":"sku-1","productId":"prod-1","name":"Safe Item","quantity":2,"seller":"1","price":5000,"sellingPrice":4995,"availability":"available"}],
			"totalizers":[{"id":"Items","value":9990}]
		}`))
	}))
	defer server.Close()

	client := checkout.NewHTTPClient(checkout.HTTPClientConfig{BaseURL: server.URL}, server.Client())
	ctx := httpclient.ContextWithRequestMetadata(context.Background(), httpclient.RequestMetadata{RequestID: "req-add-items"})
	got, err := client.AddItems(ctx, checkout.AddItemsInput{
		Brand:       "exito",
		OrderFormID: "of-123",
		Items:       []checkout.AddItemInput{{SKU: "sku-1", Quantity: 2, Seller: "1"}},
	})
	if err != nil {
		t.Fatalf("AddItems() error = %v", err)
	}
	if gotRequestID != "req-add-items" {
		t.Fatalf("request ID header = %q, want req-add-items", gotRequestID)
	}
	if len(gotBody.OrderItems) != 1 || gotBody.OrderItems[0].ID != "sku-1" || gotBody.OrderItems[0].Quantity != 2 || gotBody.OrderItems[0].Index != 0 {
		t.Fatalf("request body = %#v, want mapped orderItems", gotBody)
	}
	if got.ID != "of-123" || got.ItemCount != 1 || got.Items[0].ID != "sku-1" {
		t.Fatalf("summary = %#v, want updated orderForm", got)
	}
	if got.Diagnostics.RequestPath != "/api/checkout/pub/orderForm/of-123/items?allowedOutdatedData=false" {
		t.Fatalf("diagnostics = %#v, want add-items request path", got.Diagnostics)
	}
}
