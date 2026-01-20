package cmd

import (
	"fmt"
	"os"
	"sort"

	"camp/internal/project"
	"camp/internal/shell"

	"github.com/spf13/cobra"
)

var shellOverride string

func activateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Output environment variables for shell evaluation",
		Long: `Output shell export statements for environment variables defined in project camp.yaml.

This command reads the camp.yaml file in the current directory and outputs
shell-compatible export statements that can be evaluated to set environment
variables.

Variables that were previously exported but are no longer in the config will
be unset automatically.

Usage:
  eval "$(camp project activate)"

The shell type is auto-detected from the SHELL environment variable.
Supported shells: bash, zsh, fish.

Example camp.yaml:
  env:
    DATABASE_URL: "postgres://localhost:5432/mydb"
    API_KEY: "secret-key"
    DEBUG: "true"`,
		RunE: runActivate,
		// Silence usage on error to keep output clean for eval
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&shellOverride, "shell", "",
		"Override shell type detection (bash, zsh, fish)")

	return cmd
}

func runActivate(cmd *cobra.Command, args []string) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Determine shell type
	var shellType shell.ShellType
	if shellOverride != "" {
		shellType = shell.ShellType(shellOverride)
	} else {
		shellType = shell.DetectShell()
	}

	// Load previously tracked vars
	trackedVars, err := project.LoadTrackedVars(cwd)
	if err != nil {
		return fmt.Errorf("failed to load tracked vars: %w", err)
	}

	// Check if project config exists
	if !project.HasProjectConfig(cwd) {
		// No config - unset all previously tracked vars and clear tracking
		for _, varName := range trackedVars {
			fmt.Fprintln(cmd.OutOrStdout(), shell.FormatUnset(shellType, varName))
		}
		// Clear tracking file if there were tracked vars
		if len(trackedVars) > 0 {
			if err := project.SaveTrackedVars(cwd, []string{}); err != nil {
				return fmt.Errorf("failed to clear tracked vars: %w", err)
			}
		}
		return nil
	}

	// Load project config
	config, err := project.LoadProjectConfig(cwd)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Get removed vars and output unset statements
	removedVars := project.GetRemovedVars(trackedVars, config.Env)
	sort.Strings(removedVars)
	for _, varName := range removedVars {
		fmt.Fprintln(cmd.OutOrStdout(), shell.FormatUnset(shellType, varName))
	}

	// Sort current keys for deterministic output
	keys := make([]string, 0, len(config.Env))
	for k := range config.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Output export statements
	for _, key := range keys {
		value := config.Env[key]
		fmt.Fprintln(cmd.OutOrStdout(), shell.FormatExport(shellType, key, value))
	}

	// Save current var names for future tracking
	if err := project.SaveTrackedVars(cwd, keys); err != nil {
		return fmt.Errorf("failed to save tracked vars: %w", err)
	}

	return nil
}
