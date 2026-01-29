# Roadmap

Development phases for ArithmeGo.

---

## Overview

| Phase | Focus | Status |
|-------|-------|--------|
| 0 | Foundation | Complete |
| 1 | Core Game Loop | Planned |
| 2 | Basic TUI | Planned |
| 3 | Modes System | Planned |
| 4 | Scoring | Planned |
| 5 | Statistics | Planned |
| 6 | Quick Play | Planned |
| 7 | Practice Mode | Planned |
| 8 | Settings & Onboarding | Planned |
| 9 | CLI Commands | Planned |
| 10 | Polish | Planned |
| 11 | Distribution | Planned |

Each phase builds on the previous. Each delivers something testable.

---

## Phase Descriptions

### Phase 0: Foundation

Project setup. Go module, directory structure, dependencies, Makefile.

### Phase 1: Core Game Loop

Question generation and answer validation. No UI—pure logic with unit tests.

Includes:
- Interface-based Operation system (see `.local/operations-design.md`)
- All 12 operations: basic (4), power (4), advanced (4)
- Intelligent difficulty scoring (see `.local/difficulty-design.md`)
- Question struct with answer validation
- Comprehensive test coverage

### Phase 2: Basic TUI

Playable game in terminal. Menu, game screen, results screen.

### Phase 3: Modes System

Multiple modes available. Mode selection screen, mode launch with settings.

### Phase 4: Scoring

Gamified scoring. Difficulty multipliers, time bonuses, streaks, arcade-style display.

### Phase 5: Statistics

Track and display performance. Per-session and per-question data. Statistics screen.

### Phase 6: Quick Play

Remember last played mode. Quick Play option on main menu for returning users.

### Phase 7: Practice Mode

Sandbox for drilling. No timer, no score. Live controls to change difficulty and operation.

### Phase 8: Settings & Onboarding

First-time user experience. Preferences screen. Guided setup flow.

### Phase 9: CLI Commands

Direct access via subcommands: `arithmego play`, `arithmego statistics`, etc.

### Phase 10: Polish

Edge cases, pause/quit flow, multiple choice input, error handling, terminal resize.

### Phase 11: Distribution

Install script, GoReleaser config, GitHub Actions, README, LICENSE.

---

## Post-MVP

Ideas documented for future consideration:

- Adaptive difficulty (adjusts based on user performance)
- Adaptive placement test
- Custom user-created modes
- Post-game analysis
- Trend graphs
- Leaderboards (requires backend)
- PvP mode (requires backend)
- Compound operations (PEMDAS, fractions)
- Data export

---

## Development Approach

- Each phase is a working state
- Write tests alongside implementation
- Commit frequently
- Tag releases at milestones
