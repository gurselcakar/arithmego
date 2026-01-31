# ArithmeGo

Command-line arithmetic game for developers. Built with Go, Bubble Tea, and Cobra.

**Status:** Phase 10 complete (CLI Commands)

## Docs

- `docs/DESIGN.md` — Vision, mechanics, game modes
- `docs/ARCHITECTURE.md` — Tech stack, project structure
- `docs/ROADMAP.md` — Development phases

For detailed specs and session logs, see `.local/` (gitignored).

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
