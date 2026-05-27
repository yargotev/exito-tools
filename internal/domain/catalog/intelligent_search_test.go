package catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIntelligentSearchProductsCapabilityBuildsTypedMultiIDQuery(t *testing.T) {
	searcher := recordingIntelligentSearcher{}
	useCase := NewIntelligentSearchProductsUseCase(searcher)

	result, err := useCase.Execute(context.Background(), IntelligentSearchProductsInput{
		Brand:       "carulla",
		TradePolicy: "1",
		By:          "sku-id",
		Values:      []string{"123", "456"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Query != "sku.id:123;456" {
		t.Fatalf("query = %q, want typed multi-ID query", result.Query)
	}
}

func TestIntelligentSearchProductsRejectsAmbiguousQueryModes(t *testing.T) {
	useCase := NewIntelligentSearchProductsUseCase(recordingIntelligentSearcher{})

	_, err := useCase.Execute(context.Background(), IntelligentSearchProductsInput{
		TradePolicy: "1",
		Text:        "leche",
		Query:       "sku.id:123",
	})
	if err == nil || !strings.Contains(err.Error(), "only one query mode") {
		t.Fatalf("Execute() error = %v, want ambiguous query mode error", err)
	}
}

func TestHTTPIntelligentSearcherBuildsRequestAndRedactsCookies(t *testing.T) {
	var gotPath string
	var gotQuery string
	var gotCookie string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		gotQuery = r.URL.RawQuery
		gotCookie = r.Header.Get("Cookie")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"products":[{"productId":"534690","productName":"Milk","items":[{"itemId":"912350","name":"Milk SKU","ean":"770"}]}],"total":1}`))
	}))
	defer server.Close()

	searcher := NewHTTPIntelligentSearcher(HTTPIntelligentSearcherConfig{BaseURL: server.URL}, server.Client())
	result, err := searcher.IntelligentSearchProducts(context.Background(), IntelligentSearchProductsInput{
		Brand:              "exito",
		TradePolicy:        "1",
		By:                 "sku-id",
		Values:             []string{"912350"},
		Facets:             []string{"category-1=lacteos", "color=not:red"},
		Page:               2,
		Count:              12,
		Sort:               "price:asc",
		HideUnavailable:    boolPtr(true),
		SimulationBehavior: "skip",
		Cookies:            []string{"vtex_segment=secret-segment"},
	})
	if err != nil {
		t.Fatalf("IntelligentSearchProducts() error = %v", err)
	}
	if gotPath != "/api/io/_v/api/intelligent-search/product_search/trade-policy/1/category-1/lacteos/color/not:red" {
		t.Fatalf("path = %q", gotPath)
	}
	for _, want := range []string{"query=sku.id%3A912350", "page=2", "count=12", "sort=price%3Aasc", "hideUnavailableItems=true", "simulationBehavior=skip"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %s", gotQuery, want)
		}
	}
	if gotCookie != "vtex_segment=secret-segment" {
		t.Fatalf("cookie header = %q", gotCookie)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), "secret-segment") {
		t.Fatalf("result leaked cookie: %s", string(encoded))
	}
	if result.Diagnostics.CookieNames[0] != "vtex_segment" || len(result.Products) != 1 || result.Products[0].ProductID != "534690" {
		t.Fatalf("result = %#v", result)
	}
}

type recordingIntelligentSearcher struct{}

func (recordingIntelligentSearcher) IntelligentSearchProducts(_ context.Context, input IntelligentSearchProductsInput) (IntelligentSearchProductsResult, error) {
	return IntelligentSearchProductsResult{Brand: input.Brand, Query: intelligentQuery(input)}, nil
}

func boolPtr(value bool) *bool {
	return &value
}

func TestCreateVTEXSegmentDefinitionRequiresConfirmation(t *testing.T) {
	definition := CreateVTEXSegmentDefinition()
	if definition.ID != CapabilityCreateVTEXSegmentID {
		t.Fatalf("ID = %q, want %s", definition.ID, CapabilityCreateVTEXSegmentID)
	}
	if definition.Risk != "safe-write" || !definition.RequiresConfirmation {
		t.Fatalf("definition risk/confirmation = %s/%v, want safe-write confirmation-required", definition.Risk, definition.RequiresConfirmation)
	}
}

func TestCreateVTEXSegmentRejectsMissingInputs(t *testing.T) {
	useCase := NewCreateVTEXSegmentUseCase(recordingSegmentCreator{})

	_, err := useCase.Execute(context.Background(), CreateVTEXSegmentInput{RegionID: "REGION"})
	if err == nil || !strings.Contains(err.Error(), "salesChannel is required") {
		t.Fatalf("Execute() error = %v, want missing salesChannel error", err)
	}
}

func TestHTTPVTEXSegmentCreatorPostsSessionAndRedactsToken(t *testing.T) {
	var gotPath string
	var gotMethod string
	var gotPayload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("Decode request body error = %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"segmentToken":"secret-token","expiresIn":86400,"namespaces":{"segment":{"token":"nested-secret"}}}`))
	}))
	defer server.Close()

	creator := NewHTTPVTEXSegmentCreator(HTTPVTEXSegmentCreatorConfig{BaseURL: server.URL}, server.Client())
	result, err := creator.CreateVTEXSegment(context.Background(), CreateVTEXSegmentInput{
		Brand:         "exito",
		RegionID:      "REGION_ID",
		SalesChannel:  "1",
		IncludeCookie: true,
	})
	if err != nil {
		t.Fatalf("CreateVTEXSegment() error = %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/io/api/sessions" {
		t.Fatalf("request = %s %s, want POST /io/api/sessions", gotMethod, gotPath)
	}
	public, ok := gotPayload["public"].(map[string]any)
	if !ok {
		t.Fatalf("payload public = %#v", gotPayload["public"])
	}
	region, ok := public["regionId"].(map[string]any)
	if !ok || region["value"] != "REGION_ID" {
		t.Fatalf("payload regionId = %#v", public["regionId"])
	}
	sc, ok := public["sc"].(map[string]any)
	if !ok || sc["value"] != "1" {
		t.Fatalf("payload sc = %#v", public["sc"])
	}
	if !result.TokenSet || result.TokenLength != len("secret-token") || result.Cookie != "vtex_segment=secret-token" {
		t.Fatalf("result token metadata = %#v", result)
	}
	encoded, err := json.Marshal(result.Diagnostics)
	if err != nil {
		t.Fatalf("Marshal diagnostics error = %v", err)
	}
	for _, leaked := range []string{"secret-token", "nested-secret"} {
		if strings.Contains(string(encoded), leaked) {
			t.Fatalf("diagnostics leaked token %q: %s", leaked, string(encoded))
		}
	}
	if result.Diagnostics.ProviderPayload["segmentToken"] != redactedValue {
		t.Fatalf("redacted segmentToken = %#v", result.Diagnostics.ProviderPayload["segmentToken"])
	}
}

type recordingSegmentCreator struct{}

func (recordingSegmentCreator) CreateVTEXSegment(_ context.Context, input CreateVTEXSegmentInput) (CreateVTEXSegmentResult, error) {
	return CreateVTEXSegmentResult{Brand: input.Brand, RegionID: input.RegionID, SalesChannel: input.SalesChannel}, nil
}
