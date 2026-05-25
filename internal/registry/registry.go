package registry

import (
	"errors"

	"github.com/yargotev/exito-tools/internal/capability"
)

var ErrRegistryFinalized = errors.New("capability registry is finalized")

// Builder collects capabilities during application boot.
type Builder struct {
	definitions []capability.Definition
	finalized   bool
}

// Registry exposes an immutable capability inventory to surfaces.
type Registry interface {
	All() []capability.Definition
}

type finalizedRegistry struct {
	definitions []capability.Definition
}

// NewBuilder creates a boot-time capability registry builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// Register adds a capability definition during boot.
func (b *Builder) Register(def capability.Definition) error {
	if b.finalized {
		return ErrRegistryFinalized
	}

	b.definitions = append(b.definitions, def)
	return nil
}

// Finalize freezes the registry and returns an immutable snapshot.
func (b *Builder) Finalize() Registry {
	b.finalized = true

	return finalizedRegistry{definitions: cloneDefinitions(b.definitions)}
}

// All returns a defensive copy of the finalized capability inventory.
func (r finalizedRegistry) All() []capability.Definition {
	return cloneDefinitions(r.definitions)
}

func cloneDefinitions(definitions []capability.Definition) []capability.Definition {
	if len(definitions) == 0 {
		return []capability.Definition{}
	}

	cloned := make([]capability.Definition, len(definitions))
	copy(cloned, definitions)
	return cloned
}
