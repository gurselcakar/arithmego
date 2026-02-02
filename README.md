<div align="center">

```
 █████╗ ██████╗ ██╗████████╗██╗  ██╗███╗   ███╗███████╗ ██████╗  ██████╗
██╔══██╗██╔══██╗██║╚══██╔══╝██║  ██║████╗ ████║██╔════╝██╔════╝ ██╔═══██╗
███████║██████╔╝██║   ██║   ███████║██╔████╔██║█████╗  ██║  ███╗██║   ██║
██╔══██║██╔══██╗██║   ██║   ██╔══██║██║╚██╔╝██║██╔══╝  ██║   ██║██║   ██║
██║  ██║██║  ██║██║   ██║   ██║  ██║██║ ╚═╝ ██║███████╗╚██████╔╝╚██████╔╝
╚═╝  ╚═╝╚═╝  ╚═╝╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝ ╚═════╝  ╚═════╝
```

**Your AI is thinking. You should too.**

[![Build](https://img.shields.io/github/actions/workflow/status/gurselcakar/arithmego/ci.yml?branch=main)](https://github.com/gurselcakar/arithmego/actions)
[![Release](https://img.shields.io/github/v/release/gurselcakar/arithmego)](https://github.com/gurselcakar/arithmego/releases)
[![License](https://img.shields.io/github/license/gurselcakar/arithmego)](LICENSE)

</div>

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

### macOS / Linux

```bash
curl -fsSL https://arithmego.com/install.sh | bash
```

### Windows

Download the latest `.zip` from [Releases](https://github.com/gurselcakar/arithmego/releases), extract, and add to your PATH.

**Supported platforms:**

| OS | Architecture |
|----|--------------|
| macOS | arm64, amd64 |
| Linux | arm64, amd64 |
| Windows | arm64, amd64 |

---

## Usage

```
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

```
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
