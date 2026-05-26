package httpclient_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

func TestApplyRequestMetadataSetsHeaders(t *testing.T) {
	t.Parallel()

	ctx := httpclient.ContextWithRequestMetadata(context.Background(), httpclient.RequestMetadata{
		RequestID:     "req_test",
		CorrelationID: "corr-123",
	})
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}

	httpclient.ApplyRequestMetadata(ctx, request)

	if request.Header.Get(httpclient.HeaderRequestID) != "req_test" {
		t.Fatalf("%s = %q, want req_test", httpclient.HeaderRequestID, request.Header.Get(httpclient.HeaderRequestID))
	}
	if request.Header.Get(httpclient.HeaderCorrelationID) != "corr-123" {
		t.Fatalf("%s = %q, want corr-123", httpclient.HeaderCorrelationID, request.Header.Get(httpclient.HeaderCorrelationID))
	}
}

func TestApplyRequestMetadataOmitsMissingCorrelationID(t *testing.T) {
	t.Parallel()

	ctx := httpclient.ContextWithRequestMetadata(context.Background(), httpclient.RequestMetadata{RequestID: "req_test"})
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}

	httpclient.ApplyRequestMetadata(ctx, request)

	if request.Header.Get(httpclient.HeaderRequestID) != "req_test" {
		t.Fatalf("%s = %q, want req_test", httpclient.HeaderRequestID, request.Header.Get(httpclient.HeaderRequestID))
	}
	if request.Header.Get(httpclient.HeaderCorrelationID) != "" {
		t.Fatalf("%s = %q, want empty", httpclient.HeaderCorrelationID, request.Header.Get(httpclient.HeaderCorrelationID))
	}
}

func TestRequestMetadataFromContextReportsAbsence(t *testing.T) {
	t.Parallel()

	if metadata, ok := httpclient.RequestMetadataFromContext(context.Background()); ok {
		t.Fatalf("RequestMetadataFromContext() = %#v, true; want false", metadata)
	}
}
