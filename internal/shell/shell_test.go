package shell

import (
	"os"
	"testing"
)

func TestDetectShell(t *testing.T) {
	tests := []struct {
		name     string
		shellEnv string
		expected ShellType
	}{
		{"bash shell", "/bin/bash", ShellBash},
		{"zsh shell", "/usr/bin/zsh", ShellZsh},
		{"fish shell", "/usr/local/bin/fish", ShellFish},
		{"bash with path", "/usr/local/bin/bash", ShellBash},
		{"unknown shell dash", "/bin/dash", ShellPosix},
		{"unknown shell sh", "/bin/sh", ShellPosix},
		{"empty shell", "", ShellPosix},
		{"uppercase ZSH", "/bin/ZSH", ShellZsh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := os.Getenv("SHELL")
			os.Setenv("SHELL", tt.shellEnv)
			defer os.Setenv("SHELL", original)

			result := DetectShell()
			if result != tt.expected {
				t.Errorf("DetectShell() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatExport_Bash(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		value    string
		expected string
	}{
		{
			"simple value",
			"FOO",
			"bar",
			`export FOO="bar";`,
		},
		{
			"value with spaces",
			"MSG",
			"hello world",
			`export MSG="hello world";`,
		},
		{
			"value with double quotes",
			"MSG",
			`hello "world"`,
			`export MSG="hello \"world\"";`,
		},
		{
			"value with dollar sign",
			"PATH_EXT",
			"$HOME/bin",
			`export PATH_EXT="\$HOME/bin";`,
		},
		{
			"value with backtick",
			"CMD",
			"echo `date`",
			"export CMD=\"echo \\`date\\`\";",
		},
		{
			"value with backslash",
			"ESCAPED",
			`path\to\file`,
			`export ESCAPED="path\\to\\file";`,
		},
		{
			"empty value",
			"EMPTY",
			"",
			`export EMPTY="";`,
		},
		{
			"complex value",
			"COMPLEX",
			`$HOME/"test"\dir`,
			`export COMPLEX="\$HOME/\"test\"\\dir";`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatExport(ShellBash, tt.varName, tt.value)
			if result != tt.expected {
				t.Errorf("FormatExport(bash) = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatExport_Zsh(t *testing.T) {
	// Zsh uses same format as bash
	tests := []struct {
		name     string
		varName  string
		value    string
		expected string
	}{
		{
			"simple value",
			"FOO",
			"bar",
			`export FOO="bar";`,
		},
		{
			"value with special chars",
			"MSG",
			`$HOME/"test"`,
			`export MSG="\$HOME/\"test\"";`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatExport(ShellZsh, tt.varName, tt.value)
			if result != tt.expected {
				t.Errorf("FormatExport(zsh) = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatExport_Fish(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		value    string
		expected string
	}{
		{
			"simple value",
			"FOO",
			"bar",
			`set -gx FOO 'bar';`,
		},
		{
			"value with single quote",
			"MSG",
			"hello 'world'",
			`set -gx MSG 'hello \'world\'';`,
		},
		{
			"value with backslash",
			"PATH",
			`path\to\file`,
			`set -gx PATH 'path\\to\\file';`,
		},
		{
			"value with dollar sign (no escaping needed in single quotes)",
			"VAR",
			"$HOME/bin",
			`set -gx VAR '$HOME/bin';`,
		},
		{
			"empty value",
			"EMPTY",
			"",
			`set -gx EMPTY '';`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatExport(ShellFish, tt.varName, tt.value)
			if result != tt.expected {
				t.Errorf("FormatExport(fish) = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatExport_Posix(t *testing.T) {
	// Posix uses same format as bash
	result := FormatExport(ShellPosix, "FOO", "bar")
	expected := `export FOO="bar";`
	if result != expected {
		t.Errorf("FormatExport(posix) = %q, want %q", result, expected)
	}
}

func TestFormatUnset(t *testing.T) {
	tests := []struct {
		name     string
		shell    ShellType
		varName  string
		expected string
	}{
		{"bash", ShellBash, "FOO", "unset FOO;"},
		{"zsh", ShellZsh, "BAR", "unset BAR;"},
		{"posix", ShellPosix, "BAZ", "unset BAZ;"},
		{"fish", ShellFish, "QUX", "set -e QUX;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUnset(tt.shell, tt.varName)
			if result != tt.expected {
				t.Errorf("FormatUnset(%s) = %q, want %q", tt.shell, result, tt.expected)
			}
		})
	}
}
