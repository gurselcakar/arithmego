package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/gen"
	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/screens"
	"github.com/gurselcakar/arithmego/internal/update"
)

// App is the main Bubble Tea model that orchestrates all screens.
type App struct {
	screen Screen
	width  int
	height int

	// Screen models
	menuModel        screens.MenuModel
	playBrowseModel  screens.PlayBrowseModel
	playConfigModel  screens.PlayConfigModel
	gameModel        screens.GameModel
	pauseModel       screens.PauseModel
	resultsModel     screens.ResultsModel
	practiceModel    screens.PracticeModel
	statisticsModel  screens.StatisticsModel
	settingsModel    screens.SettingsModel
	onboardingModel  screens.OnboardingModel
	quitConfirmModel screens.QuitConfirmModel
	featureTourModel screens.FeatureTourModel

	// Current session state
	session         *game.Session
	currentMode     *modes.Mode
	lastDifficulty  game.Difficulty
	lastDuration    time.Duration
	lastInputMethod components.InputMethod

	// User config (for Quick Play and defaults)
	config *storage.Config

	// Error tracking
	lastSaveError error

	// CLI start mode flags
	cliStartMode StartMode

	// First game tracking (for feature tour after onboarding)
	isFirstGame bool

	// Update notification
	updateInfo         *update.Info
	autoUpdateInstalled string // Version that was auto-updated (empty if none)
}

// New creates a new App instance with default start mode.
func New() *App {
	return NewWithStartMode(StartModeMenu)
}

// NewWithStartMode creates a new App instance with the specified start mode.
func NewWithStartMode(startMode StartMode) *App {
	// Load config (ignore errors, use default config as fallback)
	config, _ := storage.LoadConfig()
	if config == nil {
		config = storage.NewConfig()
	}

	// Load practice settings from config
	var practiceSettings *screens.PracticeSettings
	if config.PracticeCategory != "" {
		practiceSettings = &screens.PracticeSettings{
			Category:    config.PracticeCategory,
			Operation:   config.PracticeOperation,
			Difficulty:  config.PracticeDifficulty,
			InputMethod: config.PracticeInputMethod,
		}
	}

	app := &App{
		menuModel:       screens.NewMenu(),
		practiceModel:   screens.NewPracticeWithSettings(practiceSettings),
		statisticsModel: screens.NewStatistics(),
		settingsModel:   screens.NewSettings(config),
		onboardingModel: screens.NewOnboarding(),
		config:          config,
	}

	// Determine starting screen based on start mode
	switch startMode {
	case StartModePlayBrowse:
		// Play browse - will be initialized in Init()
		app.screen = ScreenMenu
		app.cliStartMode = StartModePlayBrowse

	case StartModePlayConfig:
		// Play config with specific mode - will be initialized in Init()
		app.screen = ScreenMenu
		app.cliStartMode = StartModePlayConfig

	case StartModeStatistics:
		app.screen = ScreenStatistics

	case StartModeSettings:
		app.screen = ScreenSettings

	case StartModeOnboarding:
		app.screen = ScreenOnboarding
	default:
		// Default menu behavior: check onboarding and tour status
		if !config.Onboarded {
			app.screen = ScreenOnboarding
		} else if !config.TourCompleted {
			app.featureTourModel = screens.NewFeatureTour()
			app.screen = ScreenFeatureTour
		} else {
			app.screen = ScreenMenu
		}
	}

	return app
}

// Init initializes the app.
func (a *App) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Handle CLI start modes
	switch a.cliStartMode {
	case StartModePlayBrowse:
		a.cliStartMode = StartModeMenu // Reset flag
		cmds = append(cmds, func() tea.Msg {
			return cliPlayBrowseMsg{}
		})
	case StartModePlayConfig:
		a.cliStartMode = StartModeMenu // Reset flag
		cmds = append(cmds, func() tea.Msg {
			return cliPlayConfigMsg{modeID: CLIModeID}
		})
	}

	// Handle statistics screen init
	if a.screen == ScreenStatistics {
		cmds = append(cmds, a.statisticsModel.Init())
	}

	// Check for updates if auto_update is enabled
	if a.config != nil && a.config.AutoUpdate {
		cmds = append(cmds, checkForUpdateCmd())
	}

	return tea.Batch(cmds...)
}

// cliPlayBrowseMsg triggers play browse from CLI.
type cliPlayBrowseMsg struct{}

// cliPlayConfigMsg triggers play config with a specific mode from CLI.
type cliPlayConfigMsg struct {
	modeID string
}

