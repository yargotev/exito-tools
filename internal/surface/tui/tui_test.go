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

func TestCommandPaletteShowsPeopleFacingPaletteActions(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.Register(capability.Definition{
		ID:         "orders.get",
		Domain:     "orders",
		Title:      "Get order",
		Audiences:  []capability.Audience{capability.AudiencePeople},
		Visibility: []capability.Visibility{capability.VisibilityTUI, capability.VisibilityCommandPalette},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := builder.Register(capability.Definition{
		ID:         "geo.geocode-address",
		Domain:     "geo",
		Title:      "Geocode address",
		Audiences:  []capability.Audience{capability.AudiencePeople},
		Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := builder.Register(capability.Definition{
		ID:         "internal.agent-only",
		Domain:     "internal",
		Title:      "Agent only",
		Audiences:  []capability.Audience{capability.AudienceAgents},
		Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	model := tui.NewModel(&app.Application{Registry: builder.Finalize()})
	initialView := model.View()
	if strings.Contains(initialView, "Geocode address") {
		t.Fatalf("palette-only action should not appear in primary navigation\n%s", initialView)
	}

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	if cmd != nil {
		t.Fatalf("Update(/) command = %#v, want nil", cmd)
	}
	view := updated.(tui.Model).View()
	for _, fragment := range []string{"Command Palette", "Get order (orders.get)", "Geocode address (geo.geocode-address)"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("palette view missing %q\n%s", fragment, view)
		}
	}
	if strings.Contains(view, "Agent only") {
		t.Fatalf("agent-only action should not appear in command palette\n%s", view)
	}
}

func TestCommandPaletteFiltersActionsByQuery(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	for _, definition := range []capability.Definition{
		{
			ID:         "orders.get",
			Domain:     "orders",
			Title:      "Get order",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
		},
		{
			ID:         "geo.geocode-address",
			Domain:     "geo",
			Title:      "Geocode address",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
		},
	} {
		if err := builder.Register(definition); err != nil {
			t.Fatalf("Register() error = %v", err)
		}
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("geo")})

	view := model.(tui.Model).View()
	for _, fragment := range []string{"Command Palette", "Search: geo", "Geocode address (geo.geocode-address)"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("filtered palette view missing %q\n%s", fragment, view)
		}
	}
	if strings.Contains(view, "Get order") {
		t.Fatalf("filtered palette should hide non-matching action\n%s", view)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if strings.Contains(model.(tui.Model).View(), "Command Palette") {
		t.Fatalf("palette should close on esc\n%s", model.(tui.Model).View())
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
