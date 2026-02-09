package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

// All generators in this package implement game.Generator:
//   Generate(diff game.Difficulty) *game.Question
//   Label() string

// BuildQuestion creates a Question from an expression tree and label.
func BuildQuestion(e expr.Expr, label string) *game.Question {
	return &game.Question{
		Expression: e,
		Answer:     e.Eval(),
		Display:    e.Format(),
		Key:        e.Key(),
		OpLabel:    label,
	}
}

// TryGenerate attempts to generate a question using the pattern set for the given difficulty.
// Tries up to maxAttempts times, picking weighted patterns randomly.
func TryGenerate(patterns PatternSet, diff game.Difficulty, label string, maxAttempts int) *game.Question {
	wp, ok := patterns[diff]
	if !ok || len(wp) == 0 {
		return nil
	}

	for i := 0; i < maxAttempts; i++ {
		p := PickPattern(wp)
		e, valid := p(diff)
		if !valid || e == nil {
			continue
		}
		return BuildQuestion(e, label)
	}
	return nil
}
