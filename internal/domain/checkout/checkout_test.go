package checkout_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/domain/checkout"
)

type recordingCheckoutClient struct {
	getInput    checkout.GetOrderFormInput
	createInput checkout.CreateOrderFormInput
	addInput    checkout.AddItemsInput
}

func (c *recordingCheckoutClient) GetOrderForm(_ context.Context, input checkout.GetOrderFormInput) (checkout.OrderFormSummary, error) {
	c.getInput = input
	return checkout.OrderFormSummary{Brand: input.Brand, ID: input.OrderFormID}, nil
}

func (c *recordingCheckoutClient) CreateOrderForm(_ context.Context, input checkout.CreateOrderFormInput) (checkout.OrderFormSummary, error) {
	c.createInput = input
	return checkout.OrderFormSummary{Brand: input.Brand, SalesChannel: input.SalesChannel}, nil
}

func (c *recordingCheckoutClient) AddItems(_ context.Context, input checkout.AddItemsInput) (checkout.OrderFormSummary, error) {
	c.addInput = input
	return checkout.OrderFormSummary{Brand: input.Brand, ID: input.OrderFormID, Items: []checkout.ItemSummary{{ID: input.Items[0].SKU, Quantity: input.Items[0].Quantity, Seller: input.Items[0].Seller}}, ItemCount: len(input.Items)}, nil
}

func TestUseCasesDefaultBlankBrandToExito(t *testing.T) {
	t.Parallel()

	client := &recordingCheckoutClient{}
	_, err := checkout.NewGetOrderFormUseCase(client).Execute(context.Background(), checkout.GetOrderFormInput{OrderFormID: "of-1"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if client.getInput.Brand != "exito" {
		t.Fatalf("brand = %q, want exito", client.getInput.Brand)
	}
}

func TestUseCasesRejectUnsupportedBrand(t *testing.T) {
	t.Parallel()

	client := &recordingCheckoutClient{}
	_, err := checkout.NewCreateOrderFormUseCase(client).Execute(context.Background(), checkout.CreateOrderFormInput{Brand: "unknown", SalesChannel: "1"})
	if err == nil {
		t.Fatalf("Execute() error = nil, want structured invalid input")
	}
	var structured capability.StructuredError
	if !errors.As(err, &structured) {
		t.Fatalf("error = %T, want StructuredError", err)
	}
	if structured.Code != checkout.ErrorCheckoutInvalidInput {
		t.Fatalf("code = %q, want %s", structured.Code, checkout.ErrorCheckoutInvalidInput)
	}
	if client.createInput.Brand != "" {
		t.Fatalf("client was called with %#v", client.createInput)
	}
}

func TestAddItemsUseCaseNormalizesAndAddsItems(t *testing.T) {
	t.Parallel()

	client := &recordingCheckoutClient{}
	got, err := checkout.NewAddItemsUseCase(client).Execute(context.Background(), checkout.AddItemsInput{
		Brand:       " EXITO ",
		OrderFormID: " of-1 ",
		Items:       []checkout.AddItemInput{{SKU: " sku-1 ", Quantity: 2, Seller: " 1 "}},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if client.addInput.Brand != "exito" || client.addInput.OrderFormID != "of-1" {
		t.Fatalf("input = %#v, want normalized brand/orderForm", client.addInput)
	}
	if client.addInput.Items[0].SKU != "sku-1" || client.addInput.Items[0].Seller != "1" {
		t.Fatalf("item = %#v, want trimmed SKU/seller", client.addInput.Items[0])
	}
	if got.OrderForm.ItemCount != 1 || got.Items[0].SKU != "sku-1" {
		t.Fatalf("result = %#v, want updated orderForm and item diagnostics", got)
	}
}

func TestAddItemsUseCaseRejectsEmptyItems(t *testing.T) {
	t.Parallel()

	client := &recordingCheckoutClient{}
	_, err := checkout.NewAddItemsUseCase(client).Execute(context.Background(), checkout.AddItemsInput{Brand: "exito", OrderFormID: "of-1"})
	if err == nil {
		t.Fatalf("Execute() error = nil, want structured invalid input")
	}
	var structured capability.StructuredError
	if !errors.As(err, &structured) {
		t.Fatalf("error = %T, want StructuredError", err)
	}
	if structured.Code != checkout.ErrorCheckoutInvalidInput {
		t.Fatalf("code = %q, want %s", structured.Code, checkout.ErrorCheckoutInvalidInput)
	}
	if client.addInput.OrderFormID != "" {
		t.Fatalf("client was called with %#v", client.addInput)
	}
}
