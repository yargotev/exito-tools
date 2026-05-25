package registry

import (
	"errors"

	"github.com/yargotev/exito-tools/internal/capability"
)

var (
	ErrRegistryFinalized   = errors.New("capability registry is finalized")
	ErrDuplicateCapability = errors.New("capability ID is already registered")
)

// Builder collects capabilities during application boot.
type Builder struct {
	entries   []capability.Executable
	seenIDs   map[string]struct{}
	finalized bool
}

// Registry exposes an immutable capability inventory to surfaces.
type Registry interface {
	All() []capability.Definition
	Find(id string) (capability.Executable, bool)
}

type finalizedRegistry struct {
	entries []capability.Executable
	byID    map[string]capability.Executable
}

// NewBuilder creates a boot-time capability registry builder.
func NewBuilder() *Builder {
	return &Builder{seenIDs: map[string]struct{}{}}
}

// Register adds a capability definition during boot.
func (b *Builder) Register(def capability.Definition) error {
	return b.RegisterExecutable(capability.Executable{Definition: def})
}

// RegisterExecutable adds an executable capability entry during boot.
func (b *Builder) RegisterExecutable(entry capability.Executable) error {
	if b.finalized {
		return ErrRegistryFinalized
	}

	if _, exists := b.seenIDs[entry.Definition.ID]; exists {
		return ErrDuplicateCapability
	}

	b.seenIDs[entry.Definition.ID] = struct{}{}
	b.entries = append(b.entries, entry)
	return nil
}

// Finalize freezes the registry and returns an immutable snapshot.
func (b *Builder) Finalize() Registry {
	b.finalized = true
	entries := cloneEntries(b.entries)
	byID := make(map[string]capability.Executable, len(entries))
	for _, entry := range entries {
		byID[entry.Definition.ID] = entry
	}

	return finalizedRegistry{entries: entries, byID: byID}
}

// All returns a defensive copy of the finalized capability inventory.
func (r finalizedRegistry) All() []capability.Definition {
	definitions := make([]capability.Definition, len(r.entries))
	for i, entry := range r.entries {
		definitions[i] = entry.Definition
	}
	return definitions
}

// Find returns a registered executable capability by stable Capability ID.
func (r finalizedRegistry) Find(id string) (capability.Executable, bool) {
	entry, ok := r.byID[id]
	return entry, ok
}

func cloneEntries(entries []capability.Executable) []capability.Executable {
	if len(entries) == 0 {
		return []capability.Executable{}
	}

	cloned := make([]capability.Executable, len(entries))
	copy(cloned, entries)
	return cloned
}
