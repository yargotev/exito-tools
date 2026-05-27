package catalog

import "context"

// IntelligentBrandSearcher dispatches VTEX Intelligent Search requests to the configured brand searcher.
type IntelligentBrandSearcher struct {
	exito   IntelligentSearchProductsSearcher
	carulla IntelligentSearchProductsSearcher
}

// NewIntelligentBrandSearcher creates a brand-aware Intelligent Search searcher.
func NewIntelligentBrandSearcher(exito IntelligentSearchProductsSearcher, carulla IntelligentSearchProductsSearcher) IntelligentBrandSearcher {
	return IntelligentBrandSearcher{exito: exito, carulla: carulla}
}

// IntelligentSearchProducts dispatches by normalized brand.
func (s IntelligentBrandSearcher) IntelligentSearchProducts(ctx context.Context, input IntelligentSearchProductsInput) (IntelligentSearchProductsResult, error) {
	if normalizedBrand(input.Brand) == "carulla" {
		if s.carulla == nil {
			return UnavailableIntelligentSearcher{}.IntelligentSearchProducts(ctx, input)
		}
		return s.carulla.IntelligentSearchProducts(ctx, input)
	}
	if s.exito == nil {
		return UnavailableIntelligentSearcher{}.IntelligentSearchProducts(ctx, input)
	}
	return s.exito.IntelligentSearchProducts(ctx, input)
}
