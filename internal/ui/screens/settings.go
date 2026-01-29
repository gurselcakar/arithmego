package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 8: Implement user preferences
// - Input method selection
// - Default difficulty
// - Default duration
// - Display options

// SettingsModel represents the settings screen.
type SettingsModel struct {
	width  int
	height int
}

// NewSettings creates a new settings model.
func NewSettings() SettingsModel {
	return SettingsModel{}
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
		if msg.String() == "esc" {
			return m, func() tea.Msg {
				return ReturnToMenuMsg{}
			}
		}
	}

	return m, nil
}

// View renders the settings screen.
func (m SettingsModel) View() string {
	title := styles.Bold.Render("SETTINGS")
	body := "Coming soon."
	hints := components.RenderHints([]string{"Esc Back"})

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		body,
		"",
		hints,
	)

	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
