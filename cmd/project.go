package cmd

import (
	"fmt"

	"camp/internal/project"
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

// validateArgs checks if the subcommand exists in the devbox.json file
func validateArgs(args []string) error {
	proj := project.NewProject()
	if len(args) == 0 {
		return nil
	}
	subcommands := append(proj.CommandNames(), "info")

	for _, sub := range subcommands {
		if args[0] == sub {
			return nil
		}
	}
	return fmt.Errorf("unknown subcommand: %s", args[0])
}

func init() {
	projectCmd.AddCommand(infoCmd())
	projectCmd.AddCommand(installCmd())
	projectCmd.AddCommand(testCmd())
}
