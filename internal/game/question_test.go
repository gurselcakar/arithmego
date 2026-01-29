package game

import "testing"

func TestQuestionCheckAnswer(t *testing.T) {
	q := Question{
		Operands:  []int{5, 3},
		Operation: &mockOp{},
		Answer:    8,
		Display:   "5 + 3",
	}

	tests := []struct {
		userAnswer    int
		expectCorrect bool
	}{
		{8, true},
		{7, false},
		{9, false},
		{0, false},
		{-8, false},
	}

	for _, tt := range tests {
		result := q.CheckAnswer(tt.userAnswer)
		if result.Correct != tt.expectCorrect {
			t.Errorf("CheckAnswer(%d).Correct = %v, want %v", tt.userAnswer, result.Correct, tt.expectCorrect)
		}
		if result.UserAnswer != tt.userAnswer {
			t.Errorf("CheckAnswer(%d).UserAnswer = %v, want %v", tt.userAnswer, result.UserAnswer, tt.userAnswer)
		}
		if result.CorrectAnswer != 8 {
			t.Errorf("CheckAnswer(%d).CorrectAnswer = %v, want 8", tt.userAnswer, result.CorrectAnswer)
		}
	}
}
