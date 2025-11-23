# ⛺ CAMP CLI

> Your all-in-one development environment manager

[![Go Version](https://img.shields.io/badge/go-1.24.4-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)

## Description

Camp is a command-line tool designed to help developers manage their isolated development environments. Built with Go, Camp provides essential system information and environment management utilities.

Camp relies on existing tools like [direnv](https://direnv.net/) and [devbox](https://www.jetify.com/devbox) for environment setup and management.

## Features

### Current Features
- **System Information**: Get detailed information about your system architecture and operating system
- **Environment Configuration**: Manage custom environment variables via `camp.yml`
- **Nix Integration**: Rebuild development environments using Nix configurations
- **Platform Support**: Works on both macOS (via nix-darwin) and Linux (via home-manager)
- **Template System**: Dynamic Nix configuration generation with user-specific data

### Planned Features
- Environment isolation using direnv and devbox
- Development workflow automation
- Project-specific environment configuration
- Environment configuration remote sharing

## Installation

### Prerequisites
- Go 1.24.4 or later

### Build from Source
```bash
git clone https://github.com/tbueno/camp
cd camp
go build -o camp
```

### Install
```bash
go install
```

## Usage

### Commands

**Environment Information**
```bash
# Display system and environment information
camp env
```

**Environment Rebuild**
```bash
# Rebuild development environment with latest configuration
camp env rebuild
```

This command:
1. Copies Nix configuration files to `~/.camp/nix/`
2. Renders `flake.nix` with your custom environment variables from `camp.yml`
3. Executes platform-specific rebuild:
   - **macOS**: Uses `nix-darwin` to rebuild system configuration
   - **Linux**: Uses `home-manager` to rebuild user environment

**Prerequisites for Rebuild:**
- Nix package manager must be installed
- macOS: nix-darwin must be configured
- Linux: home-manager must be configured

**Bootstrap**
```bash
# Initial environment setup
camp bootstrap
```

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
  `

**Configuration Sections:**
- **`env`**: Environment variables injected into your Nix environment
- **`packages`**: Nix packages to install via home-manager
- **`flakes`**: External Nix flakes to integrate (see documentation for details)

All packages will be installed to your home environment and available in your PATH after running `camp env rebuild`.

### Command Help

Get help for any command:
```bash
camp --help
camp env --help
camp env rebuild --help
```

## Development

### Testing

Run the test suite:
```bash
go test ./...
```
## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Roadmap

Camp is in active development. Future releases will include:
- Container-based environment isolation
- Configuration management
- Integration with popular development tools
- Cross-platform environment synchronization
- Plugin system for extensibility

---

*Built with ❤️ using Go and Cobra*