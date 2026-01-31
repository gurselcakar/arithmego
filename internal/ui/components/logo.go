package components

import "github.com/gurselcakar/arithmego/internal/ui/styles"

// LogoMinWidth is the minimum terminal width for the full logo (75 chars + margin).
const LogoMinWidth = 80

// Logo returns the ASCII art logo for ArithmeGo.
// The full logo is 75 characters wide.
func Logo() string {
	return `█████╗ ██████╗ ██╗████████╗██╗  ██╗███╗   ███╗███████╗ ██████╗  ██████╗
██╔══██╗██╔══██╗██║╚══██╔══╝██║  ██║████╗ ████║██╔════╝██╔════╝ ██╔═══██╗
███████║██████╔╝██║   ██║   ███████║██╔████╔██║█████╗  ██║  ███╗██║   ██║
██╔══██║██╔══██╗██║   ██║   ██╔══██║██║╚██╔╝██║██╔══╝  ██║   ██║██║   ██║
██║  ██║██║  ██║██║   ██║   ██║  ██║██║ ╚═╝ ██║███████╗╚██████╔╝╚██████╔╝
╚═╝  ╚═╝╚═╝  ╚═╝╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝ ╚═════╝  ╚═════╝ `
}

// LogoSeparator returns the three dots separator shown below the logo.
func LogoSeparator() string {
	return "•  •  •"
}

// LogoCompact returns a smaller text-based logo for narrow terminals.
// Use this when terminal width is less than LogoMinWidth.
func LogoCompact() string {
	return styles.Bold.Render("ArithmeGo")
}

// LogoForWidth returns the appropriate logo based on terminal width.
func LogoForWidth(width int) string {
	if width < LogoMinWidth {
		return LogoCompact()
	}
	return Logo()
}

// Tagline returns the game's tagline with dimmed styling.
func Tagline() string {
	return styles.Dim.Render("Your AI is thinking. You should too.")
}
