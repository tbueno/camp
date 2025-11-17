# Camp - Golang CLI Application

## Project Overview
Camp is a command line application built with Go and the Cobra CLI framework. It serves as your helpful command line companion.

## Tech Stack
- **Language**: Go 1.24.4
- **CLI Framework**: Cobra (github.com/spf13/cobra v1.10.0)
- **Configuration**: YAML (gopkg.in/yaml.v3)
- **Testing**: Go's built-in testing package

## Project Structure
```
camp/
├── go.mod                 # Go module definition
├── main.go               # Application entry point
├── cmd/
│   ├── root.go           # Root command implementation
│   ├── root_test.go      # Unit tests for root command
│   ├── env.go            # Env command
│   └── bootstrap.go      # Bootstrap command
├── internal/
│   ├── system/
│   │   ├── types.go      # System, User, and other core types
│   │   ├── config.go     # Configuration loading (camp.yml)
│   │   ├── config_test.go # Config tests
│   │   ├── template.go   # Template rendering
│   │   ├── template_test.go # Template tests
│   │   ├── system.go     # System info functions
│   │   ├── bootstrap.go  # Bootstrap logic
│   │   └── envrc.go      # Envrc parsing
│   └── utils/
│       ├── command.go    # Command execution utilities
│       ├── io.go         # File I/O utilities
│       └── template.go   # Template file utilities
├── templates/
│   ├── initial/          # Bootstrap templates
│   └── files/            # Runtime templates for env rebuild
│       ├── flake.nix     # Main Nix flake template
│       ├── mac.nix       # macOS configuration
│       ├── linux.nix     # Linux configuration
│       └── modules/
│           └── common.nix # Shared home-manager config
└── doc/
    └── implementation/   # Implementation plans and documentation
```

## Development Commands
- **Build**: `go build`
- **Run**: `go run main.go`
- **Test**: `go test ./...`

## Code Style & Conventions
- Follow standard Go conventions. Run `go fmt` after executing commands that alter code.
- Use Cobra's best practices for CLI commands
- All commands should use `cmd.OutOrStdout()` for output to support testing
- Write comprehensive unit tests for all commands
- Keep command logic in the `cmd/` package
- Main entry point should be minimal, just calling `cmd.Execute()`
- When new features are added or changed, make sure to update the documentation in `doc/` folder.

## Testing Guidelines
- Test command names, descriptions, and output
- Use buffer capture for testing command output
- Test both success and error scenarios
- Maintain high test coverage
- Always run the tests when code is changed!

## Configuration System

Camp uses a YAML-based configuration file to customize environment settings.

### Configuration File Location
- **Primary location**: `~/.camp/camp.yml`
- **Alternative**: `~/.camp/camp.yaml`

### Configuration Structure

```yaml
# Environment variables to inject into the Nix environment
env:
  EDITOR: "nvim"
  BROWSER: "firefox"
  CUSTOM_VAR: "value"
```

### How It Works

1. **Automatic Loading**: The User type automatically loads configuration on initialization via `NewUser()`
2. **Reload Method**: Use `user.Reload()` to refresh configuration from disk
3. **Environment Variables**: Custom env vars are stored in `User.EnvVars` map
4. **Default Behavior**: If no config file exists, camp uses sensible defaults

### Configuration API

- `LoadConfig(path string)` - Load config from specific path
- `LoadUserConfig(homeDir string)` - Load config from user's home directory
- `DefaultConfig()` - Get default configuration
- `(c *CampConfig) SaveConfig(path string)` - Save configuration to file
- `(c *CampConfig) Validate()` - Validate configuration structure
- `(u *User) Reload()` - Reload user configuration from disk

### Testing Configuration

When writing tests that use User configuration:
- Use temporary directories with `t.TempDir()`
- Create test config files in `tmpDir/.camp/camp.yml`
- Test both with and without config files
- Verify default behavior when config is missing

## Template System

Camp uses Go's `text/template` package to generate Nix configuration files dynamically.

### Template Structure

Templates are stored in `templates/files/` and include:
- **flake.nix** - Main Nix flake configuration (uses Go template syntax)
- **mac.nix** - macOS system configuration
- **linux.nix** - Linux configuration
- **modules/common.nix** - Shared home-manager configuration

### Template Rendering

**Template Data Structure:**
```go
type TemplateData struct {
    Name         string            // Username
    HostName     string            // Machine hostname
    Platform     string            // OS (darwin/linux)
    Architecture string            // CPU arch (amd64/arm64)
    HomeDir      string            // User's home directory
    EnvVars      map[string]string // Custom environment variables
}
```

**Template Rendering Process:**
1. Load template file from `templates/files/`
2. Create `TemplateData` from User
3. Parse and execute template with data
4. Write rendered output to `~/.camp/nix/`

### Template API

- `NewTemplateData(user *User)` - Create template data from User
- `CompileTemplate(templatePath, data)` - Parse and render a template file
- `RenderFlakeTemplate(user *User)` - Render flake.nix with user data

### File Utilities

