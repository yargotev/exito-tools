package geo

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

// UnavailableGeocoder is the default Geo dependency until real API configuration exists.
type UnavailableGeocoder struct{}

// GeocodeAddress returns a structured configuration error without contacting external services.
func (UnavailableGeocoder) GeocodeAddress(context.Context, GeocodeAddressInput) (GeocodeAddressResult, error) {
	return GeocodeAddressResult{}, capability.StructuredError{
		Code:    ErrorGeoNotConfigured,
		Message: "Geo client is not configured.",
	}
}

// UnavailableVTEXRegionResolver is the default VTEX region dependency until provider configuration exists.
type UnavailableVTEXRegionResolver struct{}

// ResolveVTEXRegion returns a structured configuration error without contacting external services.
func (UnavailableVTEXRegionResolver) ResolveVTEXRegion(context.Context, ResolveVTEXRegionInput) (ResolveVTEXRegionResult, error) {
	return ResolveVTEXRegionResult{}, capability.StructuredError{
		Code:    ErrorGeoNotConfigured,
		Message: "VTEX region resolver is not configured.",
	}
}
