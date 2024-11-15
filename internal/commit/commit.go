package commit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasrod16/trac/internal/layout"
)

// Commit represents a commit in the repository.
type Commit struct {
	Parent    string            `json:"parent"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Changes   map[string]string `json:"changes"`
}

func New(message, parent string, stagedFiles map[string]string) *Commit {
	return &Commit{
		Parent:    parent,
		Message:   message,
		Timestamp: time.Now(),
		Changes:   stagedFiles,
	}
}

// Save writes the commit object to the repository and updates the main branch reference.
func (c *Commit) Save(l *layout.Layout) (string, error) {
	commitData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(commitData)
	commitHash := hex.EncodeToString(hash[:])

	objectDir := filepath.Join(l.Objects, commitHash[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return "", err
	}
	commitPath := filepath.Join(objectDir, commitHash[2:])
	if err := os.WriteFile(commitPath, commitData, 0644); err != nil {
		return "", err
	}
	if err := os.WriteFile(l.HeadFile, []byte(commitHash+"\n"), 0644); err != nil {
		return "", err
	}
	for filePath, contentHash := range c.Changes {
		if err := copyFileObject(contentHash, filePath, l); err != nil {
			return "", err
		}
	}
	return commitHash, nil
}

func copyFileObject(contentHash string, src string, l *layout.Layout) error {
	objectDir := filepath.Join(l.Objects, contentHash[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return err
	}
	objectPath := filepath.Join(objectDir, contentHash[2:])

	if _, err := os.Stat(objectPath); errors.Is(err, fs.ErrNotExist) {
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.OpenFile(objectPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}
	}
	return nil
}

// LoadParent loads the hash of the latest commit from HEAD.
func LoadParent(l *layout.Layout) (string, error) {
	data, err := os.ReadFile(l.HeadFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
