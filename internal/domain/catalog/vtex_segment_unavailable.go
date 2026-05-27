package catalog

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

// UnavailableVTEXSegmentCreator returns a stable configuration error when VTEX Sessions is not configured.
type UnavailableVTEXSegmentCreator struct{}

// CreateVTEXSegment reports an unavailable VTEX Sessions provider.
func (UnavailableVTEXSegmentCreator) CreateVTEXSegment(context.Context, CreateVTEXSegmentInput) (CreateVTEXSegmentResult, error) {
	return CreateVTEXSegmentResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Sessions client is not configured."}
}
