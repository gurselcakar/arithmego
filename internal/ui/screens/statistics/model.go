package statistics

import (
	"errors"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/analytics"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
)

var errLoadFailed = errors.New("failed to load statistics")

// Layout constants
const (
	statisticsHintsHeight = 3
	// minQuestionsForFullLog is the minimum number of questions needed to show the "Full Log" option.
	minQuestionsForFullLog = 10
)

// Model represents the statistics screen with all its views.
type Model struct {
	width  int
	height int

	// Viewport for scrolling
	viewport      viewport.Model
	viewportReady bool

	// Data
	stats      *storage.Statistics
	aggregates analytics.ExtendedAggregates

	// Current view
	view StatisticsView

	// Filter state (shared across views)
	filterPanel FilterPanelModel

	// Operations view state
	operationList     []OperationRow
	operationIndex    int
	selectedOperation string

	// Operation detail state
	opDetailDifficultyIdx int
	opHasMistakes         bool // cached to avoid querying on every render

	// Operation review state (review all mistakes)
	opReviewMistakes []analytics.RecentMistake

	// History view state
	sessionList   []storage.SessionRecord
	historyNav    HistoryNavigation
	filterSummary string

	// Session detail state
	selectedSession   *storage.SessionRecord
	sessionDetailMode SessionDetailMode
	sessionLogNav     SessionLogNavigation

	// Trends view state
	trendsState TrendsState

	// Loading/error state
	loading bool
	err     error
}

// New creates a new statistics model.
func New() Model {
	return Model{
		filterPanel: NewFilterPanel(),
		trendsState: NewTrendsState(),
		viewport:    viewport.New(0, 0),
	}
}

// Init initializes the model and triggers data loading.
func (m Model) Init() tea.Cmd {
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

// ReturnToMenuMsg signals return to main menu.
type ReturnToMenuMsg struct{}

// Update handles all messages for the statistics screen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case loadStatisticsMsg:
		m.loading = true
		stats, err := storage.Load()
		return m, func() tea.Msg {
			return statisticsLoadedMsg{stats: stats, err: err}
		}

	case statisticsLoadedMsg:
		m.loading = false
		m.stats = msg.stats
		m.err = msg.err

		if m.stats == nil && m.err == nil {
			m.err = errLoadFailed
		}

		if m.stats != nil {
			m.aggregates = analytics.ComputeExtendedAggregates(m.stats)
			m.rebuildLists()
		}
		m.updateViewportContent()
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	// Update viewport (for mouse scrolling)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// handleKeyPress routes key presses to the appropriate view handler.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.view {
	case ViewDashboard:
		return m.handleDashboardKeys(msg)
	case ViewOperations:
		return m.handleOperationsKeys(msg)
	case ViewOperationDetail:
		return m.handleOperationDetailKeys(msg)
	case ViewOperationReview:
		return m.handleOperationReviewKeys(msg)
	case ViewHistory:
		return m.handleHistoryKeys(msg)
	case ViewSessionDetail:
		return m.handleSessionDetailKeys(msg)
	case ViewSessionFullLog:
		return m.handleSessionLogKeys(msg)
	case ViewTrends:
		return m.handleTrendsKeys(msg)
	}

	return m, nil
}

// handleDashboardKeys handles dashboard view keys.
func (m Model) handleDashboardKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return m, func() tea.Msg { return ReturnToMenuMsg{} }
	case "o", "O":
		m.view = ViewOperations
		m.operationIndex = 0
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "h", "H":
		m.view = ViewHistory
		m.historyNav.Reset(len(m.sessionList))
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "t", "T":
		m.view = ViewTrends
		m.updateTrendsData()
		m.updateViewportContent()
		m.viewport.GotoTop()
	default:
		// Let viewport handle scrolling (up/down/pgup/pgdown/etc)
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

// handleOperationsKeys handles operations view keys.
func (m Model) handleOperationsKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = ViewDashboard
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "c", "C":
		m.filterPanel.CycleCategory()
		m.applyFilters()
		m.updateViewportContent()
	case "d", "D":
		m.filterPanel.CycleDifficulty()
		m.applyFilters()
		m.updateViewportContent()
	case "p", "P":
		m.filterPanel.CycleTimePeriod()
		m.applyFilters()
		m.updateViewportContent()
	case "up", "k":
		if m.operationIndex > 0 {
			m.operationIndex--
			m.updateViewportContent()
		}
	case "down", "j":
		if m.operationIndex < len(m.operationList)-1 {
			m.operationIndex++
			m.updateViewportContent()
		}
	case "enter":
		if len(m.operationList) > 0 {
			m.selectedOperation = GetSelectedOperation(m.operationList, m.operationIndex)
			m.view = ViewOperationDetail
			m.opDetailDifficultyIdx = 0
			m.opHasMistakes = len(analytics.GetRecentMistakes(m.stats, m.selectedOperation, 1)) > 0
			m.updateViewportContent()
			m.viewport.GotoTop()
		}
	}
	return m, nil
}

