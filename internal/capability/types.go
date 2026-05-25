package capability

// Definition is the shared metadata skeleton for a capability.
type Definition struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// StructuredError is the shared machine-readable error skeleton.
type StructuredError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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
