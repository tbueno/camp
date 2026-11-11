package cmd

import (
	"camp/internal/system"
	"camp/templates"
	"fmt"

	"github.com/spf13/cobra"
)

var dryRun bool

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap your development environment with Nix",
	Long:  "Bootstrap command sets up your development environment by installing Nix and configuring your home directory with the necessary tools and configuration files.",
	Run: func(cmd *cobra.Command, args []string) {
		if dryRun {
			fmt.Fprintln(cmd.OutOrStdout(), "Running in dry-run mode - no actual installations will be performed")
		}

		if err := system.RunBootstrapWithHome(templates.FS, cmd.OutOrStdout(), dryRun); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Bootstrap failed: %v\n", err)
			return
		}
	},
}

func init() {
	bootstrapCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be installed without actually installing")
}
