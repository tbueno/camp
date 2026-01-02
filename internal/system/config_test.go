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

// Validation tests

func TestValidateFlakes_DuplicateNames(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "duplicate",
				URL:  "github:user/flake1",
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
			{
				Name: "duplicate",
				URL:  "github:user/flake2",
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	err := config.ValidateFlakes()
	if err == nil {
		t.Error("ValidateFlakes() should error on duplicate flake names")
	}

	if err != nil && err.Error() != "duplicate flake name 'duplicate' - flake names must be unique" {
		t.Errorf("Expected duplicate name error, got: %v", err)
	}
}

func TestValidateFlakes_EmptyName(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "",
				URL:  "github:user/flake",
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	err := config.ValidateFlakes()
	if err == nil {
		t.Error("ValidateFlakes() should error on empty flake name")
	}

	if err != nil && err.Error() != "flake at index 0 has empty name" {
		t.Errorf("Expected empty name error, got: %v", err)
	}
}

func TestValidateFlakes_InvalidNixIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		flakeName string
	}{
		{"with spaces", "my flake"},
		{"with dots", "my.flake"},
		{"with special chars", "my@flake"},
		{"with slashes", "my/flake"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CampConfig{
				Flakes: []Flake{
					{
						Name: tt.flakeName,
						URL:  "github:user/flake",
						Outputs: []FlakeOutput{
							{Name: "packages", Type: OutputTypeHome},
						},
					},
				},
			}

			err := config.ValidateFlakes()
			if err == nil {
				t.Errorf("ValidateFlakes() should error on invalid Nix identifier '%s'", tt.flakeName)
			}
		})
	}
}

func TestValidateFlakes_ValidNixIdentifiers(t *testing.T) {
	validNames := []string{
		"my-flake",
		"my_flake",
		"MyFlake",
		"flake123",
		"123flake",
		"FLAKE",
		"_flake_",
		"flake-with-many-hyphens",
	}

	for _, name := range validNames {
		t.Run(name, func(t *testing.T) {
			config := &CampConfig{
				Flakes: []Flake{
					{
						Name: name,
						URL:  "github:user/flake",
						Outputs: []FlakeOutput{
							{Name: "packages", Type: OutputTypeHome},
						},
					},
				},
			}

			if err := config.ValidateFlakes(); err != nil {
				t.Errorf("ValidateFlakes() should accept valid Nix identifier '%s', got error: %v", name, err)
			}
		})
	}
}

func TestValidateFlakes_EmptyURL(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "my-flake",
				URL:  "",
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	err := config.ValidateFlakes()
	if err == nil {
		t.Error("ValidateFlakes() should error on empty URL")
	}

	if err != nil && err.Error() != "flake 'my-flake' has empty URL" {
		t.Errorf("Expected empty URL error, got: %v", err)
	}
}

func TestValidateFlakes_NoOutputs(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name:    "my-flake",
				URL:     "github:user/flake",
				Outputs: []FlakeOutput{},
			},
		},
	}

	err := config.ValidateFlakes()
	if err == nil {
		t.Error("ValidateFlakes() should error on flake with no outputs")
	}

	if err != nil && err.Error() != "flake 'my-flake' has no outputs defined - at least one output is required" {
		t.Errorf("Expected no outputs error, got: %v", err)
	}
}

func TestValidateFlakes_EmptyOutputName(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "my-flake",
				URL:  "github:user/flake",
				Outputs: []FlakeOutput{
					{Name: "", Type: OutputTypeHome},
				},
			},
		},
	}

	err := config.ValidateFlakes()
	if err == nil {
		t.Error("ValidateFlakes() should error on empty output name")
	}

	if err != nil && err.Error() != "flake 'my-flake' output at index 0 has empty name" {
		t.Errorf("Expected empty output name error, got: %v", err)
	}
}

