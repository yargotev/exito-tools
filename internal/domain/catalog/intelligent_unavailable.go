package catalog

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

// UnavailableIntelligentSearcher returns a stable configuration error when VTEX Intelligent Search is not configured.
type UnavailableIntelligentSearcher struct{}

// IntelligentSearchProducts reports an unavailable VTEX Intelligent Search provider.
func (UnavailableIntelligentSearcher) IntelligentSearchProducts(context.Context, IntelligentSearchProductsInput) (IntelligentSearchProductsResult, error) {
	return IntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Intelligent Search client is not configured."}
}
