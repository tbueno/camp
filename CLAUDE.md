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

## Package Management

Camp allows you to declaratively manage Nix packages through the configuration file.

### Configuration Structure

```yaml
# Nix packages to install via home-manager
packages:
  - git
  - neovim
  - ripgrep
  `

### Features

- **Simple String Array**: Package names are plain strings, no complex structures
- **Attribute Path Support**: Supports nested packages like `python3Packages.requests`
- **Automatic Validation**: Invalid package names are caught at config load time
- **Deduplication**: Duplicate packages are rejected with clear error messages
- **Cross-Platform**: Works on both macOS (via nix-darwin) and Linux (via home-manager)

### Package Name Validation

Package names must meet these requirements:
- Only letters, numbers, hyphens, underscores, and dots allowed
- Cannot be empty or whitespace-only
- No duplicate packages in the list
- Supports attribute paths with dot notation (e.g., `haskellPackages.pandoc`)

**Valid examples:**
- `git`
- `python3`
- `nodejs_20`
- `package-with-hyphens`
- `package_with_underscores`
- `python3Packages.requests`
- `haskellPackages.pandoc`

**Invalid examples:**
- `my package` (contains spaces)
- `my@package` (contains special characters)
- `my/package` (contains slashes)

### Integration

1. **Configuration Loading**: Packages are loaded from `camp.yml` during `User.Reload()`
2. **Validation**: `ValidatePackages()` runs automatically when loading config
3. **Template Rendering**: Packages are passed to Nix templates via `TemplateData`
4. **Nix Installation**: Packages are added to `home.packages` in `common.nix`
5. **Application**: Changes take effect on `camp env rebuild`

### How It Works Internally

**Configuration Flow:**
```go
// In config.go
type CampConfig struct {
    Env      map[string]string
    Packages []string  // Simple string array
    Flakes   []Flake
}

// In types.go
type User struct {
    // ... other fields ...
    Packages []string  // Loaded from config
}
```

**Template Integration:**
```go
// Packages are passed to templates
type TemplateData struct {
    // ... other fields ...
    Packages []string
}

// In flake.nix template
specialArgs = {
    customPackages = [
        "git"
        "neovim"
        "ripgrep"
    ];
};

// In common.nix (pure Nix, not a template)
home.packages = with pkgs; [
    devbox
    direnv
    git
] ++ (map (name: pkgs.${name}) customPackages);
```

The packages are:
1. Defined as strings in `camp.yml`
2. Validated on config load
3. Passed through `specialArgs` in the flake
4. Dynamically resolved to package objects in `common.nix` using Nix's `map` function

### Example Configuration

```yaml
env:
  EDITOR: nvim
  BROWSER: firefox

packages:
  # Core utilities
  - git
  - curl
  - wget

  # Development tools
  - neovim
  - ripgrep
  - fd
  - bat

  # Programming languages
  - python3
  - nodejs_20
  - go

  # Language-specific packages
  - python3Packages.requests
  - python3Packages.flask
  - nodePackages.typescript
```

### Testing Packages

When writing package tests:
- Test config loading with packages: `TestLoadConfig_WithPackages`
- Test validation: `TestValidatePackages_*`
- Test template rendering: `TestNewTemplateData_WithPackages`
- Test saving/loading: `TestSaveConfig_WithPackages`
- Use `t.TempDir()` for test configs
- Verify both valid and invalid package names

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

## Flake Arguments

Camp supports passing arguments to external flakes, allowing you to parameterize flake configurations without hardcoding values.

### Overview

When importing external flakes, camp automatically passes three standard arguments and allows you to define custom arguments:

**Automatic Arguments** (always passed):
- `userName` - from `User.Name`
- `hostName` - from `User.HostName`
- `home` - from `User.HomeDir`

**Custom Arguments**: User-defined values with type inference from YAML

### Configuration Schema

Add arguments to flakes in `~/.camp/camp.yml`:

```yaml
flakes:
  - name: personal-config
    url: "github:username/nix-config"
    args:
      email: "user@example.com"      # string (inferred from YAML)
      enableDevTools: true            # bool
      fontSize: 14                    # number (int)
      threshold: 3.14                 # number (float)
      packages: [vim, git, tmux]      # list of strings
    outputs:
      - name: darwinModules.default
        type: system
      - name: homeManagerModules.default
        type: home
```

### Supported Argument Types

Camp infers types from YAML values and renders them as proper Nix syntax:

| YAML Type | Example Value | Nix Output | Go Type |
|-----------|---------------|------------|---------|
| String | `"hello"` | `"hello"` | `string` |
| Boolean | `true` | `true` | `bool` |
| Integer | `42` | `42` | `int` |
| Float | `3.14` | `3.14` | `float64` |
| List | `[a, b, c]` | `[ "a" "b" "c" ]` | `[]interface{}` |

**Type inference examples**:
```yaml
args:
  name: "camp"           # String (quoted)
  enabled: true          # Boolean
  count: 42              # Integer
  ratio: 0.5             # Float
  items: ["x", "y"]      # List of strings
  ports: [8080, 9090]    # List of integers
  flags: [true, false]   # List of booleans
```

### Validation Rules

Arguments are validated to prevent errors:

1. **Valid Nix identifiers**: Names must contain only letters, numbers, hyphens, underscores
2. **Reserved names**: Cannot use `userName`, `hostName`, `home` (automatically provided)
3. **Supported types**: Only string, bool, number, and list types allowed
4. **List elements**: Lists can only contain strings, booleans, or numbers (no nested structures)

