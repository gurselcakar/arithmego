package game

// Difficulty represents the difficulty tier for question generation.
type Difficulty int

const (
	Beginner Difficulty = iota
	Easy
	Medium
	Hard
	Expert
)

// ScoreRange returns the min and max difficulty scores for this tier.
func (d Difficulty) ScoreRange() (min, max float64) {
	switch d {
	case Beginner:
		return 1.0, 2.0
	case Easy:
		return 2.0, 4.0
	case Medium:
		return 4.0, 6.0
	case Hard:
		return 6.0, 8.0
	case Expert:
		return 8.0, 10.0
	default:
		return 1.0, 2.0
	}
}

// AcceptsScore returns true if the score falls within this difficulty's range.
func (d Difficulty) AcceptsScore(score float64) bool {
	min, max := d.ScoreRange()
	return score >= min && score <= max
}

// String returns the display name.
func (d Difficulty) String() string {
	switch d {
	case Beginner:
		return "Beginner"
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	case Expert:
		return "Expert"
	default:
		return "Unknown"
	}
}

// AllDifficulties returns all difficulty tiers.
func AllDifficulties() []Difficulty {
	return []Difficulty{Beginner, Easy, Medium, Hard, Expert}
}

// ParseDifficulty converts a string to a Difficulty.
// Returns Medium as the default for unrecognized strings.
func ParseDifficulty(s string) Difficulty {
	switch s {
	case "Beginner":
		return Beginner
	case "Easy":
		return Easy
	case "Medium":
		return Medium
	case "Hard":
		return Hard
	case "Expert":
		return Expert
	default:
		return Medium
	}
}
