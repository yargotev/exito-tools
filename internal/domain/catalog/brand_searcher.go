package catalog

import "context"

// BrandSearcher dispatches catalog searches to the configured brand searcher.
type BrandSearcher struct {
	exito   Searcher
	carulla Searcher
}

// NewBrandSearcher creates a brand-aware VTEX Catalog searcher.
func NewBrandSearcher(exito Searcher, carulla Searcher) BrandSearcher {
	return BrandSearcher{exito: exito, carulla: carulla}
}

// SearchProducts dispatches by normalized brand.
func (s BrandSearcher) SearchProducts(ctx context.Context, input SearchProductsInput) (SearchProductsResult, error) {
	if normalizedBrand(input.Brand) == "carulla" {
		if s.carulla == nil {
			return UnavailableSearcher{}.SearchProducts(ctx, input)
		}
		return s.carulla.SearchProducts(ctx, input)
	}
	if s.exito == nil {
		return UnavailableSearcher{}.SearchProducts(ctx, input)
	}
	return s.exito.SearchProducts(ctx, input)
}
