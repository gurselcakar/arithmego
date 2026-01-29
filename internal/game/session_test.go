package game

import (
	"testing"
	"time"
)

// mockOperation is a simple operation for testing.
type mockOperation struct{}

func (m *mockOperation) Name() string           { return "Mock" }
func (m *mockOperation) Symbol() string         { return "?" }
func (m *mockOperation) Arity() Arity           { return Binary }
func (m *mockOperation) Category() Category     { return CategoryBasic }
func (m *mockOperation) Apply(operands []int) int { return operands[0] + operands[1] }
func (m *mockOperation) ScoreDifficulty(operands []int, answer int) float64 { return 5.0 }
func (m *mockOperation) Format(operands []int) string { return "1 ? 1" }
func (m *mockOperation) Generate(diff Difficulty) Question {
	return Question{
		Operands:  []int{1, 1},
		Operation: m,
		Answer:    2,
		Display:   "1 ? 1",
	}
}

func TestNewSession(t *testing.T) {
	op := &mockOperation{}
	ops := []Operation{op}
	s := NewSession(ops, Medium, 60*time.Second)

	if len(s.Operations) != 1 || s.Operations[0] != op {
		t.Error("session operations not set correctly")
	}
	if s.Difficulty != Medium {
		t.Error("session difficulty not set correctly")
	}
	if s.Duration != 60*time.Second {
		t.Error("session duration not set correctly")
	}
	if s.Correct != 0 || s.Incorrect != 0 || s.Skipped != 0 {
		t.Error("session counters should start at zero")
	}
}

func TestSessionStart(t *testing.T) {
	op := &mockOperation{}
	s := NewSession([]Operation{op}, Medium, 60*time.Second)
	s.Start()

	if s.StartTime.IsZero() {
		t.Error("start time should be set after Start()")
	}
	if s.Current == nil {
		t.Error("current question should be set after Start()")
	}
	if s.TimeLeft != 60*time.Second {
		t.Error("time left should equal duration after Start()")
	}
}

func TestSessionSubmitAnswer(t *testing.T) {
	op := &mockOperation{}
	s := NewSession([]Operation{op}, Medium, 60*time.Second)
	s.Start()

	// Correct answer
	correct := s.SubmitAnswer(2) // 1 + 1 = 2
	if !correct {
		t.Error("answer 2 should be correct")
	}
	if s.Correct != 1 {
		t.Errorf("correct count should be 1, got %d", s.Correct)
	}

	// Wrong answer
	correct = s.SubmitAnswer(999)
	if correct {
		t.Error("answer 999 should be incorrect")
	}
	if s.Incorrect != 1 {
		t.Errorf("incorrect count should be 1, got %d", s.Incorrect)
	}
}

func TestSessionSkip(t *testing.T) {
	op := &mockOperation{}
	s := NewSession([]Operation{op}, Medium, 60*time.Second)
	s.Start()

	s.Skip()
	if s.Skipped != 1 {
		t.Errorf("skipped count should be 1, got %d", s.Skipped)
	}
}

func TestSessionAccuracy(t *testing.T) {
	op := &mockOperation{}
	s := NewSession([]Operation{op}, Medium, 60*time.Second)
	s.Start()

	// No answers yet - should be 0
	if s.Accuracy() != 0 {
		t.Errorf("accuracy should be 0 with no answers, got %f", s.Accuracy())
	}

	// 1 correct, 1 incorrect = 50%
	s.SubmitAnswer(2)   // correct
	s.SubmitAnswer(999) // incorrect

	accuracy := s.Accuracy()
	if accuracy != 50 {
		t.Errorf("accuracy should be 50%%, got %f%%", accuracy)
	}
}

func TestSessionIsFinished(t *testing.T) {
	op := &mockOperation{}
	s := NewSession([]Operation{op}, Medium, 1*time.Millisecond)
	s.Start()

	if s.IsFinished() {
		t.Error("session should not be finished immediately")
	}

	// Wait and tick
	time.Sleep(5 * time.Millisecond)
	s.Tick()

	if !s.IsFinished() {
		t.Error("session should be finished after time expires")
	}
}