// updateCheckResultMsg carries the result of an update check.
type updateCheckResultMsg struct {
	info *update.Info
	err  error
}

// Version is the current app version, set by the CLI before starting the TUI.
var Version = "dev"

// CLIModeID is the mode ID specified via CLI, set before starting the TUI.
var CLIModeID = ""

// autoUpdateResultMsg carries the result of an auto-update attempt.
type autoUpdateResultMsg struct {
	version string
	err     error
}

// checkForUpdateCmd returns a command that checks for updates.
func checkForUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		info, err := update.Check(Version)
		return updateCheckResultMsg{info: info, err: err}
	}
}

// autoUpdateCmd returns a command that downloads and applies an update.
func autoUpdateCmd(version string) tea.Cmd {
	return func() tea.Msg {
		err := update.DownloadAndApply(version)
		return autoUpdateResultMsg{version: version, err: err}
	}
}

// Update handles all messages and routes them to the appropriate screen.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size for all screens
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		a.width = wsm.Width
		a.height = wsm.Height
	}

	// Handle CLI play browse trigger
	if _, ok := msg.(cliPlayBrowseMsg); ok {
		return a.startPlayBrowse()
	}

	// Handle CLI play config trigger
	if configMsg, ok := msg.(cliPlayConfigMsg); ok {
		return a.startPlayConfig(configMsg.modeID)
	}

	// Handle update check result
	if updateMsg, ok := msg.(updateCheckResultMsg); ok {
		if updateMsg.err == nil && updateMsg.info != nil && updateMsg.info.UpdateAvailable {
			a.updateInfo = updateMsg.info
			// Attempt auto-update in the background
			return a, autoUpdateCmd(updateMsg.info.LatestVersion)
		}
		return a, nil
	}

	// Handle auto-update result
	if updateMsg, ok := msg.(autoUpdateResultMsg); ok {
		if updateMsg.err == nil {
			// Auto-update succeeded
			a.autoUpdateInstalled = updateMsg.version
			a.menuModel.SetUpdateInstalled(updateMsg.version)
		} else {
			// Auto-update failed, fall back to manual notification
			if a.updateInfo != nil {
				a.menuModel.SetUpdateInfo(a.updateInfo.LatestVersion)
			}
		}
		return a, nil
	}

	// Route based on current screen
	switch a.screen {
	case ScreenMenu:
		return a.updateMenu(msg)
	case ScreenPlayBrowse:
		return a.updatePlayBrowse(msg)
	case ScreenPlayConfig:
		return a.updatePlayConfig(msg)
	case ScreenGame:
		return a.updateGame(msg)
	case ScreenPause:
		return a.updatePause(msg)
	case ScreenResults:
		return a.updateResults(msg)
	case ScreenPractice:
		return a.updatePractice(msg)
	case ScreenStatistics:
		return a.updateStatistics(msg)
	case ScreenSettings:
		return a.updateSettings(msg)
	case ScreenOnboarding:
		return a.updateOnboarding(msg)
	case ScreenQuitConfirm:
		return a.updateQuitConfirm(msg)
	case ScreenFeatureTour:
		return a.updateFeatureTour(msg)
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
		case screens.ActionPlay:
			a.playBrowseModel = screens.NewPlayBrowse(a.config)
			a.playBrowseModel.SetSize(a.width, a.height)
			a.screen = ScreenPlayBrowse
			return a, a.playBrowseModel.Init()
		case screens.ActionPractice:
			// Load practice settings from config
			var practiceSettings *screens.PracticeSettings
			if a.config.PracticeCategory != "" {
				practiceSettings = &screens.PracticeSettings{
					Category:    a.config.PracticeCategory,
					Operation:   a.config.PracticeOperation,
					Difficulty:  a.config.PracticeDifficulty,
					InputMethod: a.config.PracticeInputMethod,
				}
			}
			a.practiceModel = screens.NewPracticeWithSettings(practiceSettings)
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
		case screens.ActionX:
			return a, openURL("https://x.com/gurselcakar")
		}
	}

	// Check if quitting
	if a.menuModel.Quitting() {
		return a, tea.Quit
	}

	return a, cmd
}

// updatePlayBrowse handles play browse screen updates.
func (a *App) updatePlayBrowse(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.playBrowseModel, cmd = a.playBrowseModel.Update(msg)

	// Check for mode selection
	if selectMsg, ok := msg.(screens.ModeSelectedMsg); ok {
		a.playConfigModel = screens.NewPlayConfig(selectMsg.Mode, a.config)
		a.playConfigModel.SetSize(a.width, a.height)
		a.screen = ScreenPlayConfig
		return a, a.playConfigModel.Init()
	}

	// Check for return to menu
	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		return a.returnToMenu()
	}

	return a, cmd
}

