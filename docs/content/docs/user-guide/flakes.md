---
title: "Flakes"
linkTitle: "Flakes"
weight: 4
description: >
  Integrating external Nix flakes with Camp
---

Camp supports integrating external Nix flakes to extend your
development environment with custom packages, configurations, and
modules.

## What are Flakes?

Flakes are a standardized way to package and share Nix configurations. They provide:

- **Reproducible dependencies** with lock files
- **Standardized structure** for packages and modules
- **Version pinning** for consistent environments

## Basic Configuration

Define flakes in `~/.camp/camp.yml`:

```yaml
flakes:
  - name: my-tools              # Unique identifier
    url: "github:user/my-tools" # Flake location
    outputs:                    # What to import
      - name: packages
        type: home              # "home" or "system"
```

## Flake Structure

```yaml
flakes:
  - name: string            # Unique identifier (required)
    url: string             # Flake URL (required)
    follows:                # Input overrides (optional)
      nixpkgs: "nixpkgs"
    args:                   # Custom arguments (optional)
      key: value
    outputs:                # Outputs to import (required)
      - name: string        # Output name
        type: string        # "home" or "system"
```

## Supported URL Formats

### GitHub (Public)

```yaml
# Latest commit on default branch
url: "github:username/repository"

# Specific branch
url: "github:username/repository/branch-name"

# Specific tag
url: "github:username/repository/v1.2.3"
```

### GitHub (Private via SSH)

```yaml
url: "git+ssh://git@github.com/company/private-repo.git"
```

### GitLab

```yaml
url: "gitlab:username/repository"
```

### Local Path

```yaml
url: "path:/absolute/path/to/flake"
```

Useful for testing flakes during development.

### Generic Git

```yaml
url: "git+https://git.example.com/repo.git"
url: "git+ssh://git@git.example.com/repo.git"
```

## Output Types

### `type: home` (User-level)

Applies to your user environment via home-manager:

- Available on both macOS and Linux
- User-level packages and configurations
- No sudo required
- Dotfiles, shell configs, development tools

**Common output names:**

- `packages`
- `homeManagerModules.default`
- `homeManagerModules.{name}`

### `type: system` (System-level)

Applies to system configuration via nix-darwin:

- **macOS only**
- System-wide settings and services
- May require sudo/admin privileges
- System packages, LaunchDaemons, system preferences

**Common output names:**

- `darwinModules.default`
- `darwinModules.{name}`

## Input Following

Use the same nixpkgs version as Camp for consistency:

```yaml
flakes:
  - name: my-flake
    url: "github:user/repo"
    follows:
      nixpkgs: "nixpkgs"    # Use Camp's nixpkgs
    outputs:
      - name: packages
        type: home
```

This prevents multiple nixpkgs versions being downloaded.

## Flake Arguments

Pass custom arguments to parameterize flakes:

```yaml
flakes:
  - name: personal-config
    url: "github:user/nix-config"
    args:
      email: "user@example.com"
      enableDevTools: true
      fontSize: 14
      packages: [vim, git, tmux]
    outputs:
      - name: homeManagerModules.default
        type: home
```

### Automatic Arguments

Camp always passes:

- `userName` - Your system username
- `hostName` - Your machine hostname
- `home` - Your home directory path

### Supported Types

| YAML Type | Example     | Nix Output    |
| --------- | ----------- | ------------- |
| String    | `"hello"`   | `"hello"`     |
| Boolean   | `true`      | `true`        |
| Integer   | `42`        | `42`          |
| Float     | `3.14`      | `3.14`        |
| List      | `[a, b]`    | `[ "a" "b" ]` |

### External Flake Pattern

Your external flake must define outputs as functions:

```nix
{
  description = "Personal configuration";

  outputs = { nixpkgs, ... }:
  {
    homeManagerModules.default = {
      userName,      # automatic
      hostName,      # automatic
      home,          # automatic
      email,         # from args
      enableDevTools,# from args
      ...
    }@args: {
      # Use arguments in configuration
      programs.git.userEmail = email;

      home.packages = with pkgs;
        [ git ] ++ (if enableDevTools then [ gcc ] else []);
    };
  };
}
```

## Examples

### Personal Development Tools

