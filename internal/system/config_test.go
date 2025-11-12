package system

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if config.Env == nil {
		t.Error("DefaultConfig() should initialize Env map")
	}

	if len(config.Env) != 0 {
		t.Errorf("DefaultConfig() should have empty Env map, got %d entries", len(config.Env))
	}
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	config, err := LoadConfig("/nonexistent/path/camp.yml")

	if err != nil {
		t.Errorf("LoadConfig() should not error on non-existent file, got: %v", err)
	}

	if config == nil {
		t.Fatal("LoadConfig() should return default config when file doesn't exist")
	}

	if config.Env == nil {
		t.Error("LoadConfig() should initialize Env map")
	}
}

func TestLoadConfig_ValidYAML(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write valid YAML
	yamlContent := `env:
  EDITOR: nvim
  BROWSER: firefox
  CUSTOM_VAR: test_value
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify env vars
	if len(config.Env) != 3 {
		t.Errorf("Expected 3 env vars, got %d", len(config.Env))
	}

	if config.Env["EDITOR"] != "nvim" {
		t.Errorf("Expected EDITOR=nvim, got %s", config.Env["EDITOR"])
	}

	if config.Env["BROWSER"] != "firefox" {
		t.Errorf("Expected BROWSER=firefox, got %s", config.Env["BROWSER"])
	}

	if config.Env["CUSTOM_VAR"] != "test_value" {
		t.Errorf("Expected CUSTOM_VAR=test_value, got %s", config.Env["CUSTOM_VAR"])
	}
}

func TestLoadConfig_EmptyYAML(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write empty YAML
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify Env map is initialized
	if config.Env == nil {
		t.Error("LoadConfig() should initialize Env map even for empty YAML")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write invalid YAML
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config - should return error
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig() should return error for invalid YAML")
	}
}

func TestLoadUserConfig_YmlFile(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	// Write .yml file
	configPath := filepath.Join(campDir, "camp.yml")
	yamlContent := `env:
  TEST: value
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := LoadUserConfig(tmpHome)
	if err != nil {
		t.Fatalf("LoadUserConfig() failed: %v", err)
	}

	if config.Env["TEST"] != "value" {
		t.Errorf("Expected TEST=value, got %s", config.Env["TEST"])
	}
}

func TestLoadUserConfig_YamlFile(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	// Write .yaml file (not .yml)
	configPath := filepath.Join(campDir, "camp.yaml")
	yamlContent := `env:
  TEST: value
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := LoadUserConfig(tmpHome)
	if err != nil {
		t.Fatalf("LoadUserConfig() failed: %v", err)
	}

	if config.Env["TEST"] != "value" {
		t.Errorf("Expected TEST=value, got %s", config.Env["TEST"])
	}
}

func TestLoadUserConfig_NoFile(t *testing.T) {
	// Create temporary home directory without config
	tmpHome := t.TempDir()

	// Load config - should return default
	config, err := LoadUserConfig(tmpHome)
	if err != nil {
		t.Fatalf("LoadUserConfig() should not error when no config exists: %v", err)
	}

	if config == nil {
		t.Fatal("LoadUserConfig() should return default config")
	}

	if config.Env == nil {
		t.Error("LoadUserConfig() should initialize Env map")
	}
}

func TestValidate(t *testing.T) {
	config := DefaultConfig()
	if err := config.Validate(); err != nil {
		t.Errorf("Validate() should not error for default config: %v", err)
	}

	config.Env["TEST"] = "value"
	if err := config.Validate(); err != nil {
		t.Errorf("Validate() should not error for config with env vars: %v", err)
	}
}

func TestSaveConfig(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".camp", "camp.yml")

	// Create config
	config := &CampConfig{
		Env: map[string]string{
			"EDITOR": "nvim",
			"SHELL":  "/bin/zsh",
		},
	}

	// Save config
	if err := config.SaveConfig(configPath); err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("SaveConfig() did not create file")
	}

	// Load and verify
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedConfig.Env["EDITOR"] != "nvim" {
		t.Errorf("Expected EDITOR=nvim, got %s", loadedConfig.Env["EDITOR"])
	}

	if loadedConfig.Env["SHELL"] != "/bin/zsh" {
		t.Errorf("Expected SHELL=/bin/zsh, got %s", loadedConfig.Env["SHELL"])
	}
}

func TestUserReload(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	// Write config
	configPath := filepath.Join(campDir, "camp.yml")
	yamlContent := `env:
  CUSTOM_VAR: custom_value
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Create user
	user := &User{
		HomeDir: tmpHome,
		EnvVars: make(map[string]string),
	}

	// Reload config
	if err := user.Reload(); err != nil {
		t.Fatalf("Reload() failed: %v", err)
	}

	// Verify env vars were loaded
	if user.EnvVars["CUSTOM_VAR"] != "custom_value" {
		t.Errorf("Expected CUSTOM_VAR=custom_value, got %s", user.EnvVars["CUSTOM_VAR"])
	}
}

func TestUserReload_NoConfig(t *testing.T) {
	// Create temporary home directory without config
	tmpHome := t.TempDir()

	// Create user
	user := &User{
		HomeDir: tmpHome,
		EnvVars: make(map[string]string),
	}

	// Reload config - should not error
	if err := user.Reload(); err != nil {
		t.Fatalf("Reload() should not error when no config exists: %v", err)
	}

	// EnvVars should be empty but initialized
	if user.EnvVars == nil {
		t.Error("Reload() should initialize EnvVars")
	}

	if len(user.EnvVars) != 0 {
		t.Errorf("Expected empty EnvVars, got %d entries", len(user.EnvVars))
	}
}
