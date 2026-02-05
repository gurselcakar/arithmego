package screens

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/operations"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// operationEntry represents an operation for selection.
type operationEntry struct {
	op       game.Operation
	name     string
	symbol   string
	category game.Category
}

// PracticeSettings holds the practice mode configuration for persistence.
type PracticeSettings struct {
	Category    string // "basic", "power", "advanced"
	Operation   string // operation name or "Mixed"
	Difficulty  string // difficulty name
	InputMethod string // "typing" or "multiple_choice"
}

// PracticeModel represents the practice mode screen.
type PracticeModel struct {
	// Screen dimensions
	width  int
	height int

	// Category and operation selection
	categories       []game.Category    // Available categories
	categoryIdx      int                // Currently selected category index
	categoryOps      []operationEntry   // Operations for current category (including Mixed)
	operationIdx     int                // Currently selected operation within category
	selectedOp       game.Operation     // Current operation (nil = mixed)
	isMixed          bool               // True when "Mixed" is selected

	// Difficulty
	difficulty    game.Difficulty
	difficultyIdx int

	// Input method
	inputMethod components.InputMethod

	// Question state
	current   *game.Question
	input     components.InputModel
	choices   components.ChoicesModel
	showError bool // True when wrong answer submitted (typing mode)
}

// NewPractice creates a new practice model with default settings.
func NewPractice() PracticeModel {
	return NewPracticeWithSettings(nil)
}

// NewPracticeWithSettings creates a new practice model with the given settings.
// If settings is nil, uses defaults.
func NewPracticeWithSettings(settings *PracticeSettings) PracticeModel {
	m := PracticeModel{
		categories:    []game.Category{game.CategoryBasic, game.CategoryPower, game.CategoryAdvanced},
		categoryIdx:   0, // Start with Basic
		difficulty:    game.Medium,
		difficultyIdx: 2, // Medium is index 2
		inputMethod:   components.InputTyping,
		input:         components.NewInput(),
		choices:       components.NewChoices(),
	}

	// Apply saved settings if provided
	if settings != nil {
		// Restore category
		for i, cat := range m.categories {
			if string(cat) == settings.Category {
				m.categoryIdx = i
				break
			}
		}

		// Restore difficulty
		for i, diff := range game.AllDifficulties() {
			if diff.String() == settings.Difficulty {
				m.difficultyIdx = i
				m.difficulty = diff
				break
			}
		}

		// Restore input method
		m.inputMethod = components.ParseInputMethod(settings.InputMethod)
	}

	// Build operation list for current category
	m.buildCategoryOps()

	// Restore operation if settings provided, otherwise use first
	m.operationIdx = 0
	if settings != nil && settings.Operation != "" {
		for i, op := range m.categoryOps {
			if op.name == settings.Operation {
				m.operationIdx = i
				break
			}
		}
	}

	m.applySelectedOperation()

	return m
}

// Settings returns the current practice settings for persistence.
func (m PracticeModel) Settings() PracticeSettings {
	var opName string
	if m.isMixed {
		opName = "Mixed"
	} else if m.selectedOp != nil {
		opName = m.selectedOp.Name()
	} else {
		opName = "Mixed"
	}

	// Use persistence-compatible format for input method
	inputMethodStr := "typing"
	if m.inputMethod == components.InputMultipleChoice {
		inputMethodStr = "multiple_choice"
	}

	return PracticeSettings{
		Category:    string(m.categories[m.categoryIdx]),
		Operation:   opName,
		Difficulty:  m.difficulty.String(),
		InputMethod: inputMethodStr,
	}
}

