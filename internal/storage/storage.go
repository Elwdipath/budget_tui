package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Elwdipath/budget_tui/internal/budget"
)

func GenerateID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getDataFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".budget_tui.json")
}

func SaveBudget(b *budget.Budget) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getDataFilePath(), data, 0644)
}

func LoadBudget() (*budget.Budget, error) {
	filePath := getDataFilePath()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return budget.NewBudget(), nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var b budget.Budget
	err = json.Unmarshal(data, &b)
	if err != nil {
		return budget.NewBudget(), nil
	}

	return &b, nil
}
