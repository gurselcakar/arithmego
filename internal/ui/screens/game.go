package screens

import (
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Phase 4: Add scoreboard component (top left)
// - Points display with animation
// - Streak counter
// - Brief green/red flash on answer

// GameModel represents the gameplay screen.
type GameModel struct {
	session        *game.Session
	input          components.InputModel
	width          int
	height         int
	feedback       string    // "correct", "incorrect", or ""
	feedbackExpiry time.Time // when feedback should clear
}

// NewGame creates a new game model with the given session.
func NewGame(session *game.Session) GameModel {
	return GameModel{
		session: session,
		input:   components.NewInput(),
	}
}

// Init initializes the game and starts the timer.
func (m GameModel) Init() tea.Cmd {
	m.session.Start()
	return tea.Batch(
		m.input.Init(),
		TickCmd(),
	)
}

// tickMsg is sent every second to update the timer.
type tickMsg time.Time

// TickCmd returns a command that sends tick messages every second.
func TickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// GameOverMsg is sent when the game session ends.
type GameOverMsg struct {
	Session *game.Session
}

// PauseMsg is sent when the user pauses the game.
type PauseMsg struct {
	Session *game.Session
}

// Update handles game input and timer ticks.
func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.session.Tick()
		if m.session.IsFinished() {
			return m, func() tea.Msg {
				return GameOverMsg{Session: m.session}
			}
		}
		// Clear feedback if expired
		if m.feedback != "" && time.Now().After(m.feedbackExpiry) {
			m.feedback = ""
		}
		return m, TickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.submitAnswer()
		case "s", " ":
			return m.skipQuestion()
		case "p", "esc":
			return m, func() tea.Msg {
				return PauseMsg{Session: m.session}
			}
		default:
			// Pass to input
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// submitAnswer checks the current answer and moves to the next question.
func (m GameModel) submitAnswer() (GameModel, tea.Cmd) {
	val := m.input.Value()
	if val == "" {
		return m, nil
	}

	answer, err := strconv.Atoi(val)
	if err != nil {
		return m, nil
	}

	correct := m.session.SubmitAnswer(answer)
	m.input.Reset()

	if correct {
		m.feedback = "correct"
	} else {
		m.feedback = "incorrect"
	}
	m.feedbackExpiry = time.Now().Add(2 * time.Second)

	// Check if game is over
	if m.session.IsFinished() {
		return m, func() tea.Msg {
			return GameOverMsg{Session: m.session}
		}
	}

	return m, nil
}

// skipQuestion skips the current question.
func (m GameModel) skipQuestion() (GameModel, tea.Cmd) {
	m.session.Skip()
	m.input.Reset()
	m.feedback = ""
	m.feedbackExpiry = time.Time{}
	return m, nil
}

// View renders the game screen.
func (m GameModel) View() string {
	var b strings.Builder

	// Defensive check for nil session
	if m.session == nil {
		return "Loading..."
	}

	// Timer (top right)
	timer := components.FormatTimer(m.session.TimeLeft)

	// Question (center)
	var question string
	if m.session.Current != nil {
		question = components.RenderQuestion(m.session.Current.Display)
	}

	// Apply feedback styling to input area
	inputView := m.input.View()
	switch m.feedback {
	case "correct":
		inputView = styles.Correct.Render(inputView)
	case "incorrect":
		inputView = styles.Incorrect.Render(inputView)
	}

	// Hints
	hints := components.RenderHints([]string{"[S] Skip", "[P] Pause"})

	// Layout
	// Top row with timer on right
	timerLine := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Right).Render(timer)

	// Center content
	centerContent := lipgloss.JoinVertical(lipgloss.Center,
		question,
		"",
		inputView,
	)

	// Combine
	content := lipgloss.JoinVertical(lipgloss.Left,
		timerLine,
		"",
		"",
	)

	// Place centered content in middle
	if m.width > 0 && m.height > 0 {
		// Calculate available height after timer
		availHeight := m.height - 4 // timer line + spacing + hints

		centeredQuestion := lipgloss.Place(m.width, availHeight-4, lipgloss.Center, lipgloss.Center, centerContent)
		hintsLine := lipgloss.Place(m.width, 1, lipgloss.Center, lipgloss.Bottom, hints)

		content = lipgloss.JoinVertical(lipgloss.Left,
			timerLine,
			centeredQuestion,
			hintsLine,
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center,
			timerLine,
			"",
			"",
			centerContent,
			"",
			"",
			hints,
		)
	}

	b.WriteString(content)
	return b.String()
}

// Session returns the current game session.
func (m GameModel) Session() *game.Session {
	return m.session
}

// SetSession updates the game session (used when resuming from pause).
func (m *GameModel) SetSession(session *game.Session) {
	m.session = session
}

// SetSize sets the screen dimensions.
func (m *GameModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
