# Camp Documentation

This directory contains the source for Camp's documentation website, built with [Hugo](https://gohugo.io/) and the [Docsy theme](https://www.docsy.dev/).

## Quick Start

### Prerequisites

- **Hugo Extended** (v0.110.0 or later)
  ```bash
  # macOS
  brew install hugo

  # Or download from https://github.com/gohugoio/hugo/releases
  ```

- **Git** (for Docsy theme submodule)

### Initial Setup

1. **Clone with submodules** (if you haven't already):
   ```bash
   git clone --recurse-submodules https://github.com/tbueno/camp
   cd camp/docs
   ```

2. **Or initialize submodules** (if already cloned):
   ```bash
   git submodule update --init --recursive
   ```

### Local Development

Run the Hugo development server:

```bash
cd docs
hugo server
```

Then visit: http://localhost:1313/camp/

The server will auto-reload when you make changes.

### Building

Build the static site:

```bash
cd docs
hugo
```

Output will be in `docs/public/`.

## Documentation Structure

```
docs/
├── content/en/              # English documentation
│   ├── _index.html          # Homepage
│   ├── docs/                # Main documentation
│   │   ├── getting-started/ # Installation, quickstart
│   │   ├── user-guide/      # Feature guides
│   │   ├── developer-guide/ # Contributing, architecture
│   │   └── reference/       # CLI reference, schemas
│   └── blog/                # Release notes, tutorials
├── static/                  # Static assets (images, etc.)
├── layouts/                 # Custom Hugo layouts (if needed)
├── config.toml             # Hugo configuration
├── themes/docsy/           # Docsy theme (git submodule)
└── README.md               # This file
```

## Writing Documentation

### Creating a New Page

1. Create a markdown file in the appropriate section:
   ```bash
   # Example: new user guide page
   touch content/en/docs/user-guide/new-feature.md
   ```

2. Add front matter:
   ```markdown
   ---
   title: "Feature Name"
   linkTitle: "Feature Name"
   weight: 5
   description: >
     Brief description of the feature
   ---

   Content goes here...
   ```

3. Test locally with `hugo server`

### Front Matter Fields

- `title`: Full page title
- `linkTitle`: Short title for navigation (optional)
- `weight`: Order in navigation (lower numbers first)
- `description`: Brief description for SEO and listings

### Style Guidelines

See [DOCS_GUIDELINES.md](DOCS_GUIDELINES.md) for detailed guidelines on:
- Writing style
- Formatting conventions
- Code examples
- When to update docs

## Deployment

### Automatic Deployment

Documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch.

The workflow:
1. Push to `main` (or merge PR)
2. GitHub Actions builds the site
3. Deploys to `https://tbueno.github.io/camp/`

See [.github/workflows/deploy-docs.yml](../.github/workflows/deploy-docs.yml) for details.

### Manual Deployment

You can also manually trigger deployment:

1. Go to GitHub Actions
2. Select "Deploy Documentation" workflow
3. Click "Run workflow"

## CI Checks

Pull requests that modify documentation are automatically checked for:

- **Hugo build**: Ensures site builds without errors
- **Markdown linting**: Checks markdown formatting
- **Spell checking**: Catches typos
- **Link validation**: Finds broken links

Fix any CI failures before merging.

## Versioning

### Version Management

Documentation versions are managed in `config.toml`:

```toml
[[params.versions.list]]
  version = "v1.0"
  url = "https://tbueno.github.io/camp/v1.0/"
[[params.versions.list]]
  version = "main (development)"
  url = "https://tbueno.github.io/camp/"
```

### Creating a Version

When releasing a new version:

1. Update the version list in `config.toml`
2. Tag the documentation:
   ```bash
   git tag -a v1.0.0-docs -m "Documentation for v1.0.0"
   git push origin v1.0.0-docs
   ```

## Theme Customization

The Docsy theme is included as a git submodule. To customize:

### Updating the Theme

```bash
cd docs/themes/docsy
git pull origin main
cd ../..
git add themes/docsy
git commit -m "Update Docsy theme"
```

### Custom Layouts

Add custom layouts to `docs/layouts/` to override theme defaults.

### Custom Styles

Add custom CSS to `docs/assets/scss/_custom.scss` (create if needed).

## Troubleshooting

### Hugo build fails

```bash
# Check Hugo version (need extended version)
hugo version

# Should output something like:
# hugo v0.120.0+extended darwin/arm64
```

### Submodule issues

```bash
# Reset submodules
git submodule deinit -f .
git submodule update --init --recursive
```

### Server doesn't auto-reload

```bash
# Try with --disableFastRender
hugo server --disableFastRender
```

### Port already in use

```bash
# Use a different port
hugo server --port 1314
```

## Resources

- [Hugo Documentation](https://gohugo.io/documentation/)
- [Docsy Theme Documentation](https://www.docsy.dev/)
- [Markdown Guide](https://www.markdownguide.org/)
- [Camp Documentation Guidelines](DOCS_GUIDELINES.md)

## Getting Help

- Check [DOCS_GUIDELINES.md](DOCS_GUIDELINES.md) for writing guidelines
- Ask in [GitHub Discussions](https://github.com/tbueno/camp/discussions)
- Open an issue with the `documentation` label