func TestValidateFlakes_InvalidOutputType(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "my-flake",
				URL:  "github:user/flake",
				Outputs: []FlakeOutput{
					{Name: "packages", Type: "invalid"},
				},
			},
		},
	}

	err := config.ValidateFlakes()
	if err == nil {
		t.Error("ValidateFlakes() should error on invalid output type")
	}

	if err != nil && err.Error() != "flake 'my-flake' output 'packages' has invalid type 'invalid' - must be 'system' or 'home'" {
		t.Errorf("Expected invalid output type error, got: %v", err)
	}
}

func TestValidateFlakes_EmptyFlakes(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{},
	}

	if err := config.ValidateFlakes(); err != nil {
		t.Errorf("ValidateFlakes() should not error on empty flakes, got: %v", err)
	}
}

func TestValidateFlakes_NilFlakes(t *testing.T) {
	config := &CampConfig{
		Flakes: nil,
	}

	if err := config.ValidateFlakes(); err != nil {
		t.Errorf("ValidateFlakes() should not error on nil flakes, got: %v", err)
	}
}

func TestValidateFlakes_ValidConfiguration(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "flake1",
				URL:  "github:user/flake1",
				Follows: map[string]string{
					"nixpkgs": "nixpkgs",
				},
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
					{Name: "homeManagerModules.default", Type: OutputTypeHome},
				},
			},
			{
				Name: "flake2",
				URL:  "git+ssh://git@github.com/org/flake2.git",
				Outputs: []FlakeOutput{
					{Name: "darwinModules.team", Type: OutputTypeSystem},
				},
			},
		},
	}

	if err := config.ValidateFlakes(); err != nil {
		t.Errorf("ValidateFlakes() should not error on valid configuration, got: %v", err)
	}
}

