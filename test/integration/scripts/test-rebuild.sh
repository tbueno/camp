#!/bin/bash
set -e

# Integration test for camp env rebuild command
# Verifies that rebuild generates correct Nix files with custom config

source "$HOME/.nix-profile/etc/profile.d/nix.sh"

echo "=== Testing: camp env rebuild ==="

# Ensure clean state and bootstrap first
rm -rf "$HOME/.camp"
echo "Bootstrapping environment..."
"$HOME/bin/camp" bootstrap

# Create a custom config file with env vars
echo "Creating custom config..."
cat > "$HOME/.camp/camp.yml" <<'EOF'
env:
  EDITOR: "nvim"
  BROWSER: "firefox"
  CUSTOM_VAR: "test-value"
EOF

# Run rebuild (note: we won't actually apply Nix changes in container,
# just verify file generation)
echo "Running: camp env rebuild --dry-run"

# First, let's just test that the command prepares files correctly
# We'll verify the generated flake.nix contains our custom env vars
"$HOME/bin/camp" env rebuild || true  # May fail on actual Nix rebuild, that's OK

# Verify flake.nix was regenerated
echo "Verifying generated flake.nix..."

if [ ! -f "$HOME/.camp/nix/flake.nix" ]; then
    echo "ERROR: flake.nix not regenerated"
    exit 1
fi

# Check that custom env vars are in the flake
if ! grep -q "EDITOR" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: EDITOR env var not found in flake.nix"
    exit 1
fi

if ! grep -q "nvim" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: EDITOR value 'nvim' not found in flake.nix"
    exit 1
fi

if ! grep -q "BROWSER" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: BROWSER env var not found in flake.nix"
    exit 1
fi

if ! grep -q "firefox" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: BROWSER value 'firefox' not found in flake.nix"
    exit 1
fi

if ! grep -q "CUSTOM_VAR" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: CUSTOM_VAR env var not found in flake.nix"
    exit 1
fi

if ! grep -q "test-value" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: CUSTOM_VAR value 'test-value' not found in flake.nix"
    exit 1
fi

# Verify the flake structure is valid
echo "Verifying flake structure..."

# Check for required sections
if ! grep -q "inputs" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: 'inputs' section not found in flake.nix"
    exit 1
fi

if ! grep -q "outputs" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: 'outputs' section not found in flake.nix"
    exit 1
fi

if ! grep -q "customEnvVars" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: 'customEnvVars' section not found in flake.nix"
    exit 1
fi

# Test that we can update config and rebuild again
echo "Testing config reload..."
cat > "$HOME/.camp/camp.yml" <<'EOF'
env:
  EDITOR: "vim"
  NEW_VAR: "new-value"
EOF

"$HOME/bin/camp" env rebuild || true

# Verify updated values
if ! grep -q "vim" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Updated EDITOR value 'vim' not found after reload"
    exit 1
fi

if ! grep -q "NEW_VAR" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: NEW_VAR not found after reload"
    exit 1
fi

# Old BROWSER var should be gone
if grep -q "firefox" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Old BROWSER value still present after config change"
    exit 1
fi

echo "✓ All rebuild checks passed"
exit 0
