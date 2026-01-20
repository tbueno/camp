---
title: "camp project activate"
linkTitle: "project activate"
weight: 10
description: >
  Export project environment variables to your shell
---

The `project activate` command reads environment variables from a project-level `camp.yaml` file and outputs shell-compatible export statements that can be evaluated to set those variables in your current shell session.

## Usage

```bash
eval "$(camp project activate)"
```

## How It Works

1. Reads the `camp.yaml` (or `camp.yml`) file in the current directory
2. Outputs shell export statements for each environment variable
3. Tracks which variables were exported
4. On subsequent runs, outputs `unset` statements for variables that were removed from the config

## Project Configuration

Create a `camp.yaml` file in your project root:

```yaml
env:
  DATABASE_URL: "postgres://localhost:5432/mydb"
  API_KEY: "your-api-key"
  DEBUG: "true"
  APP_ENV: "development"
```

## Output Examples

### Bash/Zsh Output

```bash
export API_KEY="your-api-key";
export APP_ENV="development";
export DATABASE_URL="postgres://localhost:5432/mydb";
export DEBUG="true";
```

### Fish Output

```fish
set -gx API_KEY 'your-api-key';
set -gx APP_ENV 'development';
set -gx DATABASE_URL 'postgres://localhost:5432/mydb';
set -gx DEBUG 'true';
```

## Variable Tracking

Camp tracks which variables it exports in `.camp/env_vars` within your project directory. This enables automatic cleanup when variables are removed:

```bash
# Initial config has VAR1, VAR2, VAR3
eval "$(camp project activate)"
# Output: export VAR1=...; export VAR2=...; export VAR3=...;

# After removing VAR2 from camp.yaml
eval "$(camp project activate)"
# Output: unset VAR2; export VAR1=...; export VAR3=...;

# After deleting camp.yaml entirely
eval "$(camp project activate)"
# Output: unset VAR1; unset VAR3;
```

## Flags

| Flag | Description |
|------|-------------|
| `--shell <type>` | Override shell detection. Options: `bash`, `zsh`, `fish` |

### Shell Override Example

```bash
# Force fish shell syntax
camp project activate --shell fish
```

## Shell Integration

### Bash/Zsh

Add to your `.bashrc` or `.zshrc`:

```bash
# Auto-activate camp project env when entering a directory
camp_activate() {
    if [ -f "camp.yaml" ] || [ -f "camp.yml" ]; then
        eval "$(camp project activate)"
    fi
}

# Option 1: Manual activation
alias activate='camp_activate'

# Option 2: Automatic on directory change
cd() {
    builtin cd "$@" && camp_activate
}
```

### Fish

Add to your `~/.config/fish/config.fish`:

```fish
function camp_activate --on-variable PWD
    if test -f camp.yaml; or test -f camp.yml
        eval (camp project activate)
    end
end
```

### With direnv

You can also use Camp with direnv by adding to your `.envrc`:

```bash
eval "$(camp project activate)"
```

## Special Character Handling

Values with special characters are properly escaped:

| Character | Handling |
|-----------|----------|
| `"` | Escaped as `\"` |
| `$` | Escaped as `\$` (prevents expansion) |
| `` ` `` | Escaped as `` \` `` |
| `\` | Escaped as `\\` |

```yaml
env:
  # This value contains special characters
  CONNECTION: 'postgresql://user:p@ss$word@host/db'
```

## Use Cases

- **Per-project database URLs**: Different databases for different projects
- **API keys and secrets**: Project-specific credentials (use with caution)
- **Feature flags**: Enable debug mode or specific features per project
- **Tool configuration**: Set `EDITOR`, `PAGER`, or other tool preferences

## Best Practices

1. **Don't commit secrets**: Add sensitive values to `.gitignore` or use a secrets manager
2. **Use `.camp/` in `.gitignore`**: The tracking directory should not be committed
3. **Document required variables**: Add a `camp.yaml.example` with placeholder values

## Related Commands

- [`camp project info`](#) - Display project information
- [`camp env`](../env/) - Display system environment information
