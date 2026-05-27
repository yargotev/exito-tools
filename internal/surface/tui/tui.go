package tui

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/registry"
)

var (
	mochaBase     = lipgloss.Color("#1e1e2e")
	mochaSubtext  = lipgloss.Color("#a6adc8")
	mochaLavender = lipgloss.Color("#b4befe")
	mochaMauve    = lipgloss.Color("#cba6f7")
	mochaGreen    = lipgloss.Color("#a6e3a1")
	mochaYellow   = lipgloss.Color("#f9e2af")
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(mochaMauve).Background(mochaBase).Padding(0, 1)
	sectionStyle  = lipgloss.NewStyle().Bold(true).Foreground(mochaLavender)
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(mochaBase).Background(mochaLavender).Padding(0, 1)
	pillStyle     = lipgloss.NewStyle().Bold(true).Foreground(mochaGreen)
	mutedStyle    = lipgloss.NewStyle().Foreground(mochaSubtext)
	warnStyle     = lipgloss.NewStyle().Foreground(mochaYellow)
)

// IO contains terminal streams for the TUI Surface.
type IO struct {
	Input  io.Reader
	Output io.Writer
}

// Model is the initial Bubble Tea model for the task-first TUI Surface.
type Model struct {
	ctx            context.Context
	registry       registry.Registry
	configOptions  config.Options
	profile        string
	primaryActions []capability.Definition
	primaryIndex   int
	paletteActions []capability.Definition
	paletteOpen    bool
	paletteQuery   string
	paletteIndex   int
	form           formState
	profileForm    profileFormState
	defaultProfile defaultProfileFormState
	confirmation   confirmationState
	task           taskState
	taskCancel     context.CancelFunc
	resultFilter   resultFilterState
	vimEnabled     bool
}

type taskStatus string

const (
	taskIdle      taskStatus = ""
	taskLoading   taskStatus = "loading"
	taskSuccess   taskStatus = "success"
	taskFailure   taskStatus = "failure"
	taskCancelled taskStatus = "cancelled"
)

type taskState struct {
	Status       taskStatus
	CapabilityID string
	Data         any
	Error        *capability.StructuredError
}

type resultFilterState struct {
	Active bool
	Query  string
	Cursor int
}

type actionExecutedMsg struct {
	envelope capability.Envelope[any]
	err      error
}

type formState struct {
	Active       bool
	CapabilityID string
	Fields       []capability.InputField
	Values       []string
	Cursors      []int
	Index        int
	InsertMode   bool
}

type profileFormState struct {
	Active     bool
	Value      string
	Cursor     int
	InsertMode bool
}

type defaultProfileFormState struct {
	Active     bool
	Value      string
	Cursor     int
	InsertMode bool
	Status     string
}

type defaultProfileSavedMsg struct {
	result config.DefaultProfileWriteResult
	err    error
}

type confirmationState struct {
	Active     bool
	Definition capability.Definition
	Input      capability.Input
}

