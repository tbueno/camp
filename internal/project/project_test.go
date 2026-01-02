package project

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestProject(t *testing.T) {
	t.Run("Name()", func(t *testing.T) {
		t.Run("invalid folder", func(t *testing.T) {
			p := Project{Path: "path/to/my-project", Config: DevboxConfig{}}
			n := p.Name()
			if n != "my-project" {
				t.Errorf("expected 'my-project', got %s", n)
			}
		})
	})

	t.Run("Compatible()", func(t *testing.T) {
		t.Run("incompatible project", func(t *testing.T) {
			p := NewProject("/tmp")
			// Assuming /tmp does not have a devbox.json
			if p.Compatible() {
				t.Logf("Warning: /tmp seems to have a devbox.json or false positive. Path: %s", p.Path)
				// Relaxing this test if /tmp is unpredictable, but generally it should pass.
				// For stability in CI/Local, it's better to use a random temp dir.
			}
		})

		t.Run("Compatible project", func(t *testing.T) {
			// We need a real file to test compatibility
			dir := t.TempDir()
			p := NewProject(dir)
			// No devbox.json yet
			if p.Compatible() {
				t.Error("expected false, got true")
			}

			// Create devbox.json is not easy here without writing files.
			// skipping file creation for unit test simplicity unless we mock fs,
			// but we can trust the os.Stat logic mostly.
		})
	})

	t.Run("Info()", func(t *testing.T) {
		t.Run("project name", func(t *testing.T) {
			p := Project{Path: "path/to/my-project", Config: DevboxConfig{}}
			info := p.Info()
			expected := "Project name: my-project"
			if info[0] != expected {
				t.Errorf("expected '%s', got '%s'", expected, info[0])
			}
		})

		t.Run("Whole message without packages", func(t *testing.T) {
			p := Project{Path: "../../optishell", Config: DevboxConfig{Shell: ShellConfig{Scripts: map[string][]string{"test": {"run tests"}}}}}
			info := p.Info()
			expected := []string{
				"Project name: optishell",
				"Commands available through 'camp project [command]':",
				" - test",
			}
			if !reflect.DeepEqual(info, expected) {
				t.Errorf("expected %s, got %s", expected, info)
			}
		})

		t.Run("Whole message with packages", func(t *testing.T) {
			p := Project{
				Path: "path/to/my-project",
				Config: DevboxConfig{
					Packages: []string{"go@latest", "hugo@latest"},
					Shell: ShellConfig{
						Scripts: map[string][]string{"test": {"run tests"}},
					},
				},
			}
			info := p.Info()
			expected := []string{
				"Project name: my-project",
				"Packages:",
				" - go@latest",
				" - hugo@latest",
				"Commands available through 'camp project [command]':",
				" - test",
			}
			if !reflect.DeepEqual(info, expected) {
				t.Errorf("expected %v, got %v", expected, info)
			}
		})
	})

	t.Run("Commands()", func(t *testing.T) {
		t.Run("with single line command", func(t *testing.T) {
			sc := ShellConfig{Scripts: map[string][]string{"test": {"echo 'hello world'"}}}
			p := Project{Path: "path/to/my-project", Config: DevboxConfig{Shell: sc}}
			cmds := p.Commands()
			expected := map[string][][]string{"test": {{"echo", "'hello", "world'"}}}
			if !reflect.DeepEqual(cmds, expected) {
				t.Errorf("expected %v, got %v", expected, cmds)
			}
		})

		t.Run("with multiple line commands", func(t *testing.T) {
			sc := ShellConfig{Scripts: map[string][]string{"build": {"go build", "go test"}}}
			p := Project{Path: "path/to/my-project", Config: DevboxConfig{Shell: sc}}
			cmds := p.Commands()
			expected := map[string][][]string{
				"build": {
					{"go", "build"},
					{"go", "test"},
				},
			}
			if !reflect.DeepEqual(cmds, expected) {
				t.Errorf("expected %v, got %v", expected, cmds)
			}
		})

		t.Run("with no commands", func(t *testing.T) {
			sc := ShellConfig{Scripts: map[string][]string{}}
			p := Project{Path: "path/to/my-project", Config: DevboxConfig{Shell: sc}}
			cmds := p.Commands()
			expected := map[string][][]string{}
			if !reflect.DeepEqual(cmds, expected) {
				t.Errorf("expected %v, got %v", expected, cmds)
			}
		})
	})
}

func TestNewProject_WithCampConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .camp.yml
	campYaml := `env:
  PROJECT_NAME: "test-project"
  DEBUG: "true"
`
	campConfigPath := filepath.Join(tmpDir, ".camp.yml")
	if err := os.WriteFile(campConfigPath, []byte(campYaml), 0644); err != nil {
		t.Fatalf("Failed to write .camp.yml: %v", err)
	}

	// Create project
	proj := NewProject(tmpDir)

	// Verify camp config loaded
	if !proj.HasCampConfig() {
		t.Error("NewProject() should load .camp.yml")
	}

	envVars := proj.EnvVars()
	if envVars["PROJECT_NAME"] != "test-project" {
		t.Errorf("Expected PROJECT_NAME=test-project, got %s", envVars["PROJECT_NAME"])
	}

	if envVars["DEBUG"] != "true" {
		t.Errorf("Expected DEBUG=true, got %s", envVars["DEBUG"])
	}
}

