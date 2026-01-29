package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gurselcakar/arithmego/internal/ui"

	// Register operations
	_ "github.com/gurselcakar/arithmego/internal/game/operations"
)

// Phase 9: Replace with Cobra CLI
// - arithmego (no args) → TUI menu
// - arithmego play → Quick play
// - arithmego statistics → Statistics screen
// - arithmego version → Version info

func main() {
	app := ui.New()
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
