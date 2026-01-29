package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 8: Implement first-time user flow
// - Welcome screen
// - Input method selection
// - Session length selection
// - First mode selection

// OnboardingModel represents the onboarding screen.
type OnboardingModel struct {
	width  int
	height int
}

// NewOnboarding creates a new onboarding model.
func NewOnboarding() OnboardingModel {
	return OnboardingModel{}
}

// Init initializes the onboarding model.
func (m OnboardingModel) Init() tea.Cmd {
	return nil
}

// Update handles onboarding screen input.
func (m OnboardingModel) Update(msg tea.Msg) (OnboardingModel, tea.Cmd) {
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

// View renders the onboarding screen.
func (m OnboardingModel) View() string {
	title := styles.Bold.Render("WELCOME")
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
