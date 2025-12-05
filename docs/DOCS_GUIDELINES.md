# Documentation Guidelines

This document provides guidelines for maintaining and updating Camp's documentation.

## Documentation Structure

```
docs/
├── content/en/
│   ├── _index.html           # Homepage
│   ├── docs/                  # Main documentation
│   │   ├── getting-started/   # New user guides
│   │   ├── user-guide/        # Feature documentation
│   │   ├── developer-guide/   # Contributor guides
│   │   └── reference/         # API/CLI reference
│   └── blog/                  # Release notes, tutorials
├── config.toml               # Hugo configuration
└── themes/docsy/             # Docsy theme (git submodule)
```

## When to Update Documentation

### For New Features

When adding a new feature:

1. **User-facing docs** (always required):
   - Add to appropriate User Guide section
   - Update Getting Started if it affects initial setup
   - Add CLI reference if it's a new command
   - Include usage examples

2. **Developer docs** (if applicable):
   - Update Architecture docs for design changes
   - Add to Developer Guide for new patterns
   - Update Contributing guide if workflow changes

3. **Configuration schema** (if applicable):
   - Update Configuration Schema reference
   - Add examples to configuration guide

### For Bug Fixes

- Update relevant docs if behavior changes
- Add troubleshooting entries if it's a common issue
- Update examples if they were incorrect

### For Breaking Changes

- Update migration guides
- Add warnings to affected sections
- Update version compatibility info
- Document in changelog

## Writing Style

### General Principles

- **Clear and concise**: Use simple language
- **User-focused**: Write from the user's perspective
- **Actionable**: Provide specific steps and examples
- **Consistent**: Follow existing patterns

### Formatting

#### Headings

```markdown
# Page Title (H1 - only one per page)

## Major Section (H2)

### Subsection (H3)

#### Detail Section (H4 - use sparingly)
```

#### Code Blocks

Always specify the language:

````markdown
```bash
camp env rebuild
```

```yaml
packages:
  - git
  - neovim
```

```go
func main() {
    fmt.Println("Hello")
}
```
````

#### Admonitions

Use Hugo/Docsy shortcodes:

```markdown
{{% alert title="Note" %}}
This is important information.
{{% /alert %}}

{{% alert title="Warning" color="warning" %}}
Be careful with this operation.
{{% /alert %}}

{{% alert title="Tip" color="info" %}}
Here's a helpful tip.
{{% /alert %}}
```

#### Links

```markdown
<!-- Internal links (preferred) -->
[Configuration Guide](../configuration/)
[Getting Started](/docs/getting-started/)

<!-- External links -->
[Nix Documentation](https://nixos.org/manual/nix)
```

### Command Documentation

When documenting CLI commands:

```markdown
## Command Name

Brief description of what the command does.

### Usage

```bash
camp command [flags]
```

### Description

Detailed explanation of the command's purpose and behavior.

### Examples

```bash
# Example 1: Basic usage
camp command

# Example 2: With flags
camp command --flag value
```

### Options

- `--flag`: Description of what this flag does

### Related Commands

- [`camp other`](../other/) - Related command
```

### Configuration Documentation

For configuration options:

```markdown
## Option Name

**Type**: `string` | `boolean` | `number` | `array`
**Required**: Yes/No
**Default**: `default-value`

Description of what this option does.

### Example

```yaml
option: value
```

### Notes

- Additional information
- Edge cases
- Validation rules
```

## File Organization

### File Naming

- Use lowercase with hyphens: `getting-started.md`
- Use descriptive names: `package-management.md` not `packages.md`
- Index files: `_index.md` for section landing pages

### Front Matter

Every markdown file needs front matter:

```markdown
---
title: "Page Title"
linkTitle: "Short Title"  # Used in navigation
weight: 1                  # Order in navigation (lower = first)
description: >
  Brief description for SEO and page listing
---
```

For blog posts:

```markdown
---
title: "Post Title"
date: 2025-01-15
description: "Post description"
author: "Author Name"
---
```

## Testing Documentation

### Before Submitting

1. **Build locally**:
   ```bash
   cd docs
   hugo server
   # Visit http://localhost:1313/camp/
   ```

2. **Check for**:
   - Broken links
   - Missing images
   - Formatting errors
   - Code examples work
   - All sections render correctly

3. **Validate**:
   - Markdown lint: `markdownlint content/**/*.md`
   - Spell check: Run cspell
   - Link check: Use lychee or Hugo's built-in checker

### CI Checks

Pull requests automatically run:
- Hugo build test
- Markdown linting
- Spell checking
- Link validation

Fix any CI failures before merging.

## Release Documentation Checklist

When preparing a release:

- [ ] Update CHANGELOG.md
- [ ] Create blog post for release notes
- [ ] Update version in config.toml
- [ ] Update any version-specific references
- [ ] Add migration guide for breaking changes
- [ ] Update Getting Started if installation changed
- [ ] Verify all examples work with new version
- [ ] Tag documentation with version number

## Version Management

### Version Tags

When releasing version X.Y.Z:

1. Create version in config.toml:
   ```toml
   [[params.versions.list]]
     version = "vX.Y"
     url = "https://tbueno.github.io/camp/vX.Y/"
   ```

2. Tag documentation:
   ```bash
   git tag -a vX.Y.Z-docs -m "Documentation for vX.Y.Z"
   ```

### Development Documentation

The `main` branch documentation is always for the development version:
- Mark unstable features clearly
- Note version availability
- Link to stable docs when needed

## Common Patterns

### Adding a New User Guide Section

1. Create directory: `docs/content/en/docs/user-guide/feature-name/`
2. Add `_index.md` for the section
3. Add individual pages as needed
4. Update parent `_index.md` to link to new section
5. Test locally
6. Submit PR with "docs" label

### Adding a New Command

1. Add to `docs/content/en/docs/user-guide/commands/command-name.md`
2. Update `docs/content/en/docs/reference/cli-reference.md`
3. Add examples and use cases
4. Link from related pages
5. Test all examples work

### Adding Examples

Examples should be:
- **Complete**: Can be copy-pasted and run
- **Tested**: Verify they actually work
- **Commented**: Explain non-obvious parts
- **Realistic**: Show real-world usage

## Documentation Review

### Self-Review Checklist

Before submitting docs PR:

- [ ] Spelling and grammar checked
- [ ] Code examples tested
- [ ] Links work (internal and external)
- [ ] Images load correctly
- [ ] Front matter is complete
- [ ] Follows style guidelines
- [ ] Hugo builds without errors
- [ ] Mobile-friendly (check responsive design)

### Peer Review

When reviewing docs PRs:

- Verify technical accuracy
- Check for clarity and readability
- Test code examples
- Suggest improvements to structure
- Look for missing context

## Getting Help

- Check existing documentation for patterns
- Ask in GitHub Discussions
- Reference [Docsy documentation](https://www.docsy.dev/)
- Review [Hugo documentation](https://gohugo.io/documentation/)

## Resources

- [Docsy Theme](https://www.docsy.dev/)
- [Hugo Documentation](https://gohugo.io/documentation/)
- [Markdown Guide](https://www.markdownguide.org/)
- [Technical Writing Style Guide](https://developers.google.com/style)
