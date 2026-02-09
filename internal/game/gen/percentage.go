package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type PercentageGen struct{}

func (g *PercentageGen) Label() string { return "Percentage" }

func (g *PercentageGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(percentagePatterns, diff, g.Label(), 100)
}

var percentagePatterns = PatternSet{
	game.Beginner: {
		{pctSingle, 10},
	},
	game.Easy: {
		{pctSingle, 10},
	},
	game.Medium: {
		{pctSingle, 10},
	},
	game.Hard: {
		{pctSingle, 10},
	},
	game.Expert: {
		{pctSingle, 10},
	},
}

func pctSingle(diff game.Difficulty) (expr.Expr, bool) {
	var pool []int
	switch diff {
	case game.Beginner, game.Easy:
		pool = PercentEasy
	case game.Medium:
		pool = PercentMedium
	case game.Hard:
		pool = PercentHard
	case game.Expert:
		pool = PercentExpert
	default:
		pool = PercentEasy
	}

	percent := PickFrom(pool)
	vr := PercentValueRanges[diff]
	value := RandomInRange(vr.Min, vr.Max)
	value = AlignToCleanDivision(value, percent, vr.Max)

	// Verify clean integer result
	if (percent*value)%100 != 0 {
		return nil, false
	}

	return &expr.BinOp{Op: expr.OpPct, Left: &expr.Num{Value: percent}, Right: &expr.Num{Value: value}}, true
}
