package checkout

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

type BrandClient struct {
	exito interface {
		Getter
		Creator
	}
	carulla interface {
		Getter
		Creator
	}
}

func NewBrandClient(exito interface {
	Getter
	Creator
}, carulla interface {
	Getter
	Creator
},
) BrandClient {
	return BrandClient{exito: exito, carulla: carulla}
}

func (c BrandClient) GetOrderForm(ctx context.Context, input GetOrderFormInput) (OrderFormSummary, error) {
	return c.client(input.Brand).GetOrderForm(ctx, input)
}

func (c BrandClient) CreateOrderForm(ctx context.Context, input CreateOrderFormInput) (OrderFormSummary, error) {
	return c.client(input.Brand).CreateOrderForm(ctx, input)
}

func (c BrandClient) client(brand string) interface {
	Getter
	Creator
} {
	if normalizedBrand(brand) == "carulla" {
		if c.carulla != nil {
			return c.carulla
		}
		return UnavailableClient{}
	}
	if c.exito != nil {
		return c.exito
	}
	return UnavailableClient{}
}

type UnavailableClient struct{}

func (UnavailableClient) GetOrderForm(context.Context, GetOrderFormInput) (OrderFormSummary, error) {
	return OrderFormSummary{}, notConfiguredError()
}

func (UnavailableClient) CreateOrderForm(context.Context, CreateOrderFormInput) (OrderFormSummary, error) {
	return OrderFormSummary{}, notConfiguredError()
}

func notConfiguredError() error {
	return capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
}
