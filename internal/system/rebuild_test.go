package system

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"camp/internal/utils/utilstest"
	"camp/templates"
)

func TestPrepareEnvironment(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create .camp directory and config
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	configPath := filepath.Join(campDir, "camp.yml")
	configContent := `env:
  TEST_VAR: test_value
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create user
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
	}

	// Run PrepareEnvironment
	if err := PrepareEnvironment(templates.FS, user); err != nil {
		t.Fatalf("PrepareEnvironment() failed: %v", err)
	}

	// Verify .camp/nix directory was created
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if _, err := os.Stat(nixDir); os.IsNotExist(err) {
		t.Error("PrepareEnvironment() should create .camp/nix directory")
	}

	// Verify config files were copied
	macNixPath := filepath.Join(nixDir, "mac.nix")
	if _, err := os.Stat(macNixPath); os.IsNotExist(err) {
		t.Error("PrepareEnvironment() should copy mac.nix")
	}

	// Verify modules directory was copied
	modulesPath := filepath.Join(nixDir, "modules")
	if _, err := os.Stat(modulesPath); os.IsNotExist(err) {
		t.Error("PrepareEnvironment() should copy modules directory")
	}

	// Verify flake.nix was rendered
	flakePath := filepath.Join(nixDir, "flake.nix")
	if _, err := os.Stat(flakePath); os.IsNotExist(err) {
		t.Error("PrepareEnvironment() should render flake.nix")
	}
}

func TestCopyConfigFiles(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create user
	user := &User{
		Name:    "testuser",
		HomeDir: tmpHome,
	}

	// Run CopyConfigFiles
	if err := CopyConfigFiles(templates.FS, user); err != nil {
		t.Fatalf("CopyConfigFiles() failed: %v", err)
	}

	nixDir := filepath.Join(tmpHome, ".camp", "nix")

	// Verify mac.nix was copied
	macNixPath := filepath.Join(nixDir, "mac.nix")
	if _, err := os.Stat(macNixPath); os.IsNotExist(err) {
		t.Error("CopyConfigFiles() should copy mac.nix")
	}

	// Verify linux.nix was copied
	linuxNixPath := filepath.Join(nixDir, "linux.nix")
	if _, err := os.Stat(linuxNixPath); os.IsNotExist(err) {
		t.Error("CopyConfigFiles() should copy linux.nix")
	}

	// Verify modules directory was copied
	modulesPath := filepath.Join(nixDir, "modules", "common.nix")
	if _, err := os.Stat(modulesPath); os.IsNotExist(err) {
		t.Error("CopyConfigFiles() should copy modules/common.nix")
	}

	// Verify flake.nix was NOT copied (it should be rendered)
	flakePath := filepath.Join(nixDir, "flake.nix")
	if _, err := os.Stat(flakePath); err == nil {
		t.Error("CopyConfigFiles() should not copy flake.nix (it's rendered separately)")
	}
}

func TestCompileTemplates(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create .camp directory and config
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	configPath := filepath.Join(campDir, "camp.yml")
	configContent := `env:
  EDITOR: nvim
  BROWSER: firefox
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create user
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
	}

	// Run CompileTemplates
	if err := CompileTemplates(templates.FS, user); err != nil {
		t.Fatalf("CompileTemplates() failed: %v", err)
	}

	// Verify flake.nix was created
	flakePath := filepath.Join(tmpHome, ".camp", "nix", "flake.nix")
	if _, err := os.Stat(flakePath); os.IsNotExist(err) {
		t.Error("CompileTemplates() should create flake.nix")
	}

	// Verify user's EnvVars were reloaded
	if user.EnvVars["EDITOR"] != "nvim" {
		t.Error("CompileTemplates() should reload user config")
	}
}

