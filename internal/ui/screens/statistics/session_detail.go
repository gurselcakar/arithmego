package statistics

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

const questionsPerScreen = 10

// QuestionFilter represents the filter for question log.
type QuestionFilter int

const (
	QuestionFilterAll QuestionFilter = iota
	QuestionFilterCorrect
	QuestionFilterWrong
	QuestionFilterSkipped
)

// String returns the display name for the filter.
func (f QuestionFilter) String() string {
	switch f {
	case QuestionFilterCorrect:
		return "Correct"
	case QuestionFilterWrong:
		return "Wrong"
	case QuestionFilterSkipped:
		return "Skipped"
	default:
		return "All"
	}
}

// AllQuestionFilters returns all available question filters.
func AllQuestionFilters() []QuestionFilter {
	return []QuestionFilter{
		QuestionFilterAll,
		QuestionFilterCorrect,
		QuestionFilterWrong,
		QuestionFilterSkipped,
	}
}

// FilterQuestions filters questions based on the filter type.
func FilterQuestions(questions []storage.QuestionRecord, filter QuestionFilter) []storage.QuestionRecord {
	if filter == QuestionFilterAll {
		return questions
	}

	var result []storage.QuestionRecord
	for _, q := range questions {
		switch filter {
		case QuestionFilterCorrect:
			if q.Correct && !q.Skipped {
				result = append(result, q)
			}
		case QuestionFilterWrong:
			if !q.Correct && !q.Skipped {
				result = append(result, q)
			}
		case QuestionFilterSkipped:
			if q.Skipped {
				result = append(result, q)
			}
		}
	}
	return result
}