// NewModel builds the initial TUI model from the bootstrapped application.
func NewModel(application *app.Application) Model {
	if application == nil {
		return Model{}
	}

	return Model{
		ctx:            context.Background(),
		registry:       application.Registry,
		configOptions:  application.ConfigOptions,
		profile:        application.Config.Profile,
		vimEnabled:     true,
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
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		if m.task.Status == taskLoading && msg.Type == tea.KeyEsc {
			return m.cancelTask(), nil
		}

		if m.resultFilter.Active {
			return m.updateResultFilter(msg)
		}

		if m.form.Active {
			return m.updateForm(msg)
		}

		if m.profileForm.Active {
			return m.updateProfileForm(msg)
		}

		if m.defaultProfile.Active {
			return m.updateDefaultProfileForm(msg)
		}

		if m.confirmation.Active {
			return m.updateConfirmation(msg)
		}

		if m.paletteOpen {
			return m.updatePalette(msg)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "v":
			m.vimEnabled = !m.vimEnabled
		case "up":
			m.movePrimary(-1)
		case "down":
			m.movePrimary(1)
		case "k":
			if m.vimEnabled {
				m.movePrimary(-1)
			}
		case "j":
			if m.vimEnabled {
				m.movePrimary(1)
			}
		case "enter":
			return m.selectPrimaryAction()
		case "/":
			m.paletteOpen = true
			m.paletteQuery = ""
			m.paletteIndex = 0
		case "f":
			if len(resultRows(m.task.Data)) > 0 && m.task.Status == taskSuccess {
				m.resultFilter = resultFilterState{Active: true}
			}
		case "p":
			m.profileForm = profileFormState{Active: true, InsertMode: true}
		case "d":
			m.defaultProfile = defaultProfileFormState{Active: true, InsertMode: true}
		}
	case actionExecutedMsg:
		if m.task.Status == taskCancelled && m.task.CapabilityID == msg.envelope.Meta.CapabilityID {
			return m, nil
		}
		m.taskCancel = nil
		if msg.err != nil {
			m.task = taskState{
				Status: taskFailure,
				Error: &capability.StructuredError{
					Code:    execution.ErrorCapabilityExecutionFailed,
					Message: "Capability execution failed.",
				},
			}
			return m, nil
		}
		m.task.CapabilityID = msg.envelope.Meta.CapabilityID
		if msg.envelope.OK {
			m.task.Status = taskSuccess
			if msg.envelope.Data != nil {
				m.task.Data = *msg.envelope.Data
			} else {
				m.task.Data = nil
			}
			m.task.Error = nil
			m.resultFilter = resultFilterState{}
			return m, nil
		}
		m.task.Status = taskFailure
		m.task.Error = msg.envelope.Error
		m.task.Data = nil
		m.resultFilter = resultFilterState{}
	case defaultProfileSavedMsg:
		if msg.err != nil {
			m.defaultProfile = defaultProfileFormState{Status: fmt.Sprintf("Default Profile save failed: %s", msg.err.Error())}
			return m, nil
		}
		m.profile = msg.result.Profile
		m.defaultProfile = defaultProfileFormState{
			Status: fmt.Sprintf("Default Profile saved: %s (%s)", msg.result.Profile, msg.result.ConfigPath),
		}
	}

	return m, nil
}

