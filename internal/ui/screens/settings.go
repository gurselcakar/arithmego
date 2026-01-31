package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// SettingsField identifies which field is currently focused.
type SettingsField int

const (
	SettingsFieldDifficulty SettingsField = iota
	SettingsFieldDuration
	SettingsFieldInputMethod
	SettingsFieldAutoUpdate
	SettingsFieldSkipQuitConfirm
)

const settingsFieldCount = 5

// SettingsModel represents the settings screen.
type SettingsModel struct {
	config *storage.Config

	// UI state
	focusedField     SettingsField
	difficultyIndex  int
	durationIndex    int
	inputMethodIndex int
	width            int
	height           int
}

// NewSettings creates a new settings model.
func NewSettings(config *storage.Config) SettingsModel {
	if config == nil {
		config = storage.NewConfig()
	}

	// Find indices for current values
	diffIdx := findDifficultyIndex(config.DefaultDifficulty)
	durIdx := findDurationIndexByMs(config.DefaultDurationMs)
	inputIdx := 0
	if config.InputMethod == "multiple_choice" {
		inputIdx = 1
	}

	return SettingsModel{
		config:           config,
		difficultyIndex:  diffIdx,
		durationIndex:    durIdx,
		inputMethodIndex: inputIdx,
		focusedField:     SettingsFieldDifficulty,
	}
}

// Init initializes the settings model.
func (m SettingsModel) Init() tea.Cmd {
	return nil
}

// Update handles settings screen input.
func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
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
		case "enter", " ":
			if m.focusedField == SettingsFieldAutoUpdate {
				m.toggleAutoUpdate()
			} else if m.focusedField == SettingsFieldSkipQuitConfirm {
				m.toggleSkipQuitConfirm()
			}
		case "esc":
			return m, func() tea.Msg {
				return ReturnToMenuMsg{}
			}
		}
	}

	return m, nil
}

// focusPrev moves focus to the previous field.
func (m *SettingsModel) focusPrev() {
	if m.focusedField > 0 {
		m.focusedField--
	}
}

// focusNext moves focus to the next field.
func (m *SettingsModel) focusNext() {
	if m.focusedField < settingsFieldCount-1 {
		m.focusedField++
	}
}

// adjustValue changes the value of the focused field.
func (m *SettingsModel) adjustValue(delta int) {
	switch m.focusedField {
	case SettingsFieldDifficulty:
		diffs := game.AllDifficulties()
		m.difficultyIndex += delta
		if m.difficultyIndex < 0 {
			m.difficultyIndex = 0
		}
		if m.difficultyIndex >= len(diffs) {
			m.difficultyIndex = len(diffs) - 1
		}
		m.config.DefaultDifficulty = diffs[m.difficultyIndex].String()
		m.saveConfig()

	case SettingsFieldDuration:
		durs := modes.AllowedDurations
		m.durationIndex += delta
		if m.durationIndex < 0 {
			m.durationIndex = 0
		}
		if m.durationIndex >= len(durs) {
			m.durationIndex = len(durs) - 1
		}
		m.config.DefaultDurationMs = durs[m.durationIndex].Value.Milliseconds()
		m.saveConfig()

	case SettingsFieldInputMethod:
		m.toggleInputMethod()

	case SettingsFieldAutoUpdate:
		m.toggleAutoUpdate()

	case SettingsFieldSkipQuitConfirm:
		m.toggleSkipQuitConfirm()
	}
}

// toggleInputMethod switches between typing and multiple choice modes.
func (m *SettingsModel) toggleInputMethod() {
	if m.inputMethodIndex == 0 {
		m.inputMethodIndex = 1
		m.config.InputMethod = "multiple_choice"
	} else {
		m.inputMethodIndex = 0
		m.config.InputMethod = "typing"
	}
	m.saveConfig()
}

// toggleAutoUpdate toggles the auto-update preference.
func (m *SettingsModel) toggleAutoUpdate() {
	m.config.AutoUpdate = !m.config.AutoUpdate
	m.saveConfig()
}

// toggleSkipQuitConfirm toggles the skip quit confirmation preference.
func (m *SettingsModel) toggleSkipQuitConfirm() {
	m.config.SkipQuitConfirmation = !m.config.SkipQuitConfirmation
	m.saveConfig()
}

// saveConfig persists the current config to disk.
func (m *SettingsModel) saveConfig() {
	// Ignore errors - settings are non-critical
	_ = storage.SaveConfig(m.config)
}

