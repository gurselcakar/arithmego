package screens

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// FeatureTourStep represents the current step in the feature tour.
type FeatureTourStep int

const (
	StepPractice FeatureTourStep = iota
	StepModes
	StepStatistics
	StepFinale
)

const featureTourTotalSteps = 3

// FeatureTourCompleteMsg is sent when the feature tour is completed or skipped.
type FeatureTourCompleteMsg struct{}

// FeatureTourModel represents the feature tour screen.
type FeatureTourModel struct {
	step          FeatureTourStep
	width         int
	height        int
	viewport      viewport.Model
	viewportReady bool
}

// NewFeatureTour creates a new feature tour model.
func NewFeatureTour() FeatureTourModel {
	return FeatureTourModel{
		step:          StepPractice,
		viewport:      viewport.New(0, 0),
		viewportReady: false,
	}
}

// Init initializes the feature tour model.
func (m FeatureTourModel) Init() tea.Cmd {
	return nil
}

// Update handles feature tour screen input.
func (m FeatureTourModel) Update(msg tea.Msg) (FeatureTourModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "right", "l":
			return m.advance()
		case "s", "S":
			// Skip not available on statistics or finale
			if m.step != StepStatistics && m.step != StepFinale {
				return m.skip()
			}
		case "b", "B", "left", "h":
			// Back not available on finale
			if m.step != StepFinale {
				m.back()
				m.updateViewportContent()
			}
		}
	}

	// Update viewport (for scrolling support)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// advance moves to the next step or completes the tour.
func (m FeatureTourModel) advance() (FeatureTourModel, tea.Cmd) {
	switch m.step {
	case StepPractice:
		m.step = StepModes
		m.updateViewportContent()
		return m, nil
	case StepModes:
		m.step = StepStatistics
		m.updateViewportContent()
		return m, nil
	case StepStatistics:
		m.step = StepFinale
		m.updateViewportContent()
		return m, nil
	case StepFinale:
		return m, m.complete()
	}
	return m, nil
}

// back returns to the previous step.
func (m *FeatureTourModel) back() {
	switch m.step {
	case StepModes:
		m.step = StepPractice
	case StepStatistics:
		m.step = StepModes
	case StepFinale:
		m.step = StepStatistics
	}
	// StepPractice has no back
}

// skip jumps to the finale screen (skip always lands on finale before menu).
func (m FeatureTourModel) skip() (FeatureTourModel, tea.Cmd) {
	m.step = StepFinale
	m.updateViewportContent()
	return m, nil
}

// complete returns the completion message.
func (m FeatureTourModel) complete() tea.Cmd {
	return func() tea.Msg {
		return FeatureTourCompleteMsg{}
	}
}

// View renders the feature tour screen.
func (m FeatureTourModel) View() string {
	if !m.viewportReady {
		return "Loading..."
	}

	progress := m.getProgressForStep()
	hints := m.getHintsForStep()

	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		lipgloss.Place(m.width, progressHeight, lipgloss.Center, lipgloss.Center, progress),
		lipgloss.Place(m.width, components.HintsHeight, lipgloss.Center, lipgloss.Center, hints),
	)
}

// getProgressForStep returns the progress dots for the current step.
func (m FeatureTourModel) getProgressForStep() string {
	switch m.step {
	case StepPractice:
		return components.ProgressDotsColored(1, featureTourTotalSteps)
	case StepModes:
		return components.ProgressDotsColored(2, featureTourTotalSteps)
	case StepStatistics:
		return components.ProgressDotsColored(3, featureTourTotalSteps)
	case StepFinale:
		// No progress dots on finale screen
		return ""
	default:
		return ""
	}
}

// getHintsForStep returns the appropriate hints for the current step.
func (m FeatureTourModel) getHintsForStep() string {
	switch m.step {
	case StepPractice:
		// First step: no back button
		return components.RenderHintsResponsive([]components.Hint{
			{Key: "S", Action: "Skip"},
			{Key: "→", Action: "Continue"},
		}, m.width)
	case StepStatistics:
		// Last feature step: no skip (already at the end)
		return components.RenderHintsResponsive([]components.Hint{
			{Key: "←", Action: "Back"},
			{Key: "→", Action: "Continue"},
		}, m.width)
	case StepFinale:
		// Finale: only "Let's go" - no back or skip
		return components.RenderHintsResponsive([]components.Hint{
			{Key: "→", Action: "Let's go"},
		}, m.width)
	default:
		// Middle steps: back, skip, continue
		return components.RenderHintsResponsive([]components.Hint{
			{Key: "←", Action: "Back"},
			{Key: "S", Action: "Skip"},
			{Key: "→", Action: "Continue"},
		}, m.width)
	}
}

