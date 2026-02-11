package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 10: Add quit confirmation dialog before returning to menu

// PauseModel represents the pause screen.
type PauseModel struct {
	session *game.Session
	config  *storage.Config
	width   int
	height  int
}

// NewPause creates a new pause model.
func NewPause(session *game.Session, config *storage.Config) PauseModel {
	return PauseModel{
		session: session,
		config:  config,
	}
}

// Init initializes the pause model.
func (m PauseModel) Init() tea.Cmd {
	return nil
}

// ResumeMsg is sent when the user resumes the game.
type ResumeMsg struct {
	Session *game.Session
}

// QuitToMenuMsg is sent when the user quits to the menu.
type QuitToMenuMsg struct{}

// Update handles pause screen input.
func (m PauseModel) Update(msg tea.Msg) (PauseModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg {
				return ResumeMsg{Session: m.session}
			}
		case "q":
			return m, func() tea.Msg {
				return QuitConfirmMsg{Session: m.session}
			}
		}
	}

	return m, nil
}

// View renders the pause screen.
func (m PauseModel) View() string {
	// Title
	title := styles.Bold.Render("PAUSED")

	// Time remaining
	timer := components.FormatTimer(m.session.TimeLeft)

	// Hints
	hints := components.RenderHintsResponsive([]components.Hint{
		{Key: "Q", Action: "Quit"},
		{Key: "Enter", Action: "Resume"},
	}, m.width)

	// Main content (without hints)
	mainContent := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		timer,
	)

	// Bottom-anchored hints layout with small gap at bottom
	if m.width > 0 && m.height > 0 {
		hintsHeight := lipgloss.Height(hints)
		bottomPadding := 1
		availableHeight := m.height - hintsHeight - bottomPadding

		centeredMain := lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, mainContent)
		centeredHints := lipgloss.Place(m.width, hintsHeight+bottomPadding, lipgloss.Center, lipgloss.Top, hints)

		return lipgloss.JoinVertical(lipgloss.Left, centeredMain, centeredHints)
	}

	// Fallback for unknown dimensions
	return lipgloss.JoinVertical(lipgloss.Center, mainContent, "", "", hints)
}

// Session returns the current session.
func (m PauseModel) Session() *game.Session {
	return m.session
}

// SetSession updates the pause session.
func (m *PauseModel) SetSession(session *game.Session) {
	m.session = session
}

// SetSize sets the screen dimensions.
func (m *PauseModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
