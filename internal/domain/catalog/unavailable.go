package catalog

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

// UnavailableSearcher returns a stable configuration error when VTEX Catalog is not configured.
type UnavailableSearcher struct{}

// SearchProducts reports an unavailable VTEX Catalog provider.
func (UnavailableSearcher) SearchProducts(context.Context, SearchProductsInput) (SearchProductsResult, error) {
	return SearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Catalog client is not configured."}
}