// handleOperationDetailKeys handles operation detail view keys.
func (m Model) handleOperationDetailKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = ViewOperations
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "left", "h":
		// Cycle difficulty filter
		if m.opDetailDifficultyIdx > 0 {
			m.opDetailDifficultyIdx--
			m.updateViewportContent()
		}
	case "right", "l":
		diffs := analytics.AllDifficulties()
		if m.opDetailDifficultyIdx < len(diffs)-1 {
			m.opDetailDifficultyIdx++
			m.updateViewportContent()
		}
	case "r", "R":
		// Enter review all mistakes mode
		// Load all mistakes (not just 5)
		m.opReviewMistakes = analytics.GetRecentMistakes(m.stats, m.selectedOperation, 1000)
		if len(m.opReviewMistakes) > 0 {
			m.view = ViewOperationReview
			m.updateViewportContent()
			m.viewport.GotoTop()
		}
	default:
		// Let viewport handle scrolling
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

// handleOperationReviewKeys handles operation review (all mistakes) view keys.
func (m Model) handleOperationReviewKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = ViewOperationDetail
		m.updateViewportContent()
		m.viewport.GotoTop()
	default:
		// Let viewport handle scrolling
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

// handleHistoryKeys handles history view keys.
func (m Model) handleHistoryKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = ViewDashboard
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "c", "C":
		m.filterPanel.CycleCategory()
		m.applyFilters()
		m.updateViewportContent()
	case "d", "D":
		m.filterPanel.CycleDifficulty()
		m.applyFilters()
		m.updateViewportContent()
	case "p", "P":
		m.filterPanel.CycleTimePeriod()
		m.applyFilters()
		m.updateViewportContent()
	case "up", "k":
		m.historyNav.MoveUp()
		m.updateViewportContent()
	case "down", "j":
		m.historyNav.MoveDown()
		m.updateViewportContent()
	case "left":
		m.historyNav.PrevPage()
		m.updateViewportContent()
	case "right":
		m.historyNav.NextPage()
		m.updateViewportContent()
	case "enter":
		if len(m.sessionList) > 0 && m.historyNav.SelectedIndex < len(m.sessionList) {
			m.selectedSession = &m.sessionList[m.historyNav.SelectedIndex]
			m.sessionDetailMode = SessionModeSummary
			m.view = ViewSessionDetail
			m.updateViewportContent()
			m.viewport.GotoTop()
		}
	}
	return m, nil
}

// handleSessionDetailKeys handles session detail view keys.
func (m Model) handleSessionDetailKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "h", "H":
		m.view = ViewHistory
		m.selectedSession = nil
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "l", "L":
		if m.selectedSession != nil && len(m.selectedSession.Questions) > minQuestionsForFullLog {
			m.sessionDetailMode = SessionModeFullLog
			m.sessionLogNav.Reset(len(m.selectedSession.Questions))
			m.view = ViewSessionFullLog
			m.updateViewportContent()
			m.viewport.GotoTop()
		}
	default:
		// Let viewport handle scrolling
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

// handleSessionLogKeys handles session log view keys.
func (m Model) handleSessionLogKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "s", "S":
		m.view = ViewSessionDetail
		m.sessionDetailMode = SessionModeSummary
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "left", "h":
		m.sessionLogNav.PrevFilter()
		// Update total based on filtered count
		if m.selectedSession != nil {
			filtered := FilterQuestions(m.selectedSession.Questions, m.sessionLogNav.Filter)
			m.sessionLogNav.TotalQuestions = len(filtered)
		}
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "right", "l":
		m.sessionLogNav.NextFilter()
		if m.selectedSession != nil {
			filtered := FilterQuestions(m.selectedSession.Questions, m.sessionLogNav.Filter)
			m.sessionLogNav.TotalQuestions = len(filtered)
		}
		m.updateViewportContent()
		m.viewport.GotoTop()
	default:
		// Let viewport handle scrolling
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