// updatePlayConfig handles play config screen updates.
func (a *App) updatePlayConfig(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.playConfigModel, cmd = a.playConfigModel.Update(msg)

	// Check for start game
	if startMsg, ok := msg.(screens.StartGameMsg); ok {
		a.currentMode = startMsg.Mode
		a.lastDifficulty = startMsg.Difficulty
		a.lastDuration = startMsg.Duration
		a.lastInputMethod = startMsg.InputMethod
		return a.startGame()
	}

	// Check for back to browse — reuse existing model to preserve cursor position
	if _, ok := msg.(screens.BackToBrowseMsg); ok {
		if len(a.playBrowseModel.Modes()) == 0 {
			// Browse model was never initialized (e.g., CLI launched directly into config)
			a.playBrowseModel = screens.NewPlayBrowse(a.config)
		}
		a.playBrowseModel.SetSize(a.width, a.height)
		a.screen = ScreenPlayBrowse
		return a, nil
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
		// Use first game results if this is the onboarding game
		if a.isFirstGame {
			a.resultsModel = screens.NewResultsFirstGame(a.session, a.lastSaveError)
		} else {
			a.resultsModel = screens.NewResults(a.session, a.lastSaveError)
		}
		a.resultsModel.SetSize(a.width, a.height)
		a.screen = ScreenResults
		return a, a.resultsModel.Init()
	}

	// Check for pause
	if pm, ok := msg.(screens.PauseMsg); ok {
		a.session = pm.Session
		a.pauseModel = screens.NewPause(a.session, a.config)
		a.pauseModel.SetSize(a.width, a.height)
		a.screen = ScreenPause
		return a, a.pauseModel.Init()
	}

	// Check for quit confirmation from game
	if qm, ok := msg.(screens.QuitConfirmMsg); ok {
		a.session = qm.Session
		// Check if user has disabled quit confirmation
		if a.config != nil && a.config.SkipQuitConfirmation {
			// If first game, go to feature tour instead of menu
			if a.isFirstGame {
				return a.startFeatureTour()
			}
			a.rebuildMenu()
			a.screen = ScreenMenu
			a.session = nil
			return a, nil
		}
		// Show quit confirmation screen
		a.quitConfirmModel = screens.NewQuitConfirm(a.session, a.config, screens.QuitFromGame)
		a.quitConfirmModel.SetSize(a.width, a.height)
		a.screen = ScreenQuitConfirm
		return a, a.quitConfirmModel.Init()
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

	// Check for quit to menu (direct quit, skipping confirmation)
	if _, ok := msg.(screens.QuitToMenuMsg); ok {
		// If first game, go to feature tour instead of menu
		if a.isFirstGame {
			return a.startFeatureTour()
		}
		a.rebuildMenu()
		a.screen = ScreenMenu
		a.session = nil
		return a, nil
	}

	// Check for quit confirmation from pause
	if qm, ok := msg.(screens.QuitConfirmMsg); ok {
		a.session = qm.Session
		// Check if user has disabled quit confirmation
		if a.config != nil && a.config.SkipQuitConfirmation {
			// If first game, go to feature tour instead of menu
			if a.isFirstGame {
				return a.startFeatureTour()
			}
			a.rebuildMenu()
			a.screen = ScreenMenu
			a.session = nil
			return a, nil
		}
		// Show quit confirmation screen
		a.quitConfirmModel = screens.NewQuitConfirm(a.session, a.config, screens.QuitFromPause)
		a.quitConfirmModel.SetSize(a.width, a.height)
		a.screen = ScreenQuitConfirm
		return a, a.quitConfirmModel.Init()
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
		a.session = nil
		return a.returnToMenu()
	}

	// Check for continue to feature tour (first game only)
	if _, ok := msg.(screens.ContinueToFeatureTourMsg); ok {
		return a.startFeatureTour()
	}

	return a, cmd
}

// updatePractice handles practice screen updates.
func (a *App) updatePractice(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.practiceModel, cmd = a.practiceModel.Update(msg)

	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		// Save practice settings before returning to menu
		a.savePracticeSettings()
		return a.returnToMenu()
	}

	return a, cmd
}

// savePracticeSettings saves the current practice settings to config.
func (a *App) savePracticeSettings() {
	settings := a.practiceModel.Settings()
	a.config.PracticeCategory = settings.Category
	a.config.PracticeOperation = settings.Operation
	a.config.PracticeDifficulty = settings.Difficulty
	a.config.PracticeInputMethod = settings.InputMethod
	_ = storage.SaveConfig(a.config) // Ignore save errors for non-critical data
}