// View renders the settings screen.
func (m SettingsModel) View() string {
	var b strings.Builder

	// Title
	title := styles.Bold.Render("SETTINGS")

	// Gather all data
	diffs := game.AllDifficulties()
	durs := modes.AllowedDurations
	inputOptions := []string{"Typing", "Multiple Choice"}

	// All labels used in settings
	labels := []string{"Difficulty", "Duration", "Input", "Auto-update", "Skip quit confirm"}

	// All possible values across all selectors
	allValues := []string{}
	allValues = append(allValues, settingsDifficultyNames(diffs)...)
	allValues = append(allValues, settingsDurationNames(durs)...)
	allValues = append(allValues, inputOptions...)

	// Calculate widths dynamically
	labelWidth := maxLen(labels)
	valueWidth := maxLen(allValues)

	// Separator width: label + 2 spaces + arrow + space + value + space + arrow
	separatorWidth := labelWidth + 2 + 1 + 1 + valueWidth + 1 + 1
	separator := styles.Dim.Render(strings.Repeat("─", separatorWidth))

	difficultyRow := components.RenderSelector(m.difficultyIndex, settingsDifficultyNames(diffs), components.SelectorOptions{
		Label:      "Difficulty",
		LabelWidth: labelWidth,
		ValueWidth: valueWidth,
		Focused:    m.focusedField == SettingsFieldDifficulty,
	})

	durationRow := components.RenderSelector(m.durationIndex, settingsDurationNames(durs), components.SelectorOptions{
		Label:      "Duration",
		LabelWidth: labelWidth,
		ValueWidth: valueWidth,
		Focused:    m.focusedField == SettingsFieldDuration,
	})

	inputMethodRow := components.RenderSelector(m.inputMethodIndex, inputOptions, components.SelectorOptions{
		Label:      "Input",
		LabelWidth: labelWidth,
		ValueWidth: valueWidth,
		Focused:    m.focusedField == SettingsFieldInputMethod,
	})

	autoUpdateRow := components.RenderToggle(m.config.AutoUpdate, components.ToggleOptions{
		Label:      "Auto-update",
		LabelWidth: labelWidth,
		Focused:    m.focusedField == SettingsFieldAutoUpdate,
	})

	skipQuitConfirmRow := components.RenderToggle(m.config.SkipQuitConfirmation, components.ToggleOptions{
		Label:      "Skip quit confirm",
		LabelWidth: labelWidth,
		Focused:    m.focusedField == SettingsFieldSkipQuitConfirm,
	})

	// Context-aware hints
	var hints string
	if m.focusedField == SettingsFieldAutoUpdate || m.focusedField == SettingsFieldSkipQuitConfirm {
		// Toggle hints
		hints = components.RenderHintsStructured([]components.Hint{
			{Key: "↑↓", Action: "Navigate"},
			{Key: "Space", Action: "Toggle"},
			{Key: "Esc", Action: "Back"},
		})
	} else {
		// Selector hints
		hints = components.RenderHintsStructured([]components.Hint{
			{Key: "↑↓", Action: "Navigate"},
			{Key: "←→", Action: "Change"},
			{Key: "Esc", Action: "Back"},
		})
	}

	// Build settings block (left-aligned rows)
	settingsBlock := lipgloss.JoinVertical(lipgloss.Left,
		difficultyRow,
		durationRow,
		inputMethodRow,
		separator,
		autoUpdateRow,
		skipQuitConfirmRow,
	)

	// Build main content with centered title and settings block
	mainContent := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		settingsBlock,
	)

	// Bottom-anchored hints layout with small gap at bottom
	if m.width > 0 && m.height > 0 {
		hintsHeight := lipgloss.Height(hints)
		bottomPadding := 1
		availableHeight := m.height - hintsHeight - bottomPadding

		centeredMain := lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, mainContent)
		centeredHints := lipgloss.Place(m.width, hintsHeight+bottomPadding, lipgloss.Center, lipgloss.Top, hints)

		b.WriteString(lipgloss.JoinVertical(lipgloss.Left, centeredMain, centeredHints))
		return b.String()
	}

	// Fallback for unknown dimensions
	b.WriteString(lipgloss.JoinVertical(lipgloss.Center, mainContent, "", "", hints))
	return b.String()
}

// SetSize sets the screen dimensions.
func (m *SettingsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Config returns the current config.
func (m *SettingsModel) Config() *storage.Config {
	return m.config
}

// Helper functions (shared across screens package)

// maxLen returns the length of the longest string in the slice.
func maxLen(items []string) int {
	max := 0
	for _, s := range items {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

func findDifficultyIndex(name string) int {
	diffs := game.AllDifficulties()
	for i, d := range diffs {
		if d.String() == name {
			return i
		}
	}
	// Fallback: find the default difficulty
	for i, d := range diffs {
		if d.String() == storage.DefaultDifficulty {
			return i
		}
	}
	return 0
}

func findDurationIndexByMs(ms int64) int {
	for i, d := range modes.AllowedDurations {
		if d.Value.Milliseconds() == ms {
			return i
		}
	}
	return 1 // Default to 60s (index 1)
}

func settingsDifficultyNames(diffs []game.Difficulty) []string {
	names := make([]string, len(diffs))
	for i, d := range diffs {
		names[i] = d.String()
	}
	return names
}

func settingsDurationNames(durs []modes.Duration) []string {
	names := make([]string, len(durs))
	for i, d := range durs {
		names[i] = d.Label
	}
	return names
}
