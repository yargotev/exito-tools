package orders_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/domain/orders"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/registry"
)

func TestDefinition(t *testing.T) {
	t.Parallel()

	definition := orders.Definition()
	if definition.ID != orders.CapabilityGetID {
		t.Fatalf("ID = %q, want %q", definition.ID, orders.CapabilityGetID)
	}
	if definition.Domain != orders.DomainName {
		t.Fatalf("Domain = %q, want %q", definition.Domain, orders.DomainName)
	}
	if definition.Risk != capability.RiskReadOnly {
		t.Fatalf("Risk = %q, want %q", definition.Risk, capability.RiskReadOnly)
	}
	if len(definition.Audiences) != 2 || definition.Audiences[0] != capability.AudienceAgents || definition.Audiences[1] != capability.AudiencePeople {
		t.Fatalf("Audiences = %#v, want agents and people", definition.Audiences)
	}
	if len(definition.Visibility) != 3 || definition.Visibility[0] != capability.VisibilityCLI || definition.Visibility[1] != capability.VisibilityTUI || definition.Visibility[2] != capability.VisibilityCommandPalette {
		t.Fatalf("Visibility = %#v, want cli, tui, command-palette", definition.Visibility)
	}
	if definition.InputSchema == nil || len(definition.InputSchema.Fields) != 1 {
		t.Fatalf("InputSchema = %#v, want one field", definition.InputSchema)
	}
	field := definition.InputSchema.Fields[0]
	if field.Name != "id" || field.Type != capability.InputTypeString || !field.Required {
		t.Fatalf("Input field = %#v, want required string id", field)
	}
}

func TestGetCapabilityExecutesUseCase(t *testing.T) {
	t.Parallel()

	getter := &fakeGetter{order: orders.Order{ID: "A123", Status: "created", CreatedAt: "2026-05-26T00:00:00Z"}}
	envelope, err := pipelineWithOrdersGetter(t, getter).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: orders.CapabilityGetID,
		Input:        capability.Input{"id": "A123"},
		Profile:      "staging",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !envelope.OK {
		t.Fatalf("OK = false, want true: %#v", envelope.Error)
	}
	if getter.got.ID != "A123" {
		t.Fatalf("getter input ID = %q, want A123", getter.got.ID)
	}
	result, ok := (*envelope.Data).(orders.GetResult)
	if !ok {
		t.Fatalf("Data = %T, want orders.GetResult", *envelope.Data)
	}
	if result.Order.ID != "A123" || result.Order.Status != "created" {
		t.Fatalf("Order = %#v, want fake order", result.Order)
	}
}

func TestGetCapabilityPropagatesStructuredDomainError(t *testing.T) {
	t.Parallel()

	envelope, err := pipelineWithOrdersGetter(t, &fakeGetter{err: capability.StructuredError{Code: orders.ErrorOrderNotFound, Message: "Order not found."}}).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: orders.CapabilityGetID,
		Input:        capability.Input{"id": "missing"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if envelope.OK {
		t.Fatalf("OK = true, want false")
	}
	if envelope.Error == nil || envelope.Error.Code != orders.ErrorOrderNotFound {
		t.Fatalf("Error = %#v, want %s", envelope.Error, orders.ErrorOrderNotFound)
	}
}

func TestUnavailableGetterReturnsStructuredError(t *testing.T) {
	t.Parallel()

	_, err := orders.UnavailableGetter{}.Get(context.Background(), orders.GetInput{ID: "A123"})
	var structured capability.StructuredError
	if !errors.As(err, &structured) {
		t.Fatalf("Get() error = %T, want StructuredError", err)
	}
	if structured.Code != orders.ErrorOrdersNotConfigured {
		t.Fatalf("StructuredError.Code = %q, want %q", structured.Code, orders.ErrorOrdersNotConfigured)
	}
}

func pipelineWithOrdersGetter(t *testing.T, getter orders.Getter) execution.Pipeline {
	t.Helper()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(orders.NewGetCapability(getter)); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}
	return execution.NewPipeline(
		builder.Finalize(),
		execution.WithRequestIDGenerator(func() (string, error) { return "req_test", nil }),
	)
}

type fakeGetter struct {
	order orders.Order
	err   error
	got   orders.GetInput
}

func (g *fakeGetter) Get(_ context.Context, input orders.GetInput) (orders.Order, error) {
	g.got = input
	if g.err != nil {
		return orders.Order{}, g.err
	}
	return g.order, nil
}
