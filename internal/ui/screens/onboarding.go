package screens

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// OnboardingStep represents the current step in the onboarding flow.
type OnboardingStep int

const (
	StepWelcome OnboardingStep = iota
	StepDuration
	StepDifficulty
	StepOperation
	StepInputMode
	StepReady
)

const totalSteps = 5 // Duration, Difficulty, Operation, InputMode, Ready (Welcome has no dots)

// OnboardingCompleteMsg is sent when onboarding is completed with selections.
type OnboardingCompleteMsg struct {
	ModeID      string
	Difficulty  string
	DurationMs  int64
	InputMethod string // "typing" or "multiple_choice"
}

// OnboardingSkipMsg is sent when the user skips onboarding.
type OnboardingSkipMsg struct{}

// Duration option for onboarding.
type durationOption struct {
	Label      string
	DurationMs int64
}

// Operation option for onboarding.
type operationOption struct {
	Label  string
	ModeID string
}

var durationOptions = []durationOption{
	{"30 seconds", 30000},
	{"60 seconds", 60000},
	{"90 seconds", 90000},
	{"2 minutes", 120000},
}

var difficultyOptions = []string{
	"Beginner",
	"Easy",
	"Medium",
	"Hard",
	"Expert",
}

var operationOptions = []operationOption{
	{"Addition", modes.IDAddition},
	{"Subtraction", modes.IDSubtraction},
	{"Multiplication", modes.IDMultiplication},
	{"Division", modes.IDDivision},
	{"Mixed Basics", modes.IDMixedBasics},
}

var inputModeOptions = []string{
	"Typing",
	"Multiple Choice",
}

// OnboardingModel represents the onboarding screen.
type OnboardingModel struct {
	step            OnboardingStep
	durationIndex   int
	difficultyIndex int
	operationIndex  int
	inputModeIndex  int
	width           int
	height          int
	viewport        viewport.Model
	viewportReady   bool
}

// NewOnboarding creates a new onboarding model.
func NewOnboarding() OnboardingModel {
	return OnboardingModel{
		step:            StepWelcome,
		durationIndex:   1, // Default: 60s
		difficultyIndex: 1, // Default: Easy
		operationIndex:  0, // Default: Addition
		inputModeIndex:  0, // Default: Typing
		viewport:        viewport.New(0, 0),
		viewportReady:   false,
	}
}

// Init initializes the onboarding model.
func (m OnboardingModel) Init() tea.Cmd {
	return nil
}

// Layout constants for fixed sections
const (
	progressHeight = 2 // Height reserved for progress dots (1 line + padding)
)

// Scroll constants for viewport navigation
const (
	// linesBeforeOptions is the number of lines before the options list:
	// title (1) + spacing (1) + subtitle (1) + spacing (3) = 6
	linesBeforeOptions = 6
	// scrollPaddingTop is the minimum lines to keep above the selection
	scrollPaddingTop = 2
	// scrollPaddingBottom is the minimum lines to keep below the selection
	scrollPaddingBottom = 3
)

// Update handles onboarding screen input.
func (m OnboardingModel) Update(msg tea.Msg) (OnboardingModel, tea.Cmd) {
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
		case "up", "k":
			m.moveSelection(-1)
			m.updateViewportContent()
			m.scrollToSelection()
		case "down", "j":
			m.moveSelection(1)
			m.updateViewportContent()
			m.scrollToSelection()
		case "enter", "right", "l":
			newModel, cmd := m.advance()
			newModel.updateViewportContent()
			return newModel, cmd
		case "s", "S":
			return m.skip()
		case "b", "B", "left", "h":
			m.back()
			m.updateViewportContent()
		}
	}

	// Update viewport (for mouse scrolling if enabled)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// moveSelection navigates within the current step's options.
func (m *OnboardingModel) moveSelection(delta int) {
	index := m.currentIndex()
	maxIndex := m.maxIndexForStep()

	index += delta
	if index < 0 {
		index = 0
	}
	if index > maxIndex {
		index = maxIndex
	}

	m.setCurrentIndex(index)
}

// currentIndex returns the selected index for the current step.
func (m OnboardingModel) currentIndex() int {
	switch m.step {
	case StepDuration:
		return m.durationIndex
	case StepDifficulty:
		return m.difficultyIndex
	case StepOperation:
		return m.operationIndex
	case StepInputMode:
		return m.inputModeIndex
	default:
		return 0
	}
}

// setCurrentIndex sets the selected index for the current step.
func (m *OnboardingModel) setCurrentIndex(index int) {
	switch m.step {
	case StepDuration:
		m.durationIndex = index
	case StepDifficulty:
		m.difficultyIndex = index
	case StepOperation:
		m.operationIndex = index
	case StepInputMode:
		m.inputModeIndex = index
	}
}

