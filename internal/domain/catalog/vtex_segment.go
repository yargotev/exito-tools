package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
)

const CapabilityCreateVTEXSegmentID = "catalog.create-vtex-segment"

// CreateVTEXSegmentInput is the schema-shaped input accepted by catalog.create-vtex-segment.
type CreateVTEXSegmentInput struct {
	Brand         string
	RegionID      string
	SalesChannel  string
	IncludeCookie bool
}

// SegmentDiagnostics captures non-secret request/provider details for VTEX segment preparation.
type SegmentDiagnostics struct {
	RequestPath     string         `json:"requestPath,omitempty"`
	RequestPayload  map[string]any `json:"requestPayload,omitempty"`
	ProviderPayload map[string]any `json:"providerPayload,omitempty"`
}

// CreateVTEXSegmentResult is the stable use-case result shape for VTEX session segment preparation.
type CreateVTEXSegmentResult struct {
	Brand        string             `json:"brand"`
	RegionID     string             `json:"regionId"`
	SalesChannel string             `json:"salesChannel"`
	TokenSet     bool               `json:"tokenSet"`
	TokenLength  int                `json:"tokenLength,omitempty"`
	Cookie       string             `json:"cookie,omitempty"`
	Diagnostics  SegmentDiagnostics `json:"diagnostics,omitempty"`
}

// VTEXSegmentCreator creates VTEX segment tokens using domain-owned models.
type VTEXSegmentCreator interface {
	CreateVTEXSegment(ctx context.Context, input CreateVTEXSegmentInput) (CreateVTEXSegmentResult, error)
}

// CreateVTEXSegmentUseCase runs catalog.create-vtex-segment without surface dependencies.
type CreateVTEXSegmentUseCase struct {
	creator VTEXSegmentCreator
}

// NewCreateVTEXSegmentUseCase creates the VTEX segment preparation use case.
func NewCreateVTEXSegmentUseCase(creator VTEXSegmentCreator) CreateVTEXSegmentUseCase {
	return CreateVTEXSegmentUseCase{creator: creator}
}

// Execute creates a VTEX segment token from explicit region and sales-channel inputs.
func (u CreateVTEXSegmentUseCase) Execute(ctx context.Context, input CreateVTEXSegmentInput) (CreateVTEXSegmentResult, error) {
	if u.creator == nil {
		return CreateVTEXSegmentResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Sessions client is not configured."}
	}
	input.Brand = normalizedBrand(input.Brand)
	input.RegionID = strings.TrimSpace(input.RegionID)
	input.SalesChannel = strings.TrimSpace(input.SalesChannel)
	if err := validateCreateVTEXSegmentInput(input); err != nil {
		return CreateVTEXSegmentResult{}, err
	}
	return u.creator.CreateVTEXSegment(ctx, input)
}

// NewCreateVTEXSegmentCapability adapts the use case into a neutral executable Capability.
func NewCreateVTEXSegmentCapability(creator VTEXSegmentCreator) capability.Executable {
	useCase := NewCreateVTEXSegmentUseCase(creator)
	return capability.Executable{
		Definition: CreateVTEXSegmentDefinition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			input := createVTEXSegmentInputFromCapability(request.Input)
			result, err := useCase.Execute(ctx, input)
			if err != nil {
				return capability.ExecutionResult{}, err
			}
			return capability.ExecutionResult{Data: result}, nil
		},
	}
}

// CreateVTEXSegmentDefinition returns the neutral discovery contract.
func CreateVTEXSegmentDefinition() capability.Definition {
	return capability.Definition{
		ID:                   CapabilityCreateVTEXSegmentID,
		Domain:               DomainName,
		Version:              "1.0.0",
		Title:                "Create VTEX segment",
		Description:          "Creates a VTEX session segment token from an explicit region ID and sales channel.",
		Risk:                 capability.RiskSafeWrite,
		RequiresConfirmation: true,
		Audiences:            []capability.Audience{capability.AudienceAgents},
		Visibility:           []capability.Visibility{capability.VisibilityCLI},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."},
			{Name: "regionId", Type: capability.InputTypeString, Required: true, Description: "VTEX region ID to place in the session segment."},
			{Name: "salesChannel", Type: capability.InputTypeString, Required: true, Description: "VTEX sales channel/trade policy to place in the session segment."},
			{Name: "includeCookie", Type: capability.InputTypeBoolean, Required: false, Description: "Include an explicit vtex_segment cookie string in successful output."},
		}},
	}
}

func createVTEXSegmentInputFromCapability(input capability.Input) CreateVTEXSegmentInput {
	out := CreateVTEXSegmentInput{}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["regionId"].(string); ok {
		out.RegionID = value
	}
	if value, ok := input["salesChannel"].(string); ok {
		out.SalesChannel = value
	}
	if value, ok := input["includeCookie"].(bool); ok {
		out.IncludeCookie = value
	}
	return out
}

func validateCreateVTEXSegmentInput(input CreateVTEXSegmentInput) error {
	if input.RegionID == "" {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "regionId is required."}
	}
	if input.SalesChannel == "" {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "salesChannel is required."}
	}
	if brand := normalizedBrand(input.Brand); brand != "exito" && brand != "carulla" {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: fmt.Sprintf("Unsupported VTEX brand %q.", input.Brand)}
	}
	return nil
}
