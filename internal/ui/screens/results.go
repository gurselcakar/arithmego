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

// ResultsModel represents the results screen.
type ResultsModel struct {
	session     *game.Session
	saveError   error
	isFirstGame bool
	width       int
	height      int
}

// NewResults creates a new results model.
func NewResults(session *game.Session, saveError error) ResultsModel {
	return ResultsModel{
		session:     session,
		saveError:   saveError,
		isFirstGame: false,
	}
}

// NewResultsFirstGame creates a new results model for the first game (after onboarding).
func NewResultsFirstGame(session *game.Session, saveError error) ResultsModel {
	return ResultsModel{
		session:     session,
		saveError:   saveError,
		isFirstGame: true,
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

// ContinueToFeatureTourMsg is sent when continuing from first game results.
type ContinueToFeatureTourMsg struct{}

// Update handles results screen input.
func (m ResultsModel) Update(msg tea.Msg) (ResultsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.isFirstGame {
			// First game: only continue to feature tour
			switch msg.String() {
			case "enter", "right", "l":
				return m, func() tea.Msg {
					return ContinueToFeatureTourMsg{}
				}
			}
		} else {
			// Normal game: play again or menu
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
	}

	return m, nil
}

// View renders the results screen.
func (m ResultsModel) View() string {
	var b strings.Builder

	// Title
	title := styles.Bold.Render("RESULTS")

	// Score (prominent)
	score := components.RenderScore(m.session.Score)
	scoreLabel := styles.Dim.Render("points")

	// Separator
	separator := styles.Dim.Render("─────────────────────")

	// Stats line 1: correct count and accuracy
	correct := fmt.Sprintf("%d/%d correct", m.session.Correct, m.session.TotalAnswered())
	accuracy := fmt.Sprintf("%.0f%%", m.session.Accuracy())
	statsLine1 := correct + " · " + accuracy

	// Build detailed stats
	var statLines []string

	// Best streak (only show if > 0)
	if m.session.BestStreak > 0 {
		statLines = append(statLines, fmt.Sprintf("Best streak   %5d", m.session.BestStreak))
	}

	// Average response time
	avgTime := m.session.AvgResponseTime()
	if avgTime > 0 {
		statLines = append(statLines, fmt.Sprintf("Avg response  %5.2fs", avgTime.Seconds()))
	}

	// Fastest response time (only for correct answers)
	fastestTime := m.session.FastestResponseTime()
	if fastestTime > 0 {
		statLines = append(statLines, fmt.Sprintf("Fastest       %5.2fs", fastestTime.Seconds()))
	}

	// Skipped (only show if > 0)
	if m.session.Skipped > 0 {
		statLines = append(statLines, fmt.Sprintf("Skipped       %5d", m.session.Skipped))
	}

	// Intro message for first game (before feature tour)
	var introMessage string
	if m.isFirstGame {
		introMessage = styles.Tagline.Render("There's more to explore.")
	}

	// Hints based on game type
	var hints string
	if m.isFirstGame {
		hints = components.RenderHintsResponsive([]components.Hint{
			{Key: "→", Action: "Continue"},
		}, m.width)
	} else {
		hints = components.RenderHintsResponsive([]components.Hint{
			{Key: "M", Action: "Menu"},
			{Key: "↵", Action: "Play"},
		}, m.width)
	}

	// Save error warning (if any)
	var saveWarning string
	if m.saveError != nil {
		saveWarning = styles.Dim.Render("(Statistics could not be saved)")
	}

	// Build main content
	var contentParts []string
	contentParts = append(contentParts, title, "", "")
	contentParts = append(contentParts, score, scoreLabel, "")
	contentParts = append(contentParts, separator, "")
	contentParts = append(contentParts, statsLine1, "")

	// Add detailed stat lines
	for _, line := range statLines {
		contentParts = append(contentParts, line)
	}

	if saveWarning != "" {
		contentParts = append(contentParts, "", saveWarning)
	}

	// Add intro message for first game
	if introMessage != "" {
		contentParts = append(contentParts, "", separator, "", introMessage)
	}

	mainContent := lipgloss.JoinVertical(lipgloss.Center, contentParts...)

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
func (m *ResultsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
