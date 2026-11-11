package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"camp/internal/system"
	"camp/internal/utils/utilstest"
)

func TestRebuildCommand(t *testing.T) {
	t.Run("command name", func(t *testing.T) {
		if rebuildCmd.Use != "rebuild" {
			t.Errorf("Expected command name to be 'rebuild', got '%s'", rebuildCmd.Use)
		}
	})

	t.Run("command description", func(t *testing.T) {
		expected := "Rebuild the development environment"
		if rebuildCmd.Short != expected {
			t.Errorf("Expected short description to be '%s', got '%s'", expected, rebuildCmd.Short)
		}
	})

	t.Run("command has long description", func(t *testing.T) {
		if rebuildCmd.Long == "" {
			t.Error("Expected long description to be set")
		}

		if !strings.Contains(rebuildCmd.Long, "Prerequisites") {
			t.Error("Expected long description to contain prerequisites")
		}
	})

	t.Run("command is subcommand of env", func(t *testing.T) {
		found := false
		for _, cmd := range envCmd.Commands() {
			if cmd.Use == "rebuild" {
				found = true
				break
			}
		}
		if !found {
			t.Error("rebuild command should be registered as subcommand of env")
		}
	})
}

func TestRebuildCommandExecution_Linux(t *testing.T) {
	exec := utilstest.StubExec(t)
	user := stubCurrentUser(t, "linux", "env:\n  TEST_VAR: test_value\n")

	output, err := executeCommand(rebuildCmd)
	if err != nil {
		t.Fatalf("rebuild returned error: %v", err)
	}

	if !strings.Contains(output, "✓ Environment prepared successfully") {
		t.Errorf("Expected output to show environment was prepared, got:\n%s", output)
	}
	if !strings.Contains(output, "✓ Environment rebuild completed successfully!") {
		t.Errorf("Expected output to show rebuild completed, got:\n%s", output)
	}

	nixDir := filepath.Join(user.HomeDir, ".camp", "nix")
	if _, err := os.Stat(filepath.Join(nixDir, "flake.nix")); err != nil {
		t.Errorf("Expected flake.nix to be created: %v", err)
	}

	if !exec.Called("home-manager", "switch", "--flake", nixDir+"#testuser") {
		t.Errorf("Expected home-manager switch to run, got:\n%s", exec)
	}
}

func TestRebuildCommandExecution_Darwin(t *testing.T) {
	exec := utilstest.StubExec(t)
	user := stubCurrentUser(t, "darwin", "env:\n  EDITOR: nvim\n")

	if _, err := executeCommand(rebuildCmd); err != nil {
		t.Fatalf("rebuild returned error: %v", err)
	}

	nixDir := filepath.Join(user.HomeDir, ".camp", "nix")
	if !exec.Called("sudo", "nix", "nix-darwin#darwin-rebuild", "switch", "--flake", nixDir+"#testhost") {
		t.Errorf("Expected darwin-rebuild switch to run, got:\n%s", exec)
	}
}

func TestRebuildCommandExecution_RefreshesLocalFlakes(t *testing.T) {
	exec := utilstest.StubExec(t)
	stubCurrentUser(t, "linux", `flakes:
  - name: dotfiles
    url: "path:/tmp/dotfiles"
    outputs:
      - name: homeManagerModules.default
        type: home
`)

	output, err := executeCommand(rebuildCmd)
	if err != nil {
		t.Fatalf("rebuild returned error: %v", err)
	}

	if !strings.Contains(output, "Refreshing local flakes: dotfiles") {
		t.Errorf("Expected output to mention refreshed local flakes, got:\n%s", output)
	}
	if len(exec.Calls) != 2 || !exec.Called("nix", "flake", "update", "dotfiles") || exec.Calls[1][0] != "home-manager" {
		t.Errorf("Expected local flake update followed by home-manager switch, got:\n%s", exec)
	}
}

func TestRebuildCommandExecution_RebuildFails(t *testing.T) {
	exec := utilstest.StubExec(t)
	exec.Err = errors.New("exit status 1")
	stubCurrentUser(t, "linux", "env:\n  EDITOR: nvim\n")

	output, err := executeCommand(rebuildCmd)
	if err == nil || !strings.Contains(err.Error(), "rebuild failed") {
		t.Errorf("Expected rebuild failure error, got: %v", err)
	}
	if strings.Contains(output, "completed successfully") {
		t.Errorf("Expected no success message on failure, got:\n%s", output)
	}
}

func TestRebuildCommandOutputFormat(t *testing.T) {
	utilstest.StubExec(t)
	stubCurrentUser(t, "linux", "env:\n  EDITOR: nvim\n")

	output, err := executeCommand(rebuildCmd)
	if err != nil {
		t.Fatalf("rebuild returned error: %v", err)
	}

	expectedPhrases := []string{
		"Starting environment rebuild",
		"Platform: linux",
		"User: testuser",
		"Hostname: testhost",
		"Preparing environment",
		"Executing rebuild command",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected output to contain '%s', got:\n%s", phrase, output)
		}
	}
}

func TestRebuildCommandMentionsLocalFlakes(t *testing.T) {
	if !strings.Contains(rebuildCmd.Long, "local flakes") {
		t.Error("Expected long description to mention local flakes refresh")
	}
}

func TestRefreshLocalFlakes_NoLocalFlakes(t *testing.T) {
	user := &system.User{
		Name:    "testuser",
		HomeDir: t.TempDir(),
		Flakes: []system.Flake{
			{Name: "remote", URL: "github:user/repo"},
		},
	}

	var output bytes.Buffer
	if err := refreshLocalFlakes(&output, user); err != nil {
		t.Fatalf("refreshLocalFlakes() returned error: %v", err)
	}

	if output.Len() != 0 {
		t.Errorf("Expected no output without local flakes, got:\n%s", output.String())
	}
}
