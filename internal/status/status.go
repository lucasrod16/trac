package status

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// repoStatus holds the state of the repository, including tracked and untracked files.
type repoStatus struct {
	tracked   []string
	untracked []string
}

func NewRepoStatus() *repoStatus {
	return &repoStatus{
		tracked:   []string{},
		untracked: []string{},
	}
}

func (rs *repoStatus) addTracked(file string) {
	if !slices.Contains(rs.tracked, file) {
		rs.tracked = append(rs.tracked, file)
	}
}

func (rs *repoStatus) addUntracked(file string) {
	if !slices.Contains(rs.untracked, file) {
		rs.untracked = append(rs.untracked, file)
	}
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
			if strings.Contains(path, ".trac") || strings.Contains(path, ".git") {
				return filepath.SkipDir
			}
			return nil
		}

		isTracked, err := isFileTracked(cwd, path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}

		if isTracked {
			rs.addTracked(relPath)
		}

		if !isTracked {
			rs.addUntracked(relPath)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// getStagedFiles finds staged files in the repository.
func (rs *repoStatus) GetStagedFiles() error {
	// TODO: implement logic to find staged files.
	// Will require implementing staging with `trac add` first.
	// Lookup from index file.
	return nil
}

// isFileTracked checks if a file's content hash exists in .trac/objects/.
func isFileTracked(repoPath, filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(content)
	hashString := fmt.Sprintf("%x", hash)

	objectPath := filepath.Join(repoPath, ".trac", "objects", hashString[:2], hashString[2:])
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		return false, nil // file is untracked
	}
	return true, nil // file is tracked
}
