package gen

import (
	"math/rand"

	"github.com/gurselcakar/arithmego/internal/game"
)

type AnythingGoesGen struct{}

func (g *AnythingGoesGen) Label() string { return "Anything Goes" }

// singleOpGenerators are all standalone operation generators.
var singleOpGenerators = []game.Generator{
	&AdditionGen{}, &SubtractionGen{}, &MultiplicationGen{}, &DivisionGen{},
	&SquareGen{}, &CubeGen{}, &SquareRootGen{}, &CubeRootGen{},
	&PowerGen{}, &ModuloGen{}, &PercentageGen{}, &FactorialGen{},
}

var mixedGenerators = []game.Generator{
	&MixedBasicsGen{}, &MixedPowersGen{}, &MixedAdvancedGen{},
}

func (g *AnythingGoesGen) Generate(diff game.Difficulty) *game.Question {
	var picked game.Generator

	switch diff {
	case game.Beginner:
		picked = singleOpGenerators[rand.Intn(len(singleOpGenerators))]
	case game.Easy:
		if rand.Intn(10) < 7 {
			picked = singleOpGenerators[rand.Intn(len(singleOpGenerators))]
		} else {
			picked = mixedGenerators[rand.Intn(len(mixedGenerators))]
		}
	case game.Medium:
		r := rand.Intn(10)
		if r < 4 {
			picked = &MixedBasicsGen{}
		} else if r < 7 {
			picked = singleOpGenerators[rand.Intn(len(singleOpGenerators))]
		} else {
			mixed := []game.Generator{&MixedPowersGen{}, &MixedAdvancedGen{}}
			picked = mixed[rand.Intn(len(mixed))]
		}
	case game.Hard:
		r := rand.Intn(4)
		if r < 2 {
			picked = &MixedBasicsGen{}
		} else if r < 3 {
			picked = &MixedPowersGen{}
		} else {
			picked = &MixedAdvancedGen{}
		}
	case game.Expert:
		r := rand.Intn(10)
		if r < 4 {
			picked = &MixedBasicsGen{}
		} else if r < 7 {
			picked = &MixedPowersGen{}
		} else if r < 9 {
			picked = &MixedAdvancedGen{}
		} else {
			picked = singleOpGenerators[rand.Intn(len(singleOpGenerators))]
		}
	default:
		picked = singleOpGenerators[rand.Intn(len(singleOpGenerators))]
	}

	q := picked.Generate(diff)
	if q != nil {
		q.OpLabel = g.Label()
	}
	return q
}
