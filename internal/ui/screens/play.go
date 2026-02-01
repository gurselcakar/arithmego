package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Input method index constants
const (
	inputMethodTyping = iota
	inputMethodChoice
)

// StartGameMsg is sent when the user starts the game.
type StartGameMsg struct {
	Mode        *modes.Mode
	Difficulty  game.Difficulty
	Duration    time.Duration
	InputMethod components.InputMethod
}

// PlayField identifies which field is focused in the settings panel.
type PlayField int

const (
	PlayFieldMode PlayField = iota
	PlayFieldDifficulty
	PlayFieldDuration
	PlayFieldInputMethod
)

// PlayModel represents the unified Play screen.
type PlayModel struct {
	width  int
	height int

	settingsOpen bool
	focusedField PlayField

	modes     []*modes.Mode
	modeIndex int

	difficultyIndex  int
	durationIndex    int
	inputMethodIndex int

	// Tracks if current settings match last played
	isLastPlayedSettings bool
}

// NewPlay creates a new Play model with settings from config.
func NewPlay(config *storage.Config) PlayModel {
	allModes := modes.All()

	m := PlayModel{
		modes:        allModes,
		modeIndex:    0,
		focusedField: PlayFieldMode,
	}

	if config != nil && config.HasLastPlayed() {
		for i, mode := range allModes {
			if mode.ID == config.LastPlayedModeID {
				m.modeIndex = i
				break
			}
		}
		m.difficultyIndex = findPlayDifficultyIndex(config.LastPlayedDifficulty)
		m.durationIndex = modes.FindDurationIndex(time.Duration(config.LastPlayedDurationMs) * time.Millisecond)
		if config.InputMethod == "multiple_choice" {
			m.inputMethodIndex = inputMethodChoice
		}
		m.isLastPlayedSettings = true
	} else if config != nil {
		m.difficultyIndex = findPlayDifficultyIndex(config.DefaultDifficulty)
		m.durationIndex = modes.FindDurationIndex(time.Duration(config.DefaultDurationMs) * time.Millisecond)
		if config.InputMethod == "multiple_choice" {
			m.inputMethodIndex = inputMethodChoice
		}
	} else {
		m.difficultyIndex = findPlayDifficultyIndex("Medium")
		m.durationIndex = 1
	}

	return m
}

// findPlayDifficultyIndex finds the index of a difficulty by name.
func findPlayDifficultyIndex(name string) int {
	diffs := game.AllDifficulties()
	for i, d := range diffs {
		if d.String() == name {
			return i
		}
	}
	// Fallback: find Medium explicitly
	for i, d := range diffs {
		if d == game.Medium {
			return i
		}
	}
	// Ultimate fallback: middle index
	if len(diffs) > 0 {
		return len(diffs) / 2
	}
	return 0
}

// Init initializes the Play model.
func (m PlayModel) Init() tea.Cmd {
	return nil
}

// Update handles Play screen input.
func (m PlayModel) Update(msg tea.Msg) (PlayModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "esc" {
			if m.settingsOpen {
				m.settingsOpen = false
				return m, nil
			}
			return m, func() tea.Msg { return ReturnToMenuMsg{} }
		}

		if m.settingsOpen {
			return m.updateSettingsPanel(msg)
		}
		return m.updateCollapsed(msg)
	}

	return m, nil
}

// updateCollapsed handles input when settings panel is closed.
func (m PlayModel) updateCollapsed(msg tea.KeyMsg) (PlayModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		return m, m.startGame()

	case "left", "h":
		m.adjustMode(-1)
		return m, nil

	case "right", "l":
		m.adjustMode(1)
		return m, nil

	case "tab":
		m.settingsOpen = true
		m.focusedField = PlayFieldMode
		return m, nil

	case "q":
		return m, func() tea.Msg { return ReturnToMenuMsg{} }
	}

	return m, nil
}

// updateSettingsPanel handles input when settings panel is open.
func (m PlayModel) updateSettingsPanel(msg tea.KeyMsg) (PlayModel, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.settingsOpen = false
		return m, nil

	case "enter":
		return m, m.startGame()

	case "up", "k":
		if m.focusedField > PlayFieldMode {
			m.focusedField--
		}
		return m, nil

	case "down", "j":
		if m.focusedField < PlayFieldInputMethod {
			m.focusedField++
		}
		return m, nil

	case "left", "h":
		m.adjustFieldValue(-1)
		return m, nil

	case "right", "l":
		m.adjustFieldValue(1)
		return m, nil
	}

	return m, nil
}

