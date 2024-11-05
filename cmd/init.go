package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/lucasrod16/trac/internal/layout"
	"github.com/spf13/cobra"
)

type initOptions struct {
	// repo path argument
	repoPath string
}

func NewInitCmd() *cobra.Command {
	opts := &initOptions{}

	return &cobra.Command{
		Use:   "init [repoPath]",
		Short: "Initialize a new trac repository",
		Long:  "Initialize a new trac repository at the specified path. If no path is provided, the current working directory will be used.",
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

			l, err := layout.New(opts.repoPath)
			if err != nil {
				return err
			}

			return initRepo(cmd.OutOrStdout(), l)
		},
	}
}

func initRepo(w io.Writer, l *layout.Layout) error {
	if l.Exists() {
		fmt.Fprintf(w, "Reinitialized existing trac repository in %s\n", l.Root)
		return nil
	}

	if err := l.Init(); err != nil {
		return err
	}

	fmt.Fprintf(w, "Initialized empty trac repository in %s\n", l.Root)
	return nil
}