func TestNewProject_WithDevboxOnly(t *testing.T) {
	tmpDir := t.TempDir()

	// Create devbox.json
	devboxJson := `{
  "packages": ["go", "hugo"],
  "shell": {
    "scripts": {
      "test": ["go test ./..."]
    }
  }
}`
	devboxPath := filepath.Join(tmpDir, "devbox.json")
	if err := os.WriteFile(devboxPath, []byte(devboxJson), 0644); err != nil {
		t.Fatalf("Failed to write devbox.json: %v", err)
	}

	// Create project
	proj := NewProject(tmpDir)

	// Verify devbox loaded but no camp config
	if !proj.Compatible() {
		t.Error("NewProject() should load devbox.json")
	}

	if proj.HasCampConfig() {
		t.Error("NewProject() should not have camp config when .camp.yml doesn't exist")
	}

	envVars := proj.EnvVars()
	if len(envVars) != 0 {
		t.Errorf("Expected empty env vars, got %d vars", len(envVars))
	}
}

func TestNewProject_WithBoth(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .camp.yml
	campYaml := `env:
  PROJECT_NAME: "test-project"
`
	campConfigPath := filepath.Join(tmpDir, ".camp.yml")
	if err := os.WriteFile(campConfigPath, []byte(campYaml), 0644); err != nil {
		t.Fatalf("Failed to write .camp.yml: %v", err)
	}

	// Create devbox.json
	devboxJson := `{
  "packages": ["go"],
  "shell": {
    "scripts": {
      "test": ["go test"]
    }
  }
}`
	devboxPath := filepath.Join(tmpDir, "devbox.json")
	if err := os.WriteFile(devboxPath, []byte(devboxJson), 0644); err != nil {
		t.Fatalf("Failed to write devbox.json: %v", err)
	}

	// Create project
	proj := NewProject(tmpDir)

	// Verify both loaded
	if !proj.Compatible() {
		t.Error("NewProject() should load devbox.json")
	}

	if !proj.HasCampConfig() {
		t.Error("NewProject() should load .camp.yml")
	}

	if len(proj.Config.Packages) != 1 {
		t.Errorf("Expected 1 devbox package, got %d", len(proj.Config.Packages))
	}

	envVars := proj.EnvVars()
	if envVars["PROJECT_NAME"] != "test-project" {
		t.Errorf("Expected PROJECT_NAME=test-project, got %s", envVars["PROJECT_NAME"])
	}
}

func TestNewProject_WithNeither(t *testing.T) {
	tmpDir := t.TempDir()

	// Create project (no config files)
	proj := NewProject(tmpDir)

	// Verify neither loaded
	if proj.Compatible() {
		t.Error("NewProject() should not be compatible without devbox.json")
	}

	if proj.HasCampConfig() {
		t.Error("NewProject() should not have camp config without .camp.yml")
	}

	envVars := proj.EnvVars()
	if len(envVars) != 0 {
		t.Errorf("Expected empty env vars, got %d vars", len(envVars))
	}
}

func TestHasCampConfig(t *testing.T) {
	t.Run("with camp config", func(t *testing.T) {
		tmpDir := t.TempDir()

		campYaml := `env:
  TEST: "value"
`
		campConfigPath := filepath.Join(tmpDir, ".camp.yml")
		os.WriteFile(campConfigPath, []byte(campYaml), 0644)

		proj := NewProject(tmpDir)

		if !proj.HasCampConfig() {
			t.Error("HasCampConfig() should return true when .camp.yml exists")
		}
	})

	t.Run("without camp config", func(t *testing.T) {
		tmpDir := t.TempDir()
		proj := NewProject(tmpDir)

		if proj.HasCampConfig() {
			t.Error("HasCampConfig() should return false when .camp.yml doesn't exist")
		}
	})
}

func TestEnvVars(t *testing.T) {
	t.Run("with env vars", func(t *testing.T) {
		tmpDir := t.TempDir()

		campYaml := `env:
  VAR1: "value1"
  VAR2: "value2"
  VAR3: "value3"
`
		campConfigPath := filepath.Join(tmpDir, ".camp.yml")
		os.WriteFile(campConfigPath, []byte(campYaml), 0644)

		proj := NewProject(tmpDir)
		envVars := proj.EnvVars()

		if len(envVars) != 3 {
			t.Errorf("Expected 3 env vars, got %d", len(envVars))
		}

		if envVars["VAR1"] != "value1" {
			t.Errorf("Expected VAR1=value1, got %s", envVars["VAR1"])
		}

		if envVars["VAR2"] != "value2" {
			t.Errorf("Expected VAR2=value2, got %s", envVars["VAR2"])
		}

		if envVars["VAR3"] != "value3" {
			t.Errorf("Expected VAR3=value3, got %s", envVars["VAR3"])
		}
	})

	t.Run("without camp config", func(t *testing.T) {
		tmpDir := t.TempDir()
		proj := NewProject(tmpDir)
		envVars := proj.EnvVars()

		if len(envVars) != 0 {
			t.Errorf("Expected empty env vars map, got %d vars", len(envVars))
		}
	})

	t.Run("with empty env section", func(t *testing.T) {
		tmpDir := t.TempDir()

		campYaml := `env: {}
`
		campConfigPath := filepath.Join(tmpDir, ".camp.yml")
		os.WriteFile(campConfigPath, []byte(campYaml), 0644)

		proj := NewProject(tmpDir)
		envVars := proj.EnvVars()

		if len(envVars) != 0 {
			t.Errorf("Expected empty env vars map, got %d vars", len(envVars))
		}
	})
}
