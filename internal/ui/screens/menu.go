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
	ActionModes MenuAction = iota
	ActionPractice
	ActionStatistics
	ActionSettings
)

// Phase 6: Add Quick Play as first menu item (conditional on hasPlayedBefore)
// Quick Play shows current mode: "> Quick Play · Addition Sprint"
// Right arrow on Quick Play enters mode switch: "[◀ Addition Sprint ▶]"

// MenuModel represents the main menu screen.
type MenuModel struct {
	items    []MenuItem
	cursor   int
	width    int
	height   int
	quitting bool
}

// NewMenu creates a new menu model.
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
	hints := components.RenderHints([]string{"↑↓ Navigate", "Enter Select"})

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		tagline,
		"",
		"",
		menu,
		"",
		"",
		hints,
	)

	// Center in terminal
	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	b.WriteString(content)
	return b.String()
}

// Quitting returns true if the user is quitting.
func (m MenuModel) Quitting() bool {
	return m.quitting
}
