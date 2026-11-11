package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"camp/internal/system"

	"github.com/spf13/cobra"
)

// stubCurrentUser makes the env commands act on a user whose home is a temporary
// directory holding the given camp.yml.
func stubCurrentUser(t *testing.T, platform, config string) *system.User {
	t.Helper()

	tmpHome := t.TempDir()
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(campDir, "camp.yml"), []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	user := &system.User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     platform,
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
	}

	original := currentUser
	currentUser = func() *system.User { return user }
	t.Cleanup(func() { currentUser = original })

	return user
}

// executeCommand runs a copy of the given command and returns its stdout
func executeCommand(source *cobra.Command) (string, error) {
	var output bytes.Buffer
	cmd := &cobra.Command{
		Use:  source.Use,
		RunE: source.RunE,
	}
	cmd.SetOut(&output)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	return output.String(), err
}
