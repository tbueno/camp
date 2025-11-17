# Flake Examples

This directory contains example flake configurations to help you get started with integrating external Nix flakes into your camp environment.

## What are Flakes?

Flakes are a way to package and share Nix configurations. They provide:
- **Reproducible dependencies** with lock files
- **Standardized structure** for packages and modules
- **Version pinning** for consistent environments across machines

## Quick Start

1. **Choose an example** from this directory
2. **Copy the flakes section** to your `~/.camp/camp.yml`
3. **Customize** the flake URLs and outputs
4. **Run** `camp env rebuild` to apply
5. **Update** dependencies with `camp env update`

## Examples

### Personal Packages (`personal-packages.yml`)
Shows how to integrate your own packages and configurations:
- Personal package collections
- Custom dotfiles
- Individual tool configurations

### Team Tools (`team-tools.yml`)
Shows how to share configurations across a team:
- Company-wide tools and scripts
- Shared development environments
- Team coding standards

## Flake URL Formats

Camp supports all standard Nix flake URL formats:

### GitHub (Public)
```yaml
flakes:
  - name: my-flake
    url: "github:username/repository"
```

With specific branch or tag:
```yaml
url: "github:username/repository/branch-name"
url: "github:username/repository/v1.2.3"
```

### GitHub (Private via SSH)
```yaml
flakes:
  - name: company-tools
    url: "git+ssh://git@github.com/company/private-repo.git"
```

### GitLab
```yaml
flakes:
  - name: gitlab-flake
    url: "gitlab:username/repository"
```

### Local Path
```yaml
flakes:
  - name: local-dev
    url: "path:/absolute/path/to/flake"
```

Useful for:
- Testing flakes during development
- Using flakes not published publicly

### Generic Git Repository
```yaml
flakes:
  - name: custom-git
    url: "git+https://git.example.com/repo.git"
    # Or with SSH:
    url: "git+ssh://git@git.example.com/repo.git"
```

## Configuration Structure

### Basic Flake
```yaml
flakes:
  - name: my-flake          # Unique identifier (alphanumeric, hyphens, underscores)
    url: "github:user/repo" # Flake location
    outputs:                 # What to import from this flake
      - name: packages       # Output name
        type: home           # Where to apply: "home" or "system"
```

### With Input Following
Use the same nixpkgs version as camp for consistency:

```yaml
flakes:
  - name: my-flake
    url: "github:user/repo"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: packages
        type: home
```

### Multiple Outputs
Import different parts of a flake:

```yaml
flakes:
  - name: comprehensive-flake
    url: "github:user/repo"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      # Import packages
      - name: packages
        type: home

      # Import home-manager module
      - name: homeManagerModules.default
        type: home

      # Import system configuration (macOS only)
      - name: darwinModules.system
        type: system
```

## Output Types

### `type: home`
Applies to **user environment** (home-manager):
- Available on both macOS and Linux
- User-level packages and configurations
- Dotfiles, shell configs, development tools

**Common output names:**
- `packages` - Package sets
- `homeManagerModules.default` - Home-manager modules
- `homeManagerModules.{name}` - Named modules

### `type: system`
Applies to **system level** (nix-darwin):
- Available on macOS only
- System-wide settings and services
- Requires sudo/admin privileges

**Common output names:**
- `darwinModules.default` - Darwin modules
- `darwinModules.{name}` - Named modules

## Common Use Cases

### 1. Personal Development Tools
```yaml
flakes:
  - name: dev-tools
    url: "github:myuser/dev-environment"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: packages
        type: home
```

### 2. Shared Team Configuration
```yaml
flakes:
  - name: team-config
    url: "git+ssh://git@github.com/company/team-nix.git"
    follows:
      nixpkgs: "nixpkgs"
    outputs:
      - name: homeManagerModules.team
        type: home
```

### 3. Language-Specific Environment
```yaml
flakes:
  - name: python-env
    url: "github:user/python-flake"
    outputs:
      - name: homeManagerModules.python
        type: home
```

### 4. macOS System Configuration
```yaml
flakes:
  - name: macos-config
    url: "github:user/darwin-config"
    outputs:
      # System-level settings
      - name: darwinModules.system
        type: system
      # User-level settings
      - name: homeManagerModules.user
        type: home
```

## Workflow

### Adding a Flake
1. Edit `~/.camp/camp.yml`
2. Add flake definition to `flakes:` section
3. Run `camp env rebuild`
4. Verify the flake is integrated

### Updating Flakes
```bash
# Update all flake dependencies to latest versions
camp env update

# Apply the updates
camp env rebuild
```

### Removing a Flake
1. Remove flake definition from `~/.camp/camp.yml`
2. Run `camp env rebuild`

## Troubleshooting

### "duplicate flake name"
Each flake must have a unique name:
```yaml
# ❌ Bad: duplicate names
flakes:
  - name: my-flake
    url: "github:user/repo1"
  - name: my-flake  # Duplicate!
    url: "github:user/repo2"

# ✅ Good: unique names
flakes:
  - name: flake-one
    url: "github:user/repo1"
  - name: flake-two
    url: "github:user/repo2"
```

### "invalid flake name"
Flake names must be valid Nix identifiers:
```yaml
# ❌ Bad: contains invalid characters
- name: "my.flake"    # No dots
- name: "my flake"    # No spaces
- name: "my@flake"    # No special chars

# ✅ Good: valid identifiers
- name: "my-flake"    # Hyphens OK
- name: "my_flake"    # Underscores OK
- name: "MyFlake123"  # Alphanumeric OK
```

### "empty URL"
Every flake must have a URL:
```yaml
# ❌ Bad: missing URL
- name: my-flake
  outputs:
    - name: packages
      type: home

# ✅ Good: URL specified
- name: my-flake
  url: "github:user/repo"
  outputs:
    - name: packages
      type: home
```

### Private Repository Access
For private repos via SSH:
1. Ensure SSH key is configured
2. Test SSH access: `ssh -T git@github.com`
3. Use `git+ssh://` URL format
4. Add to known hosts if needed

## Additional Resources

- [Nix Flakes Documentation](https://nixos.wiki/wiki/Flakes)
- [Nix Flake URL Schemas](https://nixos.org/manual/nix/stable/command-ref/new-cli/nix3-flake.html#flake-references)
- [Home Manager](https://github.com/nix-community/home-manager)
- [nix-darwin](https://github.com/LnL7/nix-darwin)

## Need Help?

If you encounter issues:
1. Check your `~/.camp/camp.yml` syntax
2. Run `camp env rebuild` and read error messages
3. Verify flake URLs are accessible
4. Check flake output names match what's exported
