package capability

import "context"

// RiskLevel classifies the operational risk of running a Capability.
type RiskLevel string

const (
	// RiskReadOnly marks a Capability that does not intentionally mutate external state.
	RiskReadOnly RiskLevel = "read-only"
	// RiskSafeWrite marks a Capability that mutates state in a routine, reversible, or low-risk way.
	RiskSafeWrite RiskLevel = "safe-write"
	// RiskDestructive marks a Capability that can delete, overwrite, or otherwise cause high-impact changes.
	RiskDestructive RiskLevel = "destructive"
)

// Audience declares the intended consumer group for a Capability.
type Audience string

const (
	// AudienceAgents marks a Capability intended for automation and agent workflows.
	AudienceAgents Audience = "agents"
	// AudiencePeople marks a Capability intended for human-facing flows.
	AudiencePeople Audience = "people"
)

// Visibility declares where a Capability may be promoted by interaction surfaces.
type Visibility string

const (
	// VisibilityCLI makes a Capability visible to the machine-first CLI surface.
	VisibilityCLI Visibility = "cli"
	// VisibilityTUI makes a Capability visible to primary TUI navigation.
	VisibilityTUI Visibility = "tui"
	// VisibilityCommandPalette makes a Capability searchable in command palette style discovery.
	VisibilityCommandPalette Visibility = "command-palette"
)

// InputType declares a neutral input field type that surfaces can adapt into flags, forms, or validation.
type InputType string

const (
	// InputTypeString marks a textual input value.
	InputTypeString InputType = "string"
	// InputTypeNumber marks a numeric input value.
	InputTypeNumber InputType = "number"
	// InputTypeBoolean marks a boolean input value.
	InputTypeBoolean InputType = "boolean"
	// InputTypeObject marks a nested object input value.
	InputTypeObject InputType = "object"
	// InputTypeArray marks a list input value.
	InputTypeArray InputType = "array"
)

// InputField describes one schema-shaped Capability input field.
type InputField struct {
	Name        string    `json:"name"`
	Type        InputType `json:"type"`
	Required    bool      `json:"required,omitempty"`
	Description string    `json:"description,omitempty"`
}

// InputSchema describes the complete neutral input object accepted by a Capability.
type InputSchema struct {
	Fields []InputField `json:"fields,omitempty"`
}

// Definition is the shared metadata skeleton for a capability.
type Definition struct {
	ID                   string       `json:"id"`
	Domain               string       `json:"domain,omitempty"`
	Version              string       `json:"version,omitempty"`
	Title                string       `json:"title"`
	Description          string       `json:"description"`
	Risk                 RiskLevel    `json:"risk,omitempty"`
	RequiresConfirmation bool         `json:"requiresConfirmation,omitempty"`
	Audiences            []Audience   `json:"audiences,omitempty"`
	Visibility           []Visibility `json:"visibility,omitempty"`
	InputSchema          *InputSchema `json:"inputSchema,omitempty"`
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
	Data       any
	Warnings   []StructuredWarning
	Pagination *PaginationMeta
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

// StructuredWarning is a non-fatal machine-readable warning included in Envelope metadata.
type StructuredWarning struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// PaginationMeta is optional cursor-based metadata for paginated Capability results.
type PaginationMeta struct {
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

// Error returns the stable message for Go error interoperability.
func (e StructuredError) Error() string {
	return e.Message
}

// EnvelopeMeta is the standard metadata container shared by CLI JSON envelopes.
type EnvelopeMeta struct {
	RequestID     string              `json:"requestId"`
	CorrelationID string              `json:"correlationId,omitempty"`
	Profile       string              `json:"profile,omitempty"`
	CapabilityID  string              `json:"capabilityId,omitempty"`
	DurationMS    int64               `json:"durationMs"`
	Warnings      []StructuredWarning `json:"warnings,omitempty"`
	Pagination    *PaginationMeta     `json:"pagination,omitempty"`
}

// Envelope is the shared JSON-envelope-shaped result skeleton.
type Envelope[T any] struct {
	OK    bool             `json:"ok"`
	Data  *T               `json:"data,omitempty"`
	Error *StructuredError `json:"error,omitempty"`
	Meta  EnvelopeMeta     `json:"meta"`
}
