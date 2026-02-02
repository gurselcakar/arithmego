package statistics

import (
	"fmt"
	"time"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// StatisticsView represents the current view in the statistics screen.
type StatisticsView int

const (
	ViewDashboard StatisticsView = iota
	ViewOperations
	ViewOperationDetail
	ViewOperationReview // Review all mistakes for an operation
	ViewHistory
	ViewSessionDetail
	ViewSessionFullLog
	ViewTrends
)

// SessionDetailMode represents summary or full log mode.
type SessionDetailMode int

const (
	SessionModeSummary SessionDetailMode = iota
	SessionModeFullLog
)

// OperationSymbol returns the symbol for an operation name.
func OperationSymbol(name string) string {
	symbols := map[string]string{
		"Addition":       "+",
		"Subtraction":    "−",
		"Multiplication": "×",
		"Division":       "÷",
		"Square":         "²",
		"Cube":           "³",
		"Square Root":    "√",
		"Cube Root":      "∛",
		"Modulo":         "%",
		"Power":          "^",
		"Percentage":     "%",
		"Factorial":      "!",
	}
	if s, ok := symbols[name]; ok {
		return s
	}
	return "?"
}

// FormatRelativeTime formats a timestamp as relative time (e.g., "2 hours ago").
func FormatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 48*time.Hour:
		return "yesterday"
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	case diff < 14*24*time.Hour:
		return "1 week ago"
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / 24 / 7)
		return fmt.Sprintf("%d weeks ago", weeks)
	default:
		return t.Format("Jan 2, 2006")
	}
}

// FormatSessionDate formats a session timestamp for display.
// Returns "Today", "Yesterday", or formatted date.
func FormatSessionDate(t time.Time) string {
	now := time.Now()

	// Compare by calendar date to avoid DST issues with duration math
	nowYear, nowMonth, nowDay := now.Date()
	tYear, tMonth, tDay := t.Date()

	// Check if same calendar day
	if nowYear == tYear && nowMonth == tMonth && nowDay == tDay {
		return "Today"
	}

	// Check if yesterday by subtracting one day from today
	yesterday := now.AddDate(0, 0, -1)
	yYear, yMonth, yDay := yesterday.Date()
	if tYear == yYear && tMonth == yMonth && tDay == yDay {
		return "Yesterday"
	}

	// Check if within last 7 days
	weekAgo := now.AddDate(0, 0, -7)
	if t.After(weekAgo) {
		return t.Format("Mon, Jan 2")
	}

	return t.Format("Jan 2, 2006")
}

// FormatTime formats a time for session display.
func FormatTime(t time.Time) string {
	return t.Format("3:04 PM")
}

// FormatResponseTime formats response time in milliseconds to seconds with one decimal.
func FormatResponseTime(ms int64) string {
	if ms == 0 {
		return "--"
	}
	return fmt.Sprintf("%.1fs", float64(ms)/1000)
}

// FormatDuration formats duration in seconds to human-readable format.
func FormatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	mins := seconds / 60
	secs := seconds % 60
	if secs == 0 {
		return fmt.Sprintf("%dm", mins)
	}
	return fmt.Sprintf("%dm %ds", mins, secs)
}

// FormatThinkingTime formats total thinking time in milliseconds to a friendly format.
func FormatThinkingTime(ms int64) string {
	if ms == 0 {
		return "--"
	}
	totalSeconds := int(ms / 1000)
	if totalSeconds < 60 {
		return fmt.Sprintf("%ds", totalSeconds)
	}
	mins := totalSeconds / 60
	secs := totalSeconds % 60
	if secs == 0 {
		return fmt.Sprintf("%dm", mins)
	}
	return fmt.Sprintf("%dm %ds", mins, secs)
}

// FormatAccuracy formats accuracy with color based on value.
func FormatAccuracy(accuracy float64) string {
	text := fmt.Sprintf("%.0f%%", accuracy)
	if accuracy >= 80 {
		return styles.Correct.Render(text)
	} else if accuracy < 60 {
		return styles.Incorrect.Render(text)
	}
	return text
}

// FormatAccuracyPlain formats accuracy without color.
func FormatAccuracyPlain(accuracy float64) string {
	return fmt.Sprintf("%.0f%%", accuracy)
}

// MaxLen returns the maximum length of strings in a slice.
func MaxLen(strs []string) int {
	max := 0
	for _, s := range strs {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}
