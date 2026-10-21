package system

import (
	"camp/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

// PrepareEnvironment prepares the environment for rebuild by copying
// config files and compiling templates
func PrepareEnvironment(user *User) error {
	// Ensure .camp/nix directory exists
	nixDir := filepath.Join(user.HomeDir, ".camp", "nix")
	if err := os.MkdirAll(nixDir, 0755); err != nil {
		return fmt.Errorf("failed to create nix directory: %w", err)
	}

	// Copy config files
	if err := CopyConfigFiles(user); err != nil {
		return fmt.Errorf("failed to copy config files: %w", err)
	}

	// Compile templates
	if err := CompileTemplates(user); err != nil {
		return fmt.Errorf("failed to compile templates: %w", err)
	}

	return nil
}

// CopyConfigFiles copies .nix configuration files from templates/files/
// to ~/.camp/nix/, excluding flake.nix which is rendered separately
func CopyConfigFiles(user *User) error {
	srcDir := "templates/files"
	destDir := filepath.Join(user.HomeDir, ".camp", "nix")

	// Read source directory
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	// Copy each entry except flake.nix (which is rendered)
	for _, entry := range entries {
		// Skip flake.nix - it will be rendered by CompileTemplates
		if entry.Name() == "flake.nix" {
			continue
		}

		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Recursively copy directory (e.g., modules/)
			if err := copyDir(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy directory %s: %w", entry.Name(), err)
			}
		} else {
			// Copy file (e.g., mac.nix, linux.nix)
			if err := utils.CopyFile(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dest string) error {
	// Get source directory info
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(dest, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Read source directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	// Copy each entry
	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(sourcePath, destPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := utils.CopyFile(sourcePath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CompileTemplates renders the flake.nix template with user data
func CompileTemplates(user *User) error {
	// Reload user config to get latest env vars
	if err := user.Reload(); err != nil {
		return fmt.Errorf("failed to reload user config: %w", err)
	}

	// Render flake.nix template
	if err := RenderFlakeTemplate(user); err != nil {
		return fmt.Errorf("failed to render flake template: %w", err)
	}

	return nil
}

// LocalFlakeNames returns the names of the flakes sourced from the local filesystem
func LocalFlakeNames(flakes []Flake) []string {
	var names []string
	for _, flake := range flakes {
		if flake.IsLocal() {
			names = append(names, flake.Name)
		}
	}
	return names
}

// localFlakeUpdateArgs builds the nix arguments to refresh only the local flake inputs.
// Returns nil when there are no local flakes.
func localFlakeUpdateArgs(nixDir string, flakes []Flake) []string {
	names := LocalFlakeNames(flakes)
	if len(names) == 0 {
		return nil
	}

	args := []string{
		"--extra-experimental-features",
		"nix-command flakes",
		"flake",
		"update",
	}
	args = append(args, names...)
	return append(args, "--flake", nixDir)
}

// UpdateLocalFlakes refreshes the lock entries of local flakes so their latest
// contents are used on rebuild. Remote inputs stay pinned.
func UpdateLocalFlakes(user *User) error {
	nixDir := filepath.Join(user.HomeDir, ".camp", "nix")
	args := localFlakeUpdateArgs(nixDir, user.Flakes)
	if args == nil {
		return nil
	}

	if err := utils.RunCommand("nix", args...); err != nil {
		return fmt.Errorf("failed to update local flakes: %w", err)
	}

	return nil
}

var darwinEtcFiles = []string{
	"/etc/bashrc",
	"/etc/zshrc",
	"/etc/zprofile",
}

func backupDarwinEtcFiles() error {
	for _, f := range darwinEtcFiles {
		if err := utils.BackupFile(f); err != nil {
			return fmt.Errorf("failed to back up %s: %w", f, err)
		}
	}
	return nil
}

// ExecuteRebuild runs the platform-specific rebuild command
func ExecuteRebuild(user *User) error {
	nixDir := filepath.Join(user.HomeDir, ".camp", "nix")

	// Check if nix directory exists
	if _, err := os.Stat(nixDir); os.IsNotExist(err) {
		return fmt.Errorf("nix directory does not exist: %s", nixDir)
	}

	var cmd string
	var args []string

	switch user.Platform {
	case "darwin":
		if err := backupDarwinEtcFiles(); err != nil {
			return fmt.Errorf("failed to back up /etc files: %w", err)
		}

		cmd = "sudo"
		args = []string{
			"nix",
			"--extra-experimental-features",
			"nix-command flakes",
			"run",
			"nix-darwin#darwin-rebuild",
			"--", // Separator: everything after this goes to darwin-rebuild
			"switch",
			"--impure",
			"--flake",
			fmt.Sprintf("%s#%s", nixDir, user.HostName),
		}

	case "linux":
		// Linux: use home-manager
		cmd = "home-manager"
		args = []string{
			"switch",
			"--impure",
			"-b",
			"backup",
			"--flake",
			fmt.Sprintf("%s#%s", nixDir, user.Name),
		}

	default:
		return fmt.Errorf("unsupported platform: %s", user.Platform)
	}

	// Execute rebuild command
	if err := utils.RunCommand(cmd, args...); err != nil {
		return fmt.Errorf("rebuild command failed: %w", err)
	}

	return nil
}
