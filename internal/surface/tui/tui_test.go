package tui_test

import (
	"context"
	"os"
	"path/filepath"
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
	for _, fragment := range []string{"Exito Tools TUI", "Profile: prod", "Primary actions: 1", "Get order (orders.get)", "Press p to change session profile. Press d to save default profile. Press q to quit."} {
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

func TestDefaultProfileFormSavesProfile(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	configPath := filepath.Join(workDir, "exito.yaml")
	var model tea.Model = tui.NewModel(&app.Application{
		Config: config.Effective{Profile: "staging"},
		ConfigOptions: config.Options{
			ConfigPath: configPath,
			WorkDir:    workDir,
			HomeDir:    t.TempDir(),
			Env:        map[string]string{},
		},
		Registry: registry.NewBuilder().Finalize(),
	})

	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if cmd != nil {
		t.Fatalf("Update(d) command = %#v, want nil", cmd)
	}
	view := model.(tui.Model).View()
	for _, fragment := range []string{"Default Profile", "Current session: staging", "> Save default as:"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("default profile form missing %q\n%s", fragment, view)
		}
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("prod")})
	model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("default profile submit command = nil, want save command")
	}
	model, _ = model.Update(cmd())

	view = model.(tui.Model).View()
	for _, fragment := range []string{"Profile: prod", "Default Profile saved: prod", configPath} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("saved default profile view missing %q\n%s", fragment, view)
		}
	}
	content, err := os.ReadFile(configPath) // #nosec G304 -- test reads the config path created in t.TempDir.
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", configPath, err)
	}
	if string(content) != "defaultProfile: prod\n" {
		t.Fatalf("config content = %q, want default profile line", string(content))
	}
}

func TestDefaultProfileFormCancelKeepsActiveProfileAndDoesNotPersist(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	configPath := filepath.Join(workDir, "exito.yaml")
	var model tea.Model = tui.NewModel(&app.Application{
		Config: config.Effective{Profile: "staging"},
		ConfigOptions: config.Options{
			ConfigPath: configPath,
			WorkDir:    workDir,
			HomeDir:    t.TempDir(),
			Env:        map[string]string{},
		},
		Registry: registry.NewBuilder().Finalize(),
	})

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("prod")})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Fatalf("default profile cancel command = %#v, want nil", cmd)
	}

	view := model.(tui.Model).View()
	if !strings.Contains(view, "Profile: staging") {
		t.Fatalf("cancel should keep original profile\n%s", view)
	}
	if strings.Contains(view, "Default Profile\nCurrent session") {
		t.Fatalf("default profile form should close on esc\n%s", view)
	}
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatalf("config file should not be persisted on cancel, stat error = %v", err)
	}
}

func TestDefaultProfileSaveFailureKeepsActiveProfile(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	configPath := filepath.Join(workDir, "as-directory")
	if err := os.Mkdir(configPath, 0o700); err != nil {
		t.Fatalf("Mkdir(%q) error = %v", configPath, err)
	}
	var model tea.Model = tui.NewModel(&app.Application{
		Config: config.Effective{Profile: "staging"},
		ConfigOptions: config.Options{
			ConfigPath: configPath,
			WorkDir:    workDir,
			HomeDir:    t.TempDir(),
			Env:        map[string]string{},
		},
		Registry: registry.NewBuilder().Finalize(),
	})

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("prod")})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("default profile submit command = nil, want save command")
	}
	model, _ = model.Update(cmd())

	view := model.(tui.Model).View()
	for _, fragment := range []string{"Profile: staging", "Default Profile save failed:"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("failed default profile view missing %q\n%s", fragment, view)
		}
	}
}

func TestSessionProfileFormChangesActiveProfile(t *testing.T) {
	t.Parallel()

	var model tea.Model = tui.NewModel(&app.Application{
		Config:   config.Effective{Profile: "staging"},
		Registry: registry.NewBuilder().Finalize(),
	})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	if cmd != nil {
		t.Fatalf("Update(p) command = %#v, want nil", cmd)
	}

	view := model.(tui.Model).View()
	for _, fragment := range []string{"Session Profile", "Current: staging", "> New profile:"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("profile form view missing %q\n%s", fragment, view)
		}
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("prod")})
	model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatalf("profile submit command = %#v, want nil", cmd)
	}

	view = model.(tui.Model).View()
	if !strings.Contains(view, "Profile: prod") {
		t.Fatalf("view should show changed session profile\n%s", view)
	}
	if strings.Contains(view, "Session Profile") {
		t.Fatalf("profile form should close after submit\n%s", view)
	}
}

