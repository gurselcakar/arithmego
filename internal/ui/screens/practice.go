package screens

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/operations"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// Feedback duration for correct/incorrect flash
const practiceFeedbackDuration = 400 * time.Millisecond

// InputMethod represents the answer input type.
type InputMethod int

const (
	InputTyping InputMethod = iota
	InputMultipleChoice
)

func (m InputMethod) String() string {
	switch m {
	case InputTyping:
		return "Typing"
	case InputMultipleChoice:
		return "Multiple Choice"
	default:
		return "Unknown"
	}
}

// PracticeSettingsField identifies which field is focused in the settings panel.
type PracticeSettingsField int

const (
	PracticeFieldOperation PracticeSettingsField = iota
	PracticeFieldDifficulty
	PracticeFieldInputMethod
)

// operationEntry represents an operation in the settings panel.
type operationEntry struct {
	op       game.Operation
	name     string
	symbol   string
	category game.Category
}

// PracticeModel represents the practice mode screen.
type PracticeModel struct {
	// Screen dimensions
	width  int
	height int

	// Settings panel state
	settingsOpen  bool
	settingsField PracticeSettingsField

	// Operation selection
	allOperations   []operationEntry // All operations organized for display
	operationIndex  int              // Currently selected operation index
	selectedOp      game.Operation   // Current operation (nil = mixed within category)
	isMixed         bool             // True when "Mixed" is selected
	mixedCategory   game.Category    // Which category to mix (empty = all basic)
	categoryIndices map[game.Category]int // Start index of each category in allOperations

	// Difficulty
	difficulty    game.Difficulty
	difficultyIdx int

	// Input method
	inputMethod InputMethod

	// Question state
	current *game.Question
	input   components.InputModel

	// Feedback state
	feedback       string    // "correct", "incorrect", or ""
	feedbackExpiry time.Time
}

// NewPractice creates a new practice model.
func NewPractice() PracticeModel {
	m := PracticeModel{
		settingsField:   PracticeFieldOperation,
		difficulty:      game.Medium,
		difficultyIdx:   2, // Medium is index 2
		inputMethod:     InputTyping,
		input:           components.NewInput(),
		categoryIndices: make(map[game.Category]int),
	}

	// Build operation list organized by category
	m.buildOperationList()

	// Start with Addition selected
	m.operationIndex = 0
	m.selectedOp = m.allOperations[0].op
	m.isMixed = false

	// Generate first question
	m.generateQuestion()

	return m
}

// buildOperationList organizes operations by category for the settings panel.
func (m *PracticeModel) buildOperationList() {
	m.allOperations = []operationEntry{}

	// Define category order and their operations
	categories := []struct {
		cat  game.Category
		name string
	}{
		{game.CategoryBasic, "Basic"},
		{game.CategoryPower, "Power"},
		{game.CategoryAdvanced, "Advanced"},
	}

	for _, c := range categories {
		m.categoryIndices[c.cat] = len(m.allOperations)
		ops := operations.ByCategory(c.cat)

		// Sort operations within category for consistent order
		sortedOps := sortOperations(ops, c.cat)

		for _, op := range sortedOps {
			m.allOperations = append(m.allOperations, operationEntry{
				op:       op,
				name:     op.Name(),
				symbol:   op.Symbol(),
				category: c.cat,
			})
		}

		// Add "Mixed" option for this category
		m.allOperations = append(m.allOperations, operationEntry{
			op:       nil, // nil indicates mixed
			name:     fmt.Sprintf("Mixed %s", c.name),
			symbol:   "*",
			category: c.cat,
		})
	}
}

