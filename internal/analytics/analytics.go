package analytics

import (
	"sort"
	"time"

	"github.com/Elwdipath/budget_tui/internal/budget"
)

type CategorySpending struct {
	Category string
	Amount   float64
	Count    int
}

func GetSpendingByCategory(b *budget.Budget) []CategorySpending {
	categoryTotals := make(map[string]CategorySpending)

	for _, t := range b.Transactions {
		if t.Type == budget.Expense {
			if existing, ok := categoryTotals[t.Category]; ok {
				existing.Amount += t.Amount
				existing.Count++
				categoryTotals[t.Category] = existing
			} else {
				categoryTotals[t.Category] = CategorySpending{
					Category: t.Category,
					Amount:   t.Amount,
					Count:    1,
				}
			}
		}
	}

	var result []CategorySpending
	for _, cs := range categoryTotals {
		result = append(result, cs)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Amount > result[j].Amount
	})

	return result
}

func GetRecentTransactions(b *budget.Budget, limit int) []budget.Transaction {
	if len(b.Transactions) <= limit {
		return b.Transactions
	}

	transactions := make([]budget.Transaction, len(b.Transactions))
	copy(transactions, b.Transactions)

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date.After(transactions[j].Date)
	})

	return transactions[:limit]
}

func GetThisMonthTotals(b *budget.Budget) (income, expenses float64) {
	now := time.Now()
	currentMonth := now.Month()
	currentYear := now.Year()

	for _, t := range b.Transactions {
		if t.Date.Month() == currentMonth && t.Date.Year() == currentYear {
			if t.Type == budget.Income {
				income += t.Amount
			} else {
				expenses += t.Amount
			}
		}
	}
	return income, expenses
}

func GetFinancialHealthStatus(b *budget.Budget) string {
	balance := b.GetBalance()
	if balance > 0 {
		return "✅ Financially Healthy"
	} else if balance > -100 {
		return "⚠️  Watch Your Spending"
	} else {
		return "🚨 Budget Alert"
	}
}
