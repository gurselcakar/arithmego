// Package cli implements the command-line interface for ArithmeGo using Cobra.
//
// The CLI provides the following commands:
//
//   - arithmego: Opens the main menu (default behavior)
//   - arithmego play [mode]: Opens play browse, or config for a specific mode
//   - arithmego statistics: Opens the statistics screen directly
//   - arithmego update: Checks for available updates
//   - arithmego version: Displays version and build information
//
// The play command accepts an optional mode argument to jump directly to
// the configuration screen for that mode. Valid modes include: addition,
// subtraction, multiplication, division, squares, cubes, square-roots,
// cube-roots, exponents, remainders, percentages, factorials, mixed-basics,
// mixed-powers, mixed-advanced, and anything-goes.
//
// Version information (Version, CommitSHA, BuildDate) is injected via
// ldflags during the build process. See the Makefile for details.
package cli
