package geo

import (
	"context"
	"net/http"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

// HTTPGeocoderConfig contains the provider settings required by HTTPGeocoder.
type HTTPGeocoderConfig struct {
	BaseURL string
	Token   string
}

// HTTPGeocoder calls the configured Geo provider and maps its DTO to domain output.
type HTTPGeocoder struct {
	client httpclient.Client
}

// NewHTTPGeocoder creates a Geo provider-backed geocoder.
func NewHTTPGeocoder(config HTTPGeocoderConfig, client *http.Client) HTTPGeocoder {
	return HTTPGeocoder{
		client: httpclient.New(httpclient.Config{
			BaseURL: config.BaseURL,
			Token:   config.Token,
			Client:  client,
		}),
	}
}

// GeocodeAddress calls the provider geocode-address endpoint and maps the response.
func (g HTTPGeocoder) GeocodeAddress(ctx context.Context, input GeocodeAddressInput) (GeocodeAddressResult, error) {
	if !g.client.Configured() {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoNotConfigured,
			Message: "Geo client is not configured.",
		}
	}

	request, err := g.client.NewJSONRequest(ctx, http.MethodPost, "/geocode-address", geoProviderRequest(input))
	if err != nil {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoNotConfigured,
			Message: "Geo provider base URL is invalid.",
		}
	}

	response, err := g.client.Do(request)
	if err != nil {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoProviderUnavailable,
			Message: "Geo provider request failed.",
		}
	}
	defer func() { _ = response.Body.Close() }()

	if !httpclient.Successful(response) {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoProviderUnavailable,
			Message: "Geo provider returned an unsuccessful response.",
		}
	}

	var providerResponse geoProviderResponse
	if err := httpclient.DecodeJSONResponse(response, &providerResponse); err != nil {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoProviderInvalidResponse,
			Message: "Geo provider returned an invalid response.",
		}
	}

	return providerResponse.toDomain(), nil
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
