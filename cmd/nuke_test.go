package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"camp/internal/utils/utilstest"

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

// stubNuke skips the confirmation prompt and reports nix as installed or not
func stubNuke(t *testing.T, nixInstalled bool) {
	t.Helper()

	originalChecker := nixInstalledChecker
	originalSkip := skipConfirmation
	nixInstalledChecker = func() bool { return nixInstalled }
	skipConfirmation = true
	t.Cleanup(func() {
		nixInstalledChecker = originalChecker
		skipConfirmation = originalSkip
	})
}

func TestNukeCommandExecution(t *testing.T) {
	t.Run("uninstalls nix and removes camp files", func(t *testing.T) {
		exec := utilstest.StubExec(t)
		stubNuke(t, true)
		user := stubCurrentUser(t, "linux", "env:\n  EDITOR: nvim\n")

		output, err := executeCommand(nukeCmd)
		if err != nil {
			t.Fatalf("nuke returned error: %v", err)
		}

		if !exec.Called("/nix/nix-installer", "uninstall") {
			t.Errorf("Expected nix-installer uninstall to run, got:\n%s", exec)
		}
		if _, err := os.Stat(filepath.Join(user.HomeDir, ".camp")); !os.IsNotExist(err) {
			t.Error("Expected ~/.camp to be removed")
		}
		if !strings.Contains(output, "Your camp environment was erased") {
			t.Errorf("Expected success message, got:\n%s", output)
		}
	})

	t.Run("cleans up leftover files when nix is not installed", func(t *testing.T) {
		exec := utilstest.StubExec(t)
		stubNuke(t, false)
		user := stubCurrentUser(t, "linux", "env:\n  EDITOR: nvim\n")

		output, err := executeCommand(nukeCmd)
		if err != nil {
			t.Fatalf("nuke should succeed without nix, got: %v", err)
		}

		if len(exec.Calls) != 0 {
			t.Errorf("Expected no uninstall commands without nix, ran:\n%s", exec)
		}
		if _, err := os.Stat(filepath.Join(user.HomeDir, ".camp")); !os.IsNotExist(err) {
			t.Error("Expected ~/.camp to be removed")
		}
		if !strings.Contains(output, "Nix is not installed, skipping Nix uninstall") {
			t.Errorf("Expected output to mention skipping the uninstall, got:\n%s", output)
		}
	})
}

func TestNukeCommandConfirmation_OmitsNixWhenNotInstalled(t *testing.T) {
	originalChecker := nixInstalledChecker
	nixInstalledChecker = func() bool { return false }
	t.Cleanup(func() { nixInstalledChecker = originalChecker })

	var output bytes.Buffer
	cmd := &cobra.Command{Use: nukeCmd.Use, RunE: nukeCmd.RunE}
	cmd.SetOut(&output)
	cmd.SetIn(strings.NewReader("n\n"))
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error when user cancels, got: %v", err)
	}
	if strings.Contains(output.String(), "Nix package manager") {
		t.Errorf("Expected the prompt not to list Nix when it is not installed, got:\n%s", output.String())
	}
}

func TestIsNixInstalled(t *testing.T) {
	t.Run("returns bool value", func(t *testing.T) {
		// Just verify the function returns a bool without panicking
		// We can't test the actual value since it depends on the system
		result := isNixInstalled()
		_ = result // Use the result to avoid unused variable warning
	})
}
