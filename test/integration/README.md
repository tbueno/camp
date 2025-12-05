# Camp Integration Tests

This directory contains the integration test suite for camp. These tests validate the complete workflow in a realistic environment using Docker containers with Nix installed.

## Overview

Integration tests verify end-to-end functionality by:
- Building the camp binary
- Running tests in isolated Docker containers
- Testing against Ubuntu with Nix pre-installed
- Validating the full workflow: bootstrap → config → rebuild → nuke

## Directory Structure

```
test/integration/
├── docker/
│   ├── Dockerfile.ubuntu    # Ubuntu + Nix container image
│   └── entrypoint.sh         # Container initialization script
├── scripts/
│   ├── run-tests.sh          # Main test orchestrator
│   ├── test-bootstrap.sh     # Tests bootstrap command
│   ├── test-rebuild.sh       # Tests env rebuild
│   ├── test-packages.sh      # Tests package management
│   ├── test-flakes.sh        # Tests flake integration
│   └── test-nuke.sh          # Tests environment cleanup
├── fixtures/
│   ├── basic-config.yml      # Minimal config (env vars only)
│   ├── with-packages.yml     # Config with packages
│   ├── with-flakes.yml       # Config with flakes
│   └── full-config.yml       # Complete configuration
└── README.md                 # This file
```

## Running Tests

### Run All Integration Tests

```bash
# From project root
make test-integration

# Or directly
./test/integration/scripts/run-tests.sh
```

### Run Individual Test Scenarios

You can run specific tests by executing the Docker container manually:

```bash
# Build Docker image first
cd test/integration/docker
docker build -t camp-integration-test:latest .

# Build camp binary
cd ../../..
go build -o camp main.go

# Run a specific test
docker run --rm \
  -v "$(pwd)/camp:/home/testuser/bin/camp:ro" \
  -v "$(pwd)/test/integration/scripts:/home/testuser/tests:ro" \
  -v "$(pwd)/test/integration/fixtures:/home/testuser/fixtures:ro" \
  camp-integration-test:latest \
  /bin/bash /home/testuser/tests/test-bootstrap.sh
```

### Available Test Scripts

| Script | Description | What It Tests |
|--------|-------------|---------------|
| `test-bootstrap.sh` | Bootstrap workflow | Directory creation, file generation, initial setup |
| `test-rebuild.sh` | Environment rebuild | Config reload, template rendering, env var injection |
| `test-packages.sh` | Package management | Package declaration, validation, Nix integration |
| `test-flakes.sh` | Flake integration | Flake inputs, outputs, arguments, follows |
| `test-nuke.sh` | Environment cleanup | File removal, idempotency, re-bootstrap |

## Test Scenarios

### 1. Bootstrap Test
**File**: `scripts/test-bootstrap.sh`

Validates:
- ✅ Creates `~/.camp/` directory
- ✅ Creates `camp.yml` config file
- ✅ Copies Nix templates to `~/.camp/nix/`
- ✅ Generates `flake.nix`, `mac.nix`, `linux.nix`
- ✅ Creates `modules/common.nix`
- ✅ Creates `.envrc` with direnv setup

### 2. Rebuild Test
**File**: `scripts/test-rebuild.sh`

Validates:
- ✅ Reloads config from `camp.yml`
- ✅ Renders templates with custom env vars
- ✅ Generates valid `flake.nix`
- ✅ Updates config and re-renders correctly
- ✅ Removes old values on config change

### 3. Package Test
**File**: `scripts/test-packages.sh`

Validates:
- ✅ Packages declared in config appear in `flake.nix`
- ✅ Duplicate package detection works
- ✅ Invalid package name validation
- ✅ Attribute path packages (e.g., `python3Packages.requests`)
- ✅ Package array rendering in Nix syntax

### 4. Flake Test
**File**: `scripts/test-flakes.sh`

Validates:
- ✅ Flake inputs appear in generated `flake.nix`
- ✅ Flake outputs are imported correctly
- ✅ Input `follows` directive works
- ✅ Flake arguments (string, bool, int, list)
- ✅ Automatic arguments (userName, hostName, home)
- ✅ Multiple flakes with different output types
- ✅ System vs home output routing

### 5. Nuke Test
**File**: `scripts/test-nuke.sh`

Validates:
- ✅ Removes `~/.camp/` directory
- ✅ Removes `.envrc` file
- ✅ Cleans up all generated files
- ✅ Idempotent (safe to run on clean system)
- ✅ Can re-bootstrap after nuke

## Test Fixtures

Pre-configured YAML files for testing different scenarios:

### `basic-config.yml`
Minimal configuration with only environment variables. Used for testing basic setup.