func TestExecuteRebuild_Darwin(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create nix directory
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if err := os.MkdirAll(nixDir, 0755); err != nil {
		t.Fatalf("Failed to create nix directory: %v", err)
	}

	// Create user
	user := &User{
		Name:     "testuser",
		HostName: "testhost",
		Platform: "darwin",
		HomeDir:  tmpHome,
	}

	exec := utilstest.StubExec(t)

	if err := ExecuteRebuild(user); err != nil {
		t.Fatalf("ExecuteRebuild() failed: %v", err)
	}

	last := exec.Calls[len(exec.Calls)-1]
	want := []string{
		"sudo", "nix", "--extra-experimental-features", "nix-command flakes",
		"run", "nix-darwin#darwin-rebuild", "--", "switch", "--impure",
		"--flake", nixDir + "#testhost",
	}
	if !slices.Equal(last, want) {
		t.Errorf("ExecuteRebuild() ran %v, want %v", last, want)
	}

	for _, call := range exec.Calls[:len(exec.Calls)-1] {
		if !slices.Contains(darwinEtcFiles, call[2]) || call[1] != "mv" {
			t.Errorf("Expected only /etc backups before the rebuild, got %v", call)
		}
	}
}

func TestExecuteRebuild_Linux(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create nix directory
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if err := os.MkdirAll(nixDir, 0755); err != nil {
		t.Fatalf("Failed to create nix directory: %v", err)
	}

	// Create user
	user := &User{
		Name:     "testuser",
		HostName: "testhost",
		Platform: "linux",
		HomeDir:  tmpHome,
	}

	exec := utilstest.StubExec(t)

	if err := ExecuteRebuild(user); err != nil {
		t.Fatalf("ExecuteRebuild() failed: %v", err)
	}

	want := [][]string{{
		"home-manager", "switch", "--impure", "-b", "backup",
		"--flake", nixDir + "#testuser",
	}}
	if !slices.EqualFunc(exec.Calls, want, slices.Equal) {
		t.Errorf("ExecuteRebuild() ran %v, want %v", exec.Calls, want)
	}
}

func TestExecuteRebuild_CommandFails(t *testing.T) {
	tmpHome := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpHome, ".camp", "nix"), 0755); err != nil {
		t.Fatalf("Failed to create nix directory: %v", err)
	}
	user := &User{Name: "testuser", Platform: "linux", HomeDir: tmpHome}

	exec := utilstest.StubExec(t)
	exec.Err = errors.New("exit status 1")

	err := ExecuteRebuild(user)
	if err == nil || !strings.Contains(err.Error(), "rebuild command failed") {
		t.Errorf("Expected rebuild command failure, got: %v", err)
	}
}

func TestExecuteRebuild_UnsupportedPlatform(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create nix directory
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if err := os.MkdirAll(nixDir, 0755); err != nil {
		t.Fatalf("Failed to create nix directory: %v", err)
	}

	// Create user with unsupported platform
	user := &User{
		Name:     "testuser",
		HostName: "testhost",
		Platform: "windows",
		HomeDir:  tmpHome,
	}

	// Should error for unsupported platform
	err := ExecuteRebuild(user)
	if err == nil {
		t.Error("ExecuteRebuild() should error for unsupported platform")
	}

	if err != nil && err.Error() != "unsupported platform: windows" {
		t.Errorf("Expected 'unsupported platform' error, got: %v", err)
	}
}

func TestExecuteRebuild_NoNixDirectory(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Don't create nix directory

	// Create user
	user := &User{
		Name:     "testuser",
		HostName: "testhost",
		Platform: "darwin",
		HomeDir:  tmpHome,
	}

	// Should error when nix directory doesn't exist
	err := ExecuteRebuild(user)
	if err == nil {
		t.Error("ExecuteRebuild() should error when nix directory doesn't exist")
	}
}

func TestBackupDarwinEtcFiles(t *testing.T) {
	// Verify the list of files to back up
	expected := []string{
		"/etc/bashrc",
		"/etc/zshrc",
		"/etc/zprofile",
	}

	if len(darwinEtcFiles) != len(expected) {
		t.Fatalf("Expected %d files, got %d", len(expected), len(darwinEtcFiles))
	}

	for i, f := range expected {
		if darwinEtcFiles[i] != f {
			t.Errorf("Expected darwinEtcFiles[%d] = %q, got %q", i, f, darwinEtcFiles[i])
		}
	}
}

