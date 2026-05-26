package orders

import (
	"context"
	"net/http"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

// HTTPGetterConfig contains the provider settings required by HTTPGetter.
type HTTPGetterConfig struct {
	BaseURL string
	Token   string
}

// HTTPGetter calls the configured Orders provider and maps its DTO to domain output.
type HTTPGetter struct {
	client httpclient.Client
}

// NewHTTPGetter creates an Orders provider-backed getter.
func NewHTTPGetter(config HTTPGetterConfig, client *http.Client) HTTPGetter {
	return HTTPGetter{
		client: httpclient.New(httpclient.Config{
			BaseURL: config.BaseURL,
			Token:   config.Token,
			Client:  client,
		}),
	}
}

// Get calls the provider orders-get endpoint and maps the response.
func (g HTTPGetter) Get(ctx context.Context, input GetInput) (Order, error) {
	if !g.client.Configured() {
		return Order{}, capability.StructuredError{
			Code:    ErrorOrdersNotConfigured,
			Message: "Orders client is not configured.",
		}
	}

	request, err := g.client.NewJSONRequest(ctx, http.MethodPost, "/orders/get", ordersProviderRequest(input))
	if err != nil {
		return Order{}, capability.StructuredError{
			Code:    ErrorOrdersNotConfigured,
			Message: "Orders provider base URL is invalid.",
		}
	}

	response, err := g.client.Do(request)
	if err != nil {
		return Order{}, capability.StructuredError{
			Code:    ErrorOrdersProviderUnavailable,
			Message: "Orders provider request failed.",
		}
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusNotFound {
		return Order{}, capability.StructuredError{
			Code:    ErrorOrderNotFound,
			Message: "Order not found.",
		}
	}
	if !httpclient.Successful(response) {
		return Order{}, capability.StructuredError{
			Code:    ErrorOrdersProviderUnavailable,
			Message: "Orders provider returned an unsuccessful response.",
		}
	}

	var providerResponse ordersProviderResponse
	if err := httpclient.DecodeJSONResponse(response, &providerResponse); err != nil {
		return Order{}, capability.StructuredError{
			Code:    ErrorOrdersProviderInvalidResponse,
			Message: "Orders provider returned an invalid response.",
		}
	}

	return providerResponse.toDomain(), nil
}

type ordersProviderRequest struct {
	ID string `json:"id"`
}

type ordersProviderResponse struct {
	Order orderDTO `json:"order"`
}

type orderDTO struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

func (r ordersProviderResponse) toDomain() Order {
	return Order{
		ID:        r.Order.ID,
		Status:    r.Order.Status,
		CreatedAt: r.Order.CreatedAt,
	}
}
