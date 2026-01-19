package project

import (
	"os"
	"path/filepath"

	"camp/internal/system"
)

// LoadProjectConfig loads camp.yaml from the specified project directory
// Returns default empty config if file doesn't exist
// Tries camp.yaml first, then camp.yml
func LoadProjectConfig(projectPath string) (*system.CampConfig, error) {
	// Try .yaml first (project convention), then .yml
	yamlPath := filepath.Join(projectPath, "camp.yaml")
	ymlPath := filepath.Join(projectPath, "camp.yml")

	// Check which files exist and load the first one found
	if _, err := os.Stat(yamlPath); err == nil {
		return system.LoadConfig(yamlPath)
	}

	if _, err := os.Stat(ymlPath); err == nil {
		return system.LoadConfig(ymlPath)
	}

	// Neither exists, return default
	return system.DefaultConfig(), nil
}

// HasProjectConfig checks if a camp.yaml or camp.yml exists in the project directory
func HasProjectConfig(projectPath string) bool {
	yamlPath := filepath.Join(projectPath, "camp.yaml")
	ymlPath := filepath.Join(projectPath, "camp.yml")

	if _, err := os.Stat(yamlPath); err == nil {
		return true
	}
	if _, err := os.Stat(ymlPath); err == nil {
		return true
	}
	return false
}