// View renders the task-first TUI shell.
func (m Model) View() string {
	var builder strings.Builder
	builder.WriteString(titleStyle.Render("Exito Tools TUI"))
	builder.WriteString("\n")
	fmt.Fprintf(&builder, "%s %s\n", mutedStyle.Render("Profile:"), pillStyle.Render(profileLabel(m.profile)))
	fmt.Fprintf(&builder, "%s %d\n", mutedStyle.Render("Primary actions:"), len(m.primaryActions))
	fmt.Fprintf(&builder, "%s %s\n", mutedStyle.Render("Keyboard:"), pillStyle.Render(m.keyboardModeLabel()))
	builder.WriteString(mutedStyle.Render("Navigate with ↑/↓ or j/k, press enter to run, / for palette."))
	builder.WriteString("\n\n")

	if len(m.primaryActions) == 0 {
		builder.WriteString(warnStyle.Render("No primary actions available."))
		builder.WriteString("\n")
	}
	for index, definition := range m.primaryActions {
		line := fmt.Sprintf("  %s (%s)", definition.Title, definition.ID)
		if index == m.clampedPrimaryIndex() {
			line = selectedStyle.Render("› " + definition.Title + " (" + definition.ID + ")")
		}
		fmt.Fprintf(&builder, "%s\n", line)
	}

	if m.task.Status != taskIdle {
		builder.WriteString("\n")
		builder.WriteString(sectionStyle.Render("Task"))
		builder.WriteString("\n")
		switch m.task.Status {
		case taskLoading:
			fmt.Fprintf(&builder, "Running %s...\n", m.task.CapabilityID)
		case taskSuccess:
			fmt.Fprintf(&builder, "Success: %s\n", m.task.CapabilityID)
			m.renderResultRows(&builder)
		case taskFailure:
			fmt.Fprintf(&builder, "Failure: %s\n", m.task.CapabilityID)
			if m.task.Error != nil {
				fmt.Fprintf(&builder, "%s: %s\n", m.task.Error.Code, m.task.Error.Message)
			}
		case taskCancelled:
			fmt.Fprintf(&builder, "Cancelled: %s\n", m.task.CapabilityID)
		}
	}

	if m.form.Active {
		builder.WriteString("\n")
		builder.WriteString(sectionStyle.Render("Input Form"))
		builder.WriteString("\n")
		fmt.Fprintf(&builder, "Action: %s\n", m.form.CapabilityID)
		for index, field := range m.form.Fields {
			prefix := " "
			if index == m.form.Index {
				prefix = ">"
			}
			fmt.Fprintf(&builder, "%s %s: %s\n", prefix, field.Name, m.form.renderedValue(index))
		}
		fmt.Fprintf(&builder, "%s\n", m.form.helpText(m.vimEnabled))
	}

	if m.profileForm.Active {
		builder.WriteString("\n")
		builder.WriteString(sectionStyle.Render("Session Profile"))
		builder.WriteString("\n")
		fmt.Fprintf(&builder, "Current: %s\n", profileLabel(m.profile))
		fmt.Fprintf(&builder, "> New profile: %s\n", renderInputValue(m.profileForm.Value, m.profileForm.Cursor, true))
		builder.WriteString(textInputHelp(m.vimEnabled, m.profileForm.InsertMode, "apply", "cancel"))
		builder.WriteString("\n")
	}

	if m.defaultProfile.Active {
		builder.WriteString("\n")
		builder.WriteString(sectionStyle.Render("Default Profile"))
		builder.WriteString("\n")
		fmt.Fprintf(&builder, "Current session: %s\n", profileLabel(m.profile))
		fmt.Fprintf(&builder, "> Save default as: %s\n", renderInputValue(m.defaultProfile.Value, m.defaultProfile.Cursor, true))
		builder.WriteString(textInputHelp(m.vimEnabled, m.defaultProfile.InsertMode, "save", "cancel"))
		builder.WriteString("\n")
	} else if m.defaultProfile.Status != "" {
		builder.WriteString("\n")
		builder.WriteString(sectionStyle.Render("Default Profile"))
		builder.WriteString("\n")
		fmt.Fprintf(&builder, "%s\n", m.defaultProfile.Status)
	}

	if m.confirmation.Active {
		builder.WriteString("\n")
		builder.WriteString(sectionStyle.Render("Confirm Action"))
		builder.WriteString("\n")
		fmt.Fprintf(&builder, "Action: %s\n", m.confirmation.Definition.Title)
		fmt.Fprintf(&builder, "Capability: %s\n", m.confirmation.Definition.ID)
		if m.confirmation.Definition.Risk != "" {
			fmt.Fprintf(&builder, "Risk: %s\n", m.confirmation.Definition.Risk)
		}
		if m.confirmation.Definition.Description != "" {
			fmt.Fprintf(&builder, "Impact: %s\n", m.confirmation.Definition.Description)
		}
		builder.WriteString("Press y or enter to confirm. Press n or esc to cancel.\n")
	}

	if m.paletteOpen {
		builder.WriteString("\n")
		builder.WriteString(sectionStyle.Render("Command Palette"))
		builder.WriteString("\n")
		fmt.Fprintf(&builder, "Search: %s\n", renderInputValue(m.paletteQuery, len([]rune(m.paletteQuery)), true))
		matches := m.filteredPaletteActions()
		if len(matches) == 0 {
			builder.WriteString("No actions found.\n")
		}
		for index, definition := range matches {
			prefix := " "
			if index == m.paletteIndex {
				prefix = ">"
			}
			fmt.Fprintf(&builder, "%s %s (%s)\n", prefix, definition.Title, definition.ID)
		}
		builder.WriteString("Press esc to close.\n")
	}

	builder.WriteString("\n")
	builder.WriteString(mutedStyle.Render("Keys: j/k or arrows navigate • h/l move cursor in Vim normal inputs • i insert • x delete • v Vim/plain • enter run/select • / palette • p session profile • d default profile • f filter results • esc cancel/close • q quit"))
	builder.WriteString("\n")
	builder.WriteString("Press p to change session profile. Press d to save default profile. Press q to quit.\n")
	return builder.String()
}

