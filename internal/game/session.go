package game

import "time"

// Session tracks the state of a single game session.
// Phase 5: Extend with per-question history (response times, answers given)
type Session struct {
	// Configuration (immutable after creation)
	Operation  Operation
	Difficulty Difficulty
	Duration   time.Duration

	// Timer state
	StartTime time.Time
	TimeLeft  time.Duration

	// Current question
	Current *Question

	// Results (Phase 5: Replace with detailed QuestionHistory slice)
	Correct   int
	Incorrect int
	Skipped   int
}

// NewSession creates a new game session with the given configuration.
func NewSession(op Operation, diff Difficulty, duration time.Duration) *Session {
	return &Session{
		Operation:  op,
		Difficulty: diff,
		Duration:   duration,
		TimeLeft:   duration,
	}
}

// Start begins the session timer and generates the first question.
func (s *Session) Start() {
	s.StartTime = time.Now()
	s.TimeLeft = s.Duration
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
	q := GenerateQuestionForOperation(s.Operation, s.Difficulty)
	s.Current = &q
}

// SubmitAnswer checks the user's answer and updates statistics.
// Returns true if the answer was correct.
func (s *Session) SubmitAnswer(answer int) bool {
	if s.Current == nil {
		return false
	}

	result := s.Current.CheckAnswer(answer)
	if result.Correct {
		s.Correct++
	} else {
		s.Incorrect++
	}

	s.NextQuestion()
	return result.Correct
}

// Skip skips the current question without answering.
func (s *Session) Skip() {
	s.Skipped++
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
