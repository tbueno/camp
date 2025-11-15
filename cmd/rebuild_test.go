package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRebuildCommand(t *testing.T) {
	t.Run("command name", func(t *testing.T) {
		if rebuildCmd.Use != "rebuild" {
			t.Errorf("Expected command name to be 'rebuild', got '%s'", rebuildCmd.Use)
		}
	})

	t.Run("command description", func(t *testing.T) {
		expected := "Rebuild the development environment"
		if rebuildCmd.Short != expected {
			t.Errorf("Expected short description to be '%s', got '%s'", expected, rebuildCmd.Short)
		}
	})

	t.Run("command has long description", func(t *testing.T) {
		if rebuildCmd.Long == "" {
			t.Error("Expected long description to be set")
		}

		if !strings.Contains(rebuildCmd.Long, "Prerequisites") {
			t.Error("Expected long description to contain prerequisites")
		}
	})

	t.Run("command is subcommand of env", func(t *testing.T) {
		found := false
		for _, cmd := range envCmd.Commands() {
			if cmd.Use == "rebuild" {
				found = true
				break
			}
		}
		if !found {
			t.Error("rebuild command should be registered as subcommand of env")
		}
	})
}

func TestRebuildCommandExecution(t *testing.T) {
	// Skip if template files don't exist
	if _, err := os.Stat("templates/files"); os.IsNotExist(err) {
		t.Skip("Skipping test: templates/files directory not found")
	}

	t.Run("execution with missing nix tools", func(t *testing.T) {
		// Create temporary home directory
		tmpHome := t.TempDir()

		// Create .camp directory and config
		campDir := filepath.Join(tmpHome, ".camp")
		if err := os.MkdirAll(campDir, 0755); err != nil {
			t.Fatalf("Failed to create .camp directory: %v", err)
		}

		configPath := filepath.Join(campDir, "camp.yml")
		configContent := `env:
  TEST_VAR: test_value
`
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		// Save original HOME and restore after test
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpHome)

		// Create command with output buffer
		var output bytes.Buffer
		var errOutput bytes.Buffer

		cmd := &cobra.Command{
			Use:   rebuildCmd.Use,
			Short: rebuildCmd.Short,
			Long:  rebuildCmd.Long,
			RunE:  rebuildCmd.RunE,
		}

		cmd.SetOut(&output)
		cmd.SetErr(&errOutput)
		cmd.SetArgs([]string{})

		// Execute command - will likely fail due to missing nix tools
		err := cmd.Execute()

		// Verify that preparation succeeded
		outputStr := output.String()
		if !strings.Contains(outputStr, "Preparing environment") {
			t.Error("Expected output to contain 'Preparing environment'")
		}

		if !strings.Contains(outputStr, "✓ Environment prepared successfully") {
			t.Error("Expected output to show environment was prepared")
		}

		// The rebuild will likely fail unless nix is actually installed
		// That's expected in test environment
		if err != nil {
			t.Logf("Rebuild failed as expected (nix tools not installed): %v", err)
		}

		// Verify files were created during preparation
		nixDir := filepath.Join(tmpHome, ".camp", "nix")
		if _, err := os.Stat(nixDir); os.IsNotExist(err) {
			t.Error("Expected .camp/nix directory to be created")
		}

		flakePath := filepath.Join(nixDir, "flake.nix")
		if _, err := os.Stat(flakePath); os.IsNotExist(err) {
			t.Error("Expected flake.nix to be created")
		}
	})
}

func TestRebuildCommandOutputFormat(t *testing.T) {
	// Skip if template files don't exist
	if _, err := os.Stat("templates/files"); os.IsNotExist(err) {
		t.Skip("Skipping test: templates/files directory not found")
	}

	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create .camp directory and config
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	configPath := filepath.Join(campDir, "camp.yml")
	configContent := `env:
  EDITOR: nvim
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpHome)

	// Create command with output buffer
	var output bytes.Buffer

	cmd := &cobra.Command{
		Use:   rebuildCmd.Use,
		Short: rebuildCmd.Short,
		Long:  rebuildCmd.Long,
		RunE:  rebuildCmd.RunE,
	}

	cmd.SetOut(&output)
	cmd.SetArgs([]string{})

	// Execute command
	_ = cmd.Execute() // Ignore error as rebuild will likely fail

	// Verify output format
	outputStr := output.String()

	expectedPhrases := []string{
		"Starting environment rebuild",
		"Platform:",
		"User:",
		"Hostname:",
		"Preparing environment",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(outputStr, phrase) {
			t.Errorf("Expected output to contain '%s', got:\n%s", phrase, outputStr)
		}
	}
}
