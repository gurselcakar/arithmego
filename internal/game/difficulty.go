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
