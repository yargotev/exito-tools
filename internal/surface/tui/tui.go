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
	profile        string
	primaryActions []capability.Definition
	paletteActions []capability.Definition
	paletteOpen    bool
	paletteQuery   string
}

// NewModel builds the initial TUI model from the bootstrapped application.
func NewModel(application *app.Application) Model {
	if application == nil {
		return Model{}
	}

	return Model{
		profile:        application.Config.Profile,
		primaryActions: primaryActions(application.Registry.All()),
		paletteActions: paletteActions(application.Registry.All()),
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
		if m.paletteOpen {
			return m.updatePalette(msg)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "/":
			m.paletteOpen = true
			m.paletteQuery = ""
		}
	}

	return m, nil
}

// View renders the initial task-first TUI shell.
func (m Model) View() string {
	var builder strings.Builder
	builder.WriteString("Exito Tools TUI\n")
	fmt.Fprintf(&builder, "Profile: %s\n", profileLabel(m.profile))
	fmt.Fprintf(&builder, "Primary actions: %d\n", len(m.primaryActions))

	for _, definition := range m.primaryActions {
		fmt.Fprintf(&builder, "- %s (%s)\n", definition.Title, definition.ID)
	}

	if m.paletteOpen {
		builder.WriteString("\nCommand Palette\n")
		fmt.Fprintf(&builder, "Search: %s\n", m.paletteQuery)
		matches := m.filteredPaletteActions()
		if len(matches) == 0 {
			builder.WriteString("No actions found.\n")
		}
		for _, definition := range matches {
			fmt.Fprintf(&builder, "- %s (%s)\n", definition.Title, definition.ID)
		}
		builder.WriteString("Press esc to close.\n")
	}

	builder.WriteString("\nPress q to quit.\n")
	return builder.String()
}

func (m Model) updatePalette(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.paletteOpen = false
		m.paletteQuery = ""
	case tea.KeyBackspace:
		if len(m.paletteQuery) > 0 {
			runes := []rune(m.paletteQuery)
			m.paletteQuery = string(runes[:len(runes)-1])
		}
	case tea.KeyRunes:
		m.paletteQuery += string(msg.Runes)
	case tea.KeyCtrlC:
		return m, tea.Quit
	}

	return m, nil
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

func paletteActions(definitions []capability.Definition) []capability.Definition {
	actions := make([]capability.Definition, 0, len(definitions))
	for _, definition := range definitions {
		if hasVisibility(definition, capability.VisibilityCommandPalette) && hasAudience(definition, capability.AudiencePeople) {
			actions = append(actions, definition)
		}
	}
	return actions
}

func (m Model) filteredPaletteActions() []capability.Definition {
	query := strings.ToLower(strings.TrimSpace(m.paletteQuery))
	if query == "" {
		return m.paletteActions
	}

	matches := make([]capability.Definition, 0, len(m.paletteActions))
	for _, definition := range m.paletteActions {
		if strings.Contains(strings.ToLower(definition.Title), query) ||
			strings.Contains(strings.ToLower(definition.ID), query) ||
			strings.Contains(strings.ToLower(definition.Domain), query) {
			matches = append(matches, definition)
		}
	}
	return matches
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
