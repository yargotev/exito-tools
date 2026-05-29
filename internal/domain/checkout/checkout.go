package checkout

import (
	"context"
	"fmt"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
)

const (
	CapabilityGetOrderFormID    = "checkout.get-order-form"
	CapabilityCreateOrderFormID = "checkout.create-order-form"
	DomainName                  = "checkout"

	ErrorCheckoutNotConfigured           = "CHECKOUT_NOT_CONFIGURED"
	ErrorCheckoutProviderUnavailable     = "CHECKOUT_PROVIDER_UNAVAILABLE"
	ErrorCheckoutProviderInvalidResponse = "CHECKOUT_PROVIDER_INVALID_RESPONSE"
	ErrorCheckoutInvalidInput            = "CHECKOUT_INVALID_INPUT"
)

type GetOrderFormInput struct {
	Brand       string
	OrderFormID string
}

type CreateOrderFormInput struct {
	Brand        string
	SalesChannel string
}

type OrderFormSummary struct {
	Brand                string             `json:"brand"`
	ID                   string             `json:"id"`
	SalesChannel         string             `json:"salesChannel,omitempty"`
	Value                int                `json:"value"`
	Totalizers           []TotalizerSummary `json:"totalizers,omitempty"`
	Items                []ItemSummary      `json:"items,omitempty"`
	ItemCount            int                `json:"itemCount"`
	ClientProfileDataSet bool               `json:"clientProfileDataSet"`
	ShippingDataSet      bool               `json:"shippingDataSet"`
	Diagnostics          Diagnostics        `json:"diagnostics,omitempty"`
}

type TotalizerSummary struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Value int    `json:"value"`
}

type ItemSummary struct {
	ID           string `json:"id,omitempty"`
	ProductID    string `json:"productId,omitempty"`
	Name         string `json:"name,omitempty"`
	Quantity     int    `json:"quantity"`
	Seller       string `json:"seller,omitempty"`
	Price        int    `json:"price,omitempty"`
	SellingPrice int    `json:"sellingPrice,omitempty"`
	Availability string `json:"availability,omitempty"`
}

type Diagnostics struct {
	RequestPath    string `json:"requestPath,omitempty"`
	ProviderStatus int    `json:"providerStatus,omitempty"`
}

type GetOrderFormResult struct {
	OrderForm OrderFormSummary `json:"orderForm"`
}

type CreateOrderFormResult struct {
	OrderForm OrderFormSummary `json:"orderForm"`
}

type Getter interface {
	GetOrderForm(ctx context.Context, input GetOrderFormInput) (OrderFormSummary, error)
}

type Creator interface {
	CreateOrderForm(ctx context.Context, input CreateOrderFormInput) (OrderFormSummary, error)
}

type (
	GetOrderFormUseCase    struct{ getter Getter }
	CreateOrderFormUseCase struct{ creator Creator }
)

func NewGetOrderFormUseCase(getter Getter) GetOrderFormUseCase {
	return GetOrderFormUseCase{getter: getter}
}

func NewCreateOrderFormUseCase(creator Creator) CreateOrderFormUseCase {
	return CreateOrderFormUseCase{creator: creator}
}

func (u GetOrderFormUseCase) Execute(ctx context.Context, input GetOrderFormInput) (GetOrderFormResult, error) {
	if u.getter == nil {
		return GetOrderFormResult{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
	}
	input.Brand = normalizeBrandInput(input.Brand)
	input.OrderFormID = strings.TrimSpace(input.OrderFormID)
	if err := validateGetOrderFormInput(input); err != nil {
		return GetOrderFormResult{}, err
	}
	orderForm, err := u.getter.GetOrderForm(ctx, input)
	if err != nil {
		return GetOrderFormResult{}, err
	}
	return GetOrderFormResult{OrderForm: orderForm}, nil
}

func (u CreateOrderFormUseCase) Execute(ctx context.Context, input CreateOrderFormInput) (CreateOrderFormResult, error) {
	if u.creator == nil {
		return CreateOrderFormResult{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
	}
	input.Brand = normalizeBrandInput(input.Brand)
	input.SalesChannel = strings.TrimSpace(input.SalesChannel)
	if err := validateCreateOrderFormInput(input); err != nil {
		return CreateOrderFormResult{}, err
	}
	orderForm, err := u.creator.CreateOrderForm(ctx, input)
	if err != nil {
		return CreateOrderFormResult{}, err
	}
	return CreateOrderFormResult{OrderForm: orderForm}, nil
}

func NewGetOrderFormCapability(getter Getter) capability.Executable {
	useCase := NewGetOrderFormUseCase(getter)
	return capability.Executable{Definition: GetOrderFormDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, getOrderFormInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		return capability.ExecutionResult{Data: result}, nil
	}}
}

func NewCreateOrderFormCapability(creator Creator) capability.Executable {
	useCase := NewCreateOrderFormUseCase(creator)
	return capability.Executable{Definition: CreateOrderFormDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, createOrderFormInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		return capability.ExecutionResult{Data: result}, nil
	}}
}

func GetOrderFormDefinition() capability.Definition {
	return capability.Definition{ID: CapabilityGetOrderFormID, Domain: DomainName, Version: "1.0.0", Title: "Get VTEX Checkout orderForm", Description: "Gets a known VTEX Checkout orderForm by ID and returns a redacted summary.", Risk: capability.RiskReadOnly, Audiences: []capability.Audience{capability.AudienceAgents, capability.AudiencePeople}, Visibility: []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette}, InputSchema: &capability.InputSchema{Fields: []capability.InputField{{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."}, {Name: "orderFormId", Type: capability.InputTypeString, Required: true, Description: "VTEX Checkout orderForm identifier."}}}}
}

func CreateOrderFormDefinition() capability.Definition {
	return capability.Definition{ID: CapabilityCreateOrderFormID, Domain: DomainName, Version: "1.0.0", Title: "Create VTEX Checkout orderForm", Description: "Creates a fresh VTEX Checkout orderForm for a brand and sales channel.", Risk: capability.RiskSafeWrite, RequiresConfirmation: true, Audiences: []capability.Audience{capability.AudienceAgents, capability.AudiencePeople}, Visibility: []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette}, InputSchema: &capability.InputSchema{Fields: []capability.InputField{{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."}, {Name: "salesChannel", Type: capability.InputTypeString, Required: true, Description: "VTEX sales channel/trade policy used as sc."}}}}
}

func getOrderFormInputFromCapability(input capability.Input) GetOrderFormInput {
	out := GetOrderFormInput{}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["orderFormId"].(string); ok {
		out.OrderFormID = value
	}
	return out
}

func createOrderFormInputFromCapability(input capability.Input) CreateOrderFormInput {
	out := CreateOrderFormInput{}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["salesChannel"].(string); ok {
		out.SalesChannel = value
	}
	return out
}

func validateGetOrderFormInput(input GetOrderFormInput) error {
	if input.OrderFormID == "" {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "orderFormId is required."}
	}
	return validateBrand(input.Brand)
}

func validateCreateOrderFormInput(input CreateOrderFormInput) error {
	if input.SalesChannel == "" {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "salesChannel is required."}
	}
	return validateBrand(input.Brand)
}

func validateBrand(brand string) error {
	if brand != "exito" && brand != "carulla" {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: fmt.Sprintf("Unsupported VTEX brand %q.", brand)}
	}
	return nil
}

func normalizeBrandInput(brand string) string {
	normalized := strings.ToLower(strings.TrimSpace(brand))
	if normalized == "" {
		return "exito"
	}
	return normalized
}

func normalizedBrand(brand string) string {
	return normalizeBrandInput(brand)
}
