package cmd

import (
	"fmt"

	"camp/internal/system"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Display development environment information",
	Long:  "Display information about the current development environment including system architecture.",
	Run: func(cmd *cobra.Command, args []string) {
		sysInfo, err := system.GetSystemInfo()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting system information: %v\n", err)
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Architecture: %s\n", sysInfo.Architecture)
		fmt.Fprintf(cmd.OutOrStdout(), "OS: %s\n", sysInfo.OS)
	},
}
