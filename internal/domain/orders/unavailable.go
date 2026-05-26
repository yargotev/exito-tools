package orders

import (
	"context"

	"github.com/yargotev/exito-tools/internal/capability"
)

// UnavailableGetter is the default Orders dependency until real API configuration exists.
type UnavailableGetter struct{}

// Get returns a structured configuration error without contacting external services.
func (UnavailableGetter) Get(context.Context, GetInput) (Order, error) {
	return Order{}, capability.StructuredError{
		Code:    ErrorOrdersNotConfigured,
		Message: "Orders client is not configured.",
	}
}
