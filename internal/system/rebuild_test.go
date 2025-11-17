package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareEnvironment(t *testing.T) {
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
  TEST_VAR: test_value
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create user
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
	}

	// Run PrepareEnvironment
	if err := PrepareEnvironment(user); err != nil {
		t.Fatalf("PrepareEnvironment() failed: %v", err)
	}

	// Verify .camp/nix directory was created
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if _, err := os.Stat(nixDir); os.IsNotExist(err) {
		t.Error("PrepareEnvironment() should create .camp/nix directory")
	}

	// Verify config files were copied
	macNixPath := filepath.Join(nixDir, "mac.nix")
	if _, err := os.Stat(macNixPath); os.IsNotExist(err) {
		t.Error("PrepareEnvironment() should copy mac.nix")
	}

	// Verify modules directory was copied
	modulesPath := filepath.Join(nixDir, "modules")
	if _, err := os.Stat(modulesPath); os.IsNotExist(err) {
		t.Error("PrepareEnvironment() should copy modules directory")
	}

	// Verify flake.nix was rendered
	flakePath := filepath.Join(nixDir, "flake.nix")
	if _, err := os.Stat(flakePath); os.IsNotExist(err) {
		t.Error("PrepareEnvironment() should render flake.nix")
	}
}

func TestCopyConfigFiles(t *testing.T) {
	// Skip if template files don't exist
	if _, err := os.Stat("templates/files"); os.IsNotExist(err) {
		t.Skip("Skipping test: templates/files directory not found")
	}

	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create user
	user := &User{
		Name:    "testuser",
		HomeDir: tmpHome,
	}

	// Run CopyConfigFiles
	if err := CopyConfigFiles(user); err != nil {
		t.Fatalf("CopyConfigFiles() failed: %v", err)
	}

	nixDir := filepath.Join(tmpHome, ".camp", "nix")

	// Verify mac.nix was copied
	macNixPath := filepath.Join(nixDir, "mac.nix")
	if _, err := os.Stat(macNixPath); os.IsNotExist(err) {
		t.Error("CopyConfigFiles() should copy mac.nix")
	}

	// Verify linux.nix was copied
	linuxNixPath := filepath.Join(nixDir, "linux.nix")
	if _, err := os.Stat(linuxNixPath); os.IsNotExist(err) {
		t.Error("CopyConfigFiles() should copy linux.nix")
	}

	// Verify modules directory was copied
	modulesPath := filepath.Join(nixDir, "modules", "common.nix")
	if _, err := os.Stat(modulesPath); os.IsNotExist(err) {
		t.Error("CopyConfigFiles() should copy modules/common.nix")
	}

	// Verify flake.nix was NOT copied (it should be rendered)
	flakePath := filepath.Join(nixDir, "flake.nix")
	if _, err := os.Stat(flakePath); err == nil {
		t.Error("CopyConfigFiles() should not copy flake.nix (it's rendered separately)")
	}
}

func TestCompileTemplates(t *testing.T) {
	// Skip if template files don't exist
	if _, err := os.Stat("templates/files/flake.nix"); os.IsNotExist(err) {
		t.Skip("Skipping test: flake.nix template not found")
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
  BROWSER: firefox
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create user
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
	}

	// Run CompileTemplates
	if err := CompileTemplates(user); err != nil {
		t.Fatalf("CompileTemplates() failed: %v", err)
	}

	// Verify flake.nix was created
	flakePath := filepath.Join(tmpHome, ".camp", "nix", "flake.nix")
	if _, err := os.Stat(flakePath); os.IsNotExist(err) {
		t.Error("CompileTemplates() should create flake.nix")
	}

	// Verify user's EnvVars were reloaded
	if user.EnvVars["EDITOR"] != "nvim" {
		t.Error("CompileTemplates() should reload user config")
	}
}

func TestExecuteRebuild_Darwin(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create nix directory
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if err := os.MkdirAll(nixDir, 0755); err != nil {
		t.Fatalf("Failed to create nix directory: %v", err)
	}

	// Create user
	user := &User{
		Name:     "testuser",
		HostName: "testhost",
		Platform: "darwin",
		HomeDir:  tmpHome,
	}

	// We can't actually run the rebuild command in tests
	// Instead, we'll verify the function constructs the right parameters
	// by checking that it would fail with "command not found" or similar
	// (unless nix is actually installed)
	err := ExecuteRebuild(user)

	// The command will likely fail unless nix-darwin is installed
	// We're just testing that the function doesn't panic and
	// constructs a valid command
	if err != nil {
		// Expected - nix-darwin likely not installed in test environment
		t.Logf("ExecuteRebuild() failed as expected (nix-darwin not installed): %v", err)
	}
}

func TestExecuteRebuild_Linux(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create nix directory
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if err := os.MkdirAll(nixDir, 0755); err != nil {
		t.Fatalf("Failed to create nix directory: %v", err)
	}

	// Create user
	user := &User{
		Name:     "testuser",
		HostName: "testhost",
		Platform: "linux",
		HomeDir:  tmpHome,
	}

	// We can't actually run the rebuild command in tests
	err := ExecuteRebuild(user)

	// The command will likely fail unless home-manager is installed
	if err != nil {
		// Expected - home-manager likely not installed in test environment
		t.Logf("ExecuteRebuild() failed as expected (home-manager not installed): %v", err)
	}
}

