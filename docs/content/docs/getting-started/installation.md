---
title: "Installation"
linkTitle: "Installation"
weight: 1
description: >
  How to install Camp on your system
---

## Prerequisites

Before installing Camp, ensure you have:

- **Go 1.24.4 or later** - Camp is built with Go
- **Nix package manager** - Required for environment management
  - macOS: Nix with nix-darwin configured
  - Linux: Nix with home-manager configured
- **Git** - For cloning the repository

## Install Nix

If you don't have Nix installed yet:

### macOS or Linux

```bash
# Install Nix with flakes enabled
sh <(curl -L https://nixos.org/nix/install) --daemon

# Enable flakes (add to ~/.config/nix/nix.conf or /etc/nix/nix.conf)
experimental-features = nix-command flakes
```

### Set up nix-darwin (macOS only)

Follow the [nix-darwin installation guide](https://github.com/LnL7/nix-darwin).

### Set up home-manager (Linux or macOS)

Follow the [home-manager installation guide](https://github.com/nix-community/home-manager).

## Install Camp

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/tbueno/camp
cd camp

# Build the binary and install it to /usr/local/bin (asks for your password)
make install
```

### Option 2: Install with Go

```bash
go install github.com/tbueno/camp@latest
```

This will install Camp to your `$GOPATH/bin` directory.

The Nix templates Camp writes to `~/.camp/nix/` are built into the binary, so you can run
`camp` from any directory once it is on your `PATH`. Rebuild and reinstall the binary to pick
up template changes.

## Verify Installation

```bash
# Check that Camp is installed
camp --version

# View available commands
camp --help
```

## Next Steps

- [Quick Start Guide](../quickstart/) - Set up your first environment
- [Configuration](../configuration/) - Learn about `camp.yml`
