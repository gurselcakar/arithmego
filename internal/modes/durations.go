package modes

import "time"

// Duration represents a selectable game duration.
type Duration struct {
	Value time.Duration
	Label string
}

// AllowedDurations returns the list of selectable durations.
// Future: This could be extended to allow custom durations.
var AllowedDurations = []Duration{
	{Value: 30 * time.Second, Label: "30 seconds"},
	{Value: 60 * time.Second, Label: "1 minute"},
	{Value: 90 * time.Second, Label: "90 seconds"},
	{Value: 2 * time.Minute, Label: "2 minutes"},
}

// FindDurationIndex returns the index of the duration in AllowedDurations.
// Returns 0 if not found (defaults to first option).
func FindDurationIndex(d time.Duration) int {
	for i, dur := range AllowedDurations {
		if dur.Value == d {
			return i
		}
	}
	return 0
}
