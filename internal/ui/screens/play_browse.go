package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// ModeSelectedMsg is sent when the user selects a mode in the browse screen.
type ModeSelectedMsg struct {
	Mode *modes.Mode
}

// Category ordering for display
var categoryOrder = []string{"Basics", "Powers", "Advanced", "Mixed"}

// categoryModes maps category names to mode IDs in display order
var categoryModes = map[string][]string{
	"Basics":   {modes.IDAddition, modes.IDSubtraction, modes.IDMultiplication, modes.IDDivision},
	"Powers":   {modes.IDSquares, modes.IDCubes, modes.IDSquareRoots, modes.IDCubeRoots},
	"Advanced": {modes.IDExponents, modes.IDRemainders, modes.IDPercentages, modes.IDFactorials},
	"Mixed":    {modes.IDMixedBasics, modes.IDMixedPowers, modes.IDMixedAdvanced, modes.IDAnythingGoes},
}

// Layout constants
const (
	playBrowseHintsHeight = 3
)

// PlayBrowseModel represents the Mode Browser screen (Step 1 of play flow).
type PlayBrowseModel struct {
	width  int
	height int

	viewport      viewport.Model
	viewportReady bool

	modes  []*modes.Mode // All modes
	cursor int           // Selected mode index (into modes slice)

	searchQuery   string
	searchFocused bool
	filteredModes []*modes.Mode // Modes after filter (nil = show all)

	lastPlayedModeID string
	config           *storage.Config
}

// NewPlayBrowse creates a new PlayBrowseModel.
func NewPlayBrowse(config *storage.Config) PlayBrowseModel {
	allModes := modes.All()

	m := PlayBrowseModel{
		modes:    allModes,
		cursor:   0,
		config:   config,
		viewport: viewport.New(0, 0),
	}

	// Pre-select last played mode if available
	if config != nil && config.LastPlayedModeID != "" {
		m.lastPlayedModeID = config.LastPlayedModeID
		for i, mode := range allModes {
			if mode.ID == config.LastPlayedModeID {
				m.cursor = i
				break
			}
		}
	}

	return m
}

// Init initializes the PlayBrowseModel.
func (m PlayBrowseModel) Init() tea.Cmd {
	return nil
}

// Update handles PlayBrowseModel input.
func (m PlayBrowseModel) Update(msg tea.Msg) (PlayBrowseModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		if m.searchFocused {
			return m.updateSearchInput(msg)
		}
		return m.updateNavigation(msg)
	}

	return m, nil
}

// updateSearchInput handles input when search is focused.
func (m PlayBrowseModel) updateSearchInput(msg tea.KeyMsg) (PlayBrowseModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Clear search and exit search mode
		m.searchFocused = false
		m.searchQuery = ""
		m.filteredModes = nil
		m.updateViewportContent()
		return m, nil

	case "enter":
		// Exit search mode, keep filter
		m.searchFocused = false
		m.updateViewportContent()
		return m, nil

	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.applyFilter()
			m.updateViewportContent()
		}
		return m, nil

	case "up", "down":
		// Allow navigation while in search mode
		m.searchFocused = false
		return m.updateNavigation(msg)

	default:
		// Add character to search
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
			m.applyFilter()
			m.updateViewportContent()
		}
		return m, nil
	}
}

// updateNavigation handles navigation input.
func (m PlayBrowseModel) updateNavigation(msg tea.KeyMsg) (PlayBrowseModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// If search query exists, clear it first
		if m.searchQuery != "" {
			m.searchQuery = ""
			m.filteredModes = nil
			m.updateViewportContent()
			return m, nil
		}
		return m, func() tea.Msg { return ReturnToMenuMsg{} }

	case "/":
		m.searchFocused = true
		m.updateViewportContent()
		return m, nil

	case "up", "k":
		m.moveCursor(-1)
		m.scrollToSelection()
		m.updateViewportContent()
		return m, nil

	case "down", "j":
		m.moveCursor(1)
		m.scrollToSelection()
		m.updateViewportContent()
		return m, nil

	case "enter", "right", "l":
		mode := m.selectedMode()
		if mode != nil {
			return m, func() tea.Msg { return ModeSelectedMsg{Mode: mode} }
		}
		return m, nil
	}

	return m, nil
}