func TestSessionProfileFormCancelKeepsActiveProfile(t *testing.T) {
	t.Parallel()

	var model tea.Model = tui.NewModel(&app.Application{
		Config:   config.Effective{Profile: "staging"},
		Registry: registry.NewBuilder().Finalize(),
	})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("prod")})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Fatalf("profile cancel command = %#v, want nil", cmd)
	}

	view := model.(tui.Model).View()
	if !strings.Contains(view, "Profile: staging") {
		t.Fatalf("cancel should keep original profile\n%s", view)
	}
	if strings.Contains(view, "Session Profile") {
		t.Fatalf("profile form should close on esc\n%s", view)
	}
}

func TestSessionProfileChangeAppliesToSubsequentActionExecution(t *testing.T) {
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
		Config:   config.Effective{Profile: "staging"},
		Registry: builder.Finalize(),
	})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("prod")})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("Update(enter) command = nil, want execution command")
	}

	model, _ = model.Update(cmd())
	view := model.(tui.Model).View()
	if !strings.Contains(view, "Success: foundation.ping") {
		t.Fatalf("success view missing completed state\n%s", view)
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

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("CL")})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("57")})
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

func TestInputFormAcceptsSpacesInAddressValues(t *testing.T) {
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
			if request.Input["city"] != "Envigado" {
				t.Fatalf("city input = %#v, want Envigado", request.Input["city"])
			}
			if request.Input["address"] != "Carrera 3A # 10 A - 22" {
				t.Fatalf("address input = %#v, want spaces preserved", request.Input["address"])
			}
			return capability.ExecutionResult{Data: map[string]any{"success": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Envigado")})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	for _, chunk := range []string{"Carrera", "3A", "#", "10", "A", "-", "22"} {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(chunk)})
		if chunk != "22" {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}})
		}
	}

	view := model.(tui.Model).View()
	if !strings.Contains(view, "> address: Carrera 3A # 10 A - 22") {
		t.Fatalf("form view should preserve visible address spaces\n%s", view)
	}

	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("address submit command = nil, want execution command")
	}
	model, _ = model.Update(cmd())
	if !strings.Contains(model.(tui.Model).View(), "Success: geo.geocode-address") {
		t.Fatalf("spaced address should execute successfully\n%s", model.(tui.Model).View())
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

func TestConfirmationRequiredActionShowsPromptBeforeExecution(t *testing.T) {
	t.Parallel()

	executions := 0
	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:                   "orders.cancel",
			Domain:               "orders",
			Title:                "Cancel order",
			Description:          "Cancels the selected order.",
			Risk:                 capability.RiskSafeWrite,
			RequiresConfirmation: true,
			Audiences:            []capability.Audience{capability.AudiencePeople},
			Visibility:           []capability.Visibility{capability.VisibilityCommandPalette},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			executions++
			return capability.ExecutionResult{Data: map[string]any{"cancelled": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatalf("Update(enter) command = %#v, want nil while confirmation prompt opens", cmd)
	}

	view := model.(tui.Model).View()
	for _, fragment := range []string{"Confirm Action", "Action: Cancel order", "Capability: orders.cancel", "Risk: safe-write", "Impact: Cancels the selected order.", "Press y or enter to confirm. Press n or esc to cancel."} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("confirmation view missing %q\n%s", fragment, view)
		}
	}
	if strings.Contains(view, "Command Palette") {
		t.Fatalf("palette should close while confirmation prompt is active\n%s", view)
	}
	if executions != 0 {
		t.Fatalf("executions = %d, want 0 before confirmation", executions)
	}
}

