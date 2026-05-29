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
}

func (c *recordingCheckoutClient) GetOrderForm(_ context.Context, input checkout.GetOrderFormInput) (checkout.OrderFormSummary, error) {
	c.getInput = input
	return checkout.OrderFormSummary{Brand: input.Brand, ID: input.OrderFormID}, nil
}

func (c *recordingCheckoutClient) CreateOrderForm(_ context.Context, input checkout.CreateOrderFormInput) (checkout.OrderFormSummary, error) {
	c.createInput = input
	return checkout.OrderFormSummary{Brand: input.Brand, SalesChannel: input.SalesChannel}, nil
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
