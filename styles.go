package main

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

	// Base styles
	baseStyle     = lipgloss.NewStyle().Padding(1).Margin(0)
	borderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(blueColor).Padding(0, 1)
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(blueColor).Underline(true)
	positiveStyle = lipgloss.NewStyle().Foreground(greenColor).Bold(true)
	negativeStyle = lipgloss.NewStyle().Foreground(redColor).Bold(true)
	neutralStyle  = lipgloss.NewStyle().Foreground(grayColor)

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

	helpStyle = lipgloss.NewStyle().Foreground(grayColor).Faint(true).MarginTop(1)
)

func renderBalance(balance float64) string {
	if balance >= 0 {
		return positiveStyle.Render("$" + formatAmount(balance))
	}
	return negativeStyle.Render("-$" + formatAmount(-balance))
}

func renderAmount(amount float64) string {
	return formatAmount(amount)
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