func (m Model) updatePalette(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.paletteOpen = false
		m.paletteQuery = ""
		m.paletteIndex = 0
	case tea.KeyEnter:
		matches := m.filteredPaletteActions()
		if len(matches) == 0 {
			return m, nil
		}
		selected := matches[m.clampedPaletteIndex(len(matches))]
		m.paletteIndex = m.clampedPaletteIndex(len(matches))
		if fields := stringFields(selected); len(fields) > 0 {
			m.paletteOpen = false
			m.paletteQuery = ""
			m.form = newFormState(selected.ID, fields)
			return m, nil
		}
		return m.confirmOrStartExecution(selected, capability.Input{})
	case tea.KeyUp:
		m.movePalette(-1)
	case tea.KeyDown:
		m.movePalette(1)
	case tea.KeyBackspace:
		if len(m.paletteQuery) > 0 {
			runes := []rune(m.paletteQuery)
			m.paletteQuery = string(runes[:len(runes)-1])
		}
		m.paletteIndex = m.clampedPaletteIndex(len(m.filteredPaletteActions()))
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "j":
			if m.vimEnabled {
				m.movePalette(1)
			} else {
				m.paletteQuery += string(msg.Runes)
				m.paletteIndex = m.clampedPaletteIndex(len(m.filteredPaletteActions()))
			}
		case "k":
			if m.vimEnabled {
				m.movePalette(-1)
			} else {
				m.paletteQuery += string(msg.Runes)
				m.paletteIndex = m.clampedPaletteIndex(len(m.filteredPaletteActions()))
			}
		default:
			m.paletteQuery += string(msg.Runes)
			m.paletteIndex = m.clampedPaletteIndex(len(m.filteredPaletteActions()))
		}
	case tea.KeyCtrlC:
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) updateResultFilter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.resultFilter = resultFilterState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyBackspace:
		m.resultFilter.Query, m.resultFilter.Cursor = deleteBeforeCursor(m.resultFilter.Query, m.resultFilter.Cursor)
	case tea.KeyDelete:
		m.resultFilter.Query = deleteAtCursor(m.resultFilter.Query, m.resultFilter.Cursor)
	case tea.KeyLeft:
		m.resultFilter.Cursor = moveStringCursor(m.resultFilter.Query, m.resultFilter.Cursor, -1)
	case tea.KeyRight:
		m.resultFilter.Cursor = moveStringCursor(m.resultFilter.Query, m.resultFilter.Cursor, 1)
	case tea.KeyRunes:
		m.resultFilter.Query, m.resultFilter.Cursor = insertAtCursor(m.resultFilter.Query, m.resultFilter.Cursor, string(msg.Runes))
	}

	return m, nil
}

func (m Model) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.vimEnabled && !m.form.InsertMode {
		return m.updateFormNormal(msg)
	}

	switch msg.Type {
	case tea.KeyEsc:
		if m.vimEnabled {
			m.form.InsertMode = false
			return m, nil
		}
		m.form = formState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyBackspace:
		m.form.deleteBeforeCursor()
	case tea.KeyDelete:
		m.form.deleteAtCursor()
	case tea.KeyLeft:
		m.form.moveCursor(-1)
	case tea.KeyRight:
		m.form.moveCursor(1)
	case tea.KeyUp:
		m.form.moveField(-1)
	case tea.KeyDown:
		m.form.moveField(1)
	case tea.KeyRunes, tea.KeySpace:
		m.form.insertText(msg.String())
	case tea.KeyEnter:
		return m.submitForm()
	}

	return m, nil
}

func (m Model) updateFormNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.form = formState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		return m.submitForm()
	case tea.KeyUp:
		m.form.moveField(-1)
	case tea.KeyDown:
		m.form.moveField(1)
	case tea.KeyLeft:
		m.form.moveCursor(-1)
	case tea.KeyRight:
		m.form.moveCursor(1)
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "i":
			m.form.InsertMode = true
		case "a":
			m.form.moveCursor(1)
			m.form.InsertMode = true
		case "h":
			m.form.moveCursor(-1)
		case "l":
			m.form.moveCursor(1)
		case "j":
			m.form.moveField(1)
		case "k":
			m.form.moveField(-1)
		case "x":
			m.form.deleteAtCursor()
		}
	}
	return m, nil
}

func (m Model) submitForm() (tea.Model, tea.Cmd) {
	if m.form.Fields[m.form.Index].Required && strings.TrimSpace(m.form.Values[m.form.Index]) == "" {
		return m, nil
	}
	if m.form.Index < len(m.form.Fields)-1 {
		m.form.Index++
		m.form.clampCursor()
		return m, nil
	}
	input := m.form.input()
	capabilityID := m.form.CapabilityID
	definition, ok := m.findDefinition(capabilityID)
	m.form = formState{}
	if !ok {
		return m.startExecution(capabilityID, input, false)
	}
	return m.confirmOrStartExecution(definition, input)
}

