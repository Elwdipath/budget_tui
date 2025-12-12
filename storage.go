package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

func generateID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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

	var budget Budget
	err = json.Unmarshal(data, &budget)
	if err != nil {
		return NewBudget(), nil
	}

	return &budget, nil
}
