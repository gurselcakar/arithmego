package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// QuestionRecord stores data for a single answered question.
type QuestionRecord struct {
	Question       string `json:"question"`
	Operation      string `json:"operation"`
	CorrectAnswer  int    `json:"correct_answer"`
	UserAnswer     int    `json:"user_answer"`
	Correct        bool   `json:"correct"`
	Skipped        bool   `json:"skipped"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	PointsEarned   int    `json:"points_earned"`
}

// SessionRecord stores data for a completed game session.
type SessionRecord struct {
	ID                 string           `json:"id"`
	Timestamp          time.Time        `json:"timestamp"`
	Mode               string           `json:"mode"`
	Difficulty         string           `json:"difficulty"`
	DurationSeconds    int              `json:"duration_seconds"`
	QuestionsAttempted int              `json:"questions_attempted"`
	QuestionsCorrect   int              `json:"questions_correct"`
	QuestionsWrong     int              `json:"questions_wrong"`
	QuestionsSkipped   int              `json:"questions_skipped"`
	Score              int              `json:"score"`
	BestStreak         int              `json:"best_streak"`
	AvgResponseTimeMs  int64            `json:"avg_response_time_ms"`
	Questions          []QuestionRecord `json:"questions"`
}

// Statistics holds all recorded sessions.
type Statistics struct {
	Sessions []SessionRecord `json:"sessions"`
}

// NewSessionRecord creates a new session record with a generated ID.
// Returns an error if mode or difficulty is empty, or if durationSeconds is negative.
func NewSessionRecord(mode, difficulty string, durationSeconds int) (SessionRecord, error) {
	if mode == "" {
		return SessionRecord{}, errors.New("mode cannot be empty")
	}
	if difficulty == "" {
		return SessionRecord{}, errors.New("difficulty cannot be empty")
	}
	if durationSeconds < 0 {
		return SessionRecord{}, errors.New("duration cannot be negative")
	}

	return SessionRecord{
		ID:              uuid.New().String(),
		Timestamp:       time.Now(),
		Mode:            mode,
		Difficulty:      difficulty,
		DurationSeconds: durationSeconds,
		Questions:       []QuestionRecord{},
	}, nil
}

// Load reads statistics from the JSON file.
// Returns empty statistics if the file doesn't exist.
func Load() (*Statistics, error) {
	path, err := StatisticsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Statistics{Sessions: []SessionRecord{}}, nil
		}
		return nil, err
	}

	var stats Statistics
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// Save writes statistics to the JSON file using atomic write.
// Writes to a temp file first, then renames to prevent corruption.
func Save(stats *Statistics) error {
	path, err := StatisticsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	// Write to temp file first for atomic operation
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "statistics-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	// Clean up temp file on any error
	shouldCleanup := true
	defer func() {
		if shouldCleanup {
			os.Remove(tmpPath)
		}
	}()

	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return err
	}

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	// Atomic rename (on most filesystems)
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	shouldCleanup = false // Prevent cleanup of successfully renamed file
	return nil
}

// AddSession appends a session and saves to disk.
func AddSession(record SessionRecord) error {
	stats, err := Load()
	if err != nil {
		return err
	}

	stats.Sessions = append(stats.Sessions, record)
	return Save(stats)
}

// Aggregates contains computed statistics across all sessions.
type Aggregates struct {
	TotalSessions     int
	TotalQuestions    int
	TotalCorrect      int
	OverallAccuracy   float64
	BestStreakEver    int
	AvgResponseTimeMs int64

	// Per-operation stats: map[operationName]accuracy
	ByOperation map[string]OperationStats

	// Per-mode stats: map[modeName]sessionCount
	ByMode map[string]int
}

// OperationStats holds statistics for a single operation.
type OperationStats struct {
	Correct  int
	Total    int
	Accuracy float64
}

// ComputeAggregates calculates aggregate statistics from all sessions.
func ComputeAggregates(stats *Statistics) Aggregates {
	agg := Aggregates{
		ByOperation: make(map[string]OperationStats),
		ByMode:      make(map[string]int),
	}

	if len(stats.Sessions) == 0 {
		return agg
	}

	var totalResponseTime int64
	var questionsWithTime int

	for _, session := range stats.Sessions {
		agg.TotalSessions++
		agg.TotalQuestions += session.QuestionsAttempted
		agg.TotalCorrect += session.QuestionsCorrect

		if session.BestStreak > agg.BestStreakEver {
			agg.BestStreakEver = session.BestStreak
		}

		agg.ByMode[session.Mode]++

		for _, q := range session.Questions {
			opStats := agg.ByOperation[q.Operation]
			if !q.Skipped {
				opStats.Total++
				if q.Correct {
					opStats.Correct++
				}
			}
			agg.ByOperation[q.Operation] = opStats

			if q.ResponseTimeMs > 0 {
				totalResponseTime += q.ResponseTimeMs
				questionsWithTime++
			}
		}
	}

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

	return agg
}