func (m Model) updateProfileForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.vimEnabled && !m.profileForm.InsertMode {
		return m.updateProfileFormNormal(msg)
	}

	switch msg.Type {
	case tea.KeyEsc:
		m.profileForm = profileFormState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyBackspace:
		m.profileForm.Value, m.profileForm.Cursor = deleteBeforeCursor(m.profileForm.Value, m.profileForm.Cursor)
	case tea.KeyDelete:
		m.profileForm.Value = deleteAtCursor(m.profileForm.Value, m.profileForm.Cursor)
	case tea.KeyLeft:
		m.profileForm.Cursor = moveStringCursor(m.profileForm.Value, m.profileForm.Cursor, -1)
	case tea.KeyRight:
		m.profileForm.Cursor = moveStringCursor(m.profileForm.Value, m.profileForm.Cursor, 1)
	case tea.KeyRunes:
		m.profileForm.Value, m.profileForm.Cursor = insertAtCursor(m.profileForm.Value, m.profileForm.Cursor, string(msg.Runes))
	case tea.KeyEnter:
		profile := strings.TrimSpace(m.profileForm.Value)
		if profile == "" {
			return m, nil
		}
		m.profile = profile
		m.profileForm = profileFormState{}
	}

	return m, nil
}

func (m Model) updateProfileFormNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.profileForm = profileFormState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		profile := strings.TrimSpace(m.profileForm.Value)
		if profile == "" {
			return m, nil
		}
		m.profile = profile
		m.profileForm = profileFormState{}
	case tea.KeyLeft:
		m.profileForm.Cursor = moveStringCursor(m.profileForm.Value, m.profileForm.Cursor, -1)
	case tea.KeyRight:
		m.profileForm.Cursor = moveStringCursor(m.profileForm.Value, m.profileForm.Cursor, 1)
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "i":
			m.profileForm.InsertMode = true
		case "a":
			m.profileForm.Cursor = moveStringCursor(m.profileForm.Value, m.profileForm.Cursor, 1)
			m.profileForm.InsertMode = true
		case "h":
			m.profileForm.Cursor = moveStringCursor(m.profileForm.Value, m.profileForm.Cursor, -1)
		case "l":
			m.profileForm.Cursor = moveStringCursor(m.profileForm.Value, m.profileForm.Cursor, 1)
		case "x":
			m.profileForm.Value = deleteAtCursor(m.profileForm.Value, m.profileForm.Cursor)
		}
	}
	return m, nil
}

func (m Model) updateDefaultProfileForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.vimEnabled && !m.defaultProfile.InsertMode {
		return m.updateDefaultProfileFormNormal(msg)
	}

	switch msg.Type {
	case tea.KeyEsc:
		m.defaultProfile = defaultProfileFormState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyBackspace:
		m.defaultProfile.Value, m.defaultProfile.Cursor = deleteBeforeCursor(m.defaultProfile.Value, m.defaultProfile.Cursor)
	case tea.KeyDelete:
		m.defaultProfile.Value = deleteAtCursor(m.defaultProfile.Value, m.defaultProfile.Cursor)
	case tea.KeyLeft:
		m.defaultProfile.Cursor = moveStringCursor(m.defaultProfile.Value, m.defaultProfile.Cursor, -1)
	case tea.KeyRight:
		m.defaultProfile.Cursor = moveStringCursor(m.defaultProfile.Value, m.defaultProfile.Cursor, 1)
	case tea.KeyRunes:
		m.defaultProfile.Value, m.defaultProfile.Cursor = insertAtCursor(m.defaultProfile.Value, m.defaultProfile.Cursor, string(msg.Runes))
	case tea.KeyEnter:
		profile := strings.TrimSpace(m.defaultProfile.Value)
		if profile == "" {
			return m, nil
		}
		m.defaultProfile = defaultProfileFormState{}
		return m, m.saveDefaultProfile(profile)
	}

	return m, nil
}

func (m Model) updateDefaultProfileFormNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.defaultProfile = defaultProfileFormState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		profile := strings.TrimSpace(m.defaultProfile.Value)
		if profile == "" {
			return m, nil
		}
		m.defaultProfile = defaultProfileFormState{}
		return m, m.saveDefaultProfile(profile)
	case tea.KeyLeft:
		m.defaultProfile.Cursor = moveStringCursor(m.defaultProfile.Value, m.defaultProfile.Cursor, -1)
	case tea.KeyRight:
		m.defaultProfile.Cursor = moveStringCursor(m.defaultProfile.Value, m.defaultProfile.Cursor, 1)
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "i":
			m.defaultProfile.InsertMode = true
		case "a":
			m.defaultProfile.Cursor = moveStringCursor(m.defaultProfile.Value, m.defaultProfile.Cursor, 1)
			m.defaultProfile.InsertMode = true
		case "h":
			m.defaultProfile.Cursor = moveStringCursor(m.defaultProfile.Value, m.defaultProfile.Cursor, -1)
		case "l":
			m.defaultProfile.Cursor = moveStringCursor(m.defaultProfile.Value, m.defaultProfile.Cursor, 1)
		case "x":
			m.defaultProfile.Value = deleteAtCursor(m.defaultProfile.Value, m.defaultProfile.Cursor)
		}
	}
	return m, nil
}

func (m Model) saveDefaultProfile(profile string) tea.Cmd {
	return func() tea.Msg {
		result, err := config.SetDefaultProfile(m.configOptions, profile)
		return defaultProfileSavedMsg{result: result, err: err}
	}
}

func (m Model) updateConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.confirmation = confirmationState{}
	case tea.KeyEnter:
		prompt := m.confirmation
		m.confirmation = confirmationState{}
		return m.startExecution(prompt.Definition.ID, prompt.Input, true)
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyRunes:
		switch strings.ToLower(string(msg.Runes)) {
		case "y":
			prompt := m.confirmation
			m.confirmation = confirmationState{}
			return m.startExecution(prompt.Definition.ID, prompt.Input, true)
		case "n":
			m.confirmation = confirmationState{}
		}
	}

	return m, nil
}

func (m Model) confirmOrStartExecution(definition capability.Definition, input capability.Input) (tea.Model, tea.Cmd) {
	if !definition.RequiresConfirmation {
		return m.startExecution(definition.ID, input, false)
	}

	m.paletteOpen = false
	m.paletteQuery = ""
	m.paletteIndex = 0
	m.form = formState{}
	m.confirmation = confirmationState{
		Active:     true,
		Definition: definition,
		Input:      input,
	}
	return m, nil
}

func (m *Model) movePrimary(delta int) {
	if len(m.primaryActions) == 0 {
		m.primaryIndex = 0
		return
	}
	m.primaryIndex += delta
	if m.primaryIndex < 0 {
		m.primaryIndex = 0
	}
	if m.primaryIndex >= len(m.primaryActions) {
		m.primaryIndex = len(m.primaryActions) - 1
	}
}

func (m Model) clampedPrimaryIndex() int {
	if len(m.primaryActions) == 0 || m.primaryIndex < 0 {
		return 0
	}
	if m.primaryIndex >= len(m.primaryActions) {
		return len(m.primaryActions) - 1
	}
	return m.primaryIndex
}

func (m Model) selectPrimaryAction() (tea.Model, tea.Cmd) {
	if len(m.primaryActions) == 0 {
		return m, nil
	}
	selected := m.primaryActions[m.clampedPrimaryIndex()]
	m.primaryIndex = m.clampedPrimaryIndex()
	if fields := stringFields(selected); len(fields) > 0 {
		m.form = newFormState(selected.ID, fields)
		return m, nil
	}
	return m.confirmOrStartExecution(selected, capability.Input{})
}

func (m *Model) movePalette(delta int) {
	length := len(m.filteredPaletteActions())
	if length == 0 {
		m.paletteIndex = 0
		return
	}
	m.paletteIndex += delta
	if m.paletteIndex < 0 {
		m.paletteIndex = 0
	}
	if m.paletteIndex >= length {
		m.paletteIndex = length - 1
	}
}

func (m Model) startExecution(capabilityID string, input capability.Input, confirmed bool) (tea.Model, tea.Cmd) {
	m.task = taskState{Status: taskLoading, CapabilityID: capabilityID}
	m.resultFilter = resultFilterState{}
	m.confirmation = confirmationState{}
	m.paletteOpen = false
	m.paletteQuery = ""
	m.paletteIndex = 0
	ctx := m.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	m.taskCancel = cancel
	return m, m.executeAction(ctx, capabilityID, input, confirmed)
}

func (m Model) cancelTask() Model {
	if m.taskCancel != nil {
		m.taskCancel()
	}
	m.taskCancel = nil
	m.task.Status = taskCancelled
	m.task.Data = nil
	m.task.Error = nil
	m.form = formState{}
	m.paletteOpen = false
	m.paletteQuery = ""
	m.paletteIndex = 0
	m.confirmation = confirmationState{}
	m.resultFilter = resultFilterState{}
	return m
}