File copying utilities in `internal/utils/template.go`:
- `CopyFile(src, dest)` - Copy a single file
- `CopyDir(src, dest)` - Recursively copy directory
- `CopyTemplateFiles(srcDir, destDir)` - Copy all template files
- `EnsureDir(path)` - Create directory if it doesn't exist

### Template Syntax

Templates use Go template syntax with the following patterns:
- `{{.Name}}` - Insert username
- `{{.HostName}}` - Insert hostname
- `{{range $key, $value := .EnvVars}}` - Iterate over environment variables

Example from flake.nix:
```nix
hostName = "{{.HostName}}";
user = "{{.Name}}";
customEnvVars = {
  {{- range $key, $value := .EnvVars }}
  "{{ $key }}" = "{{ $value }}";
  {{- end }}
};
```

### Testing Templates

When writing template tests:
- Create temporary template files in test directories
- Use `t.TempDir()` for test output directories
- Verify rendered output contains expected values
- Test both with and without environment variables
- Mock file operations where appropriate

## Environment Rebuild System

Camp provides an `env rebuild` command to rebuild the development environment using Nix configurations.

### Rebuild Process Flow

```
PrepareEnvironment()
    ↓
1. Create ~/.camp/nix/ directory
    ↓
2. CopyConfigFiles()
    ├─ Copy mac.nix
    ├─ Copy linux.nix
    └─ Copy modules/ directory
    ↓
3. CompileTemplates()
    ├─ Reload user config (camp.yml)
    ├─ Render flake.nix with user data
    └─ Save to ~/.camp/nix/flake.nix
    ↓
ExecuteRebuild()
    ├─ macOS: nix run nix-darwin#darwin-rebuild switch --flake ~/.camp/nix#<hostname>
    └─ Linux: home-manager switch --impure -b backup --flake ~/.camp/nix#<username>
```

### Rebuild Functions

**PrepareEnvironment(user *User)**
- Orchestrates the entire preparation process
- Creates necessary directories
- Calls CopyConfigFiles() and CompileTemplates()
- Returns error if any step fails

**CopyConfigFiles(user *User)**
- Copies .nix files from `templates/files/` to `~/.camp/nix/`
- Excludes `flake.nix` (rendered separately)
- Recursively copies directories (e.g., `modules/`)
- Preserves file permissions

**CompileTemplates(user *User)**
- Reloads user configuration from camp.yml
- Renders flake.nix template with current user data
- Writes output to `~/.camp/nix/flake.nix`

**ExecuteRebuild(user *User)**
- Runs platform-specific rebuild commands
- **macOS**: Uses nix-darwin for system configuration
- **Linux**: Uses home-manager for user environment
- Returns error for unsupported platforms

### Platform-Specific Behavior

**macOS (darwin):**
- Uses `nix-darwin` for system-level configuration
- Command: `nix run nix-darwin#darwin-rebuild switch --flake <path>#<hostname>`
- Requires nix-darwin to be installed
- Applies both system and home-manager configurations

**Linux:**
- Uses `home-manager` for user-level configuration
- Command: `home-manager switch --impure -b backup --flake <path>#<username>`
- Requires home-manager to be installed
- Creates backups of modified files

### Error Handling

The rebuild system handles errors at each stage:
- Missing directories are created automatically
- Template rendering errors are reported clearly
- Platform detection prevents unsupported operations
- Command execution errors are propagated with context

### Testing Rebuild Logic

When writing rebuild tests:
- Use temporary directories for all file operations
- Create mock config files in test directories
- Verify directory and file creation
- Test platform-specific logic separately
- Mock external command execution (nix, home-manager)
- Test error conditions (missing directories, unsupported platforms)

## Flake System

Camp supports integrating external Nix flakes to extend your development environment with custom packages, configurations, and modules.

### What are Flakes?

Flakes are a standardized way to package and share Nix configurations. They provide:
- **Reproducible dependencies** with lock files
- **Standardized structure** for packages and modules
- **Version pinning** for consistent environments

### Configuration

Define flakes in `~/.camp/camp.yml`:

```yaml
env:
  EDITOR: nvim

flakes:
  - name: my-tools              # Unique identifier
    url: "github:user/my-tools" # Flake location
    follows:                    # Optional: input overrides
      nixpkgs: "nixpkgs"
    outputs:                    # What to import
      - name: packages
        type: home              # "home" or "system"
```

### Flake Types

**FlakeOutputType**: Defines where a flake output is applied
- `FlakeOutputType = "system"` - System-level (nix-darwin, macOS only)
- `FlakeOutputType = "home"` - User-level (home-manager, macOS/Linux)

**Flake Structure**:
```go
type Flake struct {
    Name    string            // Unique identifier
    URL     string            // Flake URL (github:, git+ssh:, path:, etc.)
    Follows map[string]string // Input dependency overrides
    Outputs []FlakeOutput     // Which outputs to import
}

type FlakeOutput struct {
    Name string          // Output name (e.g., "packages", "homeManagerModules.default")
    Type FlakeOutputType // Where to apply: "system" or "home"
}
```

### Supported URL Formats

