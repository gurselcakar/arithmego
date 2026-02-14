# ArithmeGo

Terminal-based arithmetic game built with Go, Bubble Tea, and Cobra.

## Commands

```bash
make build          # Build for current platform
make build-release  # Build with version info (for releases)
make run            # Build and run
make test           # Run tests
make lint           # Run linter
```

## Code Style

- Use `internal/` for all packages (no public API)
- IMPORTANT: Game logic in `internal/game/` must have no UI or storage imports
- IMPORTANT: Storage logic in `internal/storage/` must have no UI or game imports

## Architecture

- See `docs/ARCHITECTURE.md` for tech stack and project structure
- 16 modes defined in `internal/modes/presets.go`, organized as Sprint (single-op) and Challenge (mixed)
- Each mode maps to a generator via `Mode.GeneratorLabel`
- Generators in `internal/game/gen/` self-register via `init()` — new generators must follow this pattern
- Expression tree model in `internal/game/expr/` — all question display/eval flows through this
- TUI uses Bubble Tea with a central App model (`internal/ui/app.go`) routing to screen models in `internal/ui/screens/`
