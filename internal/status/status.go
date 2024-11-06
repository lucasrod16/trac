package status

import (
	"os"
	"path/filepath"

	"github.com/lucasrod16/trac/internal/index"
	"github.com/lucasrod16/trac/internal/layout"
	"github.com/lucasrod16/trac/internal/utils"
)

// repoStatus holds the state of the repository, including tracked and untracked files.
type repoStatus struct {
	layout    *layout.Layout
	index     *index.Index
	tracked   []string
	untracked []string
}

func NewRepoStatus(l *layout.Layout, idx *index.Index) *repoStatus {
	return &repoStatus{
		layout:    l,
		index:     idx,
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
			if info.Name() == ".trac" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		isStaged, err := rs.isStaged(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}

		if isStaged {
			rs.addTracked(relPath)
		} else {
			rs.addUntracked(relPath)
		}
		return nil
	})
	return err
}

func (rs *repoStatus) isStaged(path string) (bool, error) {
	hash, err := utils.HashFile(path)
	if err != nil {
		return false, err
	}
	for _, existingHash := range rs.index.Staged {
		if hash == existingHash {
			return true, nil
		}
	}
	return false, nil
}