func TestCopyConfigFiles_NestedDirectories(t *testing.T) {
	tmpHome := t.TempDir()
	user := &User{Name: "testuser", HomeDir: tmpHome}
	templDir := fstest.MapFS{
		"files/flake.nix":              {Data: []byte("rendered separately")},
		"files/mac.nix":                {Data: []byte("mac")},
		"files/modules/common.nix":     {Data: []byte("common")},
		"files/modules/extra/deep.nix": {Data: []byte("deep")},
	}

	if err := CopyConfigFiles(templDir, user); err != nil {
		t.Fatalf("CopyConfigFiles() failed: %v", err)
	}

	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	expected := map[string]string{
		"mac.nix":                "mac",
		"modules/common.nix":     "common",
		"modules/extra/deep.nix": "deep",
	}
	for name, want := range expected {
		got, err := os.ReadFile(filepath.Join(nixDir, name))
		if err != nil {
			t.Errorf("CopyConfigFiles() should copy %s: %v", name, err)
			continue
		}
		if string(got) != want {
			t.Errorf("%s content = %q, want %q", name, got, want)
		}
	}

	if _, err := os.Stat(filepath.Join(nixDir, "flake.nix")); err == nil {
		t.Error("CopyConfigFiles() should not copy flake.nix")
	}
}

func TestCopyConfigFiles_MissingTemplatesDir(t *testing.T) {
	user := &User{Name: "testuser", HomeDir: t.TempDir()}

	if err := CopyConfigFiles(fstest.MapFS{}, user); err == nil {
		t.Error("CopyConfigFiles() should error when files/ is missing from the templates")
	}
}

// Integration test for flakes

