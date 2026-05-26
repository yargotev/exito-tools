package geo

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

const (
	CapabilityGeocodeAddressID = "geo.geocode-address"
	DomainName                 = "geo"

	ErrorGeoNotConfigured           = "GEO_NOT_CONFIGURED"
	ErrorGeoProviderUnavailable     = "GEO_PROVIDER_UNAVAILABLE"
	ErrorGeoProviderInvalidResponse = "GEO_PROVIDER_INVALID_RESPONSE"
)

// Location is the domain-owned location result for geocoding.
type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// GeocodeAddressResult is the stable use-case result shape for geo.geocode-address.
type GeocodeAddressResult struct {
	Message           string   `json:"message"`
	Success           bool     `json:"success"`
	Location          Location `json:"location"`
	Status            string   `json:"status"`
	NormalizedAddress string   `json:"normalizedAddress"`
	Neighborhood      string   `json:"neighborhood"`
	DANECode          string   `json:"daneCode"`
}

// GeocodeAddressInput is the schema-shaped input accepted by the geo.geocode-address use case.
type GeocodeAddressInput struct {
	City    string
	Address string
}

// Geocoder geocodes addresses using domain-owned models.
type Geocoder interface {
	GeocodeAddress(ctx context.Context, input GeocodeAddressInput) (GeocodeAddressResult, error)
}

// GeocodeAddressUseCase runs the geo.geocode-address behavior without surface dependencies.
type GeocodeAddressUseCase struct {
	geocoder Geocoder
}

// NewGeocodeAddressUseCase creates the geo.geocode-address use case.
func NewGeocodeAddressUseCase(geocoder Geocoder) GeocodeAddressUseCase {
	return GeocodeAddressUseCase{geocoder: geocoder}
}

// Execute geocodes a city/address pair.
func (u GeocodeAddressUseCase) Execute(ctx context.Context, input GeocodeAddressInput) (GeocodeAddressResult, error) {
	if u.geocoder == nil {
		return GeocodeAddressResult{}, capability.StructuredError{
			Code:    ErrorGeoNotConfigured,
			Message: "Geo client is not configured.",
		}
	}

	return u.geocoder.GeocodeAddress(ctx, input)
}

// NewGeocodeAddressCapability adapts the geo.geocode-address use case into a neutral executable Capability.
func NewGeocodeAddressCapability(geocoder Geocoder) capability.Executable {
	useCase := NewGeocodeAddressUseCase(geocoder)
	return capability.Executable{
		Definition: Definition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			result, err := useCase.Execute(ctx, GeocodeAddressInput{
				City:    request.Input["city"].(string),
				Address: request.Input["address"].(string),
			})
			if err != nil {
				return capability.ExecutionResult{}, err
			}

			return capability.ExecutionResult{Data: result}, nil
		},
	}
}

// Definition returns the neutral geo.geocode-address discovery contract.
func Definition() capability.Definition {
	return capability.Definition{
		ID:          CapabilityGeocodeAddressID,
		Domain:      DomainName,
		Version:     "1.0.0",
		Title:       "Geocode address",
		Description: "Geocodes a city/address pair using the configured Geo provider.",
		Risk:        capability.RiskReadOnly,
		Audiences:   []capability.Audience{capability.AudienceAgents, capability.AudiencePeople},
		Visibility:  []capability.Visibility{capability.VisibilityCLI, capability.VisibilityTUI, capability.VisibilityCommandPalette},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "city", Type: capability.InputTypeString, Required: true, Description: "City name accepted by the Geo provider."},
			{Name: "address", Type: capability.InputTypeString, Required: true, Description: "Address to geocode."},
		}},
	}
}
