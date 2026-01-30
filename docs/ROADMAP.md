# Roadmap

Development phases for ArithmeGo.

---

## Overview

| Phase | Focus | Status |
|-------|-------|--------|
| 0 | Foundation | Complete |
| 1 | Core Game Loop | Complete |
| 2 | Basic TUI | Complete |
| 3 | Modes System | Complete |
| 4 | Scoring | Complete |
| 5 | Statistics | Complete |
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
- Interface-based Operation system (see `.local/specs/operations-design.md`)
- All 12 operations: basic (4), power (4), advanced (4)
- Intelligent difficulty scoring (see `.local/specs/difficulty-design.md`)
- Question struct with answer validation
- Comprehensive test coverage

### Phase 2: Basic TUI

Playable game in terminal. Menu, game screen, results screen.

Includes:
- Bubble Tea TUI with alt-screen mode
- Menu screen with logo and navigation
- Game screen with timer, question display, numeric input
- Pause screen (hides question, freezes timer)
- Results screen (correct count, accuracy)
- Session management (start, tick, pause/resume, submit, skip)
- Stubs for all future screens

### Phase 3: Modes System

Multiple modes available. Mode selection screen, mode launch with settings.

Includes:
- Mode struct with operations, difficulty, duration
- 7 preset modes: 4 sprints + Mixed Operations + Speed Round + Endurance
- Modes screen with category grouping (Sprint/Challenge)
- Launch screen with difficulty and duration selectors
- Session support for multiple operations (random selection per question)
- Full Menu → Modes → Launch → Game flow

### Phase 4: Scoring

Gamified scoring. Difficulty multipliers, time bonuses, streaks, arcade-style display.

Includes:
- Scoring engine with difficulty multipliers (0.5x–2.0x)
- Time bonus system (1.5x for <2s, linear decay to 1.0x at 10s)
- Streak system with 7 tiers: None, Building, Streak, Max, Blazing, Unstoppable, Legendary
- Streak multiplier (1.0x–2.0x, +0.25 every 5 correct)
- Milestone announcements at streak thresholds (5, 10, 15, 20, 25)
- Animated score display with easing
- Visual streak bar with tier-based styling and shimmer effects
- Delta popup showing points gained/lost
- Best streak tracking per session
- Comprehensive test coverage for scoring calculations

### Phase 5: Statistics

Track and display performance. Per-session and per-question data. Statistics screen.

Includes:
- Storage package for persistent data (`~/.config/arithmego/statistics.json`)
- Per-session tracking (mode, difficulty, duration, score, streak, accuracy)
- Per-question tracking (operation, response time, correct/wrong/skipped)
- Automatic session saving on game completion
- Statistics screen with summary and detailed views
- Summary view: total sessions, accuracy, per-operation accuracy grid
- Detailed view: breakdown by operation, by mode, best streak ever
- Aggregates computed on demand (not stored)
- Comprehensive test coverage for storage operations

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
