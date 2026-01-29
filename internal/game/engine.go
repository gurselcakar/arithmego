package game

import (
	"math/rand"
)

// GenerateQuestion creates a question using a random operation from the provided list.
func GenerateQuestion(ops []Operation, diff Difficulty) Question {
	if len(ops) == 0 {
		panic("GenerateQuestion: no operations provided")
	}
	op := ops[rand.Intn(len(ops))]
	return op.Generate(diff)
}

// GenerateQuestionForOperation creates a question using the specified operation.
func GenerateQuestionForOperation(op Operation, diff Difficulty) Question {
	return op.Generate(diff)
}
