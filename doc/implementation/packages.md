# Implementation Plan: Adding Packages Array to camp.yml

## Overview
This plan details how to add a `packages` field to camp.yml that allows users to define Nix packages to be installed via home-manager. The packages will be simple strings (package names from nixpkgs) that get injected into the home-manager configuration.

---

## Phase 1: Update Config Structure to Support Packages Array

### Files to Modify:
- `internal/system/config.go`
- `internal/system/types.go`

### Changes:

#### 1.1 Update `CampConfig` struct in `config.go`
Add a `Packages` field to the `CampConfig` struct:

```go
type CampConfig struct {
    Env      map[string]string `yaml:"env"`      // Environment variables
    Packages []string          `yaml:"packages"` // Nix packages to install
    Flakes   []Flake           `yaml:"flakes"`   // External Nix flakes to integrate
}
```

#### 1.2 Update `DefaultConfig()` function in `config.go`
Initialize the Packages slice:

```go
func DefaultConfig() *CampConfig {
    return &CampConfig{
        Env:      make(map[string]string),
        Packages: []string{},
        Flakes:   []Flake{},
    }
}
```

#### 1.3 Update `LoadConfig()` function in `config.go`
Add initialization check for Packages (around line 54):

```go
// Initialize Packages slice if nil
if config.Packages == nil {
    config.Packages = []string{}
}
```

#### 1.4 Update `User` struct in `types.go`
Add a `Packages` field to the User struct:

```go
type User struct {
    Name         string
    HomeDir      string
    Platform     string
    Architecture string
    Shell        string
    HostName     string
    EnvVars      map[string]string // Custom environment variables from camp.yml
    Packages     []string          // Nix packages to install from camp.yml
    Flakes       []Flake           // External Nix flakes from camp.yml
}
```

#### 1.5 Update `NewUser()` function in `types.go`
Initialize Packages:

```go
user := &User{
    // ... existing fields ...
    EnvVars:  make(map[string]string),
    Packages: []string{},
    Flakes:   []Flake{},
}
```

#### 1.6 Update `User.Reload()` function in `types.go`
Load packages from config (around line 70):

```go
// Update Packages from config
if config.Packages != nil {
    u.Packages = config.Packages
} else {
    u.Packages = []string{}
}
```

---

## Phase 2: Add Validation for Package Names

### Files to Modify:
- `internal/system/config.go`

### Changes:

#### 2.1 Add `ValidatePackages()` method to `CampConfig`
Add after the `ValidateFlakes()` method (around line 162):

```go
// ValidatePackages validates the packages configuration
func (c *CampConfig) ValidatePackages() error {
    if c.Packages == nil || len(c.Packages) == 0 {
        // Empty packages is valid
        return nil
    }

    // Track package names to ensure uniqueness
    seen := make(map[string]bool)

    for i, pkg := range c.Packages {
        // Validate package name is not empty
        if strings.TrimSpace(pkg) == "" {
            return fmt.Errorf("package at index %d is empty or contains only whitespace", i)
        }

        // Validate package name doesn't contain invalid characters
        // Nix package names should be alphanumeric with hyphens, underscores, and dots
        if !isValidNixPackageName(pkg) {
            return fmt.Errorf("package '%s' has invalid format - must contain only letters, numbers, hyphens, underscores, and dots", pkg)
        }

        // Check for duplicates
        if seen[pkg] {
            return fmt.Errorf("duplicate package '%s' - package names must be unique", pkg)
        }
        seen[pkg] = true
    }

    return nil
}

// isValidNixPackageName checks if a string is a valid Nix package name
// Valid package names contain letters, numbers, hyphens, underscores, and dots
// They may also contain attribute paths like "python3Packages.requests"
func isValidNixPackageName(s string) bool {
    if len(s) == 0 {
        return false
    }

    for _, r := range s {
        if !((r >= 'a' && r <= 'z') ||
            (r >= 'A' && r <= 'Z') ||
            (r >= '0' && r <= '9') ||
            r == '-' ||
            r == '_' ||
            r == '.') {
            return false
        }
    }

    return true
}
```

