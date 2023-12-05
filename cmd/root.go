package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "balance_checker",
	Short: "Balance checker",
	Long:  "Balance checker",
	Run:   func(cmd *cobra.Command, args []string) {},
}

type contextKey string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
