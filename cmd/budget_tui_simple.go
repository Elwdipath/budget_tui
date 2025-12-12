package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Import all the packages
import (
	"github.com/Elwdipath/budget_tui/internal/budget"
	"github.com/Elwdipath/budget_tui/internal/storage"
	"github.com/Elwdipath/budget_tui/internal/tui"
)

type state int

const (
	dashboardState state = iota
	addIncomeState
	addExpenseState
	viewTransactionsState
)

type model struct {
	state       state
	budget      *budget.Budget
	menuChoices []string
	menuCursor  int

	// Form fields
	amountInput      string
	descriptionInput string
	categoryInput    string
	activeField      int
	formSubmitted    bool

	// Dashboard state
	dashboardCursor     int
	selectedCategory    int
	selectedTransaction int
	showHelp            bool
}

func initialModel() model {
	b, _ := storage.LoadBudget()
	return model{
		state:  dashboardState,
		budget: b,
		menuChoices: []string{
			"[i] Income  [e] Expense  [t] Transactions  [h] Help  [q] Quit",
		},
		menuCursor:          0,
		activeField:         0,
		dashboardCursor:     0,
		selectedCategory:    0,
		selectedTransaction: 0,
		showHelp:            false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case dashboardState:
			return m.updateDashboard(msg)
		case addIncomeState, addExpenseState:
			return m.updateAddTransactionForm(msg)
		case viewTransactionsState:
			return m.updateViewTransactions(msg)
		}
	}
	return m, nil
}

func (m model) updateDashboard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "i":
		m.state = addIncomeState
		m.resetForm()
	case "e":
		m.state = addExpenseState
		m.resetForm()
	case "t":
		m.state = viewTransactionsState
		m.selectedTransaction = 0
	case "h":
		m.showHelp = !m.showHelp
	case "up", "k":
		if m.dashboardCursor > 0 {
			m.dashboardCursor--
		}
	case "down", "j":
		categories := m.budget.GetSpendingByCategory()
		if m.dashboardCursor < len(categories)-1 {
			m.dashboardCursor++
		}
	}
	return m, nil
}

func (m model) updateAddTransactionForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = dashboardState
	case "tab":
		m.activeField = (m.activeField + 1) % 3
	case "shift+tab":
		m.activeField = (m.activeField - 1 + 3) % 3
	case "enter":
		if m.amountInput != "" && m.descriptionInput != "" && m.categoryInput != "" {
			var amount float64
			fmt.Sscanf(m.amountInput, "%f", &amount)

			tType := budget.Income
			if m.state == addExpenseState {
				tType = budget.Expense
			}

			m.budget.AddTransaction(amount, m.descriptionInput, m.categoryInput, tType)
			m.budget.Save()
			m.formSubmitted = true
			m.state = dashboardState
		}
	case "backspace":
		switch m.activeField {
		case 0:
			if len(m.amountInput) > 0 {
				m.amountInput = m.amountInput[:len(m.amountInput)-1]
			}
		case 1:
			if len(m.descriptionInput) > 0 {
				m.descriptionInput = m.descriptionInput[:len(m.descriptionInput)-1]
			}
		case 2:
			if len(m.categoryInput) > 0 {
				m.categoryInput = m.categoryInput[:len(m.categoryInput)-1]
			}
		}
	default:
		if len(msg.String()) == 1 {
			switch m.activeField {
			case 0:
				m.amountInput += msg.String()
			case 1:
				m.descriptionInput += msg.String()
			case 2:
				m.categoryInput += msg.String()
			}
		}
	}
	return m, nil
}

func (m model) updateViewTransactions(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = dashboardState
	case "up", "k":
		if m.selectedTransaction > 0 {
			m.selectedTransaction--
		}
	case "down", "j":
		if m.selectedTransaction < len(m.budget.Transactions)-1 {
			m.selectedTransaction++
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case dashboardState:
		return m.viewDashboard()
	case addIncomeState, addExpenseState:
		return m.viewAddTransactionForm()
	case viewTransactionsState:
		return m.viewTransactions()
	default:
		return ""
	}
}

func (m model) resetForm() {
	m.amountInput = ""
	m.descriptionInput = ""
	m.categoryInput = ""
	m.activeField = 0
	m.formSubmitted = false
}

func (m model) viewDashboard() string {
	// Use the TUI package components
	heroBanner := tui.GetHeroBanner()
	subtitle := tui.GetSubtitle()

	// Financial Summary Panel
	summaryContent := m.renderFinancialSummary()
	summaryPanel := tui.GetSummaryPanelStyle().Render(summaryContent)

	// Category Spending Panel
	categories := m.budget.GetSpendingByCategory()
	var categoryContent strings.Builder

	if len(categories) == 0 {
		categoryContent.WriteString("No expense categories yet.\n\nAdd some expenses to see breakdown.")
	} else {
		categoryContent.WriteString("Top Spending Categories:\n\n")

		totalExpenses := m.budget.GetTotalExpenses()
		maxDisplay := 5
		if len(categories) < maxDisplay {
			maxDisplay = len(categories)
		}

		for i := 0; i < maxDisplay; i++ {
			cat := categories[i]
			percentage := (cat.Amount / totalExpenses) * 100
			bar := strings.Repeat("█", int(percentage/5)) + strings.Repeat("░", 20-int(percentage/5))

			categoryContent.WriteString(fmt.Sprintf("%-15s %s %.1f%% ($%.2f)\n\n",
				cat.Category, bar, percentage, cat.Amount))
		}
	}
	categoryPanel := tui.GetCategoryPanelStyle().Render(categoryContent.String())

	// Recent Transactions Panel
	recentTransactions := m.budget.GetRecentTransactions(8)
	transactionsContent := m.renderTransactionTable(recentTransactions)
	transactionsPanel := tui.GetTransactionsPanelStyle().Render(transactionsContent)

	// Layout panels
	title := tui.GetTitleStyle().Render("📊 Budget Dashboard")

	topRow := lipgloss.JoinHorizontal(lipgloss.Left, summaryPanel)
	middleRow := lipgloss.JoinHorizontal(lipgloss.Left, categoryPanel, " ", transactionsPanel)

	dashboard := lipgloss.JoinVertical(lipgloss.Top, heroBanner, subtitle, title, topRow, "\n", middleRow)

	// Menu bar
	menuBar := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Foreground(tui.GetBlueColor()).
		Render(strings.Join(m.menuChoices, " • "))

	// Help text
	var helpText string
	if m.showHelp {
		helpText = "\n" + lipgloss.NewStyle().Foreground(tui.GetGrayColor()).Faint(true).Render("Controls: ↑↓/j/k: Navigate categories | i: Add Income | e: Add Expense | t: View All Transactions | h: Toggle help | q: Quit")
	}

	fullDashboard := lipgloss.JoinVertical(lipgloss.Top, dashboard, "\n", menuBar)

	return lipgloss.NewStyle().Padding(1, 0, 0, 0).Render(fullDashboard + helpText)
}

