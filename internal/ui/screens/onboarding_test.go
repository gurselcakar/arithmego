package screens

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gurselcakar/arithmego/internal/modes"
)

func TestNewOnboarding(t *testing.T) {
	m := NewOnboarding()

	if m.step != StepWelcome {
		t.Errorf("expected step to be StepWelcome, got %v", m.step)
	}
	if m.durationIndex != 0 {
		t.Errorf("expected durationIndex to be 0 (30s), got %d", m.durationIndex)
	}
	if m.difficultyIndex != 0 {
		t.Errorf("expected difficultyIndex to be 0 (Beginner), got %d", m.difficultyIndex)
	}
	if m.operationIndex != 0 {
		t.Errorf("expected operationIndex to be 0 (Addition), got %d", m.operationIndex)
	}
	if m.inputModeIndex != 0 {
		t.Errorf("expected inputModeIndex to be 0 (Typing), got %d", m.inputModeIndex)
	}
}

func TestOnboardingAdvance(t *testing.T) {
	tests := []struct {
		name         string
		startStep    OnboardingStep
		expectedStep OnboardingStep
		expectCmd    bool
	}{
		{"Welcome to Duration", StepWelcome, StepDuration, false},
		{"Duration to Difficulty", StepDuration, StepDifficulty, false},
		{"Difficulty to Operation", StepDifficulty, StepOperation, false},
		{"Operation to InputMode", StepOperation, StepInputMode, false},
		{"InputMode to Ready", StepInputMode, StepReady, false},
		{"Ready completes", StepReady, StepReady, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewOnboarding()
			m.step = tt.startStep

			newModel, cmd := m.advance()

			if newModel.step != tt.expectedStep {
				t.Errorf("expected step %v, got %v", tt.expectedStep, newModel.step)
			}
			if tt.expectCmd && cmd == nil {
				t.Error("expected a command (completion), got nil")
			}
			if !tt.expectCmd && cmd != nil {
				t.Error("expected no command, got one")
			}
		})
	}
}

func TestOnboardingBack(t *testing.T) {
	tests := []struct {
		name         string
		startStep    OnboardingStep
		expectedStep OnboardingStep
	}{
		{"Welcome stays", StepWelcome, StepWelcome},
		{"Duration to Welcome", StepDuration, StepWelcome},
		{"Difficulty to Duration", StepDifficulty, StepDuration},
		{"Operation to Difficulty", StepOperation, StepDifficulty},
		{"InputMode to Operation", StepInputMode, StepOperation},
		{"Ready to InputMode", StepReady, StepInputMode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewOnboarding()
			m.step = tt.startStep

			m.back()

			if m.step != tt.expectedStep {
				t.Errorf("expected step %v, got %v", tt.expectedStep, m.step)
			}
		})
	}
}

func TestOnboardingMoveSelection(t *testing.T) {
	tests := []struct {
		name          string
		step          OnboardingStep
		startIndex    int
		delta         int
		expectedIndex int
	}{
		// Duration step (4 options: 0-3)
		{"Duration move down", StepDuration, 0, 1, 1},
		{"Duration move up", StepDuration, 2, -1, 1},
		{"Duration clamp at top", StepDuration, 0, -1, 0},
		{"Duration clamp at bottom", StepDuration, 3, 1, 3},

		// Difficulty step (5 options: 0-4)
		{"Difficulty move down", StepDifficulty, 0, 1, 1},
		{"Difficulty clamp at bottom", StepDifficulty, 4, 1, 4},

		// Operation step (5 options: 0-4)
		{"Operation move down", StepOperation, 0, 1, 1},
		{"Operation clamp at bottom", StepOperation, 4, 1, 4},

		// InputMode step (2 options: 0-1)
		{"InputMode move down", StepInputMode, 0, 1, 1},
		{"InputMode clamp at bottom", StepInputMode, 1, 1, 1},
		{"InputMode clamp at top", StepInputMode, 0, -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewOnboarding()
			m.step = tt.step
			m.setCurrentIndex(tt.startIndex)

			m.moveSelection(tt.delta)

			if m.currentIndex() != tt.expectedIndex {
				t.Errorf("expected index %d, got %d", tt.expectedIndex, m.currentIndex())
			}
		})
	}
}

func TestOnboardingComplete(t *testing.T) {
	m := NewOnboarding()
	m.durationIndex = 2   // 90 seconds
	m.difficultyIndex = 3 // Hard
	m.operationIndex = 2  // Multiplication
	m.inputModeIndex = 1  // Multiple Choice

	cmd := m.complete()
	if cmd == nil {
		t.Fatal("expected a command, got nil")
	}

	msg := cmd()
	completeMsg, ok := msg.(OnboardingCompleteMsg)
	if !ok {
		t.Fatalf("expected OnboardingCompleteMsg, got %T", msg)
	}

	if completeMsg.ModeID != modes.IDMultiplication {
		t.Errorf("expected ModeID %s, got %s", modes.IDMultiplication, completeMsg.ModeID)
	}
	if completeMsg.Difficulty != "Hard" {
		t.Errorf("expected Difficulty 'Hard', got %s", completeMsg.Difficulty)
	}
	if completeMsg.DurationMs != 90000 {
		t.Errorf("expected DurationMs 90000, got %d", completeMsg.DurationMs)
	}
	if completeMsg.InputMethod != "multiple_choice" {
		t.Errorf("expected InputMethod 'multiple_choice', got %s", completeMsg.InputMethod)
	}
}

