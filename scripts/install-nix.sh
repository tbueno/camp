#!/bin/bash
set -e

# Install Nix helper script
# Installs Nix package manager with flakes and nix-command support
# Supports macOS and Linux

SCRIPT_NAME="$(basename "$0")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_usage() {
    cat <<EOF
Usage: $SCRIPT_NAME [OPTIONS]

Install Nix package manager with experimental features enabled.

OPTIONS:
    -m, --mode MODE        Installation mode: single-user or multi-user (default: auto-detect)
    -d, --daemon           Force multi-user installation with daemon
    -s, --single-user      Force single-user installation (no daemon)
    -h, --help             Show this help message

EXAMPLES:
    # Auto-detect best installation mode
    $SCRIPT_NAME

    # Force single-user installation
    $SCRIPT_NAME --single-user

    # Force multi-user installation with daemon
    $SCRIPT_NAME --daemon

NOTES:
    - Requires internet connection
    - macOS and Linux supported
    - Enables experimental features: nix-command, flakes
    - Creates ~/.config/nix/nix.conf if needed

EOF
}

# Parse command line arguments
MODE="auto"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_usage
            exit 0
            ;;
        -d|--daemon)
            MODE="multi-user"
            shift
            ;;
        -s|--single-user)
            MODE="single-user"
            shift
            ;;
        -m|--mode)
            MODE="$2"
            shift 2
            ;;
        *)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Detect platform
PLATFORM="$(uname -s)"
case "$PLATFORM" in
    Darwin)
        OS="macOS"
        ;;
    Linux)
        OS="Linux"
        ;;
    *)
        log_error "Unsupported platform: $PLATFORM"
        exit 1
        ;;
esac

log_info "Detected platform: $OS"

# Check if Nix is already installed
if command -v nix &> /dev/null; then
    NIX_VERSION=$(nix --version 2>&1 || echo "unknown")
    log_warning "Nix is already installed: $NIX_VERSION"
    read -p "Do you want to continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Installation cancelled"
        exit 0
    fi
fi

# Determine installation mode
if [ "$MODE" = "auto" ]; then
    if [ "$OS" = "macOS" ]; then
        MODE="multi-user"
        log_info "Auto-detected mode: multi-user (recommended for macOS)"
    else
        # Check if running in Docker/container
        if [ -f /.dockerenv ] || grep -q 'docker\|lxc' /proc/1/cgroup 2>/dev/null; then
            MODE="single-user"
            log_info "Auto-detected mode: single-user (container environment)"
        else
            MODE="multi-user"
            log_info "Auto-detected mode: multi-user (recommended for Linux)"
        fi
    fi
fi

# Install Nix
log_info "Installing Nix in $MODE mode..."

if [ "$MODE" = "single-user" ]; then
    # Single-user installation (no daemon)
    if ! curl -L https://nixos.org/nix/install | sh -s -- --no-daemon; then
        log_error "Nix installation failed"
        exit 1
    fi
else
    # Multi-user installation (with daemon)
    if ! curl -L https://nixos.org/nix/install | sh -s -- --daemon; then
        log_error "Nix installation failed"
        exit 1
    fi
fi

log_success "Nix installed successfully"

# Source Nix profile
log_info "Sourcing Nix profile..."

if [ "$MODE" = "single-user" ]; then
    if [ -f "$HOME/.nix-profile/etc/profile.d/nix.sh" ]; then
        . "$HOME/.nix-profile/etc/profile.d/nix.sh"
    else
        log_warning "Could not find Nix profile script"
    fi
else
    if [ -f '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh' ]; then
        . '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'
    else
        log_warning "Could not find Nix daemon profile script"
    fi
fi

# Enable experimental features
log_info "Configuring Nix experimental features..."

NIX_CONF_DIR="$HOME/.config/nix"
NIX_CONF_FILE="$NIX_CONF_DIR/nix.conf"

mkdir -p "$NIX_CONF_DIR"

if [ -f "$NIX_CONF_FILE" ]; then
    # Check if experimental features are already enabled
    if grep -q "experimental-features" "$NIX_CONF_FILE"; then
        log_info "Experimental features already configured in $NIX_CONF_FILE"
    else
        log_info "Adding experimental features to existing $NIX_CONF_FILE"
        echo "experimental-features = nix-command flakes" >> "$NIX_CONF_FILE"
    fi
else
    log_info "Creating $NIX_CONF_FILE with experimental features"
    echo "experimental-features = nix-command flakes" > "$NIX_CONF_FILE"
fi

log_success "Nix configuration complete"

# Verify installation
log_info "Verifying Nix installation..."

if command -v nix &> /dev/null; then
    NIX_VERSION=$(nix --version)
    log_success "Nix is available: $NIX_VERSION"
else
    log_error "Nix command not found after installation"
    log_info "You may need to:"
    if [ "$MODE" = "single-user" ]; then
        log_info "  source $HOME/.nix-profile/etc/profile.d/nix.sh"
    else
        log_info "  source /nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh"
    fi
    log_info "  or restart your shell"
    exit 1
fi

# Test flakes support
log_info "Testing flakes support..."
if nix flake --help &> /dev/null; then
    log_success "Flakes support is enabled"
else
    log_warning "Flakes support test failed"
    log_info "You may need to restart your shell or source the Nix profile"
fi

# Print next steps
echo ""
log_success "Nix installation complete!"
echo ""
echo "Next steps:"
echo "  1. Restart your shell or run:"
if [ "$MODE" = "single-user" ]; then
    echo "       source $HOME/.nix-profile/etc/profile.d/nix.sh"
else
    echo "       source /nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh"
fi
echo "  2. Verify installation:"
echo "       nix --version"
echo "  3. Test flakes:"
echo "       nix flake --help"
echo "  4. Update Nix channel:"
echo "       nix-channel --update"
echo ""
log_info "For more information, visit: https://nixos.org/manual/nix/stable/"

exit 0