// buildCategoryOps builds the operation list for the current category.
func (m *PracticeModel) buildCategoryOps() {
	cat := m.categories[m.categoryIdx]
	ops := operations.ByCategory(cat)
	sortedOps := sortOperations(ops, cat)

	m.categoryOps = []operationEntry{}
	for _, op := range sortedOps {
		m.categoryOps = append(m.categoryOps, operationEntry{
			op:       op,
			name:     op.Name(),
			symbol:   op.Symbol(),
			category: cat,
		})
	}

	// Add "Mixed" option at the end
	m.categoryOps = append(m.categoryOps, operationEntry{
		op:       nil,
		name:     "Mixed",
		symbol:   "*",
		category: cat,
	})
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

// Update handles practice screen input.
func (m PracticeModel) Update(msg tea.Msg) (PracticeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case components.ChoiceSelectedMsg:
		// Handle multiple choice selection with retry on wrong answer
		if m.current == nil {
			return m, nil
		}
		result := m.current.CheckAnswer(msg.Value)
		if result.Correct {
			// Correct - move to next question
			m.choices.ClearError()
			m.generateQuestion()
		} else {
			// Wrong - mark choice as error, let user retry
			m.choices.SetError(msg.Index)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, func() tea.Msg { return ReturnToMenuMsg{} }

		case "c":
			// Cycle category
			m.cycleCategory()
			m.applySelectedOperation()
			return m, nil

		case "o":
			// Cycle operation within category
			m.cycleOperation()
			m.applySelectedOperation()
			return m, nil

		case "d":
			// Cycle difficulty
			m.cycleDifficulty()
			m.generateQuestion()
			return m, nil

		case "up", "k":
			m.adjustDifficulty(1)
			m.generateQuestion()
			return m, nil

		case "down", "j":
			m.adjustDifficulty(-1)
			m.generateQuestion()
			return m, nil

		case "m":
			// Toggle input method
			m.toggleInputMethod()
			m.generateQuestion()
			return m, nil

		case "s", " ":
			m.skip()
			return m, nil

		case "enter":
			if m.inputMethod == components.InputTyping {
				return m.submitAnswer()
			}
			return m, nil

		default:
			if m.inputMethod == components.InputMultipleChoice {
				// Route to choices component for answer selection
				var cmd tea.Cmd
				m.choices, cmd = m.choices.Update(msg)
				return m, cmd
			}
			// Typing mode: pass to text input component
			oldValue := m.input.Value()
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			// Clear error state when user modifies input
			if m.showError && m.input.Value() != oldValue {
				m.showError = false
			}
			return m, cmd
		}
	}

	return m, nil
}


// cycleCategory cycles to the next category (wrapping around).
func (m *PracticeModel) cycleCategory() {
	m.categoryIdx = (m.categoryIdx + 1) % len(m.categories)
	// Rebuild operations for new category
	m.buildCategoryOps()
	// Reset operation index to first
	m.operationIdx = 0
}

// cycleOperation cycles to the next operation within current category (wrapping around).
func (m *PracticeModel) cycleOperation() {
	m.operationIdx = (m.operationIdx + 1) % len(m.categoryOps)
}

