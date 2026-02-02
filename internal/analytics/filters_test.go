package analytics

import (
	"testing"
	"time"

	"github.com/gurselcakar/arithmego/internal/storage"
)

func TestTimePeriod_Cutoff(t *testing.T) {
	now := time.Now()

	tests := []struct {
		period      TimePeriod
		expectZero  bool
		expectDays  int
	}{
		{TimePeriodAllTime, true, 0},
		{TimePeriod7Days, false, 7},
		{TimePeriod14Days, false, 14},
		{TimePeriod30Days, false, 30},
		{TimePeriod90Days, false, 90},
	}

	for _, tt := range tests {
		t.Run(tt.period.String(), func(t *testing.T) {
			cutoff := tt.period.Cutoff()

			if tt.expectZero {
				if !cutoff.IsZero() {
					t.Errorf("Cutoff() = %v, want zero time", cutoff)
				}
			} else {
				if cutoff.IsZero() {
					t.Error("Cutoff() should not be zero")
				}

				expectedCutoff := now.AddDate(0, 0, -tt.expectDays)
				diff := cutoff.Sub(expectedCutoff)
				// Allow 1 second tolerance for test execution time
				if diff < -time.Second || diff > time.Second {
					t.Errorf("Cutoff() = %v, want ~%v (diff: %v)", cutoff, expectedCutoff, diff)
				}
			}
		})
	}
}

