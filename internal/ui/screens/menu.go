package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// MenuItem represents a menu option.
type MenuItem struct {
	Label    string
	Action   MenuAction
	IsSpacer bool
}

// MenuAction identifies what happens when a menu item is selected.
type MenuAction int

const (
	ActionPlay MenuAction = iota
	ActionPractice
	ActionStatistics
	ActionSettings
	ActionX
)


// MenuModel represents the main menu screen.
type MenuModel struct {
	items            []MenuItem
	cursor           int
	width            int
	height           int
	quitting         bool
	updateVersion    string // Available update version (empty if none)
	updateInstalled  string // Auto-updated version awaiting restart (empty if none)
	viewport         viewport.Model
	viewportReady    bool
}

// NewMenu creates a new menu model.
func NewMenu() MenuModel {
	return MenuModel{
		items: []MenuItem{
			{Label: "Play", Action: ActionPlay},
			{Label: "Practice", Action: ActionPractice},
			{Label: "Statistics", Action: ActionStatistics},
			{Label: "Settings", Action: ActionSettings},
			{IsSpacer: true},
			{Label: "Follow me on X @gurselcakar", Action: ActionX},
		},
		cursor:        0,
		viewport:      viewport.New(0, 0),
		viewportReady: false,
	}
}

// Init initializes the menu model.
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles menu input.
func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-1)
			m.updateViewportContent()
		case "down", "j":
			m.moveCursor(1)
			m.updateViewportContent()
		case "enter", "right", "l":
			return m, m.selectItem()
		case "esc", "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Update viewport (for mouse scrolling if enabled)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// moveCursor moves the cursor, skipping spacers.
func (m *MenuModel) moveCursor(delta int) {
	for range m.items {
		m.cursor += delta
		if m.cursor < 0 {
			m.cursor = len(m.items) - 1
		}
		if m.cursor >= len(m.items) {
			m.cursor = 0
		}
		if !m.items[m.cursor].IsSpacer {
			return
		}
	}
}

// selectItem returns a command for the selected menu action.
func (m MenuModel) selectItem() tea.Cmd {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return nil
	}
	item := m.items[m.cursor]
	if item.IsSpacer {
		return nil
	}
	return func() tea.Msg {
		return MenuSelectMsg{Action: item.Action}
	}
}

// MenuSelectMsg is sent when a menu item is selected.
type MenuSelectMsg struct {
	Action MenuAction
}

// View renders the menu screen.
func (m MenuModel) View() string {
	if m.quitting {
		return ""
	}

	if !m.viewportReady {
		return "Loading..."
	}

	// Hints
	hints := m.getHints()

	// All screens: viewport + hints
	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		lipgloss.Place(m.width, components.HintsHeight, lipgloss.Center, lipgloss.Center, hints),
	)
}

// getHints returns the hints for the menu.
func (m MenuModel) getHints() string {
	hints := components.RenderHintsResponsive([]components.Hint{
		{Key: "↑↓", Action: "Navigate"},
		{Key: "→", Action: "Select"},
	}, m.width)

	// Auto-update installed notification (takes priority)
	if m.updateInstalled != "" {
		updateNotice := styles.Correct.Render("Updated to " + m.updateInstalled + " — restart to apply")
		return lipgloss.JoinVertical(lipgloss.Center, hints, "", updateNotice)
	}

	// Manual update notification (fallback)
	if m.updateVersion != "" {
		updateNotice := styles.Dim.Render("Update available: " + m.updateVersion + " · run 'arithmego update'")
		return lipgloss.JoinVertical(lipgloss.Center, hints, "", updateNotice)
	}

	return hints
}

// Quitting returns true if the user is quitting.
func (m MenuModel) Quitting() bool {
	return m.quitting
}

// SetSize updates the screen dimensions.
func (m *MenuModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	viewportHeight := m.calculateViewportHeight()

	components.SetViewportSize(&m.viewport, &m.viewportReady, m.width, viewportHeight)

	m.updateViewportContent()
}

// calculateViewportHeight returns the viewport height.
func (m MenuModel) calculateViewportHeight() int {
	viewportHeight := m.height - components.HintsHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	return viewportHeight
}

// updateViewportContent updates the viewport with the menu content.
func (m *MenuModel) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	content := m.renderMenuContent()
	m.viewport.SetContent(content)
}

// renderMenuContent renders the main menu content for the viewport.
func (m MenuModel) renderMenuContent() string {
	// Logo with color
	logo := components.LogoColoredForWidth(m.width)
	separator := styles.Dim.Render(components.LogoSeparator())
	tagline := components.Tagline()

	// Split menu items into main (game) items and secondary items
	var mainItems []string
	var secondaryItems []string
	isSecondary := false
	for i, item := range m.items {
		if item.IsSpacer {
			isSecondary = true
			continue
		}

		var line string
		if i == m.cursor {
			line = styles.Accent.Render("> ") + styles.Selected.Render(item.Label)
		} else {
			line = "  " + styles.Unselected.Render(item.Label)
		}
		if isSecondary {
			secondaryItems = append(secondaryItems, line)
		} else {
			mainItems = append(mainItems, line)
		}
	}

	// Render main game items as a centered block
	mainMaxWidth := 0
	for _, line := range mainItems {
		if w := lipgloss.Width(line); w > mainMaxWidth {
			mainMaxWidth = w
		}
	}
	mainMenu := lipgloss.NewStyle().Width(mainMaxWidth).Render(strings.Join(mainItems, "\n"))

	// Render secondary items as a separate centered block
	secMaxWidth := 0
	for _, line := range secondaryItems {
		if w := lipgloss.Width(line); w > secMaxWidth {
			secMaxWidth = w
		}
	}
	secondaryMenu := lipgloss.NewStyle().Width(secMaxWidth).Render(strings.Join(secondaryItems, "\n"))

	// Build main content
	content := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		separator,
		"",
		tagline,
		"",
		"",
		mainMenu,
		"",
		secondaryMenu,
	)

	// Center both horizontally and vertically within viewport
	if m.width > 0 && m.viewportReady {
		return lipgloss.Place(m.width, m.viewport.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// SetUpdateInfo sets the available update version for display.
func (m *MenuModel) SetUpdateInfo(version string) {
	m.updateVersion = version
}

// SetUpdateInstalled sets the auto-updated version for display.
func (m *MenuModel) SetUpdateInstalled(version string) {
	m.updateInstalled = version
}
