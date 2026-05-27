package orders

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

const (
	CapabilityGetID     = "orders.get"
	CapabilityGetVTEXID = "orders.get-vtex"
	DomainName          = "orders"

	ErrorOrderNotFound                 = "ORDER_NOT_FOUND"
	ErrorOrdersNotConfigured           = "ORDERS_NOT_CONFIGURED"
	ErrorOrdersProviderUnavailable     = "ORDERS_PROVIDER_UNAVAILABLE"
	ErrorOrdersProviderInvalidResponse = "ORDERS_PROVIDER_INVALID_RESPONSE"
)

// Order is the domain-owned result model for an order lookup.
type Order struct {
	ID             string         `json:"id"`
	Status         string         `json:"status"`
	CreatedAt      string         `json:"createdAt"`
	CustomerName   string         `json:"customerName,omitempty"`
	Email          string         `json:"email,omitempty"`
	OrderTotal     float64        `json:"orderTotal,omitempty"`
	StatusOrderMax string         `json:"statusOrderMax,omitempty"`
	StatusOrderMin string         `json:"statusOrderMin,omitempty"`
	Items          *OrderItems    `json:"items,omitempty"`
	Details        map[string]any `json:"details,omitempty"`
}

// OrderItems groups GEOMS line items by food and non-food buckets.
type OrderItems struct {
	Food    []map[string]any `json:"food"`
	NotFood []map[string]any `json:"notFood"`
}

// GetResult is the stable use-case result shape for orders.get.
type GetResult struct {
	Order Order `json:"order"`
}

// VTEXOMSOrder is the domain-owned result model for an order lookup in VTEX OMS.
type VTEXOMSOrder struct {
	ID                string         `json:"id"`
	Sequence          string         `json:"sequence,omitempty"`
	Status            string         `json:"status,omitempty"`
	StatusDescription string         `json:"statusDescription,omitempty"`
	CreationDate      string         `json:"creationDate,omitempty"`
	ClientName        string         `json:"clientName,omitempty"`
	Email             string         `json:"email,omitempty"`
	TotalValue        float64        `json:"totalValue,omitempty"`
	Brand             string         `json:"brand,omitempty"`
	Details           map[string]any `json:"details,omitempty"`
}

// GetVTEXOMSResult is the stable use-case result shape for orders.get-vtex.
type GetVTEXOMSResult struct {
	Order VTEXOMSOrder `json:"order"`
}

// GetInput is the schema-shaped input accepted by the orders.get use case.
type GetInput struct {
	ID        string
	OrderType string
}

// GetVTEXOMSInput is the schema-shaped input accepted by the orders.get-vtex use case.
type GetVTEXOMSInput struct {
	ID    string
	Brand string
}

// Getter retrieves orders using domain-owned models.
type Getter interface {
	Get(ctx context.Context, input GetInput) (Order, error)
}

// VTEXOMSGetterPort retrieves orders from VTEX OMS using domain-owned models.
type VTEXOMSGetterPort interface {
	GetVTEXOMS(ctx context.Context, input GetVTEXOMSInput) (VTEXOMSOrder, error)
}

// GetUseCase runs the orders.get behavior without surface dependencies.
type GetUseCase struct {
	getter Getter
}

// GetVTEXOMSUseCase runs the orders.get-vtex behavior without surface dependencies.
type GetVTEXOMSUseCase struct {
	getter VTEXOMSGetterPort
}

// NewGetUseCase creates the orders.get use case.
func NewGetUseCase(getter Getter) GetUseCase {
	return GetUseCase{getter: getter}
}

// NewGetVTEXOMSUseCase creates the orders.get-vtex use case.
func NewGetVTEXOMSUseCase(getter VTEXOMSGetterPort) GetVTEXOMSUseCase {
	return GetVTEXOMSUseCase{getter: getter}
}

