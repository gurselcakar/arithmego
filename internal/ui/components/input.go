package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputModel wraps a text input for numeric answer entry.
// Phase 10: Add multiple choice input variant (choices.go)
// This component handles typing input only
type InputModel struct {
	textInput textinput.Model
}

// NewInput creates a new input model configured for numeric entry.
func NewInput() InputModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20
	return InputModel{textInput: ti}
}

// Init initializes the input model.
func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles input messages.
func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			// Validate input: only digits allowed, except minus sign at the start.
			// This handles both single keystrokes and pasted text:
			// - Single keystroke: validates the single rune
			// - Paste: validates all runes, rejecting entire paste if any invalid
			//
			// Edge cases handled:
			// - "-5" then cursor to start and type "-": rejected (input not empty)
			// - Paste "--5": rejected at second "-" (i != 0)
			// - Paste "123": allowed (all digits)
			for i, r := range msg.Runes {
				// Allow minus sign only as first character of empty input
				if r == '-' && i == 0 && m.textInput.Value() == "" {
					continue
				}
				if r < '0' || r > '9' {
					return m, nil // Reject entire input
				}
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the input field.
func (m InputModel) View() string {
	return m.textInput.View()
}

// Value returns the current input value.
func (m InputModel) Value() string {
	return m.textInput.Value()
}

// Reset clears the input field.
func (m *InputModel) Reset() {
	m.textInput.Reset()
}

// Focus sets focus on the input.
func (m InputModel) Focus() tea.Cmd {
	return m.textInput.Focus()
}

// Blur removes focus from the input.
func (m InputModel) Blur() {
	m.textInput.Blur()
}