// maxIndexForStep returns the maximum index for the current step.
func (m OnboardingModel) maxIndexForStep() int {
	switch m.step {
	case StepDuration:
		return len(durationOptions) - 1
	case StepDifficulty:
		return len(difficultyOptions) - 1
	case StepOperation:
		return len(operationOptions) - 1
	case StepInputMode:
		return len(inputModeOptions) - 1
	default:
		return 0
	}
}

// advance moves to the next step, completing onboarding if at the final step.
func (m OnboardingModel) advance() (OnboardingModel, tea.Cmd) {
	switch m.step {
	case StepWelcome:
		m.step = StepDuration
	case StepDuration:
		m.step = StepDifficulty
	case StepDifficulty:
		m.step = StepOperation
	case StepOperation:
		m.step = StepInputMode
	case StepInputMode:
		m.step = StepReady
	case StepReady:
		return m, m.complete()
	}
	return m, nil
}

// back returns to the previous step.
func (m *OnboardingModel) back() {
	switch m.step {
	case StepDuration:
		m.step = StepWelcome
	case StepDifficulty:
		m.step = StepDuration
	case StepOperation:
		m.step = StepDifficulty
	case StepInputMode:
		m.step = StepOperation
	case StepReady:
		m.step = StepInputMode
	}
}

// skip returns the skip message to use defaults.
func (m OnboardingModel) skip() (OnboardingModel, tea.Cmd) {
	return m, func() tea.Msg {
		return OnboardingSkipMsg{}
	}
}

// complete returns the completion message with user selections.
func (m OnboardingModel) complete() tea.Cmd {
	inputMethod := "typing"
	if m.inputModeIndex == 1 {
		inputMethod = "multiple_choice"
	}

	return func() tea.Msg {
		return OnboardingCompleteMsg{
			ModeID:      operationOptions[m.operationIndex].ModeID,
			Difficulty:  difficultyOptions[m.difficultyIndex],
			DurationMs:  durationOptions[m.durationIndex].DurationMs,
			InputMethod: inputMethod,
		}
	}
}

// View renders the onboarding screen.
func (m OnboardingModel) View() string {
	if !m.viewportReady {
		return "Loading..."
	}

	// Get progress dots and hints for current step
	progress := m.getProgressForStep()
	hints := m.getHintsForStep()

	// All screens: viewport + progress dots + hints
	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		lipgloss.Place(m.width, progressHeight, lipgloss.Center, lipgloss.Center, progress),
		lipgloss.Place(m.width, components.HintsHeight, lipgloss.Center, lipgloss.Center, hints),
	)
}

// getProgressForStep returns the progress dots for the current step.
func (m OnboardingModel) getProgressForStep() string {
	switch m.step {
	case StepWelcome:
		return "" // No progress dots on welcome screen
	case StepDuration:
		return components.ProgressDotsColored(1, totalSteps)
	case StepDifficulty:
		return components.ProgressDotsColored(2, totalSteps)
	case StepOperation:
		return components.ProgressDotsColored(3, totalSteps)
	case StepInputMode:
		return components.ProgressDotsColored(4, totalSteps)
	case StepReady:
		return components.ProgressDotsColored(5, totalSteps)
	default:
		return ""
	}
}

// getHintsForStep returns the appropriate hints for the current step.
func (m OnboardingModel) getHintsForStep() string {
	switch m.step {
	case StepWelcome:
		return components.RenderHintsResponsive([]components.Hint{
			{Key: "S", Action: "Skip"},
			{Key: "→", Action: "Continue"},
		}, m.width)
	case StepReady:
		return components.RenderHintsResponsive([]components.Hint{
			{Key: "←", Action: "Back"},
			{Key: "→", Action: "Start"},
		}, m.width)
	default:
		return components.RenderHintsResponsive([]components.Hint{
			{Key: "←", Action: "Back"},
			{Key: "↑↓", Action: "Navigate"},
			{Key: "S", Action: "Skip"},
			{Key: "→", Action: "Continue"},
		}, m.width)
	}
}

// SetSize updates the screen dimensions.
func (m *OnboardingModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	viewportHeight := m.calculateViewportHeight()

	components.SetViewportSize(&m.viewport, &m.viewportReady, m.width, viewportHeight)

	m.updateViewportContent()
}

// calculateViewportHeight returns the viewport height.
func (m OnboardingModel) calculateViewportHeight() int {
	// Bottom sections are always present (progress dots + hints)
	bottomSectionHeight := components.HintsHeight + progressHeight

	viewportHeight := m.height - bottomSectionHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}

	return viewportHeight
}

