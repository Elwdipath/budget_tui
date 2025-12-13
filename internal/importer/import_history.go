package importer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ImportSession struct {
	ID         string               `json:"id"`
	FileName   string               `json:"file_name"`
	Source     string               `json:"source"`
	Status     string               `json:"status"` // "pending", "reviewing", "imported", "error"
	TotalCount int                  `json:"total_count"`
	Imported   int                  `json:"imported"`
	Skipped    int                  `json:"skipped"`
	Errors     []string             `json:"errors,omitempty"`
	Preview    []PreviewTransaction `json:"preview,omitempty"`
	Timestamp  string               `json:"timestamp"`
}

type PreviewTransaction struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Category    string  `json:"category"`
	Confidence  float64 `json:"confidence"`
}

type ImportHistory struct {
	Sessions []ImportSession `json:"sessions"`
}

func getImportHistoryPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".budget_tui_imports.json")
}

func LoadImportHistory() (*ImportHistory, error) {
	filePath := getImportHistoryPath()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &ImportHistory{Sessions: []ImportSession{}}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return &ImportHistory{Sessions: []ImportSession{}}, nil
	}

	var history ImportHistory
	err = json.Unmarshal(data, &history)
	if err != nil {
		return &ImportHistory{Sessions: []ImportSession{}}, nil
	}

	return &history, nil
}

func (h *ImportHistory) Save() error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getImportHistoryPath(), data, 0644)
}

func (h *ImportHistory) AddSession(session ImportSession) {
	h.Sessions = append(h.Sessions, session)

	// Keep only last 10 sessions
	if len(h.Sessions) > 10 {
		h.Sessions = h.Sessions[len(h.Sessions)-10:]
	}
}

func (h *ImportHistory) GetLastSession() *ImportSession {
	if len(h.Sessions) == 0 {
		return nil
	}
	return &h.Sessions[len(h.Sessions)-1]
}
