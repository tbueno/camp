package cmd

import (
	"camp/internal/project"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCmd_CreatesConfigFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Run init command
	err := initProjectConfig()
	if err != nil {
		t.Fatalf("initProjectConfig() failed: %v", err)
	}

	// Verify .camp.yml was created
	configPath := filepath.Join(tmpDir, ".camp.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("initProjectConfig() should create .camp.yml")
	}

	// Read and verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read .camp.yml: %v", err)
	}

	contentStr := string(content)

	// Check for expected sections
	if !strings.Contains(contentStr, "# Camp project configuration") {
		t.Error("Config should contain header comment")
	}

	if !strings.Contains(contentStr, "env:") {
		t.Error("Config should contain env section")
	}

	if !strings.Contains(contentStr, "PROJECT_NAME:") {
		t.Error("Config should contain PROJECT_NAME example")
	}

	if !strings.Contains(contentStr, "# Future: packages, flakes, scripts will go here") {
		t.Error("Config should contain future features comment")
	}
}

func TestInitCmd_FailsIfConfigExists(t *testing.T) {
	// Create temporary directory with existing .camp.yml
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create existing .camp.yml
	existingConfig := `env:
  EXISTING: "value"
`
	configPath := filepath.Join(tmpDir, ".camp.yml")
	if err := os.WriteFile(configPath, []byte(existingConfig), 0644); err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	// Try to run init command
	err := initProjectConfig()
	if err == nil {
		t.Error("initProjectConfig() should fail when .camp.yml already exists")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Error message should mention 'already exists', got: %v", err)
	}
}

func TestInitCmd_FailsIfConfigExistsInParent(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create .camp.yml in parent
	parentConfig := `env:
  PARENT: "value"
`
	parentConfigPath := filepath.Join(tmpDir, ".camp.yml")
	if err := os.WriteFile(parentConfigPath, []byte(parentConfig), 0644); err != nil {
		t.Fatalf("Failed to create parent config: %v", err)
	}

	// Change to subdirectory
	oldWd, _ := os.Getwd()
	os.Chdir(subDir)
	defer os.Chdir(oldWd)

	// Try to run init command
	err := initProjectConfig()
	if err == nil {
		t.Error("initProjectConfig() should fail when .camp.yml exists in parent")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Error message should mention 'already exists', got: %v", err)
	}
}

func TestInitCmd_FilePermissions(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Run init command
	err := initProjectConfig()
	if err != nil {
		t.Fatalf("initProjectConfig() failed: %v", err)
	}

	// Check file permissions
	configPath := filepath.Join(tmpDir, ".camp.yml")
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat .camp.yml: %v", err)
	}

	// Should be readable and writable by user
	mode := info.Mode()
	if mode.Perm()&0600 != 0600 {
		t.Errorf("Expected file to be readable and writable, got permissions: %o", mode.Perm())
	}
}

func TestInitCmd_ValidYAML(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Run init command
	err := initProjectConfig()
	if err != nil {
		t.Fatalf("initProjectConfig() failed: %v", err)
	}

	// Try to load the created config to verify it's valid YAML
	proj := project.NewProject(tmpDir)
	if !proj.HasCampConfig() {
		t.Error("Created config should be valid and loadable")
	}

	// Verify the template values
	envVars := proj.EnvVars()
	if envVars["PROJECT_NAME"] != "my-project" {
		t.Errorf("Expected PROJECT_NAME=my-project, got %s", envVars["PROJECT_NAME"])
	}
}
