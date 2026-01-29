# ArithmeGo

Command-line arithmetic game for developers. Built with Go, Bubble Tea, and Cobra.

**Status:** Phase 1 complete, ready to start Phase 2

## Docs

- `docs/DESIGN.md` — Vision, mechanics, game modes
- `docs/ARCHITECTURE.md` — Tech stack, project structure
- `docs/ROADMAP.md` — Development phases

For detailed specs and session logs, see `.local/` (gitignored).

## Commands

```bash
make build    # Build for current platform
make run      # Build and run
make test     # Run tests
make lint     # Run linter
```

## Code Style

- Follow standard Go conventions
- Use `internal/` for all packages (no public API)
- Game logic in `internal/game/` must have no UI imports
- UI code in `internal/ui/` must have no storage imports
