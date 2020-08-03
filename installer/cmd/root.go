package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "argo-agent",
	Short: "Codefresh argocd agent",
	Long:  `Codefresh argocd agent`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exec hugo")
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
