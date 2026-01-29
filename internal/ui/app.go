package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/ui/screens"
)

// App is the main Bubble Tea model that orchestrates all screens.
type App struct {
	screen Screen
	width  int
	height int

	// Screen models
	menuModel       screens.MenuModel
	gameModel       screens.GameModel
	pauseModel      screens.PauseModel
	resultsModel    screens.ResultsModel
	modesModel      screens.ModesModel
	launchModel     screens.LaunchModel
	practiceModel   screens.PracticeModel
	statisticsModel screens.StatisticsModel
	settingsModel   screens.SettingsModel
	onboardingModel screens.OnboardingModel

	// Current session state
	session        *game.Session
	currentMode    *modes.Mode
	lastDifficulty game.Difficulty
	lastDuration   time.Duration
}

// New creates a new App instance.
func New() *App {
	return &App{
		screen:          ScreenMenu,
		menuModel:       screens.NewMenu(),
		modesModel:      screens.NewModes(),
		practiceModel:   screens.NewPractice(),
		statisticsModel: screens.NewStatistics(),
		settingsModel:   screens.NewSettings(),
		onboardingModel: screens.NewOnboarding(),
	}
}

// Init initializes the app.
func (a *App) Init() tea.Cmd {
	return nil
}

// Update handles all messages and routes them to the appropriate screen.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size for all screens
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		a.width = wsm.Width
		a.height = wsm.Height
	}

	// Route based on current screen
	switch a.screen {
	case ScreenMenu:
		return a.updateMenu(msg)
	case ScreenGame:
		return a.updateGame(msg)
	case ScreenPause:
		return a.updatePause(msg)
	case ScreenResults:
		return a.updateResults(msg)
	case ScreenModes:
		return a.updateModes(msg)
	case ScreenLaunch:
		return a.updateLaunch(msg)
	case ScreenPractice:
		return a.updatePractice(msg)
	case ScreenStatistics:
		return a.updateStatistics(msg)
	case ScreenSettings:
		return a.updateSettings(msg)
	case ScreenOnboarding:
		return a.updateOnboarding(msg)
	}

	return a, nil
}

// updateMenu handles menu screen updates.
func (a *App) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.menuModel, cmd = a.menuModel.Update(msg)

	// Check for menu selection
	if selectMsg, ok := msg.(screens.MenuSelectMsg); ok {
		switch selectMsg.Action {
		case screens.ActionModes:
			a.modesModel = screens.NewModes()
			a.modesModel.SetSize(a.width, a.height)
			a.screen = ScreenModes
			return a, a.modesModel.Init()
		case screens.ActionPractice:
			a.screen = ScreenPractice
			return a, a.practiceModel.Init()
		case screens.ActionStatistics:
			a.screen = ScreenStatistics
			return a, a.statisticsModel.Init()
		case screens.ActionSettings:
			a.screen = ScreenSettings
			return a, a.settingsModel.Init()
		}
	}

	// Check if quitting
	if a.menuModel.Quitting() {
		return a, tea.Quit
	}

	return a, cmd
}

// updateGame handles game screen updates.
func (a *App) updateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.gameModel, cmd = a.gameModel.Update(msg)

	// Check for game over
	if gom, ok := msg.(screens.GameOverMsg); ok {
		a.session = gom.Session
		a.resultsModel = screens.NewResults(a.session)
		a.resultsModel.SetSize(a.width, a.height)
		a.screen = ScreenResults
		return a, a.resultsModel.Init()
	}

	// Check for pause
	if pm, ok := msg.(screens.PauseMsg); ok {
		a.session = pm.Session
		a.pauseModel = screens.NewPause(a.session)
		a.pauseModel.SetSize(a.width, a.height)
		a.screen = ScreenPause
		return a, a.pauseModel.Init()
	}

	return a, cmd
}

// updatePause handles pause screen updates.
func (a *App) updatePause(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.pauseModel, cmd = a.pauseModel.Update(msg)

	// Check for resume
	if rm, ok := msg.(screens.ResumeMsg); ok {
		a.session = rm.Session
		// Resume the session - need to restart the timer tick
		a.gameModel.SetSession(a.session)
		a.screen = ScreenGame
		// Restart the session timer from where it was
		a.session.Resume()
		return a, screens.TickCmd()
	}

	// Check for quit to menu
	if _, ok := msg.(screens.QuitToMenuMsg); ok {
		a.screen = ScreenMenu
		a.session = nil
		return a, nil
	}

	return a, cmd
}

