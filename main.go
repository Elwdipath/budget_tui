package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Elwdipath/budget_tui/internal/analytics"
	"github.com/Elwdipath/budget_tui/internal/budget"
	"github.com/Elwdipath/budget_tui/internal/import"
	"github.com/Elwdipath/budget_tui/internal/tui"
	"github.com/Elwdipath/budget_tui/pkg/categorizer"
)

type state int

const (
	dashboardState state = iota
	importState
	reviewState
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

	// Import state
	importFilePath     string
	importFormat       *import.CSVFormat
	importResult       *import.ImportResult
	importSession      *import.ImportSession
	importHistory      *importhistory.ImportHistory
	categorizer        *categorizer.Categorizer
	selectedPreview    int
	showImportDetails  bool
	importStatus       string
}

func initialModel() model {
	b, _ := budget.LoadBudget()
	importHistory, _ := importhistory.LoadImportHistory()
	return model{
		state:  dashboardState,
		budget: b,
		menuChoices: []string{
			"[i] Income  [e] Expense  [t] Transactions  [b] Import  [h] Help  [q] Quit",
		},
		menuCursor:          0,
		activeField:         0,
		dashboardCursor:     0,
		selectedCategory:    0,
		selectedTransaction: 0,
		showHelp:            false,
		importHistory:       importHistory,
		categorizer:         categorizer.NewCategorizer(),
		selectedPreview:     0,
		showImportDetails:   false,
		importStatus:        "ready",
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
		case importState:
			return m.updateImportState(msg)
		case reviewState:
			return m.updateReviewState(msg)
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
	case "b":
		m.state = importState
		m.resetImportState()
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
	case "left", "right":
		// Switch between dashboard panels
		if msg.String() == "left" && m.selectedTransaction > 0 {
			m.selectedTransaction--
		} else if msg.String() == "right" {
			m.selectedTransaction++
		}
	case "esc":
		// Return to dashboard from other views will be handled by their own update functions
		return m, nil
	}
	return m, nil
}

func (m model) resetForm() {
	m.amountInput = ""
	m.descriptionInput = ""
	m.categoryInput = ""
	m.activeField = 0
	m.formSubmitted = false
}

func (m model) resetImportState() {
	m.importFilePath = ""
	m.importFormat = nil
	m.importResult = nil
	m.importSession = nil
	m.selectedPreview = 0
	m.showImportDetails = false
	m.importStatus = "ready"
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

			tType := Income
			if m.state == addExpenseState {
				tType = Expense
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

func (m model) updateImportState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = dashboardState
		m.resetImportState()
	case "enter":
		if m.importFilePath != "" {
			// Start import process
			m.importStatus = "detecting"
			m.importSession = &ImportSession{
				ID:         generateID(),
				FileName:   m.importFilePath,
				Source:     "Unknown",
				Status:     "reviewing",
				TotalCount: 0,
				Imported:   0,
				Skipped:    0,
				Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
			}

			// Detect format
			format, err := DetectCSVFormat(m.importFilePath)
			if err != nil {
				m.importStatus = "error: " + err.Error()
				return m, nil
			}

			m.importFormat = format
			m.importSession.Source = format.Name

			// Parse and preview
			result, err := ParseCSV(m.importFilePath, format)
			if err != nil {
				m.importStatus = "error: " + err.Error()
				return m, nil
			}

			m.importResult = result
			m.importSession.TotalCount = len(result.Transactions)

			// Generate preview
			preview, _ := GetImportPreview(m.importFilePath, format, 10)
			m.importSession.Preview = preview

			m.state = reviewState
			m.importStatus = "ready for review"
		}
	case "tab":
		// For file path input - this is a simplified version
		// In a real implementation, you'd want a file browser
		// For now, we'll use the sample file
		if m.importFilePath == "" {
			m.importFilePath = "sample_bank_statement.csv"
		}
	}
	return m, nil
}

