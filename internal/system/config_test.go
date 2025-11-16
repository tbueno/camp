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

// Flake tests

func TestDefaultConfig_InitializesFlakes(t *testing.T) {
	config := DefaultConfig()

	if config.Flakes == nil {
		t.Error("DefaultConfig() should initialize Flakes slice")
	}

	if len(config.Flakes) != 0 {
		t.Errorf("DefaultConfig() should have empty Flakes slice, got %d entries", len(config.Flakes))
	}
}

func TestLoadConfig_WithFlakes(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write YAML with flakes
	yamlContent := `env:
  EDITOR: nvim

flakes:
  - name: personal-tools
    url: "github:user/my-tools"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: packages
        type: home
      - name: homeManagerModules.default
        type: home
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify flakes
	if len(config.Flakes) != 1 {
		t.Fatalf("Expected 1 flake, got %d", len(config.Flakes))
	}

	flake := config.Flakes[0]
	if flake.Name != "personal-tools" {
		t.Errorf("Expected name=personal-tools, got %s", flake.Name)
	}

	if flake.URL != "github:user/my-tools" {
		t.Errorf("Expected URL=github:user/my-tools, got %s", flake.URL)
	}

	if len(flake.Follows) != 1 {
		t.Errorf("Expected 1 follows entry, got %d", len(flake.Follows))
	}

	if flake.Follows["nixpkgs"] != "nixpkgs" {
		t.Errorf("Expected follows[nixpkgs]=nixpkgs, got %s", flake.Follows["nixpkgs"])
	}

	if len(flake.Outputs) != 2 {
		t.Fatalf("Expected 2 outputs, got %d", len(flake.Outputs))
	}

	if flake.Outputs[0].Name != "packages" {
		t.Errorf("Expected output[0].Name=packages, got %s", flake.Outputs[0].Name)
	}

	if flake.Outputs[0].Type != OutputTypeHome {
		t.Errorf("Expected output[0].Type=home, got %s", flake.Outputs[0].Type)
	}

	if flake.Outputs[1].Name != "homeManagerModules.default" {
		t.Errorf("Expected output[1].Name=homeManagerModules.default, got %s", flake.Outputs[1].Name)
	}
}

func TestLoadConfig_WithMultipleFlakes(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write YAML with multiple flakes
	yamlContent := `flakes:
  - name: flake1
    url: "github:user/flake1"
    outputs:
      - name: packages
        type: home

  - name: flake2
    url: "git+ssh://git@github.com/org/flake2.git"
    outputs:
      - name: darwinModules.team
        type: system

  - name: flake3
    url: "path:/Users/test/dev/flake3"
    outputs:
      - name: homeManagerModules.dev
        type: home
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify flakes count
	if len(config.Flakes) != 3 {
		t.Fatalf("Expected 3 flakes, got %d", len(config.Flakes))
	}

	// Verify flake names
	if config.Flakes[0].Name != "flake1" {
		t.Errorf("Expected flake[0].Name=flake1, got %s", config.Flakes[0].Name)
	}

	if config.Flakes[1].Name != "flake2" {
		t.Errorf("Expected flake[1].Name=flake2, got %s", config.Flakes[1].Name)
	}

	if config.Flakes[2].Name != "flake3" {
		t.Errorf("Expected flake[2].Name=flake3, got %s", config.Flakes[2].Name)
	}

	// Verify output types
	if config.Flakes[0].Outputs[0].Type != OutputTypeHome {
		t.Errorf("Expected flake[0] output type=home, got %s", config.Flakes[0].Outputs[0].Type)
	}

	if config.Flakes[1].Outputs[0].Type != OutputTypeSystem {
		t.Errorf("Expected flake[1] output type=system, got %s", config.Flakes[1].Outputs[0].Type)
	}
}

