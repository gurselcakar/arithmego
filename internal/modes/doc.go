// Package modes defines game mode configurations and their registry.
//
// A Mode represents a playable game configuration with a set of operations,
// default difficulty, and duration. Modes are organized into categories
// (Sprint for single-operation, Challenge for mixed) and registered at startup.
//
// Use [Get] to retrieve a mode by ID or [All] to list all modes.
// The [RegisterPresets] function registers the built-in modes and must be
// called after generators are registered.
package modes
