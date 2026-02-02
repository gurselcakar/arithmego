# ArithmeGo

**Your AI is thinking. You should too.**

A command-line arithmetic game for developers. Built with Go, Bubble Tea, and Cobra.

<!-- Badges -->
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.25+-blue)
![License](https://img.shields.io/badge/license-MIT-green)

---

## Features

| Category | Details |
|----------|---------|
| Operations | 12 operations across 3 categories (Basic, Power, Advanced) |
| Difficulty | 5 tiers from Beginner to Expert |
| Game Modes | 16 modes (12 single-operation + 4 mixed) |
| Sessions | Timed gameplay from 30 seconds to 2 minutes |
| Scoring | Arcade-style with streak bonuses |
| Input | Multiple choice or typing |
| Statistics | Local performance tracking |

**Additional highlights:**

- **Quick Play** — Jump in with your last settings
- **Practice Mode** — Untimed sandbox for learning
- **Terminal-native** — Inherits your terminal theme

---

## Installation

```bash
curl -fsSL https://arithmego.com/install.sh | bash
```

**Supported platforms:**

| OS | Architecture |
|----|--------------|
| macOS | arm64, amd64 |
| Linux | arm64, amd64 |

---

## Usage

```bash
arithmego            # Open main menu (TUI)
arithmego play       # Quick play with last settings
arithmego statistics # View stats directly
arithmego update     # Check for updates
arithmego version    # Show version info
```

---

## Development

### Prerequisites

- Go 1.25+

### Commands

```bash
make build    # Build for current platform
make run      # Build and run
make test     # Run tests
make lint     # Run linter
```

---

## Documentation

- [Architecture](docs/ARCHITECTURE.md) — Tech stack and project structure
- [Contributing](docs/CONTRIBUTING.md) — Guidelines for contributors

---

## License

MIT License. See [LICENSE](LICENSE) for details.
