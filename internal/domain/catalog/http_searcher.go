package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

const maxCatalogResponseBodyBytes = 16 << 20

// HTTPSearcherConfig contains the provider settings required by HTTPSearcher.
type HTTPSearcherConfig struct {
	BaseURL string
}

// HTTPSearcher calls VTEX public Catalog Search and maps its DTO to domain output.
type HTTPSearcher struct {
	baseURL string
	client  httpclient.Client
}

// NewHTTPSearcher creates a VTEX Catalog-backed searcher.
func NewHTTPSearcher(config HTTPSearcherConfig, client *http.Client) HTTPSearcher {
	return HTTPSearcher{
		baseURL: strings.TrimSpace(config.BaseURL),
		client: httpclient.New(httpclient.Config{
			BaseURL: config.BaseURL,
			Client:  client,
		}),
	}
}

// SearchProducts calls VTEX product search using simple or advanced filters.
func (s HTTPSearcher) SearchProducts(ctx context.Context, input SearchProductsInput) (SearchProductsResult, error) {
	if strings.TrimSpace(s.baseURL) == "" {
		return SearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Catalog client is not configured."}
	}

	path, query, err := searchPathAndQuery(input)
	if err != nil {
		return SearchProductsResult{}, err
	}
	request, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return SearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Catalog provider base URL is invalid."}
	}
	request.URL.RawQuery = query.Encode()

	response, err := s.client.Do(request)
	if err != nil {
		return SearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogProviderUnavailable, Message: "VTEX Catalog provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()

	if !httpclient.Successful(response) {
		return SearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogProviderUnavailable, Message: "VTEX Catalog provider returned an unsuccessful response."}
	}

	var providerResponse []map[string]any
	limited := io.LimitReader(response.Body, maxCatalogResponseBodyBytes)
	if err := json.NewDecoder(limited).Decode(&providerResponse); err != nil {
		return SearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogProviderInvalidResponse, Message: "VTEX Catalog provider returned an invalid response."}
	}

	resources := response.Header.Get("resources")
	rangeStart, rangeEnd, total := parseResources(resources)
	products := make([]Product, 0, len(providerResponse))
	for _, fields := range providerResponse {
		products = append(products, productFromMap(fields))
	}
	return SearchProductsResult{
		Brand:      normalizedBrand(input.Brand),
		Products:   products,
		Count:      len(products),
		RangeStart: rangeStart,
		RangeEnd:   rangeEnd,
		Total:      total,
		Resources:  resources,
	}, nil
}

func searchPathAndQuery(input SearchProductsInput) (string, url.Values, error) {
	query := url.Values{}
	if input.By == "slug" {
		return "/api/catalog_system/pub/products/search/" + strings.Trim(strings.TrimSpace(input.Value), "/") + "/p", query, nil
	}

	filters, ft, err := searchFilters(input)
	if err != nil {
		return "", nil, err
	}
	for _, filter := range filters {
		query.Add("fq", filter)
	}
	if strings.TrimSpace(ft) != "" {
		query.Set("ft", strings.TrimSpace(ft))
	}
	if strings.TrimSpace(input.Order) != "" {
		query.Set("O", strings.TrimSpace(input.Order))
	}
	query.Set("_from", strconv.Itoa(input.From))
	query.Set("_to", strconv.Itoa(input.To))
	return "/api/catalog_system/pub/products/search", query, nil
}

func searchFilters(input SearchProductsInput) ([]string, string, error) {
	filters := append([]string{}, nonBlank(input.FQ)...)
	ft := input.FT
	if input.By == "" {
		return filters, ft, nil
	}
	value := strings.TrimSpace(input.Value)
	switch input.By {
	case "sku-id":
		filters = append(filters, "skuId:"+value)
	case "product-id":
		filters = append(filters, "productId:"+value)
	case "ref-id":
		filters = append(filters, "alternateIds_RefId:"+value)
	case "ean":
		filters = append(filters, "alternateIds_Ean:"+value)
	case "seller-id":
		filters = append(filters, "sellerId:"+value)
	case "category":
		filters = append(filters, "C:"+normalizeCategoryValue(value))
	case "brand-id":
		filters = append(filters, "B:"+value)
	case "collection-id":
		filters = append(filters, "productClusterIds:"+value)
	case "text":
		ft = value
	default:
		return nil, "", capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: fmt.Sprintf("Unsupported catalog lookup mode %q.", input.By)}
	}
	return filters, ft, nil
}

func normalizeCategoryValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(trimmed, "/") {
		return trimmed
	}
	return "/" + strings.Trim(trimmed, "/") + "/"
}

func parseResources(resources string) (*int, *int, *int) {
	trimmed := strings.TrimSpace(resources)
	if trimmed == "" {
		return nil, nil, nil
	}
	rangePart, totalPart, ok := strings.Cut(trimmed, "/")
	if !ok {
		return nil, nil, nil
	}
	startPart, endPart, ok := strings.Cut(rangePart, "-")
	if !ok {
		return nil, nil, nil
	}
	start, errStart := strconv.Atoi(strings.TrimSpace(startPart))
	end, errEnd := strconv.Atoi(strings.TrimSpace(endPart))
	total, errTotal := strconv.Atoi(strings.TrimSpace(totalPart))
	if errStart != nil || errEnd != nil || errTotal != nil {
		return nil, nil, nil
	}
	return &start, &end, &total
}

func productFromMap(fields map[string]any) Product {
	return Product{
		ProductID:        firstString(fields, "productId", "productID", "id"),
		ProductName:      firstString(fields, "productName", "name"),
		Brand:            firstString(fields, "brand"),
		BrandID:          intFromAny(fields["brandId"]),
		LinkText:         firstString(fields, "linkText"),
		ProductReference: firstString(fields, "productReference", "productReferenceCode"),
		CategoryID:       firstString(fields, "categoryId"),
		Link:             firstString(fields, "link"),
		Items:            skusFromAny(fields["items"]),
		Details:          fields,
	}
}

func skusFromAny(raw any) []SKU {
	values, ok := raw.([]any)
	if !ok {
		return nil
	}
	items := make([]SKU, 0, len(values))
	for _, value := range values {
		fields, ok := value.(map[string]any)
		if !ok {
			continue
		}
		items = append(items, SKU{
			ItemID:       firstString(fields, "itemId", "id"),
			Name:         firstString(fields, "name"),
			NameComplete: firstString(fields, "nameComplete"),
			EAN:          firstString(fields, "ean"),
			ReferenceIDs: refIDsFromAny(fields["referenceId"]),
			Sellers:      sellersFromAny(fields["sellers"]),
		})
	}
	return items
}

func refIDsFromAny(raw any) []RefID {
	values, ok := raw.([]any)
	if !ok {
		return nil
	}
	refs := make([]RefID, 0, len(values))
	for _, value := range values {
		fields, ok := value.(map[string]any)
		if !ok {
			continue
		}
		refs = append(refs, RefID{Key: firstString(fields, "Key", "key"), Value: firstString(fields, "Value", "value")})
	}
	return refs
}

func sellersFromAny(raw any) []Seller {
	values, ok := raw.([]any)
	if !ok {
		return nil
	}
	sellers := make([]Seller, 0, len(values))
	for _, value := range values {
		fields, ok := value.(map[string]any)
		if !ok {
			continue
		}
		offer, _ := fields["commertialOffer"].(map[string]any)
		sellers = append(sellers, Seller{
			SellerID:      firstString(fields, "sellerId"),
			SellerName:    firstString(fields, "sellerName"),
			SellerDefault: boolFromAny(fields["sellerDefault"]),
			Price:         floatFromAny(offer["Price"]),
			ListPrice:     floatFromAny(offer["ListPrice"]),
			Available:     boolFromAny(offer["IsAvailable"]),
			Quantity:      intFromAny(offer["AvailableQuantity"]),
		})
	}
	return sellers
}

func firstString(fields map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := fields[key]; ok {
			switch typed := value.(type) {
			case string:
				if strings.TrimSpace(typed) != "" {
					return typed
				}
			case float64:
				return strconv.FormatFloat(typed, 'f', -1, 64)
			case int:
				return strconv.Itoa(typed)
			}
		}
	}
	return ""
}

func floatFromAny(value any) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	case json.Number:
		value, _ := typed.Float64()
		return value
	default:
		return 0
	}
}

func intFromAny(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case float64:
		return int(typed)
	case json.Number:
		value, _ := typed.Int64()
		return int(value)
	default:
		return 0
	}
}

func boolFromAny(value any) bool {
	boolValue, _ := value.(bool)
	return boolValue
}
