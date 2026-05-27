package orders

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

// VTEXOMSBrandGetter routes VTEX OMS lookups to the configured brand client.
type VTEXOMSBrandGetter struct {
	exito   VTEXOMSGetterPort
	carulla VTEXOMSGetterPort
}

// NewVTEXOMSBrandGetter creates a brand-aware VTEX OMS getter.
func NewVTEXOMSBrandGetter(exito VTEXOMSGetterPort, carulla VTEXOMSGetterPort) VTEXOMSBrandGetter {
	return VTEXOMSBrandGetter{exito: exito, carulla: carulla}
}

// GetVTEXOMS routes by input brand, defaulting to exito.
func (g VTEXOMSBrandGetter) GetVTEXOMS(ctx context.Context, input GetVTEXOMSInput) (VTEXOMSOrder, error) {
	switch normalizedVTEXOMSBrand(input.Brand) {
	case "carulla":
		if g.carulla == nil {
			return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrdersNotConfigured, Message: "VTEX OMS Carulla client is not configured."}
		}
		return g.carulla.GetVTEXOMS(ctx, input)
	default:
		if g.exito == nil {
			return VTEXOMSOrder{}, capability.StructuredError{Code: ErrorOrdersNotConfigured, Message: "VTEX OMS Exito client is not configured."}
		}
		return g.exito.GetVTEXOMS(ctx, input)
	}
}
