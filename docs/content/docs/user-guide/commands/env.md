---
title: "camp env"
linkTitle: "env"
weight: 1
description: >
  Display system and environment information
---

The `env` command displays detailed information about your system and Camp configuration.

## Usage

```bash
camp env
```

## Output

The command displays:

- **Username**: Your system username
- **Home Directory**: Path to your home directory
- **Platform**: Operating system (darwin/linux)
- **Architecture**: CPU architecture (amd64/arm64)
- **Hostname**: Your machine's hostname
- **Shell**: Your default shell
- **Environment Variables**: Custom variables from your `camp.yml`
- **Packages**: Nix packages configured for installation
- **Flakes**: External flakes you've configured

## Example Output

```
User Information:
  Name: bueno
  Home Directory: /Users/bueno
  Platform: darwin
  Architecture: arm64
  Hostname: macbook
  Shell: /bin/zsh

Environment Variables:
  EDITOR: nvim
  BROWSER: firefox

Packages:
  git
  neovim
  ripgrep
```

## Use Cases

- Verify your Camp configuration before rebuilding
- Check system information for debugging
- Confirm environment variables are set correctly
- Review configured packages and flakes

## Related Commands

- [`camp env rebuild`](../rebuild/) - Apply your configuration
- [`camp bootstrap`](../bootstrap/) - Initial setup
