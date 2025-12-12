package main

import (
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

func (b *Budget) AddTransaction(amount float64, description, category string, tType TransactionType) {
	transaction := Transaction{
		ID:          generateID(),
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
