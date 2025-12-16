# Setting Up Camp Documentation

This guide explains how to set up the Camp documentation system using Hugo and the Docsy theme.

## Initial Setup (For Repository Maintainers)

### Step 1: Add Docsy Theme as Submodule

The Docsy theme needs to be added as a git submodule:

```bash
# From the repository root
cd docs

# Add Docsy theme as a submodule
git submodule add https://github.com/google/docsy.git themes/docsy

# Initialize and update the submodule
cd themes/docsy
git submodule update --init --recursive
cd ../../..

# Commit the submodule
git add .gitmodules docs/themes/docsy
git commit -m "Add Docsy theme as submodule"
```

### Step 2: Verify Configuration

The `docs/config.toml` file should already be configured to use the Docsy theme:

```toml
theme = "docsy"
```

### Step 3: Test Locally

```bash
# Install Hugo Extended (required for Docsy)
brew install hugo  # macOS
# or download from https://github.com/gohugoio/hugo/releases

# Verify Hugo is installed with extended support
hugo version
# Should show: hugo v0.1XX.X+extended

# Run the development server
cd docs
hugo server

# Visit http://localhost:1313/camp/
```

## For Contributors (Cloning the Repository)

### Option 1: Clone with Submodules

```bash
# Clone the repository with all submodules
git clone --recurse-submodules https://github.com/tbueno/camp.git
cd camp/docs
hugo server
```

### Option 2: Clone Then Initialize Submodules

```bash
# Clone normally
git clone https://github.com/tbueno/camp.git
cd camp

# Initialize submodules
git submodule update --init --recursive

# Run Hugo
cd docs
hugo server
```

## Troubleshooting

### Theme Not Found Error

```
Error: module "docsy" not found
```

**Solution**: Initialize the submodule:

```bash
cd docs
git submodule update --init --recursive
```

### Hugo Version Issues

```
Error: this feature requires the extended version of Hugo
```

**Solution**: Install Hugo Extended:

```bash
# macOS
brew uninstall hugo
brew install hugo

# Or download extended version from GitHub releases
```

### Submodule Update Issues

If the Docsy theme is out of date:

```bash
cd docs/themes/docsy
git pull origin main
cd ../../..
git add docs/themes/docsy
git commit -m "Update Docsy theme"
```

## Directory Structure After Setup

```
camp/
├── docs/
│   ├── content/en/           # Documentation content
│   ├── themes/
│   │   └── docsy/            # Docsy theme (git submodule)
│   ├── config.toml           # Hugo configuration
│   ├── README.md             # Documentation overview
│   ├── SETUP.md              # This file
│   └── DOCS_GUIDELINES.md    # Writing guidelines
├── .github/
│   └── workflows/
│       ├── deploy-docs.yml   # Auto-deploy to GitHub Pages
│       └── docs-pr-check.yml # PR validation
└── .gitmodules               # Git submodule configuration
```

## GitHub Pages Configuration

### Enable GitHub Pages

1. Go to repository Settings
2. Navigate to "Pages" section
3. Set source to "GitHub Actions"
4. The deploy workflow will handle the rest

### Custom Domain (Optional)

If you want to use a custom domain:

1. Add CNAME file to `docs/static/CNAME`:
   ```
   docs.yourcamp.com
   ```

2. Configure DNS:
   ```
   CNAME docs -> tbueno.github.io
   ```

3. Update baseURL in `docs/config.toml`:
   ```toml
   baseURL = "https://docs.yourcamp.com/"
   ```

## Updating the Theme

To update the Docsy theme to the latest version:

```bash
cd docs/themes/docsy
git checkout main
git pull origin main
cd ../../..
git add docs/themes/docsy
git commit -m "Update Docsy theme to latest version"
git push
```

## Development Workflow

### Working on Documentation

1. **Create/edit content**:
   ```bash
   cd docs/content/en/docs
   # Edit files or create new ones
   ```

2. **Test locally**:
   ```bash
   cd docs
   hugo server
   # Visit http://localhost:1313/camp/
   ```

3. **Commit and push**:
   ```bash
   git add docs/
   git commit -m "Update documentation for feature X"
   git push
   ```

4. **Automatic deployment**:
   - GitHub Actions will automatically deploy to GitHub Pages when merged to main

### Adding a New Section

1. **Create directory and index**:
   ```bash
   mkdir -p docs/content/en/docs/new-section
   touch docs/content/en/docs/new-section/_index.md
   ```

2. **Add front matter**:
   ```markdown
   ---
   title: "Section Name"
   linkTitle: "Section Name"
   weight: 10
   description: >
     Section description
   ---
   ```

3. **Add content pages**:
   ```bash
   touch docs/content/en/docs/new-section/page1.md
   ```

## CI/CD Pipeline

### Deployment Workflow

File: `.github/workflows/deploy-docs.yml`

- **Triggers**: Push to main, manual dispatch
- **Actions**:
  - Checkout code with submodules
  - Setup Hugo Extended
  - Build site
  - Deploy to GitHub Pages

### PR Check Workflow

File: `.github/workflows/docs-pr-check.yml`

- **Triggers**: Pull requests modifying `docs/`
- **Checks**:
  - Hugo build validation
  - Markdown linting
  - Spell checking
  - Broken link detection

## Version Management

### Creating a Documentation Version

When releasing a new version (e.g., v1.0.0):

1. **Update version menu in config.toml**:
   ```toml
   [[params.versions.list]]
     version = "v1.0"
     url = "https://tbueno.github.io/camp/v1.0/"
   ```

2. **Tag the documentation**:
   ```bash
   git tag -a v1.0.0-docs -m "Documentation for v1.0.0"
   git push origin v1.0.0-docs
   ```

3. **(Optional) Create versioned build**:
   - Create a new branch: `docs-v1.0`
   - Update baseURL for versioned path
   - Deploy to versioned subdirectory

## Resources

- [Hugo Documentation](https://gohugo.io/documentation/)
- [Docsy Theme Documentation](https://www.docsy.dev/docs/)
- [Hugo Installation](https://gohugo.io/installation/)
- [Git Submodules Documentation](https://git-scm.com/book/en/v2/Git-Tools-Submodules)

## Getting Help

- **Documentation issues**: See [DOCS_GUIDELINES.md](DOCS_GUIDELINES.md)
- **Hugo questions**: Check [Hugo Documentation](https://gohugo.io/documentation/)
- **Docsy questions**: See [Docsy Documentation](https://www.docsy.dev/)
- **Repository questions**: Open an issue on GitHub