// SetSize updates the screen dimensions.
func (m *FeatureTourModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	viewportHeight := m.calculateViewportHeight()

	components.SetViewportSize(&m.viewport, &m.viewportReady, m.width, viewportHeight)

	m.updateViewportContent()
}

// calculateViewportHeight returns the viewport height.
func (m FeatureTourModel) calculateViewportHeight() int {
	bottomSectionHeight := components.HintsHeight + progressHeight

	viewportHeight := m.height - bottomSectionHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}

	return viewportHeight
}

// updateViewportContent updates the viewport with the current step's content.
func (m *FeatureTourModel) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	content := m.getViewportContent()
	m.viewport.SetContent(content)
}

// getViewportContent returns the content for the current step.
func (m FeatureTourModel) getViewportContent() string {
	switch m.step {
	case StepPractice:
		return m.renderPracticeContent()
	case StepModes:
		return m.renderModesContent()
	case StepStatistics:
		return m.renderStatisticsContent()
	case StepFinale:
		return m.renderFinaleContent()
	default:
		return ""
	}
}

// renderPracticeContent renders the practice mode introduction.
func (m FeatureTourModel) renderPracticeContent() string {
	title := styles.Logo.Render("PRACTICE MODE")

	description := lipgloss.JoinVertical(lipgloss.Center,
		styles.Subtle.Render("Want to improve without the clock ticking?"),
		"",
		"Practice lets you focus on specific",
		"operations at your own pace.",
	)

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		"",
		description,
	)

	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// renderModesContent renders the game modes introduction with grouped list.
func (m FeatureTourModel) renderModesContent() string {
	title := styles.Logo.Render("GAME MODES")

	subtitle := styles.Subtle.Render("Here's what you can explore:")

	// Build the grouped modes box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(0, 2)

	// Sprint modes (single operations)
	sprintLabel := styles.Bold.Render("SPRINT")
	sprintModes := styles.Dim.Render("+  −  ×  ÷  x²  x³  √x  ³√x  xⁿ  mod  %  n!")

	// Challenge modes (mixed operations)
	challengeLabel := styles.Bold.Render("CHALLENGE")
	challengeModes := styles.Dim.Render("Mixed Basics · Mixed Powers · Mixed Advanced · Anything Goes")

	modesContent := lipgloss.JoinVertical(lipgloss.Left,
		sprintLabel,
		sprintModes,
		"",
		challengeLabel,
		challengeModes,
	)

	modesBox := boxStyle.Render(modesContent)

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		subtitle,
		"",
		"",
		modesBox,
	)

	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// renderStatisticsContent renders the statistics introduction.
func (m FeatureTourModel) renderStatisticsContent() string {
	title := styles.Logo.Render("STATISTICS")

	description := lipgloss.JoinVertical(lipgloss.Center,
		styles.Subtle.Render("Track your progress over time."),
		"",
		"See your accuracy, best scores, and",
		"how you're improving at each mode.",
	)

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		"",
		description,
	)

	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// renderFinaleContent renders the finale screen with logo, tagline, and personal note.
func (m FeatureTourModel) renderFinaleContent() string {
	logo := components.LogoColoredForWidth(m.width)
	logoSeparator := styles.Dim.Render(components.LogoSeparator())
	tagline := components.Tagline()

	// Personal note
	noteSeparator := styles.Dim.Render("─────────────────")
	quote := styles.Tagline.Render("\"In a world of instant answers, take a moment to think.\"")
	byLine := styles.Tagline.Render("by @gurselcakar")
	claudeLine := styles.Tagline.Render("+ Claude Code")

	content := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		logoSeparator,
		"",
		tagline,
		"",
		"",
		noteSeparator,
		"",
		quote,
		"",
		byLine,
		claudeLine,
	)

	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}
