package tui

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/registry"
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
	Index        int
}

type profileFormState struct {
	Active bool
	Value  string
}

type defaultProfileFormState struct {
	Active bool
	Value  string
	Status string
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
		case "/":
			m.paletteOpen = true
			m.paletteQuery = ""
			m.paletteIndex = 0
		case "f":
			if len(resultRows(m.task.Data)) > 0 && m.task.Status == taskSuccess {
				m.resultFilter = resultFilterState{Active: true}
			}
		case "p":
			m.profileForm = profileFormState{Active: true}
		case "d":
			m.defaultProfile = defaultProfileFormState{Active: true}
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

// View renders the initial task-first TUI shell.
func (m Model) View() string {
	var builder strings.Builder
	builder.WriteString("Exito Tools TUI\n")
	fmt.Fprintf(&builder, "Profile: %s\n", profileLabel(m.profile))
	fmt.Fprintf(&builder, "Primary actions: %d\n", len(m.primaryActions))

	for _, definition := range m.primaryActions {
		fmt.Fprintf(&builder, "- %s (%s)\n", definition.Title, definition.ID)
	}

	if m.task.Status != taskIdle {
		builder.WriteString("\nTask\n")
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
		builder.WriteString("\nInput Form\n")
		fmt.Fprintf(&builder, "Action: %s\n", m.form.CapabilityID)
		for index, field := range m.form.Fields {
			prefix := " "
			if index == m.form.Index {
				prefix = ">"
			}
			fmt.Fprintf(&builder, "%s %s: %s\n", prefix, field.Name, m.form.Values[index])
		}
		builder.WriteString("Press enter to continue.\n")
	}

	if m.profileForm.Active {
		builder.WriteString("\nSession Profile\n")
		fmt.Fprintf(&builder, "Current: %s\n", profileLabel(m.profile))
		fmt.Fprintf(&builder, "> New profile: %s\n", m.profileForm.Value)
		builder.WriteString("Press enter to apply or esc to cancel.\n")
	}

	if m.defaultProfile.Active {
		builder.WriteString("\nDefault Profile\n")
		fmt.Fprintf(&builder, "Current session: %s\n", profileLabel(m.profile))
		fmt.Fprintf(&builder, "> Save default as: %s\n", m.defaultProfile.Value)
		builder.WriteString("Press enter to save or esc to cancel.\n")
	} else if m.defaultProfile.Status != "" {
		builder.WriteString("\nDefault Profile\n")
		fmt.Fprintf(&builder, "%s\n", m.defaultProfile.Status)
	}

	if m.confirmation.Active {
		builder.WriteString("\nConfirm Action\n")
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
		builder.WriteString("\nCommand Palette\n")
		fmt.Fprintf(&builder, "Search: %s\n", m.paletteQuery)
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

	builder.WriteString("\nPress p to change session profile. Press d to save default profile. Press q to quit.\n")
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
		if fields := requiredStringFields(selected); len(fields) > 0 {
			m.paletteOpen = false
			m.paletteQuery = ""
			m.form = newFormState(selected.ID, fields)
			return m, nil
		}
		return m.confirmOrStartExecution(selected, capability.Input{})
	case tea.KeyUp:
		if m.paletteIndex > 0 {
			m.paletteIndex--
		}
	case tea.KeyDown:
		if m.paletteIndex < len(m.filteredPaletteActions())-1 {
			m.paletteIndex++
		}
	case tea.KeyBackspace:
		if len(m.paletteQuery) > 0 {
			runes := []rune(m.paletteQuery)
			m.paletteQuery = string(runes[:len(runes)-1])
		}
		m.paletteIndex = m.clampedPaletteIndex(len(m.filteredPaletteActions()))
	case tea.KeyRunes:
		m.paletteQuery += string(msg.Runes)
		m.paletteIndex = m.clampedPaletteIndex(len(m.filteredPaletteActions()))
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
		if len(m.resultFilter.Query) > 0 {
			runes := []rune(m.resultFilter.Query)
			m.resultFilter.Query = string(runes[:len(runes)-1])
		}
	case tea.KeyRunes:
		m.resultFilter.Query += string(msg.Runes)
	}

	return m, nil
}

func (m Model) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.form = formState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyBackspace:
		current := []rune(m.form.Values[m.form.Index])
		if len(current) > 0 {
			m.form.Values[m.form.Index] = string(current[:len(current)-1])
		}
	case tea.KeyRunes:
		m.form.Values[m.form.Index] += string(msg.Runes)
	case tea.KeyEnter:
		if strings.TrimSpace(m.form.Values[m.form.Index]) == "" {
			return m, nil
		}
		if m.form.Index < len(m.form.Fields)-1 {
			m.form.Index++
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

	return m, nil
}

func (m Model) updateProfileForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.profileForm = profileFormState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyBackspace:
		current := []rune(m.profileForm.Value)
		if len(current) > 0 {
			m.profileForm.Value = string(current[:len(current)-1])
		}
	case tea.KeyRunes:
		m.profileForm.Value += string(msg.Runes)
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

func (m Model) updateDefaultProfileForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.defaultProfile = defaultProfileFormState{}
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyBackspace:
		current := []rune(m.defaultProfile.Value)
		if len(current) > 0 {
			m.defaultProfile.Value = string(current[:len(current)-1])
		}
	case tea.KeyRunes:
		m.defaultProfile.Value += string(msg.Runes)
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
		fmt.Fprintf(builder, "Result Filter: %s\n", m.resultFilter.Query)
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

func requiredStringFields(definition capability.Definition) []capability.InputField {
	if definition.InputSchema == nil {
		return nil
	}
	fields := make([]capability.InputField, 0, len(definition.InputSchema.Fields))
	for _, field := range definition.InputSchema.Fields {
		if field.Required && field.Type == capability.InputTypeString {
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
	}
}

func (f formState) input() capability.Input {
	input := make(capability.Input, len(f.Fields))
	for index, field := range f.Fields {
		input[field.Name] = f.Values[index]
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
