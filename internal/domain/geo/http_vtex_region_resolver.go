package geo

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

const maxVTEXRegionResponseBodyBytes = 4 << 20

// HTTPVTEXRegionResolverConfig contains provider settings for VTEX Checkout Regions.
type HTTPVTEXRegionResolverConfig struct {
	BaseURL string
}

// HTTPVTEXRegionResolver calls VTEX Checkout Regions and maps its DTO to Geo output.
type HTTPVTEXRegionResolver struct {
	baseURL string
	client  httpclient.Client
}

// NewHTTPVTEXRegionResolver creates a VTEX Checkout Regions-backed resolver.
func NewHTTPVTEXRegionResolver(config HTTPVTEXRegionResolverConfig, client *http.Client) HTTPVTEXRegionResolver {
	return HTTPVTEXRegionResolver{
		baseURL: strings.TrimSpace(config.BaseURL),
		client:  httpclient.New(httpclient.Config{BaseURL: config.BaseURL, Client: client}),
	}
}

// ResolveVTEXRegion calls VTEX Checkout Regions using geoCoordinates=longitude;latitude.
func (r HTTPVTEXRegionResolver) ResolveVTEXRegion(ctx context.Context, input ResolveVTEXRegionInput) (ResolveVTEXRegionResult, error) {
	if strings.TrimSpace(r.baseURL) == "" {
		return ResolveVTEXRegionResult{}, capability.StructuredError{Code: ErrorGeoNotConfigured, Message: "VTEX region resolver is not configured."}
	}

	path := "/api/checkout/pub/regions"
	query := url.Values{}
	query.Set("country", strings.TrimSpace(input.Country))
	query.Set("sc", strings.TrimSpace(input.SalesChannel))
	query.Set("geoCoordinates", strings.TrimSpace(input.Longitude)+";"+strings.TrimSpace(input.Latitude))

	request, err := r.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ResolveVTEXRegionResult{}, capability.StructuredError{Code: ErrorGeoNotConfigured, Message: "VTEX region provider base URL is invalid."}
	}
	request.URL.RawQuery = query.Encode()

	response, err := r.client.Do(request)
	if err != nil {
		return ResolveVTEXRegionResult{}, capability.StructuredError{Code: ErrorGeoProviderUnavailable, Message: "VTEX region provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()

	if !httpclient.Successful(response) {
		return ResolveVTEXRegionResult{}, capability.StructuredError{Code: ErrorGeoProviderUnavailable, Message: "VTEX region provider returned an unsuccessful response."}
	}

	var providerPayload any
	limited := io.LimitReader(response.Body, maxVTEXRegionResponseBodyBytes)
	decoder := json.NewDecoder(limited)
	decoder.UseNumber()
	if err := decoder.Decode(&providerPayload); err != nil {
		return ResolveVTEXRegionResult{}, capability.StructuredError{Code: ErrorGeoProviderInvalidResponse, Message: "VTEX region provider returned an invalid response."}
	}

	sellers := sellersFromRegionPayload(providerPayload)
	requestQuery := map[string]string{}
	for key, values := range query {
		if len(values) > 0 {
			requestQuery[key] = values[0]
		}
	}
	return ResolveVTEXRegionResult{
		Brand:        normalizeBrand(input.Brand),
		Country:      strings.TrimSpace(input.Country),
		SalesChannel: strings.TrimSpace(input.SalesChannel),
		Coordinates:  Coordinates{Longitude: strings.TrimSpace(input.Longitude), Latitude: strings.TrimSpace(input.Latitude)},
		HasCoverage:  hasCoverage(sellers),
		Sellers:      sellers,
		Diagnostics: RegionDiagnostics{
			RequestPath:     path,
			RequestQuery:    requestQuery,
			ProviderPayload: providerPayload,
		},
	}, nil
}

func sellersFromRegionPayload(payload any) []RegionSeller {
	var sellers []RegionSeller
	walkRegionPayload(payload, &sellers)
	return dedupeSellers(sellers)
}

func walkRegionPayload(value any, sellers *[]RegionSeller) {
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			walkRegionPayload(item, sellers)
		}
	case map[string]any:
		if raw, ok := typed["sellers"]; ok {
			*sellers = append(*sellers, sellersFromAny(raw)...)
		}
		for _, child := range typed {
			walkRegionPayload(child, sellers)
		}
	}
}

func sellersFromAny(raw any) []RegionSeller {
	values, ok := raw.([]any)
	if !ok {
		return nil
	}
	sellers := make([]RegionSeller, 0, len(values))
	for _, value := range values {
		fields, ok := value.(map[string]any)
		if !ok {
			continue
		}
		id := firstString(fields, "id", "sellerId")
		if strings.TrimSpace(id) == "" {
			continue
		}
		sellers = append(sellers, RegionSeller{ID: id, Name: firstString(fields, "name", "sellerName"), Raw: fields})
	}
	return sellers
}

func dedupeSellers(values []RegionSeller) []RegionSeller {
	seen := map[string]bool{}
	out := make([]RegionSeller, 0, len(values))
	for _, seller := range values {
		key := strings.ToLower(strings.TrimSpace(seller.ID))
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, seller)
	}
	return out
}

func hasCoverage(sellers []RegionSeller) bool {
	// Historical Exito storefront logic treated coverage as true only when a seller
	// differed from the account/brand. That was useful for product-price flows where
	// seller "exitocol" did not identify the exact white-label store fulfilling the
	// item. This CLI reports VTEX Checkout Regions coverage diagnostics instead, so
	// any seller returned by Regions counts as coverage and the seller list remains
	// available for downstream business interpretation.
	return len(sellers) > 0
}

func firstString(fields map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := fields[key]; ok {
			switch typed := value.(type) {
			case string:
				if strings.TrimSpace(typed) != "" {
					return typed
				}
			case json.Number:
				return typed.String()
			case float64:
				return strconvFormatFloat(typed)
			case int:
				return strconv.Itoa(typed)
			}
		}
	}
	return ""
}

func strconvFormatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