// moveCursor moves the cursor by delta, skipping category headers.
func (m *PlayBrowseModel) moveCursor(delta int) {
	modeList := m.getDisplayModes()
	if len(modeList) == 0 {
		return
	}

	newCursor := m.cursor + delta
	if newCursor < 0 {
		newCursor = 0
	}
	if newCursor >= len(modeList) {
		newCursor = len(modeList) - 1
	}
	m.cursor = newCursor
}

// selectedMode returns the currently selected mode.
func (m PlayBrowseModel) selectedMode() *modes.Mode {
	modeList := m.getDisplayModes()
	if m.cursor >= 0 && m.cursor < len(modeList) {
		return modeList[m.cursor]
	}
	return nil
}

// getDisplayModes returns the modes to display (filtered or all).
func (m PlayBrowseModel) getDisplayModes() []*modes.Mode {
	if m.filteredModes != nil {
		return m.filteredModes
	}
	return m.modes
}

// applyFilter filters modes based on search query.
func (m *PlayBrowseModel) applyFilter() {
	if m.searchQuery == "" {
		m.filteredModes = nil
		return
	}

	query := strings.ToLower(m.searchQuery)
	var filtered []*modes.Mode

	for _, mode := range m.modes {
		nameLower := strings.ToLower(mode.Name)
		descLower := strings.ToLower(mode.Description)
		if strings.Contains(nameLower, query) || strings.Contains(descLower, query) {
			filtered = append(filtered, mode)
		}
	}

	m.filteredModes = filtered

	// Reset cursor if out of bounds
	if m.cursor >= len(filtered) {
		m.cursor = 0
	}
}

// scrollToSelection adjusts viewport to keep selection visible.
func (m *PlayBrowseModel) scrollToSelection() {
	if !m.viewportReady {
		return
	}

	// Calculate approximate line number of selection
	// Each category header takes 2 lines, each mode takes 1 line
	modeList := m.getDisplayModes()
	if len(modeList) == 0 {
		return
	}

	// Find the line number of the selected mode
	lineNum := 3 // Title + blank line + search

	if m.filteredModes == nil {
		// With categories
		for _, catName := range categoryOrder {
			modeIDs := categoryModes[catName]
			lineNum += 2 // Category header + blank line

			for _, modeID := range modeIDs {
				for i, mode := range modeList {
					if mode.ID == modeID && i == m.cursor {
						goto found
					}
					if mode.ID == modeID {
						lineNum++
					}
				}
			}
			lineNum++ // Blank line after category
		}
	} else {
		// Without categories (filtered)
		lineNum += 2 // Search results header
		lineNum += m.cursor
	}

found:
	// Scroll viewport to show selection
	if lineNum < m.viewport.YOffset+2 {
		m.viewport.SetYOffset(lineNum - 2)
	} else if lineNum > m.viewport.YOffset+m.viewport.Height-3 {
		m.viewport.SetYOffset(lineNum - m.viewport.Height + 3)
	}
}

// View renders the PlayBrowseModel.
func (m PlayBrowseModel) View() string {
	if !m.viewportReady {
		return "Loading..."
	}

	hints := m.getHints()

	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		lipgloss.Place(m.width, playBrowseHintsHeight, lipgloss.Center, lipgloss.Center, hints),
	)
}

// getHints returns the context-aware hints.
func (m PlayBrowseModel) getHints() string {
	if m.searchFocused {
		return components.RenderHintsStructured([]components.Hint{
			{Key: "Esc", Action: "Cancel"},
			{Key: "Enter", Action: "Apply"},
			{Key: "↑↓", Action: "Navigate"},
		})
	}
	return components.RenderHintsStructured([]components.Hint{
		{Key: "Esc", Action: "Back"},
		{Key: "↑↓", Action: "Navigate"},
		{Key: "/", Action: "Search"},
		{Key: "→", Action: "Select"},
	})
}

// SetSize sets the screen dimensions.
func (m *PlayBrowseModel) SetSize(width, height int) {
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

	// Reset scroll position if content fits within viewport
	totalLines := m.viewport.TotalLineCount()
	if totalLines <= viewportHeight {
		m.viewport.SetYOffset(0)
	} else if m.viewport.YOffset > totalLines-viewportHeight {
		// Clamp YOffset if it's now past the end of content
		newOffset := totalLines - viewportHeight
		if newOffset < 0 {
			newOffset = 0
		}
		m.viewport.SetYOffset(newOffset)
	}

	m.scrollToSelection()
}

// calculateViewportHeight returns the viewport height.
func (m PlayBrowseModel) calculateViewportHeight() int {
	viewportHeight := m.height - playBrowseHintsHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	return viewportHeight
}