// handleTrendsKeys handles trends view keys.
func (m Model) handleTrendsKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = ViewDashboard
		m.updateViewportContent()
		m.viewport.GotoTop()
	case "m", "M":
		m.trendsState.NextMetric()
		m.updateTrendsData()
		m.updateViewportContent()
	case "p", "P":
		m.trendsState.NextPeriod()
		m.updateTrendsData()
		m.updateViewportContent()
	default:
		// Let viewport handle scrolling
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

// applyFilters recomputes data after filter changes.
func (m *Model) applyFilters() {
	if m.stats == nil {
		return
	}

	filter := m.filterPanel.GetFilters()
	m.aggregates = analytics.ComputeFilteredAggregates(m.stats, filter)
	m.rebuildLists()
}

// rebuildLists rebuilds operation and session lists based on current filters.
func (m *Model) rebuildLists() {
	if m.stats == nil {
		return
	}

	filter := m.filterPanel.GetFilters()

	// Rebuild operation list
	m.operationList = BuildOperationList(m.aggregates, m.stats, filter)
	if m.operationIndex >= len(m.operationList) {
		m.operationIndex = 0
	}

	// Rebuild session list
	m.sessionList = analytics.GetSessionsByFilter(m.stats, filter)
	m.historyNav.Reset(len(m.sessionList))
}

// updateTrendsData recomputes trend data.
func (m *Model) updateTrendsData() {
	if m.stats == nil {
		return
	}

	m.trendsState.TrendData = analytics.ComputeTrendData(m.stats, m.trendsState.Period)
	m.trendsState.Insights = analytics.GenerateInsights(m.stats, m.trendsState.Period)
}

// View renders the current view.
func (m Model) View() string {
	if m.err != nil {
		return RenderError(m.width, m.height)
	}

	if m.loading || m.stats == nil {
		return RenderLoading(m.width, m.height)
	}

	if !m.viewportReady {
		return "Loading..."
	}

	hints := m.getHints()
	content := m.renderCurrentViewContent()
	contentHeight := lipgloss.Height(content)
	viewportHeight := m.viewport.Height

	// Center content vertically if it's shorter than the viewport
	var mainArea string
	if contentHeight < viewportHeight {
		mainArea = lipgloss.Place(m.width, viewportHeight, lipgloss.Center, lipgloss.Center, content)
	} else {
		mainArea = m.viewport.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		mainArea,
		lipgloss.Place(m.width, statisticsHintsHeight, lipgloss.Center, lipgloss.Center, hints),
	)
}

// getHints returns the context-aware hints for the current view.
func (m Model) getHints() string {
	switch m.view {
	case ViewDashboard:
		return components.RenderHintsStructured([]components.Hint{
			{Key: "Esc", Action: "Back"},
			{Key: "O", Action: "Operations"},
			{Key: "H", Action: "History"},
			{Key: "T", Action: "Trends"},
		})

	case ViewOperations:
		hintList := []components.Hint{
			{Key: "Esc", Action: "Back"},
			{Key: "c", Action: "Category"},
			{Key: "d", Action: "Difficulty"},
			{Key: "p", Action: "Period"},
		}
		if len(m.operationList) > 0 {
			hintList = append(hintList,
				components.Hint{Key: "↑↓", Action: "Navigate"},
				components.Hint{Key: "Enter", Action: "Details"},
			)
		}
		return components.RenderHintsStructured(hintList)

	case ViewOperationDetail:
		hintList := []components.Hint{
			{Key: "Esc", Action: "Back"},
			{Key: "←→", Action: "Filter"},
		}
		if m.opHasMistakes {
			hintList = append(hintList, components.Hint{Key: "R", Action: "Review All"})
		}
		return components.RenderHintsStructured(hintList)

	case ViewOperationReview:
		return components.RenderHintsStructured([]components.Hint{
			{Key: "Esc", Action: "Back"},
			{Key: "↑↓", Action: "Scroll"},
			{Key: "PgUp/Dn", Action: "Jump"},
		})

	case ViewHistory:
		hintList := []components.Hint{
			{Key: "Esc", Action: "Back"},
			{Key: "c", Action: "Category"},
			{Key: "d", Action: "Difficulty"},
			{Key: "p", Action: "Period"},
			{Key: "↑↓", Action: "Navigate"},
			{Key: "Enter", Action: "Details"},
		}
		return components.RenderHintsStructured(hintList)

	case ViewSessionDetail:
		if m.selectedSession != nil && len(m.selectedSession.Questions) > minQuestionsForFullLog {
			return components.RenderHintsStructured([]components.Hint{
				{Key: "Esc", Action: "Back"},
				{Key: "L", Action: "Full Log"},
				{Key: "H", Action: "History"},
			})
		}
		return components.RenderHintsStructured([]components.Hint{
			{Key: "Esc", Action: "Back"},
			{Key: "H", Action: "History"},
		})

	case ViewSessionFullLog:
		return components.RenderHintsStructured([]components.Hint{
			{Key: "Esc", Action: "Back"},
			{Key: "↑↓", Action: "Scroll"},
			{Key: "←→", Action: "Filter"},
			{Key: "PgUp/Dn", Action: "Jump"},
			{Key: "S", Action: "Summary"},
		})

	case ViewTrends:
		return components.RenderHintsStructured([]components.Hint{
			{Key: "Esc", Action: "Back"},
			{Key: "m", Action: "Metric"},
			{Key: "p", Action: "Period"},
		})

	default:
		return components.RenderHintsStructured([]components.Hint{
			{Key: "Esc", Action: "Back"},
		})
	}
}

