package checkout

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

type BrandClient struct {
	exito interface {
		Getter
		Creator
		Adder
		ClientProfileUpdater
	}
	carulla interface {
		Getter
		Creator
		Adder
		ClientProfileUpdater
	}
}

func NewBrandClient(exito interface {
	Getter
	Creator
	Adder
	ClientProfileUpdater
}, carulla interface {
	Getter
	Creator
	Adder
	ClientProfileUpdater
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

func (c BrandClient) AddItems(ctx context.Context, input AddItemsInput) (OrderFormSummary, error) {
	return c.client(input.Brand).AddItems(ctx, input)
}

func (c BrandClient) UpdateClientProfile(ctx context.Context, input UpdateClientProfileInput) (OrderFormSummary, error) {
	return c.client(input.Brand).UpdateClientProfile(ctx, input)
}

func (c BrandClient) client(brand string) interface {
	Getter
	Creator
	Adder
	ClientProfileUpdater
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

func (UnavailableClient) AddItems(context.Context, AddItemsInput) (OrderFormSummary, error) {
	return OrderFormSummary{}, notConfiguredError()
}

func (UnavailableClient) UpdateClientProfile(context.Context, UpdateClientProfileInput) (OrderFormSummary, error) {
	return OrderFormSummary{}, notConfiguredError()
}

func notConfiguredError() error {
	return capability.StructuredError{Code: ErrorCheckoutNotConfigured, Message: "VTEX Checkout client is not configured."}
}