// sortOperations returns operations in a consistent display order.
func sortOperations(ops []game.Operation, cat game.Category) []game.Operation {
	// Define preferred order by operation name
	order := map[string]int{
		// Basic
		"Addition": 0, "Subtraction": 1, "Multiplication": 2, "Division": 3,
		// Power
		"Square": 0, "Cube": 1, "Square Root": 2, "Cube Root": 3,
		// Advanced
		"Modulo": 0, "Power": 1, "Percentage": 2, "Factorial": 3,
	}

	sorted := make([]game.Operation, len(ops))
	copy(sorted, ops)

	// Simple insertion sort (small list)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && order[sorted[j].Name()] < order[sorted[j-1].Name()]; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}

	return sorted
}

// Init initializes the practice model.
// Note: First question is generated in NewPractice() since Init() has a value receiver.
func (m PracticeModel) Init() tea.Cmd {
	return m.input.Init()
}

// feedbackTickMsg is sent to clear feedback after duration.
type feedbackTickMsg struct{}

// feedbackTickCmd returns a command to clear feedback after the duration.
func feedbackTickCmd() tea.Cmd {
	return tea.Tick(practiceFeedbackDuration, func(t time.Time) tea.Msg {
		return feedbackTickMsg{}
	})
}

// Update handles practice screen input.
func (m PracticeModel) Update(msg tea.Msg) (PracticeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case feedbackTickMsg:
		if time.Now().After(m.feedbackExpiry) {
			m.feedback = ""
		}
		return m, nil

	case tea.KeyMsg:
		// Global keys that work regardless of settings panel state
		switch msg.String() {
		case "q":
			if !m.settingsOpen {
				return m, func() tea.Msg { return ReturnToMenuMsg{} }
			}
		case "esc":
			if m.settingsOpen {
				m.settingsOpen = false
				return m, m.input.Focus()
			}
			return m, func() tea.Msg { return ReturnToMenuMsg{} }
		}

		if m.settingsOpen {
			return m.updateSettingsPanel(msg)
		}
		return m.updatePractice(msg)
	}

	return m, nil
}

