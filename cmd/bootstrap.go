package cmd

import (
	"camp/internal/system"
	"fmt"

	"github.com/spf13/cobra"
)

var dryRun bool

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Install applications needed for camp environment management",
	Long:  "Bootstrap command installs a list of applications that are required for camp environment management to work properly.",
	Run: func(cmd *cobra.Command, args []string) {
		config := system.GetDefaultBootstrapConfig()

		if dryRun {
			fmt.Fprintln(cmd.OutOrStdout(), "Running in dry-run mode - no actual installations will be performed")
		}

		err := system.RunBootstrap(config, cmd.OutOrStdout(), dryRun)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Bootstrap failed: %v\n", err)
			return
		}
	},
}

func init() {
	bootstrapCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be installed without actually installing")
}