// adjustMode changes the mode by delta.
func (m *PlayModel) adjustMode(delta int) {
	if len(m.modes) == 0 {
		return
	}
	m.modeIndex += delta
	if m.modeIndex < 0 {
		m.modeIndex = 0
	}
	if m.modeIndex >= len(m.modes) {
		m.modeIndex = len(m.modes) - 1
	}
	m.isLastPlayedSettings = false
}

// adjustFieldValue changes the focused field value by delta.
func (m *PlayModel) adjustFieldValue(delta int) {
	switch m.focusedField {
	case PlayFieldMode:
		m.adjustMode(delta)

	case PlayFieldDifficulty:
		diffs := game.AllDifficulties()
		if len(diffs) == 0 {
			return
		}
		m.difficultyIndex += delta
		if m.difficultyIndex < 0 {
			m.difficultyIndex = 0
		}
		if m.difficultyIndex >= len(diffs) {
			m.difficultyIndex = len(diffs) - 1
		}
		m.isLastPlayedSettings = false

	case PlayFieldDuration:
		durs := modes.AllowedDurations
		if len(durs) == 0 {
			return
		}
		m.durationIndex += delta
		if m.durationIndex < 0 {
			m.durationIndex = 0
		}
		if m.durationIndex >= len(durs) {
			m.durationIndex = len(durs) - 1
		}
		m.isLastPlayedSettings = false

	case PlayFieldInputMethod:
		if m.inputMethodIndex == inputMethodTyping {
			m.inputMethodIndex = inputMethodChoice
		} else {
			m.inputMethodIndex = inputMethodTyping
		}
		m.isLastPlayedSettings = false
	}
}

// startGame creates the StartGameMsg with current settings.
func (m PlayModel) startGame() tea.Cmd {
	if len(m.modes) == 0 || m.modeIndex >= len(m.modes) {
		return nil
	}

	diffs := game.AllDifficulties()
	durs := modes.AllowedDurations

	// Bounds checking
	if len(diffs) == 0 || len(durs) == 0 {
		return nil
	}
	if m.difficultyIndex >= len(diffs) {
		m.difficultyIndex = len(diffs) - 1
	}
	if m.durationIndex >= len(durs) {
		m.durationIndex = len(durs) - 1
	}

	return func() tea.Msg {
		return StartGameMsg{
			Mode:        m.modes[m.modeIndex],
			Difficulty:  diffs[m.difficultyIndex],
			Duration:    durs[m.durationIndex].Value,
			InputMethod: m.currentInputMethod(),
		}
	}
}

// currentInputMethod returns the currently selected input method.
func (m PlayModel) currentInputMethod() components.InputMethod {
	if m.inputMethodIndex == inputMethodChoice {
		return components.InputMultipleChoice
	}
	return components.InputTyping
}

// inputMethodLabel returns the display label for the current input method.
func (m PlayModel) inputMethodLabel() string {
	if m.inputMethodIndex == inputMethodChoice {
		return "Choice"
	}
	return "Typing"
}

// settingsSummary returns the formatted settings summary string.
func (m PlayModel) settingsSummary() (diffName, durShort, inputLabel string) {
	diffs := game.AllDifficulties()
	durs := modes.AllowedDurations

	if m.difficultyIndex < len(diffs) {
		diffName = diffs[m.difficultyIndex].String()
	}
	if m.durationIndex < len(durs) {
		durShort = m.formatDurationShort(durs[m.durationIndex].Value)
	}
	inputLabel = m.inputMethodLabel()
	return
}

// formatDurationShort returns a short format for the duration (e.g., "60s").
func (m PlayModel) formatDurationShort(d time.Duration) string {
	return fmt.Sprintf("%ds", int(d.Seconds()))
}

// View renders the Play screen.
func (m PlayModel) View() string {
	if m.settingsOpen {
		return m.viewWithSettingsPanel()
	}
	return m.viewCollapsed()
}

