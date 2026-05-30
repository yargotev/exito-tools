package checkout

import (
	"context"
	"fmt"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
)

const (
	CapabilityGetOrderFormID        = "checkout.get-order-form"
	CapabilityCreateOrderFormID     = "checkout.create-order-form"
	CapabilityAddItemsID            = "checkout.add-items"
	CapabilityUpdateClientProfileID = "checkout.update-client-profile"
	CapabilityUpdateShippingDataID  = "checkout.update-shipping-data"
	DomainName                      = "checkout"

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

type AddItemsInput struct {
	Brand       string
	OrderFormID string
	Items       []AddItemInput
}

type UpdateClientProfileInput struct {
	Brand         string
	OrderFormID   string
	ClientProfile ClientProfileInput
}

type UpdateShippingDataInput struct {
	Brand        string
	OrderFormID  string
	ShippingData ShippingDataInput
}

type ClientProfileInput struct {
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	DocumentType string `json:"documentType"`
	Document     string `json:"document"`
	Phone        string `json:"phone"`
}

type ShippingDataInput struct {
	SelectedAddresses                []ShippingAddressInput `json:"selectedAddresses"`
	LogisticsInfo                    []LogisticsInfoInput   `json:"logisticsInfo"`
	ClearAddressIfPostalCodeNotFound bool                   `json:"clearAddressIfPostalCodeNotFound"`
}

type ShippingAddressInput struct {
	AddressType    string    `json:"addressType,omitempty"`
	ReceiverName   string    `json:"receiverName,omitempty"`
	PostalCode     string    `json:"postalCode,omitempty"`
	City           string    `json:"city,omitempty"`
	State          string    `json:"state,omitempty"`
	Country        string    `json:"country,omitempty"`
	Street         string    `json:"street,omitempty"`
	Number         string    `json:"number,omitempty"`
	Neighborhood   string    `json:"neighborhood,omitempty"`
	Complement     string    `json:"complement,omitempty"`
	Reference      string    `json:"reference,omitempty"`
	GeoCoordinates []float64 `json:"geoCoordinates,omitempty"`
}

type LogisticsInfoInput struct {
	ItemIndex               int    `json:"itemIndex"`
	SelectedSLA             string `json:"selectedSla"`
	SelectedDeliveryChannel string `json:"selectedDeliveryChannel"`
}

type AddItemInput struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Seller   string `json:"seller,omitempty"`
}

