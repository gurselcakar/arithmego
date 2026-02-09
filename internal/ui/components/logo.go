package components

import (
	"strings"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

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

// LogoColored returns the logo with the brand color applied.
func LogoColored() string {
	lines := strings.Split(Logo(), "\n")
	for i, line := range lines {
		lines[i] = styles.Logo.Render(line)
	}
	return strings.Join(lines, "\n")
}

// LogoColoredForWidth returns the colored logo based on terminal width.
func LogoColoredForWidth(width int) string {
	if width < LogoMinWidth {
		return styles.Logo.Render(styles.Bold.Render("ArithmeGo"))
	}
	return LogoColored()
}

// Tagline returns the game's tagline.
func Tagline() string {
	return styles.Tagline.Render("Your AI is thinking. You should too.")
}
