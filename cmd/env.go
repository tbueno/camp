package cmd

import (
	"fmt"
	"io"
	"os"

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

		err = printDirenvVars(cmd.OutOrStdout())
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error reading .envrc: %v\n", err)
		}
	},
}

func printDirenvVars(out io.Writer) error {
	file, err := os.Open(".envrc")
	if err != nil {
		return err
	} else {
		defer file.Close()
		envVars, err := system.GetExportedVars(file)
		if err != nil {
			return err
		}

		if len(envVars) > 0 {
			fmt.Fprintf(out, "\nDirenv variables:\n")
			for _, envVar := range envVars {
				fmt.Fprintf(out, "%s=%s\n", envVar.Name, envVar.Value)
			}
		}
		return nil
	}
}