// updateViewportContent updates the viewport with the current step's content.
func (m *OnboardingModel) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	content := m.getViewportContent()
	m.viewport.SetContent(content)
}

// scrollToSelection scrolls the viewport to keep the selected option visible.
func (m *OnboardingModel) scrollToSelection() {
	if !m.viewportReady || m.step == StepWelcome {
		return
	}

	optionsCount := m.maxIndexForStep() + 1
	contentHeight := linesBeforeOptions + optionsCount

	// Calculate where the content starts (centered)
	contentStart := (m.viewport.Height - contentHeight) / 2
	if contentStart < 0 {
		contentStart = 0
	}

	// Selection line within the viewport
	selectionLine := contentStart + linesBeforeOptions + m.currentIndex()

	// Get viewport bounds
	viewportTop := m.viewport.YOffset
	viewportBottom := viewportTop + m.viewport.Height

	// Scroll if selection is outside visible area with padding
	if selectionLine < viewportTop+scrollPaddingTop {
		m.viewport.SetYOffset(selectionLine - scrollPaddingTop)
	} else if selectionLine > viewportBottom-scrollPaddingBottom {
		m.viewport.SetYOffset(selectionLine - m.viewport.Height + scrollPaddingBottom)
	}
}

// getViewportContent returns the content for the current step.
func (m OnboardingModel) getViewportContent() string {
	switch m.step {
	case StepWelcome:
		return m.renderWelcomeContent()
	case StepDuration:
		return m.renderDurationContent()
	case StepDifficulty:
		return m.renderDifficultyContent()
	case StepOperation:
		return m.renderOperationContent()
	case StepInputMode:
		return m.renderInputModeContent()
	case StepReady:
		return m.renderReadyContent()
	default:
		return ""
	}
}

// renderWelcomeContent renders the welcome step content.
func (m OnboardingModel) renderWelcomeContent() string {
	logo := components.LogoColoredForWidth(m.width)
	separator := styles.Dim.Render(components.LogoSeparator())
	tagline := components.Tagline()
	setup := styles.Tagline.Render("Let's get you set up.")

	content := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		separator,
		"",
		tagline,
		"",
		"",
		"",
		setup,
	)

	// Center both horizontally and vertically within viewport
	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// renderDurationContent renders the duration step content.
func (m OnboardingModel) renderDurationContent() string {
	title := styles.Logo.Render("SESSION LENGTH")
	subtitle := styles.Subtle.Render("How long do you want to play?")

	var options []string
	for i, opt := range durationOptions {
		if i == m.durationIndex {
			options = append(options, styles.Accent.Render("> ")+styles.Bold.Render(opt.Label))
		} else {
			options = append(options, "  "+styles.Unselected.Render(opt.Label))
		}
	}
	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	return m.renderStepContent(title, subtitle, optionsList)
}

// renderDifficultyContent renders the difficulty step content.
func (m OnboardingModel) renderDifficultyContent() string {
	title := styles.Logo.Render("DIFFICULTY")
	subtitle := styles.Subtle.Render("What difficulty level?")

	var options []string
	for i, opt := range difficultyOptions {
		if i == m.difficultyIndex {
			options = append(options, styles.Accent.Render("> ")+styles.Bold.Render(opt))
		} else {
			options = append(options, "  "+styles.Unselected.Render(opt))
		}
	}
	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	return m.renderStepContent(title, subtitle, optionsList)
}

// renderOperationContent renders the operation step content.
func (m OnboardingModel) renderOperationContent() string {
	title := styles.Logo.Render("OPERATION")
	subtitle := styles.Subtle.Render("What do you want to play?")

	var options []string
	for i, opt := range operationOptions {
		if i == m.operationIndex {
			options = append(options, styles.Accent.Render("> ")+styles.Bold.Render(opt.Label))
		} else {
			options = append(options, "  "+styles.Unselected.Render(opt.Label))
		}
	}
	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	return m.renderStepContent(title, subtitle, optionsList)
}

// renderInputModeContent renders the input mode step content.
func (m OnboardingModel) renderInputModeContent() string {
	title := styles.Logo.Render("INPUT MODE")
	subtitle := styles.Subtle.Render("How do you want to answer?")

	var options []string
	for i, opt := range inputModeOptions {
		if i == m.inputModeIndex {
			options = append(options, styles.Accent.Render("> ")+styles.Bold.Render(opt))
		} else {
			options = append(options, "  "+styles.Unselected.Render(opt))
		}
	}
	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	// Preview based on selected input mode
	preview := m.renderInputModePreview()

	return m.renderStepContentWithPreview(title, subtitle, optionsList, preview)
}

// Preview box dimensions (fixed size for consistent layout)
const (
	previewBoxHeight = 5 // Inner content height (excludes border)
)

