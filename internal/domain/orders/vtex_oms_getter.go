package orders

import (
	"context"
	"net/http"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

// VTEXOMSGetterConfig contains the provider settings required by VTEXOMSGetter.
type VTEXOMSGetterConfig struct {
	BaseURL  string
	AppKey   string
	AppToken string
}

// VTEXOMSGetter calls VTEX OMS and maps its DTO to domain output.
type VTEXOMSGetter struct {
	baseURL  string
	appKey   string
	appToken string
	client   httpclient.Client
}

// NewVTEXOMSGetter creates a VTEX OMS-backed getter.
func NewVTEXOMSGetter(config VTEXOMSGetterConfig, client *http.Client) VTEXOMSGetter {
	return VTEXOMSGetter{
		baseURL:  strings.TrimSpace(config.BaseURL),
		appKey:   strings.TrimSpace(config.AppKey),
		appToken: strings.TrimSpace(config.AppToken),
		client: httpclient.New(httpclient.Config{
			BaseURL: config.BaseURL,
			Client:  client,
		}),
	}
}

// GetVTEXOMS calls VTEX OMS order detail by order identifier.
func (g VTEXOMSGetter) GetVTEXOMS(ctx context.Context, input GetVTEXOMSInput) (VTEXOMSOrder, error) {
	if strings.TrimSpace(g.baseURL) == "" || strings.TrimSpace(g.appKey) == "" || strings.TrimSpace(g.appToken) == "" {
		return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrdersNotConfigured, Message: "VTEX OMS client is not configured."}
	}

	request, err := g.client.NewRequest(ctx, http.MethodGet, "/api/oms/pvt/orders/"+strings.TrimSpace(input.ID), nil)
	if err != nil {
		return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrdersNotConfigured, Message: "VTEX OMS provider base URL is invalid."}
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-VTEX-API-AppKey", g.appKey)
	request.Header.Set("X-VTEX-API-AppToken", g.appToken)

	response, err := g.client.Do(request)
	if err != nil {
		return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrdersProviderUnavailable, Message: "VTEX OMS provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusNotFound {
		return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrderNotFound, Message: "Order not found."}
	}
	if !httpclient.Successful(response) {
		return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrdersProviderUnavailable, Message: "VTEX OMS provider returned an unsuccessful response."}
	}

	var providerResponse map[string]any
	if err := httpclient.DecodeJSONResponse(response, &providerResponse); err != nil {
		return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrdersProviderInvalidResponse, Message: "VTEX OMS provider returned an invalid response."}
	}

	order := vtexOMSOrderFromMap(providerResponse, input.Brand)
	if order.ID == "" {
		return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrdersProviderInvalidResponse, Message: "VTEX OMS provider returned an invalid response."}
	}
	return order, nil
}

func vtexOMSOrderFromMap(fields map[string]any, brand string) VTEXOMSOrder {
	return VTEXOMSOrder{
		ID:                firstString(fields, "orderId", "id", "orderID"),
		Sequence:          firstString(fields, "sequence"),
		Status:            firstString(fields, "status"),
		StatusDescription: firstString(fields, "statusDescription"),
		CreationDate:      firstString(fields, "creationDate", "createdAt"),
		ClientName:        vtexOMSClientName(fields),
		Email:             vtexOMSEmail(fields),
		TotalValue:        firstFloat(fields, "value", "totalValue"),
		Brand:             normalizedVTEXOMSBrand(brand),
		Details:           fields,
	}
}

func vtexOMSClientName(fields map[string]any) string {
	clientProfileData, ok := fields["clientProfileData"].(map[string]any)
	if !ok {
		return ""
	}
	firstName := firstString(clientProfileData, "firstName")
	lastName := firstString(clientProfileData, "lastName")
	return strings.TrimSpace(firstName + " " + lastName)
}

func vtexOMSEmail(fields map[string]any) string {
	clientProfileData, ok := fields["clientProfileData"].(map[string]any)
	if !ok {
		return firstString(fields, "email")
	}
	return firstString(clientProfileData, "email")
}
