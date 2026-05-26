package tui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/registry"
	"github.com/yargotev/exito-tools/internal/surface/tui"
)

func TestModelViewShowsProfileAndPrimaryActions(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.Register(capability.Definition{
		ID:         "orders.get",
		Title:      "Get order",
		Audiences:  []capability.Audience{capability.AudienceAgents, capability.AudiencePeople},
		Visibility: []capability.Visibility{capability.VisibilityCLI, capability.VisibilityTUI},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := builder.Register(capability.Definition{
		ID:         "internal.agent-only",
		Title:      "Agent only",
		Audiences:  []capability.Audience{capability.AudienceAgents},
		Visibility: []capability.Visibility{capability.VisibilityCLI, capability.VisibilityTUI},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	model := tui.NewModel(&app.Application{
		Config:   config.Effective{Profile: "prod"},
		Registry: builder.Finalize(),
	})

	view := model.View()
	for _, fragment := range []string{"Exito Tools TUI", "Profile: prod", "Primary actions: 1", "Get order (orders.get)", "Press q to quit."} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("view missing %q\n%s", fragment, view)
		}
	}
	if strings.Contains(view, "Agent only") {
		t.Fatalf("agent-only capability should not be promoted as a primary action\n%s", view)
	}
}

func TestModelQuitKeysExitProgram(t *testing.T) {
	t.Parallel()

	model := tui.NewModel(&app.Application{Registry: registry.NewBuilder().Finalize()})
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatalf("Update(q) command = nil, want quit command")
	}
}
