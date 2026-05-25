package execution

import (
	"context"
	"errors"
	"time"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/registry"
)

const (
	ErrorCapabilityNotFound        = "CAPABILITY_NOT_FOUND"
	ErrorCapabilityExecutionFailed = "CAPABILITY_EXECUTION_FAILED"
)

// Pipeline runs registered Capabilities through a shared surface-independent path.
type Pipeline struct {
	registry     registry.Registry
	newRequestID func() (string, error)
	now          func() time.Time
}

// PipelineOption customizes Pipeline behavior, primarily for deterministic tests.
type PipelineOption func(*Pipeline)

// WithRequestIDGenerator sets the request ID generator used by the Pipeline.
func WithRequestIDGenerator(generator func() (string, error)) PipelineOption {
	return func(p *Pipeline) {
		p.newRequestID = generator
	}
}

// WithClock sets the clock used by the Pipeline.
func WithClock(clock func() time.Time) PipelineOption {
	return func(p *Pipeline) {
		p.now = clock
	}
}

// NewPipeline builds a Capability execution Pipeline over an immutable Registry.
func NewPipeline(registry registry.Registry, options ...PipelineOption) Pipeline {
	pipeline := Pipeline{
		registry:     registry,
		newRequestID: NewRequestID,
		now:          time.Now,
	}

	for _, option := range options {
		option(&pipeline)
	}

	return pipeline
}

// ExecuteRequest contains neutral inputs for one Capability execution.
type ExecuteRequest struct {
	CapabilityID  string
	Input         capability.Input
	Profile       string
	CorrelationID string
}

// Execute runs a registered Capability and returns a standard JSON Envelope-shaped result.
func (p Pipeline) Execute(ctx context.Context, request ExecuteRequest) (capability.Envelope[any], error) {
	startedAt := p.now()
	requestID, err := p.newRequestID()
	if err != nil {
		return capability.Envelope[any]{}, err
	}

	entry, ok := p.registry.Find(request.CapabilityID)
	if !ok || entry.Handler == nil {
		return p.failureEnvelope(request, requestID, startedAt, capability.StructuredError{
			Code:    ErrorCapabilityNotFound,
			Message: "Capability not found.",
		}), nil
	}

	executionRequest := capability.ExecutionRequest{
		Input: request.Input,
		Context: capability.ExecutionContext{
			RequestID:     requestID,
			CorrelationID: request.CorrelationID,
			Profile:       request.Profile,
			CapabilityID:  request.CapabilityID,
		},
	}

	result, err := entry.Handler(ctx, executionRequest)
	if err != nil {
		return p.failureEnvelope(request, requestID, startedAt, structuredError(err)), nil
	}

	data := result.Data
	metadata := NewMetadata(requestID, request.CorrelationID, startedAt, p.now())
	return capability.Envelope[any]{
		OK:   true,
		Data: &data,
		Meta: metadata.EnvelopeMeta(request.Profile, request.CapabilityID),
	}, nil
}

func (p Pipeline) failureEnvelope(request ExecuteRequest, requestID string, startedAt time.Time, structured capability.StructuredError) capability.Envelope[any] {
	metadata := NewMetadata(requestID, request.CorrelationID, startedAt, p.now())
	return capability.Envelope[any]{
		OK:    false,
		Error: &structured,
		Meta:  metadata.EnvelopeMeta(request.Profile, request.CapabilityID),
	}
}

func structuredError(err error) capability.StructuredError {
	var structured capability.StructuredError
	if errors.As(err, &structured) {
		return structured
	}

	return capability.StructuredError{
		Code:    ErrorCapabilityExecutionFailed,
		Message: "Capability execution failed.",
	}
}