func TestLoadConfig_ValidationError(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write YAML with invalid flakes (duplicate names)
	yamlContent := `flakes:
  - name: duplicate
    url: "github:user/flake1"
    outputs:
      - name: packages
        type: home

  - name: duplicate
    url: "github:user/flake2"
    outputs:
      - name: packages
        type: home
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config - should return validation error
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig() should return validation error for invalid flakes")
	}

	if err != nil && !os.IsNotExist(err) {
		// Check that it's a validation error
		expectedMsg := "invalid configuration: duplicate flake name 'duplicate' - flake names must be unique"
		if err.Error() != expectedMsg {
			t.Errorf("Expected validation error, got: %v", err)
		}
	}
}

// Flake Arguments Validation Tests

func TestValidateFlakeArgs_ValidArgs(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "test-flake",
				URL:  "github:user/flake",
				Args: map[string]interface{}{
					"email":          "test@example.com",
					"enableFeature":  true,
					"fontSize":       14,
					"threshold":      3.14,
					"packages":       []interface{}{"vim", "git", "tmux"},
					"ports":          []interface{}{8080, 9090, 3000},
					"enabledModules": []interface{}{true, false, true},
				},
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	if err := config.ValidateFlakes(); err != nil {
		t.Errorf("ValidateFlakes() should not error on valid args, got: %v", err)
	}
}

func TestValidateFlakeArgs_EmptyArgs(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "test-flake",
				URL:  "github:user/flake",
				Args: map[string]interface{}{},
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	if err := config.ValidateFlakes(); err != nil {
		t.Errorf("ValidateFlakes() should not error on empty args map, got: %v", err)
	}
}

func TestValidateFlakeArgs_NilArgs(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "test-flake",
				URL:  "github:user/flake",
				Args: nil,
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	if err := config.ValidateFlakes(); err != nil {
		t.Errorf("ValidateFlakes() should not error on nil args, got: %v", err)
	}
}

func TestValidateFlakeArgs_ReservedNames(t *testing.T) {
	tests := []struct {
		name        string
		argName     string
		expectedMsg string
	}{
		{
			name:        "userName reserved",
			argName:     "userName",
			expectedMsg: "flake 'test-flake' argument 'userName' uses a reserved name - userName, hostName, and home are automatically provided",
		},
		{
			name:        "hostName reserved",
			argName:     "hostName",
			expectedMsg: "flake 'test-flake' argument 'hostName' uses a reserved name - userName, hostName, and home are automatically provided",
		},
		{
			name:        "home reserved",
			argName:     "home",
			expectedMsg: "flake 'test-flake' argument 'home' uses a reserved name - userName, hostName, and home are automatically provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CampConfig{
				Flakes: []Flake{
					{
						Name: "test-flake",
						URL:  "github:user/flake",
						Args: map[string]interface{}{
							tt.argName: "some-value",
						},
						Outputs: []FlakeOutput{
							{Name: "packages", Type: OutputTypeHome},
						},
					},
				},
			}

			err := config.ValidateFlakes()
			if err == nil {
				t.Errorf("ValidateFlakes() should error on reserved arg name '%s'", tt.argName)
			}

			if err != nil && err.Error() != tt.expectedMsg {
				t.Errorf("Expected error '%s', got: %v", tt.expectedMsg, err)
			}
		})
	}
}

func TestValidateFlakeArgs_InvalidArgNames(t *testing.T) {
	tests := []struct {
		name    string
		argName string
	}{
		{"with spaces", "my arg"},
		{"with dots", "my.arg"},
		{"with special chars", "my@arg"},
		{"with slashes", "my/arg"},
		{"with colons", "my:arg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CampConfig{
				Flakes: []Flake{
					{
						Name: "test-flake",
						URL:  "github:user/flake",
						Args: map[string]interface{}{
							tt.argName: "value",
						},
						Outputs: []FlakeOutput{
							{Name: "packages", Type: OutputTypeHome},
						},
					},
				},
			}

			err := config.ValidateFlakes()
			if err == nil {
				t.Errorf("ValidateFlakes() should error on invalid arg name '%s'", tt.argName)
			}
		})
	}
}

func TestValidateFlakeArgs_ValidArgNames(t *testing.T) {
	validNames := []string{
		"email",
		"enable-feature",
		"enable_feature",
		"FONT_SIZE",
		"arg123",
		"_private",
		"arg-with-many-hyphens",
		"arg_with_underscores",
	}

	for _, argName := range validNames {
		t.Run(argName, func(t *testing.T) {
			config := &CampConfig{
				Flakes: []Flake{
					{
						Name: "test-flake",
						URL:  "github:user/flake",
						Args: map[string]interface{}{
							argName: "value",
						},
						Outputs: []FlakeOutput{
							{Name: "packages", Type: OutputTypeHome},
						},
					},
				},
			}

			if err := config.ValidateFlakes(); err != nil {
				t.Errorf("ValidateFlakes() should accept valid arg name '%s', got error: %v", argName, err)
			}
		})
	}
}

func TestValidateFlakeArgs_UnsupportedTypes(t *testing.T) {
	tests := []struct {
		name        string
		argValue    interface{}
		expectedMsg string
	}{
		{
			name:        "map type",
			argValue:    map[string]string{"key": "value"},
			expectedMsg: "flake 'test-flake' argument 'testArg' has unsupported type - only string, bool, number, and list types are supported",
		},
		{
			name:        "nil value",
			argValue:    nil,
			expectedMsg: "flake 'test-flake' argument 'testArg' has unsupported type - only string, bool, number, and list types are supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CampConfig{
				Flakes: []Flake{
					{
						Name: "test-flake",
						URL:  "github:user/flake",
						Args: map[string]interface{}{
							"testArg": tt.argValue,
						},
						Outputs: []FlakeOutput{
							{Name: "packages", Type: OutputTypeHome},
						},
					},
				},
			}

			err := config.ValidateFlakes()
			if err == nil {
				t.Errorf("ValidateFlakes() should error on unsupported type %T", tt.argValue)
			}

			if err != nil && err.Error() != tt.expectedMsg {
				t.Errorf("Expected error '%s', got: %v", tt.expectedMsg, err)
			}
		})
	}
}

func TestValidateFlakeArgs_ListWithInvalidElements(t *testing.T) {
	tests := []struct {
		name        string
		listValue   []interface{}
		expectedMsg string
	}{
		{
			name:        "list with map element",
			listValue:   []interface{}{"string", map[string]string{"key": "value"}},
			expectedMsg: "flake 'test-flake' argument 'myList' list element at index 1 has unsupported type (only string, bool, number are supported in lists)",
		},
		{
			name:        "list with nil element",
			listValue:   []interface{}{"string", nil},
			expectedMsg: "flake 'test-flake' argument 'myList' list element at index 1 has unsupported type (only string, bool, number are supported in lists)",
		},
		{
			name:        "list with nested list",
			listValue:   []interface{}{"string", []interface{}{"nested"}},
			expectedMsg: "flake 'test-flake' argument 'myList' list element at index 1 has unsupported type (only string, bool, number are supported in lists)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CampConfig{
				Flakes: []Flake{
					{
						Name: "test-flake",
						URL:  "github:user/flake",
						Args: map[string]interface{}{
							"myList": tt.listValue,
						},
						Outputs: []FlakeOutput{
							{Name: "packages", Type: OutputTypeHome},
						},
					},
				},
			}

			err := config.ValidateFlakes()
			if err == nil {
				t.Error("ValidateFlakes() should error on list with invalid element types")
			}

			if err != nil && err.Error() != tt.expectedMsg {
				t.Errorf("Expected error '%s', got: %v", tt.expectedMsg, err)
			}
		})
	}
}

func TestValidateFlakeArgs_EmptyList(t *testing.T) {
	config := &CampConfig{
		Flakes: []Flake{
			{
				Name: "test-flake",
				URL:  "github:user/flake",
				Args: map[string]interface{}{
					"emptyList": []interface{}{},
				},
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	if err := config.ValidateFlakes(); err != nil {
		t.Errorf("ValidateFlakes() should not error on empty list, got: %v", err)
	}
}

func TestLoadConfig_WithFlakeArgs(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Write YAML with flake args
	yamlContent := `env:
  EDITOR: nvim

flakes:
  - name: personal-config
    url: "github:user/config"
    args:
      email: "test@example.com"
      enableDevTools: true
      fontSize: 14
      packages:
        - vim
        - git
        - tmux
    outputs:
      - name: darwinModules.default
        type: system
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

	// Verify flake args were loaded
	if len(config.Flakes) != 1 {
		t.Fatalf("Expected 1 flake, got %d", len(config.Flakes))
	}

	flake := config.Flakes[0]
	if flake.Args == nil {
		t.Fatal("Flake Args should not be nil")
	}

	// Check string arg
	if email, ok := flake.Args["email"].(string); !ok || email != "test@example.com" {
		t.Errorf("Expected email='test@example.com', got: %v", flake.Args["email"])
	}

	// Check bool arg
	if enabled, ok := flake.Args["enableDevTools"].(bool); !ok || !enabled {
		t.Errorf("Expected enableDevTools=true, got: %v", flake.Args["enableDevTools"])
	}

	// Check int arg
	if fontSize, ok := flake.Args["fontSize"].(int); !ok || fontSize != 14 {
		t.Errorf("Expected fontSize=14, got: %v", flake.Args["fontSize"])
	}

	// Check list arg
	if packages, ok := flake.Args["packages"].([]interface{}); !ok {
		t.Errorf("Expected packages to be a list, got: %T", flake.Args["packages"])
	} else {
		if len(packages) != 3 {
			t.Errorf("Expected 3 packages, got %d", len(packages))
		}
		if pkg, ok := packages[0].(string); !ok || pkg != "vim" {
			t.Errorf("Expected first package='vim', got: %v", packages[0])
		}
	}
}