// viewCollapsed renders the collapsed view with mode selector and settings summary.
func (m PlayModel) viewCollapsed() string {
	contentWidth := m.width
	if contentWidth == 0 {
		contentWidth = 80
	}

	title := styles.Bold.Render("PLAY")

	var modeDesc string
	if len(m.modes) > 0 && m.modeIndex < len(m.modes) {
		modeDesc = m.modes[m.modeIndex].Description
	}

	modeSelector := components.RenderSelector(m.modeIndex, m.modeNames(), components.SelectorOptions{
		Focused: true,
	})
	// Fixed-width container prevents layout shift when mode names vary in length
	modeSelectorStyled := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(modeSelector)

	modeDescStyled := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Faint(true).
		Render(modeDesc)

	diffName, durShort, inputLabel := m.settingsSummary()
	summary := styles.Subtle.Render(fmt.Sprintf("%s \u00b7 %s \u00b7 %s", diffName, durShort, inputLabel))

	// Reserve space to maintain consistent height between views
	var lastPlayedLine string
	if m.isLastPlayedSettings {
		lastPlayedLine = styles.Dim.Render("\u21bb last played")
	} else {
		lastPlayedLine = " "
	}

	startButton := styles.Selected.Render("[ START GAME ]")

	var centerParts []string
	centerParts = append(centerParts, title, "", modeSelectorStyled, modeDescStyled, "", summary)
	centerParts = append(centerParts, lastPlayedLine)
	centerParts = append(centerParts, "", startButton)

	centerContent := lipgloss.JoinVertical(lipgloss.Center, centerParts...)

	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "Esc", Action: "Back"},
		{Key: "\u2190\u2192", Action: "Mode"},
		{Key: "Tab", Action: "Settings"},
		{Key: "Enter", Action: "Start"},
	})

	if m.width > 0 && m.height > 0 {
		hintsHeight := lipgloss.Height(hints)
		bottomPadding := 1
		availableHeight := m.height - hintsHeight - bottomPadding

		centeredMain := lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, centerContent)
		centeredHints := lipgloss.Place(m.width, hintsHeight+bottomPadding, lipgloss.Center, lipgloss.Top, hints)

		return lipgloss.JoinVertical(lipgloss.Left, centeredMain, centeredHints)
	}

	return lipgloss.JoinVertical(lipgloss.Center, centerContent, "", "", hints)
}

// viewWithSettingsPanel renders the view with the settings panel overlay.
func (m PlayModel) viewWithSettingsPanel() string {
	panel := m.renderSettingsPanel()

	contentWidth := m.width
	if contentWidth == 0 {
		contentWidth = 80
	}

	title := styles.Bold.Render("PLAY")

	var modeDesc string
	if len(m.modes) > 0 && m.modeIndex < len(m.modes) {
		modeDesc = m.modes[m.modeIndex].Description
	}

	modeSelector := components.RenderSelector(m.modeIndex, m.modeNames(), components.SelectorOptions{
		Focused: false,
	})
	modeSelectorStyled := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(modeSelector)

	// Render without Faint styling to avoid ANSI corruption in overlay byte-slicing
	modeDescStyled := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(modeDesc)

	diffName, durShort, inputLabel := m.settingsSummary()
	summary := fmt.Sprintf("%s \u00b7 %s \u00b7 %s", diffName, durShort, inputLabel)

	// Reserve space to maintain consistent height between views
	lastPlayedLine := " "

	startButton := "[ START GAME ]"

	centerContent := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		modeSelectorStyled,
		modeDescStyled,
		"",
		summary,
		lastPlayedLine,
		"",
		startButton,
	)

	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "Esc", Action: "Close"},
		{Key: "\u2191\u2193", Action: "Field"},
		{Key: "\u2190\u2192", Action: "Change"},
		{Key: "Enter", Action: "Start"},
	})

	if m.width > 0 && m.height > 0 {
		hintsHeight := lipgloss.Height(hints)
		bottomPadding := 1
		panelWidth := 30
		availableHeight := m.height - hintsHeight - bottomPadding

		panelContent := lipgloss.NewStyle().PaddingLeft(2).Render(panel)
		panelStyled := lipgloss.Place(panelWidth, availableHeight, lipgloss.Left, lipgloss.Center, panelContent)

		centeredContent := lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, centerContent)

		// Overlay panel on left side of centered content
		panelLines := strings.Split(panelStyled, "\n")
		contentLines := strings.Split(centeredContent, "\n")

		var resultLines []string
		for i := 0; i < len(contentLines); i++ {
			panelLine := ""
			if i < len(panelLines) {
				panelLine = panelLines[i]
			}
			contentLine := contentLines[i]

			panelActualWidth := lipgloss.Width(panelLine)
			if panelActualWidth < panelWidth {
				panelLine = panelLine + strings.Repeat(" ", panelWidth-panelActualWidth)
			}

			if lipgloss.Width(contentLine) > panelWidth {
				resultLines = append(resultLines, panelLine+contentLine[panelWidth:])
			} else {
				resultLines = append(resultLines, panelLine)
			}
		}
		mainArea := strings.Join(resultLines, "\n")

		centeredHints := lipgloss.Place(m.width, hintsHeight+bottomPadding, lipgloss.Center, lipgloss.Top, hints)

		return lipgloss.JoinVertical(lipgloss.Left, mainArea, centeredHints)
	}

	return lipgloss.JoinVertical(lipgloss.Center, panel, "", centerContent, "", hints)
}

