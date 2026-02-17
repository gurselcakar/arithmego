<div align="center">

<img src="assets/logo-tagline.png" alt="ArithmeGo" width="600">

[![Build](https://img.shields.io/github/actions/workflow/status/gurselcakar/arithmego/ci.yml?branch=main)](https://github.com/gurselcakar/arithmego/actions)
[![Release](https://img.shields.io/github/v/release/gurselcakar/arithmego)](https://github.com/gurselcakar/arithmego/releases)
[![Downloads](https://img.shields.io/github/downloads/gurselcakar/arithmego/total)](https://github.com/gurselcakar/arithmego/releases)
[![Views](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgurselcakar%2Farithmego&count_bg=%2379C83D&title_bg=%23555555&title=views&edge_flat=false)](https://hits.seeyoufarm.com)
[![License](https://img.shields.io/github/license/gurselcakar/arithmego)](LICENSE)

</div>

## Install

**macOS / Linux**
```bash
curl -fsSL https://arithmego.com/install.sh | bash
```

**Windows (PowerShell)**
```powershell
irm https://arithmego.com/install.ps1 | iex
```

## About

ArithmeGo is an arithmetic game that runs in your terminal.

AI agents are getting better at handling longer tasks. While they work,
you wait. ArithmeGo was built to fill that gap with something useful:
mental math practice, right where you already are.

It covers basic arithmetic, powers and roots, and advanced operations
like modulo and factorials. Five difficulty levels from beginner to
expert. Timed sprints with scoring and streaks, or untimed practice
at your own pace. All progress is tracked locally.

## Usage

```
arithmego            # Open main menu
arithmego play       # Browse and pick a game mode
arithmego practice   # Start practice mode
arithmego statistics # View your stats
arithmego settings   # Adjust your preferences
```

## Development

Requires Go 1.25+.

```
make build    # Build for current platform
make run      # Build and run
make test     # Run tests
make lint     # Run linter
```

## Docs

- [Architecture](docs/ARCHITECTURE.md) — Tech stack and project structure
- [Website](docs/WEBSITE.md) — Hugo site structure and content editing

## License

MIT. See [LICENSE](LICENSE).
