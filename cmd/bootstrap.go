package cmd

import (
	"camp/internal/system"
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

var dryRun bool

// fsDirWrapper wraps fs.FS to implement utils.TemplDir interface
type fsDirWrapper struct {
	fsys fs.FS
}

func (w fsDirWrapper) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(w.fsys, name)
}

func (w fsDirWrapper) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(w.fsys, name)
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap your development environment with Nix",
	Long:  "Bootstrap command sets up your development environment by installing Nix and configuring your home directory with the necessary tools and configuration files.",
	Run: func(cmd *cobra.Command, args []string) {
		if dryRun {
			fmt.Fprintln(cmd.OutOrStdout(), "Running in dry-run mode - no actual installations will be performed")
		}

		// Get the current working directory to locate templates
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Failed to get current directory: %v\n", err)
			return
		}
		templatesFS := fsDirWrapper{fsys: os.DirFS(cwd)}
		err = system.RunBootstrapWithHome(templatesFS, cmd.OutOrStdout(), dryRun)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Bootstrap failed: %v\n", err)
			return
		}
	},
}

func init() {
	bootstrapCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be installed without actually installing")
}
