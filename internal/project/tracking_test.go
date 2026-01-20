package project

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestGetTrackingFilePath(t *testing.T) {
	path := GetTrackingFilePath("/home/user/project")
	expected := "/home/user/project/.camp/env_vars"
	if path != expected {
		t.Errorf("GetTrackingFilePath() = %q, want %q", path, expected)
	}
}

func TestLoadTrackedVars(t *testing.T) {
	t.Run("with existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		trackingDir := filepath.Join(tmpDir, ".camp")
		os.MkdirAll(trackingDir, 0755)

		trackingPath := filepath.Join(trackingDir, "env_vars")
		content := "FOO\nBAR\nBAZ\n"
		if err := os.WriteFile(trackingPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write tracking file: %v", err)
		}

		vars, err := LoadTrackedVars(tmpDir)
		if err != nil {
			t.Fatalf("LoadTrackedVars failed: %v", err)
		}

		expected := []string{"FOO", "BAR", "BAZ"}
		if !reflect.DeepEqual(vars, expected) {
			t.Errorf("LoadTrackedVars() = %v, want %v", vars, expected)
		}
	})

	t.Run("with non-existent file", func(t *testing.T) {
		tmpDir := t.TempDir()

		vars, err := LoadTrackedVars(tmpDir)
		if err != nil {
			t.Fatalf("LoadTrackedVars should not fail for missing file: %v", err)
		}

		if len(vars) != 0 {
			t.Errorf("Expected empty slice, got %v", vars)
		}
	})

	t.Run("with empty lines", func(t *testing.T) {
		tmpDir := t.TempDir()
		trackingDir := filepath.Join(tmpDir, ".camp")
		os.MkdirAll(trackingDir, 0755)

		trackingPath := filepath.Join(trackingDir, "env_vars")
		content := "FOO\n\nBAR\n  \nBAZ\n"
		if err := os.WriteFile(trackingPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write tracking file: %v", err)
		}

		vars, err := LoadTrackedVars(tmpDir)
		if err != nil {
			t.Fatalf("LoadTrackedVars failed: %v", err)
		}

		expected := []string{"FOO", "BAR", "BAZ"}
		if !reflect.DeepEqual(vars, expected) {
			t.Errorf("LoadTrackedVars() = %v, want %v", vars, expected)
		}
	})
}

func TestSaveTrackedVars(t *testing.T) {
	t.Run("creates directory and file", func(t *testing.T) {
		tmpDir := t.TempDir()

		vars := []string{"FOO", "BAR", "BAZ"}
		if err := SaveTrackedVars(tmpDir, vars); err != nil {
			t.Fatalf("SaveTrackedVars failed: %v", err)
		}

		// Verify file was created
		trackingPath := GetTrackingFilePath(tmpDir)
		if _, err := os.Stat(trackingPath); os.IsNotExist(err) {
			t.Error("Tracking file was not created")
		}

		// Verify content
		loaded, err := LoadTrackedVars(tmpDir)
		if err != nil {
			t.Fatalf("Failed to load saved vars: %v", err)
		}

		if !reflect.DeepEqual(loaded, vars) {
			t.Errorf("Loaded vars = %v, want %v", loaded, vars)
		}
	})

	t.Run("overwrites existing file", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Save initial vars
		if err := SaveTrackedVars(tmpDir, []string{"OLD1", "OLD2"}); err != nil {
			t.Fatalf("First SaveTrackedVars failed: %v", err)
		}

		// Save new vars
		newVars := []string{"NEW1", "NEW2", "NEW3"}
		if err := SaveTrackedVars(tmpDir, newVars); err != nil {
			t.Fatalf("Second SaveTrackedVars failed: %v", err)
		}

		// Verify new content
		loaded, err := LoadTrackedVars(tmpDir)
		if err != nil {
			t.Fatalf("Failed to load saved vars: %v", err)
		}

		if !reflect.DeepEqual(loaded, newVars) {
			t.Errorf("Loaded vars = %v, want %v", loaded, newVars)
		}
	})

	t.Run("with empty vars", func(t *testing.T) {
		tmpDir := t.TempDir()

		if err := SaveTrackedVars(tmpDir, []string{}); err != nil {
			t.Fatalf("SaveTrackedVars failed: %v", err)
		}

		loaded, err := LoadTrackedVars(tmpDir)
		if err != nil {
			t.Fatalf("Failed to load saved vars: %v", err)
		}

		if len(loaded) != 0 {
			t.Errorf("Expected empty slice, got %v", loaded)
		}
	})
}

func TestGetRemovedVars(t *testing.T) {
	t.Run("some vars removed", func(t *testing.T) {
		tracked := []string{"FOO", "BAR", "BAZ", "QUX"}
		current := map[string]string{
			"FOO": "value1",
			"QUX": "value2",
		}

		removed := GetRemovedVars(tracked, current)
		sort.Strings(removed)

		expected := []string{"BAR", "BAZ"}
		if !reflect.DeepEqual(removed, expected) {
			t.Errorf("GetRemovedVars() = %v, want %v", removed, expected)
		}
	})

	t.Run("no vars removed", func(t *testing.T) {
		tracked := []string{"FOO", "BAR"}
		current := map[string]string{
			"FOO": "value1",
			"BAR": "value2",
			"BAZ": "value3",
		}

		removed := GetRemovedVars(tracked, current)

		if len(removed) != 0 {
			t.Errorf("Expected no removed vars, got %v", removed)
		}
	})

	t.Run("all vars removed", func(t *testing.T) {
		tracked := []string{"FOO", "BAR"}
		current := map[string]string{}

		removed := GetRemovedVars(tracked, current)
		sort.Strings(removed)

		expected := []string{"BAR", "FOO"}
		if !reflect.DeepEqual(removed, expected) {
			t.Errorf("GetRemovedVars() = %v, want %v", removed, expected)
		}
	})

	t.Run("empty tracked", func(t *testing.T) {
		tracked := []string{}
		current := map[string]string{
			"FOO": "value1",
		}

		removed := GetRemovedVars(tracked, current)

		if len(removed) != 0 {
			t.Errorf("Expected no removed vars, got %v", removed)
		}
	})
}
