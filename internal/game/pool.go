package game

import "math/rand"

const defaultBatchSize = 50

// Generator produces questions for a specific game mode.
type Generator interface {
	Generate(diff Difficulty) *Question
	Label() string
}

// QuestionPool pre-generates and deduplicates questions for a session.
type QuestionPool struct {
	generator  Generator
	difficulty Difficulty
	questions  []*Question
	cursor     int
	seen       map[string]bool
}

// NewQuestionPool creates a pool that pre-generates a batch of questions.
func NewQuestionPool(g Generator, diff Difficulty) *QuestionPool {
	p := &QuestionPool{
		generator:  g,
		difficulty: diff,
		seen:       make(map[string]bool),
	}
	p.fill()
	return p
}

// Next returns the next question from the pool.
// Refills the pool when exhausted.
func (p *QuestionPool) Next() *Question {
	if p.cursor >= len(p.questions) {
		p.fill()
	}
	if len(p.questions) == 0 {
		// Truly exhausted — clear seen map and try again
		p.seen = make(map[string]bool)
		p.fill()
	}
	if len(p.questions) == 0 {
		// Generator is broken — return a fallback question
		return nil
	}
	q := p.questions[p.cursor]
	p.cursor++
	return q
}

// fill generates a batch of unique questions.
func (p *QuestionPool) fill() {
	p.questions = p.questions[:0]
	p.cursor = 0

	maxAttempts := defaultBatchSize * 3
	for i := 0; i < maxAttempts && len(p.questions) < defaultBatchSize; i++ {
		q := p.generator.Generate(p.difficulty)
		if q == nil {
			continue
		}
		if p.seen[q.Key] {
			continue
		}
		p.seen[q.Key] = true
		p.questions = append(p.questions, q)
	}

	// Shuffle for freshness
	rand.Shuffle(len(p.questions), func(i, j int) {
		p.questions[i], p.questions[j] = p.questions[j], p.questions[i]
	})
}
