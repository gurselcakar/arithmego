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
	blazingChar  = "▓"
	legendaryBar = "◆◆◆◆◆◆◆◆◆◆"
)

// RenderMultiplier renders the streak multiplier display (e.g., "×1.5").
// The multiplier is styled with the Multiplier style (yellow).
func RenderMultiplier(multiplier float64) string {
	text := fmt.Sprintf("×%.1f", multiplier)
	return styles.Multiplier.Render(text)
}

// RenderStreakBar renders a progress bar that evolves based on streak tier.
// The bar transforms visually as the player maintains higher streaks:
// - TierNone (0): empty bar
// - TierBuilding (1-4): filling bar, dim
// - TierStreak (5-9): filling bar, bright
// - TierMax (10-14): full bar with count
// - TierBlazing (15-19): shimmer animation
// - TierUnstoppable (20-24): different brackets, shimmer
// - TierLegendary (25+): diamond pattern, final form
// The tick parameter drives the shimmer animation for higher tiers.
func RenderStreakBar(streak int, tick int) string {
	tier := game.GetStreakTier(streak)

	switch tier {
	case game.TierNone:
		// Empty bar
		return styles.StreakNone.Render("[" + strings.Repeat(emptyChar, streakBarWidth) + "]")

	case game.TierBuilding:
		// 1-4: filling up, dim
		filled := streak
		empty := streakBarWidth - filled
		bar := strings.Repeat(filledChar, filled) + strings.Repeat(emptyChar, empty)
		return styles.StreakBuilding.Render("[" + bar + "]")

	case game.TierStreak:
		// 5-9: filling up, bright green
		filled := streak
		empty := streakBarWidth - filled
		bar := strings.Repeat(filledChar, filled) + strings.Repeat(emptyChar, empty)
		return styles.StreakActive.Render("[" + bar + "]")

	case game.TierMax:
		// 10-14: full bar with streak count
		bar := strings.Repeat(filledChar, streakBarWidth)
		return styles.StreakMax.Render("["+bar+"]") + styles.Dim.Render(fmt.Sprintf(" %d", streak))

	case game.TierBlazing:
		// 15-19: transformed bar with shimmer
		bar := renderShimmerBar(tick, blazingChar, filledChar)
		return styles.StreakBlazing.Render("["+bar+"]") + styles.Dim.Render(fmt.Sprintf(" %d", streak))

	case game.TierUnstoppable:
		// 20-24: another transformation
		bar := renderShimmerBar(tick, blazingChar, filledChar)
		return styles.StreakUnstoppable.Render("«"+bar+"»") + styles.Dim.Render(fmt.Sprintf(" %d", streak))

	case game.TierLegendary:
		// 25+: final form
		return styles.StreakLegendary.Render("<"+legendaryBar+">") + styles.Dim.Render(fmt.Sprintf(" %d", streak))

	default:
		return ""
	}
}

// renderShimmerBar creates alternating pattern for shimmer effect
func renderShimmerBar(tick int, char1, char2 string) string {
	var bar strings.Builder
	for i := 0; i < streakBarWidth; i++ {
		if (i+tick)%2 == 0 {
			bar.WriteString(char1)
		} else {
			bar.WriteString(char2)
		}
	}
	return bar.String()
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