**Validation examples**:
```yaml
# ✅ Valid
args:
  my-arg: "value"
  my_arg: 123
  ARG_NAME: true

# ❌ Invalid
args:
  "my.arg": "value"      # Error: dots not allowed in names
  userName: "override"   # Error: reserved name
  nested: {key: "val"}   # Error: maps not supported
```

### External Flake Pattern

Your external flake must define outputs as functions accepting parameters:

**Example flake** (`github:username/nix-config/flake.nix`):
```nix
{
  description = "Personal Nix configuration";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs }:
  {
    # System-level module (nix-darwin) - accepts parameters
    darwinModules.default = { userName, hostName, home, email, ... }@args: {
      networking.hostName = hostName;

      users.users.${userName} = {
        name = userName;
        home = home;
      };

      # Use custom arguments
      environment.systemPackages = [
        # packages based on email domain, etc.
      ];
    };

    # User-level module (home-manager) - accepts parameters
    homeManagerModules.default = { userName, home, email, enableDevTools, ... }@args: {
      home.username = userName;
      home.homeDirectory = home;

      programs.git = {
        enable = true;
        userEmail = email;
      };

      # Conditional configuration based on arguments
      programs.neovim.enable = enableDevTools;
    };
  };
}
```

**Key points**:
- Modules are functions, not attribute sets
- Use `{ userName, hostName, home, customArg, ... }@args:` pattern
- Include `...` to accept additional arguments
- Access automatic args (userName, hostName, home) and custom args

### How It Works

**Argument passing flow**:
```
1. User defines args in camp.yml
         ↓
2. LoadConfig() unmarshals YAML with native types
         ↓
3. ValidateFlakeArgs() checks names and types
         ↓
4. CompileTemplates() renders flake.nix:
   - Merges automatic args + custom args
   - Calls renderNixValue() for type-safe rendering
         ↓
5. Generated flake.nix calls external flake outputs as functions:
   (my-flake.darwinModules.default {
     userName = "user";
     hostName = "host";
     home = "/Users/user";
     email = "user@example.com";
     enableDevTools = true;
   })
         ↓
6. ExecuteRebuild() applies configuration
```

### Generated Template Output

Camp generates function calls in `~/.camp/nix/flake.nix`:

**For system-level outputs**:
```nix
modules = [
  ./mac.nix

  # User's external flake, called with arguments
  (personal-config.darwinModules.default {
    userName = "user";
    hostName = "macbook";
    home = "/Users/user";
    email = "user@example.com";
    enableDevTools = true;
    fontSize = 14;
  })
];
```

**For home-level outputs**:
```nix
home-manager.users.${user} = {
  imports = [
    ./modules/common.nix

    # User's external flake, called with arguments
    (personal-config.homeManagerModules.default {
      userName = "user";
      hostName = "macbook";
      home = "/Users/user";
      email = "user@example.com";
      packages = [ "vim" "git" "tmux" ];
    })
  ];
};
```

### String Escaping

Camp properly escapes special characters in string arguments:

| Character | YAML Input | Nix Output |
|-----------|------------|------------|
| Quote | `"hello \"world\""` | `"hello \"world\""` |
| Backslash | `"path\\to\\file"` | `"path\\to\\file"` |
| Newline | `"line1\nline2"` | `"line1\nline2"` |

### Complete Example

**camp.yml**:
```yaml
env:
  EDITOR: nvim

flakes:
  - name: personal-config
    url: "github:tbueno/nix-config"
    args:
      email: "tbueno@gmail.com"
      enableDevTools: true
      fontSize: 14
      packages: [vim, git, tmux, ripgrep]
    outputs:
      - name: darwinModules.default
        type: system
      - name: homeManagerModules.default
        type: home
```

**External flake** (`github:tbueno/nix-config`):
```nix
{
  outputs = { nixpkgs, ... }:
  {
    darwinModules.default = { userName, hostName, email, enableDevTools, fontSize, packages, ... }: {
      networking.hostName = hostName;

      users.users.${userName}.description = email;

      environment.systemPackages = with pkgs;
        packages ++ (if enableDevTools then [ gcc cmake ] else []);
    };

    homeManagerModules.default = { email, fontSize, packages, ... }: {
      programs.git.userEmail = email;

      programs.alacritty.settings.font.size = fontSize;

      home.packages = with pkgs; packages;
    };
  };
}
```

**Result**: Running `camp env rebuild` generates a flake that calls these modules with your specified arguments, creating a fully parameterized environment.

### Testing Flake Arguments

When writing tests for flake arguments:
- Test YAML type inference (string, bool, int, float, list)
- Verify validation catches invalid names and types
- Test Nix value rendering with `renderNixValue()`
- Verify template integration with flake args
- Test automatic args (userName, hostName, home) are passed
- Verify custom args render with correct Nix syntax

**Example test**:
```go
func TestFlakeWithArgs(t *testing.T) {
    config := &CampConfig{
        Flakes: []Flake{
            {
                Name: "test",
                URL:  "github:user/flake",
                Args: map[string]interface{}{
                    "email": "test@example.com",
                    "enabled": true,
                },
                Outputs: []FlakeOutput{
                    {Name: "packages", Type: OutputTypeHome},
                },
            },
        },
    }

    // Validate
    if err := config.ValidateFlakes(); err != nil {
        t.Fatalf("Validation failed: %v", err)
    }

    // Render template and verify args are passed correctly
    // ...
}
```

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
