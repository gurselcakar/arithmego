package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/gurselcakar/arithmego/internal/ui/screens/statistics"
)

// StatisticsModel wraps the statistics screen implementation.
type StatisticsModel struct {
	model statistics.Model
}

// NewStatistics creates a new statistics model.
func NewStatistics() StatisticsModel {
	return StatisticsModel{
		model: statistics.New(),
	}
}

// Init initializes the statistics model and loads data.
func (m StatisticsModel) Init() tea.Cmd {
	return m.model.Init()
}

// Update handles statistics screen input.
func (m StatisticsModel) Update(msg tea.Msg) (StatisticsModel, tea.Cmd) {
	// Handle the ReturnToMenuMsg from the statistics package
	if _, ok := msg.(statistics.ReturnToMenuMsg); ok {
		return m, func() tea.Msg {
			return ReturnToMenuMsg{}
		}
	}

	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

// View renders the statistics screen.
func (m StatisticsModel) View() string {
	return m.model.View()
}

// SetSize sets the screen dimensions.
func (m *StatisticsModel) SetSize(width, height int) {
	m.model.SetSize(width, height)
}
