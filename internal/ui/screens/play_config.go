package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/gen"
	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// StartGameMsg is sent when the user starts the game.
type StartGameMsg struct {
	Mode        *modes.Mode
	Difficulty  game.Difficulty
	Duration    time.Duration
	InputMethod components.InputMethod
}

// BackToBrowseMsg is sent when the user wants to go back to the browse screen.
type BackToBrowseMsg struct{}

// PlayConfigField identifies which field is focused.
type PlayConfigField int

const (
	PlayConfigFieldDifficulty PlayConfigField = iota
	PlayConfigFieldDuration
	PlayConfigFieldInputMethod
)

const playConfigFieldCount = 3


// PlayConfigModel represents the Configure & Start screen (Step 2 of play flow).
type PlayConfigModel struct {
	width  int
	height int

	viewport      viewport.Model
	viewportReady bool

	selectedMode *modes.Mode

	difficultyIndex  int
	durationIndex    int
	inputMethodIndex int
	focusedField     PlayConfigField

	sampleQuestion string // Preview equation for current difficulty

	config *storage.Config
}

// NewPlayConfig creates a new PlayConfigModel for the given mode.
func NewPlayConfig(mode *modes.Mode, config *storage.Config) PlayConfigModel {
	m := PlayConfigModel{
		selectedMode: mode,
		config:       config,
		focusedField: PlayConfigFieldDifficulty,
		viewport:     viewport.New(0, 0),
	}

	// Initialize from config or defaults
	if config != nil && config.HasLastPlayed() && config.LastPlayedModeID == mode.ID {
		// Restore last played settings for this mode
		m.difficultyIndex = findDifficultyIndex(config.LastPlayedDifficulty)
		m.durationIndex = modes.FindDurationIndex(time.Duration(config.LastPlayedDurationMs) * time.Millisecond)
		if config.InputMethod == "multiple_choice" {
			m.inputMethodIndex = 1
		}
	} else if config != nil {
		// Use default settings
		m.difficultyIndex = findDifficultyIndex(config.DefaultDifficulty)
		m.durationIndex = modes.FindDurationIndex(time.Duration(config.DefaultDurationMs) * time.Millisecond)
		if config.InputMethod == "multiple_choice" {
			m.inputMethodIndex = 1
		}
	} else {
		// Fallback defaults
		m.difficultyIndex = findDifficultyIndex("Medium")
		m.durationIndex = 1 // 60s
		m.inputMethodIndex = 0 // Typing
	}

	// Generate initial sample question
	m.generateSampleQuestion()

	return m
}


// generateSampleQuestion creates a sample question for the current difficulty.
func (m *PlayConfigModel) generateSampleQuestion() {
	if m.selectedMode == nil || m.selectedMode.GeneratorLabel == "" {
		m.sampleQuestion = ""
		return
	}

	diffs := game.AllDifficulties()
	if m.difficultyIndex >= len(diffs) {
		m.sampleQuestion = ""
		return
	}

	diff := diffs[m.difficultyIndex]
	g, ok := gen.Get(m.selectedMode.GeneratorLabel)
	if !ok {
		m.sampleQuestion = ""
		return
	}
	q := g.Generate(diff)
	if q == nil {
		m.sampleQuestion = ""
		return
	}
	m.sampleQuestion = q.Display
}

// Init initializes the PlayConfigModel.
func (m PlayConfigModel) Init() tea.Cmd {
	return nil
}

// Update handles PlayConfigModel input.
func (m PlayConfigModel) Update(msg tea.Msg) (PlayConfigModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back to browse
			return m, func() tea.Msg { return BackToBrowseMsg{} }

		case "up", "k":
			m.focusPrev()
			m.updateViewportContent()
			return m, nil

		case "down", "j":
			m.focusNext()
			m.updateViewportContent()
			return m, nil

		case "left", "h":
			m.adjustValue(-1)
			m.updateViewportContent()
			return m, nil

		case "right", "l":
			m.adjustValue(1)
			m.updateViewportContent()
			return m, nil

		case "enter":
			return m, m.startGame()
		}
	}

	return m, nil
}

// focusPrev moves focus to the previous field.
func (m *PlayConfigModel) focusPrev() {
	if m.focusedField > 0 {
		m.focusedField--
	}
}

// focusNext moves focus to the next field.
func (m *PlayConfigModel) focusNext() {
	if m.focusedField < playConfigFieldCount-1 {
		m.focusedField++
	}
}

// adjustValue changes the value of the focused field.
func (m *PlayConfigModel) adjustValue(delta int) {
	switch m.focusedField {
	case PlayConfigFieldDifficulty:
		diffs := game.AllDifficulties()
		oldIndex := m.difficultyIndex
		m.difficultyIndex += delta
		if m.difficultyIndex < 0 {
			m.difficultyIndex = 0
		}
		if m.difficultyIndex >= len(diffs) {
			m.difficultyIndex = len(diffs) - 1
		}
		// Regenerate sample if difficulty changed
		if m.difficultyIndex != oldIndex {
			m.generateSampleQuestion()
		}

	case PlayConfigFieldDuration:
		durs := modes.AllowedDurations
		m.durationIndex += delta
		if m.durationIndex < 0 {
			m.durationIndex = 0
		}
		if m.durationIndex >= len(durs) {
			m.durationIndex = len(durs) - 1
		}

	case PlayConfigFieldInputMethod:
		if m.inputMethodIndex == 0 {
			m.inputMethodIndex = 1
		} else {
			m.inputMethodIndex = 0
		}
	}
}

