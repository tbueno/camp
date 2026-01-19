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

	// Check if project config exists
	if !project.HasProjectConfig(cwd) {
		// Output nothing if no config - this allows safe eval even without config
		return nil
	}

	// Load project config
	config, err := project.LoadProjectConfig(cwd)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Check if there are any env vars to export
	if len(config.Env) == 0 {
		return nil
	}

	// Determine shell type
	var shellType shell.ShellType
	if shellOverride != "" {
		shellType = shell.ShellType(shellOverride)
	} else {
		shellType = shell.DetectShell()
	}

	// Sort keys for deterministic output
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

	return nil
}
