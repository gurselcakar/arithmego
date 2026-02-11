# ArithmeGo

Terminal-based arithmetic game for developers. Built with Go, Bubble Tea, and Cobra.

**Status:** v0.1.0

## Docs

- `docs/ARCHITECTURE.md` — Tech stack, project structure

For detailed specs, design docs, and session logs, see `.local/` (gitignored).

## Commands

```bash
make build          # Build for current platform
make build-release  # Build with version info (for releases)
make run            # Build and run
make test           # Run tests
make lint           # Run linter
```

## CLI

```bash
arithmego            # Open main menu (TUI)
arithmego play       # Quick play with last settings
arithmego statistics # View stats directly
arithmego update     # Check for updates
arithmego version    # Show version info
```

## Code Style

- Follow standard Go conventions
- Use `internal/` for all packages (no public API)
- Game logic in `internal/game/` must have no UI or storage imports
- Storage logic in `internal/storage/` must have no UI or game imports

## Question Generation

Questions use an expression tree model (`internal/game/expr/`) with node types: `Num`, `BinOp`, `Paren`, `UnaryPrefix`, `UnarySuffix`, `Pow`.

16 generators in `internal/game/gen/` implement the `Generator` interface (defined in `game/pool.go`). Each generator uses weighted patterns per difficulty level. Generators self-register via `init()` in `gen/registry.go`.

`QuestionPool` handles batch pre-generation (50 at a time) with session-level deduplication via expression keys.

## Mode Categories

Modes are organized into two categories (Sprint / Challenge):

- **Sprint** (single-operation): Addition, Subtraction, Multiplication, Division, Squares, Cubes, Square Roots, Cube Roots, Exponents, Remainders, Percentages, Factorials
- **Challenge** (mixed): Mixed Basics, Mixed Powers, Mixed Advanced, Anything Goes

Each mode maps to a generator via `Mode.GeneratorLabel`. Modes are defined in `internal/modes/presets.go`.

## UI Screens

The TUI uses Bubble Tea with a central App model routing to screen models:

- **Menu** — Main menu with Play, Practice, Statistics, Settings
- **Play Browse** — Mode selection with categories (Basics, Powers, Advanced, Mixed)
- **Play Config** — Game settings (difficulty, duration, input) with live equation preview
- **Game** — Active gameplay with timer, score, streak display
- **Pause** — Paused state with resume/quit options
- **Results** — Post-game summary with score and stats
- **Practice** — Untimed sandbox mode with live settings
- **Statistics** — Performance history and aggregates with session detail view
- **Settings** — User preferences (defaults, auto-update, quit confirmation)
- **Onboarding** — First-run setup flow
- **Feature Tour** — Post-onboarding feature introduction
- **Quit Confirm** — Quit confirmation dialog

Shared helpers are in `internal/ui/screens/helpers.go`. Components (selectors, toggles, hints) are in `internal/ui/components/`.
