package game

import "time"

// QuestionHistory stores data for a single answered question.
type QuestionHistory struct {
	Question      string
	Operation     string
	CorrectAnswer int
	UserAnswer    int
	Correct       bool
	Skipped       bool
	ResponseTime  time.Duration
	PointsEarned  int
}

// Session tracks the state of a single game session.
type Session struct {
	// Configuration (immutable after creation)
	Operations []Operation
	Difficulty Difficulty
	Duration   time.Duration

	// Timer state
	StartTime time.Time
	TimeLeft  time.Duration

	// Current question
	Current       *Question
	QuestionStart time.Time // When current question was shown

	// Results
	Correct   int
	Incorrect int
	Skipped   int

	// Per-question history
	History []QuestionHistory

	// Scoring
	Score      int          // Running total
	Streak     int          // Current consecutive correct answers
	BestStreak int          // Session high streak
	LastResult *ScoreResult // Result of last answer (for UI feedback)
}

// NewSession creates a new game session with the given configuration.
// Accepts one or more operations; questions will randomly use any of them.
func NewSession(ops []Operation, diff Difficulty, duration time.Duration) *Session {
	return &Session{
		Operations: ops,
		Difficulty: diff,
		Duration:   duration,
		TimeLeft:   duration,
		History:    []QuestionHistory{},
	}
}

// Start begins the session timer and generates the first question.
func (s *Session) Start() {
	s.StartTime = time.Now()
	s.TimeLeft = s.Duration
	s.Score = 0
	s.Streak = 0
	s.BestStreak = 0
	s.LastResult = nil
	s.NextQuestion()
}

// Tick updates the time remaining. Called each second.
func (s *Session) Tick() {
	elapsed := time.Since(s.StartTime)
	s.TimeLeft = s.Duration - elapsed
	if s.TimeLeft < 0 {
		s.TimeLeft = 0
	}
}

// IsFinished returns true if the session time has expired.
func (s *Session) IsFinished() bool {
	return s.TimeLeft <= 0
}

// NextQuestion generates and sets a new question.
func (s *Session) NextQuestion() {
	q := GenerateQuestion(s.Operations, s.Difficulty)
	s.Current = &q
	s.QuestionStart = time.Now()
}

// SubmitAnswer checks the user's answer and updates statistics.
// Returns true if the answer was correct.
func (s *Session) SubmitAnswer(answer int) bool {
	if s.Current == nil {
		return false
	}

	result := s.Current.CheckAnswer(answer)
	responseTime := time.Since(s.QuestionStart)

	var points int
	if result.Correct {
		s.Correct++
		scoreResult := CalculateCorrectAnswer(s.Difficulty, responseTime, s.Streak)
		s.Score += scoreResult.Points
		s.Streak = scoreResult.NewStreak
		if s.Streak > s.BestStreak {
			s.BestStreak = s.Streak
		}
		s.LastResult = &scoreResult
		points = scoreResult.Points
	} else {
		s.Incorrect++
		scoreResult := CalculateWrongAnswer()
		s.Score += scoreResult.Points
		s.Streak = 0
		s.LastResult = &scoreResult
		points = scoreResult.Points
	}

	// Record question history
	s.History = append(s.History, QuestionHistory{
		Question:      s.Current.Display,
		Operation:     s.Current.Operation.Name(),
		CorrectAnswer: s.Current.Answer,
		UserAnswer:    answer,
		Correct:       result.Correct,
		Skipped:       false,
		ResponseTime:  responseTime,
		PointsEarned:  points,
	})

	s.NextQuestion()
	return result.Correct
}

// Skip skips the current question without answering.
func (s *Session) Skip() {
	// Record skipped question before moving to next
	if s.Current != nil {
		s.History = append(s.History, QuestionHistory{
			Question:      s.Current.Display,
			Operation:     s.Current.Operation.Name(),
			CorrectAnswer: s.Current.Answer,
			UserAnswer:    0,
			Correct:       false,
			Skipped:       true,
			ResponseTime:  time.Since(s.QuestionStart),
			PointsEarned:  0,
		})
	}

	s.Skipped++
	scoreResult := CalculateSkip()
	s.Score += scoreResult.Points
	s.Streak = 0
	s.LastResult = &scoreResult
	s.NextQuestion()
}

// Resume restarts the session timer after a pause.
// It adjusts StartTime so that elapsed time calculations remain correct.
func (s *Session) Resume() {
	s.StartTime = time.Now().Add(-(s.Duration - s.TimeLeft))
}

// Accuracy returns the accuracy percentage (0-100).
func (s *Session) Accuracy() float64 {
	total := s.TotalAnswered()
	if total == 0 {
		return 0
	}
	return float64(s.Correct) / float64(total) * 100
}

// TotalAnswered returns the total number of questions answered (correct + incorrect).
func (s *Session) TotalAnswered() int {
	return s.Correct + s.Incorrect
}

// StreakTier returns the current streak tier.
func (s *Session) StreakTier() StreakTier {
	return GetStreakTier(s.Streak)
}

// Multiplier returns the current streak bonus multiplier.
func (s *Session) Multiplier() float64 {
	return StreakBonus(s.Streak)
}

// ClearLastResult clears the last result after UI has displayed it.
func (s *Session) ClearLastResult() {
	s.LastResult = nil
}

// AvgResponseTime returns the average response time for answered questions.
func (s *Session) AvgResponseTime() time.Duration {
	if len(s.History) == 0 {
		return 0
	}

	var total time.Duration
	var count int
	for _, h := range s.History {
		if !h.Skipped {
			total += h.ResponseTime
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}

// FastestResponseTime returns the fastest response time for correct answers.
// Returns 0 if no correct answers exist.
func (s *Session) FastestResponseTime() time.Duration {
	var fastest time.Duration
	for _, h := range s.History {
		if h.Correct && !h.Skipped {
			if fastest == 0 || h.ResponseTime < fastest {
				fastest = h.ResponseTime
			}
		}
	}
	return fastest
}
