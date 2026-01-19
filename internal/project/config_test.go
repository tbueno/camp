package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProjectConfig(t *testing.T) {
	t.Run("with camp.yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		yamlContent := `env:
  DATABASE_URL: "postgres://localhost:5432/db"
  DEBUG: "true"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		config, err := LoadProjectConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadProjectConfig failed: %v", err)
		}

		if config.Env["DATABASE_URL"] != "postgres://localhost:5432/db" {
			t.Errorf("Expected DATABASE_URL to be 'postgres://localhost:5432/db', got '%s'", config.Env["DATABASE_URL"])
		}

		if config.Env["DEBUG"] != "true" {
			t.Errorf("Expected DEBUG to be 'true', got '%s'", config.Env["DEBUG"])
		}
	})

	t.Run("with camp.yml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yml")

		yamlContent := `env:
  API_KEY: "secret"
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		config, err := LoadProjectConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadProjectConfig failed: %v", err)
		}

		if config.Env["API_KEY"] != "secret" {
			t.Errorf("Expected API_KEY to be 'secret', got '%s'", config.Env["API_KEY"])
		}
	})

	t.Run("yaml takes precedence over yml", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create both files with different content
		yamlPath := filepath.Join(tmpDir, "camp.yaml")
		ymlPath := filepath.Join(tmpDir, "camp.yml")

		yamlContent := `env:
  SOURCE: "yaml"
`
		ymlContent := `env:
  SOURCE: "yml"
`
		if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write camp.yaml: %v", err)
		}
		if err := os.WriteFile(ymlPath, []byte(ymlContent), 0644); err != nil {
			t.Fatalf("Failed to write camp.yml: %v", err)
		}

		config, err := LoadProjectConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadProjectConfig failed: %v", err)
		}

		// camp.yaml should take precedence
		if config.Env["SOURCE"] != "yaml" {
			t.Errorf("Expected SOURCE to be 'yaml' (from camp.yaml), got '%s'", config.Env["SOURCE"])
		}
	})

	t.Run("without config file", func(t *testing.T) {
		tmpDir := t.TempDir()

		config, err := LoadProjectConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadProjectConfig should not fail without config: %v", err)
		}

		// Should return default config with empty env
		if len(config.Env) != 0 {
			t.Errorf("Expected empty Env map, got %v", config.Env)
		}
	})

	t.Run("with invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		invalidYaml := `env:
  invalid yaml content
    - bad indentation
`
		if err := os.WriteFile(configPath, []byte(invalidYaml), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		_, err := LoadProjectConfig(tmpDir)
		if err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}
	})
}

func TestHasProjectConfig(t *testing.T) {
	t.Run("with camp.yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yaml")

		if err := os.WriteFile(configPath, []byte("env: {}"), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		if !HasProjectConfig(tmpDir) {
			t.Error("Expected HasProjectConfig to return true for directory with camp.yaml")
		}
	})

	t.Run("with camp.yml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "camp.yml")

		if err := os.WriteFile(configPath, []byte("env: {}"), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		if !HasProjectConfig(tmpDir) {
			t.Error("Expected HasProjectConfig to return true for directory with camp.yml")
		}
	})

	t.Run("without config file", func(t *testing.T) {
		tmpDir := t.TempDir()

		if HasProjectConfig(tmpDir) {
			t.Error("Expected HasProjectConfig to return false for directory without config")
		}
	})

	t.Run("with both files", func(t *testing.T) {
		tmpDir := t.TempDir()

		yamlPath := filepath.Join(tmpDir, "camp.yaml")
		ymlPath := filepath.Join(tmpDir, "camp.yml")

		if err := os.WriteFile(yamlPath, []byte("env: {}"), 0644); err != nil {
			t.Fatalf("Failed to write camp.yaml: %v", err)
		}
		if err := os.WriteFile(ymlPath, []byte("env: {}"), 0644); err != nil {
			t.Fatalf("Failed to write camp.yml: %v", err)
		}

		if !HasProjectConfig(tmpDir) {
			t.Error("Expected HasProjectConfig to return true for directory with both files")
		}
	})
}