func (m model) renderFinancialSummary() string {
	balance := m.budget.GetBalance()
	income := m.budget.GetTotalIncome()
	expenses := m.budget.GetTotalExpenses()

	var sb strings.Builder

	// Balance row with appropriate styling
	balanceStr := m.renderBalance(balance)
	sb.WriteString(fmt.Sprintf("Balance:   %s\n", balanceStr))

	// Income row
	sb.WriteString(fmt.Sprintf("Income:    %s\n", tui.GetPositiveStyle().Render("$"+tui.FormatAmount(income))))

	// Expenses row
	sb.WriteString(fmt.Sprintf("Expenses:  %s\n", tui.GetNegativeStyle().Render("$"+tui.FormatAmount(expenses))))

	// Financial health status
	status := "✅ Financially Healthy"
	if balance < 0 {
		status = "🚨 Budget Alert"
	} else if balance < 100 {
		status = "⚠️  Watch Your Spending"
	}

	var statusStyle lipgloss.Style
	if balance > 0 {
		statusStyle = tui.GetPositiveStyle()
	} else if balance > -100 {
		statusStyle = lipgloss.NewStyle().Foreground(tui.GetGrayColor())
	} else {
		statusStyle = tui.GetNegativeStyle()
	}

	sb.WriteString(fmt.Sprintf("\n%s", statusStyle.Render(status)))

	return sb.String()
}

func (m model) renderBalance(balance float64) string {
	if balance >= 0 {
		return tui.GetPositiveStyle().Render("$" + tui.FormatAmount(balance))
	}
	return tui.GetNegativeStyle().Render("-$" + tui.FormatAmount(-balance))
}

func (m model) renderTransactionTable(transactions []budget.Transaction) string {
	if len(transactions) == 0 {
		return "No transactions yet.\n\nAdd your first transaction to get started!"
	}

	var sb strings.Builder

	// Header
	sb.WriteString("Date       Type     Amount    Description          Category\n")
	sb.WriteString("──────────────────────────────────────────────────────────\n")

	for i, t := range transactions {
		cursor := " "
		if i == m.selectedTransaction {
			cursor = ">"
		}

		typeStr := "📈 Inc"
		if t.Type == budget.Expense {
			typeStr = "📉 Exp"
		}

		amountStr := fmt.Sprintf("$%8.2f", t.Amount)
		if t.Type == budget.Expense {
			amountStr = tui.GetNegativeStyle().Render("-" + fmt.Sprintf("$%7.2f", t.Amount))
		} else {
			amountStr = tui.GetPositiveStyle().Render(amountStr)
		}

		// Truncate description if too long
		description := t.Description
		if len(description) > 20 {
			description = description[:17] + "..."
		}

		sb.WriteString(fmt.Sprintf("%s %s %s %s %-20s %s\n",
			cursor,
			t.Date.Format("Jan 02"),
			typeStr,
			amountStr,
			description,
			t.Category))
	}

	return sb.String()
}

func (m model) viewAddTransactionForm() string {
	title := "Add Income"
	if m.state == addExpenseState {
		title = "Add Expense"
	}

	s := fmt.Sprintf("➕ %s\n\n", title)

	fields := []struct {
		label  string
		value  string
		active bool
	}{
		{"Amount", m.amountInput, m.activeField == 0},
		{"Description", m.descriptionInput, m.activeField == 1},
		{"Category", m.categoryInput, m.activeField == 2},
	}

	for _, field := range fields {
		prefix := " "
		if field.active {
			prefix = ">"
		}
		s += fmt.Sprintf("%s %s: %s\n", prefix, field.label, field.value)
	}

	s += "\nTab: switch fields • Enter: save • q/esc: return to dashboard"
	return s
}

func (m model) viewTransactions() string {
	s := "📝 Recent Transactions\n\n"

	if len(m.budget.Transactions) == 0 {
		s += "No transactions yet.\n"
	} else {
		for i, t := range m.budget.Transactions {
			cursor := " "
			if i == m.selectedTransaction {
				cursor = ">"
			}

			symbol := "📈"
			if t.Type == budget.Expense {
				symbol = "📉"
			}

			s += fmt.Sprintf("%s %s $%.2f - %s (%s)\n",
				cursor, symbol, t.Amount, t.Description, t.Category)
		}
	}

	s += "\n↑↓/j/k: navigate • q/esc: return to dashboard"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
