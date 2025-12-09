# Camp Documentation

This directory contains the Camp project documentation built with [Hugo](https://gohugo.io/) and the [Docsy](https://www.docsy.dev/) theme.

## Prerequisites

- [Hugo](https://gohugo.io/installation/) (extended version)
- [Node.js](https://nodejs.org/) v20 or later
- npm (comes with Node.js)

## Setup

1. Install npm dependencies:
   ```bash
   npm ci
   ```

2. Initialize git submodules (for the Docsy theme):
   ```bash
   git submodule update --init --recursive
   ```

## Building

### Quick Build

Use the provided build script:

```bash
./build.sh
```

This script automatically:
- Patches the Docsy theme to disable Hugo Module imports (we use npm packages instead)
- Builds the documentation with minification

### Manual Build

If you prefer to build manually:

```bash
# Patch Docsy theme (only needed once after git submodule update)
sed -i.bak 's/disable: false/disable: true/g' themes/docsy/hugo.yaml

# Build
hugo --minify
```

## Development Server

To run a local development server with live reload:

```bash
# Make sure the theme is patched first (see above)
./build.sh  # patches theme if needed

# Then run the server
hugo server
```

Visit http://localhost:1313/camp/ to view the documentation.

## Architecture Notes

### Why patch the Docsy theme?

The Docsy theme's `hugo.yaml` defines Hugo Module imports for Bootstrap and Font-Awesome. However, we use npm packages for these dependencies instead of Hugo Modules. The patch disables these module imports to prevent Hugo from trying to fetch them as Hugo Modules, which would fail in CI.

### npm vs Hugo Modules

We use npm packages because:
- More stable and predictable dependency management
- Works consistently across local development and CI
- Easier to manage versions with package-lock.json

The module mounts in `config.toml` map the npm packages to the vendor directories expected by Docsy.

## CI/CD

GitHub Actions workflows automatically:
1. Check out the repository with submodules
2. Patch the Docsy theme
3. Install npm dependencies
4. Build and deploy the documentation

See `.github/workflows/deploy-docs.yml` for details.
