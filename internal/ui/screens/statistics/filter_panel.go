package statistics

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gurselcakar/arithmego/internal/analytics"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// FilterField represents which filter field is focused.
type FilterField int

const (
	FilterFieldCategory FilterField = iota
	FilterFieldDifficulty
	FilterFieldTimePeriod
)

// FilterPanelModel manages the filter panel overlay.
type FilterPanelModel struct {
	open         bool
	focusedField FilterField

	// Current selected indices
	categoryIdx   int
	difficultyIdx int
	timePeriodIdx int

	// Pending values (applied on close)
	pendingCategory   int
	pendingDifficulty int
	pendingTimePeriod int

	// Available options
	categories   []string
	difficulties []string
	timePeriods  []analytics.TimePeriod
}

// NewFilterPanel creates a new filter panel model.
func NewFilterPanel() FilterPanelModel {
	return FilterPanelModel{
		categories:   analytics.AllCategories(),
		difficulties: analytics.AllDifficulties(),
		timePeriods:  analytics.AllTimePeriods(),
	}
}

// IsOpen returns whether the panel is open.
func (m FilterPanelModel) IsOpen() bool {
	return m.open
}

// Open opens the filter panel and sets pending values to current.
func (m *FilterPanelModel) Open() {
	m.open = true
	m.focusedField = FilterFieldCategory
	m.pendingCategory = m.categoryIdx
	m.pendingDifficulty = m.difficultyIdx
	m.pendingTimePeriod = m.timePeriodIdx
}

// Close closes the filter panel without applying changes.
func (m *FilterPanelModel) Close() {
	m.open = false
}

// Apply closes the panel and applies the pending changes.
func (m *FilterPanelModel) Apply() {
	m.categoryIdx = m.pendingCategory
	m.difficultyIdx = m.pendingDifficulty
	m.timePeriodIdx = m.pendingTimePeriod
	m.open = false
}

// GetFilters returns the current filter configuration.
func (m FilterPanelModel) GetFilters() analytics.AggregateFilter {
	var category, difficulty string
	var timePeriod analytics.TimePeriod

	if m.categoryIdx >= 0 && m.categoryIdx < len(m.categories) {
		category = m.categories[m.categoryIdx]
	}
	if m.difficultyIdx >= 0 && m.difficultyIdx < len(m.difficulties) {
		difficulty = m.difficulties[m.difficultyIdx]
	}
	if m.timePeriodIdx >= 0 && m.timePeriodIdx < len(m.timePeriods) {
		timePeriod = m.timePeriods[m.timePeriodIdx]
	}

	return analytics.AggregateFilter{
		Category:   category,
		Difficulty: difficulty,
		TimePeriod: timePeriod,
	}
}

// Update handles keyboard input for the filter panel.
func (m FilterPanelModel) Update(msg tea.Msg) (FilterPanelModel, tea.Cmd) {
	if !m.open {
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.focusedField > 0 {
				m.focusedField--
			}
		case "down", "j":
			if m.focusedField < FilterFieldTimePeriod {
				m.focusedField++
			}
		case "left", "h":
			m.decrementField()
		case "right", "l":
			m.incrementField()
		case "tab", "enter":
			m.Apply()
		case "esc":
			m.Close()
		}
	}

	return m, nil
}

func (m *FilterPanelModel) incrementField() {
	switch m.focusedField {
	case FilterFieldCategory:
		if m.pendingCategory < len(m.categories)-1 {
			m.pendingCategory++
		}
	case FilterFieldDifficulty:
		if m.pendingDifficulty < len(m.difficulties)-1 {
			m.pendingDifficulty++
		}
	case FilterFieldTimePeriod:
		if m.pendingTimePeriod < len(m.timePeriods)-1 {
			m.pendingTimePeriod++
		}
	}
}

func (m *FilterPanelModel) decrementField() {
	switch m.focusedField {
	case FilterFieldCategory:
		if m.pendingCategory > 0 {
			m.pendingCategory--
		}
	case FilterFieldDifficulty:
		if m.pendingDifficulty > 0 {
			m.pendingDifficulty--
		}
	case FilterFieldTimePeriod:
		if m.pendingTimePeriod > 0 {
			m.pendingTimePeriod--
		}
	}
}

