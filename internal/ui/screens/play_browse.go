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

// PlayBrowseModel represents the Mode Browser screen (Step 1 of play flow).
type PlayBrowseModel struct {
	width  int
	height int

	viewport      viewport.Model
	viewportReady bool

	modes  []*modes.Mode // All modes
	cursor int           // Selected mode index (into modes slice)

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
		return m.updateNavigation(msg)
	}

	return m, nil
}

// updateNavigation handles navigation input.
func (m PlayBrowseModel) updateNavigation(msg tea.KeyMsg) (PlayBrowseModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return m, func() tea.Msg { return ReturnToMenuMsg{} }

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
	if len(m.modes) == 0 {
		return
	}

	newCursor := m.cursor + delta
	if newCursor < 0 {
		newCursor = len(m.modes) - 1
	}
	if newCursor >= len(m.modes) {
		newCursor = 0
	}
	m.cursor = newCursor
}

// Modes returns the list of modes.
func (m PlayBrowseModel) Modes() []*modes.Mode {
	return m.modes
}

// selectedMode returns the currently selected mode.
func (m PlayBrowseModel) selectedMode() *modes.Mode {
	if m.cursor >= 0 && m.cursor < len(m.modes) {
		return m.modes[m.cursor]
	}
	return nil
}

// scrollToSelection adjusts viewport to keep selection visible.
func (m *PlayBrowseModel) scrollToSelection() {
	if !m.viewportReady {
		return
	}

	// Calculate approximate line number of selection
	// Each category header takes 2 lines, each mode takes 1 line
	if len(m.modes) == 0 {
		return
	}

	// Find the line number of the selected mode
	lineNum := 0

	for _, catName := range categoryOrder {
		modeIDs := categoryModes[catName]
		lineNum += 2 // Category header + blank line

		for _, modeID := range modeIDs {
			for i, mode := range m.modes {
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
		lipgloss.Place(m.width, components.HintsHeight, lipgloss.Center, lipgloss.Center, hints),
	)
}

// getHints returns the context-aware hints.
func (m PlayBrowseModel) getHints() string {
	return components.RenderHintsResponsive([]components.Hint{
		{Key: "Esc", Action: "Back"},
		{Key: "↑↓", Action: "Navigate"},
		{Key: "→", Action: "Select"},
	}, m.width)
}

// SetSize sets the screen dimensions.
func (m *PlayBrowseModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	viewportHeight := m.calculateViewportHeight()

	components.SetViewportSize(&m.viewport, &m.viewportReady, m.width, viewportHeight)

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
	viewportHeight := m.height - components.HintsHeight
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
	title := lipgloss.Place(m.width, 1, lipgloss.Center, lipgloss.Center, styles.Logo.Render("PLAY"))
	lines = append(lines, title)
	lines = append(lines, "")

	lines = append(lines, m.renderCategorizedModes()...)

	// Vertically center content if it fits within viewport
	contentHeight := len(lines)
	viewportHeight := m.calculateViewportHeight()
	if contentHeight < viewportHeight {
		topPadding := (viewportHeight - contentHeight) / 2
		paddingLines := make([]string, topPadding)
		lines = append(paddingLines, lines...)
	}

	return strings.Join(lines, "\n")
}

// renderCategorizedModes renders modes grouped by category.
func (m PlayBrowseModel) renderCategorizedModes() []string {
	var lines []string

	// Find the longest mode name for alignment
	maxNameLen := 0
	for _, mode := range m.modes {
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
			for i, mode := range m.modes {
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
