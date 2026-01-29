# ArithmeGo

Command-line arithmetic game for developers. Built with Go, Bubble Tea, and Cobra.

**Status:** Pre-development (planning complete, ready to start Phase 0)

## Project Context

- Read @docs/DESIGN.md for vision, mechanics, and game modes
- Read @docs/ARCHITECTURE.md for tech stack and project structure
- Read @docs/ROADMAP.md for development phases

## Detailed Planning (Private)

For deeper context, check `.local/`:

- `.local/mechanics.md` — Detailed game rules, scoring formulas, statistics schema
- `.local/ui-design.md` — All screen mockups and UI specifications
- `.local/conversation.md` — First planning session decisions
- `.local/conversation-N.md` — Subsequent session logs

## Commands

No code yet. When development begins:

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

## Conventions

- Document each coding session in `.local/conversation-N.md`
- Follow phases in @docs/ROADMAP.md sequentially
- Write tests alongside implementation
- Commit frequently with clear messages
