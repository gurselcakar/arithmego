package components

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// InputMethod represents the answer input type.
type InputMethod int

const (
	InputTyping InputMethod = iota
	InputMultipleChoice
)

func (m InputMethod) String() string {
	switch m {
	case InputTyping:
		return "Typing"
	case InputMultipleChoice:
		return "Multiple Choice"
	default:
		return "Unknown"
	}
}

// ParseInputMethod converts a string to an InputMethod.
func ParseInputMethod(s string) InputMethod {
	switch s {
	case "Multiple Choice", "multiple_choice":
		return InputMultipleChoice
	default:
		return InputTyping
	}
}

// ChoiceSelectedMsg is sent when a choice is selected.
type ChoiceSelectedMsg struct {
	Value int
	Index int
}

// ChoicesModel handles multiple choice input.
type ChoicesModel struct {
	choices      []int
	correctIndex int
	selected     int  // -1 = none, 0-3 = selected choice
	focused      bool
	errorIndex   int  // -1 = none, 0-3 = wrong choice (shown in red)
}

// NewChoices creates a new choices model.
func NewChoices() ChoicesModel {
	return ChoicesModel{
		choices:      make([]int, 4),
		correctIndex: 0,
		selected:     -1,
		focused:      true,
		errorIndex:   -1,
	}
}

// Init initializes the choices model.
func (m ChoicesModel) Init() tea.Cmd {
	return nil
}

// Update handles choice selection messages.
func (m ChoicesModel) Update(msg tea.Msg) (ChoicesModel, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			if len(m.choices) > 0 {
				m.selected = 0
				return m, m.selectChoice(0)
			}
		case "2":
			if len(m.choices) > 1 {
				m.selected = 1
				return m, m.selectChoice(1)
			}
		case "3":
			if len(m.choices) > 2 {
				m.selected = 2
				return m, m.selectChoice(2)
			}
		case "4":
			if len(m.choices) > 3 {
				m.selected = 3
				return m, m.selectChoice(3)
			}
		}
	}

	return m, nil
}

// selectChoice returns a command that sends a ChoiceSelectedMsg.
func (m ChoicesModel) selectChoice(index int) tea.Cmd {
	value := m.choices[index]
	return func() tea.Msg {
		return ChoiceSelectedMsg{Value: value, Index: index}
	}
}

// View renders the choices as a horizontal row.
func (m ChoicesModel) View() string {
	if len(m.choices) == 0 {
		return ""
	}

	var parts []string
	for i, choice := range m.choices {
		keyLabel := fmt.Sprintf("[%d]", i+1)
		valueStr := strconv.Itoa(choice)

		var style lipgloss.Style
		if m.errorIndex == i {
			// Wrong choice - show in red
			style = styles.Incorrect
		} else if m.selected == i {
			style = styles.Selected
		} else if m.focused {
			style = styles.Normal
		} else {
			style = styles.Dim
		}

		part := style.Render(fmt.Sprintf("%s %s", keyLabel, valueStr))
		parts = append(parts, part)
	}

	// Join with spacing
	return lipgloss.JoinHorizontal(lipgloss.Center, joinWithSpacing(parts, "  ")...)
}

// joinWithSpacing inserts spacing between parts.
func joinWithSpacing(parts []string, spacing string) []string {
	if len(parts) == 0 {
		return parts
	}

	result := make([]string, 0, len(parts)*2-1)
	for i, part := range parts {
		result = append(result, part)
		if i < len(parts)-1 {
			result = append(result, spacing)
		}
	}
	return result
}

// Value returns the selected choice value as a string, or empty if none selected.
func (m ChoicesModel) Value() string {
	if m.selected >= 0 && m.selected < len(m.choices) {
		return strconv.Itoa(m.choices[m.selected])
	}
	return ""
}

// Reset prepares for a new question.
func (m *ChoicesModel) Reset() {
	m.selected = -1
	m.errorIndex = -1
}

// SetError marks a choice as incorrect (shown in red).
func (m *ChoicesModel) SetError(index int) {
	m.errorIndex = index
}

// ClearError removes the error marking.
func (m *ChoicesModel) ClearError() {
	m.errorIndex = -1
}

// SetChoices updates the available choices and correct index.
// If correctIndex is out of bounds, it defaults to 0.
func (m *ChoicesModel) SetChoices(choices []int, correctIndex int) {
	m.choices = choices
	if correctIndex < 0 || correctIndex >= len(choices) {
		correctIndex = 0
	}
	m.correctIndex = correctIndex
	m.selected = -1
}

// Focus sets focus on the choices.
func (m *ChoicesModel) Focus() tea.Cmd {
	m.focused = true
	return nil
}

// Blur removes focus from the choices.
func (m *ChoicesModel) Blur() {
	m.focused = false
}

// IsCorrect returns whether the selected choice is correct.
func (m ChoicesModel) IsCorrect() bool {
	return m.selected == m.correctIndex
}

// SelectedIndex returns the selected index, or -1 if none.
func (m ChoicesModel) SelectedIndex() int {
	return m.selected
}

// CorrectIndex returns the index of the correct answer.
func (m ChoicesModel) CorrectIndex() int {
	return m.correctIndex
}