// RenderSessionSummaryContent renders the session detail summary content for viewport.
func RenderSessionSummaryContent(session storage.SessionRecord, width int) string {
	var b strings.Builder

	// Title
	title := fmt.Sprintf("SESSION · %s", session.Mode)
	b.WriteString(styles.Bold.Render(title))
	b.WriteString("\n\n")

	// Session metadata
	dateStr := FormatSessionDate(session.Timestamp) + ", " + FormatTime(session.Timestamp)
	b.WriteString(styles.Dim.Render("Date: " + dateStr))
	b.WriteString("\n")
	metaLine := fmt.Sprintf("Mode: %s  •  Difficulty: %s  •  %s",
		session.Mode, session.Difficulty, FormatDuration(session.DurationSeconds))
	b.WriteString(styles.Dim.Render(metaLine))
	b.WriteString("\n\n")

	// Separator
	separatorWidth := 56
	if width > 0 && width-10 < separatorWidth {
		separatorWidth = width - 10
	}
	b.WriteString(styles.Dim.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n\n")

	// Results section
	b.WriteString(styles.Bold.Render("RESULTS"))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render("───────"))
	b.WriteString("\n")

	labelWidth := 16
	valueWidth := 18 // Fixed width for consistent alignment

	b.WriteString(fmt.Sprintf("%-*s %-*s\n", labelWidth, "Score", valueWidth, fmt.Sprintf("%d points", session.Score)))
	b.WriteString(fmt.Sprintf("%-*s %-*s\n", labelWidth, "Questions", valueWidth, fmt.Sprintf("%d attempted", session.QuestionsAttempted)))

	accuracy := float64(0)
	if session.QuestionsAttempted > 0 {
		accuracy = float64(session.QuestionsCorrect) / float64(session.QuestionsAttempted) * 100
	}
	b.WriteString(fmt.Sprintf("%-*s %-*s\n", labelWidth, "Correct", valueWidth, fmt.Sprintf("%d (%s)", session.QuestionsCorrect, FormatAccuracy(accuracy))))
	b.WriteString(fmt.Sprintf("%-*s %-*d\n", labelWidth, "Wrong", valueWidth, session.QuestionsWrong))
	b.WriteString(fmt.Sprintf("%-*s %-*d\n", labelWidth, "Skipped", valueWidth, session.QuestionsSkipped))
	b.WriteString(fmt.Sprintf("%-*s %-*d\n", labelWidth, "Best Streak", valueWidth, session.BestStreak))
	if session.AvgResponseTimeMs > 0 {
		b.WriteString(fmt.Sprintf("%-*s %-*s\n", labelWidth, "Avg Time", valueWidth, FormatResponseTime(session.AvgResponseTimeMs)))
	}
	b.WriteString("\n")

	// Mistakes section
	mistakes := getMistakes(session)
	b.WriteString(styles.Bold.Render(fmt.Sprintf("MISTAKES (%d)", len(mistakes))))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render("────────────"))
	b.WriteString("\n")

	if len(mistakes) == 0 && session.QuestionsAttempted > 0 {
		b.WriteString(styles.Correct.Render("Perfect session - no mistakes!"))
		b.WriteString("\n")
	} else if len(mistakes) == 0 {
		b.WriteString(styles.Dim.Render("No questions attempted"))
		b.WriteString("\n")
	} else {
		// Show all mistakes - viewport handles scrolling
		for _, m := range mistakes {
			// Format: "#12  33 + 28 = ?    You: 60   Correct: 61    2.3s"
			line := fmt.Sprintf("#%-3d %-20s You: %-5d Correct: %-5d %s",
				m.Index,
				m.Question,
				m.UserAnswer,
				m.CorrectAnswer,
				FormatResponseTime(m.ResponseTimeMs),
			)
			b.WriteString(styles.Incorrect.Render(line))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// mistake represents a wrong answer with its index.
type mistake struct {
	Index          int
	Question       string
	UserAnswer     int
	CorrectAnswer  int
	ResponseTimeMs int64
}

// getMistakes extracts mistakes from a session.
func getMistakes(session storage.SessionRecord) []mistake {
	var mistakes []mistake
	for i, q := range session.Questions {
		if !q.Correct && !q.Skipped {
			mistakes = append(mistakes, mistake{
				Index:          i + 1,
				Question:       q.Question,
				UserAnswer:     q.UserAnswer,
				CorrectAnswer:  q.CorrectAnswer,
				ResponseTimeMs: q.ResponseTimeMs,
			})
		}
	}
	return mistakes
}

// RenderSessionFullLogContent renders the full question log content for viewport.
func RenderSessionFullLogContent(
	session storage.SessionRecord,
	filter QuestionFilter,
	width int,
) string {
	var b strings.Builder

	// Title with full breadcrumb
	title := fmt.Sprintf("SESSION · %s · Full Log", session.Mode)
	b.WriteString(styles.Bold.Render(title))
	b.WriteString("\n\n")

	// Filter selector
	filterLine := fmt.Sprintf("Filter: ◀ %s ▶    (All / Correct / Wrong / Skipped)", filter.String())
	b.WriteString(styles.Dim.Render(filterLine))
	b.WriteString("\n\n")

	// Column headers
	headerLine := "  #    Question              Answer   Correct   Time"
	b.WriteString(styles.Dim.Render(headerLine))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render(strings.Repeat("─", lipgloss.Width(headerLine)+10)))
	b.WriteString("\n")

	// Filter questions
	filteredQuestions := FilterQuestions(session.Questions, filter)

	if len(filteredQuestions) == 0 {
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render("No questions match this filter."))
		b.WriteString("\n")
	} else {
		// Render all questions - viewport handles scrolling
		for i, q := range filteredQuestions {
			b.WriteString(renderQuestionRow(q, i+1, filter == QuestionFilterAll))
			b.WriteString("\n")
		}

		// Total info
		b.WriteString(styles.Dim.Render(strings.Repeat("─", lipgloss.Width(headerLine)+10)))
		b.WriteString("\n")
		totalInfo := fmt.Sprintf("Total: %d questions", len(filteredQuestions))
		b.WriteString(lipgloss.Place(lipgloss.Width(headerLine)+10, 1, lipgloss.Center, lipgloss.Center,
			styles.Dim.Render(totalInfo)))
	}

	return b.String()
}

