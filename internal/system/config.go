package system

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CampConfig represents the camp.yml configuration file
type CampConfig struct {
	Env    map[string]string `yaml:"env"`    // Environment variables
	Flakes []Flake           `yaml:"flakes"` // External Nix flakes to integrate
}

// DefaultConfig returns a CampConfig with sensible defaults
func DefaultConfig() *CampConfig {
	return &CampConfig{
		Env:    make(map[string]string),
		Flakes: []Flake{},
	}
}

// LoadConfig loads the camp configuration from the specified path
// If the file doesn't exist, returns a default config without error
func LoadConfig(path string) (*CampConfig, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config CampConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Initialize Env map if nil
	if config.Env == nil {
		config.Env = make(map[string]string)
	}

	// Initialize Flakes slice if nil
	if config.Flakes == nil {
		config.Flakes = []Flake{}
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// LoadUserConfig loads the camp configuration from the user's home directory
// Looks for ~/.camp/camp.yml or ~/.camp/camp.yaml
func LoadUserConfig(homeDir string) (*CampConfig, error) {
	// Try .yml first, then .yaml
	ymlPath := filepath.Join(homeDir, ".camp", "camp.yml")
	yamlPath := filepath.Join(homeDir, ".camp", "camp.yaml")

	// Check which files exist
	ymlExists := false
	yamlExists := false

	if _, err := os.Stat(ymlPath); err == nil {
		ymlExists = true
	}
	if _, err := os.Stat(yamlPath); err == nil {
		yamlExists = true
	}

	// If .yml exists, use it
	if ymlExists {
		return LoadConfig(ymlPath)
	}

	// Otherwise try .yaml
	if yamlExists {
		return LoadConfig(yamlPath)
	}

	// Neither exists, return default
	return DefaultConfig(), nil
}

// Validate checks if the configuration is valid
func (c *CampConfig) Validate() error {
	// Validate flakes configuration
	if err := c.ValidateFlakes(); err != nil {
		return err
	}
	return nil
}

// ValidateFlakes validates the flakes configuration
func (c *CampConfig) ValidateFlakes() error {
	if c.Flakes == nil || len(c.Flakes) == 0 {
		// Empty flakes is valid
		return nil
	}

	// Track flake names to ensure uniqueness
	names := make(map[string]bool)

	for i, flake := range c.Flakes {
		// Validate name is not empty
		if flake.Name == "" {
			return fmt.Errorf("flake at index %d has empty name", i)
		}

		// Validate name is unique
		if names[flake.Name] {
			return fmt.Errorf("duplicate flake name '%s' - flake names must be unique", flake.Name)
		}
		names[flake.Name] = true

		// Validate name is a valid Nix identifier
		if !isValidNixIdentifier(flake.Name) {
			return fmt.Errorf("flake '%s' has invalid name - must contain only letters, numbers, hyphens, and underscores", flake.Name)
		}

		// Validate URL is not empty
		if flake.URL == "" {
			return fmt.Errorf("flake '%s' has empty URL", flake.Name)
		}

		// Validate outputs
		if len(flake.Outputs) == 0 {
			return fmt.Errorf("flake '%s' has no outputs defined - at least one output is required", flake.Name)
		}

		for j, output := range flake.Outputs {
			// Validate output name is not empty
			if output.Name == "" {
				return fmt.Errorf("flake '%s' output at index %d has empty name", flake.Name, j)
			}

			// Validate output type is valid
			if output.Type != OutputTypeSystem && output.Type != OutputTypeHome {
				return fmt.Errorf("flake '%s' output '%s' has invalid type '%s' - must be 'system' or 'home'",
					flake.Name, output.Name, output.Type)
			}
		}
	}

	return nil
}

// isValidNixIdentifier checks if a string is a valid Nix identifier
// Valid identifiers contain only letters, numbers, hyphens, and underscores
func isValidNixIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' ||
			r == '_') {
			return false
		}
	}

	return true
}

// SaveConfig saves the configuration to the specified path
func (c *CampConfig) SaveConfig(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
