---
title: "Package Management"
linkTitle: "Packages"
weight: 3
description: >
  Managing Nix packages with Camp
---

Camp allows you to declaratively manage Nix packages through your configuration file.

## Overview

Packages are defined as a simple list in your `camp.yml`:

```yaml
packages:
  - git
  - neovim
  - ripgrep
```

After running `camp env rebuild`, these packages are:
- Installed via home-manager
- Available in your PATH
- Managed declaratively (no manual `nix-env -i` needed)

## Adding Packages

1. Edit your configuration:

```yaml
packages:
  - git
  - curl
  - wget
  - neovim
  - ripgrep
  - fd
  - bat
```

2. Apply the changes:

```bash
camp env rebuild
```

3. Verify installation:

```bash
which nvim
# Output: /nix/store/.../bin/nvim
```

## Package Names

Package names must be valid Nix package identifiers from nixpkgs.

### Valid Examples

```yaml
packages:
  - git                      # Simple package name
  - python3                  # Version-specific
  - nodejs_20                # Underscore for versions
  - package-with-hyphens     # Hyphens allowed
  - package_with_underscores # Underscores allowed
```

### Attribute Paths

You can specify nested packages using dot notation:

```yaml
packages:
  # Python packages
  - python3Packages.requests
  - python3Packages.flask
  - python3Packages.numpy

  # Node packages
  - nodePackages.typescript
  - nodePackages.prettier

  # Haskell packages
  - haskellPackages.pandoc
```

### Invalid Examples

```yaml
packages:
  - "my package"    # ❌ Spaces not allowed
  - my@package      # ❌ Special characters not allowed
  - my/package      # ❌ Slashes not allowed
```

## Searching for Packages

Find packages in nixpkgs:

```bash
# Search for a package
nix search nixpkgs package-name

# Example: search for ripgrep
nix search nixpkgs ripgrep

# Browse packages online
# https://search.nixos.org/packages
```

## Package Validation

Camp validates your package list:

- **No duplicates**: Each package must be unique
- **Valid identifiers**: Only letters, numbers, hyphens, underscores, dots
- **Not empty**: Package names cannot be empty or whitespace

Invalid configurations are caught when you run `camp env rebuild`.

## Common Package Categories

### Terminal Utilities

```yaml
packages:
  - git
  - curl
  - wget
  - tree
  - jq
  - yq
```

### Modern CLI Replacements

```yaml
packages:
  - ripgrep  # grep alternative
  - fd       # find alternative
  - bat      # cat alternative
  - exa      # ls alternative
  - dust     # du alternative
  - zoxide   # cd alternative
```

### Development Tools

```yaml
packages:
  - neovim
  - tmux
  - direnv
  - devbox
  - gh       # GitHub CLI
```

### Programming Languages

```yaml
packages:
  # Go
  - go

  # Python
  - python3
  - python3Packages.pip
  - python3Packages.virtualenv

  # Node.js
  - nodejs_20
  - nodePackages.npm

  # Rust
  - rustc
  - cargo
```

### Language Servers (LSP)

```yaml
packages:
  - gopls                      # Go
  - pyright                    # Python
  - nodePackages.typescript-language-server  # TypeScript
  - rust-analyzer              # Rust
  - lua-language-server        # Lua
```

## Platform Differences

### macOS (nix-darwin)

Packages are installed via home-manager within nix-darwin:

```nix
home.packages = with pkgs; [
  # Your configured packages
];
```

### Linux (home-manager)

Packages are installed directly via home-manager:

```nix
home.packages = with pkgs; [
  # Your configured packages
];
```

Both platforms work identically from the user's perspective.

## Managing Package Versions

### Using Specific Versions

Some packages have version-specific variants:

```yaml
packages:
  - nodejs_20  # Node.js 20.x
  - nodejs_18  # Node.js 18.x
  - python311  # Python 3.11
  - python312  # Python 3.12
```

Check available versions:

```bash
nix search nixpkgs nodejs
```

### Pinning Nixpkgs Version

Camp uses nixpkgs-unstable by default. To use a specific nixpkgs version, you would need to modify the generated flake.nix (advanced).

## Removing Packages

Simply remove from your configuration and rebuild:

```yaml
packages:
  - git
  - neovim
  # ripgrep removed
```

```bash
camp env rebuild
```

The package remains in the Nix store but is no longer in your PATH.

## Example Configurations

### Minimal Setup

```yaml
packages:
  - git
  - curl
  - wget
```

### Web Developer

```yaml
packages:
  - git
  - nodejs_20
  - nodePackages.npm
  - nodePackages.typescript
  - nodePackages.prettier
  - nodePackages.eslint
```

### Python Developer

```yaml
packages:
  - git
  - python3
  - python3Packages.pip
  - python3Packages.virtualenv
  - python3Packages.ipython
  - pyright
```

### Go Developer

```yaml
packages:
  - git
  - go
  - gopls
  - golangci-lint
  - delve
```

### DevOps Engineer

```yaml
packages:
  - git
  - docker
  - kubectl
  - terraform
  - ansible
  - jq
  - yq
```

## Troubleshooting

### Package Not Found

```
error: attribute 'packagename' missing
```

**Solution**: Verify the package name exists in nixpkgs:

```bash
nix search nixpkgs packagename
```

### Duplicate Package Error

```
Error: duplicate package 'git' - package names must be unique
```

**Solution**: Remove the duplicate from your `camp.yml`.

### Invalid Package Name

```
Error: package 'my package' has invalid format
```

**Solution**: Remove spaces and special characters. Use hyphens or underscores instead.

## Best Practices

1. **Keep it declarative**: Don't use `nix-env -i`, manage all packages in `camp.yml`
2. **Group logically**: Organize packages by category with comments
3. **Use attribute paths**: Prefer `python3Packages.requests` over separate tools
4. **Search before adding**: Verify package names with `nix search`
5. **Start minimal**: Add packages as you need them
6. **Document custom choices**: Add comments for non-obvious packages

## Next Steps

- Learn about [Flakes](/docs/user-guide/flakes/) for more advanced package management
- See the [Configuration Schema](/docs/reference/configuration-schema/) for all options
- Check out [Commands](/docs/user-guide/commands/) for rebuild workflow
