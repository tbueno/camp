#!/bin/bash
set -e

# Integration test for camp env nuke command
# Verifies that nuke properly cleans up the camp environment

source "$HOME/.nix-profile/etc/profile.d/nix.sh"

echo "=== Testing: camp env nuke ==="

# First, set up a full environment
rm -rf "$HOME/.camp" "$HOME/.envrc"

echo "Setting up environment..."
"$HOME/bin/camp" bootstrap

# Create a more complete config to test cleanup
cat > "$HOME/.camp/camp.yml" <<'EOF'
env:
  EDITOR: "nvim"
  BROWSER: "firefox"

packages:
  - git
  - ripgrep
  - neovim

flakes:
  - name: test-flake
    url: "github:test/flake"
    outputs:
      - name: packages
        type: home
EOF

# Run rebuild to ensure all files are generated
"$HOME/bin/camp" env rebuild || true

# Verify environment exists before nuking
echo "Verifying environment exists..."
if [ ! -d "$HOME/.camp" ]; then
    echo "ERROR: .camp directory doesn't exist"
    exit 1
fi

if [ ! -f "$HOME/.camp/camp.yml" ]; then
    echo "ERROR: camp.yml doesn't exist"
    exit 1
fi

if [ ! -d "$HOME/.camp/nix" ]; then
    echo "ERROR: .camp/nix directory doesn't exist"
    exit 1
fi

if [ ! -f "$HOME/.envrc" ]; then
    echo "ERROR: .envrc doesn't exist"
    exit 1
fi

# Count files before nuke
FILE_COUNT_BEFORE=$(find "$HOME/.camp" -type f | wc -l)
echo "Files before nuke: $FILE_COUNT_BEFORE"

# Run nuke command
echo "Running: camp env nuke --yes"
"$HOME/bin/camp" env nuke --yes

# Verify everything is cleaned up
echo "Verifying cleanup..."

if [ -d "$HOME/.camp" ]; then
    # Check if directory is empty or has leftover files
    REMAINING_FILES=$(find "$HOME/.camp" -type f 2>/dev/null | wc -l || echo "0")
    if [ "$REMAINING_FILES" -gt 0 ]; then
        echo "ERROR: .camp directory still has files after nuke"
        find "$HOME/.camp" -type f
        exit 1
    fi
    echo "Note: .camp directory exists but is empty (acceptable)"
else
    echo "✓ .camp directory removed"
fi

if [ -f "$HOME/.envrc" ]; then
    echo "ERROR: .envrc still exists after nuke"
    exit 1
fi

echo "✓ Environment cleaned up successfully"

# Test nuke on non-existent environment (should not error)
echo "Testing nuke on clean system..."
rm -rf "$HOME/.camp"  # Ensure it's gone

if "$HOME/bin/camp" env nuke --yes; then
    echo "✓ Nuke on clean system succeeded (idempotent)"
else
    echo "ERROR: Nuke failed on clean system"
    exit 1
fi

# Test that we can bootstrap again after nuke
echo "Testing re-bootstrap after nuke..."
if "$HOME/bin/camp" bootstrap; then
    echo "✓ Can bootstrap again after nuke"
else
    echo "ERROR: Failed to bootstrap after nuke"
    exit 1
fi

# Verify the new bootstrap created files
if [ ! -d "$HOME/.camp" ] || [ ! -f "$HOME/.camp/camp.yml" ]; then
    echo "ERROR: Re-bootstrap didn't create expected files"
    exit 1
fi

echo "✓ All nuke checks passed"
exit 0
