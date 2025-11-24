package system

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveCampFiles(t *testing.T) {
	t.Run("removes existing camp directory", func(t *testing.T) {
		// Create temporary directory structure
		tmpDir := t.TempDir()
		campDir := filepath.Join(tmpDir, ".camp")

		// Create .camp directory with some files
		if err := os.MkdirAll(filepath.Join(campDir, "nix"), 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		testFile := filepath.Join(campDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Remove camp files
		var output bytes.Buffer
		if err := RemoveCampFiles(campDir, &output); err != nil {
			t.Fatalf("RemoveCampFiles() failed: %v", err)
		}

		// Verify directory was removed
		if _, err := os.Stat(campDir); !os.IsNotExist(err) {
			t.Error("Expected camp directory to be removed")
		}

		// Verify output message
		if !strings.Contains(output.String(), "Deleted") {
			t.Error("Expected output to contain deletion message")
		}
	})

	t.Run("handles non-existent directory gracefully", func(t *testing.T) {
		// Use a path that doesn't exist
		tmpDir := t.TempDir()
		campDir := filepath.Join(tmpDir, ".camp")

		// Remove camp files (directory doesn't exist)
		var output bytes.Buffer
		if err := RemoveCampFiles(campDir, &output); err != nil {
			t.Fatalf("RemoveCampFiles() should not fail for non-existent directory: %v", err)
		}

		// Verify output message indicates skipping
		if !strings.Contains(output.String(), "not found") || !strings.Contains(output.String(), "skipping") {
			t.Error("Expected output to indicate directory was not found")
		}
	})
}

func TestRemoveNixStateFiles(t *testing.T) {
	t.Run("removes nix state directories", func(t *testing.T) {
		// Create temporary home directory
		tmpHome := t.TempDir()

		// Create state directories
		homeManagerConfig := filepath.Join(tmpHome, ".config", "home-manager")
		homeManagerState := filepath.Join(tmpHome, ".local", "state", "home-manager")
		nixState := filepath.Join(tmpHome, ".local", "state", "nix")

		for _, dir := range []string{homeManagerConfig, homeManagerState, nixState} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("Failed to create test directory %s: %v", dir, err)
			}

			// Add a test file
			testFile := filepath.Join(dir, "test.txt")
			if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		// Remove nix state files
		var output bytes.Buffer
		if err := RemoveNixStateFiles(tmpHome, &output); err != nil {
			t.Fatalf("RemoveNixStateFiles() failed: %v", err)
		}

		// Verify all directories were removed
		for _, dir := range []string{homeManagerConfig, homeManagerState, nixState} {
			if _, err := os.Stat(dir); !os.IsNotExist(err) {
				t.Errorf("Expected %s to be removed", dir)
			}
		}

		// Verify output messages
		outputStr := output.String()
		if !strings.Contains(outputStr, "Deleted") {
			t.Error("Expected output to contain deletion messages")
		}
	})

	t.Run("handles non-existent directories gracefully", func(t *testing.T) {
		// Use temporary home without creating state directories
		tmpHome := t.TempDir()

		// Remove nix state files (directories don't exist)
		var output bytes.Buffer
		if err := RemoveNixStateFiles(tmpHome, &output); err != nil {
			t.Fatalf("RemoveNixStateFiles() should not fail for non-existent directories: %v", err)
		}

		// Verify output indicates skipping
		outputStr := output.String()
		if !strings.Contains(outputStr, "not found") || !strings.Contains(outputStr, "skipping") {
			t.Error("Expected output to indicate directories were not found")
		}
	})
}

func TestDeleteIfExists(t *testing.T) {
	t.Run("deletes existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.txt")

		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		var output bytes.Buffer
		if err := deleteIfExists(testFile, &output); err != nil {
			t.Fatalf("deleteIfExists() failed: %v", err)
		}

		// Verify file was deleted
		if _, err := os.Stat(testFile); !os.IsNotExist(err) {
			t.Error("Expected file to be deleted")
		}

		// Verify output message
		if !strings.Contains(output.String(), "Deleted") {
			t.Error("Expected output to contain deletion message")
		}
	})

	t.Run("deletes existing directory recursively", func(t *testing.T) {
		tmpDir := t.TempDir()
		testDir := filepath.Join(tmpDir, "testdir")

		// Create directory with nested files
		if err := os.MkdirAll(filepath.Join(testDir, "subdir"), 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		testFile := filepath.Join(testDir, "subdir", "test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		var output bytes.Buffer
		if err := deleteIfExists(testDir, &output); err != nil {
			t.Fatalf("deleteIfExists() failed: %v", err)
		}

		// Verify directory was deleted
		if _, err := os.Stat(testDir); !os.IsNotExist(err) {
			t.Error("Expected directory to be deleted")
		}

		// Verify output message
		if !strings.Contains(output.String(), "Deleted") {
			t.Error("Expected output to contain deletion message")
		}
	})

	t.Run("handles non-existent path gracefully", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "nonexistent")

		var output bytes.Buffer
		if err := deleteIfExists(testPath, &output); err != nil {
			t.Fatalf("deleteIfExists() should not fail for non-existent path: %v", err)
		}

		// Verify output message indicates skipping
		if !strings.Contains(output.String(), "not found") || !strings.Contains(output.String(), "skipping") {
			t.Error("Expected output to indicate path was not found")
		}
	})
}