func TestSaveAndLoadConfig_PreservesArgTypes(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	// Create config with various arg types
	config := &CampConfig{
		Env: map[string]string{
			"EDITOR": "nvim",
		},
		Flakes: []Flake{
			{
				Name: "test-flake",
				URL:  "github:user/flake",
				Args: map[string]interface{}{
					"stringArg": "hello",
					"boolArg":   true,
					"intArg":    42,
					"floatArg":  3.14,
					"listArg":   []interface{}{"a", "b", "c"},
				},
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	// Save config
	if err := config.SaveConfig(configPath); err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Load config
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify args were preserved with correct types
	args := loadedConfig.Flakes[0].Args

	if str, ok := args["stringArg"].(string); !ok || str != "hello" {
		t.Errorf("String arg not preserved correctly: %v (%T)", args["stringArg"], args["stringArg"])
	}

	if b, ok := args["boolArg"].(bool); !ok || !b {
		t.Errorf("Bool arg not preserved correctly: %v (%T)", args["boolArg"], args["boolArg"])
	}

	if i, ok := args["intArg"].(int); !ok || i != 42 {
		t.Errorf("Int arg not preserved correctly: %v (%T)", args["intArg"], args["intArg"])
	}

	if f, ok := args["floatArg"].(float64); !ok || f != 3.14 {
		t.Errorf("Float arg not preserved correctly: %v (%T)", args["floatArg"], args["floatArg"])
	}

	if list, ok := args["listArg"].([]interface{}); !ok || len(list) != 3 {
		t.Errorf("List arg not preserved correctly: %v (%T)", args["listArg"], args["listArg"])
	}
}

// ============================================================================
// Package Tests
// ============================================================================

func TestLoadConfig_WithPackages(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	yamlContent := `env:
  EDITOR: nvim

packages:
  - git
  - neovim
  - ripgrep
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if len(config.Packages) != 3 {
		t.Fatalf("Expected 3 packages, got %d", len(config.Packages))
	}

	if config.Packages[0] != "git" {
		t.Errorf("Expected package[0]=git, got %s", config.Packages[0])
	}

	if config.Packages[1] != "neovim" {
		t.Errorf("Expected package[1]=neovim, got %s", config.Packages[1])
	}

	if config.Packages[2] != "ripgrep" {
		t.Errorf("Expected package[2]=ripgrep, got %s", config.Packages[2])
	}
}

func TestLoadConfig_EmptyPackages(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	yamlContent := `env:
  TEST: value

packages: []
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.Packages == nil {
		t.Error("LoadConfig() should initialize Packages slice even when empty")
	}

	if len(config.Packages) != 0 {
		t.Errorf("Expected 0 packages, got %d", len(config.Packages))
	}
}

func TestLoadConfig_NoPackagesSection(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "camp.yml")

	yamlContent := `env:
  TEST: value
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.Packages == nil {
		t.Error("LoadConfig() should initialize Packages slice when section is missing")
	}
}

func TestValidatePackages_DuplicatePackages(t *testing.T) {
	config := &CampConfig{
		Packages: []string{"git", "neovim", "git"},
	}

	err := config.ValidatePackages()
	if err == nil {
		t.Error("ValidatePackages() should error on duplicate package names")
	}

	if err != nil && !contains(err.Error(), "duplicate package") {
		t.Errorf("Expected duplicate package error, got: %v", err)
	}
}

func TestValidatePackages_EmptyPackageName(t *testing.T) {
	config := &CampConfig{
		Packages: []string{"git", "", "neovim"},
	}

	err := config.ValidatePackages()
	if err == nil {
		t.Error("ValidatePackages() should error on empty package name")
	}

	if err != nil && !contains(err.Error(), "empty or contains only whitespace") {
		t.Errorf("Expected empty package name error, got: %v", err)
	}
}

func TestValidatePackages_WhitespaceOnlyPackageName(t *testing.T) {
	config := &CampConfig{
		Packages: []string{"git", "   ", "neovim"},
	}

	err := config.ValidatePackages()
	if err == nil {
		t.Error("ValidatePackages() should error on whitespace-only package name")
	}

	if err != nil && !contains(err.Error(), "empty or contains only whitespace") {
		t.Errorf("Expected whitespace-only package name error, got: %v", err)
	}
}

func TestValidatePackages_InvalidCharacters(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
	}{
		{"with spaces", "my package"},
		{"with special chars", "my@package"},
		{"with slashes", "my/package"},
		{"with colons", "my:package"},
		{"with brackets", "my[package]"},
		{"with parens", "my(package)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CampConfig{
				Packages: []string{tt.packageName},
			}

			err := config.ValidatePackages()
			if err == nil {
				t.Errorf("ValidatePackages() should error on invalid package name '%s'", tt.packageName)
			}
		})
	}
}

func TestValidatePackages_ValidPackageNames(t *testing.T) {
	validPackages := []string{
		"git",
		"neovim",
		"ripgrep",
		"python3",
		"nodejs_20",
		"python3Packages.requests",
		"haskellPackages.pandoc",
		"package-with-hyphens",
		"package_with_underscores",
		"UPPERCASE",
		"123numbers",
		"some.package.with.dots",
	}

	config := &CampConfig{
		Packages: validPackages,
	}

	if err := config.ValidatePackages(); err != nil {
		t.Errorf("ValidatePackages() should accept valid package names, got error: %v", err)
	}
}

func TestValidatePackages_EmptyList(t *testing.T) {
	config := &CampConfig{
		Packages: []string{},
	}

	if err := config.ValidatePackages(); err != nil {
		t.Errorf("ValidatePackages() should accept empty package list, got error: %v", err)
	}
}

func TestValidatePackages_NilList(t *testing.T) {
	config := &CampConfig{
		Packages: nil,
	}

	if err := config.ValidatePackages(); err != nil {
		t.Errorf("ValidatePackages() should accept nil package list, got error: %v", err)
	}
}

func TestSaveConfig_WithPackages(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".camp", "camp.yml")

	config := &CampConfig{
		Env: map[string]string{
			"EDITOR": "nvim",
		},
		Packages: []string{"git", "neovim", "ripgrep"},
	}

	if err := config.SaveConfig(configPath); err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if len(loadedConfig.Packages) != 3 {
		t.Fatalf("Expected 3 packages, got %d", len(loadedConfig.Packages))
	}

	if loadedConfig.Packages[0] != "git" {
		t.Errorf("Expected package[0]=git, got %s", loadedConfig.Packages[0])
	}

	if loadedConfig.Packages[1] != "neovim" {
		t.Errorf("Expected package[1]=neovim, got %s", loadedConfig.Packages[1])
	}

	if loadedConfig.Packages[2] != "ripgrep" {
		t.Errorf("Expected package[2]=ripgrep, got %s", loadedConfig.Packages[2])
	}
}

func TestFindProjectConfigPath_CurrentDir(t *testing.T) {
	// Create temporary directory with .camp.yml
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".camp.yml")

	// Create .camp.yml
	yamlContent := `env:
  TEST_VAR: "test"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Find config starting from tmpDir
	foundPath := FindProjectConfigPath(tmpDir)

	if foundPath == "" {
		t.Error("FindProjectConfigPath() should find .camp.yml in current directory")
	}

	if foundPath != configPath {
		t.Errorf("FindProjectConfigPath() = %q, want %q", foundPath, configPath)
	}
}

func TestFindProjectConfigPath_ParentDir(t *testing.T) {
	// Create temporary directory structure:
	// tmpDir/
	//   .camp.yml
	//   subdir/
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".camp.yml")
	subDir := filepath.Join(tmpDir, "subdir")

	// Create .camp.yml in parent
	yamlContent := `env:
  TEST_VAR: "test"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create subdirectory
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Find config starting from subdir
	foundPath := FindProjectConfigPath(subDir)

	if foundPath == "" {
		t.Error("FindProjectConfigPath() should find .camp.yml in parent directory")
	}

	if foundPath != configPath {
		t.Errorf("FindProjectConfigPath() = %q, want %q", foundPath, configPath)
	}
}

func TestFindProjectConfigPath_MultipleParents(t *testing.T) {
	// Create temporary directory structure:
	// tmpDir/
	//   .camp.yml
	//   level1/
	//     level2/
	//       level3/
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".camp.yml")
	level3Dir := filepath.Join(tmpDir, "level1", "level2", "level3")

	// Create .camp.yml at root
	yamlContent := `env:
  TEST_VAR: "test"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create nested directories
	if err := os.MkdirAll(level3Dir, 0755); err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}

	// Find config starting from level3
	foundPath := FindProjectConfigPath(level3Dir)

	if foundPath == "" {
		t.Error("FindProjectConfigPath() should find .camp.yml walking up multiple levels")
	}

	if foundPath != configPath {
		t.Errorf("FindProjectConfigPath() = %q, want %q", foundPath, configPath)
	}
}

