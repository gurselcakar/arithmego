package game

import (
	"testing"
	"time"
)

func TestGetStreakTier(t *testing.T) {
	tests := []struct {
		streak int
		expect StreakTier
	}{
		{0, TierNone},
		{1, TierBuilding},
		{4, TierBuilding},
		{5, TierStreak},
		{9, TierStreak},
		{10, TierMax},
		{14, TierMax},
		{15, TierBlazing},
		{19, TierBlazing},
		{20, TierUnstoppable},
		{24, TierUnstoppable},
		{25, TierLegendary},
		{100, TierLegendary},
	}

	for _, tt := range tests {
		if got := GetStreakTier(tt.streak); got != tt.expect {
			t.Errorf("GetStreakTier(%d) = %v, want %v", tt.streak, got, tt.expect)
		}
	}
}

func TestStreakTierString(t *testing.T) {
	tests := []struct {
		tier   StreakTier
		expect string
	}{
		{TierNone, ""},
		{TierBuilding, ""},
		{TierStreak, "STREAK"},
		{TierMax, "MAX"},
		{TierBlazing, "BLAZING"},
		{TierUnstoppable, "UNSTOPPABLE"},
		{TierLegendary, "LEGENDARY"},
	}

	for _, tt := range tests {
		if got := tt.tier.String(); got != tt.expect {
			t.Errorf("StreakTier(%d).String() = %q, want %q", tt.tier, got, tt.expect)
		}
	}
}

func TestGetMilestoneAnnouncement(t *testing.T) {
	tests := []struct {
		streak int
		expect string
	}{
		{0, ""},
		{4, ""},
		{5, "×1.25"},
		{6, ""},
		{10, "×1.5"},
		{15, "×1.75"},
		{20, "×2.0 MAX"},
		{25, "LEGENDARY"},
		{30, ""},
	}

	for _, tt := range tests {
		if got := GetMilestoneAnnouncement(tt.streak); got != tt.expect {
			t.Errorf("GetMilestoneAnnouncement(%d) = %q, want %q", tt.streak, got, tt.expect)
		}
	}
}

func TestDifficultyMultiplier(t *testing.T) {
	tests := []struct {
		diff   Difficulty
		expect float64
	}{
		{Beginner, 0.5},
		{Easy, 0.75},
		{Medium, 1.0},
		{Hard, 1.5},
		{Expert, 2.0},
	}

	for _, tt := range tests {
		if got := DifficultyMultiplier(tt.diff); got != tt.expect {
			t.Errorf("DifficultyMultiplier(%s) = %v, want %v", tt.diff, got, tt.expect)
		}
	}
}

func TestTimeBonus(t *testing.T) {
	tests := []struct {
		name         string
		responseTime time.Duration
		expect       float64
	}{
		{"instant (0s)", 0, 1.5},
		{"instant (1s)", 1 * time.Second, 1.5},
		{"instant boundary (2s)", 2 * time.Second, 1.5},
		{"decay start (3s)", 3 * time.Second, 1.4375},
		{"midpoint (6s)", 6 * time.Second, 1.25},
		{"decay end (10s)", 10 * time.Second, 1.0},
		{"no bonus (15s)", 15 * time.Second, 1.0},
		{"no bonus (60s)", 60 * time.Second, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimeBonus(tt.responseTime)
			// Allow small floating point tolerance
			if diff := got - tt.expect; diff > 0.001 || diff < -0.001 {
				t.Errorf("TimeBonus(%v) = %v, want %v", tt.responseTime, got, tt.expect)
			}
		})
	}
}

