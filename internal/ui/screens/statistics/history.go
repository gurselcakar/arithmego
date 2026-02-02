package statistics

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

const defaultSessionsPerPage = 10

// RenderHistoryContent renders the history view content for viewport.
func RenderHistoryContent(
	sessions []storage.SessionRecord,
	selectedIdx int,
	currentPage int,
	sessionsPerPage int,
	filterPanel FilterPanelModel,
	width int,
) string {
	if sessionsPerPage <= 0 {
		sessionsPerPage = defaultSessionsPerPage
	}

	var b strings.Builder

	// Title
	b.WriteString(styles.Bold.Render("STATISTICS · HISTORY"))
	b.WriteString("\n\n")

	// Filter status
	filterLine := fmt.Sprintf("%s  •  %s  •  %s",
		filterPanel.GetCategoryDisplay(),
		filterPanel.GetDifficultyDisplay(),
		filterPanel.GetTimePeriodDisplay(),
	)
	b.WriteString(styles.Dim.Render(filterLine))
	b.WriteString("\n\n")

	// Separator
	separatorWidth := 60
	if width > 0 && width-10 < separatorWidth {
		separatorWidth = width - 10
	}
	b.WriteString(styles.Dim.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n\n")

	// Empty state
	if len(sessions) == 0 {
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render("No sessions found."))
		b.WriteString("\n\n")
		b.WriteString(styles.Dim.Render("Try different filters or play more games!"))
		b.WriteString("\n")

		return b.String()
	}

	// Column headers
	headerLine := fmt.Sprintf("     %-35s  %5s  %4s  %6s",
		"SESSION", "SCORE", "ACC", "STREAK")
	b.WriteString(styles.Dim.Render(headerLine))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render(strings.Repeat("─", lipgloss.Width(headerLine))))
	b.WriteString("\n")

	// Get visible sessions for current page
	start := currentPage * sessionsPerPage
	if start > len(sessions) {
		start = len(sessions)
	}
	end := start + sessionsPerPage
	if end > len(sessions) {
		end = len(sessions)
	}
	visibleSessions := sessions[start:end]

	// Group by date and render
	currentDate := ""
	for i, session := range visibleSessions {
		globalIdx := start + i
		dateStr := FormatSessionDate(session.Timestamp)

		// Date group header
		if dateStr != currentDate {
			if currentDate != "" {
				b.WriteString("\n")
			}
			currentDate = dateStr
			b.WriteString(styles.Dim.Render(dateStr))
			b.WriteString("\n")
		}

		// Session row
		b.WriteString(renderSessionRow(session, globalIdx == selectedIdx))
		b.WriteString("\n")
	}

	// Pagination
	totalPages := (len(sessions) + sessionsPerPage - 1) / sessionsPerPage
	if totalPages > 1 {
		b.WriteString("\n")
		pageInfo := fmt.Sprintf("Page %d of %d", currentPage+1, totalPages)
		b.WriteString(lipgloss.Place(separatorWidth, 1, lipgloss.Center, lipgloss.Center, styles.Dim.Render(pageInfo)))
	}

	return b.String()
}

// renderSessionRow renders a single session row.
func renderSessionRow(session storage.SessionRecord, selected bool) string {
	// Selection indicator
	prefix := "  "
	if selected {
		prefix = styles.Accent.Render("▸ ")
	}

	// Format: "▸ Addition (Medium)                248    92%      8"
	modeInfo := fmt.Sprintf("%s (%s)", session.Mode, session.Difficulty)
	if len(modeInfo) > 30 {
		modeInfo = modeInfo[:30]
	}

	accuracy := float64(0)
	if session.QuestionsAttempted > 0 {
		accuracy = float64(session.QuestionsCorrect) / float64(session.QuestionsAttempted) * 100
	}

	line := fmt.Sprintf("%s%-32s  %5d  %3.0f%%  %6d",
		prefix,
		modeInfo,
		session.Score,
		accuracy,
		session.BestStreak,
	)

	if selected {
		return styles.Bold.Render(line)
	}
	return line
}

// HistoryNavigation handles navigation within the history view.
type HistoryNavigation struct {
	TotalSessions   int
	SessionsPerPage int
	CurrentPage     int
	SelectedIndex   int // Global index (not page-relative)
}

// NewHistoryNavigation creates a new history navigation.
func NewHistoryNavigation(totalSessions int) HistoryNavigation {
	return HistoryNavigation{
		TotalSessions:   totalSessions,
		SessionsPerPage: defaultSessionsPerPage,
	}
}

// TotalPages returns the total number of pages.
func (n HistoryNavigation) TotalPages() int {
	if n.TotalSessions == 0 || n.SessionsPerPage == 0 {
		return 1
	}
	return (n.TotalSessions + n.SessionsPerPage - 1) / n.SessionsPerPage
}

// PageRelativeIndex returns the index within the current page.
func (n HistoryNavigation) PageRelativeIndex() int {
	return n.SelectedIndex - (n.CurrentPage * n.SessionsPerPage)
}

// MoveUp moves selection up, handling page boundaries.
func (n *HistoryNavigation) MoveUp() {
	if n.SelectedIndex > 0 {
		n.SelectedIndex--
		// Adjust page if needed
		if n.SessionsPerPage > 0 && n.SelectedIndex < n.CurrentPage*n.SessionsPerPage {
			n.CurrentPage--
		}
	}
}

// MoveDown moves selection down, handling page boundaries.
func (n *HistoryNavigation) MoveDown() {
	if n.SelectedIndex < n.TotalSessions-1 {
		n.SelectedIndex++
		// Adjust page if needed
		if n.SessionsPerPage > 0 && n.SelectedIndex >= (n.CurrentPage+1)*n.SessionsPerPage {
			n.CurrentPage++
		}
	}
}

// NextPage moves to the next page.
func (n *HistoryNavigation) NextPage() {
	if n.CurrentPage < n.TotalPages()-1 {
		n.CurrentPage++
		n.SelectedIndex = n.CurrentPage * n.SessionsPerPage
	}
}

// PrevPage moves to the previous page.
func (n *HistoryNavigation) PrevPage() {
	if n.CurrentPage > 0 {
		n.CurrentPage--
		n.SelectedIndex = n.CurrentPage * n.SessionsPerPage
	}
}

// Reset resets navigation to the beginning.
func (n *HistoryNavigation) Reset(totalSessions int) {
	n.TotalSessions = totalSessions
	n.CurrentPage = 0
	n.SelectedIndex = 0
	if n.SessionsPerPage == 0 {
		n.SessionsPerPage = defaultSessionsPerPage
	}
}
