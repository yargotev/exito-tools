package geo

import (
	"context"
)

// VTEXBrandRegionResolver routes VTEX region diagnostics to the selected brand provider.
type VTEXBrandRegionResolver struct {
	exito   VTEXRegionResolver
	carulla VTEXRegionResolver
}

// NewVTEXBrandRegionResolver creates a brand-aware VTEX region resolver.
func NewVTEXBrandRegionResolver(exito VTEXRegionResolver, carulla VTEXRegionResolver) VTEXBrandRegionResolver {
	return VTEXBrandRegionResolver{exito: exito, carulla: carulla}
}

// ResolveVTEXRegion dispatches by normalized brand.
func (r VTEXBrandRegionResolver) ResolveVTEXRegion(ctx context.Context, input ResolveVTEXRegionInput) (ResolveVTEXRegionResult, error) {
	if normalizeBrand(input.Brand) == "carulla" {
		return r.carulla.ResolveVTEXRegion(ctx, input)
	}
	return r.exito.ResolveVTEXRegion(ctx, input)
}
