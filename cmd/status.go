package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lucasrod16/trac/internal/index"
	"github.com/lucasrod16/trac/internal/layout"
	"github.com/lucasrod16/trac/internal/status"
	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [repoPath]",
		Short: "Show the status of the trac repository",
		Long:  "Display the current status of the trac repository, listing untracked files if any.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			l, err := layout.New(cwd)
			if err != nil {
				return err
			}
			if err := l.ValidateIsRepo(); err != nil {
				return err
			}
			return showRepoStatus(cmd.OutOrStdout(), l)
		},
	}
}

// showRepoStatus outputs the current status of the repository.
func showRepoStatus(w io.Writer, l *layout.Layout) error {
	idx := index.New()
	err := idx.Load(l)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	repoStatus, err := status.Get(idx)
	if err != nil {
		return err
	}
	if !repoStatus.HasTracked() && !repoStatus.HasUntracked() {
		fmt.Fprintln(w, "nothing to commit (create/copy files and use \"trac add\" to track)")
		return nil
	}
	if repoStatus.HasUntracked() {
		fmt.Fprintln(w, "Untracked files:")
		for _, path := range repoStatus.Untracked() {
			untracked := preparePath(path)
			untrackedColor := color.New(color.FgHiRed)
			untrackedColor.Fprintf(w, "\t%s\n", untracked)
		}
	}
	if repoStatus.HasTracked() {
		fmt.Fprintln(w, "\nChanges to be committed:")
		for _, filepath := range repoStatus.Tracked() {
			trackedColor := color.New(color.FgHiGreen)
			trackedColor.Fprintf(w, "\tnew file:   %s\n", filepath)
		}
	}
	return nil
}

func preparePath(path string) string {
	parts := strings.Split(path, string(filepath.Separator))
	root := parts[0]
	if len(parts) > 1 {
		root += string(filepath.Separator)
	}
	return root
}
