package system

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"camp/internal/utils/utilstest"
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

func TestRemoveCampEnvrc(t *testing.T) {
	t.Run("removes the envrc created by camp bootstrap", func(t *testing.T) {
		tmpHome := t.TempDir()
		user := &User{HomeDir: tmpHome}
		if err := createEnvrc(user, filepath.Join(tmpHome, ".camp"), &bytes.Buffer{}, false); err != nil {
			t.Fatalf("createEnvrc() failed: %v", err)
		}

		if err := RemoveCampEnvrc(tmpHome, &bytes.Buffer{}); err != nil {
			t.Fatalf("RemoveCampEnvrc() failed: %v", err)
		}

		if _, err := os.Stat(filepath.Join(tmpHome, ".envrc")); !os.IsNotExist(err) {
			t.Error("Expected camp's .envrc to be removed")
		}
	})

	t.Run("keeps an envrc not created by camp", func(t *testing.T) {
		tmpHome := t.TempDir()
		envrcPath := filepath.Join(tmpHome, ".envrc")
		userContent := "use nix\nexport FOO=bar\n"
		if err := os.WriteFile(envrcPath, []byte(userContent), 0644); err != nil {
			t.Fatalf("Failed to write .envrc: %v", err)
		}

		var output bytes.Buffer
		if err := RemoveCampEnvrc(tmpHome, &output); err != nil {
			t.Fatalf("RemoveCampEnvrc() failed: %v", err)
		}

		content, err := os.ReadFile(envrcPath)
		if err != nil || string(content) != userContent {
			t.Errorf("Expected user's .envrc to be kept unchanged, got %q (err: %v)", content, err)
		}
		if !strings.Contains(output.String(), "not created by camp, keeping") {
			t.Errorf("Expected output to explain the .envrc was kept, got:\n%s", output.String())
		}
	})

	t.Run("succeeds when there is no envrc", func(t *testing.T) {
		if err := RemoveCampEnvrc(t.TempDir(), &bytes.Buffer{}); err != nil {
			t.Errorf("RemoveCampEnvrc() should succeed without .envrc, got: %v", err)
		}
	})
}

