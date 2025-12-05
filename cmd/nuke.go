package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"camp/internal/system"

	"github.com/spf13/cobra"
)

// nixInstalledChecker is a variable that can be overridden in tests
var nixInstalledChecker = isNixInstalled

var nukeCmd = &cobra.Command{
	Use:     "nuke",
	Aliases: []string{"destroy"},
	Short:   "Remove camp and Nix completely from the system",
	Long: `Remove camp and Nix completely from the system.

WARNING: This is a destructive operation that cannot be undone!

This command will:
  1. Uninstall Nix package manager (including nix-darwin on macOS)
  2. Remove all camp configuration files from ~/.camp/
  3. Remove home-manager configuration and state
  4. Remove Nix state directories

After running this command, you will need to restart your terminal and
run 'camp bootstrap' again if you want to use camp in the future.`,
	RunE: runNuke,
}

var (
	skipConfirmation bool
)

func init() {
	envCmd.AddCommand(nukeCmd)
	nukeCmd.Flags().BoolVarP(&skipConfirmation, "yes", "y", false, "Skip confirmation prompt")
}

func runNuke(cmd *cobra.Command, args []string) error {
	// Check if Nix is installed
	if !nixInstalledChecker() {
		return fmt.Errorf("nix is not installed. No need to run nuke")
	}

	// Prompt user for confirmation (unless --yes flag is used)
	if !skipConfirmation {
		fmt.Fprintf(cmd.OutOrStdout(), "⚠️  WARNING: This will completely remove camp and Nix from your system!\n\n")
		fmt.Fprintf(cmd.OutOrStdout(), "This will delete:\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  • Nix package manager\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  • All camp configuration (~/.camp/)\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  • home-manager configuration and state\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  • Nix state directories\n\n")
		fmt.Fprintf(cmd.OutOrStdout(), "Are you sure you want to continue? [y/N]: ")

		reader := bufio.NewReader(cmd.InOrStdin())
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Fprintf(cmd.OutOrStdout(), "\nNuke operation cancelled.\n")
			return nil
		}
	}

	// Get current user context
	user := system.NewUser()

	// Execute nuke
	fmt.Fprintf(cmd.OutOrStdout(), "\nStarting nuke process...\n")
	if err := system.NukeEnvironment(user, cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("nuke failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\n✓ Your camp environment was erased. Please restart your terminal to complete the process.\n")
	return nil
}

func isNixInstalled() bool {
	cmd := exec.Command("nix", "--version")
	err := cmd.Run()
	return err == nil
}
