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
	StepModes FeatureTourStep = iota
	StepPractice
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
		step:          StepModes,
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
		case "ctrl+c":
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
	case StepModes:
		m.step = StepPractice
		m.updateViewportContent()
		return m, nil
	case StepPractice:
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
	case StepPractice:
		m.step = StepModes
	case StepStatistics:
		m.step = StepPractice
	}
	// StepModes has no back (first step)
	// StepFinale back is blocked in Update()
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
	case StepModes:
		return components.ProgressDotsColored(1, featureTourTotalSteps)
	case StepPractice:
		return components.ProgressDotsColored(2, featureTourTotalSteps)
	case StepStatistics:
		return components.ProgressDotsColored(3, featureTourTotalSteps)
	case StepFinale:
		return ""
	default:
		return ""
	}
}

// getHintsForStep returns the appropriate hints for the current step.
func (m FeatureTourModel) getHintsForStep() string {
	switch m.step {
	case StepModes:
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
	case StepModes:
		return m.renderModesContent()
	case StepPractice:
		return m.renderPracticeContent()
	case StepStatistics:
		return m.renderStatisticsContent()
	case StepFinale:
		return m.renderFinaleContent()
	default:
		return ""
	}
}

// renderModesContent renders the game modes introduction with grouped list.
func (m FeatureTourModel) renderModesContent() string {
	title := styles.Logo.Render("GAME MODES")

	subtitle := styles.Subtle.Render("16 modes. Two categories. Pick your challenge.")

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
		styled := lipgloss.NewStyle().MarginTop(titleTopPadding(m.viewport.Height)).Render(content)
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Top, styled)
	}
	return content
}

// renderPracticeContent renders the practice mode introduction.
func (m FeatureTourModel) renderPracticeContent() string {
	title := styles.Logo.Render("PRACTICE MODE")

	description := lipgloss.JoinVertical(lipgloss.Center,
		styles.Subtle.Render("No timer. No pressure."),
		"",
		"Pick any operation and practice",
		"at your own pace.",
	)

	preview := m.renderPracticePreview()

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		"",
		description,
		"",
		"",
		preview,
	)

	if m.width > 0 && m.viewportReady {
		styled := lipgloss.NewStyle().MarginTop(titleTopPadding(m.viewport.Height)).Render(content)
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Top, styled)
	}
	return content
}

// renderPracticePreview renders a mock preview of the practice screen.
func (m FeatureTourModel) renderPracticePreview() string {
	previewBoxWidth := min(m.width-4, 34)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(previewBoxWidth)

	centerStyle := lipgloss.NewStyle().Width(previewBoxWidth).Align(lipgloss.Center)

	header := centerStyle.Render(styles.Dim.Render("Basic · Addition · Medium"))
	question := centerStyle.Render(styles.Bold.Render("15 + 8 = ?"))
	input := centerStyle.Render(styles.Dim.Render("> ") + styles.Accent.Render("23") + styles.Dim.Render("█"))

	inner := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		question,
		input,
		"",
	)

	return boxStyle.Render(inner)
}

// renderStatisticsContent renders the statistics introduction.
func (m FeatureTourModel) renderStatisticsContent() string {
	title := styles.Logo.Render("STATISTICS")

	description := lipgloss.JoinVertical(lipgloss.Center,
		styles.Subtle.Render("Every session counts."),
		"",
		"Track accuracy, streaks, and response",
		"times across all modes.",
		"",
		styles.Dim.Render("Dive into history, trends, and session details."),
	)

	preview := m.renderStatisticsPreview()

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		"",
		description,
		"",
		"",
		preview,
	)

	if m.width > 0 && m.viewportReady {
		styled := lipgloss.NewStyle().MarginTop(titleTopPadding(m.viewport.Height)).Render(content)
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Top, styled)
	}
	return content
}

// renderStatisticsPreview renders a mock preview of the statistics dashboard.
func (m FeatureTourModel) renderStatisticsPreview() string {
	previewBoxWidth := min(m.width-4, 34)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(previewBoxWidth).
		Padding(0, 1)

	// Mock operation rows matching real dashboard format: symbol Name  XX%  bar
	row1 := "+  Addition        " + styles.Correct.Render("92%") + "  " + styles.Correct.Render("█████████") + styles.Dim.Render("░")
	row2 := "−  Subtraction     " + styles.Correct.Render("85%") + "  " + styles.Correct.Render("████████") + styles.Dim.Render("░░")
	row3 := "×  Multiplication  " + "68%" + "  " + "██████" + styles.Dim.Render("░░░░")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		"",
		row1,
		row2,
		row3,
		"",
	)

	return boxStyle.Render(inner)
}

// renderFinaleContent renders the finale screen with creator attribution.
func (m FeatureTourModel) renderFinaleContent() string {
	byLine := styles.Dim.Render("A game by ") + styles.Accent.Render("@gurselcakar")
	claudeLine := styles.Dim.Render("+ Claude Code")

	content := lipgloss.JoinVertical(lipgloss.Center,
		byLine,
		claudeLine,
	)

	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}
