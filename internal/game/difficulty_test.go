package game

import "testing"

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
