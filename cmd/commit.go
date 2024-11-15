package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/lucasrod16/trac/internal/commit"
	"github.com/lucasrod16/trac/internal/index"
	"github.com/lucasrod16/trac/internal/layout"
	"github.com/spf13/cobra"
)

type commitOptions struct {
	message string // -m, --message
}

func NewCommitCmd() *cobra.Command {
	opts := &commitOptions{}

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Record changes to the repository",
		Long: `
	Create a new commit containing the current contents of the index and the given log message describing the changes.
	The new commit is a direct child of HEAD, usually the tip of the current branch, and the branch is updated to point to it.
	`,
		Args: cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommit(cmd.OutOrStdout(), opts)
		},
	}
	cmd.Flags().StringVarP(&opts.message, "message", "m", "", "Commit message")
	cmd.MarkFlagRequired("message")
	return cmd
}

func runCommit(w io.Writer, opts *commitOptions) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	layout, err := layout.New(cwd)
	if err != nil {
		return err
	}
	if err := layout.ValidateIsRepo(); err != nil {
		return err
	}
	idx := index.New()
	if err := idx.Load(layout); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return commit.ErrNothingAddedToCommit
		}
		return err
	}
	parentHash, err := commit.GetParentHash(layout)
	if err != nil {
		return err
	}
	newCommit := commit.New(opts.message, parentHash, idx.Staged)
	commitHash, err := newCommit.Save(layout)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Created commit %s\n", commitHash)
	return nil
}
