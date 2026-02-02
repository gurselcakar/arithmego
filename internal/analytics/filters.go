package analytics

import (
	"time"

	"github.com/gurselcakar/arithmego/internal/storage"
)

// TimePeriod represents a time range filter for statistics.
type TimePeriod int

const (
	TimePeriodAllTime TimePeriod = iota
	TimePeriod7Days
	TimePeriod14Days
	TimePeriod30Days
	TimePeriod90Days
)

// Cutoff returns the time cutoff for this period.
// Returns zero time for AllTime (meaning no cutoff).
func (t TimePeriod) Cutoff() time.Time {
	now := time.Now()
	switch t {
	case TimePeriod7Days:
		return now.AddDate(0, 0, -7)
	case TimePeriod14Days:
		return now.AddDate(0, 0, -14)
	case TimePeriod30Days:
		return now.AddDate(0, 0, -30)
	case TimePeriod90Days:
		return now.AddDate(0, 0, -90)
	default:
		return time.Time{} // Zero time = all time
	}
}

// String returns the display name for this period.
func (t TimePeriod) String() string {
	switch t {
	case TimePeriod7Days:
		return "Last 7 Days"
	case TimePeriod14Days:
		return "Last 14 Days"
	case TimePeriod30Days:
		return "Last 30 Days"
	case TimePeriod90Days:
		return "Last 90 Days"
	default:
		return "All Time"
	}
}

// AllTimePeriods returns all available time period options.
func AllTimePeriods() []TimePeriod {
	return []TimePeriod{
		TimePeriodAllTime,
		TimePeriod7Days,
		TimePeriod14Days,
		TimePeriod30Days,
		TimePeriod90Days,
	}
}

// AggregateFilter specifies filters for computing aggregates.
type AggregateFilter struct {
	// Category filter: "", "Basic", "Power", "Advanced"
	Category string

	// Difficulty filter: "", "Beginner", "Easy", "Medium", "Hard", "Expert"
	Difficulty string

	// Time range filter
	TimePeriod TimePeriod

	// Mode filter: "" for all, or specific mode name
	Mode string

	// Operation filter: "" for all, or specific operation name
	Operation string
}

// IsEmpty returns true if no filters are active.
func (f AggregateFilter) IsEmpty() bool {
	return f.Category == "" &&
		f.Difficulty == "" &&
		f.TimePeriod == TimePeriodAllTime &&
		f.Mode == "" &&
		f.Operation == ""
}

// operationCategories maps operation names to their categories.
var operationCategories = map[string]string{
	// Basic
	"Addition":       "Basic",
	"Subtraction":    "Basic",
	"Multiplication": "Basic",
	"Division":       "Basic",

	// Power
	"Square":      "Power",
	"Cube":        "Power",
	"Square Root": "Power",
	"Cube Root":   "Power",

	// Advanced
	"Modulo":     "Advanced",
	"Power":      "Advanced",
	"Percentage": "Advanced",
	"Factorial":  "Advanced",
}

// GetOperationCategory returns the category for an operation name.
// Returns empty string if operation is unknown.
func GetOperationCategory(operation string) string {
	return operationCategories[operation]
}

// AllCategories returns all available category options.
func AllCategories() []string {
	return []string{"", "Basic", "Power", "Advanced"}
}

// CategoryDisplayName returns a display-friendly name for a category.
func CategoryDisplayName(category string) string {
	if category == "" {
		return "All Categories"
	}
	return category
}

// AllDifficulties returns all available difficulty options.
func AllDifficulties() []string {
	return []string{"", "Beginner", "Easy", "Medium", "Hard", "Expert"}
}

// DifficultyDisplayName returns a display-friendly name for a difficulty.
func DifficultyDisplayName(difficulty string) string {
	if difficulty == "" {
		return "All Difficulties"
	}
	return difficulty
}

// SessionMatchesFilter checks if a session matches the given filter.
func SessionMatchesFilter(s storage.SessionRecord, f AggregateFilter) bool {
	// Check time period
	if f.TimePeriod != TimePeriodAllTime {
		cutoff := f.TimePeriod.Cutoff()
		if s.Timestamp.Before(cutoff) {
			return false
		}
	}

	// Check difficulty
	if f.Difficulty != "" && s.Difficulty != f.Difficulty {
		return false
	}

	// Check mode
	if f.Mode != "" && s.Mode != f.Mode {
		return false
	}

	return true
}

// QuestionMatchesFilter checks if a question matches the given filter.
func QuestionMatchesFilter(q storage.QuestionRecord, f AggregateFilter) bool {
	// Check category
	if f.Category != "" {
		category := GetOperationCategory(q.Operation)
		if category != f.Category {
			return false
		}
	}

	// Check operation
	if f.Operation != "" && q.Operation != f.Operation {
		return false
	}

	return true
}