// renderInputModePreview renders a preview box showing how the selected input mode looks.
func (m OnboardingModel) renderInputModePreview() string {
	previewBoxWidth := min(m.width-4, 34)

	// Fixed-size box style for consistent layout
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(previewBoxWidth).
		Height(previewBoxHeight)

	// Question on its own line (matches game screen layout)
	question := styles.Bold.Render("1 + 1 =")

	var input string
	if m.inputModeIndex == 0 {
		// Typing preview: prompt, answer, and cursor (matches game screen)
		prompt := styles.Dim.Render("> ")
		answer := styles.Accent.Render("2")
		cursor := styles.Dim.Render("█")
		input = prompt + answer + cursor
	} else {
		// Multiple choice preview: four options
		input = lipgloss.JoinHorizontal(lipgloss.Center,
			styles.Dim.Render("[1] ")+"0",
			"  ",
			styles.Accent.Render("[2] ")+"2",
			"  ",
			styles.Dim.Render("[3] ")+"3",
			"  ",
			styles.Dim.Render("[4] ")+"1",
		)
	}

	// Center the question horizontally
	questionCentered := lipgloss.NewStyle().Width(previewBoxWidth).Align(lipgloss.Center).Render(question)

	// Input alignment: both centered
	inputStyled := lipgloss.NewStyle().Width(previewBoxWidth).Align(lipgloss.Center).Render(input)

	// Fixed layout: empty line, question, empty line, input, empty line
	content := lipgloss.JoinVertical(lipgloss.Left,
		"",
		questionCentered,
		"",
		inputStyled,
		"",
	)

	return boxStyle.Render(content)
}

// renderStepContentWithPreview renders title, subtitle, options, and a preview centered in the viewport.
func (m OnboardingModel) renderStepContentWithPreview(title, subtitle, options, preview string) string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		subtitle,
		"",
		"",
		"",
		options,
		"",
		"",
		preview,
	)

	// Center both horizontally and vertically within viewport
	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// renderStepContent renders title, subtitle, and options centered in the viewport.
func (m OnboardingModel) renderStepContent(title, subtitle, options string) string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		subtitle,
		"",
		"",
		"",
		options,
	)

	// Center both horizontally and vertically within viewport
	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// renderStepContentWithNote renders title, subtitle, options, and a note centered in the viewport.
func (m OnboardingModel) renderStepContentWithNote(title, subtitle, options, note string) string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		subtitle,
		"",
		"",
		"",
		options,
		"",
		"",
		note,
	)

	// Center both horizontally and vertically within viewport
	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// renderReadyContent renders the ready screen with summary and controls.
func (m OnboardingModel) renderReadyContent() string {
	title := styles.Logo.Render("READY TO PLAY")

	// Build summary of selections
	mode := operationOptions[m.operationIndex].Label
	difficulty := difficultyOptions[m.difficultyIndex]
	duration := durationOptions[m.durationIndex].Label
	inputMode := inputModeOptions[m.inputModeIndex]

	summary := lipgloss.JoinVertical(lipgloss.Left,
		styles.Dim.Render("Mode:       ")+styles.Bold.Render(mode),
		styles.Dim.Render("Difficulty: ")+styles.Bold.Render(difficulty),
		styles.Dim.Render("Duration:   ")+styles.Bold.Render(duration),
		styles.Dim.Render("Input:      ")+styles.Bold.Render(inputMode),
	)

	// Build controls section based on input mode
	var controls string
	if m.inputModeIndex == 1 { // Multiple Choice
		controls = lipgloss.JoinVertical(lipgloss.Left,
			styles.Dim.Render("[1-4] ")+"Select answer",
			styles.Dim.Render("[S]   ")+"Skip question",
			styles.Dim.Render("[P]   ")+"Pause game",
		)
	} else { // Typing
		controls = lipgloss.JoinVertical(lipgloss.Left,
			styles.Dim.Render("[Enter] ")+"Submit answer",
			styles.Dim.Render("[S]     ")+"Skip question",
			styles.Dim.Render("[P]     ")+"Pause game",
		)
	}

	controlsSection := lipgloss.JoinVertical(lipgloss.Left,
		styles.Subtle.Render("Controls:"),
		"",
		controls,
	)

	// Combine summary and controls in a left-aligned block
	infoBlock := lipgloss.JoinVertical(lipgloss.Left,
		summary,
		"",
		"",
		controlsSection,
	)

	// Styled start instruction (dim like key hints)
	startInstruction := styles.Dim.Render("─── Press → to Start ───")

	// Combine all parts - title centered, info block centered as a unit
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		infoBlock,
		"",
		"",
		startInstruction,
	)

	// Center both horizontally and vertically within viewport
	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}
