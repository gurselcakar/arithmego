# ArithmeGo

> Your AI is thinking. You should too.

A command-line arithmetic game designed for developers using agentic coding tools.

---

## About

ArithmeGo fills the small gaps while your AI agent works. Short sessions. Minimal friction. Never leave the terminal.

**Status:** In development (Phase 5 complete)

---

## Features

### Implemented

- 12 operations across three categories:
  - **Basic:** Addition, Subtraction, Multiplication, Division
  - **Power:** Square, Cube, Square Root, Cube Root
  - **Advanced:** Modulo, Power, Percentage, Factorial
- 5 difficulty tiers (Beginner to Expert) with intelligent difficulty scoring
- Timed sessions (30s, 60s, 90s, 2min)
- 16 game modes (12 single-operation + 4 mixed modes)
- Arcade-style scoring with streaks and multipliers
- Local statistics tracking with per-session and per-question data
- Terminal-native design (inherits your theme)

### Coming Soon

- Quick Play (instant launch with last mode)
- Practice mode (no timer, no pressure)
- Settings and onboarding
- CLI subcommands (`arithmego play`, `arithmego stats`)

---

## Installation

Coming soon.

```bash
curl -fsSL https://arithmego.com/install.sh | bash
```

---

## Development

```bash
make build    # Build for current platform
make run      # Build and run
make test     # Run tests
make lint     # Run linter
```

---

## Documentation

- [Design](docs/DESIGN.md) — Vision and game mechanics
- [Architecture](docs/ARCHITECTURE.md) — Tech stack and project structure
- [Roadmap](docs/ROADMAP.md) — Development phases

---

## License

MIT
