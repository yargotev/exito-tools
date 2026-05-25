package capability

import "context"

// Definition is the shared metadata skeleton for a capability.
type Definition struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Input is a neutral, schema-shaped object supplied to a Capability execution.
type Input map[string]any

// ExecutionContext contains request-scoped values supplied to Capability handlers.
type ExecutionContext struct {
	RequestID     string
	CorrelationID string
	Profile       string
	CapabilityID  string
}

// ExecutionRequest is the neutral request passed to a Capability handler.
type ExecutionRequest struct {
	Input   Input
	Context ExecutionContext
}

// ExecutionResult is the format-neutral successful result returned by a Capability handler.
type ExecutionResult struct {
	Data any
}

// Handler runs a Capability use case without depending on an Interaction Surface.
type Handler func(context.Context, ExecutionRequest) (ExecutionResult, error)

// Executable is a registered Capability plus the handler that executes its use case.
type Executable struct {
	Definition Definition
	Handler    Handler
}

// StructuredError is the shared machine-readable error skeleton.
type StructuredError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error returns the stable message for Go error interoperability.
func (e StructuredError) Error() string {
	return e.Message
}

// EnvelopeMeta is the standard metadata container shared by CLI JSON envelopes.
type EnvelopeMeta struct {
	RequestID     string `json:"requestId"`
	CorrelationID string `json:"correlationId,omitempty"`
	Profile       string `json:"profile,omitempty"`
	CapabilityID  string `json:"capabilityId,omitempty"`
	DurationMS    int64  `json:"durationMs"`
}

// Envelope is the shared JSON-envelope-shaped result skeleton.
type Envelope[T any] struct {
	OK    bool             `json:"ok"`
	Data  *T               `json:"data,omitempty"`
	Error *StructuredError `json:"error,omitempty"`
	Meta  EnvelopeMeta     `json:"meta"`
}