func (m model) updateReviewState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = dashboardState
		m.resetImportState()
	case "c":
		// Confirm and complete import
		if m.importResult != nil {
			m.importStatus = "importing"

			// Add transactions to budget
			importedCount := 0
			for _, t := range m.importResult.Transactions {
				// Categorize transaction
				category, confidence := m.categorizer.CategorizeTransaction(t.Description, t.Amount, t.Type)
				t.Category = category
				t.Confidence = confidence

				m.budget.Transactions = append(m.budget.Transactions, t)
				importedCount++
			}

			// Save budget
			m.budget.Save()

			// Update import session
			if m.importSession != nil {
				m.importSession.Status = "imported"
				m.importSession.Imported = importedCount
				m.importSession.Skipped = len(m.importResult.Transactions) - importedCount
				m.importHistory.AddSession(*m.importSession)
				m.importHistory.Save()
			}

			m.importStatus = fmt.Sprintf("imported %d transactions", importedCount)
			m.state = dashboardState
			m.resetImportState()
		}
	case "up", "k":
		if m.selectedPreview > 0 {
			m.selectedPreview--
		}
	case "down", "j":
		if m.importSession != nil && m.selectedPreview < len(m.importSession.Preview)-1 {
			m.selectedPreview++
		}
	case "d":
		m.showImportDetails = !m.showImportDetails
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
	case importState:
		return m.viewImportState()
	case reviewState:
		return m.viewReviewState()
	case addIncomeState, addExpenseState:
		return m.viewAddTransactionForm()
	case viewTransactionsState:
		return m.viewTransactions()
	default:
		return ""
	}
}

func (m model) viewDashboard() string {
	// Hero Banner
	heroBanner := getHeroBanner()
	subtitle := getSubtitle()

	// Financial Summary Panel
	summaryContent := renderFinancialSummary(m.budget)
	summaryPanel := summaryPanelStyle.Render(summaryContent)

	// Category Spending Panel
	categories := m.budget.GetSpendingByCategory()
	var categoryContent string

	if len(categories) == 0 {
		categoryContent = "No expense categories yet.\n\nAdd some expenses to see breakdown."
	} else {
		var sb strings.Builder
		sb.WriteString("Top Spending Categories:\n\n")

		totalExpenses := m.budget.GetTotalExpenses()
		maxDisplay := 5
		if len(categories) < maxDisplay {
			maxDisplay = len(categories)
		}

		for i := 0; i < maxDisplay; i++ {
			cat := categories[i]
			bar := renderCategoryBar(cat.Category, cat.Amount, totalExpenses, 15)

			if i == m.dashboardCursor {
				bar = positiveStyle.Render("►") + bar[1:]
			}

			sb.WriteString(bar + "\n")
			sb.WriteString(fmt.Sprintf("           $%s (%d items)\n\n", formatAmount(cat.Amount), cat.Count))
		}

		if len(categories) > maxDisplay {
			sb.WriteString(fmt.Sprintf("... and %d more categories\n", len(categories)-maxDisplay))
		}

		categoryContent = sb.String()
	}
	categoryPanel := categoryPanelStyle.Render(categoryContent)

	// Import Status Panel (small)
	importContent := m.renderImportStatus()
	importPanel := borderStyle.Render(importContent)

	// Recent Transactions Panel
	recentTransactions := m.budget.GetRecentTransactions(10)
	transactionsContent := renderTransactionTable(recentTransactions, m.selectedTransaction)
	transactionsPanel := transactionsPanelStyle.Render(transactionsContent)

	// Layout panels
	// Top row: Summary panel
	topRow := lipgloss.JoinHorizontal(lipgloss.Left, summaryPanel)

	// Middle row: Categories and Import Status
	middleRow := lipgloss.JoinHorizontal(lipgloss.Left, categoryPanel, " ", importPanel)

	// Bottom row: Transactions
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Left, transactionsPanel)

	// Combine all rows
	dashboard := lipgloss.JoinVertical(lipgloss.Top, heroBanner, subtitle, topRow, "\n", middleRow, "\n", bottomRow)

	// Create menu bar at the bottom
	menuBar := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Foreground(blueColor).
		Render(strings.Join(m.menuChoices, " • "))

	// Add help text if requested
	var helpText string
	if m.showHelp {
		helpText = "\n" + helpStyle.Render("Controls: ↑↓/j/k: Navigate categories | Use shortcuts below menu | h: Toggle help")
	}

	// Combine dashboard and menu
	fullDashboard := lipgloss.JoinVertical(lipgloss.Top, dashboard, "\n", menuBar)

	return lipgloss.NewStyle().Padding(1, 0, 0, 0).Render(fullDashboard + helpText)
}

