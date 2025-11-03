package system

import (
	"bytes"
	"strings"
	"testing"
)

func TestGetDefaultBootstrapConfig(t *testing.T) {
	config := GetDefaultBootstrapConfig()

	if config == nil {
		t.Fatal("GetDefaultBootstrapConfig returned nil")
	}

	if len(config.Applications) == 0 {
		t.Fatal("Expected at least one application in default config")
	}

	// Check that all applications have required fields
	for i, app := range config.Applications {
		if app.Name == "" {
			t.Errorf("Application at index %d has empty name", i)
		}
		if app.InstallCommand == "" {
			t.Errorf("Application %s has empty install command", app.Name)
		}
	}

	foundDirenv := false
	foundNix := false
	foundDevbox := false
	for _, app := range config.Applications {
		if app.Name == "direnv" {
			foundDirenv = true
		}
		if app.Name == "nix" {
			foundNix = true
		}
		if app.Name == "devbox" {
			foundDevbox = true
		}
	}

	if !foundDirenv {
		t.Error("Expected direnv to be in default applications")
	}

	if !foundNix {
		t.Error("Expected nix to be in default applications")
	}

	if !foundDevbox {
		t.Error("Expected devbox to be in default applications")
	}
}

func TestRunBootstrap_DryRun(t *testing.T) {
	config := &BootstrapConfig{
		Applications: []Application{
			{Name: "test-app1", InstallCommand: "echo 'installing app1'"},
			{Name: "test-app2", InstallCommand: "echo 'installing app2'"},
		},
	}

	var output bytes.Buffer
	err := RunBootstrap(config, &output, true)

	if err != nil {
		t.Fatalf("RunBootstrap in dry-run mode failed: %v", err)
	}

	outputStr := output.String()

	if !strings.Contains(outputStr, "[DRY RUN]") {
		t.Error("Expected [DRY RUN] message in output")
	}

	if !strings.Contains(outputStr, "test-app1") {
		t.Error("Expected test-app1 to be mentioned in output")
	}
	if !strings.Contains(outputStr, "test-app2") {
		t.Error("Expected test-app2 to be mentioned in output")
	}

	if !strings.Contains(outputStr, "[1/2]") || !strings.Contains(outputStr, "[2/2]") {
		t.Error("Expected progress indicators in output")
	}

	if !strings.Contains(outputStr, "🎉 Bootstrap process completed successfully!") {
		t.Error("Expected completion message in output")
	}
}

