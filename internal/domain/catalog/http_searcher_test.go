package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

func TestHTTPSearcherSearchProductsBySKUID(t *testing.T) {
	t.Parallel()

	var gotRequestID string
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = r.Header.Get(httpclient.HeaderRequestID)
		gotQuery = r.URL.RawQuery
		if r.URL.Path != "/api/catalog_system/pub/products/search" {
			t.Fatalf("path = %q, want product search", r.URL.Path)
		}
		if r.URL.Query()["fq"][0] != "skuId:912350" {
			t.Fatalf("fq = %#v, want sku filter", r.URL.Query()["fq"])
		}
		if r.URL.Query().Get("_from") != "0" || r.URL.Query().Get("_to") != "0" {
			t.Fatalf("pagination query = %q", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("resources", "0-0/1")
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write([]byte(`[{
			"productId":"534690",
			"productName":"Minibar Abba",
			"brand":"ABBA",
			"brandId":6,
			"linkText":"nevera-minibar-97-litros-abba-534690",
			"productReference":"534690",
			"categoryId":"34185508",
			"link":"https://example.test/p",
			"items":[{"itemId":"912350","name":"Minibar","nameComplete":"Minibar complete","ean":"7706060050094","referenceId":[{"Key":"RefId","Value":"912350"}],"sellers":[{"sellerId":"VMIABBA","sellerName":"Exito","sellerDefault":true,"commertialOffer":{"Price":99900,"ListPrice":99900,"IsAvailable":true,"AvailableQuantity":7}}]}]
		}]`))
	}))
	defer server.Close()

	ctx := httpclient.ContextWithRequestMetadata(context.Background(), httpclient.RequestMetadata{RequestID: "req-catalog"})
	result, err := NewHTTPSearcher(HTTPSearcherConfig{BaseURL: server.URL}, server.Client()).SearchProducts(ctx, SearchProductsInput{Brand: "exito", By: "sku-id", Value: "912350", From: 0, To: 0})
	if err != nil {
		t.Fatalf("SearchProducts() error = %v", err)
	}
	if gotRequestID != "req-catalog" {
		t.Fatalf("request id = %q, want req-catalog", gotRequestID)
	}
	if gotQuery == "" {
		t.Fatalf("query should be sent")
	}
	if result.Count != 1 || result.Products[0].ProductID != "534690" || result.Products[0].Items[0].ItemID != "912350" {
		t.Fatalf("result = %#v, want mapped product/SKU", result)
	}
	if result.Total == nil || *result.Total != 1 || result.RangeStart == nil || *result.RangeStart != 0 {
		t.Fatalf("resources metadata not parsed: %#v", result)
	}
	if result.Products[0].Details["productId"] != "534690" {
		t.Fatalf("provider details not preserved: %#v", result.Products[0].Details)
	}
}

func TestHTTPSearcherSearchProductsSlugUsesProductURLPath(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/catalog_system/pub/products/search/nevera-minibar-97-litros-abba-534690/p" {
			t.Fatalf("path = %q, want slug path", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("query = %q, want empty", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	_, err := NewHTTPSearcher(HTTPSearcherConfig{BaseURL: server.URL}, server.Client()).SearchProducts(context.Background(), SearchProductsInput{Brand: "exito", By: "slug", Value: "nevera-minibar-97-litros-abba-534690", From: 0, To: 9})
	if err != nil {
		t.Fatalf("SearchProducts() error = %v", err)
	}
}

func TestSearchProductsUseCaseValidatesInput(t *testing.T) {
	t.Parallel()

	_, err := NewSearchProductsUseCase(UnavailableSearcher{}).Execute(context.Background(), SearchProductsInput{By: "sku-id"})
	if err == nil {
		t.Fatalf("Execute() error = nil, want validation error")
	}
}
