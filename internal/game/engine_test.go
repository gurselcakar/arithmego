package game

import (
	"testing"
)

func TestGenerateQuestion(t *testing.T) {
	ops := []Operation{&mockOp{name: "Mock1"}, &mockOp{name: "Mock2"}}

	// Should not panic and should return a question
	q := GenerateQuestion(ops, Beginner)
	if q.Answer != 3 {
		t.Errorf("GenerateQuestion returned wrong answer: %d", q.Answer)
	}
}

func TestGenerateQuestionPanicsWithEmptyOps(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("GenerateQuestion did not panic with empty operations")
		}
	}()

	GenerateQuestion([]Operation{}, Beginner)
}

func TestGenerateQuestionForOperation(t *testing.T) {
	op := &mockOp{name: "TestOp"}
	q := GenerateQuestionForOperation(op, Easy)

	if q.Operation != op {
		t.Error("GenerateQuestionForOperation returned wrong operation")
	}
	if q.Answer != 3 {
		t.Errorf("GenerateQuestionForOperation returned wrong answer: %d", q.Answer)
	}
}
