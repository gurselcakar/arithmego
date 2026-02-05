package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

const (
	streakBarWidth = 10
)

// Streak bar characters for each tier
var (
	emptyChar    = "░"
	filledChar   = "█"
	legendaryBar = "◆◆◆◆◆◆◆◆◆◆"
)

// RenderMultiplier renders the streak multiplier display (e.g., "×1.5").
// The multiplier is styled with the Multiplier style (yellow).
func RenderMultiplier(multiplier float64) string {
	text := fmt.Sprintf("×%.1f", multiplier)
	return styles.Multiplier.Render(text)
}

// RenderStreakBar renders a progress bar showing progress toward the next tier.
// The bar fills within each tier (5 streaks = full bar), then resets when a
// new multiplier milestone is reached. Styling upgrades with each tier:
// - TierNone (0): empty bar, dim
// - TierBuilding (1-4): filling, white
// - TierStreak (5-9): filling, green
// - TierMax (10-14): filling, bold green
// - TierBlazing (15-19): filling, bold yellow
// - TierUnstoppable (20-24): filling, bold magenta, «» brackets
// - TierLegendary (25+): full diamond bar, final form
func RenderStreakBar(streak int, tick int) string {
	tier := game.GetStreakTier(streak)

	if tier == game.TierNone {
		return styles.StreakNone.Render("[" + strings.Repeat(emptyChar, streakBarWidth) + "]")
	}

	if tier == game.TierLegendary {
		return styles.StreakLegendary.Render("<" + legendaryBar + ">")
	}

	// Progress within current tier: each tier spans 5 streaks, bar has 10 slots
	// so each correct answer fills 2 segments
	progress := streak % game.StreakMilestoneSize
	filled := progress * (streakBarWidth / game.StreakMilestoneSize)
	empty := streakBarWidth - filled
	bar := strings.Repeat(filledChar, filled) + strings.Repeat(emptyChar, empty)

	open, close := "[", "]"
	var style lipgloss.Style
	switch tier {
	case game.TierBuilding:
		style = styles.StreakBuilding
	case game.TierStreak:
		style = styles.StreakActive
	case game.TierMax:
		style = styles.StreakMax
	case game.TierBlazing:
		style = styles.StreakBlazing
	case game.TierUnstoppable:
		style = styles.StreakUnstoppable
		open, close = "«", "»"
	default:
		style = styles.StreakNone
	}

	return style.Render(open + bar + close)
}


// RenderScore renders the score with comma formatting (e.g., "1,234").
func RenderScore(score int) string {
	text := formatNumber(score)
	return styles.Score.Render(text)
}

// RenderScoreLarge renders the score in a prominent style for the game screen.
// Uses bright white color and comma formatting.
func RenderScoreLarge(score int) string {
	text := formatNumber(score)
	return styles.ScoreLarge.Render(text)
}

// RenderScoreDelta renders the points gained or lost from the last answer.
// Positive deltas show as green "+150", negative as red "-25".
// Returns empty string for zero delta (e.g., skipped questions).
func RenderScoreDelta(delta int) string {
	var text string
	var style lipgloss.Style

	if delta > 0 {
		text = fmt.Sprintf("+%d", delta)
		style = styles.Correct.Bold(true)
	} else if delta < 0 {
		text = fmt.Sprintf("%d", delta)
		style = styles.Incorrect.Bold(true)
	} else {
		return ""
	}

	return style.Render(text)
}

// RenderScoreboard renders the complete scoreboard (multiplier + streak bar).
// This is the top-left component shown during gameplay.
func RenderScoreboard(streak int, tick int) string {
	multiplier := game.StreakBonus(streak)

	mult := RenderMultiplier(multiplier)
	if streak > 0 {
		mult += styles.Dim.Render(" · ") + styles.Multiplier.Render(fmt.Sprintf("Streak %d", streak))
	}
	bar := RenderStreakBar(streak, tick)

	return lipgloss.JoinVertical(lipgloss.Left, mult, bar)
}

// formatNumber formats a number with comma separators.
func formatNumber(n int) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}

	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return sign + str
	}

	var result strings.Builder
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}
	return sign + result.String()
}
