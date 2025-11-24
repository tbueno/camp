package system

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// NukeEnvironment removes camp and Nix completely from the system
func NukeEnvironment(user *User, out io.Writer) error {
	// Step 1: Uninstall Nix (platform-specific)
	fmt.Fprintf(out, "Uninstalling Nix...\n")
	if err := UninstallNix(user.Platform, out); err != nil {
		// Log error but continue with file cleanup
		fmt.Fprintf(out, "⚠️  Warning: Nix uninstall encountered an error: %v\n", err)
		fmt.Fprintf(out, "Continuing with file cleanup...\n")
	} else {
		fmt.Fprintf(out, "✓ Nix uninstalled successfully\n")
	}

	// Step 2: Remove camp files
	fmt.Fprintf(out, "Removing camp files...\n")
	campDir := filepath.Join(user.HomeDir, ".camp")
	if err := RemoveCampFiles(campDir, out); err != nil {
		return fmt.Errorf("failed to remove camp files: %w", err)
	}
	fmt.Fprintf(out, "✓ Camp files removed successfully\n")

	// Step 3: Remove Nix state files
	fmt.Fprintf(out, "Removing Nix state files...\n")
	if err := RemoveNixStateFiles(user.HomeDir, out); err != nil {
		// Log error but don't fail - some files may not exist
		fmt.Fprintf(out, "⚠️  Warning: Some Nix state files could not be removed: %v\n", err)
	} else {
		fmt.Fprintf(out, "✓ Nix state files removed successfully\n")
	}

	return nil
}

// UninstallNix uninstalls Nix package manager based on the platform
func UninstallNix(platform string, out io.Writer) error {
	// On macOS, first try to uninstall nix-darwin
	if platform == "darwin" {
		fmt.Fprintf(out, "  Uninstalling nix-darwin...\n")
		cmd := exec.Command("sudo", "darwin-uninstaller")
		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(out, "  ⚠️  darwin-uninstaller not found or failed (this is OK if not using nix-darwin)\n")
		}
	}

	// Run the Nix installer's uninstall command
	fmt.Fprintf(out, "  Running Nix uninstaller...\n")
	cmd := exec.Command("/nix/nix-installer", "uninstall")
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("nix-installer uninstall failed: %w", err)
	}

	return nil
}

// RemoveCampFiles removes the ~/.camp directory and its contents
func RemoveCampFiles(campDir string, out io.Writer) error {
	return deleteIfExists(campDir, out)
}

// RemoveNixStateFiles removes home-manager and Nix state directories
func RemoveNixStateFiles(homeDir string, out io.Writer) error {
	// List of directories to remove
	dirsToRemove := []string{
		filepath.Join(homeDir, ".config", "home-manager"),
		filepath.Join(homeDir, ".local", "state", "home-manager"),
		filepath.Join(homeDir, ".local", "state", "nix"),
	}

	var firstError error
	for _, dir := range dirsToRemove {
		if err := deleteIfExists(dir, out); err != nil && firstError == nil {
			firstError = err
		}
	}

	return firstError
}

// deleteIfExists deletes a file or directory if it exists
func deleteIfExists(path string, out io.Writer) error {
	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(out, "  • %s (not found, skipping)\n", path)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check %s: %w", path, err)
	}

	// Remove the path
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to delete %s: %w", path, err)
	}

	fmt.Fprintf(out, "  • Deleted: %s\n", path)
	return nil
}
