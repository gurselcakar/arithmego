package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// QuitConfirmSource indicates where the quit confirmation was triggered from.
type QuitConfirmSource int

const (
	QuitFromGame QuitConfirmSource = iota
	QuitFromPause
)

// QuitConfirmModel represents the quit confirmation screen.
type QuitConfirmModel struct {
	session      *game.Session
	config       *storage.Config
	source       QuitConfirmSource
	dontAskAgain bool
	selectedYes  bool // true = Yes selected, false = No selected
	focusedRow   int  // 0 = buttons row, 1 = checkbox row
	width        int
	height       int
}

// NewQuitConfirm creates a new quit confirmation model.
func NewQuitConfirm(session *game.Session, config *storage.Config, source QuitConfirmSource) QuitConfirmModel {
	return QuitConfirmModel{
		session:     session,
		config:      config,
		source:      source,
		selectedYes: false, // Default to "No" (safer option)
	}
}

// Init initializes the quit confirmation model.
func (m QuitConfirmModel) Init() tea.Cmd {
	return nil
}

// QuitConfirmCancelMsg is sent when the user cancels the quit.
type QuitConfirmCancelMsg struct {
	Session *game.Session
	Source  QuitConfirmSource
}

// QuitConfirmAcceptMsg is sent when the user confirms the quit.
type QuitConfirmAcceptMsg struct {
	DontAskAgain bool
}

// Update handles quit confirmation input.
func (m QuitConfirmModel) Update(msg tea.Msg) (QuitConfirmModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			// Quick key for Yes
			return m, func() tea.Msg {
				return QuitConfirmAcceptMsg{DontAskAgain: m.dontAskAgain}
			}
		case "n", "esc":
			// Quick key for No / Cancel
			return m, func() tea.Msg {
				return QuitConfirmCancelMsg{Session: m.session, Source: m.source}
			}
		case "up", "k":
			// Move to buttons row
			if m.focusedRow > 0 {
				m.focusedRow--
			}
		case "down", "j":
			// Move to checkbox row
			if m.focusedRow < 1 {
				m.focusedRow++
			}
		case "left", "h":
			if m.focusedRow == 0 {
				// On buttons row: select Yes
				m.selectedYes = true
			}
		case "right", "l":
			if m.focusedRow == 0 {
				// On buttons row: select No
				m.selectedYes = false
			}
		case " ", "enter":
			if m.focusedRow == 1 {
				// On checkbox row: toggle checkbox
				m.dontAskAgain = !m.dontAskAgain
			} else {
				// On buttons row: confirm selection
				if m.selectedYes {
					return m, func() tea.Msg {
						return QuitConfirmAcceptMsg{DontAskAgain: m.dontAskAgain}
					}
				}
				return m, func() tea.Msg {
					return QuitConfirmCancelMsg{Session: m.session, Source: m.source}
				}
			}
		}
	}

	return m, nil
}

// View renders the quit confirmation screen.
func (m QuitConfirmModel) View() string {
	var b strings.Builder

	// Title
	title := styles.Bold.Render("QUIT GAME?")

	// Warning message
	warning := styles.Subtle.Render("Your progress will not be saved.")

	// Yes/No buttons
	var yesBtn, noBtn string
	if m.focusedRow == 0 {
		// Buttons row is focused
		if m.selectedYes {
			yesBtn = styles.Selected.Render("[ Yes ]")
			noBtn = styles.Unselected.Render("  No  ")
		} else {
			yesBtn = styles.Unselected.Render("  Yes  ")
			noBtn = styles.Selected.Render("[ No ]")
		}
	} else {
		// Buttons row not focused - show current selection dimmed
		if m.selectedYes {
			yesBtn = styles.Subtle.Render("[ Yes ]")
			noBtn = styles.Dim.Render("  No  ")
		} else {
			yesBtn = styles.Dim.Render("  Yes  ")
			noBtn = styles.Subtle.Render("[ No ]")
		}
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yesBtn, "    ", noBtn)

	// Checkbox for "Don't ask again"
	var checkbox string
	if m.focusedRow == 1 {
		// Checkbox row is focused
		if m.dontAskAgain {
			checkbox = styles.Selected.Render("[x] Don't ask again")
		} else {
			checkbox = styles.Selected.Render("[ ] Don't ask again")
		}
	} else {
		// Checkbox row not focused
		if m.dontAskAgain {
			checkbox = styles.Accent.Render("[x] Don't ask again")
		} else {
			checkbox = styles.Subtle.Render("[ ] Don't ask again")
		}
	}

	// Hints
	hints := components.RenderHints([]string{"Y/N Quick select", "↑↓ Navigate", "Enter Select"})

	// Combine
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		warning,
		"",
		"",
		buttons,
		"",
		checkbox,
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

// Session returns the current session.
func (m QuitConfirmModel) Session() *game.Session {
	return m.session
}

// Source returns where the quit was triggered from.
func (m QuitConfirmModel) Source() QuitConfirmSource {
	return m.source
}

// SetSize sets the screen dimensions.
func (m *QuitConfirmModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
