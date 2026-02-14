# Game Logic & Question Generation

This document explains how ArithmeGo generates questions, manages difficulty, scores players, and runs game sessions. It's written for contributors who want to understand the current system and find areas to improve.

**Key files:**

| Area | Path |
|------|------|
| Expression tree | `internal/game/expr/` |
| Generators | `internal/game/gen/` |
| Question pool | `internal/game/pool.go` |
| Session & scoring | `internal/game/session.go`, `scoring.go` |
| Difficulty | `internal/game/difficulty.go` |
| Modes | `internal/modes/presets.go` |

## Table of Contents

- [Expression Tree](#expression-tree)
- [Question Generation](#question-generation)
  - [Generator Framework](#generator-framework)
  - [Pattern System](#pattern-system)
  - [Operand Ranges](#operand-ranges)
- [The 16 Generators](#the-16-generators)
  - [Sprint Modes (Single-Operation)](#sprint-modes-single-operation)
  - [Challenge Modes (Mixed)](#challenge-modes-mixed)
- [Difficulty System](#difficulty-system)
- [Question Pool & Deduplication](#question-pool--deduplication)
- [Session Lifecycle](#session-lifecycle)
- [Scoring System](#scoring-system)
- [Multiple Choice](#multiple-choice)
- [Areas for Improvement](#areas-for-improvement)

---

## Expression Tree

All questions are represented as expression trees built from these node types:

| Node | Description | Example Display |
|------|-------------|-----------------|
| `Num` | Integer literal | `5` |
| `BinOp` | Binary operator (+, -, x, /, mod, % of) | `5 + 3` |
| `Paren` | Display-only parentheses | `(5 + 3)` |
| `UnaryPrefix` | Square root, cube root | `√49`, `∛27` |
| `UnarySuffix` | Square, cube, factorial | `7²`, `3³`, `5!` |
| `Pow` | Exponentiation | `2⁴` |

Every node implements three methods:

- **`Eval()`** — Computes the integer answer
- **`Format()`** — Renders with Unicode math symbols (×, ÷, √, superscript digits)
- **`Key()`** — Produces a canonical prefix-notation string for deduplication (e.g., `(+ 5 (* 3 2))`)

The formatting layer automatically handles PEMDAS parenthesization — a `BinOp` wraps its children in parentheses when their precedence is lower than the parent's.

---

## Question Generation

### Generator Framework

Each of the 16 game modes has a dedicated generator that implements:

```go
type Generator interface {
    Generate(diff Difficulty) *Question
    Label() string
}
```

Generators self-register via `init()` functions, so adding a new generator only requires creating a new file in `internal/game/gen/` — no other code needs to change.

The core generation loop (`TryGenerate`) works like this:

1. Look up the weighted patterns for the given difficulty
2. Loop up to 100 attempts
3. Each attempt: pick a random pattern (weighted), call it, check if it produced a valid expression
4. If valid, build a `Question` from the expression tree and return it
5. Return `nil` if all attempts fail

### Pattern System

Each generator defines a `PatternSet` — a map from difficulty level to a list of weighted patterns:

```
PatternSet = map[Difficulty][]WeightedPattern
```

A `WeightedPattern` pairs a generation function with a selection weight. Higher weights mean that pattern is picked more often. For example, at Easy difficulty, Addition might have:

- `addTwoOperands` (weight 8) — 80% chance
- `addThreeOperands` (weight 2) — 20% chance

At Expert difficulty, the weights shift toward more complex patterns.

### Operand Ranges

Each generator has range tables that scale with difficulty. For example, Addition ranges:

| Difficulty | Operand Range |
|------------|---------------|
| Beginner | 1–9 |
| Easy | 10–50 |
| Medium | 20–200 |
| Hard | 100–500 |
| Expert | 200–999 |

Multi-operand patterns use a primary range for the first operand and a smaller secondary range for subsequent operands to keep results manageable.

All range tables are defined in `internal/game/gen/ranges.go`.

---

## The 16 Generators

### Sprint Modes (Single-Operation)

#### Addition

Patterns: 2, 3, 4, or 5 operands chained with `+`. Higher difficulties increase both the number ranges and the likelihood of multi-operand expressions. At Expert, there's a 10% chance of a five-operand expression like `a + b + c + d + e`.

#### Subtraction

Patterns: 2–5 operand chains, plus mixed add/sub patterns like `a + b - c` and `a - b + c - d`.

At Beginner and Easy, results are guaranteed non-negative (the subtrahend is constrained to be smaller). From Medium onward, negative results are allowed.

#### Multiplication

Patterns: 2, 3, or 4 operand chains. Uses asymmetric ranges — one factor can be large while the other stays small to keep products reasonable. Multi-operand patterns use a smaller "multi range" for the 3rd/4th factors.

#### Division

**Uses backward generation** to guarantee clean integer results: picks a divisor and quotient first, then computes `dividend = divisor × quotient`.

Chain division (`a / b / c`) picks a final quotient, two divisors, and computes the dividend as `quotient × d1 × d2`.

#### Square / Cube

Single value patterns (`n²`, `n³`) plus composite patterns at higher difficulties:
- `n² + m²`, `n² - m²`, `n² + m² + p²`
- `n³ + m³`, `n³ - m³`

Composite subtraction patterns ensure the first value is larger than the second.

#### Square Root / Cube Root

**Uses backward generation** for perfect results: picks the root value, then computes the radicand (`n²` or `n³`). Composite patterns similarly guarantee clean results.

#### Power (Exponents)

Patterns: `a^n` with increasing base and exponent ranges. Composite patterns add/subtract two powers (`a^n + b^m`). All patterns check for overflow against a 1,000,000 cap.

| Difficulty | Base Range | Exponent Range |
|------------|-----------|----------------|
| Beginner | 2–10 | 2 |
| Easy | 2–12 | 2–3 |
| Medium | 2–10 | 2–4 |
| Hard | 2–8 | 3–5 |
| Expert | 2–6 | 4–6 |

#### Modulo (Remainders)

Single pattern (`a mod b`) at most difficulties. Expert adds a composite `(a mod b) + c`. Ensures the dividend is always larger than the divisor.

#### Percentage

Uses curated percent pools per difficulty:

| Difficulty | Example Percentages |
|------------|-------------------|
| Beginner/Easy | 10%, 20%, 25%, 50%, 100% |
| Medium | 5%, 10%, 15%, 20%, 25%, 30%, 40%, 50%, 75% |
| Hard | 5%, 10%, 12%, 15%, ..., 80% (14 values) |
| Expert | 2%, 3%, 4%, ..., 39% (30 non-round values) |

Uses `AlignToCleanDivision` to adjust the base value so the result is always a clean integer.

#### Factorial

Patterns: single `n!`, division `n! / m!`, and addition `n! + m!`. The division gap between n and m is capped (3–4) to keep results reasonable. All patterns check against a 1,000,000 overflow cap.

| Difficulty | Range |
|------------|-------|
| Beginner | 1!–4! |
| Easy | 3!–5! |
| Medium | 4!–6! |
| Hard | 5!–8! |
| Expert | 7!–10! |

### Challenge Modes (Mixed)

#### Mixed Basics

Combines +, -, ×, ÷ in progressively complex patterns:

| Difficulty | Key Patterns |
|------------|-------------|
| Beginner | Single random operation, same-precedence chains (`a + b - c`) |
| Easy | Parenthesized expressions like `(a + b) × c` |
| Medium | PEMDAS expressions like `a + b × c` |
| Hard | 3–4 operations with PEMDAS (`a + b × c - d / e`) |
| Expert | 4–5 operations, parallel mul/div (`a × b + c / d`) |

Division in mixed patterns always uses backward generation for clean results.

#### Mixed Powers

Randomly combines squares, cubes, square roots, and cube roots:

| Difficulty | Key Patterns |
|------------|-------------|
| Beginner | Single random power/root operation |
| Easy | Single or simple composite (`n² + m²`) |
| Medium | Sum/difference of two power terms |
| Hard | Includes multiplication of power terms |
| Expert | Three-term composites (`n² + ∛m + p³`) |

#### Mixed Advanced

Combines modulo, factorial, percentage, and power operations:

| Difficulty | Key Patterns |
|------------|-------------|
| Beginner | Single random advanced operation |
| Easy | `n! ± m` composites |
| Medium | `n! / m!`, `n! + a^n`, `a mod b + c` |
| Hard | More composite patterns |
| Expert | Three-term expressions like `n!/m! ± a²` and `a^n mod b ± c!` |

#### Anything Goes

A meta-generator that delegates to other generators rather than having its own patterns:

| Difficulty | Selection Logic |
|------------|----------------|
| Beginner | 100% single-operation generators |
| Easy | 70% single-op, 30% mixed |
| Medium | 40% Mixed Basics, 30% single-op, 30% powers/advanced |
| Hard | 50% Mixed Basics, 25% Mixed Powers, 25% Mixed Advanced |
| Expert | 40% Mixed Basics, 30% Mixed Powers, 20% Mixed Advanced, 10% single-op |

---

## Difficulty System

Five levels with increasing complexity:

| Level | Scoring Multiplier | Effect |
|-------|-------------------|--------|
| Beginner | 0.5x | Smallest numbers, simplest patterns |
| Easy | 0.75x | Slightly larger numbers, occasional multi-operand |
| Medium | 1.0x | Moderate numbers, regular multi-operand |
| Hard | 1.5x | Large numbers, complex patterns frequent |
| Expert | 2.0x | Largest numbers, most complex patterns dominant |

Difficulty affects two things in each generator:
1. **Operand ranges** — larger numbers at higher levels
2. **Pattern weights** — complex patterns (multi-operand, composite) become more likely

---

## Question Pool & Deduplication

The `QuestionPool` sits between generators and the session:

- **Batch size:** 50 questions pre-generated at a time
- **Dedup:** A `seen` map (keyed by expression `Key()`) prevents duplicate questions within a session
- **Refill:** When all 50 are consumed, generates a new batch (up to 150 attempts to fill 50 slots, skipping dupes)
- **Exhaustion recovery:** If all generated questions are duplicates, clears the seen map and tries again
- **Shuffle:** Each batch is shuffled before serving

---

## Session Lifecycle

1. **Create** — `NewSession(generator, difficulty, duration)` builds a question pool
2. **Start** — Sets the timer, loads the first question
3. **Play loop** — Each tick updates `TimeLeft`. Player submits answers or skips
4. **Submit** — Checks answer, calculates score, records history, loads next question
5. **End** — When `TimeLeft ≤ 0`, the session is finished

### Allowed Durations

30 seconds, 60 seconds (default), 90 seconds, 2 minutes.

---

## Scoring System

### Points Formula

```
points = 100 × difficultyMultiplier × timeBonus × streakBonus
```

### Time Bonus

| Response Time | Multiplier |
|--------------|------------|
| Under 2 seconds | 1.5x |
| 2–10 seconds | Linear decay from 1.5x to 1.0x |
| Over 10 seconds | 1.0x |

### Streak Bonus

| Streak | Multiplier | Label |
|--------|------------|-------|
| 0–4 | 1.0x | — |
| 5–9 | 1.25x | STREAK |
| 10–14 | 1.5x | MAX |
| 15–19 | 1.75x | BLAZING |
| 20–24 | 2.0x | UNSTOPPABLE |
| 25+ | 2.0x | LEGENDARY |

### Penalties

- **Wrong answer:** -25 points (score cannot go below 0), streak resets to 0
- **Skip:** 0 points, streak resets to 0

---

## Multiple Choice

When multiple choice input is enabled, `GenerateChoices` produces 4 options (1 correct + 3 distractors):

- Small answers (< 20): Fixed offsets of 1–5
- Larger answers: Percentage-based offsets (10–30% of the answer)
- Difficulty affects spread: Beginner/Easy distractors are more spread out (1.5x offset), Hard/Expert are tighter (0.7x offset)
- Negative distractors are rejected for non-negative answers

---

## Areas for Improvement

This section highlights areas where the current system could be enhanced. These are opportunities for contributors.

### Question Generation

- **Smarter difficulty scaling:** Currently, difficulty only affects number ranges and pattern weights. Generators could incorporate more nuanced scaling — for example, requiring carrying/borrowing in addition, or using numbers near common mistake boundaries.
- **Adaptive difficulty:** The system uses a fixed difficulty level per session. An adaptive mode that adjusts based on player performance (accuracy, response time) could improve engagement.
- **More pattern variety:** Several generators (Modulo, Percentage) only have a single pattern across most difficulties. Adding composite patterns (e.g., `a mod b + c mod d`, or chained percentages) would increase variety.
- **Better distractor generation:** Multiple choice distractors use simple offset-based algorithms. Distractors based on common mistakes (e.g., for `5 + 3 × 2`, offering `16` as `(5+3)×2`) would be more educationally valuable.
- **Decimal/fraction support:** The expression tree only handles integers. Supporting decimal or fractional results could expand question types significantly.
- **Negative number operations:** Currently only subtraction can produce negative intermediate values. Generators could incorporate negative operands explicitly.

### Per-Mode Opportunities

| Mode | Potential Improvements |
|------|----------------------|
| Addition | Carry-focused problems, number bond patterns, complement-to-100 |
| Subtraction | Borrowing-focused problems, complement subtraction |
| Multiplication | Times table focused drills, multiplying by 10/100/1000 |
| Division | Remainders as a variant, long division patterns |
| Squares | Perfect square recognition (reverse: "what number squared is 144?") |
| Cubes | Sums of consecutive cubes, recognizing cube patterns |
| Square/Cube Roots | Estimation questions for non-perfect squares |
| Power | Patterns like `2^n` sequences, power of 10 |
| Modulo | Clock arithmetic contexts, modular addition |
| Percentage | Chained percentages (20% of 50% of 200), percentage increase/decrease, tip calculation |
| Factorial | Permutation/combination word problems |
| Mixed Basics | Longer expression chains, nested parentheses |
| Mixed Powers | Pythagorean-style (`a² + b² = ?`) |
| Mixed Advanced | Cross-category composites (factorial mod n) |
| Anything Goes | Difficulty-weighted category selection could be tuned |

### Scoring & Session

- **Score balancing:** The current formula hasn't been extensively playtested across all modes and difficulties. Some modes may consistently award more/fewer points than others.
- **Partial credit:** Currently answers are binary (correct/incorrect). For estimation-type questions, partial credit based on proximity could work.
- **Session variety:** Only timed sessions exist. Untimed practice, question-count-based sessions, or challenge modes (no wrong answers allowed) could add variety.
- **Streak system tuning:** The streak thresholds and multipliers are somewhat arbitrary. Playtesting data could inform better values.

### Technical

- **Generator testing:** Individual generators have limited test coverage. Property-based tests (e.g., "division always produces integer results") would increase confidence.
- **Performance:** `TryGenerate` can waste attempts when valid expressions are rare at extreme difficulties. A more targeted generation approach (constraint solving) could help.
- **Question quality validation:** No system-level check ensures generated questions are "interesting" (e.g., avoiding trivial `1 + 1` at Beginner or astronomically large results at Expert).
