package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color scheme
	greenColor  = lipgloss.Color("42")
	redColor    = lipgloss.Color("196")
	yellowColor = lipgloss.Color("226")
	blueColor   = lipgloss.Color("39")
	grayColor   = lipgloss.Color("245")

	// Styles
	positiveStyle = lipgloss.NewStyle().Foreground(greenColor).Bold(true)
	negativeStyle = lipgloss.NewStyle().Foreground(redColor).Bold(true)
	neutralStyle  = lipgloss.NewStyle().Foreground(yellowColor)
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(blueColor).Underline(true)
	borderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
	helpStyle     = lipgloss.NewStyle().Foreground(grayColor)

	// Panel styles
	summaryPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1, 2).
				Width(60).
				Height(7)

	categoryPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1).
				Width(45).
				Height(12)

	transactionsPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1).
				Width(55).
				Height(12)
)

func GetHeroBanner() string {
	bannerText := `
███████╗ █████╗ ██╗   ██╗████████╗ ██████╗ ██████╗ ██╗     ██╗
██╔════╝██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔══██╗██║     ██║
█████╗  ███████║██║   ██║   ██║   ██║   ██║██████╔╝██║     ██║
██╔══╝  ██╔══██║╚██╗ ██╔╝   ██║   ██║   ██║██╔══██╗██║     ██║
██║     ██║  ██║ ╚████╔╝    ██║   ╚██████╔╝██║  ██║███████╗██║
╚═╝     ╚═╝  ╚═╝  ╚═══╝     ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝
`

	heroStyle := lipgloss.NewStyle().
		Foreground(blueColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(80)

	return heroStyle.Render(bannerText)
}

func GetSubtitle() string {
	subtitleStyle := lipgloss.NewStyle().
		Foreground(grayColor).
		Italic(true).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(2)

	return subtitleStyle.Render("Your Personal Finance Dashboard")
}

// Export styles and colors for use in main
func GetPositiveStyle() lipgloss.Style          { return positiveStyle }
func GetNegativeStyle() lipgloss.Style          { return negativeStyle }
func GetNeutralStyle() lipgloss.Style           { return neutralStyle }
func GetTitleStyle() lipgloss.Style             { return titleStyle }
func GetBorderStyle() lipgloss.Style            { return borderStyle }
func GetHelpStyle() lipgloss.Style              { return helpStyle }
func GetSummaryPanelStyle() lipgloss.Style      { return summaryPanelStyle }
func GetCategoryPanelStyle() lipgloss.Style     { return categoryPanelStyle }
func GetTransactionsPanelStyle() lipgloss.Style { return transactionsPanelStyle }
func GetBlueColor() lipgloss.Color              { return blueColor }
func GetGrayColor() lipgloss.Color              { return grayColor }

func FormatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
