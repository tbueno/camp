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
	// Currently no validation needed for simple env vars
	// Future: validate flakes structure, env var names, etc.
	return nil
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
