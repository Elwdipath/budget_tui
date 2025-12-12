package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type CategorizationRule struct {
	Pattern         string          `json:"pattern"`
	Category        string          `json:"category"`
	MinAmount       float64         `json:"min_amount,omitempty"`
	MaxAmount       float64         `json:"max_amount,omitempty"`
	Keywords        []string        `json:"keywords,omitempty"`
	Priority        int             `json:"priority"`
	IsActive        bool            `json:"is_active"`
	TransactionType TransactionType `json:"transaction_type,omitempty"`
}

type CategoryConfig struct {
	Rules []CategorizationRule `json:"rules"`
}

type Categorizer struct {
	rules []CategorizationRule
}

func NewCategorizer() *Categorizer {
	categorizer := &Categorizer{}
	categorizer.loadRules()
	return categorizer
}

func (c *Categorizer) loadRules() {
	// Load custom rules if they exist
	customRules, err := c.loadCustomRules()
	if err == nil {
		c.rules = append(c.rules, customRules...)
	}

	// Add default rules
	c.rules = append(c.rules, c.getDefaultRules()...)

	// Sort by priority (higher priority first)
	for i := 0; i < len(c.rules)-1; i++ {
		for j := i + 1; j < len(c.rules); j++ {
			if c.rules[j].Priority > c.rules[i].Priority {
				c.rules[i], c.rules[j] = c.rules[j], c.rules[i]
			}
		}
	}
}

