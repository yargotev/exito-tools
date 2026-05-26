package tui

import (
	"context"
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
)

// IO contains terminal streams for the TUI Surface.
type IO struct {
	Input  io.Reader
	Output io.Writer
}

// Model is the initial Bubble Tea model for the task-first TUI Surface.
type Model struct {
	profile      string
	capabilities []capability.Definition
}

// NewModel builds the initial TUI model from the bootstrapped application.
func NewModel(application *app.Application) Model {
	if application == nil {
		return Model{}
	}

	return Model{
		profile:      application.Config.Profile,
		capabilities: primaryActions(application.Registry.All()),
	}
}

// Init starts without side effects; capability execution is introduced by later task slices.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles foundational navigation keys.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the initial task-first TUI shell.
func (m Model) View() string {
	var builder strings.Builder
	builder.WriteString("Exito Tools TUI\n")
	fmt.Fprintf(&builder, "Profile: %s\n", profileLabel(m.profile))
	fmt.Fprintf(&builder, "Primary actions: %d\n", len(m.capabilities))

	for _, definition := range m.capabilities {
		fmt.Fprintf(&builder, "- %s (%s)\n", definition.Title, definition.ID)
	}

	builder.WriteString("\nPress q to quit.\n")
	return builder.String()
}

// Run starts the Bubble Tea TUI program.
func Run(ctx context.Context, application *app.Application, ioStreams IO) error {
	options := []tea.ProgramOption{tea.WithContext(ctx)}
	if ioStreams.Input != nil {
		options = append(options, tea.WithInput(ioStreams.Input))
	}
	if ioStreams.Output != nil {
		options = append(options, tea.WithOutput(ioStreams.Output))
	}

	_, err := tea.NewProgram(NewModel(application), options...).Run()
	return err
}

func primaryActions(definitions []capability.Definition) []capability.Definition {
	actions := make([]capability.Definition, 0, len(definitions))
	for _, definition := range definitions {
		if hasVisibility(definition, capability.VisibilityTUI) && hasAudience(definition, capability.AudiencePeople) {
			actions = append(actions, definition)
		}
	}
	return actions
}

func hasVisibility(definition capability.Definition, visibility capability.Visibility) bool {
	for _, candidate := range definition.Visibility {
		if candidate == visibility {
			return true
		}
	}
	return false
}

func hasAudience(definition capability.Definition, audience capability.Audience) bool {
	for _, candidate := range definition.Audiences {
		if candidate == audience {
			return true
		}
	}
	return false
}

func profileLabel(profile string) string {
	if profile == "" {
		return "staging"
	}
	return profile
}
