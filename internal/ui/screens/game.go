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

// Display durations
const (
	deltaDisplayTime  = 1 * time.Second
	milestoneShowTime = 2 * time.Second
)

// Score animation constants
const (
	// scoreAnimInterval: 30ms = ~33 FPS, provides smooth visual updates without excessive CPU use
	scoreAnimInterval = 30 * time.Millisecond
	// scoreAnimEasing: 30% per tick creates natural deceleration (fast start, slow finish)
	// Results in ~10 ticks (~300ms) to reach target from typical score changes
	scoreAnimEasing = 0.3
	// scoreAnimMinStep: prevents infinitely slow final ticks when remaining distance < 17
	// (since 0.3 * 16 = 4.8 rounds to 4, then 3, 2, 1... taking many ticks)
	scoreAnimMinStep = 5
)

// GameModel represents the gameplay screen.
type GameModel struct {
	session *game.Session
	input   components.InputModel
	width   int
	height  int

	// Input method
	inputMethod components.InputMethod
	choices     components.ChoicesModel

	// Scoring display state
	tick            int       // for shimmer animation (increments each second)
	scoreDelta      int       // points from last answer (for delta popup)
	deltaExpiry     time.Time // when delta popup should clear
	milestone       string    // milestone text (e.g., "×1.25", "×2.0 MAX")
	milestoneExpiry time.Time // when milestone should clear

	// Score animation state
	displayScore int  // currently displayed score (animates toward actual)
	animating    bool // whether score animation is in progress
}

// NewGame creates a new game model with the given session and input method.
// Note: displayScore starts at 0 to match session's initial score.
// The rendering logic handles any sync issues via fallback to session.Score.
func NewGame(session *game.Session, inputMethod components.InputMethod) GameModel {
	return GameModel{
		session:      session,
		input:        components.NewInput(),
		inputMethod:  inputMethod,
		choices:      components.NewChoices(),
		displayScore: 0, // Explicit: matches session.Score after Start()
	}
}

// gameStartMsg is sent to trigger initial setup that requires model mutation.
type gameStartMsg struct{}

