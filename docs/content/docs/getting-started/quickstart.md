---
title: "Quick Start"
linkTitle: "Quick Start"
weight: 2
description: >
  Set up your first Camp environment in minutes
---

This guide walks you through setting up your first development environment with Camp.

## 1. Bootstrap Your Environment

The `bootstrap` command sets up your initial Camp environment:

```bash
camp bootstrap
```

This will:
- Create the `~/.camp` directory
- Generate initial configuration files
- Set up Nix flake templates

## 2. View System Information

Check your system details:

```bash
camp env
```

This displays:
- Username and home directory
- System platform and architecture
- Hostname
- Shell
- Environment variables from your configuration

## 3. Create Your Configuration

Create or edit `~/.camp/camp.yml`:

```yaml
# Environment variables to inject
env:
  EDITOR: nvim
  BROWSER: firefox

# Nix packages to install
packages:
  - git
  - neovim
  - ripgrep
  - fd
  - bat
```

## 4. Rebuild Your Environment

Apply your configuration:

```bash
camp env rebuild
```

This command:
1. Copies Nix configuration files to `~/.camp/nix/`
2. Renders `flake.nix` with your custom settings
3. Executes the platform-specific rebuild:
   - **macOS**: Uses `nix-darwin`
   - **Linux**: Uses `home-manager`

The rebuild process may take a few minutes on first run as it downloads and builds packages.

## 5. Verify Your Setup

After the rebuild completes:

```bash
# Check that packages are available
which nvim
which rg

# Verify environment variables
echo $EDITOR
```

## Next Steps

Now that you have a basic environment set up, explore more features:

- **[Configuration Guide](../configuration/)** - Learn about all configuration options
- **[User Guide](/docs/user-guide/)** - Dive deeper into Camp's features
- **[Package Management](/docs/user-guide/packages/)** - Add more packages
- **[Flakes](/docs/user-guide/flakes/)** - Extend with external Nix flakes

## Common Issues

### "Nix not found"

Make sure Nix is installed and in your PATH. See the [Installation guide](../installation/).

### "nix-darwin not configured" (macOS)

You need to set up nix-darwin before using Camp. Follow the [nix-darwin installation guide](https://github.com/LnL7/nix-darwin).

### "home-manager not found" (Linux)

Install home-manager first. See the [home-manager installation guide](https://github.com/nix-community/home-manager).
