package cmd

import (
	"bytes"
	"strings"
	"testing"

	"camp/internal/system"
	"github.com/spf13/cobra"
)

func TestEnvCommand(t *testing.T) {
	t.Run("command name", func(t *testing.T) {
		if envCmd.Use != "env" {
			t.Errorf("Expected command name to be 'env', got '%s'", envCmd.Use)
		}
	})

	t.Run("command description", func(t *testing.T) {
		expected := "Display development environment information"
		if envCmd.Short != expected {
			t.Errorf("Expected short description to be '%s', got '%s'", expected, envCmd.Short)
		}
	})

	t.Run("command execution and output", func(t *testing.T) {
		var output bytes.Buffer

		cmd := &cobra.Command{
			Use:   envCmd.Use,
			Short: envCmd.Short,
			Long:  envCmd.Long,
			Run:   envCmd.Run,
		}

		cmd.SetOut(&output)
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		if err != nil {
			t.Errorf("Command execution failed: %v", err)
		}

		actualOutput := strings.TrimSpace(output.String())
		lines := strings.Split(actualOutput, "\n")

		if len(lines) < 2 {
			t.Errorf("Expected at least 2 lines of output, got %d", len(lines))
		}

		if !strings.HasPrefix(lines[0], "Architecture: ") {
			t.Errorf("Expected first line to start with 'Architecture: ', got '%s'", lines[0])
		}

		if !strings.HasPrefix(lines[1], "OS: ") {
			t.Errorf("Expected second line to start with 'OS: ', got '%s'", lines[1])
		}

		arch := strings.TrimPrefix(lines[0], "Architecture: ")
		if arch == "" {
			t.Error("Expected architecture value to not be empty")
		}

		os := strings.TrimPrefix(lines[1], "OS: ")
		if os == "" {
			t.Error("Expected OS value to not be empty")
		}
	})
}

func TestGetSystemInfo(t *testing.T) {
	sysInfo, err := system.GetSystemInfo()
	if err != nil {
		t.Errorf("GetSystemInfo() failed: %v", err)
	}

	if sysInfo == nil {
		t.Fatal("Expected system info to not be nil")
	}

	if sysInfo.Architecture == "" {
		t.Error("Expected architecture to not be empty")
	}

	if sysInfo.OS == "" {
		t.Error("Expected operating system to not be empty")
	}

	validArchs := []string{"x86_64", "arm64", "aarch64", "i386", "armv7l"}
	isValidArch := false
	for _, validArch := range validArchs {
		if sysInfo.Architecture == validArch {
			isValidArch = true
			break
		}
	}

	if !isValidArch {
		t.Logf("Got architecture: %s (may be valid but not in our test list)", sysInfo.Architecture)
	}

	validOS := []string{"darwin", "linux"}
	isValidOS := false
	for _, validOSName := range validOS {
		if sysInfo.OS == validOSName {
			isValidOS = true
			break
		}
	}

	if !isValidOS {
		t.Logf("Got OS: %s (may be valid but not in our test list)", sysInfo.OS)
	}
}