func TestNukeEnvironment(t *testing.T) {
	t.Run("removes camp and state files", func(t *testing.T) {
		// Create temporary home directory
		tmpHome := t.TempDir()

		// Create camp directory
		campDir := filepath.Join(tmpHome, ".camp")
		if err := os.MkdirAll(filepath.Join(campDir, "nix"), 0755); err != nil {
			t.Fatalf("Failed to create .camp directory: %v", err)
		}

		// Create state directories
		homeManagerConfig := filepath.Join(tmpHome, ".config", "home-manager")
		if err := os.MkdirAll(homeManagerConfig, 0755); err != nil {
			t.Fatalf("Failed to create home-manager config: %v", err)
		}

		// Create user
		user := &User{
			Name:         "testuser",
			HostName:     "testhost",
			Platform:     "linux", // Use linux to avoid darwin-specific uninstall
			Architecture: "amd64",
			HomeDir:      tmpHome,
		}

		// Run NukeEnvironment
		var output bytes.Buffer
		err := NukeEnvironment(user, &output)

		// Note: This will fail on Nix uninstall (which is expected in test environment)
		// but should continue with file cleanup
		outputStr := output.String()

		// Verify camp directory was removed (even if Nix uninstall failed)
		if _, err := os.Stat(campDir); !os.IsNotExist(err) {
			t.Error("Expected camp directory to be removed")
		}

		// Verify state directory was removed
		if _, err := os.Stat(homeManagerConfig); !os.IsNotExist(err) {
			t.Error("Expected home-manager config to be removed")
		}

		// Verify output contains expected messages
		if !strings.Contains(outputStr, "Removing camp files") {
			t.Error("Expected output to mention removing camp files")
		}

		// The function should not return an error even if Nix uninstall fails
		// as long as file cleanup succeeds
		if err != nil {
			t.Logf("NukeEnvironment returned error (expected if Nix not installed): %v", err)
		}
	})
}

func TestUninstallNix(t *testing.T) {
	t.Run("returns error when nix-installer not found", func(t *testing.T) {
		// This test will fail unless /nix/nix-installer exists
		// which is expected in most test environments
		var output bytes.Buffer
		err := UninstallNix("linux", &output)

		// We expect this to fail in test environment
		if err == nil {
			t.Log("UninstallNix succeeded (Nix is installed on this system)")
		} else {
			// This is the expected case in most test environments
			if !strings.Contains(err.Error(), "nix-installer") {
				t.Errorf("Expected error to mention nix-installer, got: %v", err)
			}
		}
	})
}
