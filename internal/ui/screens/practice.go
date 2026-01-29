package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 7: Implement sandbox mode with live controls
// - No timer, no score
// - Live difficulty adjustment
// - Operation switching

// PracticeModel represents the practice screen.
type PracticeModel struct {
	width  int
	height int
}

// NewPractice creates a new practice model.
func NewPractice() PracticeModel {
	return PracticeModel{}
}

// Init initializes the practice model.
func (m PracticeModel) Init() tea.Cmd {
	return nil
}

// Update handles practice screen input.
func (m PracticeModel) Update(msg tea.Msg) (PracticeModel, tea.Cmd) {
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

// View renders the practice screen.
func (m PracticeModel) View() string {
	title := styles.Bold.Render("PRACTICE")
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
