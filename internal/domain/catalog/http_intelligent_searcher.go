package catalog

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

const maxIntelligentSearchResponseBodyBytes = 16 << 20

// HTTPIntelligentSearcherConfig contains the provider settings required by HTTPIntelligentSearcher.
type HTTPIntelligentSearcherConfig struct {
	BaseURL string
}

// HTTPIntelligentSearcher calls VTEX Intelligent Search and maps its DTO to domain output.
type HTTPIntelligentSearcher struct {
	baseURL string
	client  httpclient.Client
}

// NewHTTPIntelligentSearcher creates a VTEX Intelligent Search-backed product searcher.
func NewHTTPIntelligentSearcher(config HTTPIntelligentSearcherConfig, client *http.Client) HTTPIntelligentSearcher {
	return HTTPIntelligentSearcher{
		baseURL: strings.TrimSpace(config.BaseURL),
		client:  httpclient.New(httpclient.Config{BaseURL: config.BaseURL, Client: client}),
	}
}

// IntelligentSearchProducts calls VTEX Intelligent Search product_search.
func (s HTTPIntelligentSearcher) IntelligentSearchProducts(ctx context.Context, input IntelligentSearchProductsInput) (IntelligentSearchProductsResult, error) {
	if strings.TrimSpace(s.baseURL) == "" {
		return IntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Intelligent Search client is not configured."}
	}

	path, query, facets, err := intelligentPathAndQuery(input)
	if err != nil {
		return IntelligentSearchProductsResult{}, err
	}
	request, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return IntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Intelligent Search provider base URL is invalid."}
	}
	request.URL.RawQuery = query.Encode()
	if cookieHeader := cookieHeader(input.Cookies); cookieHeader != "" {
		request.Header.Set("Cookie", cookieHeader)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return IntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogProviderUnavailable, Message: "VTEX Intelligent Search provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()

	if !httpclient.Successful(response) {
		return IntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogProviderUnavailable, Message: "VTEX Intelligent Search provider returned an unsuccessful response."}
	}

	var providerResponse map[string]any
	limited := io.LimitReader(response.Body, maxIntelligentSearchResponseBodyBytes)
	if err := json.NewDecoder(limited).Decode(&providerResponse); err != nil {
		return IntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogProviderInvalidResponse, Message: "VTEX Intelligent Search provider returned an invalid response."}
	}

	productsPayload := intelligentProductsPayload(providerResponse)
	products := make([]Product, 0, len(productsPayload))
	for _, fields := range productsPayload {
		products = append(products, productFromMap(fields))
	}
	requestQuery := map[string]string{}
	for key, values := range query {
		if len(values) > 0 {
			requestQuery[key] = values[0]
		}
	}
	return IntelligentSearchProductsResult{
		Brand:    normalizedBrand(input.Brand),
		Query:    intelligentQuery(input),
		Facets:   facets,
		Page:     input.Page,
		Count:    input.Count,
		Sort:     sortOrRelevance(input.Sort),
		Products: products,
		Total:    totalFromIntelligentResponse(providerResponse),
		Diagnostics: RequestDiagnostics{
			RequestPath:      path,
			RequestQuery:     requestQuery,
			CookieNames:      cookieNames(input.Cookies),
			ProviderPayload:  providerResponse,
			ProviderProducts: productsPayload,
		},
	}, nil
}

func intelligentPathAndQuery(input IntelligentSearchProductsInput) (string, url.Values, []Facet, error) {
	facets, err := intelligentFacets(input)
	if err != nil {
		return "", nil, nil, err
	}
	segments := make([]string, 0, len(facets)*2)
	for _, facet := range facets {
		segments = append(segments, url.PathEscape(facet.Key), url.PathEscape(facet.Value))
	}
	path := "/api/io/_v/api/intelligent-search/product_search/" + strings.Join(segments, "/")
	query := url.Values{}
	if searchQuery := intelligentQuery(input); searchQuery != "" {
		query.Set("query", searchQuery)
	}
	query.Set("page", strconv.Itoa(input.Page))
	query.Set("count", strconv.Itoa(input.Count))
	if strings.TrimSpace(input.Sort) != "" {
		query.Set("sort", strings.TrimSpace(input.Sort))
	}
	if strings.TrimSpace(input.Locale) != "" {
		query.Set("locale", strings.TrimSpace(input.Locale))
	}
	if input.HideUnavailable != nil {
		query.Set("hideUnavailableItems", boolString(*input.HideUnavailable))
	}
	if strings.TrimSpace(input.SimulationBehavior) != "" {
		query.Set("simulationBehavior", strings.TrimSpace(input.SimulationBehavior))
	}
	return path, query, facets, nil
}

func intelligentFacets(input IntelligentSearchProductsInput) ([]Facet, error) {
	facets := []Facet{{Key: "trade-policy", Value: strings.TrimSpace(input.TradePolicy)}}
	for _, raw := range nonBlank(input.Facets) {
		key, value, ok := strings.Cut(raw, "=")
		if !ok || strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return nil, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "facet must use key=value format."}
		}
		facets = append(facets, Facet{Key: strings.TrimSpace(key), Value: strings.TrimSpace(value)})
	}
	return facets, nil
}

func intelligentProductsPayload(response map[string]any) []map[string]any {
	for _, key := range []string{"products", "items"} {
		if products := mapsFromAny(response[key]); len(products) > 0 {
			return products
		}
	}
	if records, ok := response["records"].(map[string]any); ok {
		return mapsFromAny(records["products"])
	}
	return nil
}

func mapsFromAny(raw any) []map[string]any {
	values, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(values))
	for _, value := range values {
		fields, ok := value.(map[string]any)
		if ok {
			out = append(out, fields)
		}
	}
	return out
}

func totalFromIntelligentResponse(response map[string]any) *int {
	for _, key := range []string{"total", "totalCount", "recordsFiltered"} {
		if value := intFromAny(response[key]); value > 0 {
			return &value
		}
	}
	if records, ok := response["records"].(map[string]any); ok {
		for _, key := range []string{"total", "totalCount", "filtered"} {
			if value := intFromAny(records[key]); value > 0 {
				return &value
			}
		}
	}
	return nil
}

func cookieHeader(cookies []string) string {
	return strings.Join(nonBlank(cookies), "; ")
}

func cookieNames(cookies []string) []string {
	values := nonBlank(cookies)
	names := make([]string, 0, len(values))
	for _, raw := range values {
		name, _, _ := strings.Cut(raw, "=")
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			names = append(names, trimmed)
		}
	}
	return names
}

func sortOrRelevance(sort string) string {
	if strings.TrimSpace(sort) == "" {
		return "relevance"
	}
	return strings.TrimSpace(sort)
}
