package execution_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
	"github.com/yargotev/exito-tools/internal/registry"
)

func TestPipelineExecutesRegisteredCapability(t *testing.T) {
	t.Parallel()

	input := capability.Input{"id": "123"}
	var gotRequest capability.ExecutionRequest
	var gotContext context.Context

	reg := registryWithExecutable(t, capability.Executable{
		Definition: capability.Definition{ID: "orders.get"},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			gotContext = ctx
			gotRequest = request
			return capability.ExecutionResult{Data: map[string]any{"orderId": request.Input["id"]}}, nil
		},
	})

	ctx := context.WithValue(context.Background(), testContextKey{}, "context-value")
	pipeline := deterministicPipeline(reg)
	envelope, err := pipeline.Execute(ctx, execution.ExecuteRequest{
		CapabilityID:  "orders.get",
		Input:         input,
		Profile:       "staging",
		CorrelationID: "corr-123",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !envelope.OK {
		t.Fatalf("OK = false, want true")
	}
	if gotContext.Value(testContextKey{}) != "context-value" {
		t.Fatalf("handler did not receive original context")
	}
	gotHTTPMetadata, ok := httpclient.RequestMetadataFromContext(gotContext)
	if !ok {
		t.Fatalf("handler context missing HTTP request metadata")
	}
	if gotHTTPMetadata.RequestID != "req_test" || gotHTTPMetadata.CorrelationID != "corr-123" {
		t.Fatalf("HTTP request metadata = %#v, want request/correlation IDs", gotHTTPMetadata)
	}
	if !reflect.DeepEqual(gotRequest.Input, input) {
		t.Fatalf("handler input = %#v, want %#v", gotRequest.Input, input)
	}
	if gotRequest.Context.RequestID != "req_test" {
		t.Fatalf("handler request ID = %q, want req_test", gotRequest.Context.RequestID)
	}
	if gotRequest.Context.CorrelationID != "corr-123" {
		t.Fatalf("handler correlation ID = %q, want corr-123", gotRequest.Context.CorrelationID)
	}
	if gotRequest.Context.Profile != "staging" {
		t.Fatalf("handler profile = %q, want staging", gotRequest.Context.Profile)
	}
	if gotRequest.Context.CapabilityID != "orders.get" {
		t.Fatalf("handler capability ID = %q, want orders.get", gotRequest.Context.CapabilityID)
	}
	if envelope.Data == nil {
		t.Fatalf("Data is nil")
	}
	data, ok := (*envelope.Data).(map[string]any)
	if !ok || data["orderId"] != "123" {
		t.Fatalf("Data = %#v, want orderId 123", envelope.Data)
	}
	assertMeta(t, envelope.Meta, "corr-123", "orders.get")
}

func TestPipelineValidatesCapabilityInputSchema(t *testing.T) {
	t.Parallel()

	handlerCalled := false
	reg := registryWithExecutable(t, capability.Executable{
		Definition: capability.Definition{
			ID: "orders.get",
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "id", Type: capability.InputTypeString, Required: true},
			}},
		},
		Handler: func(context.Context, capability.ExecutionRequest) (capability.ExecutionResult, error) {
			handlerCalled = true
			return capability.ExecutionResult{Data: map[string]any{"unreachable": true}}, nil
		},
	})

	envelope, err := deterministicPipeline(reg).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: "orders.get",
		Input:        capability.Input{},
		Profile:      "staging",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if handlerCalled {
		t.Fatalf("handler was called for invalid input")
	}
	if envelope.OK {
		t.Fatalf("OK = true, want false")
	}
	if envelope.Error == nil || envelope.Error.Code != execution.ErrorInvalidInput {
		t.Fatalf("Error = %#v, want %s", envelope.Error, execution.ErrorInvalidInput)
	}
	assertMeta(t, envelope.Meta, "", "orders.get")
}