func TestStreakBonus(t *testing.T) {
	tests := []struct {
		streak int
		expect float64
	}{
		{0, 1.0},
		{-1, 1.0},  // negative should be treated as 0
		{1, 1.0},   // streak 1-4: ×1.0
		{4, 1.0},   // still in first tier
		{5, 1.25},  // streak 5-9: ×1.25
		{9, 1.25},  // still in second tier
		{10, 1.5},  // streak 10-14: ×1.5
		{14, 1.5},  // still in third tier
		{15, 1.75}, // streak 15-19: ×1.75
		{19, 1.75}, // still in fourth tier
		{20, 2.0},  // streak 20+: ×2.0 (cap)
		{25, 2.0},  // still capped
		{100, 2.0}, // still capped
	}

	for _, tt := range tests {
		got := StreakBonus(tt.streak)
		if diff := got - tt.expect; diff > 0.001 || diff < -0.001 {
			t.Errorf("StreakBonus(%d) = %v, want %v", tt.streak, got, tt.expect)
		}
	}
}

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name         string
		difficulty   Difficulty
		responseTime time.Duration
		streak       int
		expect       int
	}{
		// Base case: Medium, no bonuses, no streak
		{"base case", Medium, 10 * time.Second, 0, 100},
		// Beginner halves the points
		{"beginner", Beginner, 10 * time.Second, 0, 50},
		// Expert doubles the points
		{"expert", Expert, 10 * time.Second, 0, 200},
		// Streak 5: ×1.25 multiplier: 100 * 1.0 * 1.0 * 1.25 = 125
		{"streak 5", Medium, 10 * time.Second, 5, 125},
		// Streak 10: ×1.5 multiplier: 100 * 1.0 * 1.0 * 1.5 = 150
		{"streak 10", Medium, 10 * time.Second, 10, 150},
		// Streak 15: ×1.75 multiplier: 100 * 1.0 * 1.0 * 1.75 = 175
		{"streak 15", Medium, 10 * time.Second, 15, 175},
		// Streak 20: ×2.0 multiplier (max): 100 * 1.0 * 1.0 * 2.0 = 200
		{"streak 20 max", Medium, 10 * time.Second, 20, 200},
		// Instant answer with max streak: 100 * 1.0 * 1.5 * 2.0 = 300
		{"max bonuses", Medium, 0, 20, 300},
		// Expert + instant + max streak: 100 * 2.0 * 1.5 * 2.0 = 600
		{"expert max", Expert, 0, 20, 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePoints(tt.difficulty, tt.responseTime, tt.streak)
			if got != tt.expect {
				t.Errorf("CalculatePoints(%s, %v, %d) = %d, want %d",
					tt.difficulty, tt.responseTime, tt.streak, got, tt.expect)
			}
		})
	}
}

func TestCalculateCorrectAnswer(t *testing.T) {
	tests := []struct {
		name            string
		difficulty      Difficulty
		responseTime    time.Duration
		currentStreak   int
		expectStreak    int
		expectMilestone bool
	}{
		{"first correct", Medium, 5 * time.Second, 0, 1, false},
		{"building streak", Medium, 5 * time.Second, 3, 4, false},
		{"hit 5 milestone", Medium, 5 * time.Second, 4, 5, true},   // ×1.25
		{"between milestones", Medium, 5 * time.Second, 6, 7, false},
		{"hit 10 milestone", Medium, 5 * time.Second, 9, 10, true}, // ×1.5
		{"hit 15 milestone", Medium, 5 * time.Second, 14, 15, true}, // ×1.75
		{"hit 20 milestone", Medium, 5 * time.Second, 19, 20, true},  // ×2.0 MAX
		{"between 20 and 25", Medium, 5 * time.Second, 22, 23, false},
		{"hit 25 legendary", Medium, 5 * time.Second, 24, 25, true}, // LEGENDARY
		{"after legendary", Medium, 5 * time.Second, 29, 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateCorrectAnswer(tt.difficulty, tt.responseTime, tt.currentStreak)
			if result.NewStreak != tt.expectStreak {
				t.Errorf("NewStreak = %d, want %d", result.NewStreak, tt.expectStreak)
			}
			if result.IsMilestone != tt.expectMilestone {
				t.Errorf("IsMilestone = %v, want %v", result.IsMilestone, tt.expectMilestone)
			}
			if result.Points <= 0 {
				t.Errorf("Points = %d, should be positive for correct answer", result.Points)
			}
		})
	}
}

func TestCalculateWrongAnswer(t *testing.T) {
	result := CalculateWrongAnswer()
	if result.Points != BasePointsWrong {
		t.Errorf("Points = %d, want %d", result.Points, BasePointsWrong)
	}
	if result.NewStreak != 0 {
		t.Errorf("NewStreak = %d, want 0", result.NewStreak)
	}
	if result.IsMilestone {
		t.Errorf("IsMilestone should be false for wrong answer")
	}
}

func TestCalculateSkip(t *testing.T) {
	result := CalculateSkip()
	if result.Points != BasePointsSkip {
		t.Errorf("Points = %d, want %d", result.Points, BasePointsSkip)
	}
	if result.NewStreak != 0 {
		t.Errorf("NewStreak = %d, want 0", result.NewStreak)
	}
	if result.IsMilestone {
		t.Errorf("IsMilestone should be false for skip")
	}
}