// Init initializes the game and starts the timer.
func (m GameModel) Init() tea.Cmd {
	m.session.Start()

	return tea.Batch(
		m.input.Init(),
		TickCmd(),
		func() tea.Msg { return gameStartMsg{} }, // Trigger choice generation in Update
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

// scoreAnimMsg is sent rapidly during score animation.
type scoreAnimMsg time.Time

// ScoreAnimCmd returns a command for the fast score animation tick.
func ScoreAnimCmd() tea.Cmd {
	return tea.Tick(scoreAnimInterval, func(t time.Time) tea.Msg {
		return scoreAnimMsg(t)
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

// QuitConfirmMsg is sent when the user wants to quit from the game.
type QuitConfirmMsg struct {
	Session *game.Session
}

// Update handles game input and timer ticks.
func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	switch msg := msg.(type) {
	case gameStartMsg:
		// Generate initial choices for multiple choice mode
		if m.inputMethod == components.InputMultipleChoice && m.session.Current != nil {
			choices, correctIndex := game.GenerateChoices(m.session.Current.Answer, m.session.Difficulty)
			m.choices.SetChoices(choices, correctIndex)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case components.ChoiceSelectedMsg:
		// Auto-submit on choice selection (multiple choice mode)
		return m.submitAnswerValue(msg.Value)

	case tickMsg:
		m.session.Tick()
		m.tick++ // increment for shimmer animation

		if m.session.IsFinished() {
			// Stop animation and sync display score before transitioning.
			// Not returning TickCmd() stops the timer tick loop.
			m.animating = false
			m.displayScore = m.session.Score
			return m, func() tea.Msg {
				return GameOverMsg{Session: m.session}
			}
		}

		// Clear expired display states
		now := time.Now()
		if m.scoreDelta != 0 && now.After(m.deltaExpiry) {
			m.scoreDelta = 0
		}
		if m.milestone != "" && now.After(m.milestoneExpiry) {
			m.milestone = ""
		}

		return m, TickCmd()

	case scoreAnimMsg:
		if !m.animating {
			return m, nil
		}

		target := m.session.Score
		diff := target - m.displayScore

		if diff == 0 {
			m.animating = false
			return m, nil
		}

		// Calculate step with easing (move percentage of remaining distance)
		step := int(float64(diff) * scoreAnimEasing)

		// Ensure minimum step to avoid slow crawl at the end
		if step == 0 {
			if diff > 0 {
				step = min(scoreAnimMinStep, diff)
			} else {
				step = max(-scoreAnimMinStep, diff)
			}
		}

		m.displayScore += step

		// Check if we've reached target
		if (diff > 0 && m.displayScore >= target) || (diff < 0 && m.displayScore <= target) {
			m.displayScore = target
			m.animating = false
			return m, nil
		}

		return m, ScoreAnimCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.inputMethod == components.InputTyping {
				return m.submitAnswer()
			}
			return m, nil
		case "s", " ":
			return m.skipQuestion()
		case "p":
			return m, func() tea.Msg {
				return PauseMsg{Session: m.session}
			}
		case "q", "esc":
			return m, func() tea.Msg {
				return QuitConfirmMsg{Session: m.session}
			}
		default:
			// Route to appropriate input component
			if m.inputMethod == components.InputMultipleChoice {
				var cmd tea.Cmd
				m.choices, cmd = m.choices.Update(msg)
				return m, cmd
			}
			// Pass to text input
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
		return m, nil // User must type something; "0" for zero answers
	}

	answer, err := strconv.Atoi(val)
	if err != nil {
		return m, nil // Safety net for edge cases like "-" only
	}

	return m.submitAnswerValue(answer)
}

// submitAnswerValue submits an answer and handles animation.
// Used by both typing mode (via submitAnswer) and multiple choice mode.
func (m GameModel) submitAnswerValue(answer int) (GameModel, tea.Cmd) {
	_ = m.session.SubmitAnswer(answer)

	// Reset input components
	m.input.Reset()
	m.choices.Reset()

	// Set score delta for display and start animation
	var cmd tea.Cmd
	if m.session.LastResult != nil {
		m.scoreDelta = m.session.LastResult.Points
		m.deltaExpiry = time.Now().Add(deltaDisplayTime)

		// Start score animation
		if m.session.Score != m.displayScore {
			m.animating = true
			cmd = ScoreAnimCmd()
		}

		// Check for milestone
		if m.session.LastResult.IsMilestone {
			m.milestone = game.GetMilestoneAnnouncement(m.session.LastResult.NewStreak)
			m.milestoneExpiry = time.Now().Add(milestoneShowTime)
		}
	}

	// Check if game is over
	if m.session.IsFinished() {
		m.animating = false
		m.displayScore = m.session.Score
		return m, func() tea.Msg {
			return GameOverMsg{Session: m.session}
		}
	}

	// Generate new choices for the next question
	if m.inputMethod == components.InputMultipleChoice && m.session.Current != nil {
		choices, correctIndex := game.GenerateChoices(m.session.Current.Answer, m.session.Difficulty)
		m.choices.SetChoices(choices, correctIndex)
	}

	return m, cmd
}

// skipQuestion skips the current question.
func (m GameModel) skipQuestion() (GameModel, tea.Cmd) {
	m.session.Skip()
	m.input.Reset()
	m.choices.Reset()
	m.scoreDelta = 0
	m.deltaExpiry = time.Time{}
	// Sync animation state (skip doesn't change score, but stop any in-progress animation)
	m.animating = false
	m.displayScore = m.session.Score

	// Generate new choices for the next question
	if m.inputMethod == components.InputMultipleChoice && m.session.Current != nil {
		choices, correctIndex := game.GenerateChoices(m.session.Current.Answer, m.session.Difficulty)
		m.choices.SetChoices(choices, correctIndex)
	}

	return m, nil
}

// View renders the game screen.
func (m GameModel) View() string {
	var b strings.Builder

	// Defensive check for nil session
	if m.session == nil {
		return "Loading..."
	}

	// Build top row: scoreboard (left) | score+delta (center) | timer (right)
	topRow := m.renderTopRow()

	// Question (center)
	var question string
	if m.session.Current != nil {
		if m.inputMethod == components.InputMultipleChoice {
			question = components.RenderQuestionWithAnswer(m.session.Current.Display)
		} else {
			question = components.RenderQuestion(m.session.Current.Display)
		}
	}

	// Input area
	var inputView string
	if m.inputMethod == components.InputMultipleChoice {
		inputView = m.choices.View()
	} else {
		inputView = m.input.View()
	}

	// Hints - differ based on input method
	var hints string
	if m.inputMethod == components.InputMultipleChoice {
		hints = components.RenderHintsResponsive([]components.Hint{
			{Key: "1-4", Action: "Select"},
			{Key: "S", Action: "Skip"},
			{Key: "P", Action: "Pause"},
			{Key: "Esc", Action: "Quit"},
		}, m.width)
	} else {
		hints = components.RenderHintsResponsive([]components.Hint{
			{Key: "S", Action: "Skip"},
			{Key: "P", Action: "Pause"},
			{Key: "Esc", Action: "Quit"},
		}, m.width)
	}

	// Center content (milestone is now shown above score in top row)
	centerContent := lipgloss.JoinVertical(lipgloss.Center,
		question,
		"",
		inputView,
	)

	// Layout with bottom-anchored hints and small gap at bottom
	if m.width > 0 && m.height > 0 {
		hintsHeight := lipgloss.Height(hints)
		bottomPadding := 1
		topPadding := 2
		topRowHeight := lipgloss.Height(topRow)
		availableHeight := m.height - topPadding - topRowHeight - hintsHeight - bottomPadding

		topRow = strings.Repeat("\n", topPadding) + topRow

		centeredQuestion := lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, centerContent)
		centeredHints := lipgloss.Place(m.width, hintsHeight+bottomPadding, lipgloss.Center, lipgloss.Top, hints)

		content := lipgloss.JoinVertical(lipgloss.Left,
			topRow,
			centeredQuestion,
			centeredHints,
		)
		b.WriteString(content)
	} else {
		content := lipgloss.JoinVertical(lipgloss.Center,
			topRow,
			"",
			"",
			centerContent,
			"",
			"",
			hints,
		)
		b.WriteString(content)
	}

	return b.String()
}

// renderTopRow renders the top status bar with scoreboard, score, and timer.
func (m GameModel) renderTopRow() string {
	if m.width < 40 {
		// Fallback for unknown or very narrow width
		return m.renderTopRowSimple()
	}

	// Horizontal margin to keep HUD content away from terminal edges
	const margin = 4

	// Left: Scoreboard (multiplier + streak bar)
	scoreboard := components.RenderScoreboard(m.session.Streak, m.tick)

	// Center: Score with delta popup
	scoreDisplay := m.renderScoreWithDelta()

	// Right: Timer with "remaining" label
	timer := components.FormatTimer(m.session.TimeLeft)
	timerWithLabel := lipgloss.JoinVertical(lipgloss.Right,
		styles.Dim.Render("Remaining"),
		timer,
	)

	// Calculate column widths from usable area (excluding margins)
	hudWidth := m.width - 2*margin
	leftWidth := hudWidth / 4
	rightWidth := hudWidth / 4
	centerWidth := hudWidth - leftWidth - rightWidth

	// Style each column
	leftCol := lipgloss.NewStyle().Width(leftWidth).Align(lipgloss.Left).Render(scoreboard)
	centerCol := lipgloss.NewStyle().Width(centerWidth).Align(lipgloss.Center).Render(scoreDisplay)
	rightCol := lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Right).Render(timerWithLabel)

	row := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, centerCol, rightCol)
	return lipgloss.NewStyle().PaddingLeft(margin).Render(row)
}

// renderTopRowSimple renders a simpler top row when width is unknown.
func (m GameModel) renderTopRowSimple() string {
	scoreboard := components.RenderScoreboard(m.session.Streak, m.tick)
	// Use displayScore during animation, actual score otherwise
	scoreValue := m.displayScore
	if !m.animating {
		scoreValue = m.session.Score
	}
	score := components.RenderScore(scoreValue)
	timer := components.FormatTimer(m.session.TimeLeft)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		scoreboard,
		"    ",
		score,
		"    ",
		timer,
	)
}

// renderScoreWithDelta renders the score with label, delta popup, and score number.
func (m GameModel) renderScoreWithDelta() string {
	// During animation, show the animating displayScore.
	// When not animating, displayScore should equal session.Score.
	// The fallback to session.Score handles edge cases (e.g., if animation
	// state gets out of sync due to rapid state changes).
	scoreValue := m.displayScore
	if !m.animating {
		scoreValue = m.session.Score
	}
	score := components.RenderScoreLarge(scoreValue)

	// Build the display from top to bottom: label/milestone, score, delta
	var parts []string

	// Top line: milestone replaces "Score" label when active
	if m.milestone != "" {
		parts = append(parts, styles.Milestone.Render(m.milestone))
	} else {
		parts = append(parts, styles.Dim.Render("Score"))
	}

	// Score number
	parts = append(parts, score)

	// Delta popup (+150 or -25)
	if m.scoreDelta != 0 {
		parts = append(parts, components.RenderScoreDelta(m.scoreDelta))
	} else {
		parts = append(parts, "") // empty line for consistent height
	}

	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

// Session returns the current game session.
func (m GameModel) Session() *game.Session {
	return m.session
}

// SetSession updates the game session (used when resuming from pause).
// Panics if session is nil.
func (m *GameModel) SetSession(session *game.Session) {
	if session == nil {
		panic("SetSession: session cannot be nil")
	}
	m.session = session
}

// SetSize sets the screen dimensions.
func (m *GameModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
