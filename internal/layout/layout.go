package layout

import (
	"fmt"
	"os"
	"path/filepath"
)

// Layout represents the filesystem structure of a trac repository.
type Layout struct {
	RootPath      string // Path to the .trac/ directory
	ObjectsPath   string // Path to the objects/ directory
	RefsPath      string // Path to the refs/ directory
	HeadsPath     string // Path to the refs/heads/ directory
	HeadFilePath  string // Path to the HEAD file
	MainBranchRef string // Path to the main branch reference file (refs/heads/main)
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
		RootPath:      absPath,
		ObjectsPath:   filepath.Join(absPath, "objects"),
		RefsPath:      filepath.Join(absPath, "refs"),
		HeadsPath:     filepath.Join(absPath, "refs", "heads"),
		HeadFilePath:  filepath.Join(absPath, "HEAD"),
		MainBranchRef: filepath.Join(absPath, "refs", "heads", "main"),
	}, nil
}

// Init initializes an empty repository by creating necessary folders and files.
func (l *Layout) Init() error {
	directories := []string{
		l.ObjectsPath,
		l.HeadsPath,
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	if err := os.WriteFile(l.HeadFilePath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return fmt.Errorf("failed to write HEAD file: %w", err)
	}

	if err := os.WriteFile(l.MainBranchRef, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to write main branch reference: %w", err)
	}

	return nil
}

// Exists checks if the .trac directory already exists.
func (l *Layout) Exists() bool {
	_, err := os.Stat(l.RootPath)
	return !os.IsNotExist(err)
}

func (l *Layout) ValidateIsRepo() error {
	if !l.Exists() {
		return fmt.Errorf("not a trac repository (or any of the parent directories): .trac")
	}
	return nil
}