// startGame creates the StartGameMsg with current settings.
func (m PlayConfigModel) startGame() tea.Cmd {
	if m.selectedMode == nil {
		return nil
	}

	diffs := game.AllDifficulties()
	durs := modes.AllowedDurations

	if len(diffs) == 0 || len(durs) == 0 {
		return nil
	}

	diffIndex := m.difficultyIndex
	if diffIndex >= len(diffs) {
		diffIndex = len(diffs) - 1
	}

	durIndex := m.durationIndex
	if durIndex >= len(durs) {
		durIndex = len(durs) - 1
	}

	inputMethod := components.InputTyping
	if m.inputMethodIndex == 1 {
		inputMethod = components.InputMultipleChoice
	}

	return func() tea.Msg {
		return StartGameMsg{
			Mode:        m.selectedMode,
			Difficulty:  diffs[diffIndex],
			Duration:    durs[durIndex].Value,
			InputMethod: inputMethod,
		}
	}
}

// View renders the PlayConfigModel.
func (m PlayConfigModel) View() string {
	if !m.viewportReady {
		return "Loading..."
	}

	hints := m.getHints()

	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		lipgloss.Place(m.width, components.HintsHeight, lipgloss.Center, lipgloss.Center, hints),
	)
}

// getHints returns the hints for the config screen.
func (m PlayConfigModel) getHints() string {
	return components.RenderHintsResponsive([]components.Hint{
		{Key: "Esc", Action: "Back"},
		{Key: "↑↓", Action: "Navigate"},
		{Key: "←→", Action: "Change"},
		{Key: "Enter", Action: "Start"},
	}, m.width)
}

// SetSize sets the screen dimensions.
func (m *PlayConfigModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	viewportHeight := m.calculateViewportHeight()

	components.SetViewportSize(&m.viewport, &m.viewportReady, m.width, viewportHeight)

	m.updateViewportContent()
}

// calculateViewportHeight returns the viewport height.
func (m PlayConfigModel) calculateViewportHeight() int {
	viewportHeight := m.height - components.HintsHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	return viewportHeight
}

// updateViewportContent updates the viewport with the current content.
func (m *PlayConfigModel) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	content := m.renderContent()
	m.viewport.SetContent(content)
}

// renderContent renders the main content for the viewport.
func (m PlayConfigModel) renderContent() string {
	if m.selectedMode == nil {
		return "No mode selected"
	}

	// Title (mode name) - prominent uppercase
	title := styles.Logo.Render(strings.ToUpper(m.selectedMode.Name))

	// Description
	desc := styles.Dim.Render(m.selectedMode.Description)

	// Settings
	diffs := game.AllDifficulties()
	durs := modes.AllowedDurations
	inputOptions := []string{"Typing", "Choice"}

	// Calculate widths for alignment
	labels := []string{"Difficulty", "Duration", "Input"}
	labelWidth := maxLen(labels)

	allValues := []string{}
	allValues = append(allValues, difficultyNames(diffs)...)
	allValues = append(allValues, durationShortNames(durs)...)
	allValues = append(allValues, inputOptions...)
	valueWidth := maxLen(allValues)

	// Focus prefix helper
	focusPrefix := func(focused bool) string {
		if focused {
			return styles.Accent.Render("> ")
		}
		return "  "
	}

	// Difficulty row
	difficultyRow := focusPrefix(m.focusedField == PlayConfigFieldDifficulty) +
		components.RenderSelector(m.difficultyIndex, difficultyNames(diffs), components.SelectorOptions{
			Label:      "Difficulty",
			LabelWidth: labelWidth,
			ValueWidth: valueWidth,
			Focused:    m.focusedField == PlayConfigFieldDifficulty,
		})

	// Duration row
	durationRow := focusPrefix(m.focusedField == PlayConfigFieldDuration) +
		components.RenderSelector(m.durationIndex, durationShortNames(durs), components.SelectorOptions{
			Label:      "Duration",
			LabelWidth: labelWidth,
			ValueWidth: valueWidth,
			Focused:    m.focusedField == PlayConfigFieldDuration,
		})

	// Input row
	inputRow := focusPrefix(m.focusedField == PlayConfigFieldInputMethod) +
		components.RenderSelector(m.inputMethodIndex, inputOptions, components.SelectorOptions{
			Label:      "Input",
			LabelWidth: labelWidth,
			ValueWidth: valueWidth,
			Focused:    m.focusedField == PlayConfigFieldInputMethod,
		})

	// Settings block (no box)
	settingsBlock := lipgloss.JoinVertical(lipgloss.Left,
		difficultyRow,
		durationRow,
		inputRow,
	)

	// Preview box with sample equation
	previewBoxWidth := min(m.width-4, 30)
	previewBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(previewBoxWidth).
		Padding(1, 2).
		Align(lipgloss.Center).
		Render(styles.Bold.Render(m.sampleQuestion))

	// Build content with better spacing
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		desc,
		"",
		"",
		settingsBlock,
		"",
		previewBox,
	)

	// Center in viewport
	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// Helper functions

func durationShortNames(durs []modes.Duration) []string {
	names := make([]string, len(durs))
	for i, d := range durs {
		names[i] = fmt.Sprintf("%ds", int(d.Value.Seconds()))
	}
	return names
}
