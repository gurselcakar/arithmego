package screens

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

var errLoadFailed = errors.New("failed to load statistics")

// StatisticsModel represents the statistics screen.
type StatisticsModel struct {
	width      int
	height     int
	stats      *storage.Statistics
	aggregates storage.Aggregates
	detailed   bool // false = summary, true = detailed
	err        error
}

// NewStatistics creates a new statistics model.
func NewStatistics() StatisticsModel {
	return StatisticsModel{}
}

// Init initializes the statistics model and loads data.
func (m StatisticsModel) Init() tea.Cmd {
	return func() tea.Msg {
		return loadStatisticsMsg{}
	}
}

// loadStatisticsMsg triggers statistics loading.
type loadStatisticsMsg struct{}

// statisticsLoadedMsg carries loaded statistics.
type statisticsLoadedMsg struct {
	stats *storage.Statistics
	err   error
}

// Update handles statistics screen input.
func (m StatisticsModel) Update(msg tea.Msg) (StatisticsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case loadStatisticsMsg:
		stats, err := storage.Load()
		return m, func() tea.Msg {
			return statisticsLoadedMsg{stats: stats, err: err}
		}

	case statisticsLoadedMsg:
		m.stats = msg.stats
		m.err = msg.err
		// Defensive: treat nil stats with nil error as an error
		if m.stats == nil && m.err == nil {
			m.err = errLoadFailed
		}
		if m.stats != nil {
			m.aggregates = storage.ComputeAggregates(m.stats)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return ReturnToMenuMsg{}
			}
		case "d", "D":
			if !m.detailed {
				m.detailed = true
			}
			return m, nil
		case "s", "S":
			if m.detailed {
				m.detailed = false
			}
			return m, nil
		}
	}

	return m, nil
}

// View renders the statistics screen.
func (m StatisticsModel) View() string {
	if m.err != nil {
		return m.renderError()
	}

	if m.stats == nil {
		return m.renderLoading()
	}

	var content string
	if m.detailed {
		content = m.renderDetailed()
	} else {
		content = m.renderSummary()
	}

	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

// renderLoading shows a loading state.
func (m StatisticsModel) renderLoading() string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		styles.Bold.Render("STATISTICS"),
		"",
		"Loading...",
	)

	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

// renderError shows an error state.
func (m StatisticsModel) renderError() string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		styles.Bold.Render("STATISTICS"),
		"",
		styles.Incorrect.Render("Error loading statistics"),
		"",
		components.RenderHints([]string{"Esc Back"}),
	)

	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

// renderSummary renders the summary view.
func (m StatisticsModel) renderSummary() string {
	agg := m.aggregates

	// Title
	title := styles.Bold.Render("STATISTICS")

	// Empty state
	if agg.TotalSessions == 0 {
		return lipgloss.JoinVertical(lipgloss.Center,
			title,
			"",
			"",
			styles.Dim.Render("No sessions yet."),
			"",
			styles.Dim.Render("Play a game to see your stats!"),
			"",
			"",
			components.RenderHints([]string{"Esc Back"}),
		)
	}

	// Stats
	sessions := fmt.Sprintf("%d sessions", agg.TotalSessions)
	accuracy := fmt.Sprintf("%.0f%% accuracy", agg.OverallAccuracy)

	separator := styles.Dim.Render(strings.Repeat("─", 25))

	// Per-operation accuracy (basic 4 operations)
	opGrid := m.renderOperationGrid()

	// Hints
	hints := components.RenderHints([]string{"[D] Details", "Esc Back"})

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		sessions,
		accuracy,
		"",
		separator,
		"",
		opGrid,
		"",
		"",
		hints,
	)
}

