# Architecture

This document describes the high-level architecture of ArithmeGo.

## Tech Stack

- **Go** - Primary language
- **Bubble Tea** - Terminal UI framework (Elm-inspired Model-View-Update)
- **Lip Gloss** - Styling and layout
- **Cobra** - CLI command framework

## Project Structure

```
cmd/arithmego/main.go     Entry point

internal/
  cli/                    Cobra commands
  game/                   Core game logic
    expr/                 Expression tree (nodes, eval, format, key)
    gen/                  16 question generators + framework
  modes/                  Game mode definitions (Sprint / Challenge)
  ui/                     Bubble Tea UI layer
    screens/              Screen models
      statistics/         Statistics sub-screens
    components/           Reusable UI components
    styles/               Styling constants
  storage/                Local persistence
  analytics/              Statistics computation
  update/                 Update checking
```

## Design Principles

### Separation of Concerns

The codebase enforces strict boundaries between layers:

- **game/** contains pure game logic with no UI or storage imports
- **storage/** handles persistence with no game logic imports
- **ui/** depends on game/ and storage/ but they don't depend on it

### Expression Tree Model

Questions are built using an expression tree (`game/expr/`):

- **Node types**: `Num`, `BinOp`, `Paren`, `UnaryPrefix`, `UnarySuffix`, `Pow`
- **Evaluation**: `Eval()` computes the integer result
- **Formatting**: `Format()` renders with Unicode math symbols (×, ÷, √, ², ³)
- **Deduplication**: `Key()` produces a canonical string for duplicate detection

### Generator-Based Question System

Generators implement a common interface (defined in `game/pool.go`):

```go
type Generator interface {
    Generate(diff Difficulty) *Question
    Label() string
}
```

Each of the 16 generators uses weighted patterns per difficulty level. Generators self-register via `init()` in `gen/registry.go`. This enables:
- Easy addition of new question types
- Consistent difficulty scaling via pattern weights
- Composable mixed-mode generators that delegate to single-operation generators
- Multi-operand expressions (e.g., 3+4+5) and PEMDAS-aware expressions (e.g., 5+3×2)

### Question Pool

`QuestionPool` handles batch pre-generation of 50 questions at a time with session-level deduplication via expression keys. The pool auto-refills when exhausted.

### Difficulty Levels

Five difficulty levels (Beginner, Easy, Medium, Hard, Expert) affect number ranges and expression complexity via weighted pattern selection. Each generator defines its own pattern set.

## UI Architecture

### Model-View-Update Pattern

The UI follows Bubble Tea's MVU architecture:

1. **Model** - Application state
2. **Update** - State transitions based on messages
3. **View** - Render state to terminal output

### Screen Routing

A central `App` model manages screen transitions:

```
App
 ├── Menu
 ├── Play Browse → Play Config → Game → Pause / Results
 ├── Practice
 ├── Statistics → Session Detail
 ├── Settings
 ├── Onboarding → Feature Tour
 └── Quit Confirm
```

Each screen is a self-contained Bubble Tea model.

## CLI Commands

| Command | Description |
|---------|-------------|
| `arithmego` | Opens the TUI main menu |
| `arithmego play` | Quick play with last used settings |
| `arithmego play [mode]` | Jump to config for a specific mode |
| `arithmego statistics` | View performance statistics |
| `arithmego update` | Check for updates |
| `arithmego version` | Show version information |

## Data Storage

All user data is stored locally in the system config directory:

| OS | Location |
|----|----------|
| macOS | `~/Library/Application Support/arithmego/` |
| Linux | `~/.config/arithmego/` |
| Windows | `%AppData%\arithmego\` |

**Files:**
- `config.json` — User preferences
- `statistics.json` — Game history and stats

No data is sent externally. The update checker only fetches release metadata from GitHub.
