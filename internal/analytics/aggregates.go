package analytics

import (
	"sort"
	"time"

	"github.com/gurselcakar/arithmego/internal/storage"
)

// OperationStats holds basic statistics for a single operation.
type OperationStats struct {
	Correct  int
	Total    int
	Accuracy float64
}

// ExtendedAggregates contains comprehensive computed statistics.
type ExtendedAggregates struct {
	// Basic counts
	TotalSessions     int
	TotalQuestions    int
	TotalCorrect      int
	OverallAccuracy   float64
	BestStreakEver    int
	AvgResponseTimeMs int64

	// Per-operation stats
	ByOperation map[string]OperationStats

	// Per-mode stats: map[modeName]sessionCount
	ByMode map[string]int

	// Cumulative points across all sessions
	TotalPoints int

	// Personal bests
	PersonalBests PersonalBests

	// Extended operation stats with response time and difficulty breakdown
	ByOperationExtended map[string]ExtendedOperationStats

	// Most recent session timestamp
	LastPlayedAt time.Time

	// Fastest response ever (across all questions)
	FastestResponseMs int64

	// Total time spent (sum of all response times)
	TotalResponseTimeMs int64
}

// PersonalBests tracks personal best achievements.
type PersonalBests struct {
	BestStreak     int     // Highest streak ever
	BestScore      int     // Highest single-session score
	BestAccuracy   float64 // Highest single-session accuracy (min 10 questions)
	FastestAvgTime int64   // Fastest average response time in a session (ms)
}

// ExtendedOperationStats holds detailed statistics for a single operation.
type ExtendedOperationStats struct {
	Correct           int
	Total             int
	Accuracy          float64
	AvgResponseTimeMs int64
	FastestTimeMs     int64

	// Breakdown by difficulty
	ByDifficulty map[string]DifficultyStats
}

// DifficultyStats holds statistics for a specific difficulty level.
type DifficultyStats struct {
	Correct  int
	Total    int
	Accuracy float64
}

// RecentMistake represents a recent wrong answer.
type RecentMistake struct {
	Question       string
	UserAnswer     int
	CorrectAnswer  int
	Operation      string
	SessionDate    time.Time
	ResponseTimeMs int64
}

// ComputeExtendedAggregates computes all aggregate statistics from all sessions.
func ComputeExtendedAggregates(stats *storage.Statistics) ExtendedAggregates {
	return ComputeFilteredAggregates(stats, AggregateFilter{})
}

