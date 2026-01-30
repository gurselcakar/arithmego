# ArithmeGo

Command-line arithmetic game for developers. Built with Go, Bubble Tea, and Cobra.

**Status:** Phase 5 complete, ready to start Phase 6

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
- Game logic in `internal/game/` must have no UI or storage imports
- Storage logic in `internal/storage/` must have no UI or game imports