func TestPrepareEnvironment_WithFlakes(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create .camp directory and config with flakes
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	configPath := filepath.Join(campDir, "camp.yml")
	configContent := `env:
  EDITOR: nvim

flakes:
  - name: test-flake
    url: "github:test/flake"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: packages
        type: home
      - name: darwinModules.test
        type: system
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create user
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
		Flakes:       []Flake{},
	}

	// Run PrepareEnvironment
	if err := PrepareEnvironment(templates.FS, user); err != nil {
		t.Fatalf("PrepareEnvironment() failed: %v", err)
	}

	// Verify flake.nix was created
	flakePath := filepath.Join(tmpHome, ".camp", "nix", "flake.nix")
	if _, err := os.Stat(flakePath); os.IsNotExist(err) {
		t.Fatal("PrepareEnvironment() should create flake.nix")
	}

	// Read generated flake.nix
	content, err := os.ReadFile(flakePath)
	if err != nil {
		t.Fatalf("Failed to read generated flake.nix: %v", err)
	}

	contentStr := string(content)

	// Verify flake is in inputs section
	if !strings.Contains(contentStr, "test-flake = {") {
		t.Error("Generated flake.nix should contain test-flake in inputs")
	}

	if !strings.Contains(contentStr, `url = "github:test/flake";`) {
		t.Error("Generated flake.nix should contain flake URL")
	}

	if !strings.Contains(contentStr, `inputs.nixpkgs.follows = "nixpkgs";`) {
		t.Error("Generated flake.nix should contain follows declaration")
	}

	// Verify flake is in outputs function signature
	if !strings.Contains(contentStr, "test-flake,") {
		t.Error("Generated flake.nix should include test-flake in outputs signature")
	}

	// Verify system-level output is injected into darwin modules
	if !strings.Contains(contentStr, "test-flake.darwinModules.test") {
		t.Error("Generated flake.nix should inject system output into darwin modules")
	}

	// Verify home-level output is injected into home-manager modules
	if !strings.Contains(contentStr, "test-flake.packages") {
		t.Error("Generated flake.nix should inject home output into home-manager modules")
	}

	// Verify user data is present
	if !strings.Contains(contentStr, "testuser") {
		t.Error("Generated flake.nix should contain username")
	}

	if !strings.Contains(contentStr, "testhost") {
		t.Error("Generated flake.nix should contain hostname")
	}
}

func TestLocalFlakeNames(t *testing.T) {
	flakes := []Flake{
		{Name: "dotfiles", URL: "path:/Users/me/dotfiles"},
		{Name: "remote", URL: "github:user/repo"},
		{Name: "work", URL: "/Users/me/work"},
	}

	got := LocalFlakeNames(flakes)
	want := []string{"dotfiles", "work"}

	if !slices.Equal(got, want) {
		t.Errorf("LocalFlakeNames() = %v, want %v", got, want)
	}
}

func TestLocalFlakeNames_NoLocalFlakes(t *testing.T) {
	flakes := []Flake{
		{Name: "remote", URL: "github:user/repo"},
	}

	if got := LocalFlakeNames(flakes); len(got) != 0 {
		t.Errorf("LocalFlakeNames() = %v, want empty", got)
	}
}

func TestLocalFlakeUpdateArgs_OnlyLocalFlakes(t *testing.T) {
	flakes := []Flake{
		{Name: "dotfiles", URL: "path:/Users/me/dotfiles"},
		{Name: "remote", URL: "github:user/repo"},
		{Name: "work", URL: "git+file:///Users/me/work"},
	}

	got := localFlakeUpdateArgs("/home/me/.camp/nix", flakes)
	want := []string{
		"--extra-experimental-features",
		"nix-command flakes",
		"flake",
		"update",
		"dotfiles",
		"work",
		"--flake",
		"/home/me/.camp/nix",
	}

	if !slices.Equal(got, want) {
		t.Errorf("localFlakeUpdateArgs() = %v, want %v", got, want)
	}
}

func TestLocalFlakeUpdateArgs_NoLocalFlakes(t *testing.T) {
	flakes := []Flake{
		{Name: "remote", URL: "github:user/repo"},
	}

	if got := localFlakeUpdateArgs("/home/me/.camp/nix", flakes); got != nil {
		t.Errorf("localFlakeUpdateArgs() = %v, want nil", got)
	}
}

func TestLocalFlakeUpdateArgs_NoFlakes(t *testing.T) {
	if got := localFlakeUpdateArgs("/home/me/.camp/nix", nil); got != nil {
		t.Errorf("localFlakeUpdateArgs() = %v, want nil", got)
	}
}

func TestUpdateLocalFlakes_RunsNixFlakeUpdate(t *testing.T) {
	exec := utilstest.StubExec(t)
	user := &User{
		Name:    "testuser",
		HomeDir: t.TempDir(),
		Flakes: []Flake{
			{Name: "dotfiles", URL: "path:/Users/me/dotfiles"},
			{Name: "remote", URL: "github:user/repo"},
		},
	}

	if err := UpdateLocalFlakes(user); err != nil {
		t.Fatalf("UpdateLocalFlakes() failed: %v", err)
	}

	nixDir := filepath.Join(user.HomeDir, ".camp", "nix")
	want := [][]string{append([]string{"nix"}, localFlakeUpdateArgs(nixDir, user.Flakes)...)}
	if !slices.EqualFunc(exec.Calls, want, slices.Equal) {
		t.Errorf("UpdateLocalFlakes() ran %v, want %v", exec.Calls, want)
	}
}

func TestUpdateLocalFlakes_NoLocalFlakesIsNoop(t *testing.T) {
	exec := utilstest.StubExec(t)
	user := &User{
		Name:    "testuser",
		HomeDir: t.TempDir(),
		Flakes: []Flake{
			{Name: "remote", URL: "github:user/repo"},
		},
	}

	if err := UpdateLocalFlakes(user); err != nil {
		t.Errorf("UpdateLocalFlakes() should be a no-op without local flakes, got: %v", err)
	}
	if len(exec.Calls) != 0 {
		t.Errorf("UpdateLocalFlakes() should run nothing without local flakes, ran %v", exec.Calls)
	}
}
