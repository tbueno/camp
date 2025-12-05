#!/bin/bash
set -e

# Entrypoint script for camp integration test container
# Initializes Nix environment and prepares for test execution

# Source Nix profile to ensure commands are available
if [ -f "$HOME/.nix-profile/etc/profile.d/nix.sh" ]; then
    source "$HOME/.nix-profile/etc/profile.d/nix.sh"
fi

# Verify Nix is available
if ! command -v nix &> /dev/null; then
    echo "ERROR: Nix is not available in PATH"
    exit 1
fi

# Display environment info for debugging
echo "=== Container Environment ==="
echo "User: $(whoami)"
echo "Home: $HOME"
echo "Nix version: $(nix --version)"
echo "PATH: $PATH"
echo "============================"

# Execute the command passed to the container
exec "$@"
