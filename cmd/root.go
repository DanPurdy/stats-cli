package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitlab-cli",
	Short: "A CLI application to interact with GitLab API",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}