func (m model) viewImportState() string {
	title := titleStyle.Render("📁 Import Bank Statement")

	var content strings.Builder

	content.WriteString("Import CSV bank statements to automatically categorize transactions.\n\n")

	// File path input
	content.WriteString("File Path:\n")
	if m.importFilePath == "" {
		content.WriteString("  Press Tab to use sample file\n")
	} else {
		content.WriteString(fmt.Sprintf("  > %s\n", m.importFilePath))
	}

	// Status
	content.WriteString(fmt.Sprintf("\nStatus: %s\n", m.importStatus))

	// Instructions
	instructions := neutralStyle.Render(`
Instructions:
1. Press Tab to set a sample file path (or enter your own)
2. Press Enter to detect format and preview
3. Review the imported transactions
4. Confirm to add to your budget

Note: This is a demo version. In production, you'll be able to browse files.
`)

	content.WriteString(instructions)

	// Navigation
	nav := helpStyle.Render("Tab: Set file path • Enter: Import • q/esc: Back to dashboard")

	panel := borderStyle.Render(content.String())

	return lipgloss.JoinVertical(lipgloss.Top, title, panel, nav)
}

func (m model) viewReviewState() string {
	title := titleStyle.Render("🔍 Review Import")

	var content strings.Builder

	if m.importSession == nil || m.importResult == nil {
		content.WriteString("No import session to review.\n")
	} else {
		// Import summary
		content.WriteString(fmt.Sprintf("File: %s\n", m.importSession.FileName))
		content.WriteString(fmt.Sprintf("Format: %s\n", m.importSession.Source))
		content.WriteString(fmt.Sprintf("Total transactions: %d\n", len(m.importResult.Transactions)))
		content.WriteString(fmt.Sprintf("Parse errors: %d\n\n", len(m.importResult.Errors)))

		// Preview transactions
		content.WriteString("Preview (first 10 transactions):\n\n")

		for i, preview := range m.importSession.Preview {
			cursor := " "
			if i == m.selectedPreview {
				cursor = ">"
			}

			amountStr := fmt.Sprintf("$%.2f", preview.Amount)
			if preview.Amount < 0 {
				amountStr = negativeStyle.Render(amountStr)
			} else {
				amountStr = positiveStyle.Render(amountStr)
			}

			confidenceBar := getConfidenceBar(preview.Confidence)

			content.WriteString(fmt.Sprintf("%s %s %-20s %s %s\n",
				cursor,
				preview.Date,
				preview.Description[:min(20, len(preview.Description))],
				amountStr,
				confidenceBar))

			content.WriteString(fmt.Sprintf("  %s (%.0f%% confident)\n",
				preview.Category,
				preview.Confidence*100))
		}

		// Details toggle
		if m.showImportDetails {
			content.WriteString("\n--- Import Details ---\n")
			if len(m.importResult.Errors) > 0 {
				content.WriteString("Parse errors:\n")
				for _, err := range m.importResult.Errors {
					content.WriteString(fmt.Sprintf("  - %s\n", err))
				}
			} else {
				content.WriteString("No parse errors detected.\n")
			}
		}
	}

	// Instructions
	instructions := `
Navigation:
↑↓/j/k: Navigate transactions  
c: Confirm and import
d: Toggle details
q/esc: Cancel and return to dashboard
`

	content.WriteString(helpStyle.Render(instructions))

	panel := borderStyle.Render(content.String())

	return lipgloss.JoinVertical(lipgloss.Top, title, panel)
}

func getConfidenceBar(confidence float64) string {
	width := 10
	filled := int(confidence * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	var style lipgloss.Style
	if confidence >= 0.8 {
		style = positiveStyle
	} else if confidence >= 0.5 {
		style = neutralStyle
	} else {
		style = negativeStyle
	}

	return style.Render(bar)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m model) renderImportStatus() string {
	var sb strings.Builder

	sb.WriteString("📊 Import Status\n\n")

	if m.importHistory != nil && len(m.importHistory.Sessions) > 0 {
		lastSession := m.importHistory.GetLastSession()
		sb.WriteString(fmt.Sprintf("Last Import: %s\n", lastSession.Timestamp))
		sb.WriteString(fmt.Sprintf("File: %s\n", lastSession.FileName))
		sb.WriteString(fmt.Sprintf("Status: %s\n", lastSession.Status))
		sb.WriteString(fmt.Sprintf("Imported: %d transactions\n", lastSession.Imported))

		if lastSession.Status == "imported" {
			sb.WriteString("\n" + positiveStyle.Render("✅ Import completed successfully"))
		} else if lastSession.Status == "error" {
			sb.WriteString("\n" + negativeStyle.Render("❌ Last import had errors"))
		}
	} else {
		sb.WriteString("No imports yet\n")
		sb.WriteString("\nPress [b] to import your first bank statement")
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
			if t.Type == Expense {
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
