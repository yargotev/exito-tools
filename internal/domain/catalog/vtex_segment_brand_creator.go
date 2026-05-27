package catalog

import "context"

// VTEXSegmentBrandCreator dispatches VTEX segment requests to the configured brand creator.
type VTEXSegmentBrandCreator struct {
	exito   VTEXSegmentCreator
	carulla VTEXSegmentCreator
}

// NewVTEXSegmentBrandCreator creates a brand-aware VTEX segment creator.
func NewVTEXSegmentBrandCreator(exito VTEXSegmentCreator, carulla VTEXSegmentCreator) VTEXSegmentBrandCreator {
	return VTEXSegmentBrandCreator{exito: exito, carulla: carulla}
}

// CreateVTEXSegment dispatches by normalized brand.
func (c VTEXSegmentBrandCreator) CreateVTEXSegment(ctx context.Context, input CreateVTEXSegmentInput) (CreateVTEXSegmentResult, error) {
	if normalizedBrand(input.Brand) == "carulla" {
		if c.carulla == nil {
			return UnavailableVTEXSegmentCreator{}.CreateVTEXSegment(ctx, input)
		}
		return c.carulla.CreateVTEXSegment(ctx, input)
	}
	if c.exito == nil {
		return UnavailableVTEXSegmentCreator{}.CreateVTEXSegment(ctx, input)
	}
	return c.exito.CreateVTEXSegment(ctx, input)
}
