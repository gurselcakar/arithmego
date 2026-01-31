package screens

import (
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
)

const totalSteps = 5

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
	{"Addition", modes.IDAdditionSprint},
	{"Subtraction", modes.IDSubtractionSprint},
	{"Multiplication", modes.IDMultiplicationSprint},
	{"Division", modes.IDDivisionSprint},
	{"Mixed Operations", modes.IDMixedOperations},
}

var inputModeOptions = []string{
	"Typing",
	"Multiple Choice",
}

// OnboardingModel represents the onboarding screen.
type OnboardingModel struct {
	step            OnboardingStep
	cursor          int
	durationIndex   int
	difficultyIndex int
	operationIndex  int
	inputModeIndex  int
	width           int
	height          int
}

// NewOnboarding creates a new onboarding model.
func NewOnboarding() OnboardingModel {
	return OnboardingModel{
		step:            StepWelcome,
		cursor:          0,
		durationIndex:   1, // Default: 60s
		difficultyIndex: 1, // Default: Easy
		operationIndex:  0, // Default: Addition
		inputModeIndex:  0, // Default: Typing
	}
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
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case "enter":
			return m.advance()
		case "s", "S":
			return m.skip()
		case "b", "B":
			m.back()
		}
	}

	return m, nil
}

// moveCursor navigates within the current step's options.
func (m *OnboardingModel) moveCursor(delta int) {
	maxCursor := m.maxCursorForStep()
	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}
}

// maxCursorForStep returns the maximum cursor position for the current step.
func (m OnboardingModel) maxCursorForStep() int {
	switch m.step {
	case StepWelcome:
		return 0
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

// advance saves the current selection and moves to the next step.
func (m OnboardingModel) advance() (OnboardingModel, tea.Cmd) {
	switch m.step {
	case StepWelcome:
		m.step = StepDuration
		m.cursor = m.durationIndex
	case StepDuration:
		m.durationIndex = m.cursor
		m.step = StepDifficulty
		m.cursor = m.difficultyIndex
	case StepDifficulty:
		m.difficultyIndex = m.cursor
		m.step = StepOperation
		m.cursor = m.operationIndex
	case StepOperation:
		m.operationIndex = m.cursor
		m.step = StepInputMode
		m.cursor = m.inputModeIndex
	case StepInputMode:
		m.inputModeIndex = m.cursor
		return m, m.complete()
	}
	return m, nil
}

// back returns to the previous step.
func (m *OnboardingModel) back() {
	switch m.step {
	case StepDuration:
		m.step = StepWelcome
		m.cursor = 0
	case StepDifficulty:
		m.step = StepDuration
		m.cursor = m.durationIndex
	case StepOperation:
		m.step = StepDifficulty
		m.cursor = m.difficultyIndex
	case StepInputMode:
		m.step = StepOperation
		m.cursor = m.operationIndex
	}
	// No-op on Welcome step
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
	var content string

	switch m.step {
	case StepWelcome:
		content = m.viewWelcome()
	case StepDuration:
		content = m.viewDuration()
	case StepDifficulty:
		content = m.viewDifficulty()
	case StepOperation:
		content = m.viewOperation()
	case StepInputMode:
		content = m.viewInputMode()
	}

	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

// viewWelcome renders the welcome step.
func (m OnboardingModel) viewWelcome() string {
	logo := components.LogoForWidth(m.width)
	tagline := components.Tagline()

	intro := "Welcome! Let's set up your first session."

	progress := styles.Dim.Render(components.ProgressDots(1, totalSteps))
	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "S", Action: "Skip"},
		{Key: "Enter", Action: "Continue"},
	})

	return lipgloss.JoinVertical(lipgloss.Center,
		logo,
		tagline,
		"",
		"",
		intro,
		"",
		"",
		progress,
		"",
		hints,
	)
}

// viewDuration renders the duration selection step.
func (m OnboardingModel) viewDuration() string {
	title := styles.Bold.Render("SESSION LENGTH")
	subtitle := styles.Subtle.Render("How long do you want to play?")

	var options []string
	for i, opt := range durationOptions {
		if i == m.cursor {
			options = append(options, styles.Selected.Render("> "+opt.Label))
		} else {
			options = append(options, styles.Unselected.Render("  "+opt.Label))
		}
	}

	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	progress := styles.Dim.Render(components.ProgressDots(2, totalSteps))
	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "B", Action: "Back"},
		{Key: "↑↓", Action: "Navigate"},
		{Key: "S", Action: "Skip"},
		{Key: "Enter", Action: "Continue"},
	})

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		subtitle,
		"",
		"",
		optionsList,
		"",
		"",
		progress,
		"",
		hints,
	)
}

// viewDifficulty renders the difficulty selection step.
func (m OnboardingModel) viewDifficulty() string {
	title := styles.Bold.Render("DIFFICULTY")
	subtitle := styles.Subtle.Render("Select your starting level")

	var options []string
	for i, opt := range difficultyOptions {
		if i == m.cursor {
			options = append(options, styles.Selected.Render("> "+opt))
		} else {
			options = append(options, styles.Unselected.Render("  "+opt))
		}
	}

	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	progress := styles.Dim.Render(components.ProgressDots(3, totalSteps))
	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "B", Action: "Back"},
		{Key: "↑↓", Action: "Navigate"},
		{Key: "S", Action: "Skip"},
		{Key: "Enter", Action: "Continue"},
	})

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		subtitle,
		"",
		"",
		optionsList,
		"",
		"",
		progress,
		"",
		hints,
	)
}

// viewOperation renders the operation selection step.
func (m OnboardingModel) viewOperation() string {
	title := styles.Bold.Render("OPERATION")
	subtitle := styles.Subtle.Render("What would you like to practice?")

	var options []string
	for i, opt := range operationOptions {
		if i == m.cursor {
			options = append(options, styles.Selected.Render("> "+opt.Label))
		} else {
			options = append(options, styles.Unselected.Render("  "+opt.Label))
		}
	}

	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	progress := styles.Dim.Render(components.ProgressDots(4, totalSteps))
	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "B", Action: "Back"},
		{Key: "↑↓", Action: "Navigate"},
		{Key: "S", Action: "Skip"},
		{Key: "Enter", Action: "Continue"},
	})

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		subtitle,
		"",
		"",
		optionsList,
		"",
		"",
		progress,
		"",
		hints,
	)
}

// viewInputMode renders the input mode selection step.
func (m OnboardingModel) viewInputMode() string {
	title := styles.Bold.Render("INPUT MODE")
	subtitle := styles.Subtle.Render("How would you like to answer?")

	var options []string
	for i, opt := range inputModeOptions {
		label := opt
		if opt == "Typing" {
			label = opt + " - Type your answers"
		} else if opt == "Multiple Choice" {
			label = opt + " - Select from options"
		}
		if i == m.cursor {
			options = append(options, styles.Selected.Render("> "+label))
		} else {
			options = append(options, styles.Unselected.Render("  "+label))
		}
	}

	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	progress := styles.Dim.Render(components.ProgressDots(5, totalSteps))
	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "B", Action: "Back"},
		{Key: "↑↓", Action: "Navigate"},
		{Key: "S", Action: "Skip"},
		{Key: "Enter", Action: "Start Game"},
	})

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		subtitle,
		"",
		"",
		optionsList,
		"",
		"",
		progress,
		"",
		hints,
	)
}

// SetSize updates the screen dimensions.
func (m *OnboardingModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
