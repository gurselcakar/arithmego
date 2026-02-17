# Architecture

This document describes the high-level architecture of ArithmeGo.

## Tech Stack

- **Go** - Primary language
- **Bubble Tea** - Terminal UI framework (Elm-inspired Model-View-Update)
- **Bubbles** - Reusable Bubble Tea components (viewport, text input)
- **Lip Gloss** - Styling and layout
- **Cobra** - CLI command framework

## Project Structure

```
cmd/arithmego/main.go     Entry point

internal/
  cli/                    Cobra commands (root, play, practice, statistics, settings, update, version)
  game/                   Core game logic
    expr/                 Expression tree (nodes, eval, format, key)
    gen/                  16 question generators + framework
  modes/                  Game mode definitions (Sprint / Challenge)
  ui/                     Bubble Tea UI layer
    screens/              Screen models
      statistics/         Statistics sub-screens (dashboard, operations, history, trends, charts)
    components/           Reusable UI components (timer, input, choices, scoreboard, keyhints, etc.)
    styles/               Styling constants
  storage/                Local persistence (config, statistics, paths)
  analytics/              Statistics computation (aggregates, filters, trends)
  update/                 Update checking and auto-update

website/                  Hugo static site (arithmego.com)
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
- **Binary ops**: `+`, `−`, `×`, `÷`, `mod`, `% of`
- **Unary ops**: `√`, `∛` (prefix); `²`, `³`, `!` (suffix)
- **Evaluation**: `Eval()` computes the integer result
- **Formatting**: `Format()` renders with Unicode math symbols, auto-parenthesizes based on PEMDAS precedence
- **Deduplication**: `Key()` produces a canonical prefix-notation string for duplicate detection

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

`QuestionPool` handles batch pre-generation of 50 questions at a time with session-level deduplication via expression keys. The pool auto-refills when exhausted, and clears the seen set if truly exhausted to avoid deadlock.

### Scoring System

Points are calculated as: `base × difficulty × time bonus × streak bonus`.

- **Difficulty multiplier**: 0.5× (Beginner) to 2.0× (Expert)
- **Time bonus**: 1.5× for instant answers (< 2s), linear decay to 1.0× at 10s
- **Streak bonus**: +0.25× every 5 correct answers, capped at 2.0× at streak 20
- **Streak tiers**: Building → Streak → Max → Blazing → Unstoppable → Legendary

### Multiple Choice

The `game/choices.go` module generates distractor answers for multiple choice input mode. Distractors are produced using offset-based algorithms to create plausible wrong answers.

### Difficulty Levels

Five difficulty levels (Beginner, Easy, Medium, Hard, Expert) affect number ranges and expression complexity via weighted pattern selection. Each generator defines its own pattern set.

## UI Architecture

### Model-View-Update Pattern

The UI follows Bubble Tea's MVU architecture:

1. **Model** - Application state
2. **Update** - State transitions based on messages
3. **View** - Render state to terminal output

### Screen Routing

A central `App` model (`ui/app.go`) manages screen transitions:

```
App
 ├── Menu
 ├── Play Browse → Play Config → Game → Pause / Results
 ├── Practice
 ├── Statistics (Dashboard → Operations → Operation Detail → Operation Review)
 │              (Dashboard → History → Session Detail → Session Full Log)
 │              (Dashboard → Trends)
 ├── Settings
 ├── Onboarding → Game → Feature Tour
 └── Quit Confirm
```

Each screen is a self-contained Bubble Tea model. The statistics screen uses a sub-model architecture with multiple views sharing a single model.

## CLI Commands

| Command | Description |
|---------|-------------|
| `arithmego` | Opens the TUI main menu |
| `arithmego play` | Browse all game modes |
| `arithmego play [mode]` | Jump to config for a specific mode |
| `arithmego practice` | Start practice mode |
| `arithmego statistics` | View performance statistics |
| `arithmego settings` | Open settings |
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
- `config.json` — User preferences (last played mode, input method, onboarding state, etc.)
- `statistics.json` — Game session history and per-question records

No data is sent externally. The update module fetches release metadata from GitHub and can auto-download binary updates.

## Website

A Hugo static site in `website/` serves as the project homepage. It includes video showcases of gameplay, a download page with install instructions, and is deployed to `arithmego.com`.