func TestRunBootstrap_WithNilConfig(t *testing.T) {
	var output bytes.Buffer
	err := RunBootstrap(nil, &output, true)

	if err != nil {
		t.Fatalf("RunBootstrap with nil config failed: %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "Starting bootstrap process") {
		t.Error("Expected bootstrap process to start with default config")
	}
}

func TestRunBootstrap_SuccessfulInstallation(t *testing.T) {
	config := &BootstrapConfig{
		Applications: []Application{
			{Name: "test-success", InstallCommand: "echo 'success'"},
		},
	}

	var output bytes.Buffer
	err := RunBootstrap(config, &output, false)

	if err != nil {
		t.Fatalf("RunBootstrap failed: %v", err)
	}

	outputStr := output.String()

	if !strings.Contains(outputStr, "✅ test-success installation completed") {
		t.Error("Expected success checkmark and completion message")
	}

	if !strings.Contains(outputStr, "success") {
		t.Error("Expected command output to be included")
	}
}

func TestRunBootstrap_FailedInstallation(t *testing.T) {
	config := &BootstrapConfig{
		Applications: []Application{
			{Name: "test-fail", InstallCommand: "exit 1"},
		},
	}

	var output bytes.Buffer
	err := RunBootstrap(config, &output, false)

	if err == nil {
		t.Fatal("Expected RunBootstrap to fail with failing command")
	}

	outputStr := output.String()

	if !strings.Contains(outputStr, "❌ Failed to install test-fail") {
		t.Error("Expected failure indicator in output")
	}

	if !strings.Contains(err.Error(), "test-fail") {
		t.Error("Expected error message to contain application name")
	}
}

func TestRunBootstrap_EmptyApplicationsList(t *testing.T) {
	config := &BootstrapConfig{
		Applications: []Application{},
	}

	var output bytes.Buffer
	err := RunBootstrap(config, &output, true)

	if err != nil {
		t.Fatalf("RunBootstrap with empty applications failed: %v", err)
	}

	outputStr := output.String()

	if !strings.Contains(outputStr, "Starting bootstrap process for 0 applications") {
		t.Error("Expected message about 0 applications")
	}

	if !strings.Contains(outputStr, "🎉 Bootstrap process completed successfully!") {
		t.Error("Expected completion message even with 0 applications")
	}
}

func TestExecuteInstallCommand_SimpleCommand(t *testing.T) {
	var output bytes.Buffer
	err := executeInstallCommand("echo 'hello world'", &output)

	if err != nil {
		t.Fatalf("executeInstallCommand failed: %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "hello world") {
		t.Error("Expected command output in result")
	}
}

func TestExecuteInstallCommand_CompoundCommand(t *testing.T) {
	var output bytes.Buffer
	err := executeInstallCommand("echo 'first' && echo 'second'", &output)

	if err != nil {
		t.Fatalf("executeInstallCommand with compound command failed: %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "first") || !strings.Contains(outputStr, "second") {
		t.Error("Expected both parts of compound command in output")
	}
}

func TestExecuteInstallCommand_PipeCommand(t *testing.T) {
	var output bytes.Buffer
	err := executeInstallCommand("echo 'test' | cat", &output)

	if err != nil {
		t.Fatalf("executeInstallCommand with pipe failed: %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "test") {
		t.Error("Expected piped command output in result")
	}
}

func TestExecuteInstallCommand_EmptyCommand(t *testing.T) {
	var output bytes.Buffer
	err := executeInstallCommand("", &output)

	if err == nil {
		t.Fatal("Expected executeInstallCommand to fail with empty command")
	}

	if !strings.Contains(err.Error(), "empty install command") {
		t.Error("Expected specific error message for empty command")
	}
}

func TestExecuteInstallCommand_WhitespaceOnlyCommand(t *testing.T) {
	var output bytes.Buffer
	err := executeInstallCommand("   \t\n   ", &output)

	if err == nil {
		t.Fatal("Expected executeInstallCommand to fail with whitespace-only command")
	}

	if !strings.Contains(err.Error(), "empty install command") {
		t.Error("Expected specific error message for whitespace-only command")
	}
}

func TestExecuteInstallCommand_FailingCommand(t *testing.T) {
	var output bytes.Buffer
	err := executeInstallCommand("exit 42", &output)

	if err == nil {
		t.Fatal("Expected executeInstallCommand to fail with failing command")
	}

	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestIsCommandAvailable_ExistingCommand(t *testing.T) {
	// Test with a command that should exist on all systems
	if !isCommandAvailable("sh") {
		t.Error("Expected 'sh' to be available in PATH")
	}
}

func TestIsCommandAvailable_NonExistingCommand(t *testing.T) {
	// Test with a command that definitely doesn't exist
	if isCommandAvailable("this-command-definitely-does-not-exist-12345") {
		t.Error("Expected non-existing command to return false")
	}
}

func TestRunBootstrap_SkipsAlreadyInstalledCommands(t *testing.T) {
	// Use 'echo' as a command that's guaranteed to exist
	config := &BootstrapConfig{
		Applications: []Application{
			{Name: "echo", InstallCommand: "echo 'should not run'"},
			{Name: "this-does-not-exist-xyz", InstallCommand: "echo 'this would install'"},
		},
	}

	var output bytes.Buffer
	err := RunBootstrap(config, &output, true)

	if err != nil {
		t.Fatalf("RunBootstrap failed: %v", err)
	}

	outputStr := output.String()

	// Should skip echo because it exists
	if !strings.Contains(outputStr, "echo is already installed, skipping") {
		t.Error("Expected message about skipping already installed command")
	}

	// Should attempt to install the non-existing command
	if !strings.Contains(outputStr, "this-does-not-exist-xyz") {
		t.Error("Expected non-existing command to be processed")
	}

	// Should not actually run the install command for echo
	if strings.Contains(outputStr, "should not run") {
		t.Error("Should not have executed install command for already installed application")
	}
}
