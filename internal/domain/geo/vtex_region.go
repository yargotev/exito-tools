package geo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
)

const CapabilityResolveVTEXRegionID = "geo.resolve-vtex-region"

const ErrorGeoInvalidInput = "GEO_INVALID_INPUT"

// ResolveVTEXRegionInput is the schema-shaped input accepted by geo.resolve-vtex-region.
type ResolveVTEXRegionInput struct {
	Brand        string
	Country      string
	SalesChannel string
	Longitude    string
	Latitude     string
}

// Coordinates describes a longitude/latitude pair in VTEX request order.
type Coordinates struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// RegionSeller describes a seller returned by VTEX Checkout Regions.
type RegionSeller struct {
	ID   string         `json:"id"`
	Name string         `json:"name,omitempty"`
	Raw  map[string]any `json:"raw,omitempty"`
}

// Region describes a VTEX checkout region returned for coordinates.
type Region struct {
	ID      string         `json:"id"`
	Sellers []RegionSeller `json:"sellers,omitempty"`
	Raw     map[string]any `json:"raw,omitempty"`
}

// RegionDiagnostics captures non-secret request and provider details.
type RegionDiagnostics struct {
	RequestPath     string            `json:"requestPath,omitempty"`
	RequestQuery    map[string]string `json:"requestQuery,omitempty"`
	ProviderPayload any               `json:"providerPayload,omitempty"`
}

// ResolveVTEXRegionResult is the stable use-case result shape for VTEX region diagnostics.
type ResolveVTEXRegionResult struct {
	Brand        string            `json:"brand"`
	Country      string            `json:"country"`
	SalesChannel string            `json:"salesChannel"`
	Coordinates  Coordinates       `json:"coordinates"`
	HasCoverage  bool              `json:"hasCoverage"`
	Regions      []Region          `json:"regions,omitempty"`
	Sellers      []RegionSeller    `json:"sellers"`
	Diagnostics  RegionDiagnostics `json:"diagnostics,omitempty"`
}

// VTEXRegionResolver resolves VTEX checkout region coverage using domain-owned models.
type VTEXRegionResolver interface {
	ResolveVTEXRegion(ctx context.Context, input ResolveVTEXRegionInput) (ResolveVTEXRegionResult, error)
}

// ResolveVTEXRegionUseCase runs geo.resolve-vtex-region without surface dependencies.
type ResolveVTEXRegionUseCase struct {
	resolver VTEXRegionResolver
}

// NewResolveVTEXRegionUseCase creates the VTEX region use case.
func NewResolveVTEXRegionUseCase(resolver VTEXRegionResolver) ResolveVTEXRegionUseCase {
	return ResolveVTEXRegionUseCase{resolver: resolver}
}

// Execute resolves VTEX region coverage diagnostics.
func (u ResolveVTEXRegionUseCase) Execute(ctx context.Context, input ResolveVTEXRegionInput) (ResolveVTEXRegionResult, error) {
	if u.resolver == nil {
		return ResolveVTEXRegionResult{}, capability.StructuredError{Code: ErrorGeoNotConfigured, Message: "VTEX region resolver is not configured."}
	}
	input.Brand = normalizeBrand(input.Brand)
	if strings.TrimSpace(input.Country) == "" {
		input.Country = "COL"
	}
	if err := validateResolveVTEXRegionInput(input); err != nil {
		return ResolveVTEXRegionResult{}, err
	}
	return u.resolver.ResolveVTEXRegion(ctx, input)
}

// NewResolveVTEXRegionCapability adapts the use case into a neutral executable Capability.
func NewResolveVTEXRegionCapability(resolver VTEXRegionResolver) capability.Executable {
	useCase := NewResolveVTEXRegionUseCase(resolver)
	return capability.Executable{
		Definition: ResolveVTEXRegionDefinition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			input, err := resolveVTEXRegionInputFromCapability(request.Input)
			if err != nil {
				return capability.ExecutionResult{}, err
			}
			result, err := useCase.Execute(ctx, input)
			if err != nil {
				return capability.ExecutionResult{}, err
			}
			return capability.ExecutionResult{Data: result}, nil
		},
	}
}

// ResolveVTEXRegionDefinition returns the neutral discovery contract.
func ResolveVTEXRegionDefinition() capability.Definition {
	return capability.Definition{
		ID:          CapabilityResolveVTEXRegionID,
		Domain:      DomainName,
		Version:     "1.0.0",
		Title:       "Resolve VTEX region coverage",
		Description: "Checks VTEX Checkout Regions coverage for known longitude/latitude coordinates without creating sessions or writing address data.",
		Risk:        capability.RiskReadOnly,
		Audiences:   []capability.Audience{capability.AudienceAgents, capability.AudiencePeople},
		Visibility:  []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."},
			{Name: "country", Type: capability.InputTypeString, Required: false, Description: "VTEX country code. Defaults to COL."},
			{Name: "salesChannel", Type: capability.InputTypeString, Required: true, Description: "VTEX sales channel/trade policy passed as sc."},
			{Name: "longitude", Type: capability.InputTypeString, Required: true, Description: "Longitude coordinate; sent before latitude in geoCoordinates."},
			{Name: "latitude", Type: capability.InputTypeString, Required: true, Description: "Latitude coordinate; sent after longitude in geoCoordinates."},
		}},
	}
}

func resolveVTEXRegionInputFromCapability(input capability.Input) (ResolveVTEXRegionInput, error) {
	out := ResolveVTEXRegionInput{}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["country"].(string); ok {
		out.Country = value
	}
	if value, ok := input["salesChannel"].(string); ok {
		out.SalesChannel = value
	}
	longitude, err := coordinateString(input["longitude"], "longitude")
	if err != nil {
		return ResolveVTEXRegionInput{}, err
	}
	latitude, err := coordinateString(input["latitude"], "latitude")
	if err != nil {
		return ResolveVTEXRegionInput{}, err
	}
	out.Longitude = longitude
	out.Latitude = latitude
	return out, nil
}

func validateResolveVTEXRegionInput(input ResolveVTEXRegionInput) error {
	if strings.TrimSpace(input.SalesChannel) == "" {
		return capability.StructuredError{Code: ErrorGeoInvalidInput, Message: "salesChannel is required."}
	}
	if strings.TrimSpace(input.Longitude) == "" || strings.TrimSpace(input.Latitude) == "" {
		return capability.StructuredError{Code: ErrorGeoInvalidInput, Message: "longitude and latitude are required."}
	}
	return nil
}

func coordinateString(value any, field string) (string, error) {
	switch typed := value.(type) {
	case nil:
		return "", nil
	case string:
		return strings.TrimSpace(typed), nil
	case int:
		return strconv.Itoa(typed), nil
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64), nil
	default:
		return "", capability.StructuredError{Code: ErrorGeoInvalidInput, Message: fmt.Sprintf("%s must be a string or number.", field)}
	}
}

func normalizeBrand(brand string) string {
	switch strings.ToLower(strings.TrimSpace(brand)) {
	case "carulla":
		return "carulla"
	default:
		return "exito"
	}
}
