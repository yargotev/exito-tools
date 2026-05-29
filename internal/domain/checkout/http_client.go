package checkout

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

const checkoutOrderFormPath = "/api/checkout/pub/orderForm"

type HTTPClientConfig struct {
	BaseURL string
}

type HTTPClient struct {
	baseURL string
	client  httpclient.Client
}

func NewHTTPClient(config HTTPClientConfig, client *http.Client) HTTPClient {
	return HTTPClient{baseURL: strings.TrimSpace(config.BaseURL), client: httpclient.New(httpclient.Config{BaseURL: config.BaseURL, Client: client})}
}

func (c HTTPClient) GetOrderForm(ctx context.Context, input GetOrderFormInput) (OrderFormSummary, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
	}
	path := checkoutOrderFormPath + "/" + url.PathEscape(input.OrderFormID)
	request, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout provider base URL is invalid."}
	}
	return c.doOrderForm(request, input.Brand, path)
}

func (c HTTPClient) CreateOrderForm(ctx context.Context, input CreateOrderFormInput) (OrderFormSummary, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
	}
	path := checkoutOrderFormPath
	request, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout provider base URL is invalid."}
	}
	query := request.URL.Query()
	query.Set("forceNewCart", "true")
	query.Set("sc", input.SalesChannel)
	request.URL.RawQuery = query.Encode()
	return c.doOrderForm(request, input.Brand, path+"?forceNewCart=true&sc="+url.QueryEscape(input.SalesChannel))
}

func (c HTTPClient) AddItems(ctx context.Context, input AddItemsInput) (OrderFormSummary, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
	}
	path := checkoutOrderFormPath + "/" + url.PathEscape(input.OrderFormID) + "/items"
	body, err := json.Marshal(addItemsRequestFromInput(input))
	if err != nil {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutInvalidInput, Message: "Checkout add-items request is invalid."}
	}
	request, err := c.client.NewRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout provider base URL is invalid."}
	}
	request.Header.Set("Content-Type", "application/json")
	query := request.URL.Query()
	query.Set("allowedOutdatedData", "false")
	request.URL.RawQuery = query.Encode()
	return c.doOrderForm(request, input.Brand, path+"?allowedOutdatedData=false")
}

func (c HTTPClient) doOrderForm(request *http.Request, brand string, requestPath string) (OrderFormSummary, error) {
	response, err := c.client.Do(request)
	if err != nil {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutProviderUnavailable, Message: "VTEX Checkout provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()
	if !httpclient.Successful(response) {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutProviderUnavailable, Message: "VTEX Checkout provider returned an unsuccessful response."}
	}
	var payload orderFormDTO
	if err := httpclient.DecodeJSONResponse(response, &payload); err != nil {
		return OrderFormSummary{}, capability.StructuredError{Code: ErrorCheckoutProviderInvalidResponse, Message: "VTEX Checkout provider returned an invalid response."}
	}
	return mapOrderFormDTO(payload, brand, requestPath, response.StatusCode), nil
}

type orderFormDTO struct {
	ID                string         `json:"orderFormId"`
	SalesChannel      string         `json:"salesChannel"`
	Value             int            `json:"value"`
	Totalizers        []totalizerDTO `json:"totalizers"`
	Items             []itemDTO      `json:"items"`
	ClientProfileData any            `json:"clientProfileData"`
	ShippingData      any            `json:"shippingData"`
}

type totalizerDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}
type itemDTO struct {
	ID           string `json:"id"`
	ProductID    string `json:"productId"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	Seller       string `json:"seller"`
	Price        int    `json:"price"`
	SellingPrice int    `json:"sellingPrice"`
	Availability string `json:"availability"`
}

type addItemsRequestDTO struct {
	OrderItems []addItemRequestDTO `json:"orderItems"`
}

type addItemRequestDTO struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Seller   string `json:"seller"`
	Index    int    `json:"index"`
}

func addItemsRequestFromInput(input AddItemsInput) addItemsRequestDTO {
	items := make([]addItemRequestDTO, 0, len(input.Items))
	for index, item := range input.Items {
		items = append(items, addItemRequestDTO{ID: item.SKU, Quantity: item.Quantity, Seller: item.Seller, Index: index})
	}
	return addItemsRequestDTO{OrderItems: items}
}

func mapOrderFormDTO(dto orderFormDTO, brand string, requestPath string, status int) OrderFormSummary {
	items := make([]ItemSummary, 0, len(dto.Items))
	for _, item := range dto.Items {
		items = append(items, ItemSummary(item))
	}
	totalizers := make([]TotalizerSummary, 0, len(dto.Totalizers))
	for _, totalizer := range dto.Totalizers {
		totalizers = append(totalizers, TotalizerSummary(totalizer))
	}
	return OrderFormSummary{Brand: normalizedBrand(brand), ID: dto.ID, SalesChannel: dto.SalesChannel, Value: dto.Value, Totalizers: totalizers, Items: items, ItemCount: len(items), ClientProfileDataSet: isSet(dto.ClientProfileData), ShippingDataSet: isSet(dto.ShippingData), Diagnostics: Diagnostics{RequestPath: requestPath, ProviderStatus: status}}
}

func isSet(value any) bool {
	if value == nil {
		return false
	}
	if mapped, ok := value.(map[string]any); ok {
		return len(mapped) > 0
	}
	return true
}
