# ArithmeGo

> Your AI is thinking. You should too.

A command-line arithmetic game designed for developers using agentic coding tools.

---

## About

ArithmeGo fills the small gaps while your AI agent works. Short sessions. Minimal friction. Never leave the terminal.

**Status:** In development (Phase 11 — Polish)

---

## Features

- 12 operations across three categories:
  - **Basic:** Addition, Subtraction, Multiplication, Division
  - **Power:** Square, Cube, Square Root, Cube Root
  - **Advanced:** Modulo, Power, Percentage, Factorial
- 5 difficulty tiers (Beginner to Expert) with intelligent difficulty scoring
- Timed sessions (30s, 60s, 90s, 2min)
- 16 game modes (12 single-operation + 4 mixed modes)
- Arcade-style scoring with streaks and multipliers
- Quick Play (instant launch with last settings)
- Practice mode (untimed sandbox)
- Multiple choice and typing input modes
- First-run onboarding and feature tour
- Local statistics with history and insights
- CLI subcommands (`arithmego play`, `arithmego statistics`, `arithmego update`)
- Auto-update checking
- Terminal-native design (inherits your theme)

---

## Installation

```bash
curl -fsSL https://arithmego.com/install.sh | bash
```

Supports macOS and Linux (arm64, amd64). See [arithmego.com](https://arithmego.com) for details.

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

- [arithmego.com](https://arithmego.com) — Website and user docs
- [Design](docs/DESIGN.md) — Vision and game mechanics
- [Architecture](docs/ARCHITECTURE.md) — Tech stack and project structure
- [Roadmap](docs/ROADMAP.md) — Development phases

---

## License

MIT