func TestPipelinePreservesStructuredError(t *testing.T) {
	t.Parallel()

	reg := registryWithExecutable(t, capability.Executable{
		Definition: capability.Definition{ID: "orders.get"},
		Handler: func(context.Context, capability.ExecutionRequest) (capability.ExecutionResult, error) {
			return capability.ExecutionResult{}, capability.StructuredError{Code: "ORDER_NOT_FOUND", Message: "Order not found."}
		},
	})

	envelope, err := deterministicPipeline(reg).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: "orders.get",
		Profile:      "staging",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if envelope.OK {
		t.Fatalf("OK = true, want false")
	}
	if envelope.Error == nil || envelope.Error.Code != "ORDER_NOT_FOUND" || envelope.Error.Message != "Order not found." {
		t.Fatalf("Error = %#v, want structured order error", envelope.Error)
	}
	if envelope.Data != nil {
		t.Fatalf("Data = %#v, want nil on failure", envelope.Data)
	}
	assertMeta(t, envelope.Meta, "", "orders.get")
}

func TestPipelineTranslatesUnknownError(t *testing.T) {
	t.Parallel()

	reg := registryWithExecutable(t, capability.Executable{
		Definition: capability.Definition{ID: "orders.get"},
		Handler: func(context.Context, capability.ExecutionRequest) (capability.ExecutionResult, error) {
			return capability.ExecutionResult{}, errors.New("database unavailable")
		},
	})

	envelope, err := deterministicPipeline(reg).Execute(context.Background(), execution.ExecuteRequest{CapabilityID: "orders.get"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if envelope.OK {
		t.Fatalf("OK = true, want false")
	}
	if envelope.Error == nil || envelope.Error.Code != execution.ErrorCapabilityExecutionFailed {
		t.Fatalf("Error = %#v, want %s", envelope.Error, execution.ErrorCapabilityExecutionFailed)
	}
}

func TestPipelineReturnsCapabilityNotFound(t *testing.T) {
	t.Parallel()

	envelope, err := deterministicPipeline(registry.NewBuilder().Finalize()).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID:  "missing.example",
		Profile:       "staging",
		CorrelationID: "corr-123",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if envelope.OK {
		t.Fatalf("OK = true, want false")
	}
	if envelope.Error == nil || envelope.Error.Code != execution.ErrorCapabilityNotFound {
		t.Fatalf("Error = %#v, want %s", envelope.Error, execution.ErrorCapabilityNotFound)
	}
	assertMeta(t, envelope.Meta, "corr-123", "missing.example")
}

type testContextKey struct{}

func deterministicPipeline(reg registry.Registry) execution.Pipeline {
	times := []time.Time{
		time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 25, 10, 0, 1, 500*int(time.Millisecond), time.UTC),
	}
	index := 0
	return execution.NewPipeline(
		reg,
		execution.WithRequestIDGenerator(func() (string, error) { return "req_test", nil }),
		execution.WithClock(func() time.Time {
			if index >= len(times) {
				return times[len(times)-1]
			}
			current := times[index]
			index++
			return current
		}),
	)
}

func registryWithExecutable(t *testing.T, entry capability.Executable) registry.Registry {
	t.Helper()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(entry); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}
	return builder.Finalize()
}

func assertMeta(t *testing.T, got capability.EnvelopeMeta, correlationID string, capabilityID string) {
	t.Helper()

	if got.RequestID != "req_test" {
		t.Fatalf("RequestID = %q, want req_test", got.RequestID)
	}
	if got.CorrelationID != correlationID {
		t.Fatalf("CorrelationID = %q, want %q", got.CorrelationID, correlationID)
	}
	if got.Profile != "staging" {
		t.Fatalf("Profile = %q, want staging", got.Profile)
	}
	if got.CapabilityID != capabilityID {
		t.Fatalf("CapabilityID = %q, want %q", got.CapabilityID, capabilityID)
	}
	if got.DurationMS != 1500 {
		t.Fatalf("DurationMS = %d, want 1500", got.DurationMS)
	}
}
