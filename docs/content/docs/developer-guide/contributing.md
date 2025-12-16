---
title: "Contributing"
linkTitle: "Contributing"
weight: 1
description: >
  How to contribute to Camp
---

Thank you for your interest in contributing to Camp! This guide will
help you get started.

## Ways to Contribute

- **Report bugs**: Open an issue describing the problem
- **Request features**: Suggest new features or improvements
- **Write code**: Submit pull requests for bug fixes or features
- **Improve documentation**: Help make our docs better
- **Help others**: Answer questions in discussions

## Getting Started

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/camp
cd camp
```

### 2. Set Up Development Environment

<!-- See the [Development Setup](../development-setup/) guide. -->

Requirements:

- Go 1.24.4 or later
- Nix package manager
- Git

### 3. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 4. Make Changes

- Follow the [code style](#code-style)
- Write or update tests
- Update documentation

### 5. Test Your Changes

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./cmd
go test ./internal/system

# Format code
go fmt ./...
```

### 6. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "Add feature: support for custom Nix channels"
```

**Good commit messages:**

- Start with a verb (Add, Fix, Update, Remove)
- Be specific about what changed
- Reference issues when applicable

```text
Add support for custom environment variables

Implements #42. Users can now define custom environment variables
in camp.yml that will be injected into their Nix environment.
```

### 7. Push and Create Pull Request

```bash
git push origin your-branch-name
```

Then create a pull request on GitHub.

## Pull Request Guidelines

### PR Requirements

- [ ] Tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] PR description explains the changes

### PR Template

```markdown
## Description
Brief description of what this PR does

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How has this been tested?

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] Code formatted
```

### Review Process

1. A maintainer will review your PR
2. Address any requested changes
3. Once approved, a maintainer will merge it
4. Your changes will be included in the next release

## Code Style

### Go Conventions

Follow standard Go conventions:

- Run `go fmt` after making changes
- Use meaningful variable names
- Write clear comments for exported functions
- Keep functions focused and concise

### Command Structure

All CLI commands should:

```go
// Use cmd.OutOrStdout() for output (supports testing)
fmt.Fprintln(cmd.OutOrStdout(), "Output here")

// Not: fmt.Println() - this breaks testing
```

### Package Organization

```text
camp/
├── cmd/                 # CLI commands (Cobra)
│   ├── root.go
│   ├── env.go
│   └── bootstrap.go
├── internal/            # Internal packages
│   ├── system/          # System info, config, templates
│   └── utils/           # Utilities
└── main.go              # Entry point (minimal)
```

- Keep command logic in `cmd/` package
- System logic in `internal/system/`
- Utilities in `internal/utils/`
- Main entry point should just call `cmd.Execute()`

## Testing Guidelines

### Test Coverage

- Write tests for all new features
- Maintain existing test coverage
- Test both success and error cases

### Test Organization

```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := setupInput()

    // Act
    result := FunctionName(input)

    // Assert
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Using Temporary Directories

```go
func TestWithConfig(t *testing.T) {
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, ".camp", "camp.yml")

    // Create test config
    // ...

    // Test runs, tmpDir is cleaned up automatically
}
```

### Running Tests

```bash
# All tests
go test ./...

# Verbose output
go test -v ./...

# Specific package
go test ./cmd

# With coverage
go test -cover ./...
```

## Documentation Updates

When adding features or changing behavior:

1. **Update user docs**: Add to appropriate docs section
2. **Update CLAUDE.md**: Update developer instructions
3. **Update README**: If it affects quick start or overview
4. **Add examples**: Include usage examples

### Documentation Style

- Use clear, simple language
- Include code examples
- Add troubleshooting tips
- Link to related docs

## Feature Development Workflow

### Planning

1. **Open an issue** to discuss the feature
2. **Get feedback** from maintainers
3. **Create implementation plan** (for large features)

### Implementation

1. **Create branch**: `feature/feature-name`
2. **Write tests first** (TDD approach encouraged)
3. **Implement feature**
4. **Update documentation**
5. **Submit PR**

### Example: Adding a New Command

```go
// 1. Create cmd/newcommand.go
package cmd

import "github.com/spf13/cobra"

var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}

// 2. Create cmd/newcommand_test.go
func TestNewCommand(t *testing.T) {
    // Test implementation
}

// 3. Update documentation in docs/
// 4. Update CLAUDE.md if needed
```

## Reporting Bugs

### Before Reporting

1. Search existing issues
2. Verify it's reproducible
3. Test on latest version

### Bug Report Template

```markdown
**Description**
Clear description of the bug

**To Reproduce**
1. Step one
2. Step two
3. Bug occurs

**Expected Behavior**
What should happen

**Environment**
- OS: macOS/Linux
- Camp version: vX.Y.Z
- Go version: 1.24.4

**Additional Context**
Logs, screenshots, etc.
```

## Questions?

- Check the [FAQ](/docs/faq/)
- Ask in [GitHub Discussions](https://github.com/tbueno/camp/discussions)
- Open an issue with the question label

## License

By contributing, you agree that your contributions will be licensed
under the Apache License 2.0.