// cycleDifficulty cycles to the next difficulty (wrapping around).
func (m *PracticeModel) cycleDifficulty() {
	diffs := game.AllDifficulties()
	m.difficultyIdx = (m.difficultyIdx + 1) % len(diffs)
	m.difficulty = diffs[m.difficultyIdx]
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

// applySelectedOperation applies the currently selected operation and generates a new question.
func (m *PracticeModel) applySelectedOperation() {
	if m.operationIdx >= len(m.categoryOps) {
		return
	}
	entry := m.categoryOps[m.operationIdx]
	if entry.op == nil {
		// Mixed mode for this category
		m.selectedOp = nil
		m.isMixed = true
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
		// Mixed mode: use all operations from current category
		ops = operations.ByCategory(m.categories[m.categoryIdx])
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
	m.choices.Reset()

	// Generate choices if in multiple choice mode
	if m.inputMethod == components.InputMultipleChoice {
		choices, correctIndex := game.GenerateChoices(q.Answer, m.difficulty)
		m.choices.SetChoices(choices, correctIndex)
	}
}

// toggleInputMethod switches between typing and multiple choice modes.
func (m *PracticeModel) toggleInputMethod() {
	if m.inputMethod == components.InputTyping {
		m.inputMethod = components.InputMultipleChoice
		// Generate choices for current question
		if m.current != nil {
			choices, correctIndex := game.GenerateChoices(m.current.Answer, m.difficulty)
			m.choices.SetChoices(choices, correctIndex)
		}
	} else {
		m.inputMethod = components.InputTyping
	}
}

// submitAnswer checks the answer and handles correct/incorrect.
func (m PracticeModel) submitAnswer() (PracticeModel, tea.Cmd) {
	val := m.input.Value()
	if val == "" {
		return m, nil
	}

	answer, err := strconv.Atoi(val)
	if err != nil {
		return m, nil
	}

	return m.submitAnswerValue(answer)
}

// submitAnswerValue submits an answer.
// Correct: advance to next question. Incorrect: show error, let user retry.
func (m PracticeModel) submitAnswerValue(answer int) (PracticeModel, tea.Cmd) {
	if m.current == nil {
		return m, nil
	}

	result := m.current.CheckAnswer(answer)

	if result.Correct {
		// Correct - move to next question
		m.showError = false
		m.generateQuestion()
	} else {
		// Incorrect - show error, keep input for retry
		m.showError = true
	}

	return m, nil
}

// skip moves to the next question without answering.
func (m *PracticeModel) skip() {
	m.showError = false
	m.generateQuestion()
}

// View renders the practice screen.
func (m PracticeModel) View() string {
	return m.viewPractice()
}

// viewPractice renders the practice view.
func (m PracticeModel) viewPractice() string {
	// Header bar with current settings
	var opName string
	if m.isMixed {
		opName = "Mixed"
	} else if m.selectedOp != nil {
		opName = m.selectedOp.Name()
	} else {
		opName = "Mixed"
	}
	catName := categoryDisplayName(m.categories[m.categoryIdx])
	inputMethodName := "Typing"
	if m.inputMethod == components.InputMultipleChoice {
		inputMethodName = "Choice"
	}

	// Clean header without shortcuts
	header := fmt.Sprintf("%s • %s • %s • %s", catName, opName, m.difficulty.String(), inputMethodName)
	headerStyled := styles.Subtle.Render(header)

	// Question (center)
	var questionView string
	if m.current != nil {
		questionView = components.RenderQuestion(m.current.Display)
	}

	// Input (with red styling on error for typing mode)
	var inputContent string
	if m.inputMethod == components.InputMultipleChoice {
		inputContent = m.choices.View()
	} else {
		inputContent = m.input.View()
		if m.showError {
			inputContent = styles.Incorrect.Render(inputContent)
		}
	}
	// Use screen width to prevent layout shift when switching input methods
	inputWidth := m.width
	if inputWidth == 0 {
		inputWidth = 80 // fallback
	}
	inputView := lipgloss.NewStyle().
		Width(inputWidth).
		Align(lipgloss.Center).
		Render(inputContent)

	// Center content
	centerContent := lipgloss.JoinVertical(lipgloss.Center,
		questionView,
		"",
		inputView,
	)

	// Hints - include settings shortcuts
	hintsWidth := m.width
	if hintsWidth == 0 {
		hintsWidth = 80
	}
	var hints string
	if m.inputMethod == components.InputMultipleChoice {
		hints = components.RenderHintsResponsive([]components.Hint{
			{Key: "C", Action: "Category"},
			{Key: "O", Action: "Operation"},
			{Key: "D", Action: "Difficulty"},
			{Key: "M", Action: "Input"},
			{Key: "1-4", Action: "Answer"},
			{Key: "S", Action: "Skip"},
			{Key: "Q", Action: "Quit"},
		}, hintsWidth)
	} else {
		hints = components.RenderHintsResponsive([]components.Hint{
			{Key: "C", Action: "Category"},
			{Key: "O", Action: "Operation"},
			{Key: "D", Action: "Difficulty"},
			{Key: "M", Action: "Input"},
			{Key: "S", Action: "Skip"},
			{Key: "Q", Action: "Quit"},
		}, hintsWidth)
	}

	// Layout with header at top, content centered, hints at bottom
	if m.width > 0 && m.height > 0 {
		hintsHeight := lipgloss.Height(hints)
		headerHeight := 2 // header + padding
		bottomPadding := 1
		availableHeight := m.height - hintsHeight - bottomPadding - headerHeight

		// Header centered at top
		centeredHeader := lipgloss.Place(m.width, headerHeight, lipgloss.Center, lipgloss.Top, headerStyled)

		// Center content on full width
		centeredContent := lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, centerContent)

		// Center hints at bottom with padding
		centeredHints := lipgloss.Place(m.width, hintsHeight+bottomPadding, lipgloss.Center, lipgloss.Top, hints)

		return lipgloss.JoinVertical(lipgloss.Left, centeredHeader, centeredContent, centeredHints)
	}

	// Fallback for unknown dimensions
	return lipgloss.JoinVertical(lipgloss.Center,
		headerStyled,
		"",
		"",
		centerContent,
		"",
		"",
		hints,
	)
}

// SetSize updates the screen dimensions.
func (m *PracticeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

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

