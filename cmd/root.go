package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "camp",
	Short: "Camp is your all-in-one dev environment manager",
	Long:  "Camp is a command line application helps you managing your isolated development environment.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "Hello! Welcome to camp - your dev environment manager!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(bootstrapCmd)
	rootCmd.AddCommand(projectCmd)
}
