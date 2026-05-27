package orders

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

// UnavailableVTEXOMSGetter is the default VTEX OMS dependency until configuration exists.
type UnavailableVTEXOMSGetter struct{}

// GetVTEXOMS returns a stable not-configured structured error.
func (UnavailableVTEXOMSGetter) GetVTEXOMS(context.Context, GetVTEXOMSInput) (VTEXOMSOrder, error) {
	return VTEXOMSOrder{}, capability.StructuredError{
		Code:    ErrorOrdersNotConfigured,
		Message: "VTEX OMS client is not configured.",
	}
}
