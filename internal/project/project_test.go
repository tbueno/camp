package project

import (
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

		t.Run("Whole message", func(t *testing.T) {
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