func (c *Categorizer) loadCustomRules() ([]CategorizationRule, error) {
	filePath := filepath.Join(os.Getenv("HOME"), ".budget_tui_rules.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []CategorizationRule{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config CategoryConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config.Rules, nil
}

func (c *Categorizer) getDefaultRules() []CategorizationRule {
	return []CategorizationRule{
		// High priority specific merchants
		{
			Pattern:  ".*Netflix.*",
			Category: "Entertainment",
			Priority: 100,
			IsActive: true,
			Keywords: []string{"netflix"},
		},
		{
			Pattern:  ".*Spotify.*",
			Category: "Entertainment",
			Priority: 100,
			IsActive: true,
			Keywords: []string{"spotify"},
		},
		{
			Pattern:  ".*Amazon.*",
			Category: "Shopping",
			Priority: 95,
			IsActive: true,
			Keywords: []string{"amazon"},
		},
		{
			Pattern:  ".*Walmart.*",
			Category: "Shopping",
			Priority: 95,
			IsActive: true,
			Keywords: []string{"walmart"},
		},
		{
			Pattern:  ".*Target.*",
			Category: "Shopping",
			Priority: 95,
			IsActive: true,
			Keywords: []string{"target"},
		},

		// Food & Dining
		{
			Pattern:  ".*McDonald.*",
			Category: "Food & Dining",
			Priority: 90,
			IsActive: true,
			Keywords: []string{"mcdonald"},
		},
		{
			Pattern:  ".*Starbucks.*",
			Category: "Food & Dining",
			Priority: 90,
			IsActive: true,
			Keywords: []string{"starbucks"},
		},
		{
			Pattern:  ".*Restaurant|Dining|Cafe.*",
			Category: "Food & Dining",
			Priority: 80,
			IsActive: true,
			Keywords: []string{"restaurant", "dining", "cafe", "bistro"},
		},
		{
			Pattern:  ".*Grocery|Supermarket.*",
			Category: "Groceries",
			Priority: 85,
			IsActive: true,
			Keywords: []string{"grocery", "supermarket", "kroger", "safeway", "whole foods"},
		},

		// Housing & Utilities
		{
			Pattern:   ".*Rent|Mortgage.*",
			Category:  "Housing",
			Priority:  100,
			IsActive:  true,
			Keywords:  []string{"rent", "mortgage"},
			MinAmount: 500,
		},
		{
			Pattern:  ".*Electric|Gas|Water|Utility.*",
			Category: "Utilities",
			Priority: 90,
			IsActive: true,
			Keywords: []string{"electric", "gas", "water", "utility", "power"},
		},
		{
			Pattern:  ".*Internet|Cable|Phone.*",
			Category: "Utilities",
			Priority: 85,
			IsActive: true,
			Keywords: []string{"internet", "cable", "phone", "verizon", "comcast", "at&t"},
		},

		// Transportation
		{
			Pattern:  ".*Gas Station|Shell|Exxon|Chevron.*",
			Category: "Transportation",
			Priority: 90,
			IsActive: true,
			Keywords: []string{"gas station", "shell", "exxon", "chevron", "bp"},
		},
		{
			Pattern:  ".*Uber|Lyft|Taxi.*",
			Category: "Transportation",
			Priority: 90,
			IsActive: true,
			Keywords: []string{"uber", "lyft", "taxi", "rideshare"},
		},
		{
			Pattern:  ".*Parking|Toll.*",
			Category: "Transportation",
			Priority: 85,
			IsActive: true,
			Keywords: []string{"parking", "toll"},
		},

		// Healthcare
		{
			Pattern:  ".*CVS|Walgreens|Pharmacy.*",
			Category: "Healthcare",
			Priority: 90,
			IsActive: true,
			Keywords: []string{"cvs", "walgreens", "pharmacy"},
		},
		{
			Pattern:  ".*Hospital|Doctor|Medical.*",
			Category: "Healthcare",
			Priority: 85,
			IsActive: true,
			Keywords: []string{"hospital", "doctor", "medical", "clinic"},
		},

		// Financial
		{
			Pattern:  ".*ATM.*",
			Category: "Cash & ATM",
			Priority: 90,
			IsActive: true,
			Keywords: []string{"atm"},
		},
		{
			Pattern:  ".*Bank Fee|Interest Charge.*",
			Category: "Bank Fees",
			Priority: 85,
			IsActive: true,
			Keywords: []string{"bank fee", "interest charge", "overdraft", "maintenance"},
		},

		// Income specific rules
		{
			Pattern:         ".*Salary|Payroll|Paycheck.*",
			Category:        "Salary",
			Priority:        100,
			IsActive:        true,
			TransactionType: Income,
			Keywords:        []string{"salary", "payroll", "paycheck"},
		},
		{
			Pattern:         ".*Deposit.*",
			Category:        "Deposits",
			Priority:        80,
			IsActive:        true,
			TransactionType: Income,
			Keywords:        []string{"deposit"},
		},

		// General rules (lower priority)
		{
			Pattern:  ".*Transfer.*",
			Category: "Transfers",
			Priority: 60,
			IsActive: true,
			Keywords: []string{"transfer", "xfer"},
		},
		{
			Pattern:  ".*Payment.*",
			Category: "Bills & Payments",
			Priority: 70,
			IsActive: true,
			Keywords: []string{"payment"},
		},
	}
}

func (c *Categorizer) CategorizeTransaction(description string, amount float64, transType TransactionType) (string, float64) {
	description = strings.ToLower(strings.TrimSpace(description))
	bestMatch := "Uncategorized"
	bestConfidence := 0.0

	for _, rule := range c.rules {
		if !rule.IsActive {
			continue
		}

		// Check transaction type if specified
		if rule.TransactionType != "" && rule.TransactionType != transType {
			continue
		}

		// Check amount constraints
		if rule.MinAmount > 0 && amount < rule.MinAmount {
			continue
		}
		if rule.MaxAmount > 0 && amount > rule.MaxAmount {
			continue
		}

		var confidence float64
		var matched bool

		// Try regex pattern first
		if rule.Pattern != "" {
			matched, _ = regexp.MatchString(strings.ToLower(rule.Pattern), description)
			if matched {
				confidence = 0.9
			}
		}

		// Try keywords if pattern didn't match
		if !matched && len(rule.Keywords) > 0 {
			for _, keyword := range rule.Keywords {
				if strings.Contains(description, strings.ToLower(keyword)) {
					matched = true
					confidence = 0.8
					break
				}
			}
		}

		if matched {
			// Adjust confidence based on rule priority
			confidence = confidence * (float64(rule.Priority) / 100.0)

			// Add confidence for exact matches
			if strings.Contains(description, "amazon") && rule.Category == "Shopping" {
				confidence = math.Min(confidence+0.1, 1.0)
			}

			if confidence > bestConfidence {
				bestMatch = rule.Category
				bestConfidence = confidence
			}
		}
	}

	return bestMatch, math.Min(bestConfidence, 1.0)
}

func (c *Categorizer) AddCustomRule(rule CategorizationRule) error {
	c.rules = append(c.rules, rule)
	return c.saveCustomRules()
}

func (c *Categorizer) saveCustomRules() error {
	filePath := filepath.Join(os.Getenv("HOME"), ".budget_tui_rules.json")

	// Separate custom rules (non-default ones)
	var customRules []CategorizationRule
	defaultRules := c.getDefaultRules()

	for _, rule := range c.rules {
		isDefault := false
		for _, defaultRule := range defaultRules {
			if rule.Pattern == defaultRule.Pattern && rule.Category == defaultRule.Category {
				isDefault = true
				break
			}
		}
		if !isDefault {
			customRules = append(customRules, rule)
		}
	}

	config := CategoryConfig{Rules: customRules}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func (c *Categorizer) GetAllCategories() []string {
	categories := make(map[string]bool)
	categories["Uncategorized"] = true

	for _, rule := range c.rules {
		categories[rule.Category] = true
	}

	var result []string
	for category := range categories {
		result = append(result, category)
	}

	return result
}
