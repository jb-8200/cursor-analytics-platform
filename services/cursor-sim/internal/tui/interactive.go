package tui

import (
	"fmt"
	"strconv"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormModel represents the state of the interactive configuration form.
// It implements the Bubble Tea Model interface for interactive prompts.
type FormModel struct {
	// Form fields
	Developers int
	Months     int
	MaxCommits int

	// Input state
	developerInput string
	monthInput     string
	commitInput    string

	// UI state
	focusedField int
	submitted    bool
	cancelled    bool
	error        string
	mu           sync.RWMutex
}

// NewFormModel creates a new FormModel with default values.
func NewFormModel() *FormModel {
	return &FormModel{
		Developers: 10,  // Default
		Months:     6,   // Default
		MaxCommits: 500, // Default
		focusedField: 0,
		submitted:    false,
		cancelled:    false,
	}
}

// Init implements the Tea Model interface.
func (m *FormModel) Init() tea.Cmd {
	return nil
}

// Update handles user input and state changes.
func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg), nil
	}

	return m, nil
}

// handleKeyMsg processes keyboard input.
func (m *FormModel) handleKeyMsg(msg tea.KeyMsg) *FormModel {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.Cancel()
		return m
	case tea.KeyEscape:
		m.Cancel()
		return m
	case tea.KeyTab, tea.KeyDown:
		m.NextField()
	case tea.KeyShiftTab, tea.KeyUp:
		m.PrevField()
	case tea.KeyEnter:
		if m.focusedField == 2 { // At last field, submit
			if m.ValidateAll() {
				m.Submit()
			}
		} else {
			m.NextField()
		}
	case tea.KeyBackspace:
		m.Backspace()
	default:
		if msg.Runes != nil && len(msg.Runes) > 0 {
			for _, r := range msg.Runes {
				m.AddChar(r)
			}
		}
	}

	return m
}

// View returns the rendered form.
func (m *FormModel) View() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var output string

	// Banner
	output += TitleStyle.Render("Cursor Simulator - Configuration") + "\n"
	output += SubtitleStyle.Render("Configure generation parameters") + "\n\n"

	// Instructions
	output += HelpStyle.Render("Use Tab/Shift+Tab to navigate | Enter to submit | ESC to cancel") + "\n\n"

	// Developer field
	output += m.renderField(0, "Number of Developers", m.developerInput, "1-100 (default: 10)")
	output += "\n"

	// Months field
	output += m.renderField(1, "Period (Months)", m.monthInput, "1-24 (default: 6)")
	output += "\n"

	// Max commits field
	output += m.renderField(2, "Max Commits per Developer", m.commitInput, "100-2000 (default: 500)")
	output += "\n"

	// Error message
	if m.error != "" {
		output += ErrorStyle.Render("✗ "+m.error) + "\n\n"
	} else {
		output += "\n"
	}

	// Summary
	if m.IsValid() {
		output += SuccessStyle.Render("✓ Configuration valid") + "\n"
		output += m.GetSummary()
	}

	// Submit/Cancel info
	output += "\n" + SubtitleStyle.Render("Press Enter on Max Commits field to submit, or ESC to cancel")

	return output
}

// renderField renders a single form field with focus highlighting.
func (m *FormModel) renderField(fieldIndex int, label, input, help string) string {
	var value string
	switch fieldIndex {
	case 0:
		value = input
		if value == "" {
			value = strconv.Itoa(m.Developers)
		}
	case 1:
		value = input
		if value == "" {
			value = strconv.Itoa(m.Months)
		}
	case 2:
		value = input
		if value == "" {
			value = strconv.Itoa(m.MaxCommits)
		}
	}

	var fieldStyle lipgloss.Style
	if m.focusedField == fieldIndex {
		fieldStyle = PromptStyle.Bold(true).Background(lipgloss.Color("235"))
	} else {
		fieldStyle = lipgloss.NewStyle()
	}

	labelStr := fmt.Sprintf("%-35s ", label)
	inputStr := fmt.Sprintf("[ %-40s ]", value)
	helpStr := fmt.Sprintf("  %s", help)

	return fieldStyle.Render(labelStr+inputStr) + "\n" + helpStr
}

// NextField moves focus to the next field.
func (m *FormModel) NextField() {
	if m.focusedField < 2 {
		m.focusedField++
		m.error = ""
	}
}

// PrevField moves focus to the previous field.
func (m *FormModel) PrevField() {
	if m.focusedField > 0 {
		m.focusedField--
		m.error = ""
	}
}

