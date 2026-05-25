package app

import "github.com/yargotev/exito-tools/internal/registry"

// Application is the explicit application wiring seam shared by surfaces.
type Application struct {
	Registry registry.Registry
}

// New builds the minimal application scaffold and finalizes the registry.
func New() (*Application, error) {
	builder := registry.NewBuilder()

	return &Application{
		Registry: builder.Finalize(),
	}, nil
}
