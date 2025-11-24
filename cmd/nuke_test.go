package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNukeCommand(t *testing.T) {
	t.Run("command name", func(t *testing.T) {
		if nukeCmd.Use != "nuke" {
			t.Errorf("Expected command name to be 'nuke', got '%s'", nukeCmd.Use)
		}
	})

	t.Run("command has destroy alias", func(t *testing.T) {
		found := false
		for _, alias := range nukeCmd.Aliases {
			if alias == "destroy" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'destroy' to be an alias for nuke command")
		}
	})

	t.Run("command description", func(t *testing.T) {
		expected := "Remove camp and Nix completely from the system"
		if nukeCmd.Short != expected {
			t.Errorf("Expected short description to be '%s', got '%s'", expected, nukeCmd.Short)
		}
	})

	t.Run("command has long description with warnings", func(t *testing.T) {
		if nukeCmd.Long == "" {
			t.Error("Expected long description to be set")
		}

		requiredPhrases := []string{
			"WARNING",
			"destructive",
			"cannot be undone",
		}

		for _, phrase := range requiredPhrases {
			if !strings.Contains(nukeCmd.Long, phrase) {
				t.Errorf("Expected long description to contain '%s'", phrase)
			}
		}
	})

	t.Run("command is subcommand of env", func(t *testing.T) {
		found := false
		for _, cmd := range envCmd.Commands() {
			if cmd.Use == "nuke" {
				found = true
				break
			}
		}
		if !found {
			t.Error("nuke command should be registered as subcommand of env")
		}
	})
}

func TestNukeCommandConfirmation(t *testing.T) {
	t.Run("aborts when user says no", func(t *testing.T) {
		// Mock Nix as installed
		originalChecker := nixInstalledChecker
		nixInstalledChecker = func() bool { return true }
		defer func() { nixInstalledChecker = originalChecker }()

		// Create command with input/output buffers
		var output bytes.Buffer
		var input bytes.Buffer

		// Create a new command instance to avoid state issues
		cmd := &cobra.Command{
			Use:   nukeCmd.Use,
			Short: nukeCmd.Short,
			Long:  nukeCmd.Long,
			RunE:  nukeCmd.RunE,
		}
		cmd.SetOut(&output)
		cmd.SetIn(&input)
		cmd.SetArgs([]string{})

		// Simulate user typing 'n' for no
		input.WriteString("n\n")

		// Execute command
		err := cmd.Execute()

		// Should not return error when cancelled
		if err != nil {
			t.Errorf("Expected no error when user cancels, got: %v", err)
		}

		outputStr := output.String()

		// Should show warning
		if !strings.Contains(outputStr, "WARNING") {
			t.Errorf("Expected output to contain warning, got:\n%s", outputStr)
		}

		// Should show cancellation message
		if !strings.Contains(outputStr, "cancelled") {
			t.Errorf("Expected output to contain cancellation message, got:\n%s", outputStr)
		}
	})

	t.Run("aborts when user presses enter without input", func(t *testing.T) {
		// Mock Nix as installed
		originalChecker := nixInstalledChecker
		nixInstalledChecker = func() bool { return true }
		defer func() { nixInstalledChecker = originalChecker }()

		// Create command with input/output buffers
		var output bytes.Buffer
		var input bytes.Buffer

		// Create a new command instance to avoid state issues
		cmd := &cobra.Command{
			Use:   nukeCmd.Use,
			Short: nukeCmd.Short,
			Long:  nukeCmd.Long,
			RunE:  nukeCmd.RunE,
		}
		cmd.SetOut(&output)
		cmd.SetIn(&input)
		cmd.SetArgs([]string{})

		// Simulate user pressing enter without typing anything
		input.WriteString("\n")

		// Execute command
		err := cmd.Execute()

		// Should not return error when cancelled
		if err != nil {
			t.Errorf("Expected no error when user cancels, got: %v", err)
		}

		outputStr := output.String()

		// Should show cancellation message
		if !strings.Contains(outputStr, "cancelled") {
			t.Errorf("Expected output to contain cancellation message, got:\n%s", outputStr)
		}
	})
}

func TestNukeCommandOutputFormat(t *testing.T) {
	t.Run("displays warning and items to be deleted", func(t *testing.T) {
		// Mock Nix as installed
		originalChecker := nixInstalledChecker
		nixInstalledChecker = func() bool { return true }
		defer func() { nixInstalledChecker = originalChecker }()

		// Create command with input/output buffers
		var output bytes.Buffer
		var input bytes.Buffer

		// Create a new command instance to avoid state issues
		cmd := &cobra.Command{
			Use:   nukeCmd.Use,
			Short: nukeCmd.Short,
			Long:  nukeCmd.Long,
			RunE:  nukeCmd.RunE,
		}
		cmd.SetOut(&output)
		cmd.SetIn(&input)
		cmd.SetArgs([]string{})

		// Simulate user typing 'n' to cancel (so we don't actually run nuke)
		input.WriteString("n\n")

		// Execute command
		_ = cmd.Execute()

		outputStr := output.String()

		// Verify output format contains key information
		expectedPhrases := []string{
			"WARNING",
			"Nix package manager",
			"~/.camp/",
			"home-manager",
			"Are you sure",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(outputStr, phrase) {
				t.Errorf("Expected output to contain '%s', got:\n%s", phrase, outputStr)
			}
		}
	})
}

func TestIsNixInstalled(t *testing.T) {
	t.Run("returns bool value", func(t *testing.T) {
		// Just verify the function returns a bool without panicking
		// We can't test the actual value since it depends on the system
		result := isNixInstalled()
		_ = result // Use the result to avoid unused variable warning
	})
}
