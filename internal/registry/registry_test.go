package registry_test

import (
	"errors"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/registry"
)

func TestBuilderLifecycle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "register before finalize persists capability in snapshot",
			run: func(t *testing.T) {
				t.Helper()

				builder := registry.NewBuilder()
				definition := capability.Definition{
					ID:          "foundation.example",
					Title:       "Foundation Example",
					Description: "Registered during application boot.",
				}

				if err := builder.Register(definition); err != nil {
					t.Fatalf("Register() error = %v", err)
				}

				finalized := builder.Finalize()
				got := finalized.All()

				if len(got) != 1 {
					t.Fatalf("All() length = %d, want 1", len(got))
				}

				if got[0] != definition {
					t.Fatalf("All()[0] = %#v, want %#v", got[0], definition)
				}

				got[0].Title = "mutated outside registry"
				again := finalized.All()
				if again[0] != definition {
					t.Fatalf("All() should return a defensive copy, got %#v want %#v", again[0], definition)
				}
			},
		},
		{
			name: "empty finalize returns empty registry",
			run: func(t *testing.T) {
				t.Helper()

				builder := registry.NewBuilder()
				finalized := builder.Finalize()

				got := finalized.All()
				if len(got) != 0 {
					t.Fatalf("All() length = %d, want 0", len(got))
				}
			},
		},
		{
			name: "register after finalize returns stable error",
			run: func(t *testing.T) {
				t.Helper()

				builder := registry.NewBuilder()
				_ = builder.Finalize()

				err := builder.Register(capability.Definition{ID: "foundation.late"})
				if !errors.Is(err, registry.ErrRegistryFinalized) {
					t.Fatalf("Register() error = %v, want %v", err, registry.ErrRegistryFinalized)
				}
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.run(t)
		})
	}
}
