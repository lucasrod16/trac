package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasrod16/trac/internal/index"
	"github.com/lucasrod16/trac/internal/layout"
	"github.com/spf13/cobra"
)

type addOptions struct {
	// Files to stage
	files []string
}

func NewAddCmd() *cobra.Command {
	opts := &addOptions{}

	return &cobra.Command{
		Use:   "add [file...]",
		Short: "Add file contents to the index",
		Long: `
	This command updates the index using the current content found in the working tree, to prepare the content staged for the next commit.

	The "index" holds a snapshot of the content of the working tree, and it is this snapshot that is taken as the contents of the next commit. Thus after
	making any changes to the working tree, and before running the commit command, you must use the add command to add any new or modified files to the
	index.

	This command can be performed multiple times before a commit. It only adds the content of the specified file(s) at the time the add command is run; if
	you want subsequent changes included in the next commit, then you must run trac add again to add the new content to the index.

	The trac status command can be used to obtain a summary of which files have changes that are staged for the next commit.
	`,
		Args: cobra.MinimumNArgs(1),
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

			if args[0] == "." {
				files, err := getFilesRecursively(cwd, l)
				if err != nil {
					return err
				}
				opts.files = files
			} else {
				opts.files = args
			}

			if err := stageFiles(l, opts); err != nil {
				return err
			}

			return nil
		},
	}
}

func getFilesRecursively(rootDir string, l *layout.Layout) ([]string, error) {
	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.Contains(path, l.Root) || strings.Contains(path, ".git") {
				return filepath.SkipDir
			}
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func stageFiles(l *layout.Layout, opts *addOptions) error {
	idx := index.New()
	err := idx.Load(l)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	for _, file := range opts.files {
		absPath, err := filepath.Abs(file)
		if err != nil {
			return err
		}
		if err := idx.Add(absPath); err != nil {
			return fmt.Errorf("failed to add file %s: %w", file, err)
		}
	}

	if err := idx.Write(l); err != nil {
		return fmt.Errorf("failed to write updated index: %w", err)
	}

	return nil
}