### `with-packages.yml`
Configuration with multiple packages including:
- Core utilities (git, curl, wget)
- Development tools (neovim, ripgrep, fd)
- Programming languages (python3, nodejs, go)

### `with-flakes.yml`
Configuration with external flakes demonstrating:
- Home-level outputs
- System-level outputs
- Flake arguments
- Input follows

### `full-config.yml`
Complete configuration combining:
- Extensive environment variables
- 20+ packages
- Multiple flakes with arguments
- All feature combinations

## Docker Environment

### Image: `camp-integration-test:latest`

**Base**: Ubuntu latest
**User**: `testuser` (non-root)
**Installed**:
- Nix (single-user mode)
- Git, curl, xz-utils
- Direnv
- Nix experimental features enabled (flakes, nix-command)

**Mounts**:
- `/home/testuser/bin/camp` - Camp binary (read-only)
- `/home/testuser/tests` - Test scripts (read-only)
- `/home/testuser/fixtures` - Config fixtures (read-only)

## CI/CD Integration

Integration tests run in GitHub Actions on:
- Push to main branch
- Pull requests
- Manual workflow dispatch

**Platforms tested**:
- Ubuntu (Docker container)
- macOS (native runner)

See `.github/workflows/integration.yml` for configuration.

## Debugging Tests

### View Container Environment
```bash
docker run --rm -it \
  -v "$(pwd)/camp:/home/testuser/bin/camp:ro" \
  camp-integration-test:latest \
  /bin/bash
```

### Run Tests with Verbose Output
```bash
# Edit test script to add 'set -x' at the top for verbose mode
docker run --rm \
  -v "$(pwd)/camp:/home/testuser/bin/camp:ro" \
  -v "$(pwd)/test/integration/scripts:/home/testuser/tests:ro" \
  camp-integration-test:latest \
  /bin/bash -x /home/testuser/tests/test-bootstrap.sh
```

### Check Generated Files
```bash
# Run test and keep container alive
docker run --rm -it \
  -v "$(pwd)/camp:/home/testuser/bin/camp:ro" \
  -v "$(pwd)/test/integration/scripts:/home/testuser/tests:ro" \
  camp-integration-test:latest \
  /bin/bash

# Inside container:
testuser@container:~$ /home/testuser/bin/camp bootstrap
testuser@container:~$ cat ~/.camp/nix/flake.nix
testuser@container:~$ tree ~/.camp
```

## Adding New Tests

To add a new integration test:

1. Create test script in `scripts/`:
   ```bash
   touch test/integration/scripts/test-myfeature.sh
   chmod +x test/integration/scripts/test-myfeature.sh
   ```

2. Write test following the pattern:
   ```bash
   #!/bin/bash
   set -e
   source "$HOME/.nix-profile/etc/profile.d/nix.sh"

   echo "=== Testing: My Feature ==="

   # Setup
   rm -rf "$HOME/.camp"
   "$HOME/bin/camp" bootstrap

   # Test logic
   # ...

   # Assertions
   if [ condition ]; then
       echo "✓ Check passed"
   else
       echo "ERROR: Check failed"
       exit 1
   fi

   exit 0
   ```

3. Add test to `run-tests.sh`:
   ```bash
   run_test "My Feature" "test-myfeature.sh"
   ```

4. Test locally:
   ```bash
   make test-integration
   ```

## Performance

Typical test execution times:
- Docker image build: ~3-5 minutes (first time)
- Camp binary build: ~5-10 seconds
- Individual test: ~10-30 seconds
- Full test suite: ~2-3 minutes

**Optimizations**:
- Nix is pre-installed in Docker image (cached)
- Docker image layers are cached
- Tests run in parallel where possible
- Minimal package installation in tests

## Troubleshooting

### Tests fail with "Nix not found"
Ensure the Docker image was built correctly and Nix is installed:
```bash
cd test/integration/docker
docker build -t camp-integration-test:latest .
```

### Tests fail with "Permission denied"
Ensure test scripts are executable:
```bash
chmod +x test/integration/scripts/*.sh
chmod +x test/integration/docker/entrypoint.sh
```

### Tests fail with "camp: command not found"
Ensure camp binary is built before running tests:
```bash
go build -o camp main.go
```

### Docker build fails
Check Docker is running and you have sufficient disk space:
```bash
docker info
df -h
```

## Future Improvements

Potential enhancements:
- [ ] Add macOS container support (when available)
- [ ] Test actual Nix rebuilds (not just file generation)
- [ ] Add performance benchmarks
- [ ] Test concurrent operations
- [ ] Add stress tests (large configs)
- [ ] Test error recovery scenarios
- [ ] Add network-dependent flake tests
- [ ] Test migration from old versions

## References

- [Camp Documentation](../../CLAUDE.md)
- [Nix Documentation](https://nixos.org/manual/nix/stable/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
