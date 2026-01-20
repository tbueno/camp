package shell

import (
	"os"
	"path/filepath"
	"strings"
)

// ShellType represents different shell types
type ShellType string

const (
	ShellBash  ShellType = "bash"
	ShellZsh   ShellType = "zsh"
	ShellFish  ShellType = "fish"
	ShellPosix ShellType = "posix" // fallback for unknown shells
)

// DetectShell determines the current shell type from the SHELL environment variable
func DetectShell() ShellType {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return ShellPosix
	}

	shellName := strings.ToLower(filepath.Base(shellPath))

	switch shellName {
	case "bash":
		return ShellBash
	case "zsh":
		return ShellZsh
	case "fish":
		return ShellFish
	default:
		return ShellPosix
	}
}

// FormatExport formats an environment variable export statement for the given shell
func FormatExport(shell ShellType, name, value string) string {
	escapedValue := escapeValue(shell, value)

	switch shell {
	case ShellFish:
		return "set -gx " + name + " " + escapedValue + ";"
	default:
		// Bash, Zsh, and POSIX all use export VAR="value"
		return "export " + name + "=" + escapedValue + ";"
	}
}

// FormatUnset formats an environment variable unset statement for the given shell
func FormatUnset(shell ShellType, name string) string {
	switch shell {
	case ShellFish:
		return "set -e " + name + ";"
	default:
		// Bash, Zsh, and POSIX all use unset VAR
		return "unset " + name + ";"
	}
}

// escapeValue escapes special characters in the value for the given shell
func escapeValue(shell ShellType, value string) string {
	switch shell {
	case ShellFish:
		// Fish uses single quotes with escaped single quotes
		escaped := strings.ReplaceAll(value, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "'", "\\'")
		return "'" + escaped + "'"
	default:
		// For bash/zsh/posix, use double quotes and escape special chars
		escaped := value
		escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		escaped = strings.ReplaceAll(escaped, "$", "\\$")
		escaped = strings.ReplaceAll(escaped, "`", "\\`")
		return "\"" + escaped + "\""
	}
}
