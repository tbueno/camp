package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"camp/internal/utils/utilstest"
)

func TestUpdateCommand(t *testing.T) {
	t.Run("command name", func(t *testing.T) {
		if updateCmd.Use != "update" {
			t.Errorf("Expected command name to be 'update', got '%s'", updateCmd.Use)
		}
	})

	t.Run("command description", func(t *testing.T) {
		expected := "Update flake dependencies to their latest versions"
		if updateCmd.Short != expected {
			t.Errorf("Expected short description to be '%s', got '%s'", expected, updateCmd.Short)
		}
	})

	t.Run("command has long description", func(t *testing.T) {
		if updateCmd.Long == "" {
			t.Error("Expected long description to be set")
		}

		if !strings.Contains(updateCmd.Long, "flake.lock") {
			t.Error("Expected long description to mention flake.lock")
		}

		if !strings.Contains(updateCmd.Long, "Prerequisites") {
			t.Error("Expected long description to contain prerequisites")
		}
	})

	t.Run("command is subcommand of env", func(t *testing.T) {
		found := false
		for _, cmd := range envCmd.Commands() {
			if cmd.Use == "update" {
				found = true
				break
			}
		}
		if !found {
			t.Error("update command should be registered as subcommand of env")
		}
	})
}

func TestUpdateCommandExecution(t *testing.T) {
	exec := utilstest.StubExec(t)
	user := stubCurrentUser(t, "linux", `env:
  EDITOR: nvim

flakes:
  - name: test-flake
    url: "github:test/flake"
    outputs:
      - name: packages
        type: home
`)

	output, err := executeCommand(updateCmd)
	if err != nil {
		t.Fatalf("update returned error: %v", err)
	}

	if !strings.Contains(output, "✓ Environment prepared successfully") {
		t.Errorf("Expected output to show environment was prepared, got:\n%s", output)
	}
	if !strings.Contains(output, "✓ Flake dependencies updated successfully!") {
		t.Errorf("Expected output to show flakes were updated, got:\n%s", output)
	}

	nixDir := filepath.Join(user.HomeDir, ".camp", "nix")
	if _, err := os.Stat(filepath.Join(nixDir, "flake.nix")); err != nil {
		t.Errorf("Expected flake.nix to be created: %v", err)
	}

	if len(exec.Calls) != 1 || !exec.Called("nix", "flake", "update", "--flake", nixDir) {
		t.Errorf("Expected a single nix flake update for %s, got:\n%s", nixDir, exec)
	}
}

func TestUpdateCommandExecution_UpdateFails(t *testing.T) {
	exec := utilstest.StubExec(t)
	exec.Err = errors.New("exit status 1")
	stubCurrentUser(t, "linux", "env:\n  EDITOR: nvim\n")

	output, err := executeCommand(updateCmd)
	if err == nil || !strings.Contains(err.Error(), "nix flake update failed") {
		t.Errorf("Expected nix flake update failure error, got: %v", err)
	}
	if strings.Contains(output, "updated successfully") {
		t.Errorf("Expected no success message on failure, got:\n%s", output)
	}
}

func TestUpdateCommandOutputFormat(t *testing.T) {
	utilstest.StubExec(t)
	user := stubCurrentUser(t, "linux", "env:\n  EDITOR: nvim\n")

	output, err := executeCommand(updateCmd)
	if err != nil {
		t.Fatalf("update returned error: %v", err)
	}

	expectedPhrases := []string{
		"Starting flake update",
		"User: testuser",
		"Nix directory: " + user.HomeDir + "/.camp/nix",
		"Preparing environment",
		"Updating flake dependencies",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected output to contain '%s', got:\n%s", phrase, output)
		}
	}
}

func TestUpdateCommandMentionsAutomaticLocalRefresh(t *testing.T) {
	if !strings.Contains(updateCmd.Long, "local flakes") {
		t.Error("Expected long description to explain local flakes are refreshed on rebuild")
	}
}
