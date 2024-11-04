package status

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasrod16/trac/internal/layout"
)

// repoStatus holds the state of the repository, including tracked and untracked files.
type repoStatus struct {
	layout    *layout.Layout
	tracked   []string
	untracked []string
}

func NewRepoStatus(l *layout.Layout) *repoStatus {
	return &repoStatus{
		layout:    l,
		tracked:   []string{},
		untracked: []string{},
	}
}

func (rs *repoStatus) addTracked(path string) {
	rs.tracked = append(rs.tracked, path)
}

func (rs *repoStatus) addUntracked(path string) {
	rs.untracked = append(rs.untracked, path)
}

func (rs *repoStatus) GetTracked() []string {
	return rs.tracked
}

func (rs *repoStatus) GetUntracked() []string {
	return rs.untracked
}

func (rs *repoStatus) HasTracked() bool {
	return len(rs.tracked) > 0
}

func (rs *repoStatus) HasUntracked() bool {
	return len(rs.untracked) > 0
}

// DetectTrackedStatus scans the current working directory for tracked and untracked files.
func (rs *repoStatus) DetectTrackedStatus() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.Contains(path, rs.layout.RootPath) || strings.Contains(path, ".git") {
				return filepath.SkipDir
			}
			return nil
		}

		isTracked, err := rs.isFileTracked(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}

		if isTracked {
			rs.addTracked(relPath)
		} else {
			rs.addUntracked(relPath)
		}
		return nil
	})
	return err
}

// GetStagedFiles finds staged files in the repository.
func (rs *repoStatus) GetStagedFiles() error {
	// TODO: implement logic to find staged files.
	// Will require implementing staging with `trac add` first.
	// Lookup from index file.
	return nil
}

// isFileTracked checks if a file's content hash exists in the .trac/objects directory.
func (rs *repoStatus) isFileTracked(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(content)
	hashString := fmt.Sprintf("%x", hash)

	objectPath := filepath.Join(rs.layout.ObjectsPath, hashString[:2], hashString[2:])
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		return false, nil // file is untracked
	}
	return true, nil // file is tracked
}
