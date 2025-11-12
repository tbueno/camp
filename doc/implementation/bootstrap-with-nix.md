# Bootstrap with Nix - Implementation Plan

## Overview
Port the bootstrap feature from optishell to camp, replacing any mention of "optishell" with "camp" and adapting it to the current camp structure.

## Current State Analysis
- **Camp**: Has a simple bootstrap that installs tools (direnv, nix, devbox) but doesn't set up home directory structure or Nix configuration
- **Optishell**: Has a full bootstrap that copies files to `$HOME/.optishell`, sets up Nix configuration, and runs platform-specific setups

## Implementation Strategy

### 1. Create Templates Structure (`templates/initial/`)
- **Include:** `flake.nix`, `darwin.nix`, `home.nix` (no yml files)
- **Exclude:** `camp.yml.example` (implement later)
- Replace "optishell" references with "camp"
- Include `bin/` directory with setup scripts

### 2. Create Utils Package (`internal/utils/`)
- Port file operations: SaveFile, ReplaceInContent, BackupFile
- Port command execution: RunCommand, CommandReturn
- Port TemplDir interface for template handling
- Port HostName function

### 3. Extend User Type (`internal/system/types.go`)
- Add fields: `HomeDir`, `Platform`, `Architecture`, `Shell`, `HostName`
- Create NewUser() constructor
- **Skip:** EnvVars, Flakes fields (yml-related)

### 4. Update Current Bootstrap Tool Installation
- **Remove:** direnv and devbox from default applications
- **Keep only:** nix installation
- These tools will be installed via initial flake.nix instead

### 5. Enhance Bootstrap Command (`cmd/bootstrap.go`)
- ~~Add `--setup-home` flag for full environment setup~~ (Changed: Always set up home)
- Bootstrap always sets up full home directory with Nix configuration
- No flag needed - simplified user experience

### 6. Create macOS-Specific Bootstrap Logic (`internal/system/bootstrap.go`)
- Add `bootstrapHome()` - creates `$HOME/.camp/{nix,bin}`
- Add `bootstrapMac()` - sets up nix-darwin configuration
- **Linux support:** Add TODO comments, skip implementation
- Copy and template Nix files for macOS

### 7. Add Bin Scripts (`templates/initial/bin/`)
- Port `bootstrap` and `install_nix` scripts
- Replace optishell paths with camp paths (`$HOME/.camp`)
- Make executable during copy

### 8. File Mappings
```
$HOME/.optishell → $HOME/.camp
All "optishell" strings → "camp"
github.com/Optibus/optishell → camp/internal
Remove: camp.yml.example, EnvVars, Flakes handling
Keep: Only nix in default bootstrap applications
```

### 9. Linux Handling
- Add `bootstrapLinux()` function with TODO comment
- Skip Linux-specific logic implementation
- Return appropriate "not implemented" error for Linux

## Integration Points
- Bootstrap now always sets up home directory (no optional flag)
- Simplified user experience - one command does everything
- Use `cmd.OutOrStdout()` for all output (testing compatibility)
- `--dry-run` flag available for testing without making changes

## Implementation Notes
1. **No YML support:** Skip yml file handling for now, implement later
2. **Linux TODO:** Leave TODO comments for Linux parts, do not implement yet
3. **Simplified tools:** Change current tool installation to install only nix, remove direnv and devbox (will be installed via flake.nix)
4. **Simplified bootstrap:** Removed `--setup-home` flag - bootstrap always sets up complete environment

This simplified approach focuses on macOS + nix-darwin setup while keeping the foundation ready for yml configuration and Linux support later.

## Final Implementation Status

### ✅ Completed
- Templates structure with Nix configurations
- Utils package with file and command operations
- Extended User type
- Bootstrap always sets up home directory
- macOS bootstrap with nix-darwin
- Linux stub with TODO
- Comprehensive tests
- All "optishell" references replaced with "camp"

### Usage
```bash
# Bootstrap your environment (sets up everything)
camp bootstrap

# Preview what would be done
camp bootstrap --dry-run
```