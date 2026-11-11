package system

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"camp/templates"
)

func TestCopyBinFiles_CopiesScriptsAsExecutable(t *testing.T) {
	campPath := t.TempDir()
	if err := os.MkdirAll(filepath.Join(campPath, "bin"), 0755); err != nil {
		t.Fatal(err)
	}
	templDir := fstest.MapFS{
		"initial/bin/bootstrap":   {Data: []byte("#!/bin/sh\necho bootstrap\n")},
		"initial/bin/install_nix": {Data: []byte("#!/bin/sh\necho install\n")},
	}

	if err := copyBinFiles(templDir, campPath, &bytes.Buffer{}, false); err != nil {
		t.Fatalf("copyBinFiles failed: %v", err)
	}

	for _, name := range []string{"bootstrap", "install_nix"} {
		path := filepath.Join(campPath, "bin", name)
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("expected %s to be copied: %v", name, err)
		}
		if info.Mode().Perm() != 0755 {
			t.Errorf("expected %s to have mode 0755, got %v", name, info.Mode().Perm())
		}
		content, _ := os.ReadFile(path)
		if string(content) != string(templDir["initial/bin/"+name].Data) {
			t.Errorf("unexpected content for %s: %q", name, content)
		}
	}
}

func TestCopyBinFiles_DryRunWritesNothing(t *testing.T) {
	campPath := t.TempDir()
	templDir := fstest.MapFS{
		"initial/bin/bootstrap": {Data: []byte("#!/bin/sh\n")},
	}
	var output bytes.Buffer

	if err := copyBinFiles(templDir, campPath, &output, true); err != nil {
		t.Fatalf("copyBinFiles failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(campPath, "bin", "bootstrap")); !os.IsNotExist(err) {
		t.Error("expected no file to be written in dry-run mode")
	}
	if !strings.Contains(output.String(), "[DRY RUN] Would copy bin file") {
		t.Errorf("expected dry-run message, got: %s", output.String())
	}
}

func TestCopyBinFiles_FromEmbeddedTemplates(t *testing.T) {
	campPath := t.TempDir()
	if err := os.MkdirAll(filepath.Join(campPath, "bin"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := copyBinFiles(templates.FS, campPath, &bytes.Buffer{}, false); err != nil {
		t.Fatalf("copyBinFiles failed: %v", err)
	}

	for _, name := range []string{"bootstrap", "install_nix"} {
		if _, err := os.Stat(filepath.Join(campPath, "bin", name)); err != nil {
			t.Errorf("expected %s to be copied from embedded templates: %v", name, err)
		}
	}
}

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

	foundNix := false
	for _, app := range config.Applications {
		if app.Name == "nix" {
			foundNix = true
		}
	}

	if !foundNix {
		t.Error("Expected nix to be in default applications")
	}

	// Verify that only nix is included (direnv and devbox are installed via flake.nix)
	if len(config.Applications) != 1 {
		t.Errorf("Expected exactly 1 application (nix only), got %d", len(config.Applications))
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

func TestCreateEnvrc_WritesWhenMissing(t *testing.T) {
	home := t.TempDir()
	campPath := filepath.Join(home, ".camp")
	var out bytes.Buffer

	if err := createEnvrc(&User{HomeDir: home}, campPath, &out, false); err != nil {
		t.Fatalf("createEnvrc() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(home, ".envrc"))
	if err != nil {
		t.Fatalf("expected .envrc to be created: %v", err)
	}
	want := envrcMarker + "\nuse flake " + campPath + "/nix\n"
	if string(content) != want {
		t.Errorf(".envrc = %q, want %q", content, want)
	}
}

func TestCreateEnvrc_RewritesCampEnvrc(t *testing.T) {
	home := t.TempDir()
	campPath := filepath.Join(home, ".camp")
	envrcPath := filepath.Join(home, ".envrc")
	if err := os.WriteFile(envrcPath, []byte(envrcMarker+"\nuse flake /old/path\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer

	if err := createEnvrc(&User{HomeDir: home}, campPath, &out, false); err != nil {
		t.Fatalf("createEnvrc() error = %v", err)
	}

	content, _ := os.ReadFile(envrcPath)
	if !strings.Contains(string(content), "use flake "+campPath+"/nix") {
		t.Errorf("expected camp .envrc to be rewritten, got %q", content)
	}
}

func TestCreateEnvrc_KeepsForeignEnvrc(t *testing.T) {
	home := t.TempDir()
	campPath := filepath.Join(home, ".camp")
	envrcPath := filepath.Join(home, ".envrc")
	userContent := "export FOO=bar\n"
	if err := os.WriteFile(envrcPath, []byte(userContent), 0644); err != nil {
		t.Fatal(err)
	}

	for _, dryRun := range []bool{false, true} {
		var out bytes.Buffer

		if err := createEnvrc(&User{HomeDir: home}, campPath, &out, dryRun); err != nil {
			t.Fatalf("createEnvrc(dryRun=%v) error = %v", dryRun, err)
		}

		content, _ := os.ReadFile(envrcPath)
		if string(content) != userContent {
			t.Errorf("dryRun=%v: expected .envrc to be kept, got %q", dryRun, content)
		}
		if !strings.Contains(out.String(), "use flake "+campPath+"/nix") {
			t.Errorf("dryRun=%v: expected output to show the use flake line, got %q", dryRun, out.String())
		}
	}
}

func TestCreateEnvrc_DryRunWritesNothing(t *testing.T) {
	home := t.TempDir()
	var out bytes.Buffer

	if err := createEnvrc(&User{HomeDir: home}, filepath.Join(home, ".camp"), &out, true); err != nil {
		t.Fatalf("createEnvrc() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(home, ".envrc")); !os.IsNotExist(err) {
		t.Error("expected dry run not to create .envrc")
	}
	if !strings.Contains(out.String(), "[DRY RUN]") {
		t.Errorf("expected dry run output, got %q", out.String())
	}
}

func TestBootstrapPlatforms_DryRunReportsEnvrc(t *testing.T) {
	platforms := map[string]func(string, *User, *bytes.Buffer) error{
		"mac": func(campPath string, user *User, out *bytes.Buffer) error {
			return bootstrapMac(campPath, templates.FS, user, out, true)
		},
		"linux": func(campPath string, user *User, out *bytes.Buffer) error {
			return bootstrapLinux(campPath, templates.FS, user, out, true)
		},
	}

	for name, bootstrap := range platforms {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			campPath := filepath.Join(home, ".camp")
			user := &User{Name: "tester", HomeDir: home, Architecture: "arm64"}
			var out bytes.Buffer

			if err := bootstrap(campPath, user, &out); err != nil {
				t.Fatalf("dry run error = %v", err)
			}

			if !strings.Contains(out.String(), "[DRY RUN] Would create .envrc file: "+filepath.Join(home, ".envrc")) {
				t.Errorf("expected dry run to report the .envrc step, got:\n%s", out.String())
			}
			entries, err := os.ReadDir(home)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Errorf("expected dry run to write nothing, found %d entries in home", len(entries))
			}
		})
	}
}
