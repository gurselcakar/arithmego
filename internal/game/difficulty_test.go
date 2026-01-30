package game

import "testing"

func TestDifficultyScoreRange(t *testing.T) {
	tests := []struct {
		diff       Difficulty
		expectMin  float64
		expectMax  float64
	}{
		{Beginner, 1.0, 2.0},
		{Easy, 2.0, 4.0},
		{Medium, 4.0, 6.0},
		{Hard, 6.0, 8.0},
		{Expert, 8.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.diff.String(), func(t *testing.T) {
			min, max := tt.diff.ScoreRange()
			if min != tt.expectMin {
				t.Errorf("min = %v, want %v", min, tt.expectMin)
			}
			if max != tt.expectMax {
				t.Errorf("max = %v, want %v", max, tt.expectMax)
			}
		})
	}
}

func TestDifficultyAcceptsScore(t *testing.T) {
	tests := []struct {
		diff   Difficulty
		score  float64
		expect bool
	}{
		{Beginner, 1.0, true},
		{Beginner, 1.5, true},
		{Beginner, 2.0, true},
		{Beginner, 2.5, false},
		{Easy, 2.0, true},
		{Easy, 3.0, true},
		{Easy, 4.0, true},
		{Easy, 4.5, false},
		{Medium, 5.0, true},
		{Hard, 7.0, true},
		{Expert, 9.0, true},
		{Expert, 10.0, true},
	}

	for _, tt := range tests {
		result := tt.diff.AcceptsScore(tt.score)
		if result != tt.expect {
			t.Errorf("%s.AcceptsScore(%v) = %v, want %v", tt.diff, tt.score, result, tt.expect)
		}
	}
}

func TestDifficultyString(t *testing.T) {
	tests := []struct {
		diff   Difficulty
		expect string
	}{
		{Beginner, "Beginner"},
		{Easy, "Easy"},
		{Medium, "Medium"},
		{Hard, "Hard"},
		{Expert, "Expert"},
	}

	for _, tt := range tests {
		if got := tt.diff.String(); got != tt.expect {
			t.Errorf("%d.String() = %v, want %v", tt.diff, got, tt.expect)
		}
	}
}

func TestAllDifficulties(t *testing.T) {
	diffs := AllDifficulties()
	if len(diffs) != 5 {
		t.Errorf("AllDifficulties() returned %d difficulties, want 5", len(diffs))
	}

	expected := []Difficulty{Beginner, Easy, Medium, Hard, Expert}
	for i, d := range diffs {
		if d != expected[i] {
			t.Errorf("AllDifficulties()[%d] = %v, want %v", i, d, expected[i])
		}
	}
}

func TestParseDifficulty(t *testing.T) {
	tests := []struct {
		input  string
		expect Difficulty
	}{
		{"Beginner", Beginner},
		{"Easy", Easy},
		{"Medium", Medium},
		{"Hard", Hard},
		{"Expert", Expert},
		{"", Medium},         // default
		{"invalid", Medium},  // default
		{"beginner", Medium}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseDifficulty(tt.input); got != tt.expect {
				t.Errorf("ParseDifficulty(%q) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}
