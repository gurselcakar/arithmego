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
    operations/           Operation implementations
  modes/                  Game mode definitions
  ui/                     Bubble Tea UI layer
    screens/              Screen models
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

### Interface-Based Operations

Operations (addition, subtraction, multiplication, etc.) implement a common interface:

```go
type Operation interface {
    Generate(difficulty int) Equation
    Apply(a, b int) int
    ScoreDifficulty(equation Equation) int
}
```

This enables:
- Easy addition of new operation types
- Consistent difficulty scaling
- Composable mixed-mode games

### Cognitive Difficulty Scoring

Difficulty is based on cognitive complexity rather than just number magnitude. Factors include:
- Number of digits
- Carrying/borrowing requirements
- Mental calculation steps

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
 ├── Play Browse → Play Config → Game → Results
 ├── Practice
 ├── Statistics
 └── Settings
```

Each screen is a self-contained Bubble Tea model.

## CLI Commands

| Command | Description |
|---------|-------------|
| `arithmego` | Opens the TUI main menu |
| `arithmego play` | Quick play with last used settings |
| `arithmego statistics` | View performance statistics |
| `arithmego update` | Check for updates |
| `arithmego version` | Show version information |

## Data Storage

All user data is stored locally:

```
~/.config/arithmego/
  config.json       User preferences
  statistics.json   Game history and stats
```

No data is sent externally. The update checker only fetches release metadata from GitHub.
