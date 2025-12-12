package import

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CSVFormat struct {
	Name              string
	DateColumn        int
	DescriptionColumn int
	AmountColumn      int
	DateFormat        string
	AmountIsNegative  bool
	HasHeader         bool
	Delimiter         rune
}

// Common CSV formats for different banks
var CommonFormats = []CSVFormat{
	{
		Name:              "Chase",
		DateColumn:        0,
		DescriptionColumn: 2,
		AmountColumn:      3,
		DateFormat:        "01/02/2006",
		AmountIsNegative:  true,
		HasHeader:         true,
		Delimiter:         ',',
	},
	{
		Name:              "Bank of America",
		DateColumn:        0,
		DescriptionColumn: 1,
		AmountColumn:      2,
		DateFormat:        "01/02/2006",
		AmountIsNegative:  false,
		HasHeader:         true,
		Delimiter:         ',',
	},
	{
		Name:              "Wells Fargo",
		DateColumn:        1,
		DescriptionColumn: 4,
		AmountColumn:      2,
		DateFormat:        "01/02/06",
		AmountIsNegative:  true,
		HasHeader:         true,
		Delimiter:         ',',
	},
	{
		Name:              "Generic",
		DateColumn:        0,
		DescriptionColumn: 1,
		AmountColumn:      2,
		DateFormat:        "2006-01-02",
		AmountIsNegative:  true,
		HasHeader:         false,
		Delimiter:         ',',
	},
}

type ImportResult struct {
	Transactions []Transaction `json:"transactions"`
	Format       CSVFormat     `json:"format"`
	Errors       []string      `json:"errors"`
	TotalRows    int           `json:"total_rows"`
	SuccessCount int           `json:"success_count"`
}

func DetectCSVFormat(filePath string) (*CSVFormat, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Try each format
	for _, format := range CommonFormats {
		file.Seek(0, 0)
		reader.Comma = format.Delimiter

		records, err := reader.ReadAll()
		if err != nil {
			continue
		}

		if len(records) < 2 {
			continue
		}

		// Test parsing a few rows
		successCount := 0
		testRows := records
		if format.HasHeader && len(records) > 1 {
			testRows = records[1:]
		}

		for i := 0; i < len(testRows) && i < 5; i++ {
			row := testRows[i]
			if len(row) <= max(format.DateColumn, max(format.DescriptionColumn, format.AmountColumn)) {
				continue
			}

			// Test date parsing
			_, err := time.Parse(format.DateFormat, strings.TrimSpace(row[format.DateColumn]))
			if err != nil {
				continue
			}

			// Test amount parsing
			amountStr := strings.TrimSpace(row[format.AmountColumn])
			_, err = strconv.ParseFloat(amountStr, 64)
			if err != nil {
				// Try removing common formatting
				re := regexp.MustCompile(`[$,]`)
				amountStr = re.ReplaceAllString(amountStr, "")
				_, err = strconv.ParseFloat(amountStr, 64)
				if err != nil {
					continue
				}
			}

			successCount++
		}

		if successCount >= 3 {
			return &format, nil
		}
	}

	// Return generic format as fallback
	return &CommonFormats[len(CommonFormats)-1], nil
}

func ParseCSV(filePath string, format *CSVFormat) (*ImportResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = format.Delimiter

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %v", err)
	}

	result := &ImportResult{
		Transactions: []Transaction{},
		Format:       *format,
		Errors:       []string{},
		TotalRows:    len(records),
	}

	startRow := 0
	if format.HasHeader && len(records) > 0 {
		startRow = 1
	}

	for i := startRow; i < len(records); i++ {
		row := records[i]

		if len(row) <= max(format.DateColumn, max(format.DescriptionColumn, format.AmountColumn)) {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: insufficient columns", i+1))
			continue
		}

		// Parse date
		dateStr := strings.TrimSpace(row[format.DateColumn])
		date, err := time.Parse(format.DateFormat, dateStr)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: invalid date '%s'", i+1, dateStr))
			continue
		}

		// Parse amount
		amountStr := strings.TrimSpace(row[format.AmountColumn])
		re := regexp.MustCompile(`[$,() ]`)
		amountStr = re.ReplaceAllString(amountStr, "")

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: invalid amount '%s'", i+1, row[format.AmountColumn]))
			continue
		}

		// Determine transaction type
		var transType TransactionType
		description := strings.TrimSpace(row[format.DescriptionColumn])

		if format.AmountIsNegative {
			if amount < 0 {
				transType = Expense
				amount = -amount
			} else {
				transType = Income
			}
		} else {
			// For formats where expenses are positive but we need to determine type from description
			if isIncomeDescription(description) {
				transType = Income
			} else {
				transType = Expense
			}
		}

		transaction := Transaction{
			ID:                  generateID(),
			Amount:              amount,
			Description:         description,
			OriginalDescription: description,
			Category:            "Uncategorized",
			Type:                transType,
			Date:                date,
			ImportSource:        format.Name,
			IsImported:          true,
		}

		result.Transactions = append(result.Transactions, transaction)
		result.SuccessCount++
	}

	return result, nil
}

func isIncomeDescription(description string) bool {
	description = strings.ToLower(description)
	incomeKeywords := []string{
		"deposit", "salary", "payroll", "income", "payment", "credit",
		"refund", "transfer in", "direct deposit", "interest", "dividend",
		"bonus", "commission", "cash back",
	}

	for _, keyword := range incomeKeywords {
		if strings.Contains(description, keyword) {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetImportPreview(filePath string, format *CSVFormat, maxRows int) ([]PreviewTransaction, error) {
	result, err := ParseCSV(filePath, format)
	if err != nil {
		return nil, err
	}

	preview := []PreviewTransaction{}
	count := maxRows
	if len(result.Transactions) < count {
		count = len(result.Transactions)
	}

	categorizer := NewCategorizer()

	for i := 0; i < count; i++ {
		t := result.Transactions[i]
		category, confidence := categorizer.CategorizeTransaction(t.Description, t.Amount, t.Type)

		preview = append(preview, PreviewTransaction{
			Amount:      t.Amount,
			Description: t.Description,
			Date:        t.Date.Format("Jan 02"),
			Category:    category,
			Confidence:  confidence,
		})
	}

	return preview, nil
}
