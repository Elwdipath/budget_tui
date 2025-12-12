package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func getHeroBanner() string {
	bannerText := `
███████╗ █████╗ ██╗   ██╗████████╗ ██████╗ ██████╗ ██╗     ██╗
██╔════╝██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔══██╗██║     ██║
█████╗  ███████║██║   ██║   ██║   ██║   ██║██████╔╝██║     ██║
██╔══╝  ██╔══██║╚██╗ ██╔╝   ██║   ██║   ██║██╔══██╗██║     ██║
██║     ██║  ██║ ╚████╔╝    ██║   ╚██████╔╝██║  ██║███████╗██║
╚═╝     ╚═╝  ╚═╝  ╚═══╝     ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝
`

	// Style for the hero banner
	heroStyle := lipgloss.NewStyle().
		Foreground(blueColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(80)

	return heroStyle.Render(bannerText)
}

func getSubtitle() string {
	subtitleStyle := lipgloss.NewStyle().
		Foreground(grayColor).
		Italic(true).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(2)

	return subtitleStyle.Render("Your Personal Finance Dashboard")
}
