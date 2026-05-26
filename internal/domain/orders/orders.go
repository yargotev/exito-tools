package orders

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

const (
	CapabilityGetID = "orders.get"
	DomainName      = "orders"

	ErrorOrderNotFound                 = "ORDER_NOT_FOUND"
	ErrorOrdersNotConfigured           = "ORDERS_NOT_CONFIGURED"
	ErrorOrdersProviderUnavailable     = "ORDERS_PROVIDER_UNAVAILABLE"
	ErrorOrdersProviderInvalidResponse = "ORDERS_PROVIDER_INVALID_RESPONSE"
)

// Order is the domain-owned result model for an order lookup.
type Order struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

// GetResult is the stable use-case result shape for orders.get.
type GetResult struct {
	Order Order `json:"order"`
}

// GetInput is the schema-shaped input accepted by the orders.get use case.
type GetInput struct {
	ID string
}

// Getter retrieves orders using domain-owned models.
type Getter interface {
	Get(ctx context.Context, input GetInput) (Order, error)
}

// GetUseCase runs the orders.get behavior without surface dependencies.
type GetUseCase struct {
	getter Getter
}

// NewGetUseCase creates the orders.get use case.
func NewGetUseCase(getter Getter) GetUseCase {
	return GetUseCase{getter: getter}
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

// NewGetCapability adapts the orders.get use case into a neutral executable Capability.
func NewGetCapability(getter Getter) capability.Executable {
	useCase := NewGetUseCase(getter)
	return capability.Executable{
		Definition: Definition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			result, err := useCase.Execute(ctx, GetInput{ID: request.Input["id"].(string)})
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
		}},
	}
}