func TestNukeEnvironment(t *testing.T) {
	t.Run("uninstalls nix and removes camp and state files", func(t *testing.T) {
		exec := utilstest.StubExec(t)
		tmpHome := t.TempDir()

		campDir := filepath.Join(tmpHome, ".camp")
		if err := os.MkdirAll(filepath.Join(campDir, "nix"), 0755); err != nil {
			t.Fatalf("Failed to create .camp directory: %v", err)
		}

		homeManagerConfig := filepath.Join(tmpHome, ".config", "home-manager")
		if err := os.MkdirAll(homeManagerConfig, 0755); err != nil {
			t.Fatalf("Failed to create home-manager config: %v", err)
		}

		envrcPath := filepath.Join(tmpHome, ".envrc")
		if err := os.WriteFile(envrcPath, []byte(envrcMarker+"\nuse flake "+campDir+"/nix\n"), 0644); err != nil {
			t.Fatalf("Failed to write .envrc: %v", err)
		}

		user := &User{
			Name:         "testuser",
			HostName:     "testhost",
			Platform:     "linux",
			Architecture: "amd64",
			HomeDir:      tmpHome,
		}

		var output bytes.Buffer
		if err := NukeEnvironment(user, &output, true); err != nil {
			t.Fatalf("NukeEnvironment() failed: %v", err)
		}

		if !exec.Called("/nix/nix-installer", "uninstall") {
			t.Errorf("Expected nix-installer uninstall to run, got:\n%s", exec)
		}
		if _, err := os.Stat(campDir); !os.IsNotExist(err) {
			t.Error("Expected camp directory to be removed")
		}
		if _, err := os.Stat(homeManagerConfig); !os.IsNotExist(err) {
			t.Error("Expected home-manager config to be removed")
		}
		if _, err := os.Stat(envrcPath); !os.IsNotExist(err) {
			t.Error("Expected camp's .envrc to be removed")
		}
		if !strings.Contains(output.String(), "✓ Nix uninstalled successfully") {
			t.Errorf("Expected output to report nix uninstalled, got:\n%s", output.String())
		}
	})

	t.Run("continues with file cleanup when nix uninstall fails", func(t *testing.T) {
		exec := utilstest.StubExec(t)
		exec.Err = errors.New("exit status 1")
		tmpHome := t.TempDir()

		campDir := filepath.Join(tmpHome, ".camp")
		if err := os.MkdirAll(campDir, 0755); err != nil {
			t.Fatalf("Failed to create .camp directory: %v", err)
		}

		user := &User{Name: "testuser", Platform: "linux", HomeDir: tmpHome}

		var output bytes.Buffer
		if err := NukeEnvironment(user, &output, true); err != nil {
			t.Fatalf("NukeEnvironment() should not fail when only the nix uninstall fails: %v", err)
		}

		if !strings.Contains(output.String(), "Continuing with file cleanup") {
			t.Errorf("Expected output to mention continuing cleanup, got:\n%s", output.String())
		}
		if _, err := os.Stat(campDir); !os.IsNotExist(err) {
			t.Error("Expected camp directory to be removed")
		}
	})

	t.Run("only removes leftover files when nix is not installed", func(t *testing.T) {
		exec := utilstest.StubExec(t)
		tmpHome := t.TempDir()

		campDir := filepath.Join(tmpHome, ".camp")
		if err := os.MkdirAll(campDir, 0755); err != nil {
			t.Fatalf("Failed to create .camp directory: %v", err)
		}

		user := &User{Name: "testuser", Platform: "darwin", HomeDir: tmpHome}

		var output bytes.Buffer
		if err := NukeEnvironment(user, &output, false); err != nil {
			t.Fatalf("NukeEnvironment() failed: %v", err)
		}

		if len(exec.Calls) != 0 {
			t.Errorf("Expected no uninstall commands without nix, ran:\n%s", exec)
		}
		if !strings.Contains(output.String(), "Nix is not installed, skipping Nix uninstall") {
			t.Errorf("Expected output to mention skipping the uninstall, got:\n%s", output.String())
		}
		if _, err := os.Stat(campDir); !os.IsNotExist(err) {
			t.Error("Expected camp directory to be removed")
		}
	})

	t.Run("succeeds on a clean system", func(t *testing.T) {
		utilstest.StubExec(t)
		user := &User{Name: "testuser", Platform: "linux", HomeDir: t.TempDir()}

		if err := NukeEnvironment(user, &bytes.Buffer{}, false); err != nil {
			t.Errorf("NukeEnvironment() should succeed with nothing to remove, got: %v", err)
		}
	})
}

func TestUninstallNix(t *testing.T) {
	t.Run("linux runs only the nix installer uninstall", func(t *testing.T) {
		exec := utilstest.StubExec(t)

		if err := UninstallNix("linux", &bytes.Buffer{}); err != nil {
			t.Fatalf("UninstallNix() failed: %v", err)
		}

		want := [][]string{{"/nix/nix-installer", "uninstall"}}
		if !slices.EqualFunc(exec.Calls, want, slices.Equal) {
			t.Errorf("UninstallNix() ran %v, want %v", exec.Calls, want)
		}
	})

	t.Run("darwin uninstalls nix-darwin first", func(t *testing.T) {
		exec := utilstest.StubExec(t)

		if err := UninstallNix("darwin", &bytes.Buffer{}); err != nil {
			t.Fatalf("UninstallNix() failed: %v", err)
		}

		want := [][]string{
			{"sudo", "darwin-uninstaller"},
			{"/nix/nix-installer", "uninstall"},
		}
		if !slices.EqualFunc(exec.Calls, want, slices.Equal) {
			t.Errorf("UninstallNix() ran %v, want %v", exec.Calls, want)
		}
	})

	t.Run("returns error when nix-installer fails", func(t *testing.T) {
		exec := utilstest.StubExec(t)
		exec.Err = errors.New("exit status 1")

		var output bytes.Buffer
		err := UninstallNix("darwin", &output)
		if err == nil || !strings.Contains(err.Error(), "nix-installer uninstall failed") {
			t.Errorf("Expected nix-installer failure, got: %v", err)
		}
		if !strings.Contains(output.String(), "darwin-uninstaller not found or failed") {
			t.Errorf("Expected darwin-uninstaller failure to be reported as a warning, got:\n%s", output.String())
		}
	})
}
