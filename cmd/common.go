package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func validateIsRepo() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		tracPath := filepath.Join(cwd, ".trac")
		if _, err := os.Stat(tracPath); err == nil {
			return nil
		}

		parentDir := filepath.Dir(cwd)
		if parentDir == cwd {
			break
		}
		cwd = parentDir
	}

	return fmt.Errorf("not a trac repository (or any of the parent directories): .trac")
}
