package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type initOptions struct {
	// repo path argument
	repoPath string
}

func NewInitCmd() *cobra.Command {
	opts := &initOptions{}

	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new repository",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.repoPath = args[0]
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					return err
				}
				opts.repoPath = cwd
			}
			return initRepo(cmd.OutOrStdout(), opts)
		},
	}
}

func initRepo(w io.Writer, opts *initOptions) error {
	absTracPath, err := filepath.Abs(filepath.Join(opts.repoPath, ".trac"))
	if err != nil {
		return err
	}

	directories := []string{
		filepath.Join(absTracPath, "objects"),
		filepath.Join(absTracPath, "refs", "heads"),
	}

	if _, err := os.Stat(absTracPath); err != nil {
		if os.IsNotExist(err) {
			for _, dir := range directories {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", dir, err)
				}
			}

			headPath := filepath.Join(absTracPath, "HEAD")
			if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
				return fmt.Errorf("failed to write HEAD file: %w", err)
			}

			mainBranchPath := filepath.Join(absTracPath, "refs", "heads", "main")
			if err := os.WriteFile(mainBranchPath, []byte{}, 0644); err != nil {
				return fmt.Errorf("failed to write main branch reference: %w", err)
			}

			fmt.Fprintf(w, "Initialized empty trac repository in %s\n", absTracPath)
			return nil
		}
		return fmt.Errorf("failed to check for existing repository at %s: %w", absTracPath, err)
	}

	fmt.Fprintf(w, "Reinitialized existing trac repository in %s\n", absTracPath)
	return nil
}
