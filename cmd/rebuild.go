package cmd

import (
	"fmt"
	"io"
	"strings"

	"camp/internal/system"

	"github.com/spf13/cobra"
)

var rebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild the development environment",
	Long: `Rebuild the development environment based on the latest configuration.

This command:
  1. Copies Nix configuration files from templates to ~/.camp/nix/
  2. Renders flake.nix with your configuration from camp.yml
  3. Refreshes local flakes (path:, git+file: or absolute path URLs) so
     their latest changes are applied. Remote flakes stay pinned until
     you run 'camp env update'
  4. Executes the platform-specific rebuild command:
     - macOS: Uses nix-darwin to rebuild system configuration
     - Linux: Uses home-manager to rebuild user environment

Prerequisites:
  - Nix package manager must be installed
  - macOS: nix-darwin must be configured (requires sudo/admin privileges)
  - Linux: home-manager must be configured

Note: On macOS, this command requires sudo privileges and will prompt for your password.`,
	RunE: runRebuild,
}

func init() {
	envCmd.AddCommand(rebuildCmd)
}

func runRebuild(cmd *cobra.Command, args []string) error {
	// Get current user context
	user := system.NewUser()

	// Output rebuild start message
	fmt.Fprintf(cmd.OutOrStdout(), "Starting environment rebuild...\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Platform: %s\n", user.Platform)
	fmt.Fprintf(cmd.OutOrStdout(), "User: %s\n", user.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "Hostname: %s\n\n", user.HostName)

	// Prepare environment (copy files and render templates)
	fmt.Fprintf(cmd.OutOrStdout(), "Preparing environment...\n")
	if err := system.PrepareEnvironment(user); err != nil {
		return fmt.Errorf("failed to prepare environment: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "✓ Environment prepared successfully\n\n")

	if err := refreshLocalFlakes(cmd.OutOrStdout(), user); err != nil {
		return err
	}

	// Execute rebuild
	fmt.Fprintf(cmd.OutOrStdout(), "Executing rebuild command...\n")
	if err := system.ExecuteRebuild(user); err != nil {
		return fmt.Errorf("rebuild failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\n✓ Environment rebuild completed successfully!\n")
	return nil
}

// refreshLocalFlakes updates the lock entries of local flakes, skipping silently when there are none
func refreshLocalFlakes(out io.Writer, user *system.User) error {
	names := system.LocalFlakeNames(user.Flakes)
	if len(names) == 0 {
		return nil
	}

	fmt.Fprintf(out, "Refreshing local flakes: %s\n", strings.Join(names, ", "))
	if err := system.UpdateLocalFlakes(user); err != nil {
		return err
	}
	fmt.Fprintf(out, "✓ Local flakes refreshed\n\n")

	return nil
}
