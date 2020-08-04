package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "argo-agent",
	Short: "Codefresh argocd agent",
	Long:  `Codefresh argocd agent`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