// renderQuestionRow renders a single question row in the full log.
func renderQuestionRow(q storage.QuestionRecord, index int, showIndex bool) string {
	var prefix string
	if showIndex {
		prefix = fmt.Sprintf(" %3d   ", index)
	} else {
		prefix = "       "
	}

	// Question text (truncate if needed)
	questionText := q.Question
	if len(questionText) > 18 {
		questionText = questionText[:15] + "..."
	}
	questionText = fmt.Sprintf("%-18s", questionText)

	// Answer
	answerStr := fmt.Sprintf("%-7d", q.UserAnswer)
	if q.Skipped {
		answerStr = "--     "
	}

	// Result indicator
	var resultStr string
	if q.Skipped {
		resultStr = styles.Dim.Render("⊘")
	} else if q.Correct {
		resultStr = styles.Correct.Render("✓") + "        "
	} else {
		resultStr = styles.Incorrect.Render("✗") + " " + fmt.Sprintf("%-5d", q.CorrectAnswer)
	}

	// Time
	timeStr := FormatResponseTime(q.ResponseTimeMs)
	if q.Skipped {
		timeStr = "--"
	}

	line := fmt.Sprintf("%s%s  %s  %s  %s", prefix, questionText, answerStr, resultStr, timeStr)

	// Highlight wrong answers
	if !q.Correct && !q.Skipped {
		return line + "   ←"
	}

	return line
}

// SessionLogNavigation handles scrolling within the question log.
type SessionLogNavigation struct {
	TotalQuestions int
	ScrollOffset   int
	Filter         QuestionFilter
}

// NewSessionLogNavigation creates a new session log navigation.
func NewSessionLogNavigation(totalQuestions int) SessionLogNavigation {
	return SessionLogNavigation{
		TotalQuestions: totalQuestions,
	}
}

// ScrollUp scrolls up by one line.
func (n *SessionLogNavigation) ScrollUp() {
	if n.ScrollOffset > 0 {
		n.ScrollOffset--
	}
}

// ScrollDown scrolls down by one line.
func (n *SessionLogNavigation) ScrollDown() {
	maxScroll := n.TotalQuestions - questionsPerScreen
	if maxScroll < 0 {
		maxScroll = 0
	}
	if n.ScrollOffset < maxScroll {
		n.ScrollOffset++
	}
}

// PageUp scrolls up by a page.
func (n *SessionLogNavigation) PageUp() {
	n.ScrollOffset -= questionsPerScreen
	if n.ScrollOffset < 0 {
		n.ScrollOffset = 0
	}
}

// PageDown scrolls down by a page.
func (n *SessionLogNavigation) PageDown() {
	maxScroll := n.TotalQuestions - questionsPerScreen
	if maxScroll < 0 {
		maxScroll = 0
	}
	n.ScrollOffset += questionsPerScreen
	if n.ScrollOffset > maxScroll {
		n.ScrollOffset = maxScroll
	}
}

// NextFilter cycles to the next filter.
func (n *SessionLogNavigation) NextFilter() {
	n.Filter = (n.Filter + 1) % 4
	n.ScrollOffset = 0 // Reset scroll on filter change
}

// PrevFilter cycles to the previous filter.
func (n *SessionLogNavigation) PrevFilter() {
	if n.Filter == 0 {
		n.Filter = 3
	} else {
		n.Filter--
	}
	n.ScrollOffset = 0 // Reset scroll on filter change
}

// Reset resets navigation.
func (n *SessionLogNavigation) Reset(totalQuestions int) {
	n.TotalQuestions = totalQuestions
	n.ScrollOffset = 0
	n.Filter = QuestionFilterAll
}
