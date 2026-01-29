package game

// Question represents a single arithmetic question.
type Question struct {
	Operands  []int
	Operation Operation
	Answer    int
	Display   string
}

// CheckAnswer validates a user's answer.
func (q Question) CheckAnswer(userAnswer int) AnswerResult {
	return AnswerResult{
		Correct:       userAnswer == q.Answer,
		UserAnswer:    userAnswer,
		CorrectAnswer: q.Answer,
	}
}

// AnswerResult represents the result of checking an answer.
type AnswerResult struct {
	Correct       bool
	UserAnswer    int
	CorrectAnswer int
}
