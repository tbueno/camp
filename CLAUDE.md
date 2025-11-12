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

## Current Commands

The source code commands can be executes with the command `go run main.go <command>`. Running the command `go run main.go` will display the help message for all the commands.

## Migration from Optishell

This project is a new version of the original Optishell project. The original project can be found in `./tmp/optishell`. In that project, any reference to optishell should be considered a reference to `camp` when ported to this new version.

### Important considerations
- The current project is supposed to be a more modern solution, so if something needs to adapted, the current project should take precedence.
-