#### 2.2 Update `Validate()` method in `config.go`
Add package validation call (around line 98):

```go
func (c *CampConfig) Validate() error {
    // Validate flakes configuration
    if err := c.ValidateFlakes(); err != nil {
        return err
    }

    // Validate packages configuration
    if err := c.ValidatePackages(); err != nil {
        return err
    }

    return nil
}
```

---

## Phase 3: Update Template Rendering to Pass Packages

### Files to Modify:
- `internal/system/template.go`

### Changes:

#### 3.1 Update `TemplateData` struct in `template.go`
Add a `Packages` field (around line 13):

```go
type TemplateData struct {
    Name         string            // Username
    HostName     string            // Machine hostname
    Platform     string            // OS (darwin/linux)
    Architecture string            // CPU arch (amd64/arm64)
    HomeDir      string            // User's home directory
    EnvVars      map[string]string // Custom environment variables
    Packages     []string          // Nix packages to install
    Flakes       []Flake           // External Nix flakes to integrate
}
```

#### 3.2 Update `NewTemplateData()` function in `template.go`
Include packages in the template data (around line 24):

```go
func NewTemplateData(user *User) *TemplateData {
    return &TemplateData{
        Name:         user.Name,
        HostName:     user.HostName,
        Platform:     user.Platform,
        Architecture: user.Architecture,
        HomeDir:      user.HomeDir,
        EnvVars:      user.EnvVars,
        Packages:     user.Packages,
        Flakes:       user.Flakes,
    }
}
```

---

## Phase 4: Update common.nix Template to Use Packages

### Files to Modify:
- `templates/files/modules/common.nix`

### Changes:

#### 4.1 Update common.nix to inject packages dynamically
Replace the hardcoded packages list (around line 15) with template rendering:

```nix
{ config, lib, pkgs, user, usersPath, customEnvVars, customPackages, ... }:

{
  programs.home-manager.enable = true;

  programs.direnv = {
    enable = true;
    enableZshIntegration = true;   # Enable zsh integration for direnv
    enableBashIntegration = false;
    nix-direnv.enable = true;
  };

  home = {
    homeDirectory = "${usersPath}";
    packages = with pkgs; [
      devbox
      direnv
      git  # Add git from Nix to ensure it's available
      {{- range .Packages }}
      {{ . }}
      {{- end }}
    ];
    stateVersion = "24.05";
    username = user;

    # Session variables managed by home-manager through zsh
    sessionVariables = customEnvVars;
  };

  # Enable zsh management with dotDir approach
  programs.zsh = {
    enable = true;
    dotDir = ".camp";  # Relative path from home directory

    # Source user's original .zshrc after camp's config loads
    initExtra = ''
      [ -f ~/.zshrc ] && source ~/.zshrc
    '';
  };
}
```

**Note:** The `{{- range .Packages }}` syntax with the `-` trims whitespace for cleaner output.

---

## Phase 5: Update Flake Template to Pass Packages to common.nix

### Files to Modify:
- `templates/files/flake.nix`

### Changes:

#### 5.1 Update specialArgs in flake.nix
Add `customPackages` to the specialArgs (around line 48):

```nix
# Define variables that will be injected in other templates
specialArgs = {
  inherit hostName user usersPath;
  customEnvVars = {
    {{- range $key, $value := .EnvVars }}
    "{{ $key }}" = "{{ $value }}";
    {{- end }}
  };
  customPackages = [
    {{- range .Packages }}
    "{{ . }}"
    {{- end }}
  ];
};
```

---

## Phase 6: Testing Strategy

### Files to Create/Modify:
- `internal/system/config_test.go` (add new tests)
- `internal/system/template_test.go` (add new tests)

### Test Cases:

#### 6.1 Config Tests (add to `config_test.go`)

