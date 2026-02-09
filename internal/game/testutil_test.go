package game

import "fmt"

// mockGen implements Generator for testing across test files.
type mockGen struct {
	name    string
	counter int
}

func (m *mockGen) Generate(diff Difficulty) *Question {
	m.counter++
	return &Question{
		Key:     fmt.Sprintf("%s-%d", m.Label(), m.counter),
		OpLabel: m.Label(),
		Answer:  3,
		Display: "1 + 2",
	}
}

func (m *mockGen) Label() string {
	if m.name == "" {
		return "Mock"
	}
	return m.name
}
