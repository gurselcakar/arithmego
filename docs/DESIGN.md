# Design

> Your AI is thinking. You should too.

---

## Vision

Developers using agentic coding tools experience small gaps while agents work. Most fill this time unproductively—checking phones, scrolling social media, context-switching away.

ArithmeGo is a command-line arithmetic game designed for these moments. Short sessions. Minimal friction. Never leave the terminal.

The goal: restore the mathematical brain, one quick session at a time.

---

## Why CLI

2025 was the year developers shifted to terminal-based workflows. Agentic coding tools moved development into the command line. This is where developers now live.

ArithmeGo meets them there:
- No browser tab to open
- No phone to pick up
- No app to download
- Just launch, play, return to work

---

## Target Audience

**Primary:** Developers using agentic coding tools who want to stay sharp during agent wait times.

**Secondary:** Anyone who lives in the terminal and wants a minimal way to exercise their brain.

---

## Design Principles

1. **Minimal friction** — Start playing in seconds
2. **Clean and elegant** — No visual noise
3. **Respect user time** — Short sessions by design
4. **Stay in flow** — Never pull the user away from their environment
5. **Depth for those who want it** — Simple by default, configurable for power users

---

## Game Mechanics

### Operations

12 operations across three categories:

| Category | Operations |
|----------|------------|
| **Basic** | Addition, Subtraction, Multiplication, Division |
| **Power** | Square, Cube, Square Root, Cube Root |
| **Advanced** | Modulo, Power, Percentage, Factorial |

### Difficulty Tiers

Difficulty is based on cognitive complexity, not just number size:

| Tier | Description |
|------|-------------|
| Beginner | Trivial, instant recall |
| Easy | Simple computation |
| Medium | Requires focus |
| Hard | Challenging |
| Expert | Demanding |

### Session Lengths

Short by design:
- 30 seconds
- 1 minute
- 90 seconds
- 2 minutes

### Input Methods

Two options:
- **Typing** — User types the numerical answer
- **Multiple choice** — Four options, select one

### Scoring

Points earned per correct answer, scaled by:
- Difficulty multiplier
- Response time bonus
- Streak multiplier

Wrong answers incur a small penalty. Skips reset streak but no penalty.

---

## Game Modes

### Core Modes

- **Addition Sprint** — Addition only
- **Subtraction Sprint** — Subtraction only
- **Multiplication Sprint** — Multiplication only
- **Division Sprint** — Division only
- **Mixed Operations** — All four operators
- **Speed Round** — 30 seconds, all operators
- **Endurance** — 2 minutes, all operators

### Practice Mode

A sandbox with no timer and no score. Change difficulty and operation on-the-fly.

---

## Menu Structure

```
Quick Play · [Mode]    → Instant launch (returning users)
Modes                  → Browse and select modes
Practice               → Sandbox mode

Statistics             → Performance data
Settings               → Preferences
```

Quick Play only appears after the first game.

---

## Statistics

Local storage of performance data:
- Sessions played
- Accuracy overall and by operation
- Best streaks
- Response time trends

All data stays local. User owns their data.

---

## Future Possibilities

Documented for later, not current scope:
- Adaptive difficulty (adjusts based on performance)
- Leaderboards
- PvP mode
- Custom level builder
- Compound operations (PEMDAS, fractions)