**Tests to add:**
1. `TestLoadConfig_WithPackages` - Load config with packages array
2. `TestLoadConfig_EmptyPackages` - Load config with empty packages array
3. `TestLoadConfig_NoPackagesSection` - Load config without packages section
4. `TestValidatePackages_DuplicatePackages` - Reject duplicate package names
5. `TestValidatePackages_EmptyPackageName` - Reject empty package names
6. `TestValidatePackages_InvalidCharacters` - Reject invalid characters in package names
7. `TestValidatePackages_ValidPackageNames` - Accept valid package names including attribute paths
8. `TestUserReload_WithPackages` - User reload loads packages from config
9. `TestSaveConfig_WithPackages` - Save and load config preserves packages

#### 6.2 Template Tests (add to `template_test.go`)

**Tests to add:**
1. `TestNewTemplateData_WithPackages` - Template data includes packages
2. `TestCompileTemplate_WithPackages` - Template compiles with packages
3. `TestCompileTemplate_WithEmptyPackages` - Template handles empty packages array

#### 6.3 Integration Tests

**Tests to add:**
1. `TestPackagesIntegration` - Full end-to-end test from config to rendered templates

---

## Phase 7: Documentation Updates

### Files to Modify:
- `CLAUDE.md`
- `README.md`

### Changes:

#### 7.1 Update CLAUDE.md

Add a new section under "Configuration System":

```markdown
### Package Management

Camp allows you to declaratively manage Nix packages through the configuration file.

**Configuration Structure:**

```yaml
# Nix packages to install via home-manager
packages:
  - git
  - neovim
  - ripgrep
  - python3
  - nodejs_20
  - python3Packages.requests
```

**Features:**
- Package names must be valid Nix package identifiers
- Supports attribute paths (e.g., `python3Packages.requests`)
- Packages are deduplicated automatically
- Invalid package names are caught during validation

**Package Name Validation:**
- Must contain only letters, numbers, hyphens, underscores, and dots
- Cannot be empty or whitespace-only
- No duplicate packages allowed
- Supports nested attributes with dot notation

**Integration:**
- Packages are automatically added to home-manager's `home.packages`
- Applied on `camp env rebuild`
- Works on both macOS and Linux

**Example:**

```yaml
env:
  EDITOR: nvim

packages:
  # Core utilities
  - git
  - curl
  - wget

  # Development tools
  - neovim
  - ripgrep
  - fd

  # Programming languages
  - python3
  - nodejs_20
  - go

  # Python packages
  - python3Packages.requests
  - python3Packages.flask
```
```

#### 7.2 Update README.md

Update the configuration example section:

```markdown
### Configuration

Camp uses a YAML configuration file to customize your environment.

**Location:** `~/.camp/camp.yml` (or `.yaml`)

**Example configuration:**
```yaml
env:
  EDITOR: nvim
  BROWSER: firefox
  CUSTOM_VAR: custom_value

packages:
  - git
  - neovim
  - ripgrep
  - python3
  - nodejs_20
```

**Configuration Sections:**
- `env`: Custom environment variables injected into your Nix environment
- `packages`: Nix packages to install via home-manager
- `flakes`: External Nix flakes to integrate (see Flakes documentation)

All packages will be installed to your home environment and available in your PATH.
```

---

## Summary

This implementation plan provides:

1. **Phase 1:** Core data structure changes to support packages in config and user objects
2. **Phase 2:** Validation to ensure package names are valid and unique
3. **Phase 3:** Template data preparation to pass packages to Nix templates
4. **Phase 4:** Nix template updates to actually install the packages
5. **Phase 5:** Flake integration to pass packages through specialArgs
6. **Phase 6:** Comprehensive testing strategy covering all aspects
7. **Phase 7:** Documentation updates for users

**Key Design Decisions:**
- Packages are simple strings (no complex structures needed)
- Validation happens at config load time
- Packages are installed at home-manager level (not system-level)
- Support for attribute paths like `python3Packages.requests`
- Automatic deduplication and validation

**Testing Coverage:**
- Config loading and saving
- Validation (duplicates, empty names, invalid characters)
- User reload functionality
- Template rendering
- End-to-end integration

**User Experience:**
Users can simply add package names to their `camp.yml` and run `camp env rebuild` to have them installed automatically.
