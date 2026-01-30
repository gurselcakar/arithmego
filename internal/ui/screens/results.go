package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 5: Add detailed breakdown
// - Per-operation accuracy
// - Response time stats
// - Save to statistics storage

// ResultsModel represents the results screen.
type ResultsModel struct {
	session *game.Session
	width   int
	height  int
}

// NewResults creates a new results model.
func NewResults(session *game.Session) ResultsModel {
	return ResultsModel{
		session: session,
	}
}

// Init initializes the results model.
func (m ResultsModel) Init() tea.Cmd {
	return nil
}

// PlayAgainMsg is sent when the user wants to play again.
type PlayAgainMsg struct{}

// ReturnToMenuMsg is sent when the user returns to the menu.
type ReturnToMenuMsg struct{}

// Update handles results screen input.
func (m ResultsModel) Update(msg tea.Msg) (ResultsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg {
				return PlayAgainMsg{}
			}
		case "m", "esc":
			return m, func() tea.Msg {
				return ReturnToMenuMsg{}
			}
		}
	}

	return m, nil
}

// View renders the results screen.
func (m ResultsModel) View() string {
	var b strings.Builder

	// Title
	title := styles.Bold.Render("SESSION COMPLETE")

	// Score (prominent)
	score := components.RenderScore(m.session.Score)
	scoreLabel := styles.Dim.Render("points")

	// Separator
	separator := styles.Dim.Render("────────────────────")

	// Stats
	correct := fmt.Sprintf("%d/%d correct", m.session.Correct, m.session.TotalAnswered())
	accuracy := fmt.Sprintf("%.0f%% accuracy", m.session.Accuracy())

	// Best streak (only show if > 0)
	var bestStreak string
	if m.session.BestStreak > 0 {
		bestStreak = fmt.Sprintf("Best streak: %d", m.session.BestStreak)
	}

	// Hints
	hints := components.RenderHints([]string{"[Enter] Play again", "[M] Main menu"})

	// Combine
	var contentParts []string
	contentParts = append(contentParts, title, "", "")
	contentParts = append(contentParts, score, scoreLabel, "")
	contentParts = append(contentParts, separator, "")
	contentParts = append(contentParts, correct, accuracy)
	if bestStreak != "" {
		contentParts = append(contentParts, bestStreak)
	}
	contentParts = append(contentParts, "", "", hints)

	content := lipgloss.JoinVertical(lipgloss.Center, contentParts...)

	// Center in terminal
	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	b.WriteString(content)
	return b.String()
}

// SetSize sets the screen dimensions.
func (m *ResultsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
