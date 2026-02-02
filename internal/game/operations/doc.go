// Package operations implements arithmetic operations for the game.
//
// Each operation (addition, multiplication, square root, etc.) implements
// the [game.Operation] interface and is registered at package initialization.
// Operations are grouped into categories: Basic, Power, and Advanced.
//
// # Registry
//
// Operations self-register via init() functions. Use the registry functions
// to access them:
//
//	op, ok := operations.Get("Addition")     // by name
//	ops := operations.All()                   // all operations
//	ops := operations.BasicOperations()       // +, -, ×, ÷
//	ops := operations.PowerOperations()       // squares, cubes, roots
//	ops := operations.AdvancedOperations()    // modulo, power, percentage, factorial
//	ops := operations.ByCategory(game.CategoryBasic)
//
// # Question Generation
//
// Each operation's Generate method produces questions appropriate for the
// requested difficulty. The generation system uses a two-phase approach:
//
//  1. Standard generation with difficulty-appropriate ranges
//  2. Relaxed fallback if standard generation can't meet difficulty targets
//
// This ensures questions are always generated while maintaining quality.
//
// # Difficulty Scoring
//
// Each operation implements ScoreDifficulty to estimate cognitive load based
// on factors like digit count, memorization likelihood, and computation steps.
// These scores enable adaptive difficulty and performance analysis.
//
// # Error Handling
//
// This package follows two error handling patterns:
//
//   - Generation functions (makeCandidate, makeCandidateRelaxed) return (Candidate, bool)
//     where false indicates an unacceptable candidate that should be skipped during search.
//     This is expected during normal operation.
//
//   - Computation functions (Apply, intPow, factorial) panic on invalid input since they
//     should only receive validated operands from Generate(). Direct callers must validate
//     inputs. This catches bugs early rather than silently producing wrong results.
package operations
