package cmd

import (
	"log"

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
	rootCmd.AddCommand(NewCommitCmd())
	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
