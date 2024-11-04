package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
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
	status := status.NewRepoStatus(l)

	if err := status.DetectTrackedStatus(); err != nil {
		return err
	}

	if err := status.GetStagedFiles(); err != nil {
		return err
	}

	if !status.HasTracked() && !status.HasUntracked() {
		fmt.Fprintln(w, "nothing to commit (create/copy files and use \"trac add\" to track)")
		return nil
	}

	if status.HasUntracked() {
		fmt.Fprintln(w, "Untracked files:")
		for _, file := range status.GetUntracked() {
			untrackedColor := color.New(color.FgHiRed)
			untrackedColor.Fprintf(w, "\t%s\n", file)
		}
	}

	if status.HasTracked() {
		fmt.Fprintln(w, "\nChanges to be committed:")
		for _, file := range status.GetTracked() {
			fmt.Fprintf(w, "\tnew file:   %s\n", file)
		}
	}

	return nil
}