// AddChar adds a character to the current field if it's numeric.
func (m *FormModel) AddChar(ch rune) {
	if ch < '0' || ch > '9' {
		return
	}

	switch m.focusedField {
	case 0: // Developers (1-100, max 3 digits)
		if len(m.developerInput) < 3 {
			m.developerInput += string(ch)
			if val, err := strconv.Atoi(m.developerInput); err == nil {
				m.Developers = val
			}
		}
	case 1: // Months (1-24, max 2 digits)
		if len(m.monthInput) < 2 {
			m.monthInput += string(ch)
			if val, err := strconv.Atoi(m.monthInput); err == nil {
				m.Months = val
			}
		}
	case 2: // MaxCommits (100-2000, max 4 digits)
		if len(m.commitInput) < 4 {
			m.commitInput += string(ch)
			if val, err := strconv.Atoi(m.commitInput); err == nil {
				m.MaxCommits = val
			}
		}
	}
}

// Backspace removes the last character from the current field.
func (m *FormModel) Backspace() {
	switch m.focusedField {
	case 0:
		if len(m.developerInput) > 0 {
			m.developerInput = m.developerInput[:len(m.developerInput)-1]
			if m.developerInput == "" {
				m.Developers = 10 // Reset to default
			} else if val, err := strconv.Atoi(m.developerInput); err == nil {
				m.Developers = val
			}
		}
	case 1:
		if len(m.monthInput) > 0 {
			m.monthInput = m.monthInput[:len(m.monthInput)-1]
			if m.monthInput == "" {
				m.Months = 6 // Reset to default
			} else if val, err := strconv.Atoi(m.monthInput); err == nil {
				m.Months = val
			}
		}
	case 2:
		if len(m.commitInput) > 0 {
			m.commitInput = m.commitInput[:len(m.commitInput)-1]
			if m.commitInput == "" {
				m.MaxCommits = 500 // Reset to default
			} else if val, err := strconv.Atoi(m.commitInput); err == nil {
				m.MaxCommits = val
			}
		}
	}
}

// ClearCurrentField clears the current field's input.
func (m *FormModel) ClearCurrentField() {
	switch m.focusedField {
	case 0:
		m.developerInput = ""
		m.Developers = 10
	case 1:
		m.monthInput = ""
		m.Months = 6
	case 2:
		m.commitInput = ""
		m.MaxCommits = 500
	}
}

// ValidateDevelopers validates the developers field.
func (m *FormModel) ValidateDevelopers() bool {
	return m.Developers >= 1 && m.Developers <= 100
}

// ValidateMonths validates the months field.
func (m *FormModel) ValidateMonths() bool {
	return m.Months >= 1 && m.Months <= 24
}

// ValidateMaxCommits validates the max commits field.
func (m *FormModel) ValidateMaxCommits() bool {
	return m.MaxCommits >= 100 && m.MaxCommits <= 2000
}

// ValidateAll validates all fields.
func (m *FormModel) ValidateAll() bool {
	if !m.ValidateDevelopers() {
		m.error = "Developers must be between 1 and 100"
		return false
	}
	if !m.ValidateMonths() {
		m.error = "Period must be between 1 and 24 months"
		return false
	}
	if !m.ValidateMaxCommits() {
		m.error = "Max commits must be between 100 and 2000"
		return false
	}
	m.error = ""
	return true
}

// IsValid checks if all fields are currently valid.
func (m *FormModel) IsValid() bool {
	return m.ValidateDevelopers() && m.ValidateMonths() && m.ValidateMaxCommits()
}

// Submit marks the form as submitted.
func (m *FormModel) Submit() {
	m.submitted = true
}

// Cancel marks the form as cancelled.
func (m *FormModel) Cancel() {
	m.cancelled = true
}

// IsSubmitted returns whether the form was submitted.
func (m *FormModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns whether the form was cancelled.
func (m *FormModel) IsCancelled() bool {
	return m.cancelled
}

// GetError returns the current error message.
func (m *FormModel) GetError() string {
	return m.error
}

// SetError sets the error message.
func (m *FormModel) SetError(err string) {
	m.error = err
}

// IsFieldFocused returns whether a specific field has focus.
func (m *FormModel) IsFieldFocused(fieldIndex int) bool {
	return m.focusedField == fieldIndex
}

// GetDays converts months to days (months * 30).
func (m *FormModel) GetDays() int {
	return m.Months * 30
}

// GetSummary returns a formatted summary of the configuration.
func (m *FormModel) GetSummary() string {
	return fmt.Sprintf(`
Configuration Summary:
  Developers: %d
  Period: %d months (%d days)
  Max commits: %d
`,
		m.Developers,
		m.Months,
		m.GetDays(),
		m.MaxCommits,
	)
}
