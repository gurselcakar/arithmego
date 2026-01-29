package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 3: Implement mode selection and customization
// - List of available modes
// - Mode preview/description
// - Custom mode creation

// ModesModel represents the modes selection screen.
type ModesModel struct {
	width  int
	height int
}

// NewModes creates a new modes model.
func NewModes() ModesModel {
	return ModesModel{}
}

// Init initializes the modes model.
func (m ModesModel) Init() tea.Cmd {
	return nil
}

// Update handles modes screen input.
func (m ModesModel) Update(msg tea.Msg) (ModesModel, tea.Cmd) {
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

// View renders the modes screen.
func (m ModesModel) View() string {
	title := styles.Bold.Render("MODES")
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