// renderOperationGrid renders a 2x2 grid of basic operation accuracies.
func (m StatisticsModel) renderOperationGrid() string {
	agg := m.aggregates

	// Helper to format operation accuracy
	formatOp := func(symbol, name string) string {
		if stats, ok := agg.ByOperation[name]; ok && stats.Total > 0 {
			return fmt.Sprintf("%s  %3.0f%%", symbol, stats.Accuracy)
		}
		return fmt.Sprintf("%s   --", symbol)
	}

	// Basic 4 operations in a 2x2 layout
	add := formatOp("+", "Addition")
	sub := formatOp("−", "Subtraction")
	mul := formatOp("×", "Multiplication")
	div := formatOp("÷", "Division")

	row1 := fmt.Sprintf("%-12s  %-12s", add, sub)
	row2 := fmt.Sprintf("%-12s  %-12s", mul, div)

	return lipgloss.JoinVertical(lipgloss.Center, row1, row2)
}

// renderDetailed renders the detailed view.
func (m StatisticsModel) renderDetailed() string {
	agg := m.aggregates

	var b strings.Builder

	// Title
	title := styles.Bold.Render("STATISTICS · DETAILED")
	b.WriteString(title)
	b.WriteString("\n\n\n")

	// Overview section
	b.WriteString(styles.Bold.Render("OVERVIEW"))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render(strings.Repeat("─", 30)))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Total Sessions         %d\n", agg.TotalSessions))
	b.WriteString(fmt.Sprintf("Total Questions        %d\n", agg.TotalQuestions))
	b.WriteString(fmt.Sprintf("Overall Accuracy       %.0f%%\n", agg.OverallAccuracy))
	b.WriteString(fmt.Sprintf("Best Streak            %d\n", agg.BestStreakEver))
	if agg.AvgResponseTimeMs > 0 {
		b.WriteString(fmt.Sprintf("Avg Response Time      %.1fs\n", float64(agg.AvgResponseTimeMs)/1000))
	}
	b.WriteString("\n")

	// By Operation section
	if len(agg.ByOperation) > 0 {
		b.WriteString(styles.Bold.Render("BY OPERATION"))
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render(strings.Repeat("─", 30)))
		b.WriteString("\n")

		// Sort operations for consistent display
		ops := make([]string, 0, len(agg.ByOperation))
		for op := range agg.ByOperation {
			ops = append(ops, op)
		}
		sort.Strings(ops)

		for _, op := range ops {
			stats := agg.ByOperation[op]
			symbol := operationSymbol(op)
			b.WriteString(fmt.Sprintf("%s  %-16s %3.0f%%   (%d correct)\n",
				symbol, op, stats.Accuracy, stats.Correct))
		}
		b.WriteString("\n")
	}

	// By Mode section
	if len(agg.ByMode) > 0 {
		b.WriteString(styles.Bold.Render("BY MODE"))
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render(strings.Repeat("─", 30)))
		b.WriteString("\n")

		// Sort modes by session count (descending)
		type modeCount struct {
			name  string
			count int
		}
		modes := make([]modeCount, 0, len(agg.ByMode))
		for name, count := range agg.ByMode {
			modes = append(modes, modeCount{name, count})
		}
		sort.Slice(modes, func(i, j int) bool {
			return modes[i].count > modes[j].count
		})

		for _, mc := range modes {
			plural := "sessions"
			if mc.count == 1 {
				plural = "session"
			}
			b.WriteString(fmt.Sprintf("%-20s %d %s\n", mc.name, mc.count, plural))
		}
		b.WriteString("\n")
	}

	// Hints
	b.WriteString("\n")
	b.WriteString(components.RenderHints([]string{"[S] Summary", "Esc Back"}))

	return b.String()
}

// SetSize sets the screen dimensions.
func (m *StatisticsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// operationSymbol returns the symbol for an operation name.
func operationSymbol(name string) string {
	symbols := map[string]string{
		"Addition":       "+",
		"Subtraction":    "−",
		"Multiplication": "×",
		"Division":       "÷",
		"Square":         "²",
		"Cube":           "³",
		"Square Root":    "√",
		"Cube Root":      "∛",
		"Modulo":         "%",
		"Power":          "^",
		"Percentage":     "%",
		"Factorial":      "!",
	}
	if s, ok := symbols[name]; ok {
		return s
	}
	return "?"
}
