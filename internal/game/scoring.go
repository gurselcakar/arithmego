package game

import "time"

// Scoring constants
const (
	BasePointsCorrect = 100
	BasePointsWrong   = -25
	BasePointsSkip    = 0

	MaxTimeBonus     = 1.5
	TimeBonusFloor   = 1.0
	InstantThreshold = 2 * time.Second
	TimeBonusDecay   = 10 * time.Second

	// Streak multiplier: increases by 0.25 every 5 correct answers
	MaxStreakBonus      = 2.0
	StreakBonusStep     = 0.25
	StreakMilestoneSize = 5
)

// StreakTier represents the visual tier of a streak.
type StreakTier int

const (
	TierNone        StreakTier = iota // 0
	TierBuilding                      // 1-4
	TierStreak                        // 5-9
	TierMax                           // 10-14
	TierBlazing                       // 15-19
	TierUnstoppable                   // 20-24
	TierLegendary                     // 25+
)

// String returns the display name for the tier (used for visual styling).
func (t StreakTier) String() string {
	switch t {
	case TierStreak:
		return "STREAK"
	case TierMax:
		return "MAX"
	case TierBlazing:
		return "BLAZING"
	case TierUnstoppable:
		return "UNSTOPPABLE"
	case TierLegendary:
		return "LEGENDARY"
	default:
		return ""
	}
}

// GetMilestoneAnnouncement returns the announcement text if this streak
// is a milestone, or empty string if not.
// Milestones occur when the multiplier increases (5, 10, 15, 20) or at Legendary tier (25).
func GetMilestoneAnnouncement(streak int) string {
	switch streak {
	case 5:
		return "×1.25"
	case 10:
		return "×1.5"
	case 15:
		return "×1.75"
	case 20:
		return "×2.0 MAX"
	case 25:
		return "LEGENDARY"
	default:
		return ""
	}
}

// GetStreakTier returns the tier for a given streak count.
func GetStreakTier(streak int) StreakTier {
	switch {
	case streak == 0:
		return TierNone
	case streak < 5:
		return TierBuilding
	case streak < 10:
		return TierStreak
	case streak < 15:
		return TierMax
	case streak < 20:
		return TierBlazing
	case streak < 25:
		return TierUnstoppable
	default:
		return TierLegendary
	}
}

// DifficultyMultiplier returns the point multiplier for a difficulty level.
func DifficultyMultiplier(d Difficulty) float64 {
	switch d {
	case Beginner:
		return 0.5
	case Easy:
		return 0.75
	case Medium:
		return 1.0
	case Hard:
		return 1.5
	case Expert:
		return 2.0
	default:
		return 1.0
	}
}

// TimeBonus calculates the time bonus multiplier based on response time.
// Returns 1.5x for instant answers (< 2s), linear decay to 1.0x at 10s.
func TimeBonus(responseTime time.Duration) float64 {
	if responseTime < InstantThreshold {
		return MaxTimeBonus
	}
	if responseTime >= TimeBonusDecay {
		return TimeBonusFloor
	}

	// Linear decay from 1.5 to 1.0 between 2s and 10s
	// Formula: 1.5 - (0.5 × (responseTime - 2) / 8)
	elapsed := responseTime - InstantThreshold
	decayWindow := TimeBonusDecay - InstantThreshold
	decay := float64(elapsed) / float64(decayWindow) * (MaxTimeBonus - TimeBonusFloor)
	return MaxTimeBonus - decay
}

// StreakBonus calculates the streak multiplier.
// Increases by 0.25 every 5 correct answers, capped at 2.0.
// Streak 0-4: ×1.0, 5-9: ×1.25, 10-14: ×1.5, 15-19: ×1.75, 20+: ×2.0
func StreakBonus(streak int) float64 {
	if streak <= 0 {
		return 1.0
	}
	milestones := streak / StreakMilestoneSize
	bonus := 1.0 + float64(milestones)*StreakBonusStep
	if bonus > MaxStreakBonus {
		return MaxStreakBonus
	}
	return bonus
}

// CalculatePoints calculates points for a correct answer.
// Returns the total points after applying all multipliers.
func CalculatePoints(difficulty Difficulty, responseTime time.Duration, streak int) int {
	base := float64(BasePointsCorrect)
	diffMult := DifficultyMultiplier(difficulty)
	timeMult := TimeBonus(responseTime)
	streakMult := StreakBonus(streak)

	points := base * diffMult * timeMult * streakMult
	return int(points)
}

// ScoreResult contains the result of a scoring calculation.
type ScoreResult struct {
	Points      int
	NewStreak   int
	OldTier     StreakTier
	NewTier     StreakTier
	IsMilestone bool
}

// CalculateCorrectAnswer calculates the score for a correct answer.
// Points are calculated based on the current streak (before this answer),
// then the streak is incremented for the next question.
func CalculateCorrectAnswer(difficulty Difficulty, responseTime time.Duration, currentStreak int) ScoreResult {
	oldTier := GetStreakTier(currentStreak)
	newStreak := currentStreak + 1
	newTier := GetStreakTier(newStreak)

	// Use currentStreak for points - you earn based on streak you had, not streak you'll have
	points := CalculatePoints(difficulty, responseTime, currentStreak)

	// Milestone when hitting 5, 10, 15, or 20 streak (multiplier increases)
	isMilestone := GetMilestoneAnnouncement(newStreak) != ""

	return ScoreResult{
		Points:      points,
		NewStreak:   newStreak,
		OldTier:     oldTier,
		NewTier:     newTier,
		IsMilestone: isMilestone,
	}
}

// CalculateWrongAnswer returns the penalty for a wrong answer.
func CalculateWrongAnswer() ScoreResult {
	return ScoreResult{
		Points:      BasePointsWrong,
		NewStreak:   0,
		OldTier:     TierNone,
		NewTier:     TierNone,
		IsMilestone: false,
	}
}

// CalculateSkip returns the result for skipping a question.
func CalculateSkip() ScoreResult {
	return ScoreResult{
		Points:      BasePointsSkip,
		NewStreak:   0,
		OldTier:     TierNone,
		NewTier:     TierNone,
		IsMilestone: false,
	}
}
