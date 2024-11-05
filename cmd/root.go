package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "trac",
		Short: "A barebones version control system",
	}
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewStatusCmd())
	rootCmd.AddCommand(NewAddCmd())
	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
