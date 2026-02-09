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
//   - [Generator]: Interface for question generators (implemented in game/gen)
//   - [Question]: A generated problem with expression tree, display string, and answer
//   - [Difficulty]: Skill level affecting number ranges and complexity
//   - [Session]: Tracks state during active gameplay
//   - [QuestionPool]: Batch pre-generation with deduplication
//
// # Generators
//
// Generators are defined in the gen subpackage and implement the
// [Generator] interface. Each generator produces questions for a specific
// mode at various difficulty levels using weighted patterns.
//
// # Sessions
//
// A [Session] manages the flow of a timed game:
//
//	session := game.NewSession(generator, game.Medium, 60*time.Second)
//	session.Start()
//	session.SubmitAnswer(59) // returns true if correct
//
// # Scoring
//
// Points are awarded based on correctness and response time. Streaks of
// correct answers earn bonus multipliers. See [CalculateCorrectAnswer] for details.
//
// # Multiple Choice
//
// For multiple choice mode, use [GenerateChoices] to create plausible
// wrong answers alongside the correct one.
package game
