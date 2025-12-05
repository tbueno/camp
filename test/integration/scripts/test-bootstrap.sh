#!/bin/bash
set -e

# Integration test for camp bootstrap command
# Verifies that bootstrap creates all necessary files and directories

source "$HOME/.nix-profile/etc/profile.d/nix.sh"

echo "=== Testing: camp bootstrap ==="

# Ensure clean state
rm -rf "$HOME/.camp"

# Run bootstrap
echo "Running: camp bootstrap"
if ! "$HOME/bin/camp" bootstrap; then
    echo "ERROR: camp bootstrap failed"
    exit 1
fi

# Verify directory structure
echo "Verifying directory structure..."

if [ ! -d "$HOME/.camp" ]; then
    echo "ERROR: .camp directory not created"
    exit 1
fi

if [ ! -d "$HOME/.camp/nix" ]; then
    echo "ERROR: .camp/nix directory not created"
    exit 1
fi

# Verify configuration file
echo "Verifying config file..."

if [ ! -f "$HOME/.camp/camp.yml" ]; then
    echo "ERROR: camp.yml not created"
    exit 1
fi

# Verify Nix files
echo "Verifying Nix files..."

if [ ! -f "$HOME/.camp/nix/flake.nix" ]; then
    echo "ERROR: flake.nix not created"
    exit 1
fi

# Check that flake.nix is valid (basic syntax check)
if ! grep -q "description" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: flake.nix appears invalid"
    exit 1
fi

# Verify platform-specific files exist
if [ -f "$HOME/.camp/nix/mac.nix" ] || [ -f "$HOME/.camp/nix/linux.nix" ]; then
    echo "Platform-specific Nix files found"
else
    echo "ERROR: No platform-specific Nix files found"
    exit 1
fi

# Verify modules directory
if [ ! -d "$HOME/.camp/nix/modules" ]; then
    echo "ERROR: modules directory not created"
    exit 1
fi

if [ ! -f "$HOME/.camp/nix/modules/common.nix" ]; then
    echo "ERROR: common.nix not created"
    exit 1
fi

# Verify .envrc file
if [ ! -f "$HOME/.envrc" ]; then
    echo "ERROR: .envrc not created"
    exit 1
fi

# Check that .envrc contains direnv setup
if ! grep -q "use flake" "$HOME/.envrc"; then
    echo "ERROR: .envrc doesn't contain flake directive"
    exit 1
fi

echo "✓ All bootstrap checks passed"
echo ""
echo "Created files:"
find "$HOME/.camp" -type f | sort

exit 0