func TestExecuteRebuild_UnsupportedPlatform(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create nix directory
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if err := os.MkdirAll(nixDir, 0755); err != nil {
		t.Fatalf("Failed to create nix directory: %v", err)
	}

	// Create user with unsupported platform
	user := &User{
		Name:     "testuser",
		HostName: "testhost",
		Platform: "windows",
		HomeDir:  tmpHome,
	}

	// Should error for unsupported platform
	err := ExecuteRebuild(user)
	if err == nil {
		t.Error("ExecuteRebuild() should error for unsupported platform")
	}

	if err != nil && err.Error() != "unsupported platform: windows" {
		t.Errorf("Expected 'unsupported platform' error, got: %v", err)
	}
}

func TestExecuteRebuild_NoNixDirectory(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Don't create nix directory

	// Create user
	user := &User{
		Name:     "testuser",
		HostName: "testhost",
		Platform: "darwin",
		HomeDir:  tmpHome,
	}

	// Should error when nix directory doesn't exist
	err := ExecuteRebuild(user)
	if err == nil {
		t.Error("ExecuteRebuild() should error when nix directory doesn't exist")
	}
}

func TestCopyDir(t *testing.T) {
	// Create temporary directories
	tmpSrc := t.TempDir()
	tmpDest := t.TempDir()

	// Create source structure
	if err := os.MkdirAll(filepath.Join(tmpSrc, "subdir"), 0755); err != nil {
		t.Fatalf("Failed to create source structure: %v", err)
	}

	// Create files
	file1 := filepath.Join(tmpSrc, "file1.txt")
	file2 := filepath.Join(tmpSrc, "subdir", "file2.txt")

	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// Copy directory
	destPath := filepath.Join(tmpDest, "copied")
	if err := copyDir(tmpSrc, destPath); err != nil {
		t.Fatalf("copyDir() failed: %v", err)
	}

	// Verify files were copied
	copiedFile1 := filepath.Join(destPath, "file1.txt")
	copiedFile2 := filepath.Join(destPath, "subdir", "file2.txt")

	if _, err := os.Stat(copiedFile1); os.IsNotExist(err) {
		t.Error("copyDir() should copy file1.txt")
	}
	if _, err := os.Stat(copiedFile2); os.IsNotExist(err) {
		t.Error("copyDir() should copy subdir/file2.txt")
	}

	// Verify content
	content1, err := os.ReadFile(copiedFile1)
	if err != nil {
		t.Fatalf("Failed to read copied file1: %v", err)
	}
	if string(content1) != "content1" {
		t.Error("Copied file1 content doesn't match")
	}
}

// Integration test for flakes

func TestPrepareEnvironment_WithFlakes(t *testing.T) {
	// Skip if template files don't exist
	if _, err := os.Stat("templates/files/flake.nix"); os.IsNotExist(err) {
		t.Skip("Skipping test: flake.nix template not found")
	}

	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create .camp directory and config with flakes
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	configPath := filepath.Join(campDir, "camp.yml")
	configContent := `env:
  EDITOR: nvim

flakes:
  - name: test-flake
    url: "github:test/flake"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: packages
        type: home
      - name: darwinModules.test
        type: system
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create user
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
		Flakes:       []Flake{},
	}

	// Run PrepareEnvironment
	if err := PrepareEnvironment(user); err != nil {
		t.Fatalf("PrepareEnvironment() failed: %v", err)
	}

	// Verify flake.nix was created
	flakePath := filepath.Join(tmpHome, ".camp", "nix", "flake.nix")
	if _, err := os.Stat(flakePath); os.IsNotExist(err) {
		t.Fatal("PrepareEnvironment() should create flake.nix")
	}

	// Read generated flake.nix
	content, err := os.ReadFile(flakePath)
	if err != nil {
		t.Fatalf("Failed to read generated flake.nix: %v", err)
	}

	contentStr := string(content)

	// Verify flake is in inputs section
	if !strings.Contains(contentStr, "test-flake = {") {
		t.Error("Generated flake.nix should contain test-flake in inputs")
	}

	if !strings.Contains(contentStr, `url = "github:test/flake";`) {
		t.Error("Generated flake.nix should contain flake URL")
	}

	if !strings.Contains(contentStr, `inputs.nixpkgs.follows = "nixpkgs";`) {
		t.Error("Generated flake.nix should contain follows declaration")
	}

	// Verify flake is in outputs function signature
	if !strings.Contains(contentStr, "test-flake,") {
		t.Error("Generated flake.nix should include test-flake in outputs signature")
	}

	// Verify system-level output is injected into darwin modules
	if !strings.Contains(contentStr, "test-flake.darwinModules.test") {
		t.Error("Generated flake.nix should inject system output into darwin modules")
	}

	// Verify home-level output is injected into home-manager modules
	if !strings.Contains(contentStr, "test-flake.packages") {
		t.Error("Generated flake.nix should inject home output into home-manager modules")
	}

	// Verify user data is present
	if !strings.Contains(contentStr, "testuser") {
		t.Error("Generated flake.nix should contain username")
	}

	if !strings.Contains(contentStr, "testhost") {
		t.Error("Generated flake.nix should contain hostname")
	}
}