type OrderFormSummary struct {
	Brand                string               `json:"brand"`
	ID                   string               `json:"id"`
	SalesChannel         string               `json:"salesChannel,omitempty"`
	Value                int                  `json:"value"`
	Totalizers           []TotalizerSummary   `json:"totalizers,omitempty"`
	Items                []ItemSummary        `json:"items,omitempty"`
	ItemCount            int                  `json:"itemCount"`
	ClientProfileDataSet bool                 `json:"clientProfileDataSet"`
	ShippingDataSet      bool                 `json:"shippingDataSet"`
	ShippingTotal        int                  `json:"shippingTotal,omitempty"`
	SelectedSLAs         []SelectedSLASummary `json:"selectedSlas,omitempty"`
	Diagnostics          Diagnostics          `json:"diagnostics,omitempty"`
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

type SelectedSLASummary struct {
	ItemIndex               int    `json:"itemIndex"`
	SelectedSLA             string `json:"selectedSla,omitempty"`
	SelectedDeliveryChannel string `json:"selectedDeliveryChannel,omitempty"`
	Price                   int    `json:"price,omitempty"`
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

type AddItemsResult struct {
	OrderForm OrderFormSummary `json:"orderForm"`
	Items     []AddItemInput   `json:"items"`
}

type UpdateClientProfileResult struct {
	OrderForm OrderFormSummary `json:"orderForm"`
}

type UpdateShippingDataResult struct {
	OrderForm OrderFormSummary `json:"orderForm"`
}

type Getter interface {
	GetOrderForm(ctx context.Context, input GetOrderFormInput) (OrderFormSummary, error)
}

type Creator interface {
	CreateOrderForm(ctx context.Context, input CreateOrderFormInput) (OrderFormSummary, error)
}

type Adder interface {
	AddItems(ctx context.Context, input AddItemsInput) (OrderFormSummary, error)
}

type ClientProfileUpdater interface {
	UpdateClientProfile(ctx context.Context, input UpdateClientProfileInput) (OrderFormSummary, error)
}

type ShippingDataUpdater interface {
	UpdateShippingData(ctx context.Context, input UpdateShippingDataInput) (OrderFormSummary, error)
}

type (
	GetOrderFormUseCase        struct{ getter Getter }
	CreateOrderFormUseCase     struct{ creator Creator }
	AddItemsUseCase            struct{ adder Adder }
	UpdateClientProfileUseCase struct{ updater ClientProfileUpdater }
	UpdateShippingDataUseCase  struct{ updater ShippingDataUpdater }
)

func NewGetOrderFormUseCase(getter Getter) GetOrderFormUseCase {
	return GetOrderFormUseCase{getter: getter}
}

func NewCreateOrderFormUseCase(creator Creator) CreateOrderFormUseCase {
	return CreateOrderFormUseCase{creator: creator}
}

func NewAddItemsUseCase(adder Adder) AddItemsUseCase {
	return AddItemsUseCase{adder: adder}
}

func NewUpdateClientProfileUseCase(updater ClientProfileUpdater) UpdateClientProfileUseCase {
	return UpdateClientProfileUseCase{updater: updater}
}

func NewUpdateShippingDataUseCase(updater ShippingDataUpdater) UpdateShippingDataUseCase {
	return UpdateShippingDataUseCase{updater: updater}
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

func (u AddItemsUseCase) Execute(ctx context.Context, input AddItemsInput) (AddItemsResult, error) {
	if u.adder == nil {
		return AddItemsResult{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
	}
	input = normalizeAddItemsInput(input)
	if err := validateAddItemsInput(input); err != nil {
		return AddItemsResult{}, err
	}
	orderForm, err := u.adder.AddItems(ctx, input)
	if err != nil {
		return AddItemsResult{}, err
	}
	return AddItemsResult{OrderForm: orderForm, Items: input.Items}, nil
}

func (u UpdateClientProfileUseCase) Execute(ctx context.Context, input UpdateClientProfileInput) (UpdateClientProfileResult, error) {
	if u.updater == nil {
		return UpdateClientProfileResult{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
	}
	input = normalizeUpdateClientProfileInput(input)
	if err := validateUpdateClientProfileInput(input); err != nil {
		return UpdateClientProfileResult{}, err
	}
	orderForm, err := u.updater.UpdateClientProfile(ctx, input)
	if err != nil {
		return UpdateClientProfileResult{}, err
	}
	return UpdateClientProfileResult{OrderForm: orderForm}, nil
}

func (u UpdateShippingDataUseCase) Execute(ctx context.Context, input UpdateShippingDataInput) (UpdateShippingDataResult, error) {
	if u.updater == nil {
		return UpdateShippingDataResult{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
	}
	input = normalizeUpdateShippingDataInput(input)
	if err := validateUpdateShippingDataInput(input); err != nil {
		return UpdateShippingDataResult{}, err
	}
	orderForm, err := u.updater.UpdateShippingData(ctx, input)
	if err != nil {
		return UpdateShippingDataResult{}, err
	}
	return UpdateShippingDataResult{OrderForm: orderForm}, nil
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

func NewAddItemsCapability(adder Adder) capability.Executable {
	useCase := NewAddItemsUseCase(adder)
	return capability.Executable{Definition: AddItemsDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, addItemsInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		return capability.ExecutionResult{Data: result}, nil
	}}
}

func NewUpdateClientProfileCapability(updater ClientProfileUpdater) capability.Executable {
	useCase := NewUpdateClientProfileUseCase(updater)
	return capability.Executable{Definition: UpdateClientProfileDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, updateClientProfileInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		return capability.ExecutionResult{Data: result}, nil
	}}
}

func NewUpdateShippingDataCapability(updater ShippingDataUpdater) capability.Executable {
	useCase := NewUpdateShippingDataUseCase(updater)
	return capability.Executable{Definition: UpdateShippingDataDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, updateShippingDataInputFromCapability(request.Input))
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

func AddItemsDefinition() capability.Definition {
	return capability.Definition{ID: CapabilityAddItemsID, Domain: DomainName, Version: "1.0.0", Title: "Add items to VTEX Checkout orderForm", Description: "Adds selected SKU items to an existing VTEX Checkout orderForm.", Risk: capability.RiskSafeWrite, RequiresConfirmation: true, Audiences: []capability.Audience{capability.AudienceAgents, capability.AudiencePeople}, Visibility: []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette}, InputSchema: &capability.InputSchema{Fields: []capability.InputField{{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."}, {Name: "orderFormId", Type: capability.InputTypeString, Required: true, Description: "VTEX Checkout orderForm identifier."}, {Name: "items", Type: capability.InputTypeArray, Required: true, Description: "Items to add. Each item contains sku, quantity, and optional seller."}}}}
}

func UpdateClientProfileDefinition() capability.Definition {
	return capability.Definition{ID: CapabilityUpdateClientProfileID, Domain: DomainName, Version: "1.0.0", Title: "Update VTEX Checkout client profile", Description: "Attaches customer client profile data to an existing VTEX Checkout orderForm and returns a redacted summary.", Risk: capability.RiskSafeWrite, RequiresConfirmation: true, Audiences: []capability.Audience{capability.AudienceAgents, capability.AudiencePeople}, Visibility: []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette}, InputSchema: &capability.InputSchema{Fields: []capability.InputField{{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."}, {Name: "orderFormId", Type: capability.InputTypeString, Required: true, Description: "VTEX Checkout orderForm identifier."}, {Name: "clientProfile", Type: capability.InputTypeObject, Required: true, Description: "Client profile object with email, firstName, lastName, documentType, document, and phone."}}}}
}

func UpdateShippingDataDefinition() capability.Definition {
	return capability.Definition{ID: CapabilityUpdateShippingDataID, Domain: DomainName, Version: "1.0.0", Title: "Update VTEX Checkout shipping data", Description: "Attaches delivery address data and selected logistics options to an existing VTEX Checkout orderForm and returns a redacted summary.", Risk: capability.RiskSafeWrite, RequiresConfirmation: true, Audiences: []capability.Audience{capability.AudienceAgents, capability.AudiencePeople}, Visibility: []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette}, InputSchema: &capability.InputSchema{Fields: []capability.InputField{{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."}, {Name: "orderFormId", Type: capability.InputTypeString, Required: true, Description: "VTEX Checkout orderForm identifier."}, {Name: "shippingData", Type: capability.InputTypeObject, Required: true, Description: "Shipping data object with selectedAddresses and logisticsInfo."}}}}
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

func addItemsInputFromCapability(input capability.Input) AddItemsInput {
	out := AddItemsInput{}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["orderFormId"].(string); ok {
		out.OrderFormID = value
	}
	out.Items = addItemsFromCapabilityValue(input["items"])
	return out
}

func updateClientProfileInputFromCapability(input capability.Input) UpdateClientProfileInput {
	out := UpdateClientProfileInput{}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["orderFormId"].(string); ok {
		out.OrderFormID = value
	}
	out.ClientProfile = clientProfileFromCapabilityValue(input["clientProfile"])
	return out
}

func updateShippingDataInputFromCapability(input capability.Input) UpdateShippingDataInput {
	out := UpdateShippingDataInput{}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["orderFormId"].(string); ok {
		out.OrderFormID = value
	}
	out.ShippingData = shippingDataFromCapabilityValue(input["shippingData"])
	return out
}

func addItemsFromCapabilityValue(value any) []AddItemInput {
	switch items := value.(type) {
	case []AddItemInput:
		return append([]AddItemInput(nil), items...)
	case []any:
		out := make([]AddItemInput, 0, len(items))
		for _, item := range items {
			if mapped, ok := item.(map[string]any); ok {
				out = append(out, addItemFromMap(mapped))
			}
		}
		return out
	default:
		return nil
	}
}

func addItemFromMap(mapped map[string]any) AddItemInput {
	item := AddItemInput{}
	if value, ok := mapped["sku"].(string); ok {
		item.SKU = value
	}
	if value, ok := mapped["seller"].(string); ok {
		item.Seller = value
	}
	switch value := mapped["quantity"].(type) {
	case int:
		item.Quantity = value
	case int64:
		item.Quantity = int(value)
	case float64:
		if value == float64(int(value)) {
			item.Quantity = int(value)
		}
	}
	return item
}

func clientProfileFromCapabilityValue(value any) ClientProfileInput {
	switch profile := value.(type) {
	case ClientProfileInput:
		return profile
	case map[string]any:
		return clientProfileFromMap(profile)
	case capability.Input:
		return clientProfileFromMap(map[string]any(profile))
	default:
		return ClientProfileInput{}
	}
}

func clientProfileFromMap(mapped map[string]any) ClientProfileInput {
	profile := ClientProfileInput{}
	if value, ok := mapped["email"].(string); ok {
		profile.Email = value
	}
	if value, ok := mapped["firstName"].(string); ok {
		profile.FirstName = value
	}
	if value, ok := mapped["lastName"].(string); ok {
		profile.LastName = value
	}
	if value, ok := mapped["documentType"].(string); ok {
		profile.DocumentType = value
	}
	if value, ok := mapped["document"].(string); ok {
		profile.Document = value
	}
	if value, ok := mapped["phone"].(string); ok {
		profile.Phone = value
	}
	return profile
}

func shippingDataFromCapabilityValue(value any) ShippingDataInput {
	switch shippingData := value.(type) {
	case ShippingDataInput:
		return shippingData
	case map[string]any:
		return shippingDataFromMap(shippingData)
	case capability.Input:
		return shippingDataFromMap(map[string]any(shippingData))
	default:
		return ShippingDataInput{}
	}
}

func shippingDataFromMap(mapped map[string]any) ShippingDataInput {
	shippingData := ShippingDataInput{}
	if value, ok := mapped["clearAddressIfPostalCodeNotFound"].(bool); ok {
		shippingData.ClearAddressIfPostalCodeNotFound = value
	}
	if values, ok := mapped["selectedAddresses"].([]any); ok {
		shippingData.SelectedAddresses = make([]ShippingAddressInput, 0, len(values))
		for _, value := range values {
			if address, ok := value.(map[string]any); ok {
				shippingData.SelectedAddresses = append(shippingData.SelectedAddresses, shippingAddressFromMap(address))
			}
		}
	}
	if values, ok := mapped["logisticsInfo"].([]any); ok {
		shippingData.LogisticsInfo = make([]LogisticsInfoInput, 0, len(values))
		for _, value := range values {
			if logistics, ok := value.(map[string]any); ok {
				shippingData.LogisticsInfo = append(shippingData.LogisticsInfo, logisticsInfoFromMap(logistics))
			}
		}
	}
	return shippingData
}

func shippingAddressFromMap(mapped map[string]any) ShippingAddressInput {
	address := ShippingAddressInput{}
	if value, ok := mapped["addressType"].(string); ok {
		address.AddressType = value
	}
	if value, ok := mapped["receiverName"].(string); ok {
		address.ReceiverName = value
	}
	if value, ok := mapped["postalCode"].(string); ok {
		address.PostalCode = value
	}
	if value, ok := mapped["city"].(string); ok {
		address.City = value
	}
	if value, ok := mapped["state"].(string); ok {
		address.State = value
	}
	if value, ok := mapped["country"].(string); ok {
		address.Country = value
	}
	if value, ok := mapped["street"].(string); ok {
		address.Street = value
	}
	if value, ok := mapped["number"].(string); ok {
		address.Number = value
	}
	if value, ok := mapped["neighborhood"].(string); ok {
		address.Neighborhood = value
	}
	if value, ok := mapped["complement"].(string); ok {
		address.Complement = value
	}
	if value, ok := mapped["reference"].(string); ok {
		address.Reference = value
	}
	if values, ok := mapped["geoCoordinates"].([]any); ok {
		address.GeoCoordinates = float64SliceFromCapabilityValues(values)
	}
	return address
}

func logisticsInfoFromMap(mapped map[string]any) LogisticsInfoInput {
	logistics := LogisticsInfoInput{}
	switch value := mapped["itemIndex"].(type) {
	case int:
		logistics.ItemIndex = value
	case int64:
		logistics.ItemIndex = int(value)
	case float64:
		if value == float64(int(value)) {
			logistics.ItemIndex = int(value)
		}
	}
	if value, ok := mapped["selectedSla"].(string); ok {
		logistics.SelectedSLA = value
	}
	if value, ok := mapped["selectedDeliveryChannel"].(string); ok {
		logistics.SelectedDeliveryChannel = value
	}
	return logistics
}

func float64SliceFromCapabilityValues(values []any) []float64 {
	out := make([]float64, 0, len(values))
	for _, value := range values {
		switch number := value.(type) {
		case float64:
			out = append(out, number)
		case int:
			out = append(out, float64(number))
		case int64:
			out = append(out, float64(number))
		}
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

func validateAddItemsInput(input AddItemsInput) error {
	if input.OrderFormID == "" {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "orderFormId is required."}
	}
	if len(input.Items) == 0 {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "At least one item is required."}
	}
	for index, item := range input.Items {
		if item.SKU == "" {
			return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: fmt.Sprintf("items[%d].sku is required.", index)}
		}
		if item.Quantity <= 0 {
			return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: fmt.Sprintf("items[%d].quantity must be greater than zero.", index)}
		}
		if item.Seller == "" {
			return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: fmt.Sprintf("items[%d].seller is required.", index)}
		}
	}
	return validateBrand(input.Brand)
}

func normalizeAddItemsInput(input AddItemsInput) AddItemsInput {
	input.Brand = normalizeBrandInput(input.Brand)
	input.OrderFormID = strings.TrimSpace(input.OrderFormID)
	for index := range input.Items {
		input.Items[index].SKU = strings.TrimSpace(input.Items[index].SKU)
		input.Items[index].Seller = strings.TrimSpace(input.Items[index].Seller)
	}
	return input
}

func validateUpdateClientProfileInput(input UpdateClientProfileInput) error {
	if input.OrderFormID == "" {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "orderFormId is required."}
	}
	required := []struct {
		field string
		value string
	}{
		{field: "clientProfile.email", value: input.ClientProfile.Email},
		{field: "clientProfile.firstName", value: input.ClientProfile.FirstName},
		{field: "clientProfile.lastName", value: input.ClientProfile.LastName},
		{field: "clientProfile.documentType", value: input.ClientProfile.DocumentType},
		{field: "clientProfile.document", value: input.ClientProfile.Document},
		{field: "clientProfile.phone", value: input.ClientProfile.Phone},
	}
	for _, item := range required {
		if item.value == "" {
			return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: item.field + " is required."}
		}
	}
	return validateBrand(input.Brand)
}

func normalizeUpdateClientProfileInput(input UpdateClientProfileInput) UpdateClientProfileInput {
	input.Brand = normalizeBrandInput(input.Brand)
	input.OrderFormID = strings.TrimSpace(input.OrderFormID)
	input.ClientProfile.Email = strings.TrimSpace(input.ClientProfile.Email)
	input.ClientProfile.FirstName = strings.TrimSpace(input.ClientProfile.FirstName)
	input.ClientProfile.LastName = strings.TrimSpace(input.ClientProfile.LastName)
	input.ClientProfile.DocumentType = strings.TrimSpace(input.ClientProfile.DocumentType)
	input.ClientProfile.Document = strings.TrimSpace(input.ClientProfile.Document)
	input.ClientProfile.Phone = strings.TrimSpace(input.ClientProfile.Phone)
	return input
}

func validateUpdateShippingDataInput(input UpdateShippingDataInput) error {
	if input.OrderFormID == "" {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "orderFormId is required."}
	}
	if len(input.ShippingData.SelectedAddresses) == 0 {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "shippingData.selectedAddresses is required."}
	}
	if len(input.ShippingData.LogisticsInfo) == 0 {
		return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "shippingData.logisticsInfo is required."}
	}
	for index, address := range input.ShippingData.SelectedAddresses {
		required := []struct {
			field string
			value string
		}{
			{field: fmt.Sprintf("shippingData.selectedAddresses[%d].addressType", index), value: address.AddressType},
			{field: fmt.Sprintf("shippingData.selectedAddresses[%d].receiverName", index), value: address.ReceiverName},
			{field: fmt.Sprintf("shippingData.selectedAddresses[%d].country", index), value: address.Country},
			{field: fmt.Sprintf("shippingData.selectedAddresses[%d].city", index), value: address.City},
			{field: fmt.Sprintf("shippingData.selectedAddresses[%d].state", index), value: address.State},
			{field: fmt.Sprintf("shippingData.selectedAddresses[%d].street", index), value: address.Street},
			{field: fmt.Sprintf("shippingData.selectedAddresses[%d].number", index), value: address.Number},
			{field: fmt.Sprintf("shippingData.selectedAddresses[%d].neighborhood", index), value: address.Neighborhood},
		}
		for _, item := range required {
			if item.value == "" {
				return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: item.field + " is required."}
			}
		}
		if len(address.GeoCoordinates) != 0 && len(address.GeoCoordinates) != 2 {
			return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: fmt.Sprintf("shippingData.selectedAddresses[%d].geoCoordinates must contain longitude and latitude.", index)}
		}
	}
	for index, logistics := range input.ShippingData.LogisticsInfo {
		if logistics.ItemIndex < 0 {
			return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: fmt.Sprintf("shippingData.logisticsInfo[%d].itemIndex must be zero or greater.", index)}
		}
		if logistics.SelectedSLA == "" {
			return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: fmt.Sprintf("shippingData.logisticsInfo[%d].selectedSla is required.", index)}
		}
		if logistics.SelectedDeliveryChannel == "" {
			return capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: fmt.Sprintf("shippingData.logisticsInfo[%d].selectedDeliveryChannel is required.", index)}
		}
	}
	return validateBrand(input.Brand)
}

func normalizeUpdateShippingDataInput(input UpdateShippingDataInput) UpdateShippingDataInput {
	input.Brand = normalizeBrandInput(input.Brand)
	input.OrderFormID = strings.TrimSpace(input.OrderFormID)
	for index := range input.ShippingData.SelectedAddresses {
		address := &input.ShippingData.SelectedAddresses[index]
		address.AddressType = strings.TrimSpace(address.AddressType)
		address.ReceiverName = strings.TrimSpace(address.ReceiverName)
		address.PostalCode = strings.TrimSpace(address.PostalCode)
		address.City = strings.TrimSpace(address.City)
		address.State = strings.TrimSpace(address.State)
		address.Country = strings.TrimSpace(address.Country)
		address.Street = strings.TrimSpace(address.Street)
		address.Number = strings.TrimSpace(address.Number)
		address.Neighborhood = strings.TrimSpace(address.Neighborhood)
		address.Complement = strings.TrimSpace(address.Complement)
		address.Reference = strings.TrimSpace(address.Reference)
	}
	for index := range input.ShippingData.LogisticsInfo {
		input.ShippingData.LogisticsInfo[index].SelectedSLA = strings.TrimSpace(input.ShippingData.LogisticsInfo[index].SelectedSLA)
		input.ShippingData.LogisticsInfo[index].SelectedDeliveryChannel = strings.TrimSpace(input.ShippingData.LogisticsInfo[index].SelectedDeliveryChannel)
	}
	return input
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