// ComputeFilteredAggregates computes aggregates for sessions/questions matching the filter.
func ComputeFilteredAggregates(stats *storage.Statistics, filter AggregateFilter) ExtendedAggregates {
	agg := ExtendedAggregates{
		ByOperation:         make(map[string]OperationStats),
		ByMode:              make(map[string]int),
		ByOperationExtended: make(map[string]ExtendedOperationStats),
	}

	if len(stats.Sessions) == 0 {
		return agg
	}

	var totalResponseTime int64
	var questionsWithTime int

	// Track for extended operation stats
	opResponseTimes := make(map[string][]int64)

	for _, session := range stats.Sessions {
		// Check if session matches filter
		if !SessionMatchesFilter(session, filter) {
			continue
		}

		agg.TotalSessions++
		agg.TotalPoints += session.Score
		agg.ByMode[session.Mode]++

		// Track last played
		if session.Timestamp.After(agg.LastPlayedAt) {
			agg.LastPlayedAt = session.Timestamp
		}

		// Track personal bests
		if session.BestStreak > agg.PersonalBests.BestStreak {
			agg.PersonalBests.BestStreak = session.BestStreak
		}
		if session.Score > agg.PersonalBests.BestScore {
			agg.PersonalBests.BestScore = session.Score
		}

		// Track best accuracy (min 10 questions for meaningful stat)
		if session.QuestionsAttempted >= 10 {
			sessionAccuracy := float64(session.QuestionsCorrect) / float64(session.QuestionsAttempted) * 100
			if sessionAccuracy > agg.PersonalBests.BestAccuracy {
				agg.PersonalBests.BestAccuracy = sessionAccuracy
			}
		}

		// Track fastest avg time (only if session has valid response times)
		if session.AvgResponseTimeMs > 0 {
			if agg.PersonalBests.FastestAvgTime == 0 || session.AvgResponseTimeMs < agg.PersonalBests.FastestAvgTime {
				agg.PersonalBests.FastestAvgTime = session.AvgResponseTimeMs
			}
		}

		// Process questions
		for _, q := range session.Questions {
			// Check if question matches filter
			if !QuestionMatchesFilter(q, filter) {
				continue
			}

			// Basic operation stats
			opStats := agg.ByOperation[q.Operation]
			if !q.Skipped {
				opStats.Total++
				agg.TotalQuestions++
				if q.Correct {
					opStats.Correct++
					agg.TotalCorrect++
				}
			}
			agg.ByOperation[q.Operation] = opStats

			// Extended operation stats
			extOpStats := agg.ByOperationExtended[q.Operation]
			if extOpStats.ByDifficulty == nil {
				extOpStats.ByDifficulty = make(map[string]DifficultyStats)
			}

			if !q.Skipped {
				extOpStats.Total++
				if q.Correct {
					extOpStats.Correct++
				}

				// Track by difficulty
				diffStats := extOpStats.ByDifficulty[session.Difficulty]
				diffStats.Total++
				if q.Correct {
					diffStats.Correct++
				}
				extOpStats.ByDifficulty[session.Difficulty] = diffStats

				// Track response times
				if q.ResponseTimeMs > 0 {
					opResponseTimes[q.Operation] = append(opResponseTimes[q.Operation], q.ResponseTimeMs)

					// Track fastest response per operation
					if extOpStats.FastestTimeMs == 0 || q.ResponseTimeMs < extOpStats.FastestTimeMs {
						extOpStats.FastestTimeMs = q.ResponseTimeMs
					}

					// Track global fastest
					if agg.FastestResponseMs == 0 || q.ResponseTimeMs < agg.FastestResponseMs {
						agg.FastestResponseMs = q.ResponseTimeMs
					}
				}
			}
			agg.ByOperationExtended[q.Operation] = extOpStats

			// Track overall response time
			if q.ResponseTimeMs > 0 {
				totalResponseTime += q.ResponseTimeMs
				questionsWithTime++
				agg.TotalResponseTimeMs += q.ResponseTimeMs
			}
		}

		if session.BestStreak > agg.BestStreakEver {
			agg.BestStreakEver = session.BestStreak
		}
	}

	// Compute derived values
	if agg.TotalQuestions > 0 {
		agg.OverallAccuracy = float64(agg.TotalCorrect) / float64(agg.TotalQuestions) * 100
	}

	if questionsWithTime > 0 {
		agg.AvgResponseTimeMs = totalResponseTime / int64(questionsWithTime)
	}

	// Compute per-operation accuracy
	for op, opStats := range agg.ByOperation {
		if opStats.Total > 0 {
			opStats.Accuracy = float64(opStats.Correct) / float64(opStats.Total) * 100
		}
		agg.ByOperation[op] = opStats
	}

	// Compute extended operation stats
	for op, extOpStats := range agg.ByOperationExtended {
		if extOpStats.Total > 0 {
			extOpStats.Accuracy = float64(extOpStats.Correct) / float64(extOpStats.Total) * 100
		}

		// Compute avg response time for operation
		if times, ok := opResponseTimes[op]; ok && len(times) > 0 {
			var sum int64
			for _, t := range times {
				sum += t
			}
			extOpStats.AvgResponseTimeMs = sum / int64(len(times))
		}

		// Compute difficulty accuracy
		for diff, diffStats := range extOpStats.ByDifficulty {
			if diffStats.Total > 0 {
				diffStats.Accuracy = float64(diffStats.Correct) / float64(diffStats.Total) * 100
			}
			extOpStats.ByDifficulty[diff] = diffStats
		}

		agg.ByOperationExtended[op] = extOpStats
	}

	return agg
}