func (m Model) executeAction(ctx context.Context, capabilityID string, input capability.Input, confirmed bool) tea.Cmd {
	return func() tea.Msg {
		pipeline := execution.NewPipeline(m.registry)
		envelope, err := pipeline.Execute(ctx, execution.ExecuteRequest{
			CapabilityID: capabilityID,
			Input:        input,
			Profile:      profileLabel(m.profile),
			Confirmed:    confirmed,
		})
		return actionExecutedMsg{envelope: envelope, err: err}
	}
}

func (m Model) renderResultRows(builder *strings.Builder) {
	rows := resultRows(m.task.Data)
	if len(rows) == 0 {
		return
	}

	builder.WriteString("\nResult\n")
	if m.resultFilter.Active {
		fmt.Fprintf(builder, "Result Filter: %s\n", renderInputValue(m.resultFilter.Query, m.resultFilter.Cursor, true))
	}
	for _, row := range filteredResultRows(rows, m.resultFilter.Query) {
		fmt.Fprintf(builder, "- %s\n", row)
	}
	if m.resultFilter.Active {
		builder.WriteString("Press esc to close result filter.\n")
	} else {
		builder.WriteString("Press f to filter results.\n")
	}
}

func resultRows(data any) []string {
	if data == nil {
		return nil
	}

	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Map:
		return mapResultRows(value)
	case reflect.Slice, reflect.Array:
		rows := make([]string, 0, value.Len())
		for index := 0; index < value.Len(); index++ {
			rows = append(rows, fmt.Sprint(value.Index(index).Interface()))
		}
		return rows
	default:
		return []string{fmt.Sprint(data)}
	}
}

func mapResultRows(value reflect.Value) []string {
	keys := value.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
	})

	rows := make([]string, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, fmt.Sprintf("%v: %v", key.Interface(), value.MapIndex(key).Interface()))
	}
	return rows
}

func filteredResultRows(rows []string, query string) []string {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return rows
	}

	filtered := make([]string, 0, len(rows))
	for _, row := range rows {
		if strings.Contains(strings.ToLower(row), query) {
			filtered = append(filtered, row)
		}
	}
	return filtered
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

	model := NewModel(application)
	model.ctx = ctx
	_, err := tea.NewProgram(model, options...).Run()
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

func (m Model) clampedPaletteIndex(length int) int {
	if length <= 0 || m.paletteIndex < 0 {
		return 0
	}
	if m.paletteIndex >= length {
		return length - 1
	}
	return m.paletteIndex
}

func (m Model) findDefinition(capabilityID string) (capability.Definition, bool) {
	for _, definition := range m.registry.All() {
		if definition.ID == capabilityID {
			return definition, true
		}
	}
	return capability.Definition{}, false
}

func stringFields(definition capability.Definition) []capability.InputField {
	if definition.InputSchema == nil {
		return nil
	}
	fields := make([]capability.InputField, 0, len(definition.InputSchema.Fields))
	for _, field := range definition.InputSchema.Fields {
		if field.Type == capability.InputTypeString {
			fields = append(fields, field)
		}
	}
	return fields
}

func newFormState(capabilityID string, fields []capability.InputField) formState {
	return formState{
		Active:       true,
		CapabilityID: capabilityID,
		Fields:       fields,
		Values:       make([]string, len(fields)),
		Cursors:      make([]int, len(fields)),
		InsertMode:   true,
	}
}

func (m Model) keyboardModeLabel() string {
	if !m.vimEnabled {
		return "Plain"
	}
	mode := "normal"
	if m.form.Active && m.form.InsertMode ||
		m.profileForm.Active && m.profileForm.InsertMode ||
		m.defaultProfile.Active && m.defaultProfile.InsertMode {
		mode = "insert"
	}
	return "Vim " + mode
}

func (f formState) renderedValue(index int) string {
	if index < 0 || index >= len(f.Values) {
		return ""
	}
	cursor := 0
	if index < len(f.Cursors) {
		cursor = f.Cursors[index]
	}
	return renderInputValue(f.Values[index], cursor, index == f.Index)
}