func TestConfirmationPromptConfirmExecutesWithExplicitConfirmation(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:                   "orders.cancel",
			Domain:               "orders",
			Title:                "Cancel order",
			Risk:                 capability.RiskSafeWrite,
			RequiresConfirmation: true,
			Audiences:            []capability.Audience{capability.AudiencePeople},
			Visibility:           []capability.Visibility{capability.VisibilityCommandPalette},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			return capability.ExecutionResult{Data: map[string]any{"cancelled": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if cmd == nil {
		t.Fatalf("Update(y) command = nil, want execution command")
	}

	loadingView := model.(tui.Model).View()
	if !strings.Contains(loadingView, "Running orders.cancel...") {
		t.Fatalf("loading view missing running state\n%s", loadingView)
	}
	if strings.Contains(loadingView, "Confirm Action") {
		t.Fatalf("confirmation prompt should close after confirm\n%s", loadingView)
	}

	model, _ = model.Update(cmd())
	view := model.(tui.Model).View()
	if !strings.Contains(view, "Success: orders.cancel") {
		t.Fatalf("confirmed action should succeed through Pipeline\n%s", view)
	}
	if strings.Contains(view, "CONFIRMATION_REQUIRED") {
		t.Fatalf("confirmed TUI action should not fail confirmation policy\n%s", view)
	}
}

func TestConfirmationPromptCancelDoesNotExecute(t *testing.T) {
	t.Parallel()

	executions := 0
	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:                   "orders.cancel",
			Domain:               "orders",
			Title:                "Cancel order",
			Risk:                 capability.RiskSafeWrite,
			RequiresConfirmation: true,
			Audiences:            []capability.Audience{capability.AudiencePeople},
			Visibility:           []capability.Visibility{capability.VisibilityCommandPalette},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			executions++
			return capability.ExecutionResult{Data: nil}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if cmd != nil {
		t.Fatalf("Update(n) command = %#v, want nil when cancelling confirmation", cmd)
	}

	view := model.(tui.Model).View()
	if strings.Contains(view, "Confirm Action") {
		t.Fatalf("confirmation prompt should close after cancel\n%s", view)
	}
	if strings.Contains(view, "Running orders.cancel") || strings.Contains(view, "Success: orders.cancel") || strings.Contains(view, "Failure: orders.cancel") {
		t.Fatalf("cancelled confirmation should not start task\n%s", view)
	}
	if executions != 0 {
		t.Fatalf("executions = %d, want 0 after cancelling confirmation", executions)
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

func TestLoadingTaskEscCancelsContextAndRendersCancelledState(t *testing.T) {
	t.Parallel()

	contextCancelled := false
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
			select {
			case <-ctx.Done():
				contextCancelled = true
				return capability.ExecutionResult{}, ctx.Err()
			default:
				return capability.ExecutionResult{Data: map[string]any{"id": "A123"}}, nil
			}
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

	model, cancelCmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cancelCmd != nil {
		t.Fatalf("Update(esc) command = %#v, want nil while cancelling task", cancelCmd)
	}
	cancelledView := model.(tui.Model).View()
	if !strings.Contains(cancelledView, "Cancelled: orders.get") {
		t.Fatalf("cancelled view missing cancelled state\n%s", cancelledView)
	}
	if strings.Contains(cancelledView, "Command Palette") {
		t.Fatalf("cancelled view should close the command palette\n%s", cancelledView)
	}

	model, _ = model.Update(cmd())
	if !contextCancelled {
		t.Fatalf("handler context was not cancelled")
	}
	lateView := model.(tui.Model).View()
	if !strings.Contains(lateView, "Cancelled: orders.get") {
		t.Fatalf("late completion should not replace cancelled state\n%s", lateView)
	}
	if strings.Contains(lateView, "Failure: orders.get") || strings.Contains(lateView, "Success: orders.get") {
		t.Fatalf("late completion should not render success or failure\n%s", lateView)
	}
}

func TestPrimaryNavigationSupportsArrowsVimKeysAndEnterExecution(t *testing.T) {
	t.Parallel()

	executed := ""
	builder := registry.NewBuilder()
	for _, definition := range []capability.Definition{
		{
			ID:         "orders.get",
			Domain:     "orders",
			Title:      "Get order",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityTUI},
		},
		{
			ID:         "geo.geocode-address",
			Domain:     "geo",
			Title:      "Geocode address",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityTUI},
		},
	} {
		def := definition
		if err := builder.RegisterExecutable(capability.Executable{
			Definition: def,
			Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
				executed = request.Context.CapabilityID
				return capability.ExecutionResult{Data: map[string]any{"ok": true}}, nil
			},
		}); err != nil {
			t.Fatalf("RegisterExecutable() error = %v", err)
		}
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	initialView := model.(tui.Model).View()
	if !strings.Contains(initialView, "Navigate with ↑/↓ or j/k") || !strings.Contains(initialView, "› Get order (orders.get)") {
		t.Fatalf("initial view should expose navigable primary actions\n%s", initialView)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	view := model.(tui.Model).View()
	if !strings.Contains(view, "› Geocode address (geo.geocode-address)") {
		t.Fatalf("j should move primary selection down\n%s", view)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	view = model.(tui.Model).View()
	if !strings.Contains(view, "› Get order (orders.get)") {
		t.Fatalf("arrow up should move primary selection up\n%s", view)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("primary enter command = nil, want execution command")
	}
	model, _ = model.Update(cmd())
	if executed != "geo.geocode-address" {
		t.Fatalf("executed = %q, want geo.geocode-address", executed)
	}
	if !strings.Contains(model.(tui.Model).View(), "Success: geo.geocode-address") {
		t.Fatalf("primary action should render success\n%s", model.(tui.Model).View())
	}
}

func TestCommandPaletteVimNavigationMovesSelection(t *testing.T) {
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
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	view := model.(tui.Model).View()
	if !strings.Contains(view, "> Geocode address (geo.geocode-address)") {
		t.Fatalf("palette j should move selection down\n%s", view)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	view = model.(tui.Model).View()
	if !strings.Contains(view, "> Get order (orders.get)") {
		t.Fatalf("palette k should move selection up\n%s", view)
	}
}

func TestTUIE2EPrimaryActionCollectsOptionalInputAndRunsPipeline(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "orders.get",
			Domain:     "orders",
			Title:      "Get order",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityTUI, capability.VisibilityCommandPalette},
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "id", Type: capability.InputTypeString, Required: true},
				{Name: "orderType", Type: capability.InputTypeString, Required: false},
			}},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			if request.Input["id"] != "1611511090420" {
				t.Fatalf("id input = %#v, want 1611511090420", request.Input["id"])
			}
			if _, ok := request.Input["orderType"]; ok {
				t.Fatalf("empty optional orderType should be omitted, input = %#v", request.Input)
			}
			return capability.ExecutionResult{Data: map[string]any{"order": request.Input["id"]}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{
		Config:   config.Effective{Profile: "staging"},
		Registry: builder.Finalize(),
	})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatalf("opening primary form command = %#v, want nil", cmd)
	}
	if !strings.Contains(model.(tui.Model).View(), "Input Form") || !strings.Contains(model.(tui.Model).View(), "> id:") {
		t.Fatalf("primary action should open input form\n%s", model.(tui.Model).View())
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1611511090420")})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !strings.Contains(model.(tui.Model).View(), "> orderType:") {
		t.Fatalf("form should advance to optional orderType field\n%s", model.(tui.Model).View())
	}

	model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("submitting optional field command = nil, want execution command")
	}
	model, _ = model.Update(cmd())
	view := model.(tui.Model).View()
	for _, fragment := range []string{"Success: orders.get", "- order: 1611511090420"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("e2e view missing %q\n%s", fragment, view)
		}
	}
}

