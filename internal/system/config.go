package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// CampConfig represents the camp.yml configuration file
type CampConfig struct {
	Env      map[string]string `yaml:"env"`      // Environment variables
	Packages []string          `yaml:"packages"` // Nix packages to install
	Flakes   []Flake           `yaml:"flakes"`   // External Nix flakes to integrate
}

// DefaultConfig returns a CampConfig with sensible defaults
func DefaultConfig() *CampConfig {
	return &CampConfig{
		Env:      make(map[string]string),
		Packages: []string{},
		Flakes:   []Flake{},
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

	// Initialize Packages slice if nil
	if config.Packages == nil {
		config.Packages = []string{}
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

// FindProjectConfigPath searches for .camp.yml starting from the given directory
// and walking up the directory tree until finding one or reaching root.
// Returns full path if found, empty string if not found.
func FindProjectConfigPath(startDir string) string {
	dir := startDir

	// Handle empty start directory
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return ""
		}
	}

	// Convert to absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}
	dir = absDir

	// Walk up directory tree
	for {
		configPath := filepath.Join(dir, ".camp.yml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			return ""
		}
		dir = parent
	}
}

// Validate checks if the configuration is valid
func (c *CampConfig) Validate() error {
	// Validate flakes configuration
	if err := c.ValidateFlakes(); err != nil {
		return err
	}

	// Validate packages configuration
	if err := c.ValidatePackages(); err != nil {
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

		// Validate arguments
		if err := validateFlakeArgs(flake.Name, flake.Args); err != nil {
			return err
		}
	}

	return nil
}

// validateFlakeArgs validates the arguments for a flake
func validateFlakeArgs(flakeName string, args map[string]interface{}) error {
	if args == nil || len(args) == 0 {
		// Empty args is valid
		return nil
	}

	// Reserved argument names (automatically provided by camp)
	reservedNames := map[string]bool{
		"userName": true,
		"hostName": true,
		"home":     true,
	}

	for argName, argValue := range args {
		// Validate arg name is not empty
		if argName == "" {
			return fmt.Errorf("flake '%s' has an argument with empty name", flakeName)
		}

		// Validate arg name is a valid Nix identifier
		if !isValidNixIdentifier(argName) {
			return fmt.Errorf("flake '%s' argument '%s' has invalid name - must contain only letters, numbers, hyphens, and underscores", flakeName, argName)
		}

		// Check for reserved names
		if reservedNames[argName] {
			return fmt.Errorf("flake '%s' argument '%s' uses a reserved name - userName, hostName, and home are automatically provided", flakeName, argName)
		}

		// Validate argument type is supported
		if err := validateArgType(flakeName, argName, argValue); err != nil {
			return err
		}
	}

	return nil
}

// validateArgType validates that an argument value is a supported type
func validateArgType(flakeName, argName string, value interface{}) error {
	switch v := value.(type) {
	case string, bool, int, int64, float64:
		// Supported scalar types
		return nil
	case []interface{}:
		// Validate list elements are supported types
		for i, elem := range v {
			switch elem.(type) {
			case string, bool, int, int64, float64:
				// Supported element types
				continue
			default:
				return fmt.Errorf("flake '%s' argument '%s' list element at index %d has unsupported type (only string, bool, number are supported in lists)", flakeName, argName, i)
			}
		}
		return nil
	default:
		return fmt.Errorf("flake '%s' argument '%s' has unsupported type - only string, bool, number, and list types are supported", flakeName, argName)
	}
}

// ValidatePackages validates the packages configuration
func (c *CampConfig) ValidatePackages() error {
	if c.Packages == nil || len(c.Packages) == 0 {
		// Empty packages is valid
		return nil
	}

	// Track package names to ensure uniqueness
	seen := make(map[string]bool)

	for i, pkg := range c.Packages {
		// Validate package name is not empty or whitespace-only
		if strings.TrimSpace(pkg) == "" {
			return fmt.Errorf("package at index %d is empty or contains only whitespace", i)
		}

		// Validate package name doesn't contain invalid characters
		// Nix package names should be alphanumeric with hyphens, underscores, and dots
		if !isValidNixPackageName(pkg) {
			return fmt.Errorf("package '%s' has invalid format - must contain only letters, numbers, hyphens, underscores, and dots", pkg)
		}

		// Check for duplicates
		if seen[pkg] {
			return fmt.Errorf("duplicate package '%s' - package names must be unique", pkg)
		}
		seen[pkg] = true
	}

	return nil
}

// isValidNixPackageName checks if a string is a valid Nix package name
// Valid package names contain letters, numbers, hyphens, underscores, and dots
// They may also contain attribute paths like "python3Packages.requests"
func isValidNixPackageName(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' ||
			r == '_' ||
			r == '.') {
			return false
		}
	}

	return true
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