// updateResults handles results screen updates.
func (a *App) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.resultsModel, cmd = a.resultsModel.Update(msg)

	// Check for play again
	if _, ok := msg.(screens.PlayAgainMsg); ok {
		return a.startGame()
	}

	// Check for return to menu
	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		a.screen = ScreenMenu
		a.session = nil
		return a, nil
	}

	return a, cmd
}

// updateModes handles modes screen updates.
func (a *App) updateModes(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.modesModel, cmd = a.modesModel.Update(msg)

	// Check for mode selection
	if selectMsg, ok := msg.(screens.ModeSelectMsg); ok {
		a.currentMode = selectMsg.Mode
		a.launchModel = screens.NewLaunch(a.currentMode)
		a.launchModel.SetSize(a.width, a.height)
		a.screen = ScreenLaunch
		return a, a.launchModel.Init()
	}

	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		a.screen = ScreenMenu
		return a, nil
	}

	return a, cmd
}

// updateLaunch handles launch screen updates.
func (a *App) updateLaunch(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.launchModel, cmd = a.launchModel.Update(msg)

	// Check for start game
	if startMsg, ok := msg.(screens.StartGameMsg); ok {
		a.currentMode = startMsg.Mode
		a.lastDifficulty = startMsg.Difficulty
		a.lastDuration = startMsg.Duration
		return a.startGame()
	}

	// Check for return to modes
	if _, ok := msg.(screens.ReturnToModesMsg); ok {
		a.modesModel = screens.NewModes()
		a.modesModel.SetSize(a.width, a.height)
		a.screen = ScreenModes
		return a, a.modesModel.Init()
	}

	return a, cmd
}

// updatePractice handles practice screen updates.
func (a *App) updatePractice(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.practiceModel, cmd = a.practiceModel.Update(msg)

	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		a.screen = ScreenMenu
		return a, nil
	}

	return a, cmd
}

// updateStatistics handles statistics screen updates.
func (a *App) updateStatistics(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.statisticsModel, cmd = a.statisticsModel.Update(msg)

	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		a.screen = ScreenMenu
		return a, nil
	}

	return a, cmd
}

// updateSettings handles settings screen updates.
func (a *App) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.settingsModel, cmd = a.settingsModel.Update(msg)

	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		a.screen = ScreenMenu
		return a, nil
	}

	return a, cmd
}

// updateOnboarding handles onboarding screen updates.
func (a *App) updateOnboarding(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.onboardingModel, cmd = a.onboardingModel.Update(msg)

	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		a.screen = ScreenMenu
		return a, nil
	}

	return a, cmd
}

// startGame creates a new session and starts the game.
func (a *App) startGame() (tea.Model, tea.Cmd) {
	if a.currentMode == nil || len(a.currentMode.Operations) == 0 {
		panic("startGame: invalid mode state")
	}

	a.session = game.NewSession(a.currentMode.Operations, a.lastDifficulty, a.lastDuration)
	a.gameModel = screens.NewGame(a.session)
	a.gameModel.SetSize(a.width, a.height)
	a.screen = ScreenGame
	return a, a.gameModel.Init()
}

// View renders the current screen.
func (a *App) View() string {
	switch a.screen {
	case ScreenMenu:
		return a.menuModel.View()
	case ScreenGame:
		return a.gameModel.View()
	case ScreenPause:
		return a.pauseModel.View()
	case ScreenResults:
		return a.resultsModel.View()
	case ScreenModes:
		return a.modesModel.View()
	case ScreenLaunch:
		return a.launchModel.View()
	case ScreenPractice:
		return a.practiceModel.View()
	case ScreenStatistics:
		return a.statisticsModel.View()
	case ScreenSettings:
		return a.settingsModel.View()
	case ScreenOnboarding:
		return a.onboardingModel.View()
	default:
		return ""
	}
}

// Phase 9: Replace main.go with Cobra CLI
// - arithmego (no args) → TUI menu
// - arithmego play → Quick play
// - arithmego statistics → Statistics screen
// - arithmego version → Version info
