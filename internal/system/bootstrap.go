package system

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func GetDefaultBootstrapConfig() *BootstrapConfig {
	return &BootstrapConfig{
		Applications: []Application{
			{
				Name:           "direnv",
				InstallCommand: "curl -sfL https://direnv.net/install.sh | bash",
			},
			{
				Name:           "nix",
				InstallCommand: "curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install --determinate",
			},
			{
				Name:           "devbox",
				InstallCommand: "curl -fsSL https://get.jetify.com/devbox | bash",
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