func TestTimePeriod_String(t *testing.T) {
	tests := []struct {
		period TimePeriod
		want   string
	}{
		{TimePeriodAllTime, "All Time"},
		{TimePeriod7Days, "Last 7 Days"},
		{TimePeriod14Days, "Last 14 Days"},
		{TimePeriod30Days, "Last 30 Days"},
		{TimePeriod90Days, "Last 90 Days"},
		{TimePeriod(99), "All Time"}, // Unknown defaults to All Time
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.period.String(); got != tt.want {
				t.Errorf("String() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestAllTimePeriods(t *testing.T) {
	periods := AllTimePeriods()

	if len(periods) != 5 {
		t.Errorf("len(AllTimePeriods()) = %d, want 5", len(periods))
	}

	// Check order
	expected := []TimePeriod{
		TimePeriodAllTime,
		TimePeriod7Days,
		TimePeriod14Days,
		TimePeriod30Days,
		TimePeriod90Days,
	}

	for i, p := range periods {
		if p != expected[i] {
			t.Errorf("periods[%d] = %v, want %v", i, p, expected[i])
		}
	}
}

func TestAggregateFilter_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		filter AggregateFilter
		want   bool
	}{
		{"empty filter", AggregateFilter{}, true},
		{"with category", AggregateFilter{Category: "Basic"}, false},
		{"with difficulty", AggregateFilter{Difficulty: "Hard"}, false},
		{"with time period", AggregateFilter{TimePeriod: TimePeriod7Days}, false},
		{"with mode", AggregateFilter{Mode: "Addition"}, false},
		{"with operation", AggregateFilter{Operation: "Addition"}, false},
		{"all time is empty", AggregateFilter{TimePeriod: TimePeriodAllTime}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOperationCategory(t *testing.T) {
	tests := []struct {
		operation string
		want      string
	}{
		// Basic
		{"Addition", "Basic"},
		{"Subtraction", "Basic"},
		{"Multiplication", "Basic"},
		{"Division", "Basic"},
		// Power
		{"Square", "Power"},
		{"Cube", "Power"},
		{"Square Root", "Power"},
		{"Cube Root", "Power"},
		// Advanced
		{"Modulo", "Advanced"},
		{"Power", "Advanced"},
		{"Percentage", "Advanced"},
		{"Factorial", "Advanced"},
		// Unknown
		{"Unknown", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.operation, func(t *testing.T) {
			if got := GetOperationCategory(tt.operation); got != tt.want {
				t.Errorf("GetOperationCategory(%s) = %s, want %s", tt.operation, got, tt.want)
			}
		})
	}
}

func TestAllCategories(t *testing.T) {
	categories := AllCategories()

	if len(categories) != 4 {
		t.Errorf("len(AllCategories()) = %d, want 4", len(categories))
	}

	// First should be empty string for "All"
	if categories[0] != "" {
		t.Errorf("categories[0] = %s, want empty string", categories[0])
	}

	// Check others exist
	expected := map[string]bool{"": true, "Basic": true, "Power": true, "Advanced": true}
	for _, c := range categories {
		if !expected[c] {
			t.Errorf("Unexpected category: %s", c)
		}
	}
}

func TestCategoryDisplayName(t *testing.T) {
	tests := []struct {
		category string
		want     string
	}{
		{"", "All Categories"},
		{"Basic", "Basic"},
		{"Power", "Power"},
		{"Advanced", "Advanced"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := CategoryDisplayName(tt.category); got != tt.want {
				t.Errorf("CategoryDisplayName(%s) = %s, want %s", tt.category, got, tt.want)
			}
		})
	}
}

func TestAllDifficulties(t *testing.T) {
	difficulties := AllDifficulties()

	if len(difficulties) != 6 {
		t.Errorf("len(AllDifficulties()) = %d, want 6", len(difficulties))
	}

	// First should be empty string for "All"
	if difficulties[0] != "" {
		t.Errorf("difficulties[0] = %s, want empty string", difficulties[0])
	}

	expected := []string{"", "Beginner", "Easy", "Medium", "Hard", "Expert"}
	for i, d := range difficulties {
		if d != expected[i] {
			t.Errorf("difficulties[%d] = %s, want %s", i, d, expected[i])
		}
	}
}

func TestDifficultyDisplayName(t *testing.T) {
	tests := []struct {
		difficulty string
		want       string
	}{
		{"", "All Difficulties"},
		{"Easy", "Easy"},
		{"Hard", "Hard"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := DifficultyDisplayName(tt.difficulty); got != tt.want {
				t.Errorf("DifficultyDisplayName(%s) = %s, want %s", tt.difficulty, got, tt.want)
			}
		})
	}
}

func TestSessionMatchesFilter(t *testing.T) {
	now := time.Now()

	session := storage.SessionRecord{
		Timestamp:  now.AddDate(0, 0, -5), // 5 days ago
		Difficulty: "Medium",
		Mode:       "Addition",
	}

	tests := []struct {
		name   string
		filter AggregateFilter
		want   bool
	}{
		{"empty filter", AggregateFilter{}, true},
		{"matching difficulty", AggregateFilter{Difficulty: "Medium"}, true},
		{"non-matching difficulty", AggregateFilter{Difficulty: "Hard"}, false},
		{"matching mode", AggregateFilter{Mode: "Addition"}, true},
		{"non-matching mode", AggregateFilter{Mode: "Multiplication"}, false},
		{"within time period", AggregateFilter{TimePeriod: TimePeriod7Days}, true},
		{"outside time period", AggregateFilter{TimePeriod: TimePeriod7Days}, true}, // 5 days ago is within 7 days
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SessionMatchesFilter(session, tt.filter); got != tt.want {
				t.Errorf("SessionMatchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionMatchesFilter_TimePeriod(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		daysAgo    int
		timePeriod TimePeriod
		want       bool
	}{
		{"3 days ago, 7 day filter", 3, TimePeriod7Days, true},
		{"10 days ago, 7 day filter", 10, TimePeriod7Days, false},
		{"10 days ago, 14 day filter", 10, TimePeriod14Days, true},
		{"20 days ago, 14 day filter", 20, TimePeriod14Days, false},
		{"100 days ago, all time", 100, TimePeriodAllTime, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := storage.SessionRecord{
				Timestamp: now.AddDate(0, 0, -tt.daysAgo),
			}
			filter := AggregateFilter{TimePeriod: tt.timePeriod}

			if got := SessionMatchesFilter(session, filter); got != tt.want {
				t.Errorf("SessionMatchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionMatchesFilter(t *testing.T) {
	question := storage.QuestionRecord{
		Operation: "Addition",
	}

	tests := []struct {
		name   string
		filter AggregateFilter
		want   bool
	}{
		{"empty filter", AggregateFilter{}, true},
		{"matching operation", AggregateFilter{Operation: "Addition"}, true},
		{"non-matching operation", AggregateFilter{Operation: "Multiplication"}, false},
		{"matching category", AggregateFilter{Category: "Basic"}, true},
		{"non-matching category", AggregateFilter{Category: "Power"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QuestionMatchesFilter(question, tt.filter); got != tt.want {
				t.Errorf("QuestionMatchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionMatchesFilter_Categories(t *testing.T) {
	tests := []struct {
		operation string
		category  string
		want      bool
	}{
		{"Addition", "Basic", true},
		{"Addition", "Power", false},
		{"Square", "Power", true},
		{"Square", "Basic", false},
		{"Modulo", "Advanced", true},
		{"Modulo", "Basic", false},
	}

	for _, tt := range tests {
		t.Run(tt.operation+"_"+tt.category, func(t *testing.T) {
			question := storage.QuestionRecord{Operation: tt.operation}
			filter := AggregateFilter{Category: tt.category}

			if got := QuestionMatchesFilter(question, filter); got != tt.want {
				t.Errorf("QuestionMatchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
