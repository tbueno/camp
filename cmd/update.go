package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"camp/internal/system"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update flake dependencies to their latest versions",
	Long: `Update all flake dependencies defined in camp.yml to their latest versions.

This command:
  1. Prepares the environment (copies files and renders templates with current config)
  2. Updates flake.lock with the latest versions of all flake inputs
  3. This includes both built-in flakes (nixpkgs, nix-darwin, home-manager)
     and any custom flakes defined in your camp.yml

After running this command, you'll need to run 'camp env rebuild' to apply
the updated dependencies.

Prerequisites:
  - Nix package manager must be installed with flakes enabled`,
	RunE: runUpdate,
}

func init() {
	envCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Get current user context
	user := system.NewUser()

	// Output update start message
	fmt.Fprintf(cmd.OutOrStdout(), "Starting flake update...\n")
	fmt.Fprintf(cmd.OutOrStdout(), "User: %s\n", user.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "Nix directory: %s/.camp/nix\n\n", user.HomeDir)

	// Prepare environment (copy files and render templates)
	fmt.Fprintf(cmd.OutOrStdout(), "Preparing environment...\n")
	if err := system.PrepareEnvironment(user); err != nil {
		return fmt.Errorf("failed to prepare environment: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "✓ Environment prepared successfully\n\n")

	// Update flakes
	fmt.Fprintf(cmd.OutOrStdout(), "Updating flake dependencies...\n")
	nixDir := filepath.Join(user.HomeDir, ".camp", "nix")

	// Run nix flake update
	nixCmd := exec.Command(
		"nix",
		"--extra-experimental-features", "nix-command",
		"--extra-experimental-features", "flakes",
		"flake", "update",
		"--flake", nixDir,
	)

	// Stream output to user
	nixCmd.Stdout = cmd.OutOrStdout()
	nixCmd.Stderr = cmd.ErrOrStderr()

	if err := nixCmd.Run(); err != nil {
		return fmt.Errorf("nix flake update failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\n✓ Flake dependencies updated successfully!\n")
	fmt.Fprintf(cmd.OutOrStdout(), "\nNext step: Run 'camp env rebuild' to apply the updates.\n")
	return nil
}