// updateViewportContent updates the viewport with the current content.
func (m *PlayBrowseModel) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	content := m.renderContent()
	m.viewport.SetContent(content)
}

// renderContent renders the main content for the viewport.
func (m PlayBrowseModel) renderContent() string {
	var lines []string

	// Title
	title := styles.Logo.Render("PLAY")
	lines = append(lines, lipgloss.Place(m.width, 1, lipgloss.Center, lipgloss.Center, title))
	lines = append(lines, "")

	// Search bar
	searchBar := m.renderSearchBar()
	lines = append(lines, lipgloss.Place(m.width, 1, lipgloss.Center, lipgloss.Center, searchBar))
	lines = append(lines, "")

	modeList := m.getDisplayModes()

	if len(modeList) == 0 {
		// No results
		noResults := styles.Dim.Render("No modes match your search")
		lines = append(lines, lipgloss.Place(m.width, 1, lipgloss.Center, lipgloss.Center, noResults))
	} else if m.filteredModes != nil {
		// Filtered results (no categories)
		lines = append(lines, m.renderFilteredModes()...)
	} else {
		// All modes with categories
		lines = append(lines, m.renderCategorizedModes()...)
	}

	return strings.Join(lines, "\n")
}

// renderSearchBar renders the search input.
func (m PlayBrowseModel) renderSearchBar() string {
	prefix := styles.Dim.Render("/")
	if m.searchFocused {
		prefix = styles.Accent.Render("/")
	}

	if m.searchQuery == "" {
		if m.searchFocused {
			return prefix + styles.Dim.Render(" Search modes...") + styles.Accent.Render("_")
		}
		return prefix + styles.Dim.Render(" Search modes...")
	}

	if m.searchFocused {
		return prefix + " " + m.searchQuery + styles.Accent.Render("_")
	}
	return prefix + " " + m.searchQuery
}

// renderCategorizedModes renders modes grouped by category.
func (m PlayBrowseModel) renderCategorizedModes() []string {
	var lines []string
	modeList := m.getDisplayModes()

	// Find the longest mode name for alignment
	maxNameLen := 0
	for _, mode := range modeList {
		if len(mode.Name) > maxNameLen {
			maxNameLen = len(mode.Name)
		}
	}

	// Calculate content width for centering
	contentWidth := maxNameLen + 30 // name + description space
	leftPadding := (m.width - contentWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}
	padding := strings.Repeat(" ", leftPadding)

	for _, catName := range categoryOrder {
		modeIDs := categoryModes[catName]

		// Category header
		header := styles.Dim.Render("── " + catName + " " + strings.Repeat("─", 40))
		lines = append(lines, padding+header)
		lines = append(lines, "")

		// Modes in this category
		for _, modeID := range modeIDs {
			for i, mode := range modeList {
				if mode.ID == modeID {
					line := m.renderModeLine(mode, i == m.cursor, maxNameLen)
					lines = append(lines, padding+line)
					break
				}
			}
		}
		lines = append(lines, "")
	}

	return lines
}

// renderFilteredModes renders filtered modes without categories.
func (m PlayBrowseModel) renderFilteredModes() []string {
	var lines []string
	modeList := m.filteredModes

	// Find the longest mode name for alignment
	maxNameLen := 0
	for _, mode := range modeList {
		if len(mode.Name) > maxNameLen {
			maxNameLen = len(mode.Name)
		}
	}

	// Calculate content width for centering
	contentWidth := maxNameLen + 30
	leftPadding := (m.width - contentWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}
	padding := strings.Repeat(" ", leftPadding)

	for i, mode := range modeList {
		line := m.renderModeLine(mode, i == m.cursor, maxNameLen)
		lines = append(lines, padding+line)
	}

	return lines
}

// renderModeLine renders a single mode line.
func (m PlayBrowseModel) renderModeLine(mode *modes.Mode, selected bool, maxNameLen int) string {
	// Focus indicator
	var prefix string
	if selected {
		prefix = styles.Accent.Render("> ")
	} else {
		prefix = "  "
	}

	// Mode name with padding
	name := mode.Name
	namePadded := name + strings.Repeat(" ", maxNameLen-len(name))

	// Description
	desc := mode.Description

	if selected {
		return prefix + styles.Bold.Render(namePadded) + "    " + styles.Subtle.Render(desc)
	}
	return prefix + styles.Subtle.Render(namePadded) + "    " + styles.Dim.Render(desc)
}
