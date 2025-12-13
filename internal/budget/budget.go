package budget

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          string          `json:"id"`
	Amount      float64         `json:"amount"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Type        TransactionType `json:"type"`
	Date        time.Time       `json:"date"`
	// Import-specific fields
	OriginalDescription string  `json:"original_description,omitempty"`
	ImportSource        string  `json:"import_source,omitempty"`
	Confidence          float64 `json:"confidence,omitempty"`
	IsImported          bool    `json:"is_imported,omitempty"`
}

type Budget struct {
	Transactions []Transaction `json:"transactions"`
}

func NewBudget() *Budget {
	return &Budget{
		Transactions: []Transaction{},
	}
}

func GenerateID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (b *Budget) AddTransaction(amount float64, description, category string, tType TransactionType) {
	transaction := Transaction{
		ID:          GenerateID(),
		Amount:      amount,
		Description: description,
		Category:    category,
		Type:        tType,
		Date:        time.Now(),
	}
	b.Transactions = append(b.Transactions, transaction)
}

func (b *Budget) GetTotalIncome() float64 {
	total := 0.0
	for _, t := range b.Transactions {
		if t.Type == Income {
			total += t.Amount
		}
	}
	return total
}

func (b *Budget) GetTotalExpenses() float64 {
	total := 0.0
	for _, t := range b.Transactions {
		if t.Type == Expense {
			total += t.Amount
		}
	}
	return total
}

func (b *Budget) GetBalance() float64 {
	return b.GetTotalIncome() - b.GetTotalExpenses()
}

func (b *Budget) GetTransactionsByCategory(category string) []Transaction {
	var transactions []Transaction
	for _, t := range b.Transactions {
		if t.Category == category {
			transactions = append(transactions, t)
		}
	}
	return transactions
}

type CategorySpending struct {
	Category string
	Amount   float64
	Count    int
}

func (b *Budget) GetSpendingByCategory() []CategorySpending {
	categoryMap := make(map[string]*CategorySpending)
	for _, t := range b.Transactions {
		if t.Type == Expense {
			if _, exists := categoryMap[t.Category]; !exists {
				categoryMap[t.Category] = &CategorySpending{Category: t.Category, Amount: 0, Count: 0}
			}
			categoryMap[t.Category].Amount += t.Amount
			categoryMap[t.Category].Count++
		}
	}
	var categories []CategorySpending
	for _, cat := range categoryMap {
		categories = append(categories, *cat)
	}
	// Sort by amount descending
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Amount > categories[j].Amount
	})
	return categories
}

func (b *Budget) GetRecentTransactions(limit int) []Transaction {
	transactions := make([]Transaction, len(b.Transactions))
	copy(transactions, b.Transactions)
	// Sort by date descending
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date.After(transactions[j].Date)
	})
	if len(transactions) > limit {
		return transactions[:limit]
	}
	return transactions
}

func (b *Budget) GetFinancialHealthStatus() string {
	balance := b.GetBalance()
	if balance > 0 {
		return "✅ Positive balance"
	} else if balance == 0 {
		return "⚠️ Balanced"
	} else {
		return "❌ Negative balance"
	}
}

func getDataFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".budget_tui.json")
}

func (b *Budget) Save() error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getDataFilePath(), data, 0644)
}

func LoadBudget() (*Budget, error) {
	filePath := getDataFilePath()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return NewBudget(), nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var b Budget
	err = json.Unmarshal(data, &b)
	if err != nil {
		return NewBudget(), nil
	}

	return &b, nil
}
