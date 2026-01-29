package components

import "github.com/gurselcakar/arithmego/internal/ui/styles"

// Logo returns the ASCII art logo for ArithmeGo.
func Logo() string {
	return `█████╗ ██████╗ ██╗████████╗██╗  ██╗███╗   ███╗███████╗ ██████╗  ██████╗
██╔══██╗██╔══██╗██║╚══██╔══╝██║  ██║████╗ ████║██╔════╝██╔════╝ ██╔═══██╗
███████║██████╔╝██║   ██║   ███████║██╔████╔██║█████╗  ██║  ███╗██║   ██║
██╔══██║██╔══██╗██║   ██║   ██╔══██║██║╚██╔╝██║██╔══╝  ██║   ██║██║   ██║
██║  ██║██║  ██║██║   ██║   ██║  ██║██║ ╚═╝ ██║███████╗╚██████╔╝╚██████╔╝
╚═╝  ╚═╝╚═╝  ╚═╝╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝ ╚═════╝  ╚═════╝`
}

// Tagline returns the game's tagline with dimmed styling.
func Tagline() string {
	return styles.Dim.Render("Your AI is thinking. You should too.")
}