```yaml
# GitHub (public)
url: "github:username/repository"
url: "github:username/repository/branch-name"

# GitHub (private via SSH)
url: "git+ssh://git@github.com/company/private-repo.git"

# GitLab
url: "gitlab:username/repository"

# Local path (for development)
url: "path:/absolute/path/to/flake"

# Generic Git
url: "git+https://git.example.com/repo.git"
```

### How It Works

1. **User Configuration**: Define flakes in `~/.camp/camp.yml`
2. **Validation**: `ValidateFlakes()` checks for errors (unique names, valid URLs, etc.)
3. **Template Rendering**: Flakes are dynamically injected into generated `flake.nix`
4. **Nix Integration**: System applies flake outputs during rebuild

**Template Generation Flow**:
```
User edits ~/.camp/camp.yml
         ↓
LoadConfig() + ValidateFlakes()
         ↓
User.Reload() loads flakes
         ↓
CompileTemplates() renders flake.nix with:
  - Flake inputs section
  - Flake outputs in function signature
  - System outputs → nix-darwin modules
  - Home outputs → home-manager modules
         ↓
ExecuteRebuild() applies configuration
```

### Commands

**Rebuild with flakes**:
```bash
camp env rebuild
```
Integrates all defined flakes into your environment.

**Update flake dependencies**:
```bash
camp env update
```
Updates `flake.lock` with latest versions of all flake inputs.

### Validation Rules

The system validates flakes to prevent errors:

- **Unique names**: No duplicate flake names allowed
- **Valid Nix identifiers**: Names must be alphanumeric with hyphens/underscores only
- **Non-empty URLs**: Every flake must have a URL
- **Valid output types**: Must be "system" or "home"
- **At least one output**: Each flake must define outputs to import

**Example validation error**:
```
Error: duplicate flake name 'my-flake' - flake names must be unique
```

### Template Integration

Flakes are injected into `templates/files/flake.nix`:

**Inputs section**:
```nix
inputs = {
  nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  # User-defined flakes
  {{- range .Flakes }}
  {{ .Name }} = {
    url = "{{ .URL }}";
    {{- range $key, $value := .Follows }}
    inputs.{{ $key }}.follows = "{{ $value }}";
    {{- end }}
  };
  {{- end }}
};
```

**System outputs** (nix-darwin, macOS):
```nix
modules = [
  ./mac.nix

  # Custom system-level flake modules
  {{- range $flake := .Flakes }}
    {{- range .Outputs }}
      {{- if eq .Type "system" }}
  {{ $flake.Name }}.{{ .Name }}
      {{- end }}
    {{- end }}
  {{- end }}
];
```

**Home outputs** (home-manager, macOS/Linux):
```nix
home-manager.users.${user} = {
  imports = [
    ./modules/common.nix

    # Custom home-level flake modules
    {{- range $flake := .Flakes }}
      {{- range .Outputs }}
        {{- if eq .Type "home" }}
    {{ $flake.Name }}.{{ .Name }}
        {{- end }}
      {{- end }}
    {{- end }}
  ];
};
```

### Example Use Cases

**Personal packages**:
```yaml
flakes:
  - name: my-packages
    url: "github:username/nix-packages"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: packages
        type: home
```

**Team configuration**:
```yaml
flakes:
  - name: company-tools
    url: "git+ssh://git@github.com/company/nix-tools.git"
    outputs:
      - name: darwinModules.company
        type: system
      - name: homeManagerModules.company
        type: home
```

**Local development**:
```yaml
flakes:
  - name: local-test
    url: "path:/Users/me/projects/test-flake"
    outputs:
      - name: packages
        type: home
```

### Testing Flakes

When writing flake tests:
- Use temporary directories for config files
- Create test flakes with valid YAML
- Verify validation catches errors
- Test template rendering with flakes
- Verify system vs home output routing
- Test integration with PrepareEnvironment()

**Example test structure**:
```go
func TestFlakeIntegration(t *testing.T) {
    // Create config with flakes
    config := &CampConfig{
        Flakes: []Flake{
            {
                Name: "test-flake",
                URL:  "github:test/flake",
                Outputs: []FlakeOutput{
                    {Name: "packages", Type: OutputTypeHome},
                },
            },
        },
    }

    // Validate
    if err := config.ValidateFlakes(); err != nil {
        t.Errorf("Validation failed: %v", err)
    }

    // Render template and verify output
    // ...
}
```

### Example Templates

See `templates/flakes/` for ready-to-use examples:
- `personal-packages.yml` - Personal development tools
- `team-tools.yml` - Organization-wide configurations
- `README.md` - Comprehensive guide with all URL formats

## Current Commands

To see a list of available commands, run:
```bash
camp --help
# or during development:
go run main.go --help
```

For detailed help on any command:
```bash
camp <command> --help
# Example: camp env rebuild --help
```

## Migration from Optishell

This project is a new version of the original Optishell project. The original project can be found in `./tmp/optishell`. In that project, any reference to optishell should be considered a reference to `camp` when ported to this new version.

### Important considerations
- The current project is supposed to be a more modern solution, so if something needs to adapted, the current project should take precedence.
-