// updateStatistics handles statistics screen updates.
func (a *App) updateStatistics(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.statisticsModel, cmd = a.statisticsModel.Update(msg)

	if _, ok := msg.(screens.ReturnToMenuMsg); ok {
		return a.returnToMenu()
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
		return a.returnToMenu()
	}

	return a, cmd
}

// updateOnboarding handles onboarding screen updates.
func (a *App) updateOnboarding(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.onboardingModel, cmd = a.onboardingModel.Update(msg)

	// Handle onboarding completion with user selections
	if completeMsg, ok := msg.(screens.OnboardingCompleteMsg); ok {
		return a.completeOnboarding(completeMsg.ModeID, completeMsg.Difficulty, completeMsg.DurationMs, completeMsg.InputMethod)
	}

	// Handle onboarding skip - use defaults (Easy, 60s, Addition, Typing)
	if _, ok := msg.(screens.OnboardingSkipMsg); ok {
		return a.completeOnboarding(modes.IDAddition, "Easy", 60000, "typing")
	}

	return a, cmd
}

// updateQuitConfirm handles quit confirmation screen updates.
func (a *App) updateQuitConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.quitConfirmModel, cmd = a.quitConfirmModel.Update(msg)

	// Handle cancel - return to previous screen
	if cancelMsg, ok := msg.(screens.QuitConfirmCancelMsg); ok {
		a.session = cancelMsg.Session
		if cancelMsg.Source == screens.QuitFromGame {
			// Return to game and resume timer
			a.gameModel.SetSession(a.session)
			a.screen = ScreenGame
			a.session.Resume()
			return a, screens.TickCmd()
		}
		// Return to pause screen (timer already stopped)
		a.pauseModel.SetSession(a.session)
		a.screen = ScreenPause
		return a, nil
	}

	// Handle accept - save preference if checked and go to menu (or feature tour for first game)
	if acceptMsg, ok := msg.(screens.QuitConfirmAcceptMsg); ok {
		if acceptMsg.DontAskAgain && a.config != nil {
			a.config.SkipQuitConfirmation = true
			_ = storage.SaveConfig(a.config)
		}
		// If first game, go to feature tour instead of menu
		if a.isFirstGame {
			return a.startFeatureTour()
		}
		a.rebuildMenu()
		a.screen = ScreenMenu
		a.session = nil
		return a, nil
	}

	return a, cmd
}

// updateFeatureTour handles feature tour screen updates.
func (a *App) updateFeatureTour(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.featureTourModel, cmd = a.featureTourModel.Update(msg)

	// Check for feature tour completion (or skip)
	if _, ok := msg.(screens.FeatureTourCompleteMsg); ok {
		a.isFirstGame = false
		a.config.TourCompleted = true
		_ = storage.SaveConfig(a.config)
		a.rebuildMenu()
		a.screen = ScreenMenu
		a.session = nil
		return a, nil
	}

	return a, cmd
}

// startFeatureTour creates the feature tour and transitions to it.
func (a *App) startFeatureTour() (tea.Model, tea.Cmd) {
	// Reset first game flag - single point of truth
	a.isFirstGame = false
	// Save session if it exists (user quit mid-game)
	if a.session != nil {
		a.saveSession()
	}
	a.featureTourModel = screens.NewFeatureTour()
	a.featureTourModel.SetSize(a.width, a.height)
	a.screen = ScreenFeatureTour
	a.session = nil
	return a, a.featureTourModel.Init()
}

// completeOnboarding finishes onboarding and starts the game with selected settings.
func (a *App) completeOnboarding(modeID, difficulty string, durationMs int64, inputMethod string) (tea.Model, tea.Cmd) {
	// Mark as onboarded
	a.config.Onboarded = true

	// Set last played and default settings
	a.config.LastPlayedModeID = modeID
	a.config.LastPlayedDifficulty = difficulty
	a.config.LastPlayedDurationMs = durationMs
	a.config.DefaultDifficulty = difficulty
	a.config.DefaultDurationMs = durationMs
	a.config.InputMethod = inputMethod

	// Save config (ignore errors - config is non-critical)
	_ = storage.SaveConfig(a.config)

	// Mark this as the first game (for feature tour after completion)
	a.isFirstGame = true

	// Set up game state
	mode, ok := modes.Get(modeID)
	if !ok || mode == nil {
		// Mode doesn't exist - fall back to Play Browse screen
		a.isFirstGame = false
		a.playBrowseModel = screens.NewPlayBrowse(a.config)
		a.playBrowseModel.SetSize(a.width, a.height)
		a.screen = ScreenPlayBrowse
		return a, a.playBrowseModel.Init()
	}
	a.currentMode = mode
	a.lastDifficulty = game.ParseDifficulty(difficulty)
	a.lastDuration = time.Duration(durationMs) * time.Millisecond
	a.lastInputMethod = components.ParseInputMethod(inputMethod)

	// Start the game
	return a.startGame()
}

