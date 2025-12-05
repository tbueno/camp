---
title: "Configuration"
linkTitle: "Configuration"
weight: 3
description: >
  Understanding the Camp configuration file
---

Camp uses a YAML configuration file to manage your development environment settings.

## Configuration File Location

Camp looks for configuration in:

- **Primary**: `~/.camp/camp.yml`
- **Alternative**: `~/.camp/camp.yaml`

If no configuration file exists, Camp uses sensible defaults.

## Configuration Structure

The configuration file has three main sections:

```yaml
# Environment variables to inject into your Nix environment
env:
  EDITOR: nvim
  BROWSER: firefox
  CUSTOM_VAR: custom_value

# Nix packages to install via home-manager
packages:
  - git
  - neovim
  - ripgrep

# External Nix flakes to integrate
flakes:
  - name: my-tools
    url: "github:username/nix-tools"
    outputs:
      - name: packages
        type: home
```

## Environment Variables

The `env` section defines custom environment variables:

```yaml
env:
  EDITOR: nvim           # Set your preferred editor
  BROWSER: firefox       # Set your preferred browser
  GO_PATH: ~/go          # Custom Go workspace
  NODE_ENV: development  # Node environment
```

These variables are:
- Injected into your Nix environment
- Available in all shells after rebuild
- Managed by home-manager

## Packages

The `packages` section lists Nix packages to install:

```yaml
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

  # Language-specific packages
  - python3Packages.requests
  - nodePackages.typescript
```

Package names:
- Must be valid Nix package identifiers
- Support attribute paths (e.g., `python3Packages.requests`)
- Are deduplicated automatically
- Must be unique in the list

## Flakes

The `flakes` section integrates external Nix flakes:

```yaml
flakes:
  - name: personal-tools     # Unique identifier
    url: "github:user/repo"  # Flake location
    follows:                 # Optional: input overrides
      nixpkgs: "nixpkgs"
    outputs:                 # What to import
      - name: packages
        type: home           # "home" or "system"
```

For detailed flake configuration, see the [Flakes Guide](/docs/user-guide/flakes/).

## Applying Configuration

After editing your configuration:

```bash
# Apply changes
camp env rebuild
```

This will:
1. Reload your `camp.yml`
2. Generate Nix configurations
3. Rebuild your environment with the new settings

## Example Configuration

Here's a complete example:

```yaml
env:
  EDITOR: nvim
  BROWSER: brave
  LANG: en_US.UTF-8
  GO111MODULE: on

packages:
  # Terminal utilities
  - git
  - curl
  - wget
  - tree
  - jq

  # Modern replacements
  - ripgrep  # grep alternative
  - fd       # find alternative
  - bat      # cat alternative
  - exa      # ls alternative

  # Development
  - neovim
  - tmux
  - direnv
  - devbox

  # Languages
  - go
  - python3
  - nodejs_20

  # Python packages
  - python3Packages.pip
  - python3Packages.virtualenv

flakes:
  - name: my-config
    url: "github:myuser/nix-config"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: homeManagerModules.default
        type: home
```

## Next Steps

- Learn more about [Package Management](/docs/user-guide/packages/)
- Explore [Flakes](/docs/user-guide/flakes/) for advanced configuration
- See the [Configuration Schema](/docs/reference/configuration-schema/) for all options
