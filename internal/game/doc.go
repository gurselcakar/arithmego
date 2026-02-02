// Package game implements the core arithmetic game logic.
//
// The package provides types and functions for generating questions,
// tracking game sessions, and computing scores. It is designed to be
// UI-agnostic and storage-agnostic.
//
// # Core Types
//
// The primary types are:
//
//   - [Operation]: Interface for arithmetic operations (addition, multiplication, etc.)
//   - [Question]: A generated problem with operands, expression, and answer
//   - [Difficulty]: Skill level affecting number ranges and complexity
//   - [Session]: Tracks state during active gameplay
//
// # Operations
//
// Operations are defined in the operations subpackage and implement the
// [Operation] interface. Each operation can generate questions at various
// difficulty levels and compute difficulty scores for adaptive gameplay.
//
//	ops := operations.BasicOperations() // +, -, ×, ÷
//	question := game.GenerateQuestion(ops, game.Medium)
//	fmt.Println(question.Expression) // "42 + 17"
//
// # Sessions
//
// A [Session] manages the flow of a timed game:
//
//	session := game.NewSession(ops, game.Medium, 60*time.Second)
//	session.Start()
//	question := session.CurrentQuestion()
//	session.SubmitAnswer(59) // returns true if correct
//	results := session.End()
//
// # Scoring
//
// Points are awarded based on correctness and response time. Streaks of
// correct answers earn bonus multipliers. See [ComputeScore] for details.
//
// # Multiple Choice
//
// For multiple choice mode, use [GenerateChoices] to create plausible
// wrong answers alongside the correct one.
package game
