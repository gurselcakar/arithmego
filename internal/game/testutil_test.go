package game

// mockOp implements Operation for testing.
type mockOp struct {
	name string
}

func (m *mockOp) Name() string {
	if m.name == "" {
		return "Mock"
	}
	return m.name
}
func (m *mockOp) Symbol() string                                       { return "?" }
func (m *mockOp) Arity() Arity                                         { return Binary }
func (m *mockOp) Category() Category                                   { return CategoryBasic }
func (m *mockOp) Apply(operands []int) int                             { return operands[0] + operands[1] }
func (m *mockOp) ScoreDifficulty(operands []int, answer int) float64   { return 1.0 }
func (m *mockOp) Format(operands []int) string                         { return "mock" }
func (m *mockOp) Generate(diff Difficulty) Question {
	return Question{
		Operands:  []int{1, 2},
		Operation: m,
		Answer:    3,
		Display:   "1 + 2",
	}
}
