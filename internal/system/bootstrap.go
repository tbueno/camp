package system

import (
	"camp/internal/utils"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func GetDefaultBootstrapConfig() *BootstrapConfig {
	return &BootstrapConfig{
		Applications: []Application{
			{
				Name:           "nix",
				InstallCommand: "curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install --determinate",
			},
		},
	}
}

func RunBootstrap(config *BootstrapConfig, output io.Writer, dryRun bool) error {
	if config == nil {
		config = GetDefaultBootstrapConfig()
	}

	fmt.Fprintf(output, "Starting bootstrap process for %d applications...\n\n", len(config.Applications))

	for i, app := range config.Applications {
		fmt.Fprintf(output, "[%d/%d] Checking %s...\n", i+1, len(config.Applications), app.Name)

		// Check if the executable already exists
		if isCommandAvailable(app.Name) {
			fmt.Fprintf(output, "  ⏭️  %s is already installed, skipping\n", app.Name)
			continue
		}

		fmt.Fprintf(output, "  Installing %s...\n", app.Name)

		if dryRun {
			fmt.Fprintf(output, "  [DRY RUN] Would execute: %s\n", app.InstallCommand)
			continue
		}

		err := executeInstallCommand(app.InstallCommand, output)
		if err != nil {
			fmt.Fprintf(output, "  ❌ Failed to install %s: %v\n", app.Name, err)
			return fmt.Errorf("bootstrap failed at application %s: %w", app.Name, err)
		}

		fmt.Fprintf(output, "  ✅ %s installation completed\n", app.Name)
	}

	fmt.Fprintf(output, "\n🎉 Bootstrap process completed successfully!\n")
	return nil
}

// isCommandAvailable checks if a command is available in the system PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func executeInstallCommand(command string, output io.Writer) error {
	if strings.TrimSpace(command) == "" {
		return fmt.Errorf("empty install command")
	}

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = output
	cmd.Stderr = output

	return cmd.Run()
}

// RunBootstrapWithHome runs the bootstrap process with home directory setup
func RunBootstrapWithHome(templDir utils.TemplDir, output io.Writer, dryRun bool) error {
	user := NewUser()
	campPath := user.HomeDir + "/.camp"

	fmt.Fprintf(output, "Bootstrapping your environment...\n")

	if err := bootstrapHome(campPath, templDir, output, dryRun); err != nil {
		return fmt.Errorf("failed to bootstrap home: %w", err)
	}

	if err := utils.RunCommand("nix", "--version"); err != nil {
		fmt.Fprintf(output, "Nix is not installed. Installing Nix...\n")

		if dryRun {
			fmt.Fprintf(output, "[DRY RUN] Would execute: %s/bin/install_nix\n", campPath)
		} else {
			if err = utils.RunCommand(campPath+"/bin"+"/install_nix", campPath); err != nil {
				return fmt.Errorf("failed to install nix: %w", err)
			}
		}
	}

	if user.Platform == "linux" {
		fmt.Fprintf(output, "Bootstrapping Linux...\n")
		return bootstrapLinux(campPath, templDir, user, output, dryRun)
	} else {
		fmt.Fprintf(output, "Bootstrapping macOS...\n")
		return bootstrapMac(campPath, templDir, user, output, dryRun)
	}
}

// bootstrapHome sets up the home directory structure
func bootstrapHome(campPath string, templDir utils.TemplDir, output io.Writer, dryRun bool) error {
	folders := []string{
		campPath,
		campPath + "/nix",
		campPath + "/bin",
	}

	for _, folder := range folders {
		if dryRun {
			fmt.Fprintf(output, "[DRY RUN] Would create directory: %s\n", folder)
		} else {
			if err := os.MkdirAll(folder, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", folder, err)
			}
		}
	}

	if err := copyBinFiles(templDir, campPath, output, dryRun); err != nil {
		return fmt.Errorf("failed to copy bin files: %w", err)
	}

	return nil
}

// bootstrapMac sets up macOS-specific configuration with nix-darwin
func bootstrapMac(campPath string, templDir utils.TemplDir, user *User, output io.Writer, dryRun bool) error {
	// Read and process darwin.nix (flake.nix)
	darwinContent, err := templDir.ReadFile("templates/initial/darwin.nix")
	if err != nil {
		return err
	}
	darwinContent = utils.ReplaceInContent(darwinContent, "__USER__", user.Name)
	darwinContent = utils.ReplaceInContent(darwinContent, "__HOME__", user.HomeDir)

	// Read home.nix (no replacements needed - nix-darwin handles user info)
	homeContent, err := templDir.ReadFile("templates/initial/home.nix")
	if err != nil {
		return err
	}

	nixHome := campPath + "/nix"
	flakePath := nixHome + "/flake.nix"
	homePath := nixHome + "/home.nix"

	if dryRun {
		fmt.Fprintf(output, "[DRY RUN] Would write darwin.nix to: %s\n", flakePath)
		fmt.Fprintf(output, "[DRY RUN] Would write home.nix to: %s\n", homePath)
	} else {
		// Save flake.nix
		err = utils.SaveFile(darwinContent, flakePath)
		if err != nil {
			return err
		}

		// Save home.nix
		err = utils.SaveFile(homeContent, homePath)
		if err != nil {
			return err
		}

		utils.BackupFile("/etc/bashrc")
		utils.BackupFile("/etc/zshrc")
		utils.BackupFile("/etc/nix/nix.conf")

		fmt.Fprintf(output, "Loading nix-darwin for the first time. This may take a while...\n")
		return utils.RunCommand(campPath+"/bin"+"/bootstrap", campPath)
	}

	return nil
}

// bootstrapLinux sets up Linux-specific configuration with home-manager
func bootstrapLinux(campPath string, templDir utils.TemplDir, user *User, output io.Writer, dryRun bool) error {
	// TODO: Implement Linux bootstrap support
	// This should:
	// 1. Determine system based on architecture (x86_64-linux or aarch64-linux)
	// 2. Copy flake.nix and home.nix with proper replacements (__USER__, __HOME__, __SYSTEM__)
	// 3. Run the bootstrap script
	return fmt.Errorf("Linux bootstrap is not yet implemented")
}

// copyBinFiles copies the bin files from the templates directory to the camp bin directory
func copyBinFiles(templDir utils.TemplDir, campPath string, output io.Writer, dryRun bool) error {
	binFiles, err := templDir.ReadDir("templates/initial/bin")
	if err != nil {
		return fmt.Errorf("failed to read directory bin: %w", err)
	}

	for _, file := range binFiles {
		if !file.IsDir() {
			destPath := campPath + "/bin/" + file.Name()

			if dryRun {
				fmt.Fprintf(output, "[DRY RUN] Would copy bin file: %s\n", destPath)
				continue
			}

			content, err := templDir.ReadFile("templates/initial/bin/" + file.Name())
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
			}

			if err = utils.SaveFile(content, destPath); err != nil {
				return fmt.Errorf("failed to save file %s: %w", destPath, err)
			}

			// change file permissions to be executable
			if err = os.Chmod(destPath, 0755); err != nil {
				return fmt.Errorf("failed to change file permissions: %w", err)
			}
		}
	}

	return nil
}
