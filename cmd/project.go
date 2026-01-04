package cmd

import (
	"fmt"
	"os"

	"camp/internal/project"
	"camp/internal/system"
	"camp/internal/utils"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"proj"},
	Short:   "Interact with current project",
	Long:    "This command allows the user to interact with the current project by running pre-defined scripts.",
	Args: func(cmd *cobra.Command, args []string) error {
		if v := validateArgs(args); v != nil {
			return fmt.Errorf("%s \n run 'camp project info' for help: ", v)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return runProjectCommand(args[0])
	},
}

func installCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install project dependencies",
		Long:  "Install project dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProjectCommand("install")
		},
	}
}

func testCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Execute project test suite",
		Long:  "Execute test script defined in the project's devbox.json file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProjectCommand("test")
		},
	}
}

func infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print information about current project",
		Long:  "Print information about current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			proj := project.NewProject()
			if !proj.Compatible() {
				return fmt.Errorf("project is not compatible with camp")
			}
			for _, line := range proj.Info() {
				fmt.Println(line)
			}
			return nil
		},
	}
}

func runProjectCommand(cmd string) error {
	proj := project.NewProject()
	if !proj.Compatible() {
		return fmt.Errorf("project is not compatible with camp")
	}

	commands := proj.Commands()
	if commands[cmd] == nil {
		return fmt.Errorf("command not found: %s", cmd)
	}
	return utils.RunCommands(commands[cmd])
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize .camp.yml in current project",
		Long:  "Create a .camp.yml configuration file with example settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initProjectConfig()
		},
	}
}

func syncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync .envrc from .camp.yml",
		Long:  "Generate or update .envrc file from project configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return syncProjectEnv()
		},
	}
}

func syncProjectEnv() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find .camp.yml in current dir or parents
	configPath := system.FindProjectConfigPath(cwd)
	if configPath == "" {
		return fmt.Errorf(".camp.yml not found in current directory or parents\nRun 'camp project init' to create one")
	}

	// Load project config
	config, err := system.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load .camp.yml: %w", err)
	}

	// Generate .envrc content
	envrcContent, err := system.GenerateEnvrc(cwd, config)
	if err != nil {
		return fmt.Errorf("failed to generate .envrc: %w", err)
	}

	// Write to .envrc in current directory
	envrcPath := ".envrc"
	if err := os.WriteFile(envrcPath, []byte(envrcContent), 0644); err != nil {
		return fmt.Errorf("failed to write .envrc: %w", err)
	}

	// Count env vars
	envVarCount := len(config.Env)

	fmt.Println("✓ Generated .envrc from .camp.yml")
	if envVarCount > 0 {
		fmt.Printf("  Exported %d environment variable(s)\n", envVarCount)
	}
	fmt.Println("")
	fmt.Println("⚠ Run 'direnv allow' to activate the environment")

	return nil
}

func initProjectConfig() error {
	// Check if .camp.yml already exists in current dir or parents
	existingConfig := system.FindProjectConfigPath(".")
	if existingConfig != "" {
		return fmt.Errorf(".camp.yml already exists at %s", existingConfig)
	}

	// Create template config in current directory
	configPath := ".camp.yml"
	template := `# Camp project configuration
# See https://tbueno.github.io/camp/docs/project-config

# Project-specific environment variables
# These are loaded via direnv when entering the directory
env:
  PROJECT_NAME: "my-project"
  # Add your environment variables here
  # Example:
  # DATABASE_URL: "postgres://localhost/mydb"
  # DEBUG: "true"

# Future: packages, flakes, scripts will go here
`

	// Write to file
	if err := os.WriteFile(configPath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create .camp.yml: %w", err)
	}

	fmt.Println("✓ Created .camp.yml")
	fmt.Println("  Next steps:")
	fmt.Println("  1. Edit .camp.yml to add your environment variables")
	fmt.Println("  2. Run 'camp project sync' to generate .envrc")
	fmt.Println("  3. Run 'direnv allow' to activate the environment")

	return nil
}

// validateArgs checks if the subcommand exists in the devbox.json file
func validateArgs(args []string) error {
	proj := project.NewProject()
	if len(args) == 0 {
		return nil
	}
	subcommands := append(proj.CommandNames(), "info", "init", "sync")

	for _, sub := range subcommands {
		if args[0] == sub {
			return nil
		}
	}
	return fmt.Errorf("unknown subcommand: %s", args[0])
}

func init() {
	projectCmd.AddCommand(initCmd())
	projectCmd.AddCommand(syncCmd())
	projectCmd.AddCommand(infoCmd())
	projectCmd.AddCommand(installCmd())
	projectCmd.AddCommand(testCmd())
}