// updatePractice handles input when settings panel is closed.
func (m PracticeModel) updatePractice(msg tea.KeyMsg) (PracticeModel, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.settingsOpen = true
		m.input.Blur()
		return m, nil

	case "up", "k":
		m.adjustDifficulty(1)
		m.generateQuestion()
		return m, nil

	case "down", "j":
		m.adjustDifficulty(-1)
		m.generateQuestion()
		return m, nil

	case "1":
		m.selectOperationByName("Addition")
		return m, nil

	case "2":
		m.selectOperationByName("Subtraction")
		return m, nil

	case "3":
		m.selectOperationByName("Multiplication")
		return m, nil

	case "4":
		m.selectOperationByName("Division")
		return m, nil

	case "0", "m":
		m.selectMixedBasic()
		return m, nil

	case "s", " ":
		m.skip()
		return m, nil

	case "enter":
		return m.submitAnswer()

	default:
		// Pass to input component
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

// updateSettingsPanel handles input when settings panel is open.
func (m PracticeModel) updateSettingsPanel(msg tea.KeyMsg) (PracticeModel, tea.Cmd) {
	switch msg.String() {
	case "tab", "enter":
		if msg.String() == "enter" && m.settingsField == PracticeFieldOperation {
			// Confirm selection and close panel
			m.applySelectedOperation()
			m.settingsOpen = false
			return m, m.input.Focus()
		}
		// Tab cycles through fields
		m.settingsField = (m.settingsField + 1) % 3
		return m, nil

	case "up", "k":
		switch m.settingsField {
		case PracticeFieldOperation:
			if m.operationIndex > 0 {
				m.operationIndex--
			}
		case PracticeFieldDifficulty:
			m.adjustDifficulty(1)
		}
		return m, nil

	case "down", "j":
		switch m.settingsField {
		case PracticeFieldOperation:
			if m.operationIndex < len(m.allOperations)-1 {
				m.operationIndex++
			}
		case PracticeFieldDifficulty:
			m.adjustDifficulty(-1)
		}
		return m, nil

	case "left", "h":
		switch m.settingsField {
		case PracticeFieldDifficulty:
			m.adjustDifficulty(-1)
		case PracticeFieldInputMethod:
			// TODO: Enable when multiple choice is implemented (Phase 10)
			// For now, input method is locked to Typing
		}
		return m, nil

	case "right", "l":
		switch m.settingsField {
		case PracticeFieldDifficulty:
			m.adjustDifficulty(1)
		case PracticeFieldInputMethod:
			// TODO: Enable when multiple choice is implemented (Phase 10)
			// For now, input method is locked to Typing
		}
		return m, nil
	}

	return m, nil
}

// adjustDifficulty changes the difficulty by delta, clamping to valid range.
func (m *PracticeModel) adjustDifficulty(delta int) {
	diffs := game.AllDifficulties()
	m.difficultyIdx += delta
	if m.difficultyIdx < 0 {
		m.difficultyIdx = 0
	}
	if m.difficultyIdx >= len(diffs) {
		m.difficultyIdx = len(diffs) - 1
	}
	m.difficulty = diffs[m.difficultyIdx]
}

// selectOperationByName selects a specific operation by name.
func (m *PracticeModel) selectOperationByName(name string) {
	for i, entry := range m.allOperations {
		if entry.op != nil && entry.op.Name() == name {
			m.operationIndex = i
			m.selectedOp = entry.op
			m.isMixed = false
			m.generateQuestion()
			return
		}
	}
}

// selectMixedBasic selects mixed basic operations.
func (m *PracticeModel) selectMixedBasic() {
	// Find the "Mixed Basic" entry
	for i, entry := range m.allOperations {
		if entry.op == nil && entry.category == game.CategoryBasic {
			m.operationIndex = i
			m.selectedOp = nil
			m.isMixed = true
			m.mixedCategory = game.CategoryBasic
			m.generateQuestion()
			return
		}
	}
}

// applySelectedOperation applies the currently highlighted operation.
func (m *PracticeModel) applySelectedOperation() {
	entry := m.allOperations[m.operationIndex]
	if entry.op == nil {
		// Mixed mode for this category
		m.selectedOp = nil
		m.isMixed = true
		m.mixedCategory = entry.category
	} else {
		m.selectedOp = entry.op
		m.isMixed = false
	}
	m.generateQuestion()
}

// generateQuestion creates a new question based on current settings.
func (m *PracticeModel) generateQuestion() {
	var ops []game.Operation

	if m.isMixed {
		ops = operations.ByCategory(m.mixedCategory)
	} else if m.selectedOp != nil {
		ops = []game.Operation{m.selectedOp}
	} else {
		// Fallback to basic operations
		ops = operations.BasicOperations()
	}

	if len(ops) == 0 {
		ops = operations.BasicOperations()
	}

	q := game.GenerateQuestion(ops, m.difficulty)
	m.current = &q
	m.input.Reset()
}

// submitAnswer checks the answer and generates a new question.
func (m PracticeModel) submitAnswer() (PracticeModel, tea.Cmd) {
	val := m.input.Value()
	if val == "" {
		return m, nil
	}

	answer, err := strconv.Atoi(val)
	if err != nil {
		return m, nil
	}

	if m.current == nil {
		return m, nil
	}

	correct := m.current.CheckAnswer(answer).Correct

	if correct {
		m.feedback = "correct"
	} else {
		m.feedback = "incorrect"
		// TODO: Consider showing the correct answer here.
		// For now, we just flash red and move on.
		// The correct answer was: m.current.Answer
	}
	m.feedbackExpiry = time.Now().Add(practiceFeedbackDuration)

	m.generateQuestion()

	return m, feedbackTickCmd()
}

// skip moves to the next question without answering.
func (m *PracticeModel) skip() {
	m.feedback = ""
	m.generateQuestion()
}

// View renders the practice screen.
func (m PracticeModel) View() string {
	if m.settingsOpen {
		return m.viewWithSettingsPanel()
	}
	return m.viewClean()
}

// viewClean renders the minimal practice view.
func (m PracticeModel) viewClean() string {
	// Current operation and difficulty indicator (top-left)
	var opName string
	if m.isMixed {
		opName = fmt.Sprintf("Mixed %s", categoryDisplayName(m.mixedCategory))
	} else if m.selectedOp != nil {
		opName = m.selectedOp.Name()
	} else {
		opName = "Mixed"
	}
	indicator := fmt.Sprintf("%s\n%s", opName, m.difficulty.String())
	indicatorStyled := styles.Subtle.Render(indicator)

	// Question (center)
	var questionView string
	if m.current != nil {
		questionView = components.RenderQuestion(m.current.Display)
	}

	// Input with feedback styling
	inputView := m.input.View()
	switch m.feedback {
	case "correct":
		inputView = styles.Correct.Render(inputView)
	case "incorrect":
		inputView = styles.Incorrect.Render(inputView)
	}

	// Center content
	centerContent := lipgloss.JoinVertical(lipgloss.Center,
		questionView,
		"",
		inputView,
	)

	// Bottom bar
	separator := styles.Dim.Render(strings.Repeat("─", min(m.width-4, 78)))
	hints := components.RenderHints([]string{"Tab Settings", "↑↓ Difficulty", "1-4 Operation", "S Skip", "Q Quit"})
	bottomBar := lipgloss.JoinVertical(lipgloss.Center, separator, hints)

	// Layout
	if m.width > 0 && m.height > 0 {
		// Calculate layout dimensions
		indicatorWidth := 20
		centerWidth := m.width - indicatorWidth*2
		availHeight := m.height - 4 // Reserve space for bottom bar

		// Position indicator on left
		indicatorCol := lipgloss.NewStyle().
			Width(indicatorWidth).
			Height(availHeight).
			Align(lipgloss.Left).
			AlignVertical(lipgloss.Center).
			PaddingLeft(2).
			Render(indicatorStyled)

		// Center the question
		centerCol := lipgloss.NewStyle().
			Width(centerWidth).
			Height(availHeight).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(centerContent)

		// Empty right column for balance
		rightCol := lipgloss.NewStyle().
			Width(indicatorWidth).
			Height(availHeight).
			Render("")

		mainArea := lipgloss.JoinHorizontal(lipgloss.Top, indicatorCol, centerCol, rightCol)

		// Combine with bottom bar
		bottomBarStyled := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(bottomBar)

		return lipgloss.JoinVertical(lipgloss.Left, mainArea, bottomBarStyled)
	}

	// Fallback for unknown dimensions
	return lipgloss.JoinVertical(lipgloss.Center,
		indicatorStyled,
		"",
		"",
		centerContent,
		"",
		"",
		bottomBar,
	)
}

// viewWithSettingsPanel renders the view with the settings panel open.
func (m PracticeModel) viewWithSettingsPanel() string {
	// Build settings panel
	panel := m.renderSettingsPanel()

	// Question (stays visible)
	var questionView string
	if m.current != nil {
		questionView = components.RenderQuestion(m.current.Display)
	}

	// Input (dimmed when settings open)
	inputView := styles.Dim.Render(m.input.View())

	centerContent := lipgloss.JoinVertical(lipgloss.Center,
		questionView,
		"",
		inputView,
	)

	// Bottom bar (different hints when settings open)
	separator := styles.Dim.Render(strings.Repeat("─", min(m.width-4, 78)))
	hints := components.RenderHints([]string{"Tab Next", "↑↓ Select", "←→ Adjust", "Enter Confirm", "Esc Close"})
	bottomBar := lipgloss.JoinVertical(lipgloss.Center, separator, hints)

	// Layout
	if m.width > 0 && m.height > 0 {
		panelWidth := 28
		centerWidth := m.width - panelWidth - 4
		availHeight := m.height - 4

		// Settings panel on left
		panelCol := lipgloss.NewStyle().
			Width(panelWidth).
			Height(availHeight).
			Align(lipgloss.Left).
			AlignVertical(lipgloss.Top).
			PaddingLeft(2).
			PaddingTop(2).
			Render(panel)

		// Center the question
		centerCol := lipgloss.NewStyle().
			Width(centerWidth).
			Height(availHeight).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(centerContent)

		mainArea := lipgloss.JoinHorizontal(lipgloss.Top, panelCol, centerCol)

		bottomBarStyled := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(bottomBar)

		return lipgloss.JoinVertical(lipgloss.Left, mainArea, bottomBarStyled)
	}

	// Fallback
	return lipgloss.JoinVertical(lipgloss.Center,
		panel,
		"",
		centerContent,
		"",
		bottomBar,
	)
}

// renderSettingsPanel renders the settings panel content.
func (m PracticeModel) renderSettingsPanel() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.Bold.Render("SETTINGS"))
	b.WriteString("\n\n")

	// Operation section
	opSectionStyle := styles.Bold
	if m.settingsField != PracticeFieldOperation {
		opSectionStyle = styles.Subtle
	}
	b.WriteString(opSectionStyle.Render("Operation"))
	b.WriteString("\n")

	// Render operations grouped by category
	var currentCategory game.Category
	for i, entry := range m.allOperations {
		// Category header
		if entry.category != currentCategory {
			currentCategory = entry.category
			if i > 0 {
				b.WriteString("\n")
			}
			catName := categoryDisplayName(entry.category)
			b.WriteString(styles.Dim.Render("  " + catName))
			b.WriteString("\n")
		}

		// Operation entry
		prefix := "    "
		style := styles.Unselected
		if i == m.operationIndex && m.settingsField == PracticeFieldOperation {
			prefix = "  > "
			style = styles.Selected
		} else if i == m.operationIndex {
			prefix = "  > "
			style = styles.Normal
		}

		displayName := entry.name
		if entry.op == nil {
			displayName = "Mixed"
		}
		b.WriteString(style.Render(prefix + displayName))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Difficulty section
	diffSectionStyle := styles.Bold
	if m.settingsField != PracticeFieldDifficulty {
		diffSectionStyle = styles.Subtle
	}
	b.WriteString(diffSectionStyle.Render("Difficulty"))
	b.WriteString("\n")
	b.WriteString(m.renderHorizontalSelector(
		m.difficultyIdx,
		difficultyNames(),
		m.settingsField == PracticeFieldDifficulty,
	))
	b.WriteString("\n\n")

	// Input method section
	inputSectionStyle := styles.Bold
	if m.settingsField != PracticeFieldInputMethod {
		inputSectionStyle = styles.Subtle
	}
	b.WriteString(inputSectionStyle.Render("Input Method"))
	b.WriteString("\n")
	// TODO: Enable input method selection when multiple choice is implemented (Phase 10)
	// Currently locked to "Typing"
	b.WriteString(styles.Dim.Render("  Typing"))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render("  (Multiple choice coming soon)"))

	return b.String()
}

// renderHorizontalSelector renders a ◀ value ▶ selector.
func (m PracticeModel) renderHorizontalSelector(index int, options []string, focused bool) string {
	return components.RenderSelector(index, options, components.SelectorOptions{
		Prefix:  "  ",
		Focused: focused,
	})
}

// SetSize updates the screen dimensions.
func (m *PracticeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Helper functions

func categoryDisplayName(cat game.Category) string {
	switch cat {
	case game.CategoryBasic:
		return "Basic"
	case game.CategoryPower:
		return "Power"
	case game.CategoryAdvanced:
		return "Advanced"
	default:
		return string(cat)
	}
}

func difficultyNames() []string {
	diffs := game.AllDifficulties()
	names := make([]string, len(diffs))
	for i, d := range diffs {
		names[i] = d.String()
	}
	return names
}
