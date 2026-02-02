package screens

import (
	"github.com/charmbracelet/bubbles/viewport"
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

// Layout constants for fixed sections
const (
	settingsHintsHeight = 3 // Height reserved for hints at the bottom
)

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
	viewport         viewport.Model
	viewportReady    bool
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
		viewport:         viewport.New(0, 0),
		viewportReady:    false,
	}
}

// Init initializes the settings model.
func (m SettingsModel) Init() tea.Cmd {
	return nil
}

// Update handles settings screen input.
func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.focusPrev()
			m.updateViewportContent()
		case "down", "j":
			m.focusNext()
			m.updateViewportContent()
		case "left", "h":
			m.adjustValue(-1)
			m.updateViewportContent()
		case "right", "l":
			m.adjustValue(1)
			m.updateViewportContent()
		case "enter", " ":
			if m.focusedField == SettingsFieldAutoUpdate {
				m.toggleAutoUpdate()
				m.updateViewportContent()
			} else if m.focusedField == SettingsFieldSkipQuitConfirm {
				m.toggleSkipQuitConfirm()
				m.updateViewportContent()
			}
		case "esc":
			return m, func() tea.Msg {
				return ReturnToMenuMsg{}
			}
		}
	}

	// Update viewport (for mouse scrolling if enabled)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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
	if !m.viewportReady {
		return "Loading..."
	}

	// Hints
	hints := m.getHints()

	// All screens: viewport + hints
	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		lipgloss.Place(m.width, settingsHintsHeight, lipgloss.Center, lipgloss.Center, hints),
	)
}

// getHints returns the context-aware hints for the settings screen.
func (m SettingsModel) getHints() string {
	if m.focusedField == SettingsFieldAutoUpdate || m.focusedField == SettingsFieldSkipQuitConfirm {
		// Toggle hints
		return components.RenderHintsStructured([]components.Hint{
			{Key: "↑↓", Action: "Navigate"},
			{Key: "Space", Action: "Toggle"},
			{Key: "Esc", Action: "Back"},
		})
	}
	// Selector hints
	return components.RenderHintsStructured([]components.Hint{
		{Key: "↑↓", Action: "Navigate"},
		{Key: "←→", Action: "Change"},
		{Key: "Esc", Action: "Back"},
	})
}

// SetSize sets the screen dimensions.
func (m *SettingsModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	viewportHeight := m.calculateViewportHeight()

	if !m.viewportReady {
		m.viewport = viewport.New(m.width, viewportHeight)
		m.viewport.YPosition = 0
		m.viewportReady = true
	} else {
		m.viewport.Width = m.width
		m.viewport.Height = viewportHeight
	}

	m.updateViewportContent()
}

// calculateViewportHeight returns the viewport height.
func (m SettingsModel) calculateViewportHeight() int {
	viewportHeight := m.height - settingsHintsHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	return viewportHeight
}

// updateViewportContent updates the viewport with the settings content.
func (m *SettingsModel) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	content := m.renderSettingsContent()
	m.viewport.SetContent(content)
}

// renderSettingsContent renders the main settings content for the viewport.
func (m SettingsModel) renderSettingsContent() string {
	// Title
	title := styles.Logo.Render("SETTINGS")

	// Gather all data
	diffs := game.AllDifficulties()
	durs := modes.AllowedDurations
	inputOptions := []string{"Typing", "Multiple Choice"}

	// All labels used in settings (for width calculation)
	labels := []string{"Difficulty", "Duration", "Input", "Auto-update", "Skip quit confirm"}

	// All possible values across all selectors
	allValues := []string{}
	allValues = append(allValues, settingsDifficultyNames(diffs)...)
	allValues = append(allValues, settingsDurationNames(durs)...)
	allValues = append(allValues, inputOptions...)

	// Calculate widths dynamically
	labelWidth := maxLen(labels)
	valueWidth := maxLen(allValues)

	// Focus indicator helper
	focusPrefix := func(focused bool) string {
		if focused {
			return styles.Accent.Render("> ")
		}
		return "  "
	}

	// Section headers
	gameDefaultsHeader := styles.Dim.Render("── Game Defaults ──")
	preferencesHeader := styles.Dim.Render("── Preferences ──")

	// Game defaults rows
	difficultyRow := focusPrefix(m.focusedField == SettingsFieldDifficulty) +
		components.RenderSelector(m.difficultyIndex, settingsDifficultyNames(diffs), components.SelectorOptions{
			Label:      "Difficulty",
			LabelWidth: labelWidth,
			ValueWidth: valueWidth,
			Focused:    m.focusedField == SettingsFieldDifficulty,
		})

	durationRow := focusPrefix(m.focusedField == SettingsFieldDuration) +
		components.RenderSelector(m.durationIndex, settingsDurationNames(durs), components.SelectorOptions{
			Label:      "Duration",
			LabelWidth: labelWidth,
			ValueWidth: valueWidth,
			Focused:    m.focusedField == SettingsFieldDuration,
		})

	inputMethodRow := focusPrefix(m.focusedField == SettingsFieldInputMethod) +
		components.RenderSelector(m.inputMethodIndex, inputOptions, components.SelectorOptions{
			Label:      "Input",
			LabelWidth: labelWidth,
			ValueWidth: valueWidth,
			Focused:    m.focusedField == SettingsFieldInputMethod,
		})

	// Preferences rows
	autoUpdateRow := focusPrefix(m.focusedField == SettingsFieldAutoUpdate) +
		components.RenderToggle(m.config.AutoUpdate, components.ToggleOptions{
			Label:      "Auto-update",
			LabelWidth: labelWidth,
			Focused:    m.focusedField == SettingsFieldAutoUpdate,
		})

	skipQuitConfirmRow := focusPrefix(m.focusedField == SettingsFieldSkipQuitConfirm) +
		components.RenderToggle(m.config.SkipQuitConfirmation, components.ToggleOptions{
			Label:      "Skip quit confirm",
			LabelWidth: labelWidth,
			Focused:    m.focusedField == SettingsFieldSkipQuitConfirm,
		})

	// Build settings block with section headers
	settingsBlock := lipgloss.JoinVertical(lipgloss.Left,
		gameDefaultsHeader,
		"",
		difficultyRow,
		durationRow,
		inputMethodRow,
		"",
		preferencesHeader,
		"",
		autoUpdateRow,
		skipQuitConfirmRow,
	)

	// Build main content with centered title and settings block
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		settingsBlock,
	)

	// Center both horizontally and vertically within viewport
	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
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