// renderSettingsPanel renders the settings panel content.
func (m PlayModel) renderSettingsPanel() string {
	var b strings.Builder

	b.WriteString(styles.Bold.Render("SETTINGS"))
	b.WriteString("\n\n")

	// Mode
	modeSectionStyle := styles.Bold
	if m.focusedField != PlayFieldMode {
		modeSectionStyle = styles.Subtle
	}
	b.WriteString(modeSectionStyle.Render("Mode"))
	b.WriteString("\n")
	b.WriteString(m.renderHorizontalSelector(
		m.modeIndex,
		m.modeNames(),
		m.focusedField == PlayFieldMode,
	))
	b.WriteString("\n\n")

	// Difficulty
	diffSectionStyle := styles.Bold
	if m.focusedField != PlayFieldDifficulty {
		diffSectionStyle = styles.Subtle
	}
	b.WriteString(diffSectionStyle.Render("Difficulty"))
	b.WriteString("\n")
	b.WriteString(m.renderHorizontalSelector(
		m.difficultyIndex,
		m.difficultyNames(),
		m.focusedField == PlayFieldDifficulty,
	))
	b.WriteString("\n\n")

	// Duration
	durSectionStyle := styles.Bold
	if m.focusedField != PlayFieldDuration {
		durSectionStyle = styles.Subtle
	}
	b.WriteString(durSectionStyle.Render("Duration"))
	b.WriteString("\n")
	b.WriteString(m.renderHorizontalSelector(
		m.durationIndex,
		m.durationNames(),
		m.focusedField == PlayFieldDuration,
	))
	b.WriteString("\n\n")

	// Input
	inputSectionStyle := styles.Bold
	if m.focusedField != PlayFieldInputMethod {
		inputSectionStyle = styles.Subtle
	}
	b.WriteString(inputSectionStyle.Render("Input"))
	b.WriteString("\n")
	b.WriteString(m.renderHorizontalSelector(
		m.inputMethodIndex,
		[]string{"Typing", "Choice"},
		m.focusedField == PlayFieldInputMethod,
	))

	return b.String()
}

// renderHorizontalSelector renders a horizontal selector component.
func (m PlayModel) renderHorizontalSelector(index int, options []string, focused bool) string {
	return components.RenderSelector(index, options, components.SelectorOptions{
		Prefix:  "  ",
		Focused: focused,
	})
}

// modeNames returns the names of all modes.
func (m PlayModel) modeNames() []string {
	names := make([]string, len(m.modes))
	for i, mode := range m.modes {
		names[i] = mode.Name
	}
	return names
}

// difficultyNames returns the names of all difficulties.
func (m PlayModel) difficultyNames() []string {
	diffs := game.AllDifficulties()
	names := make([]string, len(diffs))
	for i, d := range diffs {
		names[i] = d.String()
	}
	return names
}

// durationNames returns the labels for all durations.
func (m PlayModel) durationNames() []string {
	durs := modes.AllowedDurations
	names := make([]string, len(durs))
	for i, d := range durs {
		names[i] = d.Label
	}
	return names
}

// SetSize updates the screen dimensions.
func (m *PlayModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
