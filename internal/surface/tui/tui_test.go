package tui_test

import (
	"context"
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

func TestCommandPaletteSelectionExecutesActionThroughPipeline(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "foundation.ping",
			Domain:     "foundation",
			Title:      "Ping foundation",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			if request.Context.Profile != "prod" {
				t.Fatalf("profile = %q, want prod", request.Context.Profile)
			}
			return capability.ExecutionResult{Data: map[string]any{"pong": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{
		Config:   config.Effective{Profile: "prod"},
		Registry: builder.Finalize(),
	})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("Update(enter) command = nil, want execution command")
	}

	loadingView := model.(tui.Model).View()
	if !strings.Contains(loadingView, "Running foundation.ping...") {
		t.Fatalf("loading view missing running state\n%s", loadingView)
	}

	model, _ = model.Update(cmd())
	view := model.(tui.Model).View()
	for _, fragment := range []string{"Task", "Success: foundation.ping"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("success view missing %q\n%s", fragment, view)
		}
	}
}

func TestCommandPaletteExecutionRendersStructuredFailure(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "orders.get",
			Domain:     "orders",
			Title:      "Get order",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "id", Type: capability.InputTypeNumber, Required: true},
			}},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			t.Fatal("handler should not run when required input is missing")
			return capability.ExecutionResult{}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("Update(enter) command = nil, want execution command")
	}

	model, _ = model.Update(cmd())
	view := model.(tui.Model).View()
	for _, fragment := range []string{"Failure: orders.get", "INVALID_INPUT", "Required input field \"id\" is missing."} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("failure view missing %q\n%s", fragment, view)
		}
	}
}

func TestCommandPaletteSelectionMovesWithArrowKeys(t *testing.T) {
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
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})

	view := model.(tui.Model).View()
	if !strings.Contains(view, "> Geocode address (geo.geocode-address)") {
		t.Fatalf("view should mark second action selected\n%s", view)
	}
	if !strings.Contains(view, "  Get order (orders.get)") {
		t.Fatalf("view should leave first action unselected\n%s", view)
	}
}

func TestCommandPaletteActionWithRequiredStringInputOpensForm(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "orders.get",
			Domain:     "orders",
			Title:      "Get order",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "id", Type: capability.InputTypeString, Required: true},
			}},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			t.Fatal("handler should not run before the form is submitted")
			return capability.ExecutionResult{}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatalf("Update(enter) command = %#v, want nil while form opens", cmd)
	}

	view := model.(tui.Model).View()
	for _, fragment := range []string{"Input Form", "Action: orders.get", "> id:"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("form view missing %q\n%s", fragment, view)
		}
	}
	if strings.Contains(view, "Command Palette") {
		t.Fatalf("palette should close while form is active\n%s", view)
	}
}

func TestInputFormSubmissionExecutesWithCollectedInput(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "geo.geocode-address",
			Domain:     "geo",
			Title:      "Geocode address",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "city", Type: capability.InputTypeString, Required: true},
				{Name: "address", Type: capability.InputTypeString, Required: true},
			}},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			if request.Input["city"] != "Bogota" {
				t.Fatalf("city input = %#v, want Bogota", request.Input["city"])
			}
			if request.Input["address"] != "CL 57" {
				t.Fatalf("address input = %#v, want CL 57", request.Input["address"])
			}
			return capability.ExecutionResult{Data: map[string]any{"success": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Bogota")})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	view := model.(tui.Model).View()
	if !strings.Contains(view, "> address:") {
		t.Fatalf("form should advance to address field\n%s", view)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("CL 57")})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("final field submit command = nil, want execution command")
	}

	loadingView := model.(tui.Model).View()
	if !strings.Contains(loadingView, "Running geo.geocode-address...") {
		t.Fatalf("loading view missing running state\n%s", loadingView)
	}
	if strings.Contains(loadingView, "Input Form") {
		t.Fatalf("form should close after submit\n%s", loadingView)
	}

	model, _ = model.Update(cmd())
	successView := model.(tui.Model).View()
	if !strings.Contains(successView, "Success: geo.geocode-address") {
		t.Fatalf("success view missing completed state\n%s", successView)
	}
}

func TestInputFormDoesNotSubmitEmptyRequiredStringField(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.Register(capability.Definition{
		ID:         "orders.get",
		Domain:     "orders",
		Title:      "Get order",
		Audiences:  []capability.Audience{capability.AudiencePeople},
		Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "id", Type: capability.InputTypeString, Required: true},
		}},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatalf("empty field submit command = %#v, want nil", cmd)
	}

	view := model.(tui.Model).View()
	if !strings.Contains(view, "Input Form") || !strings.Contains(view, "> id:") {
		t.Fatalf("empty submit should keep the form active\n%s", view)
	}
}

func TestResultFilterRefinesLoadedTaskDataWithoutReexecution(t *testing.T) {
	t.Parallel()

	executions := 0
	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "orders.get",
			Domain:     "orders",
			Title:      "Get order",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			executions++
			return capability.ExecutionResult{Data: map[string]any{
				"id":     "A123",
				"status": "ready",
				"city":   "Bogota",
			}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("Update(enter) command = nil, want execution command")
	}
	model, _ = model.Update(cmd())

	view := model.(tui.Model).View()
	for _, fragment := range []string{"Result", "- city: Bogota", "- id: A123", "- status: ready", "Press f to filter results."} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("success result view missing %q\n%s", fragment, view)
		}
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("ready")})

	filteredView := model.(tui.Model).View()
	for _, fragment := range []string{"Result Filter: ready", "- status: ready"} {
		if !strings.Contains(filteredView, fragment) {
			t.Fatalf("filtered result view missing %q\n%s", fragment, filteredView)
		}
	}
	for _, hidden := range []string{"- city: Bogota", "- id: A123", "Command Palette"} {
		if strings.Contains(filteredView, hidden) {
			t.Fatalf("filtered result view should not contain %q\n%s", hidden, filteredView)
		}
	}
	if executions != 1 {
		t.Fatalf("executions = %d, want 1", executions)
	}
}

func TestResultFilterEscClosesAndRestoresRows(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "geo.geocode-address",
			Domain:     "geo",
			Title:      "Geocode address",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityCommandPalette},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			return capability.ExecutionResult{Data: map[string]any{
				"address": "CL 57",
				"city":    "Bogota",
			}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = model.Update(cmd())
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("cl")})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	view := model.(tui.Model).View()
	for _, fragment := range []string{"- address: CL 57", "- city: Bogota", "Press f to filter results."} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("view after closing result filter missing %q\n%s", fragment, view)
		}
	}
	if strings.Contains(view, "Result Filter:") {
		t.Fatalf("result filter should close on esc\n%s", view)
	}
}
