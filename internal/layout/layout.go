package layout

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Layout represents the filesystem structure of a trac repository.
type Layout struct {
	Root          string // Path to the .trac/ directory
	Objects       string // Path to the objects/ directory
	Refs          string // Path to the refs/ directory
	Heads         string // Path to the refs/heads/ directory
	HeadFile      string // Path to the HEAD file
	MainBranchRef string // Path to the main branch reference file (refs/heads/main)
	Index         string // Path to the index file (index.json)
}

// New creates a new Layout instance with paths initialized based on repoPath.
func New(repoPath string) (*Layout, error) {
	if repoPath == "" {
		return nil, fmt.Errorf("repoPath argument must be provided")
	}

	absPath, err := filepath.Abs(filepath.Join(repoPath, ".trac"))
	if err != nil {
		return nil, err
	}

	return &Layout{
		Root:          absPath,
		Objects:       filepath.Join(absPath, "objects"),
		Refs:          filepath.Join(absPath, "refs"),
		Heads:         filepath.Join(absPath, "refs", "heads"),
		HeadFile:      filepath.Join(absPath, "HEAD"),
		MainBranchRef: filepath.Join(absPath, "refs", "heads", "main"),
		Index:         filepath.Join(absPath, "index.json"),
	}, nil
}

// Init initializes an empty repository by creating necessary folders and files.
func (l *Layout) Init() error {
	directories := []string{
		l.Objects,
		l.Heads,
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	if err := os.WriteFile(l.HeadFile, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return fmt.Errorf("failed to write HEAD file: %w", err)
	}

	if err := os.WriteFile(l.MainBranchRef, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to write main branch reference: %w", err)
	}

	return nil
}

// Exists checks if the .trac directory already exists.
func (l *Layout) Exists() bool {
	_, err := os.Stat(l.Root)
	return !errors.Is(err, fs.ErrNotExist)
}

func (l *Layout) ValidateIsRepo() error {
	if !l.Exists() {
		return fmt.Errorf("not a trac repository (or any of the parent directories): .trac")
	}
	return nil
}
