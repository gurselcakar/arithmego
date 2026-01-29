package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 5: Implement performance data display
// - Session history
// - Per-operation accuracy
// - Best streaks

// StatisticsModel represents the statistics screen.
type StatisticsModel struct {
	width  int
	height int
}

// NewStatistics creates a new statistics model.
func NewStatistics() StatisticsModel {
	return StatisticsModel{}
}

// Init initializes the statistics model.
func (m StatisticsModel) Init() tea.Cmd {
	return nil
}

// Update handles statistics screen input.
func (m StatisticsModel) Update(msg tea.Msg) (StatisticsModel, tea.Cmd) {
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

// View renders the statistics screen.
func (m StatisticsModel) View() string {
	title := styles.Bold.Render("STATISTICS")
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
