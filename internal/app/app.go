package app

import (
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/domain/geo"
	"github.com/yargotev/exito-tools/internal/domain/orders"
	"github.com/yargotev/exito-tools/internal/registry"
)

// Options contains application boot inputs shared by all surfaces.
type Options struct {
	Config config.Options
}

// Application is the explicit application wiring seam shared by surfaces.
type Application struct {
	Config   config.Effective
	Registry registry.Registry
}

// New builds the minimal application scaffold, resolves configuration, and finalizes the registry.
func New(options Options) (*Application, error) {
	effectiveConfig, err := config.Resolve(options.Config)
	if err != nil {
		return nil, err
	}

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(orders.NewGetCapability(orders.UnavailableGetter{})); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(geo.NewGeocodeAddressCapability(geo.UnavailableGeocoder{})); err != nil {
		return nil, err
	}

	return &Application{
		Config:   effectiveConfig,
		Registry: builder.Finalize(),
	}, nil
}
