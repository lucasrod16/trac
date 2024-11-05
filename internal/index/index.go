package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasrod16/trac/internal/layout"
	"github.com/lucasrod16/trac/internal/utils"
)

// Index represents the entire Index of staged and unstaged files.
type Index struct {
	Staged map[string]string `json:"staged"`
}

func New() *Index {
	return &Index{
		Staged: make(map[string]string),
	}
}

// Add adds an entry (file) to the Index by calculating its SHA-256 hash and storing it.
func (idx *Index) Add(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	hash, err := utils.HashFile(absPath)
	if err != nil {
		return err
	}
	idx.Staged[hash] = filePath
	return nil
}

// Write serializes the Index to a JSON file.
func (idx *Index) Write(l *layout.Layout) error {
	file, err := os.OpenFile(l.Index, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create or open Index file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(idx); err != nil {
		return fmt.Errorf("failed to write Index to file: %w", err)
	}

	return nil
}

// Load reads and deserializes the Index from a JSON file.
func (idx *Index) Load(l *layout.Layout) error {
	file, err := os.Open(l.Index)
	if err != nil {
		return fmt.Errorf("failed to open Index file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(idx); err != nil {
		return fmt.Errorf("failed to read Index from file: %w", err)
	}

	return nil
}