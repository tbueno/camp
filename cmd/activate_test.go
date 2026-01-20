package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"camp/internal/project"
)

func TestActivateCmd(t *testing.T) {
	t.Run("with valid camp.yaml", func(t *testing.T) {
		// Create temp directory with camp.yaml
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		yamlContent := `env:
  DATABASE_URL: "postgres://localhost:5432/db"
  DEBUG: "true"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		// Change to temp directory
		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		// Set shell to bash for predictable output
		origShell := os.Getenv("SHELL")
		os.Setenv("SHELL", "/bin/bash")
		defer os.Setenv("SHELL", origShell)

		// Execute command
		var output bytes.Buffer
		cmd := activateCmd()
		cmd.SetOut(&output)
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		result := output.String()
		if !strings.Contains(result, `export DATABASE_URL="postgres://localhost:5432/db";`) {
			t.Errorf("Expected DATABASE_URL export, got: %s", result)
		}
		if !strings.Contains(result, `export DEBUG="true";`) {
			t.Errorf("Expected DEBUG export, got: %s", result)
		}
	})

	t.Run("without camp.yaml", func(t *testing.T) {
		tmpDir := t.TempDir()

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		var output bytes.Buffer
		cmd := activateCmd()
		cmd.SetOut(&output)

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command should not fail without config: %v", err)
		}

		if output.String() != "" {
			t.Errorf("Expected empty output, got: %s", output.String())
		}
	})

	t.Run("with empty env section", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		yamlContent := `env: {}
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		var output bytes.Buffer
		cmd := activateCmd()
		cmd.SetOut(&output)

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command should not fail with empty env: %v", err)
		}

		if output.String() != "" {
			t.Errorf("Expected empty output for empty env, got: %s", output.String())
		}
	})

	t.Run("with shell override to fish", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		yamlContent := `env:
  FOO: "bar"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		var output bytes.Buffer
		cmd := activateCmd()
		cmd.SetOut(&output)
		cmd.SetArgs([]string{"--shell", "fish"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		if !strings.Contains(output.String(), "set -gx FOO 'bar';") {
			t.Errorf("Expected fish syntax, got: %s", output.String())
		}
	})

	t.Run("output is sorted alphabetically", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		yamlContent := `env:
  ZEBRA: "last"
  APPLE: "first"
  MIDDLE: "middle"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		origShell := os.Getenv("SHELL")
		os.Setenv("SHELL", "/bin/bash")
		defer os.Setenv("SHELL", origShell)

		var output bytes.Buffer
		cmd := activateCmd()
		cmd.SetOut(&output)

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(output.String()), "\n")
		if len(lines) != 3 {
			t.Fatalf("Expected 3 lines, got %d", len(lines))
		}

		// Check alphabetical order
		if !strings.HasPrefix(lines[0], "export APPLE=") {
			t.Errorf("Expected APPLE first, got: %s", lines[0])
		}
		if !strings.HasPrefix(lines[1], "export MIDDLE=") {
			t.Errorf("Expected MIDDLE second, got: %s", lines[1])
		}
		if !strings.HasPrefix(lines[2], "export ZEBRA=") {
			t.Errorf("Expected ZEBRA third, got: %s", lines[2])
		}
	})

	t.Run("with special characters in values", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		yamlContent := `env:
  QUOTED: 'hello "world"'
  DOLLAR: "$HOME/bin"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		origShell := os.Getenv("SHELL")
		os.Setenv("SHELL", "/bin/bash")
		defer os.Setenv("SHELL", origShell)

		var output bytes.Buffer
		cmd := activateCmd()
		cmd.SetOut(&output)

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		result := output.String()
		// Check that special characters are properly escaped
		if !strings.Contains(result, `\$HOME`) {
			t.Errorf("Expected escaped dollar sign, got: %s", result)
		}
		if !strings.Contains(result, `\"world\"`) {
			t.Errorf("Expected escaped quotes, got: %s", result)
		}
	})

	t.Run("with invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		invalidYaml := `env:
  invalid yaml content
    - bad indentation
`
		if err := os.WriteFile(configPath, []byte(invalidYaml), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		var output bytes.Buffer
		cmd := activateCmd()
		cmd.SetOut(&output)
		cmd.SetErr(&output)

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}
	})
}

func TestActivateCmdMetadata(t *testing.T) {
	cmd := activateCmd()

	t.Run("command name", func(t *testing.T) {
		if cmd.Use != "activate" {
			t.Errorf("Expected command name to be 'activate', got '%s'", cmd.Use)
		}
	})

	t.Run("has short description", func(t *testing.T) {
		if cmd.Short == "" {
			t.Error("Expected non-empty short description")
		}
	})

	t.Run("has long description", func(t *testing.T) {
		if cmd.Long == "" {
			t.Error("Expected non-empty long description")
		}
	})

	t.Run("has shell flag", func(t *testing.T) {
		flag := cmd.Flags().Lookup("shell")
		if flag == nil {
			t.Error("Expected --shell flag to be defined")
		}
	})
}

func TestActivateCmdTracking(t *testing.T) {
	t.Run("creates tracking file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		yamlContent := `env:
  FOO: "bar"
  BAZ: "qux"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		origShell := os.Getenv("SHELL")
		os.Setenv("SHELL", "/bin/bash")
		defer os.Setenv("SHELL", origShell)

		var output bytes.Buffer
		cmd := activateCmd()
		cmd.SetOut(&output)

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		// Verify tracking file was created
		trackingPath := project.GetTrackingFilePath(tmpDir)
		if _, err := os.Stat(trackingPath); os.IsNotExist(err) {
			t.Error("Tracking file was not created")
		}

		// Verify tracked vars
		tracked, err := project.LoadTrackedVars(tmpDir)
		if err != nil {
			t.Fatalf("Failed to load tracked vars: %v", err)
		}

		if len(tracked) != 2 {
			t.Errorf("Expected 2 tracked vars, got %d", len(tracked))
		}
	})

	t.Run("unsets removed vars", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		// First activation with FOO and BAR
		yamlContent := `env:
  FOO: "foo_value"
  BAR: "bar_value"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		origShell := os.Getenv("SHELL")
		os.Setenv("SHELL", "/bin/bash")
		defer os.Setenv("SHELL", origShell)

		// First run
		var output1 bytes.Buffer
		cmd1 := activateCmd()
		cmd1.SetOut(&output1)
		if err := cmd1.Execute(); err != nil {
			t.Fatalf("First command failed: %v", err)
		}

		// Update config - remove BAR, keep FOO
		yamlContent2 := `env:
  FOO: "foo_value"
`
		if err := os.WriteFile(configPath, []byte(yamlContent2), 0644); err != nil {
			t.Fatalf("Failed to write updated config: %v", err)
		}

		// Second run
		var output2 bytes.Buffer
		cmd2 := activateCmd()
		cmd2.SetOut(&output2)
		if err := cmd2.Execute(); err != nil {
			t.Fatalf("Second command failed: %v", err)
		}

		result := output2.String()
		// Should contain unset for BAR
		if !strings.Contains(result, "unset BAR;") {
			t.Errorf("Expected unset BAR, got: %s", result)
		}
		// Should still export FOO
		if !strings.Contains(result, "export FOO=") {
			t.Errorf("Expected export FOO, got: %s", result)
		}
	})

	t.Run("unsets all vars when config removed", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		// First activation
		yamlContent := `env:
  VAR1: "value1"
  VAR2: "value2"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		origShell := os.Getenv("SHELL")
		os.Setenv("SHELL", "/bin/bash")
		defer os.Setenv("SHELL", origShell)

		// First run
		var output1 bytes.Buffer
		cmd1 := activateCmd()
		cmd1.SetOut(&output1)
		if err := cmd1.Execute(); err != nil {
			t.Fatalf("First command failed: %v", err)
		}

		// Remove config file
		os.Remove(configPath)

		// Second run
		var output2 bytes.Buffer
		cmd2 := activateCmd()
		cmd2.SetOut(&output2)
		if err := cmd2.Execute(); err != nil {
			t.Fatalf("Second command failed: %v", err)
		}

		result := output2.String()
		// Should contain unset for both vars
		if !strings.Contains(result, "unset VAR1;") {
			t.Errorf("Expected unset VAR1, got: %s", result)
		}
		if !strings.Contains(result, "unset VAR2;") {
			t.Errorf("Expected unset VAR2, got: %s", result)
		}
		// Should not contain any exports
		if strings.Contains(result, "export") {
			t.Errorf("Expected no exports, got: %s", result)
		}
	})

	t.Run("fish shell unset syntax", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		// First activation
		yamlContent := `env:
  OLD_VAR: "old"
  KEEP_VAR: "keep"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		original, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(original)

		// First run with fish
		var output1 bytes.Buffer
		cmd1 := activateCmd()
		cmd1.SetOut(&output1)
		cmd1.SetArgs([]string{"--shell", "fish"})
		if err := cmd1.Execute(); err != nil {
			t.Fatalf("First command failed: %v", err)
		}

		// Update config - remove OLD_VAR
		yamlContent2 := `env:
  KEEP_VAR: "keep"
`
		if err := os.WriteFile(configPath, []byte(yamlContent2), 0644); err != nil {
			t.Fatalf("Failed to write updated config: %v", err)
		}

		// Second run with fish
		var output2 bytes.Buffer
		cmd2 := activateCmd()
		cmd2.SetOut(&output2)
		cmd2.SetArgs([]string{"--shell", "fish"})
		if err := cmd2.Execute(); err != nil {
			t.Fatalf("Second command failed: %v", err)
		}

		result := output2.String()
		// Should use fish unset syntax
		if !strings.Contains(result, "set -e OLD_VAR;") {
			t.Errorf("Expected fish unset syntax, got: %s", result)
		}
	})
}
