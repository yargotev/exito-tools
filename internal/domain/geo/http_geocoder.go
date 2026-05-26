package geo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yargotev/exito-tools/internal/capability"
)

const defaultHTTPTimeout = 10 * time.Second

// HTTPGeocoderConfig contains the provider settings required by HTTPGeocoder.
type HTTPGeocoderConfig struct {
	BaseURL string
	Token   string
}

// HTTPGeocoder calls the configured Geo provider and maps its DTO to domain output.
type HTTPGeocoder struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewHTTPGeocoder creates a Geo provider-backed geocoder.
func NewHTTPGeocoder(config HTTPGeocoderConfig, client *http.Client) HTTPGeocoder {
	if client == nil {
		client = &http.Client{Timeout: defaultHTTPTimeout}
	}

	return HTTPGeocoder{
		baseURL: strings.TrimSpace(config.BaseURL),
		token:   strings.TrimSpace(config.Token),
		client:  client,
	}
}

// GeocodeAddress calls the provider geocode-address endpoint and maps the response.
func (g HTTPGeocoder) GeocodeAddress(ctx context.Context, input GeocodeAddressInput) (GeocodeAddressResult, error) {
	if g.baseURL == "" || g.token == "" {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoNotConfigured,
			Message: "Geo client is not configured.",
		}
	}

	endpoint, err := geocodeEndpoint(g.baseURL)
	if err != nil {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoNotConfigured,
			Message: "Geo provider base URL is invalid.",
		}
	}

	body, err := json.Marshal(geoProviderRequest(input))
	if err != nil {
		return GeocodeAddressResult{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return GeocodeAddressResult{}, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+g.token)
	request.Header.Set("Content-Type", "application/json")

	response, err := g.client.Do(request)
	if err != nil {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoProviderUnavailable,
			Message: "Geo provider request failed.",
		}
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoProviderUnavailable,
			Message: "Geo provider returned an unsuccessful response.",
		}
	}

	limited := io.LimitReader(response.Body, 1<<20)
	var providerResponse geoProviderResponse
	if err := json.NewDecoder(limited).Decode(&providerResponse); err != nil {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoProviderInvalidResponse,
			Message: "Geo provider returned an invalid response.",
		}
	}

	return providerResponse.toDomain(), nil
}

func geocodeEndpoint(baseURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("absolute URL required")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/geocode-address"
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

type geoProviderRequest struct {
	City    string `json:"city"`
	Address string `json:"address"`
}

type geoProviderResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
		Status    string `json:"estado"`
		DirTrad   string `json:"dirtrad"`
		Barrio    string `json:"barrio"`
		CodDANE   string `json:"coddane"`
	} `json:"data"`
}

func (r geoProviderResponse) toDomain() GeocodeAddressResult {
	return GeocodeAddressResult{
		Message:           r.Message,
		Success:           r.Success,
		Location:          Location{Latitude: r.Data.Latitude, Longitude: r.Data.Longitude},
		Status:            r.Data.Status,
		NormalizedAddress: r.Data.DirTrad,
		Neighborhood:      r.Data.Barrio,
		DANECode:          r.Data.CodDANE,
	}
}
