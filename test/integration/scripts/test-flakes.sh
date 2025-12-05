#!/bin/bash
set -e

# Integration test for flake integration in camp
# Verifies that external flakes are properly configured with inputs, outputs, and args

source "$HOME/.nix-profile/etc/profile.d/nix.sh"

echo "=== Testing: Flake Integration ==="

# Ensure clean state and bootstrap first
rm -rf "$HOME/.camp"
echo "Bootstrapping environment..."
"$HOME/bin/camp" bootstrap

# Test 1: Basic flake with home output
echo "Test 1: Basic flake configuration..."
cat > "$HOME/.camp/camp.yml" <<'EOF'
env:
  EDITOR: "nvim"

flakes:
  - name: my-tools
    url: "github:user/my-tools"
    outputs:
      - name: packages
        type: home
EOF

"$HOME/bin/camp" env rebuild || true

# Verify flake input is in flake.nix
if ! grep -q "my-tools" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Flake 'my-tools' not found in flake.nix inputs"
    exit 1
fi

if ! grep -q "github:user/my-tools" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Flake URL not found in flake.nix"
    exit 1
fi

# Verify output is referenced in home-manager imports
if ! grep -q "my-tools.packages" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Flake output 'packages' not found in flake.nix"
    exit 1
fi

echo "✓ Basic flake configuration works"

# Test 2: Flake with follows
echo "Test 2: Flake with input follows..."
cat > "$HOME/.camp/camp.yml" <<'EOF'
flakes:
  - name: custom-flake
    url: "github:org/custom"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: homeManagerModules.default
        type: home
EOF

"$HOME/bin/camp" env rebuild || true

# Verify follows is in the flake
if ! grep -q "inputs.nixpkgs.follows" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Input follows not found in flake.nix"
    exit 1
fi

echo "✓ Flake follows configuration works"

# Test 3: Flake with arguments
echo "Test 3: Flake with arguments..."
cat > "$HOME/.camp/camp.yml" <<'EOF'
flakes:
  - name: parameterized
    url: "github:test/flake"
    args:
      email: "test@example.com"
      enableTools: true
      fontSize: 14
      packages: ["vim", "git"]
    outputs:
      - name: homeManagerModules.default
        type: home
EOF

"$HOME/bin/camp" env rebuild || true

# Verify arguments are passed in the flake
if ! grep -q "email" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Argument 'email' not found in flake.nix"
    exit 1
fi

if ! grep -q "test@example.com" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Argument value 'test@example.com' not found"
    exit 1
fi

if ! grep -q "enableTools" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Argument 'enableTools' not found"
    exit 1
fi

if ! grep -q "true" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Boolean argument value not found"
    exit 1
fi

if ! grep -q "fontSize" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Argument 'fontSize' not found"
    exit 1
fi

if ! grep -q "14" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Integer argument value not found"
    exit 1
fi

echo "✓ Flake arguments work"

# Test 4: Multiple flakes with different output types
echo "Test 4: Multiple flakes..."
cat > "$HOME/.camp/camp.yml" <<'EOF'
flakes:
  - name: flake-one
    url: "github:user/one"
    outputs:
      - name: packages
        type: home
  - name: flake-two
    url: "github:user/two"
    outputs:
      - name: darwinModules.default
        type: system
      - name: homeManagerModules.default
        type: home
EOF

"$HOME/bin/camp" env rebuild || true

# Verify both flakes are present
if ! grep -q "flake-one" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: flake-one not found"
    exit 1
fi

if ! grep -q "flake-two" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: flake-two not found"
    exit 1
fi

# Verify different output types
if ! grep -q "flake-one.packages" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: flake-one packages output not found"
    exit 1
fi

if ! grep -q "flake-two.darwinModules.default" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: flake-two system module not found"
    exit 1
fi

if ! grep -q "flake-two.homeManagerModules.default" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: flake-two home module not found"
    exit 1
fi

echo "✓ Multiple flakes work"

# Test 5: Validate automatic arguments (userName, hostName, home)
echo "Test 5: Automatic arguments..."
cat > "$HOME/.camp/camp.yml" <<'EOF'
flakes:
  - name: auto-args
    url: "github:test/auto"
    outputs:
      - name: modules.default
        type: home
EOF

"$HOME/bin/camp" env rebuild || true

# Verify automatic args are passed
if ! grep -q "userName" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Automatic argument 'userName' not found"
    exit 1
fi

if ! grep -q "hostName" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Automatic argument 'hostName' not found"
    exit 1
fi

if ! grep -q "home" "$HOME/.camp/nix/flake.nix"; then
    echo "ERROR: Automatic argument 'home' not found"
    exit 1
fi

echo "✓ Automatic arguments work"

echo "✓ All flake integration checks passed"
exit 0