```yaml
flakes:
  - name: dev-tools
    url: "github:myuser/dev-environment"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: packages
        type: home
```

### Team Configuration

```yaml
flakes:
  - name: team-config
    url: "git+ssh://git@github.com/company/team-nix.git"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: homeManagerModules.team
        type: home
```

### macOS System Configuration

```yaml
flakes:
  - name: macos-config
    url: "github:user/darwin-config"
    outputs:
      # System-level settings
      - name: darwinModules.system
        type: system
      # User-level settings
      - name: homeManagerModules.user
        type: home
```

### Parameterized Flake

```yaml
flakes:
  - name: personal-setup
    url: "github:user/nix-config"
    args:
      email: "user@example.com"
      gitSigningKey: "ABCD1234"
      enableGPG: true
      terminalFont: "JetBrains Mono"
      fontSize: 14
    outputs:
      - name: darwinModules.default
        type: system
      - name: homeManagerModules.default
        type: home
```

## Workflow

### Adding a Flake

1. Edit `~/.camp/camp.yml`
2. Add flake to `flakes:` section
3. Run `camp env rebuild`

### Updating Flakes

Update all flake dependencies:

```bash
camp env update
camp env rebuild
```

### Removing a Flake

1. Remove from `~/.camp/camp.yml`
2. Run `camp env rebuild`

## Validation

Camp validates flakes on load:

- **Unique names**: No duplicate flake names
- **Valid identifiers**: Alphanumeric, hyphens, underscores only
- **Non-empty URLs**: Every flake needs a URL
- **Valid output types**: Must be "system" or "home"
- **At least one output**: Each flake needs outputs defined

## Template Integration

Flakes are injected into the generated `~/.camp/nix/flake.nix`:

```nix
inputs = {
  nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  # User-defined flakes
  my-tools = {
    url = "github:user/my-tools";
    inputs.nixpkgs.follows = "nixpkgs";
  };
};

outputs = { nixpkgs, my-tools, ... }:
{
  # System outputs (macOS)
  darwinConfigurations."hostname" = {
    modules = [
      ./mac.nix
      my-tools.darwinModules.system  # If type: system
    ];
  };

  # Home outputs
  homeConfigurations."username" = {
    imports = [
      ./modules/common.nix
      my-tools.homeManagerModules.default  # If type: home
    ];
  };
};
```

## Troubleshooting

### Duplicate Flake Name

```yaml
Error: duplicate flake name 'my-flake'
```

Each flake must have a unique name. Rename one:

```yaml
flakes:
  - name: flake-one
    url: "github:user/repo1"
  - name: flake-two
    url: "github:user/repo2"
```

### Invalid Flake Name

```yaml
Error: invalid flake name 'my.flake'
```

Use only letters, numbers, hyphens, underscores:

```yaml
- name: "my-flake"    # ✅ Valid
- name: "my_flake"    # ✅ Valid
- name: "my.flake"    # ❌ Invalid
- name: "my flake"    # ❌ Invalid
```

### Empty URL

```yaml
Error: flake 'my-flake' has empty URL
```

Every flake needs a URL:

```yaml
- name: my-flake
  url: "github:user/repo"  # Required
  outputs:
    - name: packages
      type: home
```

### Private Repository Access

For private repos via SSH:

1. Set up SSH key authentication
2. Test: `ssh -T git@github.com`
3. Use `git+ssh://` URL format

### Output Not Found

```yaml
error: attribute 'packages' missing
```

Verify the flake exports the output you're referencing:

```bash
nix flake show github:user/repo
```

## Best Practices

1. **Use `follows`**: Keep nixpkgs versions consistent
2. **Pin versions**: Use tags or commits for stability
3. **Test locally**: Use `path:` URLs during development
4. **Document outputs**: Comment what each output does
5. **Version in URLs**: Use tags for production flakes

## Example Flake Templates

Check `templates/flakes/` in the Camp repository for:

- Personal package examples
- Team configuration templates
- Parameterized flake examples

## Next Steps

- Create your own flake (see [Developer Guide](/docs/developer-guide/))
- Browse [Flake examples](https://github.com/tbueno/camp/tree/main/templates/flakes)
- Learn about [Templates](/docs/user-guide/templates/)
