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
//   - [Question]: A generated problem with expression tree, answer, and display string
//   - [QuestionPool]: Batch pre-generation with session-level deduplication
//   - [Session]: Tracks state during active gameplay with timer, score, and streak
//   - [Difficulty]: Skill level (Beginner, Easy, Medium, Hard, Expert)
//   - [Category]: Mode grouping (basic, power, advanced)
//
// # Expression Trees
//
// Questions use expression trees from the expr subpackage. Supported node types:
// Num (literals), BinOp (+, -, ×, ÷), Paren, UnaryPrefix (√), UnarySuffix (!), Pow.
// Each expression provides Format() for display and Eval() for computation.
//
// # Generators
//
// Generators are defined in the gen subpackage and implement the [Generator] interface:
//
//	type Generator interface {
//	    Generate(diff Difficulty) *Question
//	    Label() string
//	}
//
// Each generator produces questions for a specific mode using weighted patterns
// that vary by difficulty. Generators self-register via init() in gen/registry.go.
//
// # Question Pool
//
// [QuestionPool] pre-generates batches of 50 questions with automatic deduplication
// based on Question.Key. When exhausted, it refills with fresh questions while
// maintaining the session-level dedup cache to prevent repeats.
//
// # Sessions
//
// A [Session] manages the flow of a timed game:
//
//	session := game.NewSession(generator, game.Medium, 60*time.Second)
//	session.Start()
//	correct := session.SubmitAnswer(42)
//	session.Skip()
//
// Session tracks correct/incorrect/skipped counts, score, current streak,
// best streak, and per-question history with response times and points earned.
//
// # Scoring
//
// Points are calculated via [CalculateCorrectAnswer] based on:
//   - Base points: 100
//   - Difficulty multiplier: 0.5x (Beginner) to 2.0x (Expert)
//   - Time bonus: 1.5x for <2s, linear decay to 1.0x at 10s
//   - Streak multiplier: +0.25 per 5 correct answers, capped at 2.0x
//
// Wrong answers deduct 25 points; skips award 0 points. Both reset streak to 0.
// Streaks trigger visual tiers and milestone announcements at 5, 10, 15, 20, 25.
//
// # Multiple Choice
//
// [GenerateChoices] creates 4 shuffled options with 3 distractors based on
// answer magnitude and difficulty. Returns choices and the correct index.
package game
