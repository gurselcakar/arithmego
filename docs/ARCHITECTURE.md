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
│   │   ├── operation.go            # Operation interface
│   │   ├── question.go             # Question struct
│   │   ├── difficulty.go           # Difficulty tiers
│   │   ├── engine.go               # Generation helpers
│   │   ├── scoring.go
│   │   ├── session.go
│   │   └── operations/             # Operation implementations
│   │       ├── registry.go         # Registration system
│   │       ├── helpers.go          # Shared utilities
│   │       ├── addition.go         # Basic binary
│   │       ├── subtraction.go
│   │       ├── multiplication.go
│   │       ├── division.go
│   │       ├── modulo.go           # Additional binary
│   │       ├── power.go
│   │       ├── percentage.go
│   │       ├── square.go           # Unary
│   │       ├── cube.go
│   │       ├── square_root.go
│   │       ├── cube_root.go
│   │       └── factorial.go
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
| **Extensibility** | Operation interface allows easy addition of new operations. |
| **Testability** | Core logic is pure Go, unit testable without TUI. |
| **Open/Closed** | Add operations without modifying existing code. |

---

## Operations Architecture

Operations use an interface-based design for maximum extensibility. Each operation is self-contained and encapsulates its own logic.

```go
type Operation interface {
    Name() string           // "Square Root"
    Symbol() string         // "√"
    Arity() Arity           // Unary (1) or Binary (2)
    Category() Category     // basic, power, advanced
    Apply(operands []int) int
    ScoreDifficulty(operands []int, answer int) float64
    Generate(diff Difficulty) Question
    Format(operands []int) string
}
```

### Supported Operations

| Category | Operations |
|----------|------------|
| **Basic** | Addition, Subtraction, Multiplication, Division |
| **Power** | Square, Cube, Square Root, Cube Root |
| **Advanced** | Modulo, Power, Percentage, Factorial |

### Adding New Operations

1. Create a new file in `internal/game/operations/`
2. Implement the `Operation` interface
3. Register via `init()` function
4. Add tests

No existing code needs modification. See `.local/specs/operations-design.md` for full specification.

---

## Difficulty Scoring

Difficulty uses intelligent scoring based on cognitive factors, not just number size. Each operation implements `ScoreDifficulty()` to compute problem difficulty.

| Tier | Score Range | Description |
|------|-------------|-------------|
| **Beginner** | 1.0 - 2.0 | Trivial, instant recall |
| **Easy** | 2.0 - 4.0 | Simple computation |
| **Medium** | 4.0 - 6.0 | Requires focus |
| **Hard** | 6.0 - 8.0 | Challenging |
| **Expert** | 8.0 - 10.0 | Demanding |

### Scoring Factors

Each operation considers operation-specific factors:
- **Addition**: Carries, digit count, nice numbers
- **Subtraction**: Borrows, digit count, zeros
- **Multiplication**: Digit combinations, easy multipliers
- **Division**: Times table inverse, quotient size

See `.local/specs/difficulty-design.md` for full specification including scoring weights.

---

## UI Architecture

The TUI uses Bubble Tea's Model-View-Update pattern with a central App orchestrator.

### Screen Flow

```
Menu
 ├─ Quick Play* ───────→ Game ←→ Pause
 │                         ↓
 ├─ Modes → Launch ──────→ Game ←→ Pause
 │                           ↓
 │                        Results → Menu / Play Again
 ├─ Practice
 ├─ Statistics
 └─ Settings

* Quick Play appears only for returning users (Phase 6)
  Uses last played mode, bypasses Modes/Launch screens
```

### Message Types

Screen transitions are handled via typed messages:

| Message | Trigger |
|---------|---------|
| `MenuSelectMsg` | Menu item selected |
| `GameOverMsg` | Timer expired |
| `PauseMsg` | User pressed P/Esc |
| `ResumeMsg` | User pressed Enter on pause |
| `QuitToMenuMsg` | User quit from pause |
| `PlayAgainMsg` | User pressed Enter on results |
| `ReturnToMenuMsg` | User pressed M/Esc on results |

### Components

Reusable UI elements in `internal/ui/components/`:

| Component | Purpose |
|-----------|---------|
| `logo.go` | ASCII art logo and tagline |
| `timer.go` | MM:SS countdown format |
| `question.go` | Question display |
| `input.go` | Numeric text input |
| `keyhints.go` | Navigation hints |
| `scoreboard.go` | Points/streak display (Phase 4) |
| `choices.go` | Multiple choice input (Phase 10) |

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