func TestFindProjectConfigPath_NotFound(t *testing.T) {
	// Create temporary directory without .camp.yml
	tmpDir := t.TempDir()

	// Find config should return empty string
	foundPath := FindProjectConfigPath(tmpDir)

	if foundPath != "" {
		t.Errorf("FindProjectConfigPath() should return empty string when not found, got %q", foundPath)
	}
}

func TestFindProjectConfigPath_EmptyStartDir(t *testing.T) {
	// Create temporary directory with .camp.yml
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".camp.yml")

	yamlContent := `env:
  TEST_VAR: "test"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Change to tmpDir
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Find config with empty string (should use current dir)
	foundPath := FindProjectConfigPath("")

	if foundPath == "" {
		t.Error("FindProjectConfigPath(\"\") should use current directory")
	}

	// Resolve both paths to handle symlinks (e.g., /var vs /private/var on macOS)
	foundPathResolved, _ := filepath.EvalSymlinks(foundPath)
	configPathResolved, _ := filepath.EvalSymlinks(configPath)

	if foundPathResolved != configPathResolved {
		t.Errorf("FindProjectConfigPath(\"\") = %q, want %q", foundPath, configPath)
	}
}

func TestFindProjectConfigPath_ClosestConfig(t *testing.T) {
	// Create temporary directory structure:
	// tmpDir/
	//   .camp.yml (parent config)
	//   subdir/
	//     .camp.yml (child config - should find this one)
	tmpDir := t.TempDir()
	parentConfig := filepath.Join(tmpDir, ".camp.yml")
	subDir := filepath.Join(tmpDir, "subdir")
	childConfig := filepath.Join(subDir, ".camp.yml")

	// Create parent config
	parentYaml := `env:
  LEVEL: "parent"
`
	if err := os.WriteFile(parentConfig, []byte(parentYaml), 0644); err != nil {
		t.Fatalf("Failed to write parent config: %v", err)
	}

	// Create subdirectory and child config
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	childYaml := `env:
  LEVEL: "child"
`
	if err := os.WriteFile(childConfig, []byte(childYaml), 0644); err != nil {
		t.Fatalf("Failed to write child config: %v", err)
	}

	// Find config starting from subdir - should find child config, not parent
	foundPath := FindProjectConfigPath(subDir)

	if foundPath == "" {
		t.Error("FindProjectConfigPath() should find .camp.yml")
	}

	if foundPath != childConfig {
		t.Errorf("FindProjectConfigPath() should find closest config, got %q, want %q", foundPath, childConfig)
	}
}

// Helper function for substring checking
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