func TestLoadConfig_EmptyFlakes(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write YAML with empty flakes array
	yamlContent := `env:
  TEST: value

flakes: []
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify Flakes is initialized
	if config.Flakes == nil {
		t.Error("LoadConfig() should initialize Flakes slice even when empty")
	}

	if len(config.Flakes) != 0 {
		t.Errorf("Expected 0 flakes, got %d", len(config.Flakes))
	}
}

func TestLoadConfig_NoFlakesSection(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write YAML without flakes section
	yamlContent := `env:
  TEST: value
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify Flakes is initialized
	if config.Flakes == nil {
		t.Error("LoadConfig() should initialize Flakes slice when section is missing")
	}
}

func TestUserReload_WithFlakes(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	// Write config with flakes
	configPath := filepath.Join(campDir, "camp.yml")
	yamlContent := `env:
  EDITOR: nvim

flakes:
  - name: test-flake
    url: "github:test/flake"
    outputs:
      - name: packages
        type: home
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Create user
	user := &User{
		HomeDir: tmpHome,
		EnvVars: make(map[string]string),
		Flakes:  []Flake{},
	}

	// Reload config
	if err := user.Reload(); err != nil {
		t.Fatalf("Reload() failed: %v", err)
	}

	// Verify flakes were loaded
	if len(user.Flakes) != 1 {
		t.Fatalf("Expected 1 flake, got %d", len(user.Flakes))
	}

	if user.Flakes[0].Name != "test-flake" {
		t.Errorf("Expected flake name=test-flake, got %s", user.Flakes[0].Name)
	}

	if user.Flakes[0].URL != "github:test/flake" {
		t.Errorf("Expected flake URL=github:test/flake, got %s", user.Flakes[0].URL)
	}

	// Verify env vars were also loaded
	if user.EnvVars["EDITOR"] != "nvim" {
		t.Errorf("Expected EDITOR=nvim, got %s", user.EnvVars["EDITOR"])
	}
}

func TestUserReload_NoFlakes(t *testing.T) {
	// Create temporary home directory without config
	tmpHome := t.TempDir()

	// Create user
	user := &User{
		HomeDir: tmpHome,
		EnvVars: make(map[string]string),
		Flakes:  []Flake{},
	}

	// Reload config - should not error
	if err := user.Reload(); err != nil {
		t.Fatalf("Reload() should not error when no config exists: %v", err)
	}

	// Flakes should be empty but initialized
	if user.Flakes == nil {
		t.Error("Reload() should initialize Flakes")
	}

	if len(user.Flakes) != 0 {
		t.Errorf("Expected empty Flakes, got %d entries", len(user.Flakes))
	}
}

func TestSaveConfig_WithFlakes(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".camp", "camp.yml")

	// Create config with flakes
	config := &CampConfig{
		Env: map[string]string{
			"EDITOR": "nvim",
		},
		Flakes: []Flake{
			{
				Name: "my-flake",
				URL:  "github:user/repo",
				Follows: map[string]string{
					"nixpkgs": "nixpkgs",
				},
				Outputs: []FlakeOutput{
					{
						Name: "packages",
						Type: OutputTypeHome,
					},
				},
			},
		},
	}

	// Save config
	if err := config.SaveConfig(configPath); err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Load and verify
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	// Verify flakes
	if len(loadedConfig.Flakes) != 1 {
		t.Fatalf("Expected 1 flake, got %d", len(loadedConfig.Flakes))
	}

	if loadedConfig.Flakes[0].Name != "my-flake" {
		t.Errorf("Expected flake name=my-flake, got %s", loadedConfig.Flakes[0].Name)
	}

	if loadedConfig.Flakes[0].URL != "github:user/repo" {
		t.Errorf("Expected flake URL=github:user/repo, got %s", loadedConfig.Flakes[0].URL)
	}

	if loadedConfig.Flakes[0].Follows["nixpkgs"] != "nixpkgs" {
		t.Errorf("Expected follows[nixpkgs]=nixpkgs, got %s", loadedConfig.Flakes[0].Follows["nixpkgs"])
	}
}