func TestOnboardingCompleteTypingMode(t *testing.T) {
	m := NewOnboarding()
	m.inputModeIndex = 0 // Typing

	cmd := m.complete()
	msg := cmd()
	completeMsg := msg.(OnboardingCompleteMsg)

	if completeMsg.InputMethod != "typing" {
		t.Errorf("expected InputMethod 'typing', got %s", completeMsg.InputMethod)
	}
}

func TestOnboardingSkip(t *testing.T) {
	m := NewOnboarding()

	_, cmd := m.skip()
	if cmd == nil {
		t.Fatal("expected a command, got nil")
	}

	msg := cmd()
	_, ok := msg.(OnboardingSkipMsg)
	if !ok {
		t.Fatalf("expected OnboardingSkipMsg, got %T", msg)
	}
}

func TestOnboardingCurrentIndex(t *testing.T) {
	m := NewOnboarding()
	m.durationIndex = 2
	m.difficultyIndex = 3
	m.operationIndex = 1
	m.inputModeIndex = 1

	tests := []struct {
		step     OnboardingStep
		expected int
	}{
		{StepWelcome, 0},
		{StepDuration, 2},
		{StepDifficulty, 3},
		{StepOperation, 1},
		{StepInputMode, 1},
		{StepReady, 0}, // Ready has no selection
	}

	for _, tt := range tests {
		m.step = tt.step
		if m.currentIndex() != tt.expected {
			t.Errorf("step %v: expected index %d, got %d", tt.step, tt.expected, m.currentIndex())
		}
	}
}

func TestOnboardingSetCurrentIndex(t *testing.T) {
	m := NewOnboarding()

	// Set each step's index and verify
	m.step = StepDuration
	m.setCurrentIndex(3)
	if m.durationIndex != 3 {
		t.Errorf("expected durationIndex 3, got %d", m.durationIndex)
	}

	m.step = StepDifficulty
	m.setCurrentIndex(4)
	if m.difficultyIndex != 4 {
		t.Errorf("expected difficultyIndex 4, got %d", m.difficultyIndex)
	}

	m.step = StepOperation
	m.setCurrentIndex(2)
	if m.operationIndex != 2 {
		t.Errorf("expected operationIndex 2, got %d", m.operationIndex)
	}

	m.step = StepInputMode
	m.setCurrentIndex(1)
	if m.inputModeIndex != 1 {
		t.Errorf("expected inputModeIndex 1, got %d", m.inputModeIndex)
	}
}

func TestOnboardingMaxIndexForStep(t *testing.T) {
	m := NewOnboarding()

	tests := []struct {
		step     OnboardingStep
		expected int
	}{
		{StepWelcome, 0},
		{StepDuration, 3},   // 4 options
		{StepDifficulty, 4}, // 5 options
		{StepOperation, 4},  // 5 options
		{StepInputMode, 1},  // 2 options
		{StepReady, 0},      // No options
	}

	for _, tt := range tests {
		m.step = tt.step
		if m.maxIndexForStep() != tt.expected {
			t.Errorf("step %v: expected maxIndex %d, got %d", tt.step, tt.expected, m.maxIndexForStep())
		}
	}
}

func TestOnboardingUpdateKeyNavigation(t *testing.T) {
	m := NewOnboarding()
	m.step = StepDuration
	m.durationIndex = 1
	m.SetSize(80, 24) // Initialize viewport

	// Test down navigation
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if newModel.durationIndex != 2 {
		t.Errorf("expected durationIndex 2 after 'j', got %d", newModel.durationIndex)
	}

	// Test up navigation
	newModel, _ = newModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if newModel.durationIndex != 1 {
		t.Errorf("expected durationIndex 1 after 'k', got %d", newModel.durationIndex)
	}
}

func TestOnboardingUpdateStepNavigation(t *testing.T) {
	m := NewOnboarding()
	m.SetSize(80, 24)

	// Advance from Welcome
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if newModel.step != StepDuration {
		t.Errorf("expected StepDuration after Enter, got %v", newModel.step)
	}

	// Go back
	newModel, _ = newModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	if newModel.step != StepWelcome {
		t.Errorf("expected StepWelcome after 'b', got %v", newModel.step)
	}
}

func TestOnboardingSetSize(t *testing.T) {
	m := NewOnboarding()

	if m.viewportReady {
		t.Error("viewport should not be ready before SetSize")
	}

	m.SetSize(80, 24)

	if !m.viewportReady {
		t.Error("viewport should be ready after SetSize")
	}
	if m.width != 80 {
		t.Errorf("expected width 80, got %d", m.width)
	}
	if m.height != 24 {
		t.Errorf("expected height 24, got %d", m.height)
	}
}
