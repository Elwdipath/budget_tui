package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func renderCategoryBar(category string, amount float64, total float64, width int) string {
	percentage := (amount / total) * 100
	barWidth := int((percentage / 100) * float64(width))

	if barWidth > width {
		barWidth = width
	}

	bar := strings.Repeat("█", barWidth) + strings.Repeat("░", width-barWidth)

	return fmt.Sprintf("%-15s %s %6.1f%%", category, bar, percentage)
}

func renderTransactionTable(transactions []Transaction, cursor int) string {
	if len(transactions) == 0 {
		return "No transactions yet.\n\nAdd your first transaction to get started!"
	}

	var sb strings.Builder

	// Header
	sb.WriteString("Date       Type     Amount    Description          Category\n")
	sb.WriteString("──────────────────────────────────────────────────────────\n")

	// Transactions (show max 10 for dashboard space)
	maxDisplay := 10
	if len(transactions) < maxDisplay {
		maxDisplay = len(transactions)
	}

	for i := 0; i < maxDisplay; i++ {
		t := transactions[i]
		cursorStr := " "
		if i == cursor {
			cursorStr = ">"
		}

		typeStr := "📈 Inc"
		if t.Type == Expense {
			typeStr = "📉 Exp"
		}

		amountStr := fmt.Sprintf("$%8.2f", t.Amount)
		if t.Type == Expense {
			amountStr = negativeStyle.Render("-" + fmt.Sprintf("$%7.2f", t.Amount))
		} else {
			amountStr = positiveStyle.Render(amountStr)
		}

		// Truncate description if too long
		description := t.Description
		if len(description) > 20 {
			description = description[:17] + "..."
		}

		sb.WriteString(fmt.Sprintf("%s %s %s %s %-20s %s\n",
			cursorStr,
			t.Date.Format("Jan 02"),
			typeStr,
			amountStr,
			description,
			t.Category))
	}

	if len(transactions) > maxDisplay {
		sb.WriteString(fmt.Sprintf("\n... and %d more transactions\n", len(transactions)-maxDisplay))
	}

	return sb.String()
}

func renderFinancialSummary(b *Budget) string {
	balance := b.GetBalance()
	income := b.GetTotalIncome()
	expenses := b.GetTotalExpenses()

	var sb strings.Builder

	// Balance row with appropriate styling
	balanceStr := renderBalance(balance)
	sb.WriteString(fmt.Sprintf("Balance:   %s\n", balanceStr))

	// Income row
	sb.WriteString(fmt.Sprintf("Income:    %s\n", positiveStyle.Render("$"+formatAmount(income))))

	// Expenses row
	sb.WriteString(fmt.Sprintf("Expenses:  %s\n", negativeStyle.Render("$"+formatAmount(expenses))))

	// Financial health status
	status := b.GetFinancialHealthStatus()
	var statusStyle lipgloss.Style
	switch {
	case strings.Contains(status, "✅"):
		statusStyle = positiveStyle
	case strings.Contains(status, "⚠️"):
		statusStyle = neutralStyle
	default:
		statusStyle = negativeStyle
	}

	sb.WriteString(fmt.Sprintf("\n%s", statusStyle.Render(status)))

	return sb.String()
}
