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

// PracticeSettingsField identifies which field is focused in the settings panel.
type PracticeSettingsField int

const (
	PracticeFieldCategory PracticeSettingsField = iota
	PracticeFieldOperation
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
	current *game.Question
	input   components.InputModel
	choices components.ChoicesModel

	// Feedback state
	feedback       string    // "correct", "incorrect", or ""
	feedbackExpiry time.Time
}

// NewPractice creates a new practice model.
func NewPractice() PracticeModel {
	m := PracticeModel{
		settingsField: PracticeFieldCategory,
		categories:    []game.Category{game.CategoryBasic, game.CategoryPower, game.CategoryAdvanced},
		categoryIdx:   0, // Start with Basic
		difficulty:    game.Medium,
		difficultyIdx: 2, // Medium is index 2
		inputMethod:   components.InputTyping,
		input:         components.NewInput(),
		choices:       components.NewChoices(),
	}

	// Build operation list for initial category
	m.buildCategoryOps()

	// Start with first operation (Addition) and generate first question
	m.operationIdx = 0
	m.applySelectedOperation()

	return m
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

	case components.ChoiceSelectedMsg:
		// Auto-submit on choice selection (multiple choice mode)
		return m.submitAnswerValue(msg.Value)

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
		m.choices.Blur()
		return m, nil

	case "up", "k":
		m.adjustDifficulty(1)
		m.generateQuestion()
		return m, nil

	case "down", "j":
		m.adjustDifficulty(-1)
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
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

// updateSettingsPanel handles input when settings panel is open.
func (m PracticeModel) updateSettingsPanel(msg tea.KeyMsg) (PracticeModel, tea.Cmd) {
	switch msg.String() {
	case "tab", "enter":
		// Close panel and apply changes
		m.applySelectedOperation()
		m.settingsOpen = false
		return m, m.input.Focus()

	case "up", "k":
		// Move to previous field
		if m.settingsField > 0 {
			m.settingsField--
		}
		return m, nil

	case "down", "j":
		// Move to next field
		if m.settingsField < PracticeFieldInputMethod {
			m.settingsField++
		}
		return m, nil

	case "left", "h":
		switch m.settingsField {
		case PracticeFieldCategory:
			m.adjustCategory(-1)
			m.applySelectedOperation()
		case PracticeFieldOperation:
			m.adjustOperation(-1)
			m.applySelectedOperation()
		case PracticeFieldDifficulty:
			m.adjustDifficulty(-1)
			m.generateQuestion()
		case PracticeFieldInputMethod:
			m.toggleInputMethod()
			m.generateQuestion()
		}
		return m, nil

	case "right", "l":
		switch m.settingsField {
		case PracticeFieldCategory:
			m.adjustCategory(1)
			m.applySelectedOperation()
		case PracticeFieldOperation:
			m.adjustOperation(1)
			m.applySelectedOperation()
		case PracticeFieldDifficulty:
			m.adjustDifficulty(1)
			m.generateQuestion()
		case PracticeFieldInputMethod:
			m.toggleInputMethod()
			m.generateQuestion()
		}
		return m, nil
	}

	return m, nil
}

// adjustCategory changes the category by delta, rebuilding operations.
func (m *PracticeModel) adjustCategory(delta int) {
	m.categoryIdx += delta
	if m.categoryIdx < 0 {
		m.categoryIdx = 0
	}
	if m.categoryIdx >= len(m.categories) {
		m.categoryIdx = len(m.categories) - 1
	}
	// Rebuild operations for new category
	m.buildCategoryOps()
	// Reset operation index to first
	m.operationIdx = 0
}

// adjustOperation changes the operation by delta within current category.
func (m *PracticeModel) adjustOperation(delta int) {
	m.operationIdx += delta
	if m.operationIdx < 0 {
		m.operationIdx = 0
	}
	if m.operationIdx >= len(m.categoryOps) {
		m.operationIdx = len(m.categoryOps) - 1
	}
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

	return m.submitAnswerValue(answer)
}

// submitAnswerValue submits an answer and handles feedback.
// Used by both typing mode and multiple choice mode.
func (m PracticeModel) submitAnswerValue(answer int) (PracticeModel, tea.Cmd) {
	if m.current == nil {
		return m, nil
	}

	correct := m.current.CheckAnswer(answer).Correct

	if correct {
		m.feedback = "correct"
	} else {
		m.feedback = "incorrect"
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
	header := fmt.Sprintf("%s • %s • %s • %s", catName, opName, m.difficulty.String(), inputMethodName)
	headerStyled := styles.Subtle.Render(header)

	// Question (center)
	var questionView string
	if m.current != nil {
		questionView = components.RenderQuestion(m.current.Display)
	}

	// Input with feedback styling, wrapped in fixed-width container
	var inputContent string
	if m.inputMethod == components.InputMultipleChoice {
		inputContent = m.choices.View()
	} else {
		inputContent = m.input.View()
	}
	switch m.feedback {
	case "correct":
		inputContent = styles.Correct.Render(inputContent)
	case "incorrect":
		inputContent = styles.Incorrect.Render(inputContent)
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

	// Hints - differ based on input mode
	var hints string
	if m.inputMethod == components.InputMultipleChoice {
		hints = components.RenderHintsStructured([]components.Hint{
			{Key: "Q", Action: "Quit"},
			{Key: "Tab", Action: "Settings"},
			{Key: "↑↓", Action: "Difficulty"},
			{Key: "1-4", Action: "Select"},
			{Key: "S", Action: "Skip"},
		})
	} else {
		hints = components.RenderHintsStructured([]components.Hint{
			{Key: "Q", Action: "Quit"},
			{Key: "Tab", Action: "Settings"},
			{Key: "↑↓", Action: "Difficulty"},
			{Key: "S", Action: "Skip"},
		})
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

// viewWithSettingsPanel renders the view with the settings panel open.
func (m PracticeModel) viewWithSettingsPanel() string {
	// Build settings panel
	panel := m.renderSettingsPanel()

	// Header bar with current settings (same as viewClean)
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
	header := fmt.Sprintf("%s • %s • %s • %s", catName, opName, m.difficulty.String(), inputMethodName)
	headerStyled := styles.Subtle.Render(header)

	// Question (stays visible)
	var questionView string
	if m.current != nil {
		questionView = components.RenderQuestion(m.current.Display)
	}

	// Input (dimmed when settings open), wrapped in fixed-width container
	var inputContent string
	if m.inputMethod == components.InputMultipleChoice {
		inputContent = styles.Dim.Render(m.choices.View())
	} else {
		inputContent = styles.Dim.Render(m.input.View())
	}
	inputWidth := m.width
	if inputWidth == 0 {
		inputWidth = 80 // fallback
	}
	inputView := lipgloss.NewStyle().
		Width(inputWidth).
		Align(lipgloss.Center).
		Render(inputContent)

	centerContent := lipgloss.JoinVertical(lipgloss.Center,
		questionView,
		"",
		inputView,
	)

	// Hints (different hints when settings open)
	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "Tab", Action: "Close"},
		{Key: "↑↓", Action: "Field"},
		{Key: "←→", Action: "Change"},
	})

	// Layout: panel overlaid on left, content stays centered on full screen
	if m.width > 0 && m.height > 0 {
		hintsHeight := lipgloss.Height(hints)
		headerHeight := 2
		bottomPadding := 1
		panelWidth := 30
		availableHeight := m.height - hintsHeight - bottomPadding - headerHeight

		// Header centered at top (same position as viewClean)
		centeredHeader := lipgloss.Place(m.width, headerHeight, lipgloss.Center, lipgloss.Top, headerStyled)

		// Settings panel with vertical centering
		panelContent := lipgloss.NewStyle().PaddingLeft(2).Render(panel)
		panelStyled := lipgloss.Place(panelWidth, availableHeight, lipgloss.Left, lipgloss.Center, panelContent)

		// Center content on full width (same as viewClean)
		centeredContent := lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, centerContent)

		// Overlay: panel on left, content from panelWidth onwards
		panelLines := strings.Split(panelStyled, "\n")
		contentLines := strings.Split(centeredContent, "\n")

		var resultLines []string
		for i := 0; i < len(contentLines); i++ {
			panelLine := ""
			if i < len(panelLines) {
				panelLine = panelLines[i]
			}
			contentLine := contentLines[i]

			// Pad panel to exact width
			panelActualWidth := lipgloss.Width(panelLine)
			if panelActualWidth < panelWidth {
				panelLine = panelLine + strings.Repeat(" ", panelWidth-panelActualWidth)
			}

			// Append content after panel area (content starts with spaces due to centering)
			if lipgloss.Width(contentLine) > panelWidth {
				resultLines = append(resultLines, panelLine+contentLine[panelWidth:])
			} else {
				resultLines = append(resultLines, panelLine)
			}
		}
		mainArea := strings.Join(resultLines, "\n")

		// Center hints at bottom with padding
		centeredHints := lipgloss.Place(m.width, hintsHeight+bottomPadding, lipgloss.Center, lipgloss.Top, hints)

		return lipgloss.JoinVertical(lipgloss.Left, centeredHeader, mainArea, centeredHints)
	}

	// Fallback
	return lipgloss.JoinVertical(lipgloss.Center,
		panel,
		"",
		centerContent,
		"",
		hints,
	)
}

// renderSettingsPanel renders the settings panel content.
func (m PracticeModel) renderSettingsPanel() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.Bold.Render("SETTINGS"))
	b.WriteString("\n\n")

	// Category section
	catSectionStyle := styles.Bold
	if m.settingsField != PracticeFieldCategory {
		catSectionStyle = styles.Subtle
	}
	b.WriteString(catSectionStyle.Render("Category"))
	b.WriteString("\n")
	b.WriteString(m.renderHorizontalSelector(
		m.categoryIdx,
		categoryNames(m.categories),
		m.settingsField == PracticeFieldCategory,
	))
	b.WriteString("\n\n")

	// Operation section
	opSectionStyle := styles.Bold
	if m.settingsField != PracticeFieldOperation {
		opSectionStyle = styles.Subtle
	}
	b.WriteString(opSectionStyle.Render("Operation"))
	b.WriteString("\n")
	b.WriteString(m.renderHorizontalSelector(
		m.operationIdx,
		operationNames(m.categoryOps),
		m.settingsField == PracticeFieldOperation,
	))
	b.WriteString("\n\n")

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
	inputMethodIdx := 0
	if m.inputMethod == components.InputMultipleChoice {
		inputMethodIdx = 1
	}
	b.WriteString(m.renderHorizontalSelector(
		inputMethodIdx,
		[]string{"Typing", "Choice"},
		m.settingsField == PracticeFieldInputMethod,
	))

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

func categoryNames(cats []game.Category) []string {
	names := make([]string, len(cats))
	for i, c := range cats {
		names[i] = categoryDisplayName(c)
	}
	return names
}

func operationNames(ops []operationEntry) []string {
	names := make([]string, len(ops))
	for i, e := range ops {
		names[i] = e.name
	}
	return names
}
