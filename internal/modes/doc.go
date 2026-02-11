// Package modes defines game mode configurations and their registry.
//
// A Mode represents a playable game configuration that maps to a question
// generator via GeneratorLabel. Each mode specifies default difficulty and
// duration settings. Modes are organized into categories (Sprint for
// single-operation modes, Challenge for mixed modes) for UI grouping.
//
// The package provides 16 built-in modes:
//   - 4 Basic: Addition, Subtraction, Multiplication, Division
//   - 4 Power: Squares, Cubes, Square Roots, Cube Roots
//   - 4 Advanced: Exponents, Remainders, Percentages, Factorials
//   - 4 Mixed: Mixed Basics, Mixed Powers, Mixed Advanced, Anything Goes
//
// Use [Get] to retrieve a mode by ID, [All] to list all registered modes,
// and [Register] to add custom modes. [RegisterPresets] registers all built-in
// modes and must be called after generators are registered in game/gen.
package modes