// View renders the filter panel.
func (m FilterPanelModel) View() string {
	if !m.open {
		return ""
	}

	var b strings.Builder

	// Box border (24 chars interior to fit: 2 + 1 + 1 + 16 + 1 + 1 + 2 = 24)
	b.WriteString("┌────────────────────────┐\n")
	b.WriteString("│        FILTERS         │\n")
	b.WriteString("│                        │\n")

	// Category field
	categoryLabel := "  Category              "
	categoryValue := m.formatSelector(
		analytics.CategoryDisplayName(m.categories[m.pendingCategory]),
		m.focusedField == FilterFieldCategory,
	)
	if m.focusedField == FilterFieldCategory {
		b.WriteString("│" + styles.Bold.Render(categoryLabel) + "│\n")
	} else {
		b.WriteString("│" + styles.Dim.Render(categoryLabel) + "│\n")
	}
	b.WriteString("│" + categoryValue + "│\n")
	b.WriteString("│                        │\n")

	// Difficulty field
	difficultyLabel := "  Difficulty            "
	difficultyValue := m.formatSelector(
		analytics.DifficultyDisplayName(m.difficulties[m.pendingDifficulty]),
		m.focusedField == FilterFieldDifficulty,
	)
	if m.focusedField == FilterFieldDifficulty {
		b.WriteString("│" + styles.Bold.Render(difficultyLabel) + "│\n")
	} else {
		b.WriteString("│" + styles.Dim.Render(difficultyLabel) + "│\n")
	}
	b.WriteString("│" + difficultyValue + "│\n")
	b.WriteString("│                        │\n")

	// Time Period field
	timePeriodLabel := "  Time Period           "
	timePeriodValue := m.formatSelector(
		m.timePeriods[m.pendingTimePeriod].String(),
		m.focusedField == FilterFieldTimePeriod,
	)
	if m.focusedField == FilterFieldTimePeriod {
		b.WriteString("│" + styles.Bold.Render(timePeriodLabel) + "│\n")
	} else {
		b.WriteString("│" + styles.Dim.Render(timePeriodLabel) + "│\n")
	}
	b.WriteString("│" + timePeriodValue + "│\n")
	b.WriteString("│                        │\n")

	b.WriteString("└────────────────────────┘")

	return b.String()
}

// formatSelector formats a value with left/right arrows.
func (m FilterPanelModel) formatSelector(value string, focused bool) string {
	// Pad/truncate value to fit in selector space (16 chars to fit "All Difficulties")
	const maxValueLen = 16
	if len(value) > maxValueLen {
		value = value[:maxValueLen]
	}

	// Center the value
	totalPad := maxValueLen - len(value)
	leftPad := totalPad / 2
	rightPad := totalPad - leftPad
	paddedValue := strings.Repeat(" ", leftPad) + value + strings.Repeat(" ", rightPad)

	// Total width: 2 + 1 + 1 + 16 + 1 + 1 + 2 = 24 (matches box interior)
	if focused {
		return "  " + styles.Accent.Render("◀") + " " + styles.Bold.Render(paddedValue) + " " + styles.Accent.Render("▶") + "  "
	}
	return "  " + styles.Dim.Render("◀") + " " + paddedValue + " " + styles.Dim.Render("▶") + "  "
}

// FilterSummary returns a one-line summary of active filters.
func (m FilterPanelModel) FilterSummary() string {
	var parts []string

	if m.categoryIdx >= 0 && m.categoryIdx < len(m.categories) {
		cat := m.categories[m.categoryIdx]
		if cat == "" {
			parts = append(parts, "All Categories")
		} else {
			parts = append(parts, cat)
		}
	}

	if m.difficultyIdx >= 0 && m.difficultyIdx < len(m.difficulties) {
		diff := m.difficulties[m.difficultyIdx]
		if diff == "" {
			parts = append(parts, "All Difficulties")
		} else {
			parts = append(parts, diff)
		}
	}

	if m.timePeriodIdx >= 0 && m.timePeriodIdx < len(m.timePeriods) {
		period := m.timePeriods[m.timePeriodIdx]
		if period != analytics.TimePeriodAllTime {
			parts = append(parts, period.String())
		}
	}

	return strings.Join(parts, " • ")
}

// OverlayView renders the filter panel overlaid on content.
// The panel appears on the left side of the screen.
func (m FilterPanelModel) OverlayView(content string, width int) string {
	if !m.open {
		return content
	}

	panel := m.View()
	panelLines := strings.Split(panel, "\n")
	contentLines := strings.Split(content, "\n")

	// Ensure content has enough lines
	for len(contentLines) < len(panelLines)+2 {
		contentLines = append(contentLines, "")
	}

	// Overlay panel on content (starting from line 2 to leave room for title)
	var result []string

	for i, line := range contentLines {
		panelLineIdx := i - 2 // Start panel at line 2
		if panelLineIdx >= 0 && panelLineIdx < len(panelLines) {
			// Show panel line on the left, content is obscured
			panelLine := panelLines[panelLineIdx]
			result = append(result, " "+panelLine)
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
