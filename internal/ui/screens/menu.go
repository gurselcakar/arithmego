package screens

import (
	"strings"

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
	ActionQuickPlay MenuAction = iota
	ActionModes
	ActionPractice
	ActionStatistics
	ActionSettings
)

// QuickPlayInfo contains information for displaying the Quick Play menu item.
type QuickPlayInfo struct {
	ModeName string
}

// MenuModel represents the main menu screen.
type MenuModel struct {
	items         []MenuItem
	cursor        int
	width         int
	height        int
	quitting      bool
	updateVersion string // Available update version (empty if none)
}

// NewMenu creates a new menu model without Quick Play.
func NewMenu() MenuModel {
	return MenuModel{
		items: []MenuItem{
			{Label: "Modes", Action: ActionModes},
			{Label: "Practice", Action: ActionPractice},
			{IsSpacer: true},
			{Label: "Statistics", Action: ActionStatistics},
			{Label: "Settings", Action: ActionSettings},
		},
		cursor: 0,
	}
}

// NewMenuWithQuickPlay creates a new menu model with Quick Play as the first item.
func NewMenuWithQuickPlay(quickPlay *QuickPlayInfo) MenuModel {
	label := "Quick Play"
	if quickPlay != nil && quickPlay.ModeName != "" {
		label = "Quick Play · " + quickPlay.ModeName
	}
	return MenuModel{
		items: []MenuItem{
			{Label: label, Action: ActionQuickPlay},
			{Label: "Modes", Action: ActionModes},
			{Label: "Practice", Action: ActionPractice},
			{IsSpacer: true},
			{Label: "Statistics", Action: ActionStatistics},
			{Label: "Settings", Action: ActionSettings},
		},
		cursor: 0,
	}
}

// Init initializes the menu model.
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles menu input.
func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case "enter":
			return m, m.selectItem()
		case "esc", "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// moveCursor moves the cursor, skipping spacers.
func (m *MenuModel) moveCursor(delta int) {
	for {
		m.cursor += delta
		if m.cursor < 0 {
			m.cursor = len(m.items) - 1
		}
		if m.cursor >= len(m.items) {
			m.cursor = 0
		}
		if !m.items[m.cursor].IsSpacer {
			break
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

	var b strings.Builder

	// Logo
	logo := components.Logo()
	separator := styles.Dim.Render(components.LogoSeparator())
	tagline := components.Tagline()

	// Menu items
	var menuItems []string
	for i, item := range m.items {
		if item.IsSpacer {
			menuItems = append(menuItems, "")
			continue
		}

		line := "  " + item.Label
		if i == m.cursor {
			line = styles.Selected.Render("> " + item.Label)
		} else {
			line = styles.Unselected.Render("  " + item.Label)
		}
		menuItems = append(menuItems, line)
	}
	menu := strings.Join(menuItems, "\n")

	// Hints
	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "↑↓", Action: "Navigate"},
		{Key: "Enter", Action: "Select"},
	})

	// Update notification (if available)
	var updateNotice string
	if m.updateVersion != "" {
		// TODO: Implement actual auto-update in Phase 12 (Distribution).
		// For now, we just notify the user to run the update command manually.
		updateNotice = styles.Dim.Render("Update available: " + m.updateVersion + " · run 'arithmego update'")
	}

	// Build main content (without hints)
	mainContent := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		separator,
		"",
		tagline,
		"",
		"",
		menu,
	)

	// Build bottom section with hints and optional update notice
	var bottomSection string
	if updateNotice != "" {
		bottomSection = lipgloss.JoinVertical(lipgloss.Center, hints, "", updateNotice)
	} else {
		bottomSection = hints
	}

	// Bottom-anchored hints layout with small gap at bottom
	if m.width > 0 && m.height > 0 {
		bottomHeight := lipgloss.Height(bottomSection)
		bottomPadding := 1
		availableHeight := m.height - bottomHeight - bottomPadding

		centeredMain := lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, mainContent)
		centeredBottom := lipgloss.Place(m.width, bottomHeight+bottomPadding, lipgloss.Center, lipgloss.Top, bottomSection)

		b.WriteString(lipgloss.JoinVertical(lipgloss.Left, centeredMain, centeredBottom))
		return b.String()
	}

	// Fallback for unknown dimensions
	b.WriteString(lipgloss.JoinVertical(lipgloss.Center, mainContent, "", "", bottomSection))
	return b.String()
}

// Quitting returns true if the user is quitting.
func (m MenuModel) Quitting() bool {
	return m.quitting
}

// SetSize updates the screen dimensions.
func (m *MenuModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetUpdateInfo sets the available update version for display.
func (m *MenuModel) SetUpdateInfo(version string) {
	m.updateVersion = version
}