// SetSize sets the screen dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	viewportHeight := m.calculateViewportHeight()

	if !m.viewportReady {
		m.viewport = viewport.New(m.width, viewportHeight)
		m.viewport.YPosition = 0
		m.viewportReady = true
	} else {
		m.viewport.Width = m.width
		m.viewport.Height = viewportHeight
	}

	m.updateViewportContent()
}

// calculateViewportHeight returns the viewport height.
func (m Model) calculateViewportHeight() int {
	viewportHeight := m.height - statisticsHintsHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	return viewportHeight
}

// updateViewportContent updates the viewport with the current view content.
func (m *Model) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	content := m.renderCurrentViewContent()
	m.viewport.SetContent(content)
}

// renderCurrentViewContent renders the content for the current view.
func (m Model) renderCurrentViewContent() string {
	var content string

	switch m.view {
	case ViewDashboard:
		content = RenderDashboardContent(m.aggregates, m.width)

	case ViewOperations:
		content = RenderOperationsContent(
			m.operationList,
			m.operationIndex,
			m.filterPanel,
			m.width,
		)

	case ViewOperationDetail:
		extStats, ok := m.aggregates.ByOperationExtended[m.selectedOperation]
		if !ok {
			// Operation no longer exists in aggregates (e.g., after filter change)
			return RenderOperationsContent(m.operationList, m.operationIndex, m.filterPanel, m.width)
		}
		mistakes := analytics.GetRecentMistakes(m.stats, m.selectedOperation, 5)
		diffs := analytics.AllDifficulties()
		var diffFilter string
		if m.opDetailDifficultyIdx < len(diffs) {
			diffFilter = diffs[m.opDetailDifficultyIdx]
		}
		content = RenderOperationDetailContent(
			m.selectedOperation,
			extStats,
			mistakes,
			diffFilter,
			m.width,
		)

	case ViewOperationReview:
		content = RenderOperationReviewContent(
			m.selectedOperation,
			m.opReviewMistakes,
			m.width,
		)

	case ViewHistory:
		content = RenderHistoryContent(
			m.sessionList,
			m.historyNav.SelectedIndex,
			m.historyNav.CurrentPage,
			m.historyNav.SessionsPerPage,
			m.filterPanel,
			m.width,
		)

	case ViewSessionDetail:
		if m.selectedSession == nil {
			return ""
		}
		content = RenderSessionSummaryContent(*m.selectedSession, m.width)

	case ViewSessionFullLog:
		if m.selectedSession == nil {
			return ""
		}
		content = RenderSessionFullLogContent(
			*m.selectedSession,
			m.sessionLogNav.Filter,
			m.width,
		)

	case ViewTrends:
		content = RenderTrendsContent(m.trendsState, m.aggregates, m.width)

	default:
		content = RenderDashboardContent(m.aggregates, m.width)
	}

	// Center content horizontally within viewport
	if m.width > 0 {
		return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, content)
	}
	return content
}

// RenderLoading renders the loading state.
func RenderLoading(width, height int) string {
	content := "STATISTICS\n\nLoading..."
	if width > 0 && height > 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// RenderError renders the error state.
func RenderError(width, height int) string {
	content := "STATISTICS\n\nError loading statistics\n\n[Esc] Back"
	if width > 0 && height > 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}
