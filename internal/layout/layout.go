package layout

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Layout represents the filesystem structure of a trac repository.
type Layout struct {
	Root          string // Path to the root of the repository (the directory containing .trac)
	Config        string // Path to the .trac/ directory (repository configuration)
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
		return nil, errors.New("repoPath argument must be provided")
	}
	rootPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(rootPath, ".trac")
	return &Layout{
		Root:          rootPath,
		Config:        configPath,
		Objects:       filepath.Join(configPath, "objects"),
		Refs:          filepath.Join(configPath, "refs"),
		Heads:         filepath.Join(configPath, "refs", "heads"),
		HeadFile:      filepath.Join(configPath, "HEAD"),
		MainBranchRef: filepath.Join(configPath, "refs", "heads", "main"),
		Index:         filepath.Join(configPath, "index.json"),
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
			return err
		}
	}
	if err := os.WriteFile(l.HeadFile, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(l.MainBranchRef, []byte{}, 0644); err != nil {
		return err
	}
	return nil
}

// Exists checks if the .trac directory already exists.
func (l *Layout) Exists() bool {
	_, err := os.Stat(l.Config)
	return !errors.Is(err, fs.ErrNotExist)
}

func (l *Layout) ValidateIsRepo() error {
	if !l.Exists() {
		return errors.New("not a trac repository (or any of the parent directories): .trac")
	}
	return nil
}

// ValidatePathInRepo validates if a given path is within the repository root.
func (l *Layout) ValidatePathInRepo(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absPath, l.Root) {
		return fmt.Errorf("%q is outside the repository at %q", absPath, l.Root)
	}
	return nil
}