func (f formState) helpText(vimEnabled bool) string {
	if !vimEnabled {
		return "Plain input: type to edit, ←/→ move cursor, ↑/↓ move fields, enter continues, esc cancels."
	}
	if f.InsertMode {
		return "Vim insert: type to edit, ←/→ move cursor, enter continues, esc returns to normal."
	}
	return "Vim normal: h/l move cursor, j/k move fields, i insert, a append, x delete, enter continues, esc cancels."
}

func textInputHelp(vimEnabled bool, insertMode bool, submit string, cancel string) string {
	if !vimEnabled {
		return fmt.Sprintf("Plain input: type to edit, ←/→ move cursor, enter to %s or esc to %s.", submit, cancel)
	}
	if insertMode {
		return fmt.Sprintf("Vim insert: type to edit, ←/→ move cursor, enter to %s or esc for normal mode.", submit)
	}
	return fmt.Sprintf("Vim normal: h/l move cursor, i insert, a append, x delete, enter to %s or esc to %s.", submit, cancel)
}

func renderInputValue(value string, cursor int, active bool) string {
	if !active {
		return value
	}
	runes := []rune(value)
	cursor = clamp(cursor, len(runes))
	return string(runes[:cursor]) + "▌" + string(runes[cursor:])
}

func (f *formState) insertText(text string) {
	if f.Index < 0 || f.Index >= len(f.Values) {
		return
	}
	f.Values[f.Index], f.Cursors[f.Index] = insertAtCursor(f.Values[f.Index], f.cursor(), text)
}

func (f *formState) deleteBeforeCursor() {
	if f.Index < 0 || f.Index >= len(f.Values) {
		return
	}
	f.Values[f.Index], f.Cursors[f.Index] = deleteBeforeCursor(f.Values[f.Index], f.cursor())
}

func (f *formState) deleteAtCursor() {
	if f.Index < 0 || f.Index >= len(f.Values) {
		return
	}
	f.Values[f.Index] = deleteAtCursor(f.Values[f.Index], f.cursor())
	f.clampCursor()
}

func (f *formState) moveCursor(delta int) {
	if f.Index < 0 || f.Index >= len(f.Values) {
		return
	}
	f.Cursors[f.Index] = moveStringCursor(f.Values[f.Index], f.cursor(), delta)
}

func (f *formState) moveField(delta int) {
	if len(f.Fields) == 0 {
		f.Index = 0
		return
	}
	f.Index = clamp(f.Index+delta, len(f.Fields)-1)
	f.clampCursor()
}

func (f *formState) cursor() int {
	if f.Index < 0 || f.Index >= len(f.Cursors) {
		return 0
	}
	return clamp(f.Cursors[f.Index], len([]rune(f.Values[f.Index])))
}

func (f *formState) clampCursor() {
	if f.Index < 0 || f.Index >= len(f.Cursors) {
		return
	}
	f.Cursors[f.Index] = clamp(f.Cursors[f.Index], len([]rune(f.Values[f.Index])))
}

func insertAtCursor(value string, cursor int, text string) (string, int) {
	runes := []rune(value)
	cursor = clamp(cursor, len(runes))
	insert := []rune(text)
	next := make([]rune, 0, len(runes)+len(insert))
	next = append(next, runes[:cursor]...)
	next = append(next, insert...)
	next = append(next, runes[cursor:]...)
	return string(next), cursor + len(insert)
}

func deleteBeforeCursor(value string, cursor int) (string, int) {
	runes := []rune(value)
	cursor = clamp(cursor, len(runes))
	if cursor == 0 {
		return value, cursor
	}
	next := append(append([]rune{}, runes[:cursor-1]...), runes[cursor:]...)
	return string(next), cursor - 1
}

func deleteAtCursor(value string, cursor int) string {
	runes := []rune(value)
	cursor = clamp(cursor, len(runes))
	if cursor >= len(runes) {
		return value
	}
	next := append(append([]rune{}, runes[:cursor]...), runes[cursor+1:]...)
	return string(next)
}

func moveStringCursor(value string, cursor int, delta int) int {
	return clamp(cursor+delta, len([]rune(value)))
}

func clamp(value int, maximum int) int {
	if value < 0 {
		return 0
	}
	if value > maximum {
		return maximum
	}
	return value
}

func (f formState) input() capability.Input {
	input := make(capability.Input, len(f.Fields))
	for index, field := range f.Fields {
		value := strings.TrimSpace(f.Values[index])
		if value == "" && !field.Required {
			continue
		}
		input[field.Name] = value
	}
	return input
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
