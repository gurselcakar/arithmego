# Architecture

Technical decisions and project structure for ArithmeGo.

---

## Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go** | Language |
| **Bubble Tea** | TUI framework |
| **Lip Gloss** | Styling |
| **Cobra** | CLI framework |

---

## CLI Commands

```
arithmego              → Opens main menu (TUI)
arithmego play         → Instant quick play
arithmego statistics   → Opens statistics screen
arithmego update       → Check for updates
arithmego version      → Show version
```

---

## Data Storage

Location: `~/.config/arithmego/`

```
~/.config/arithmego/
├── config.json        # User preferences
├── statistics.json    # Game history
├── modes.json         # Custom modes
└── state.json         # Last played settings
```

All files are JSON. Human-readable, easy to debug.

---

## Project Structure

```
arithmego/
├── cmd/
│   └── arithmego/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── cli/                        # Cobra commands
│   │   ├── root.go
│   │   ├── play.go
│   │   ├── statistics.go
│   │   ├── update.go
│   │   └── version.go
│   │
│   ├── game/                       # Core game logic
│   │   ├── engine.go
│   │   ├── question.go
│   │   ├── operations.go
│   │   ├── difficulty.go
│   │   ├── scoring.go
│   │   └── session.go
│   │
│   ├── modes/                      # Game modes
│   │   ├── mode.go
│   │   ├── registry.go
│   │   ├── presets.go
│   │   └── custom.go
│   │
│   ├── ui/                         # Bubble Tea UI
│   │   ├── app.go
│   │   ├── router.go
│   │   ├── screens/
│   │   │   ├── menu.go
│   │   │   ├── modes.go
│   │   │   ├── launch.go
│   │   │   ├── game.go
│   │   │   ├── pause.go
│   │   │   ├── results.go
│   │   │   ├── practice.go
│   │   │   ├── statistics.go
│   │   │   ├── settings.go
│   │   │   └── onboarding.go
│   │   ├── components/
│   │   │   ├── logo.go
│   │   │   ├── timer.go
│   │   │   ├── question.go
│   │   │   ├── input.go
│   │   │   ├── choices.go
│   │   │   ├── scoreboard.go
│   │   │   ├── keyhints.go
│   │   │   └── confirm.go
│   │   └── styles/
│   │       └── styles.go
│   │
│   ├── storage/                    # Persistence
│   │   ├── storage.go
│   │   ├── paths.go
│   │   ├── config.go
│   │   ├── state.go
│   │   ├── statistics.go
│   │   └── modes.go
│   │
│   └── analytics/                  # Statistics
│       ├── tracker.go
│       ├── aggregator.go
│       └── insights.go
│
├── scripts/
│   └── install.sh
│
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
│
├── docs/
├── .goreleaser.yaml
├── Makefile
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

---

## Design Principles

| Principle | Implementation |
|-----------|----------------|
| **Separation of concerns** | Game logic has no UI imports. UI has no storage imports. |
| **Extensibility** | Mode interface allows easy addition of new modes. |
| **Testability** | Core logic is pure Go, unit testable without TUI. |

---

## Distribution

Primary method: curl install script

```bash
curl -fsSL https://arithmego.com/install.sh | bash
```

### Supported Platforms

| OS | Architecture |
|----|--------------|
| macOS | arm64, amd64 |
| Linux | arm64, amd64 |

Windows not supported initially.

---

## Build & Release

- **Makefile** — Local development (`build`, `run`, `test`, `lint`)
- **GoReleaser** — Cross-compilation for releases
- **GitHub Actions** — CI on push, release on tag
