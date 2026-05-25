package execution_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yargotev/exito-tools/internal/execution"
)

func TestNewRequestID(t *testing.T) {
	t.Parallel()

	requestID, err := execution.NewRequestID()
	if err != nil {
		t.Fatalf("NewRequestID() error = %v", err)
	}

	if !strings.HasPrefix(requestID, "req_") {
		t.Fatalf("requestID = %q, want req_ prefix", requestID)
	}
	if len(requestID) != len("req_")+32 {
		t.Fatalf("requestID length = %d, want %d", len(requestID), len("req_")+32)
	}
}

func TestMetadataEnvelopeMeta(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(1500 * time.Millisecond)

	metadata := execution.NewMetadata("req_test", "corr-123", startedAt, finishedAt)
	meta := metadata.EnvelopeMeta("staging", "orders.get")

	if meta.RequestID != "req_test" {
		t.Fatalf("RequestID = %q, want req_test", meta.RequestID)
	}
	if meta.CorrelationID != "corr-123" {
		t.Fatalf("CorrelationID = %q, want corr-123", meta.CorrelationID)
	}
	if meta.Profile != "staging" {
		t.Fatalf("Profile = %q, want staging", meta.Profile)
	}
	if meta.CapabilityID != "orders.get" {
		t.Fatalf("CapabilityID = %q, want orders.get", meta.CapabilityID)
	}
	if meta.DurationMS != 1500 {
		t.Fatalf("DurationMS = %d, want 1500", meta.DurationMS)
	}
}

func TestMetadataDoesNotReturnNegativeDuration(t *testing.T) {
	t.Parallel()

	finishedAt := time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC)
	startedAt := finishedAt.Add(time.Second)

	metadata := execution.NewMetadata("req_test", "", startedAt, finishedAt)
	meta := metadata.EnvelopeMeta("staging", "")

	if meta.DurationMS != 0 {
		t.Fatalf("DurationMS = %d, want 0", meta.DurationMS)
	}
}