// Execute gets one order by identifier.
func (u GetUseCase) Execute(ctx context.Context, input GetInput) (GetResult, error) {
	if u.getter == nil {
		return GetResult{}, capability.StructuredError{
			Code:    ErrorOrdersNotConfigured,
			Message: "Orders client is not configured.",
		}
	}

	order, err := u.getter.Get(ctx, input)
	if err != nil {
		return GetResult{}, err
	}

	return GetResult{Order: order}, nil
}

// Execute gets one order by identifier from VTEX OMS.
func (u GetVTEXOMSUseCase) Execute(ctx context.Context, input GetVTEXOMSInput) (GetVTEXOMSResult, error) {
	if u.getter == nil {
		return GetVTEXOMSResult{}, capability.StructuredError{
			Code:    ErrorOrdersNotConfigured,
			Message: "VTEX OMS client is not configured.",
		}
	}

	input.Brand = normalizedVTEXOMSBrand(input.Brand)
	order, err := u.getter.GetVTEXOMS(ctx, input)
	if err != nil {
		return GetVTEXOMSResult{}, err
	}

	return GetVTEXOMSResult{Order: order}, nil
}

// NewGetCapability adapts the orders.get use case into a neutral executable Capability.
func NewGetCapability(getter Getter) capability.Executable {
	useCase := NewGetUseCase(getter)
	return capability.Executable{
		Definition: Definition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			input := GetInput{ID: request.Input["id"].(string)}
			if orderType, ok := request.Input["orderType"].(string); ok {
				input.OrderType = orderType
			}
			result, err := useCase.Execute(ctx, input)
			if err != nil {
				return capability.ExecutionResult{}, err
			}

			return capability.ExecutionResult{Data: result}, nil
		},
	}
}

// NewGetVTEXOMSCapability adapts the orders.get-vtex use case into a neutral executable Capability.
func NewGetVTEXOMSCapability(getter VTEXOMSGetterPort) capability.Executable {
	useCase := NewGetVTEXOMSUseCase(getter)
	return capability.Executable{
		Definition: VTEXOMSDefinition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			input := GetVTEXOMSInput{ID: request.Input["id"].(string)}
			if brand, ok := request.Input["brand"].(string); ok {
				input.Brand = brand
			}
			result, err := useCase.Execute(ctx, input)
			if err != nil {
				return capability.ExecutionResult{}, err
			}

			return capability.ExecutionResult{Data: result}, nil
		},
	}
}

// Definition returns the neutral orders.get discovery contract.
func Definition() capability.Definition {
	return capability.Definition{
		ID:          CapabilityGetID,
		Domain:      DomainName,
		Version:     "1.0.0",
		Title:       "Get order",
		Description: "Gets an order by its identifier.",
		Risk:        capability.RiskReadOnly,
		Audiences:   []capability.Audience{capability.AudienceAgents, capability.AudiencePeople},
		Visibility:  []capability.Visibility{capability.VisibilityCLI, capability.VisibilityTUI, capability.VisibilityCommandPalette},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "id", Type: capability.InputTypeString, Required: true, Description: "Order identifier."},
			{Name: "orderType", Type: capability.InputTypeString, Required: false, Description: "GEOMS order type filter, such as ExitoEcomm or CarullaEcomm."},
		}},
	}
}

// VTEXOMSDefinition returns the neutral orders.get-vtex discovery contract.
func VTEXOMSDefinition() capability.Definition {
	return capability.Definition{
		ID:          CapabilityGetVTEXID,
		Domain:      DomainName,
		Version:     "1.0.0",
		Title:       "Get VTEX OMS order",
		Description: "Gets an order by its identifier from VTEX OMS.",
		Risk:        capability.RiskReadOnly,
		Audiences:   []capability.Audience{capability.AudienceAgents, capability.AudiencePeople},
		Visibility:  []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "id", Type: capability.InputTypeString, Required: true, Description: "VTEX OMS order identifier, such as 1611511090420-01."},
			{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."},
		}},
	}
}

func normalizedVTEXOMSBrand(brand string) string {
	switch brand {
	case "carulla", "Carulla":
		return "carulla"
	default:
		return "exito"
	}
}
