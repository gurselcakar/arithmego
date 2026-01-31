package game

import (
	"math/rand"
	"sort"
)

// Distractor generation tuning constants
const (
	// Answers below this threshold use fixed offsets; above use percentage-based
	smallAnswerThreshold = 20
	// For small answers, offset randomly by 1 to this value
	smallAnswerMaxOffset = 5
	// For larger answers, offset by this percentage range of the answer
	minOffsetPercent = 10
	maxOffsetPercent = 30
	// Difficulty multipliers: easy = more obvious distractors, hard = trickier
	easyOffsetMultiplier = 1.5
	hardOffsetMultiplier = 0.7
)

// GenerateChoices creates 4 multiple choice options for an answer.
// Returns choices in shuffled order and the correct answer's index (0-3).
func GenerateChoices(answer int, difficulty Difficulty) (choices []int, correctIndex int) {
	choices = make([]int, 4)
	choices[0] = answer

	// Generate 3 distractors
	distractors := generateDistractors(answer, difficulty)
	copy(choices[1:], distractors)

	// Shuffle and find correct index
	rand.Shuffle(len(choices), func(i, j int) {
		choices[i], choices[j] = choices[j], choices[i]
	})

	for i, c := range choices {
		if c == answer {
			correctIndex = i
			break
		}
	}

	return choices, correctIndex
}

// generateDistractors creates 3 unique distractor values near the correct answer.
func generateDistractors(answer int, difficulty Difficulty) []int {
	distractors := make(map[int]bool)
	attempts := 0
	maxAttempts := 100

	for len(distractors) < 3 && attempts < maxAttempts {
		attempts++
		d := generateDistractor(answer, difficulty)

		// Skip duplicates and the correct answer
		if d == answer || distractors[d] {
			continue
		}

		// For non-negative answers, skip negative distractors
		if answer >= 0 && d < 0 {
			continue
		}

		distractors[d] = true
	}

	// Fallback: if we couldn't generate enough distractors, use simple offsets
	fallbackOffsets := []int{1, 2, 3, -1, -2, -3, 4, 5, -4, -5}
	for _, offset := range fallbackOffsets {
		if len(distractors) >= 3 {
			break
		}
		d := answer + offset
		if d != answer && !distractors[d] && (answer < 0 || d >= 0) {
			distractors[d] = true
		}
	}

	// Guaranteed fallback: incrementing positive offsets until we have 3
	// This handles edge cases like answer=0 where negative distractors are rejected
	offset := 1
	for len(distractors) < 3 {
		d := answer + offset
		// Same check as above: no negative distractors for non-negative answers
		if d != answer && !distractors[d] && (answer < 0 || d >= 0) {
			distractors[d] = true
		}
		offset++
		// Safety limit to prevent infinite loop (should never be reached)
		if offset > 1000 {
			break
		}
	}

	result := make([]int, 0, 3)
	for d := range distractors {
		result = append(result, d)
	}

	// Sort for consistent ordering before shuffle
	sort.Ints(result)

	// Return exactly 3 distractors
	if len(result) > 3 {
		result = result[:3]
	}

	return result
}

// generateDistractor creates a single distractor value based on answer magnitude.
func generateDistractor(answer int, difficulty Difficulty) int {
	absAnswer := abs(answer)
	var offset int

	if absAnswer < smallAnswerThreshold {
		offset = rand.Intn(smallAnswerMaxOffset) + 1
	} else {
		percentRange := maxOffsetPercent - minOffsetPercent + 1
		percentage := float64(rand.Intn(percentRange)+minOffsetPercent) / 100.0
		offset = max(1, int(float64(absAnswer)*percentage))
	}

	switch difficulty {
	case Beginner, Easy:
		offset = int(float64(offset) * easyOffsetMultiplier)
	case Hard, Expert:
		offset = max(1, int(float64(offset)*hardOffsetMultiplier))
	}

	// Randomly add or subtract
	if rand.Intn(2) == 0 {
		return answer + offset
	}
	return answer - offset
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
