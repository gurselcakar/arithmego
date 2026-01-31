package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/storage"
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

	// User config (for Quick Play and defaults)
	config *storage.Config

	// Error tracking
	lastSaveError error
}

// New creates a new App instance.
func New() *App {
	// Load config (ignore errors, use default config as fallback)
	config, _ := storage.LoadConfig()
	if config == nil {
		config = storage.NewConfig()
	}

	// Build menu with Quick Play if we have last played data
	var menuModel screens.MenuModel
	if config.HasLastPlayed() {
		mode, _ := modes.Get(config.LastPlayedModeID)
		menuModel = screens.NewMenuWithQuickPlay(&screens.QuickPlayInfo{
			ModeName: mode.Name,
		})
	} else {
		menuModel = screens.NewMenu()
	}

	// Determine starting screen based on onboarding status
	startScreen := ScreenMenu
	if !config.Onboarded {
		startScreen = ScreenOnboarding
	}

	return &App{
		screen:          startScreen,
		menuModel:       menuModel,
		modesModel:      screens.NewModes(),
		practiceModel:   screens.NewPractice(),
		statisticsModel: screens.NewStatistics(),
		settingsModel:   screens.NewSettings(config),
		onboardingModel: screens.NewOnboarding(),
		config:          config,
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
		case screens.ActionQuickPlay:
			return a.startQuickPlay()
		case screens.ActionModes:
			a.modesModel = screens.NewModes()
			a.modesModel.SetSize(a.width, a.height)
			a.screen = ScreenModes
			return a, a.modesModel.Init()
		case screens.ActionPractice:
			a.practiceModel = screens.NewPractice()
			a.practiceModel.SetSize(a.width, a.height)
			a.screen = ScreenPractice
			return a, a.practiceModel.Init()
		case screens.ActionStatistics:
			a.statisticsModel = screens.NewStatistics()
			a.statisticsModel.SetSize(a.width, a.height)
			a.screen = ScreenStatistics
			return a, a.statisticsModel.Init()
		case screens.ActionSettings:
			a.settingsModel = screens.NewSettings(a.config)
			a.settingsModel.SetSize(a.width, a.height)
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

	// Check for game over.
	// When this message arrives, the game model has already stopped its timer
	// and animation loops by not returning TickCmd/ScoreAnimCmd. Any stale
	// tick messages will be ignored since we transition to ScreenResults.
	if gom, ok := msg.(screens.GameOverMsg); ok {
		a.session = gom.Session
		a.saveSession()
		a.resultsModel = screens.NewResults(a.session, a.lastSaveError)
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
		a.rebuildMenu()
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
		a.launchModel = screens.NewLaunch(a.currentMode, a.config)
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

	// Keep config in sync with settings model
	a.config = a.settingsModel.Config()

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

	// Handle onboarding completion with user selections
	if completeMsg, ok := msg.(screens.OnboardingCompleteMsg); ok {
		return a.completeOnboarding(completeMsg.ModeID, completeMsg.Difficulty, completeMsg.DurationMs)
	}

	// Handle onboarding skip - use defaults (Easy, 60s, Addition)
	if _, ok := msg.(screens.OnboardingSkipMsg); ok {
		return a.completeOnboarding(modes.IDAdditionSprint, "Easy", 60000)
	}

	return a, cmd
}

// completeOnboarding finishes onboarding and starts the game with selected settings.
func (a *App) completeOnboarding(modeID, difficulty string, durationMs int64) (tea.Model, tea.Cmd) {
	// Mark as onboarded
	a.config.Onboarded = true

	// Set last played and default settings
	a.config.LastPlayedModeID = modeID
	a.config.LastPlayedDifficulty = difficulty
	a.config.LastPlayedDurationMs = durationMs
	a.config.DefaultDifficulty = difficulty
	a.config.DefaultDurationMs = durationMs

	// Save config (ignore errors - config is non-critical)
	_ = storage.SaveConfig(a.config)

	// Set up game state
	mode, _ := modes.Get(modeID)
	a.currentMode = mode
	a.lastDifficulty = game.ParseDifficulty(difficulty)
	a.lastDuration = time.Duration(durationMs) * time.Millisecond

	// Start the game
	return a.startGame()
}

// startGame creates a new session and starts the game.
// If mode is invalid, gracefully returns to the modes screen.
func (a *App) startGame() (tea.Model, tea.Cmd) {
	if a.currentMode == nil || len(a.currentMode.Operations) == 0 {
		// Gracefully recover: return to modes screen instead of crashing
		a.modesModel = screens.NewModes()
		a.modesModel.SetSize(a.width, a.height)
		a.screen = ScreenModes
		return a, a.modesModel.Init()
	}

	a.session = game.NewSession(a.currentMode.Operations, a.lastDifficulty, a.lastDuration)
	a.gameModel = screens.NewGame(a.session)
	a.gameModel.SetSize(a.width, a.height)
	a.screen = ScreenGame
	return a, a.gameModel.Init()
}

// startQuickPlay starts a game with the last played settings.
func (a *App) startQuickPlay() (tea.Model, tea.Cmd) {
	mode, _ := modes.Get(a.config.LastPlayedModeID)

	a.currentMode = mode
	a.lastDifficulty = game.ParseDifficulty(a.config.LastPlayedDifficulty)
	a.lastDuration = time.Duration(a.config.LastPlayedDurationMs) * time.Millisecond

	return a.startGame()
}

// saveLastPlayed saves the current game configuration to config.
// Must only be called from the main Bubble Tea update loop (single-threaded).
func (a *App) saveLastPlayed() {
	if a.config == nil {
		a.config = storage.NewConfig()
	}
	if a.currentMode == nil {
		return
	}

	a.config.LastPlayedModeID = a.currentMode.ID
	a.config.LastPlayedDifficulty = a.lastDifficulty.String()
	a.config.LastPlayedDurationMs = a.lastDuration.Milliseconds()

	// Ignore save errors - config is non-critical
	_ = storage.SaveConfig(a.config)
}

// rebuildMenu rebuilds the menu model with current Quick Play state.
func (a *App) rebuildMenu() {
	if a.config != nil && a.config.HasLastPlayed() {
		mode, _ := modes.Get(a.config.LastPlayedModeID)
		a.menuModel = screens.NewMenuWithQuickPlay(&screens.QuickPlayInfo{
			ModeName: mode.Name,
		})
	} else {
		a.menuModel = screens.NewMenu()
	}
	a.menuModel.SetSize(a.width, a.height)
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

// saveSession saves the completed session to statistics storage.
func (a *App) saveSession() {
	if a.session == nil || a.currentMode == nil {
		return
	}

	// Build the session record
	record, err := storage.NewSessionRecord(
		a.currentMode.Name,
		a.lastDifficulty.String(),
		int(a.lastDuration.Seconds()),
	)
	if err != nil {
		a.lastSaveError = err
		return
	}

	record.QuestionsAttempted = a.session.TotalAnswered() + a.session.Skipped
	record.QuestionsCorrect = a.session.Correct
	record.QuestionsWrong = a.session.Incorrect
	record.QuestionsSkipped = a.session.Skipped
	record.Score = a.session.Score
	record.BestStreak = a.session.BestStreak
	record.AvgResponseTimeMs = a.session.AvgResponseTime().Milliseconds()

	// Convert question history
	for _, h := range a.session.History {
		record.Questions = append(record.Questions, storage.QuestionRecord{
			Question:       h.Question,
			Operation:      h.Operation,
			CorrectAnswer:  h.CorrectAnswer,
			UserAnswer:     h.UserAnswer,
			Correct:        h.Correct,
			Skipped:        h.Skipped,
			ResponseTimeMs: h.ResponseTime.Milliseconds(),
			PointsEarned:   h.PointsEarned,
		})
	}

	// Save to storage - track error but don't disrupt gameplay flow
	a.lastSaveError = storage.AddSession(record)

	// Save last played settings for Quick Play
	a.saveLastPlayed()
}

// Phase 9: Replace main.go with Cobra CLI
// - arithmego (no args) → TUI menu
// - arithmego play → Quick play
// - arithmego statistics → Statistics screen
// - arithmego version → Version info
