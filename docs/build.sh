#!/bin/bash
# Helper script to build documentation locally
# This patches the Docsy theme to disable Hugo Module imports (we use npm instead)

set -e

cd "$(dirname "$0")"

# Patch Docsy theme if not already patched
if grep -q "disable: false" themes/docsy/hugo.yaml 2>/dev/null; then
    echo "Patching Docsy theme to disable Hugo Module imports..."
    sed -i.bak 's/disable: false/disable: true/g' themes/docsy/hugo.yaml
fi

# Build
echo "Building documentation..."
hugo --minify "$@"