// startGame creates a new session and starts the game.
// If mode is invalid, gracefully returns to the Play Browse screen.
func (a *App) startGame() (tea.Model, tea.Cmd) {
	if a.currentMode == nil || a.currentMode.GeneratorLabel == "" {
		// Gracefully recover: return to Play Browse screen instead of crashing
		a.isFirstGame = false // Reset flag on error
		a.playBrowseModel = screens.NewPlayBrowse(a.config)
		a.playBrowseModel.SetSize(a.width, a.height)
		a.screen = ScreenPlayBrowse
		return a, a.playBrowseModel.Init()
	}

	g, ok := gen.Get(a.currentMode.GeneratorLabel)
	if !ok {
		a.playBrowseModel = screens.NewPlayBrowse(a.config)
		a.playBrowseModel.SetSize(a.width, a.height)
		a.screen = ScreenPlayBrowse
		return a, a.playBrowseModel.Init()
	}
	a.session = game.NewSession(g, a.lastDifficulty, a.lastDuration)
	a.gameModel = screens.NewGame(a.session, a.lastInputMethod)
	a.gameModel.SetSize(a.width, a.height)
	a.screen = ScreenGame
	return a, a.gameModel.Init()
}

// startPlayBrowse opens the play browse screen from CLI.
func (a *App) startPlayBrowse() (tea.Model, tea.Cmd) {
	a.playBrowseModel = screens.NewPlayBrowse(a.config)
	a.playBrowseModel.SetSize(a.width, a.height)
	a.screen = ScreenPlayBrowse
	return a, a.playBrowseModel.Init()
}

// startPlayConfig opens the play config screen with a specific mode from CLI.
func (a *App) startPlayConfig(modeID string) (tea.Model, tea.Cmd) {
	mode, ok := modes.Get(modeID)
	if !ok || mode == nil {
		// Mode not found - fall back to Play Browse screen
		return a.startPlayBrowse()
	}

	a.playConfigModel = screens.NewPlayConfig(mode, a.config)
	a.playConfigModel.SetSize(a.width, a.height)
	a.screen = ScreenPlayConfig
	return a, a.playConfigModel.Init()
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
	if a.lastInputMethod == components.InputMultipleChoice {
		a.config.InputMethod = "multiple_choice"
	} else {
		a.config.InputMethod = "typing"
	}

	// Ignore save errors - config is non-critical
	_ = storage.SaveConfig(a.config)
}

// returnToMenu transitions to the menu screen, ensuring it's properly initialized.
func (a *App) returnToMenu() (tea.Model, tea.Cmd) {
	a.menuModel.SetSize(a.width, a.height)
	a.screen = ScreenMenu
	return a, nil
}

// rebuildMenu rebuilds the menu model.
func (a *App) rebuildMenu() {
	a.menuModel = screens.NewMenu()
	a.menuModel.SetSize(a.width, a.height)

	// Re-apply update state
	if a.autoUpdateInstalled != "" {
		a.menuModel.SetUpdateInstalled(a.autoUpdateInstalled)
	} else if a.updateInfo != nil && a.updateInfo.UpdateAvailable {
		a.menuModel.SetUpdateInfo(a.updateInfo.LatestVersion)
	}
}

// View renders the current screen.
func (a *App) View() string {
	switch a.screen {
	case ScreenMenu:
		return a.menuModel.View()
	case ScreenPlayBrowse:
		return a.playBrowseModel.View()
	case ScreenPlayConfig:
		return a.playConfigModel.View()
	case ScreenGame:
		return a.gameModel.View()
	case ScreenPause:
		return a.pauseModel.View()
	case ScreenResults:
		return a.resultsModel.View()
	case ScreenPractice:
		return a.practiceModel.View()
	case ScreenStatistics:
		return a.statisticsModel.View()
	case ScreenSettings:
		return a.settingsModel.View()
	case ScreenOnboarding:
		return a.onboardingModel.View()
	case ScreenQuitConfirm:
		return a.quitConfirmModel.View()
	case ScreenFeatureTour:
		return a.featureTourModel.View()
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
