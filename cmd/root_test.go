package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	t.Run("command name", func(t *testing.T) {
		if rootCmd.Use != "camp" {
			t.Errorf("Expected command name to be 'camp', got '%s'", rootCmd.Use)
		}
	})

	t.Run("command description", func(t *testing.T) {
		expected := "Camp is your all-in-one dev environment manager"
		if rootCmd.Short != expected {
			t.Errorf("Expected short description to be '%s', got '%s'", expected, rootCmd.Short)
		}
	})

	t.Run("command execution and output", func(t *testing.T) {
		var output bytes.Buffer

		// Create a copy of the root command to avoid modifying the original
		cmd := &cobra.Command{
			Use:   rootCmd.Use,
			Short: rootCmd.Short,
			Long:  rootCmd.Long,
			Run:   rootCmd.Run,
		}

		cmd.SetOut(&output)
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		if err != nil {
			t.Errorf("Command execution failed: %v", err)
		}

		expectedOutput := "Hello! Welcome to camp - your dev environment manager!"
		actualOutput := strings.TrimSpace(output.String())

		if actualOutput != expectedOutput {
			t.Errorf("Expected output '%s', got '%s'", expectedOutput, actualOutput)
		}
	})
}
