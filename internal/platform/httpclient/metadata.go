package httpclient

import (
	"context"
	"net/http"
)

const (
	// HeaderRequestID carries the Exito Tools request identifier to outbound HTTP providers.
	HeaderRequestID = "X-Request-Id"
	// HeaderCorrelationID carries the optional caller-supplied correlation identifier to outbound HTTP providers.
	HeaderCorrelationID = "X-Correlation-Id"
)

type metadataContextKey struct{}

// RequestMetadata contains provider-agnostic HTTP request metadata.
type RequestMetadata struct {
	RequestID     string
	CorrelationID string
}

// ContextWithRequestMetadata attaches outbound HTTP metadata to a context.
func ContextWithRequestMetadata(ctx context.Context, metadata RequestMetadata) context.Context {
	return context.WithValue(ctx, metadataContextKey{}, metadata)
}

// RequestMetadataFromContext returns outbound HTTP metadata attached to a context.
func RequestMetadataFromContext(ctx context.Context) (RequestMetadata, bool) {
	metadata, ok := ctx.Value(metadataContextKey{}).(RequestMetadata)
	return metadata, ok
}

// ApplyRequestMetadata copies context metadata into provider request headers.
func ApplyRequestMetadata(ctx context.Context, request *http.Request) {
	if request == nil {
		return
	}

	metadata, ok := RequestMetadataFromContext(ctx)
	if !ok {
		return
	}
	if metadata.RequestID != "" {
		request.Header.Set(HeaderRequestID, metadata.RequestID)
	}
	if metadata.CorrelationID != "" {
		request.Header.Set(HeaderCorrelationID, metadata.CorrelationID)
	}
}
