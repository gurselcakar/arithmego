package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// ModeSelectMsg is sent when a mode is selected.
type ModeSelectMsg struct {
	Mode *modes.Mode
}

// ModesModel represents the modes selection screen.
type ModesModel struct {
	modes  []*modes.Mode
	cursor int
	width  int
	height int
}

// NewModes creates a new modes model.
func NewModes() ModesModel {
	return ModesModel{
		modes:  modes.All(),
		cursor: 0,
	}
}

// Init initializes the modes model.
func (m ModesModel) Init() tea.Cmd {
	return nil
}

// Update handles modes screen input.
func (m ModesModel) Update(msg tea.Msg) (ModesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.modes)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.modes) > 0 {
				return m, func() tea.Msg {
					return ModeSelectMsg{Mode: m.modes[m.cursor]}
				}
			}
		case "esc":
			return m, func() tea.Msg {
				return ReturnToMenuMsg{}
			}
		}
	}

	return m, nil
}

// View renders the modes screen.
func (m ModesModel) View() string {
	var b strings.Builder

	title := styles.Bold.Render("SELECT MODE")

	// Group modes by category
	sprintModes := modes.ByCategory(modes.CategorySprint)
	challengeModes := modes.ByCategory(modes.CategoryChallenge)

	// Build mode list
	var modeLines []string

	// Sprint section
	modeLines = append(modeLines, styles.Subtle.Render("─── Sprint ───"))
	modeLines = append(modeLines, "")
	for _, mode := range sprintModes {
		modeLines = append(modeLines, m.renderModeItem(mode))
	}

	// Spacer between categories
	modeLines = append(modeLines, "")
	modeLines = append(modeLines, styles.Subtle.Render("─── Challenge ───"))
	modeLines = append(modeLines, "")

	// Challenge section
	for _, mode := range challengeModes {
		modeLines = append(modeLines, m.renderModeItem(mode))
	}

	modeList := strings.Join(modeLines, "\n")

	// Description of selected mode
	var description string
	if len(m.modes) > 0 && m.cursor < len(m.modes) {
		selected := m.modes[m.cursor]
		description = styles.Subtle.Render(selected.Description)
	}

	// Hints
	hints := components.RenderHintsStructured([]components.Hint{
		{Key: "Esc", Action: "Back"},
		{Key: "↑↓", Action: "Navigate"},
		{Key: "Enter", Action: "Select"},
	})

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		modeList,
		"",
		description,
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

// renderModeItem renders a single mode item.
func (m ModesModel) renderModeItem(mode *modes.Mode) string {
	// Find index of this mode in the full list
	idx := -1
	for i, md := range m.modes {
		if md.ID == mode.ID {
			idx = i
			break
		}
	}

	if idx == m.cursor {
		return styles.Selected.Render("> " + mode.Name)
	}
	return styles.Unselected.Render("  " + mode.Name)
}

// SetSize updates the screen dimensions.
func (m *ModesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
