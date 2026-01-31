package game

import (
	"testing"
)

func TestGenerateChoices_ReturnsCorrectAnswer(t *testing.T) {
	tests := []struct {
		name       string
		answer     int
		difficulty Difficulty
	}{
		{"positive small", 5, Medium},
		{"positive medium", 42, Medium},
		{"positive large", 150, Medium},
		{"zero", 0, Medium},
		{"negative", -10, Medium},
		{"beginner", 12, Beginner},
		{"expert", 75, Expert},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			choices, correctIndex := GenerateChoices(tt.answer, tt.difficulty)

			if len(choices) != 4 {
				t.Errorf("expected 4 choices, got %d", len(choices))
			}

			if correctIndex < 0 || correctIndex > 3 {
				t.Errorf("correctIndex %d out of range [0,3]", correctIndex)
			}

			if choices[correctIndex] != tt.answer {
				t.Errorf("choices[%d] = %d, want %d", correctIndex, choices[correctIndex], tt.answer)
			}
		})
	}
}

func TestGenerateChoices_NoDuplicates(t *testing.T) {
	for i := 0; i < 100; i++ {
		choices, _ := GenerateChoices(25, Medium)

		seen := make(map[int]bool)
		for _, c := range choices {
			if seen[c] {
				t.Errorf("duplicate choice found: %d in %v", c, choices)
			}
			seen[c] = true
		}
	}
}

func TestGenerateChoices_NoNegativesForPositiveAnswer(t *testing.T) {
	for i := 0; i < 100; i++ {
		// Test with small positive answers where negatives would be likely
		choices, _ := GenerateChoices(3, Medium)

		for _, c := range choices {
			if c < 0 {
				t.Errorf("negative choice %d found for positive answer 3", c)
			}
		}
	}
}

func TestGenerateChoices_ShufflesPosition(t *testing.T) {
	// Run multiple times and verify the correct answer isn't always at the same index
	positions := make(map[int]int)

	for i := 0; i < 100; i++ {
		_, correctIndex := GenerateChoices(50, Medium)
		positions[correctIndex]++
	}

	// Should have at least 2 different positions used
	if len(positions) < 2 {
		t.Errorf("correct answer always at same position: %v", positions)
	}
}

func TestGenerateChoices_DistractorsNearAnswer(t *testing.T) {
	answer := 50
	for i := 0; i < 50; i++ {
		choices, _ := GenerateChoices(answer, Medium)

		for _, c := range choices {
			if c == answer {
				continue
			}
			// Distractors should be within reasonable range
			diff := abs(c - answer)
			if diff > 50 { // Should be within 100% of the answer
				t.Errorf("distractor %d too far from answer %d (diff: %d)", c, answer, diff)
			}
		}
	}
}

func TestGenerateChoices_ZeroAnswer(t *testing.T) {
	for i := 0; i < 50; i++ {
		choices, correctIndex := GenerateChoices(0, Medium)

		if choices[correctIndex] != 0 {
			t.Errorf("correct answer should be 0, got %d", choices[correctIndex])
		}

		// All choices should be valid (distinct)
		seen := make(map[int]bool)
		for _, c := range choices {
			if seen[c] {
				t.Errorf("duplicate choice in %v", choices)
			}
			seen[c] = true
		}
	}
}

func TestGenerateChoices_NegativeAnswer(t *testing.T) {
	for i := 0; i < 50; i++ {
		choices, correctIndex := GenerateChoices(-15, Medium)

		if choices[correctIndex] != -15 {
			t.Errorf("correct answer should be -15, got %d", choices[correctIndex])
		}

		// Should have 4 distinct choices
		if len(choices) != 4 {
			t.Errorf("expected 4 choices, got %d", len(choices))
		}
	}
}

func TestGenerateChoices_LargeAnswer(t *testing.T) {
	answer := 1000
	for i := 0; i < 50; i++ {
		choices, correctIndex := GenerateChoices(answer, Medium)

		if choices[correctIndex] != answer {
			t.Errorf("correct answer should be %d, got %d", answer, choices[correctIndex])
		}

		// Distractors should be proportionally offset
		for _, c := range choices {
			if c == answer {
				continue
			}
			diff := abs(c - answer)
			// For 1000, 30% offset = 300, so all should be within this range
			if diff > 500 {
				t.Errorf("distractor %d too far from answer %d", c, answer)
			}
		}
	}
}

func TestGenerateChoices_LargeNegativeAnswer(t *testing.T) {
	answer := -1000
	for i := 0; i < 50; i++ {
		choices, correctIndex := GenerateChoices(answer, Medium)

		if choices[correctIndex] != answer {
			t.Errorf("correct answer should be %d, got %d", answer, choices[correctIndex])
		}

		// Should have 4 distinct choices
		if len(choices) != 4 {
			t.Errorf("expected 4 choices, got %d", len(choices))
		}

		// Verify all choices are unique
		seen := make(map[int]bool)
		for _, c := range choices {
			if seen[c] {
				t.Errorf("duplicate choice %d in %v", c, choices)
			}
			seen[c] = true
		}

		// Distractors should be proportionally offset (within 50% for large answers)
		for _, c := range choices {
			if c == answer {
				continue
			}
			diff := abs(c - answer)
			if diff > 500 {
				t.Errorf("distractor %d too far from answer %d (diff: %d)", c, answer, diff)
			}
		}
	}
}

func TestGenerateChoices_NoNegativeDistractorsForZero(t *testing.T) {
	// Specifically tests that the guaranteed fallback loop respects the
	// "no negative distractors for non-negative answers" rule
	for i := 0; i < 100; i++ {
		choices, _ := GenerateChoices(0, Medium)

		for _, c := range choices {
			if c < 0 {
				t.Errorf("negative choice %d found for answer 0: %v", c, choices)
			}
		}
	}
}

func TestGenerateChoices_EdgeCases_Always4Choices(t *testing.T) {
	// Edge cases that previously could fail to generate 3 distractors
	edgeCases := []int{0, 1, -1, 2, -2}

	for _, answer := range edgeCases {
		for i := 0; i < 50; i++ {
			choices, correctIndex := GenerateChoices(answer, Medium)

			if len(choices) != 4 {
				t.Errorf("answer=%d: expected 4 choices, got %d: %v", answer, len(choices), choices)
			}

			if correctIndex < 0 || correctIndex > 3 {
				t.Errorf("answer=%d: correctIndex %d out of range", answer, correctIndex)
			}

			if choices[correctIndex] != answer {
				t.Errorf("answer=%d: choices[%d]=%d, want %d", answer, correctIndex, choices[correctIndex], answer)
			}

			// Verify all choices are unique
			seen := make(map[int]bool)
			for _, c := range choices {
				if seen[c] {
					t.Errorf("answer=%d: duplicate choice %d in %v", answer, c, choices)
				}
				seen[c] = true
			}
		}
	}
}
