package index

import (
	"encoding/json"
	"os"

	"github.com/lucasrod16/trac/internal/layout"
	"github.com/lucasrod16/trac/internal/utils"
)

// Index represents the index of staged files.
type Index struct {
	Staged map[string]string `json:"staged"`
}

func New() *Index {
	return &Index{
		Staged: make(map[string]string),
	}
}

// Add adds an entry (file) to the index by calculating its SHA-256 hash and storing it.
func (idx *Index) Add(filePath string, l *layout.Layout) error {
	if err := l.ValidatePathInRepo(filePath); err != nil {
		return err
	}
	hash, err := utils.HashFile(filePath)
	if err != nil {
		return err
	}
	idx.Staged[filePath] = hash
	return nil
}

// Write serializes the index to a JSON file.
func (idx *Index) Write(l *layout.Layout) error {
	file, err := os.OpenFile(l.Index, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(idx); err != nil {
		return err
	}
	return nil
}

// Load reads and deserializes the index from a JSON file.
func (idx *Index) Load(l *layout.Layout) error {
	file, err := os.Open(l.Index)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(idx); err != nil {
		return err
	}
	return nil
}
