package capability

// Definition is the shared metadata skeleton for a capability.
type Definition struct {
	ID          string
	Title       string
	Description string
}

// StructuredError is the shared machine-readable error skeleton.
type StructuredError struct {
	Code    string
	Message string
}

// Envelope is the shared JSON-envelope-shaped result skeleton.
type Envelope[T any] struct {
	OK    bool
	Data  *T
	Error *StructuredError
}
