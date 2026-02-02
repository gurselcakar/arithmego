# ArithmeGo

Terminal-based arithmetic game for developers. Built with Go, Bubble Tea, and Cobra.

**Status:** v0.1.0

## Docs

- `docs/ARCHITECTURE.md` — Tech stack, project structure
- `docs/CONTRIBUTING.md` — Guidelines for contributors

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

## UI Screens

The TUI uses Bubble Tea with a central App model routing to screen models:

- **Menu** — Main menu with Play, Practice, Statistics, Settings
- **Play Browse** — Mode selection with categories (Basics, Powers, Advanced, Mixed) and search
- **Play Config** — Game settings (difficulty, duration, input) with live equation preview
- **Game** — Active gameplay with timer, score, streak display
- **Pause** — Paused state with resume/quit options
- **Results** — Post-game summary with score and stats
- **Practice** — Untimed sandbox mode with live settings
- **Statistics** — Performance history and aggregates
- **Settings** — User preferences (defaults, auto-update, quit confirmation)
- **Onboarding** — First-run setup flow
- **Feature Tour** — Post-onboarding feature introduction

Shared helpers are in `internal/ui/screens/helpers.go`. Components (selectors, toggles, hints) are in `internal/ui/components/`.
