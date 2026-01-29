package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 3: Implement pre-game launch screen
// - Mode settings display
// - Difficulty selection
// - Duration selection
// - Start game button

// LaunchModel represents the launch screen.
type LaunchModel struct {
	width  int
	height int
}

// NewLaunch creates a new launch model.
func NewLaunch() LaunchModel {
	return LaunchModel{}
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
		if msg.String() == "esc" {
			return m, func() tea.Msg {
				return ReturnToMenuMsg{}
			}
		}
	}

	return m, nil
}

// View renders the launch screen.
func (m LaunchModel) View() string {
	title := styles.Bold.Render("LAUNCH")
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
