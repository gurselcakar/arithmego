package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// LaunchField identifies which field is currently focused.
type LaunchField int

const (
	FieldDifficulty LaunchField = iota
	FieldDuration
	FieldStart
)

// StartGameMsg is sent when the user starts the game.
type StartGameMsg struct {
	Mode       *modes.Mode
	Difficulty game.Difficulty
	Duration   time.Duration
}

// LaunchModel represents the launch screen.
type LaunchModel struct {
	mode *modes.Mode

	// Settings
	difficultyIndex int
	durationIndex   int

	// UI state
	focusedField LaunchField
	width        int
	height       int
}

// NewLaunch creates a new launch model for the given mode.
func NewLaunch(mode *modes.Mode) LaunchModel {
	// Find the index of the mode's default difficulty
	diffIdx := 0
	for i, d := range game.AllDifficulties() {
		if d == mode.DefaultDifficulty {
			diffIdx = i
			break
		}
	}

	// Find the index of the mode's default duration
	durIdx := modes.FindDurationIndex(mode.DefaultDuration)

	return LaunchModel{
		mode:            mode,
		difficultyIndex: diffIdx,
		durationIndex:   durIdx,
		focusedField:    FieldDifficulty,
	}
}

// Init initializes the launch model.
func (m LaunchModel) Init() tea.Cmd {
	return nil
}

// Update handles launch screen input.
func (m LaunchModel) Update(msg tea.Msg) (LaunchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.focusPrev()
		case "down", "j":
			m.focusNext()
		case "left", "h":
			m.adjustValue(-1)
		case "right", "l":
			m.adjustValue(1)
		case "enter":
			if m.focusedField == FieldStart {
				return m, m.startGame()
			}
			m.focusNext()
		case "esc":
			return m, func() tea.Msg {
				return ReturnToModesMsg{}
			}
		}
	}

	return m, nil
}

// focusPrev moves focus to the previous field.
func (m *LaunchModel) focusPrev() {
	if m.focusedField > FieldDifficulty {
		m.focusedField--
	}
}

// focusNext moves focus to the next field.
func (m *LaunchModel) focusNext() {
	if m.focusedField < FieldStart {
		m.focusedField++
	}
}

// adjustValue changes the value of the focused field.
func (m *LaunchModel) adjustValue(delta int) {
	switch m.focusedField {
	case FieldDifficulty:
		diffs := game.AllDifficulties()
		m.difficultyIndex += delta
		if m.difficultyIndex < 0 {
			m.difficultyIndex = 0
		}
		if m.difficultyIndex >= len(diffs) {
			m.difficultyIndex = len(diffs) - 1
		}
	case FieldDuration:
		durs := modes.AllowedDurations
		m.durationIndex += delta
		if m.durationIndex < 0 {
			m.durationIndex = 0
		}
		if m.durationIndex >= len(durs) {
			m.durationIndex = len(durs) - 1
		}
	}
}

// startGame creates the StartGameMsg with current settings.
func (m LaunchModel) startGame() tea.Cmd {
	diffs := game.AllDifficulties()
	durs := modes.AllowedDurations

	return func() tea.Msg {
		return StartGameMsg{
			Mode:       m.mode,
			Difficulty: diffs[m.difficultyIndex],
			Duration:   durs[m.durationIndex].Value,
		}
	}
}

// View renders the launch screen.
func (m LaunchModel) View() string {
	var b strings.Builder

	// Mode name as title
	title := styles.Bold.Render(strings.ToUpper(m.mode.Name))
	description := styles.Subtle.Render(m.mode.Description)

	// Settings
	diffs := game.AllDifficulties()
	durs := modes.AllowedDurations

	difficultyRow := m.renderSelector("Difficulty", m.difficultyIndex, diffNames(diffs), m.focusedField == FieldDifficulty)
	durationRow := m.renderSelector("Duration", m.durationIndex, durNames(durs), m.focusedField == FieldDuration)

	// Start button
	startStyle := styles.Unselected
	startText := "  Start"
	if m.focusedField == FieldStart {
		startStyle = styles.Selected
		startText = "> Start"
	}
	startButton := startStyle.Render(startText)

	// Hints
	hints := components.RenderHints([]string{"↑↓ Navigate", "←→ Adjust", "Enter Confirm", "Esc Back"})

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		description,
		"",
		"",
		difficultyRow,
		"",
		durationRow,
		"",
		"",
		startButton,
		"",
		"",
		hints,
	)

	// Center in terminal
	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	b.WriteString(content)
	return b.String()
}

// renderSelector renders a horizontal selector with arrows.
func (m LaunchModel) renderSelector(label string, index int, options []string, focused bool) string {
	leftArrow := "◀"
	rightArrow := "▶"

	if index == 0 {
		leftArrow = " "
	}
	if index >= len(options)-1 {
		rightArrow = " "
	}

	value := options[index]

	// Format: Label: ◀ Value ▶
	var row string
	if focused {
		row = fmt.Sprintf("%s: %s %s %s",
			styles.Bold.Render(label),
			styles.Accent.Render(leftArrow),
			styles.Selected.Render(value),
			styles.Accent.Render(rightArrow),
		)
	} else {
		row = fmt.Sprintf("%s: %s %s %s",
			styles.Subtle.Render(label),
			styles.Subtle.Render(leftArrow),
			styles.Unselected.Render(value),
			styles.Subtle.Render(rightArrow),
		)
	}

	return row
}

// SetSize updates the screen dimensions.
func (m *LaunchModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// diffNames returns the display names for difficulties.
func diffNames(diffs []game.Difficulty) []string {
	names := make([]string, len(diffs))
	for i, d := range diffs {
		names[i] = d.String()
	}
	return names
}

// durNames returns the display names for durations.
func durNames(durs []modes.Duration) []string {
	names := make([]string, len(durs))
	for i, d := range durs {
		names[i] = d.Label
	}
	return names
}

// ReturnToModesMsg is sent when the user wants to go back to modes screen.
type ReturnToModesMsg struct{}
