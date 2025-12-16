---
title: "camp env rebuild"
linkTitle: "rebuild"
weight: 2
description: >
  Rebuild your development environment with current configuration
---

The `rebuild` command applies your Camp configuration by rebuilding your Nix environment.

## Usage

```bash
camp env rebuild
```

## What It Does

The rebuild process:

1. **Prepares the environment**:
   - Creates `~/.camp/nix/` directory if needed
   - Copies Nix configuration files from templates
   - Reloads your `camp.yml` configuration

2. **Compiles templates**:
   - Renders `flake.nix` with your custom data
   - Injects environment variables
   - Adds your package list
   - Integrates configured flakes

3. **Executes platform-specific rebuild**:
   - **macOS**: Runs `nix-darwin` to rebuild system configuration
   - **Linux**: Runs `home-manager` to rebuild user environment

## Prerequisites

### macOS

- Nix package manager installed
- nix-darwin configured
- Flakes enabled in Nix configuration

### Linux

- Nix package manager installed
- home-manager configured
- Flakes enabled in Nix configuration

## Rebuild Process Flow

```text
camp env rebuild
     ↓
Load ~/.camp/camp.yml
     ↓
Create ~/.camp/nix/ directory
     ↓
Copy static Nix files
     ↓
Render flake.nix template
     ↓
Execute rebuild (nix-darwin or home-manager)
     ↓
Environment updated!
```

## First Rebuild

The first time you run `camp env rebuild`, it may take several minutes:

- Nix downloads package definitions
- Builds packages that aren't in cache
- Sets up home-manager configuration
- Installs all configured packages

Subsequent rebuilds are much faster, only updating what changed.

## When to Rebuild

Run `camp env rebuild` when you:

- Edit `~/.camp/camp.yml`
- Add or remove packages
- Change environment variables
- Add or configure flakes
- Want to apply configuration updates

## Example

```bash
# Edit your configuration
vim ~/.camp/camp.yml

# Apply the changes
camp env rebuild

# Verify packages are installed
which nvim
```

## Troubleshooting

### "Nix not found"

Ensure Nix is installed and in your PATH:

```bash
which nix
# Should output: /nix/store/.../bin/nix
```

### "nix-darwin not configured" (macOS)

Set up nix-darwin first. See the [nix-darwin documentation](https://github.com/LnL7/nix-darwin).

### "home-manager not found" (Linux)

Install home-manager. See the [home-manager documentation](https://github.com/nix-community/home-manager).

### Rebuild fails with package errors

Check that package names are valid:

```bash
# Search for packages
nix search nixpkgs package-name
```

### Permission errors

On macOS, nix-darwin may require sudo for system-level changes.
Ensure your user has the necessary permissions.

## Performance Tips

- **Use Nix caches**: Configure binary caches to avoid rebuilding
- **Incremental changes**: Small, frequent rebuilds are faster
- **Clean old generations**: Periodically clean up old configurations

## Related Commands

- [`camp env`](../) - View environment commands
<!-- - [`camp env update`](../update/) - Update flake dependencies -->
<!-- - [`camp bootstrap`](../bootstrap/) - Initial setup -->