// GetRecentMistakes returns the last N wrong answers, optionally filtered by operation.
// Results are sorted by session date, most recent first.
func GetRecentMistakes(stats *storage.Statistics, operation string, limit int) []RecentMistake {
	var mistakes []RecentMistake

	// Iterate sessions in reverse order (most recent first)
	for i := len(stats.Sessions) - 1; i >= 0; i-- {
		session := stats.Sessions[i]

		for _, q := range session.Questions {
			// Skip correct answers and skipped questions
			if q.Correct || q.Skipped {
				continue
			}

			// Apply operation filter if specified
			if operation != "" && q.Operation != operation {
				continue
			}

			mistakes = append(mistakes, RecentMistake{
				Question:       q.Question,
				UserAnswer:     q.UserAnswer,
				CorrectAnswer:  q.CorrectAnswer,
				Operation:      q.Operation,
				SessionDate:    session.Timestamp,
				ResponseTimeMs: q.ResponseTimeMs,
			})

			// Stop once we have enough
			if len(mistakes) >= limit {
				return mistakes
			}
		}
	}

	return mistakes
}

// GetSessionsByFilter returns sessions matching the filter, sorted by timestamp (most recent first).
func GetSessionsByFilter(stats *storage.Statistics, filter AggregateFilter) []storage.SessionRecord {
	var sessions []storage.SessionRecord

	for _, session := range stats.Sessions {
		if SessionMatchesFilter(session, filter) {
			// Check if any questions match (for category/operation filters)
			if filter.Category != "" || filter.Operation != "" {
				hasMatch := false
				for _, q := range session.Questions {
					if QuestionMatchesFilter(q, filter) {
						hasMatch = true
						break
					}
				}
				if !hasMatch {
					continue
				}
			}
			sessions = append(sessions, session)
		}
	}

	// Sort by timestamp, most recent first
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Timestamp.After(sessions[j].Timestamp)
	})

	return sessions
}

// GetAllModes returns all unique mode names from sessions.
func GetAllModes(stats *storage.Statistics) []string {
	modeSet := make(map[string]bool)
	for _, session := range stats.Sessions {
		modeSet[session.Mode] = true
	}

	modes := make([]string, 0, len(modeSet))
	for mode := range modeSet {
		modes = append(modes, mode)
	}
	sort.Strings(modes)
	return modes
}

// GetAllOperations returns all unique operation names from sessions.
func GetAllOperations(stats *storage.Statistics) []string {
	opSet := make(map[string]bool)
	for _, session := range stats.Sessions {
		for _, q := range session.Questions {
			opSet[q.Operation] = true
		}
	}

	ops := make([]string, 0, len(opSet))
	for op := range opSet {
		ops = append(ops, op)
	}
	sort.Strings(ops)
	return ops
}

// GetOperationsByCategory returns operations that belong to a specific category.
// If category is empty, returns all operations.
func GetOperationsByCategory(stats *storage.Statistics, category string) []string {
	opSet := make(map[string]bool)
	for _, session := range stats.Sessions {
		for _, q := range session.Questions {
			if category == "" || GetOperationCategory(q.Operation) == category {
				opSet[q.Operation] = true
			}
		}
	}

	ops := make([]string, 0, len(opSet))
	for op := range opSet {
		ops = append(ops, op)
	}

	// Sort by category order, then alphabetically within category
	categoryOrder := map[string]int{"Basic": 0, "Power": 1, "Advanced": 2}
	sort.Slice(ops, func(i, j int) bool {
		catI := GetOperationCategory(ops[i])
		catJ := GetOperationCategory(ops[j])
		if categoryOrder[catI] != categoryOrder[catJ] {
			return categoryOrder[catI] < categoryOrder[catJ]
		}
		return ops[i] < ops[j]
	})

	return ops
}