func TestInputFormVimNormalModeMovesCursorAndDeletesUnderCursor(t *testing.T) {
	t.Parallel()

	gotQuery := ""
	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "geo.search",
			Domain:     "geo",
			Title:      "Search",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityTUI},
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "query", Type: capability.InputTypeString, Required: true},
			}},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			gotQuery, _ = request.Input["query"].(string)
			return capability.ExecutionResult{Data: map[string]any{"query": gotQuery}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("bat")})

	view := model.(tui.Model).View()
	if !strings.Contains(view, "> query: bat▌") {
		t.Fatalf("insert mode should render cursor after typed text\n%s", view)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	view = model.(tui.Model).View()
	if !strings.Contains(view, "Keyboard: Vim normal") {
		t.Fatalf("esc in form should enter Vim normal mode\n%s", view)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	view = model.(tui.Model).View()
	if !strings.Contains(view, "> query: b▌t") {
		t.Fatalf("h/h/x should remove the character under the cursor\n%s", view)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("submit command = nil, want execution command")
	}
	_, _ = model.Update(cmd())
	if gotQuery != "bot" {
		t.Fatalf("query = %q, want bot", gotQuery)
	}
}

func TestInputFormPlainModeLetsHJKLEditText(t *testing.T) {
	t.Parallel()

	gotQuery := ""
	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID:         "geo.search",
			Domain:     "geo",
			Title:      "Search",
			Audiences:  []capability.Audience{capability.AudiencePeople},
			Visibility: []capability.Visibility{capability.VisibilityTUI},
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "query", Type: capability.InputTypeString, Required: true},
			}},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			gotQuery, _ = request.Input["query"].(string)
			return capability.ExecutionResult{Data: map[string]any{"query": gotQuery}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	var model tea.Model = tui.NewModel(&app.Application{Registry: builder.Finalize()})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if !strings.Contains(model.(tui.Model).View(), "Keyboard: Plain") {
		t.Fatalf("v should toggle to plain keyboard mode\n%s", model.(tui.Model).View())
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hjkl")})
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("submit command = nil, want execution command")
	}
	_, _ = model.Update(cmd())
	if gotQuery != "hjkl" {
		t.Fatalf("query = %q, want hjkl", gotQuery)
	}
}
