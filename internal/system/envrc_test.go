package system

import (
	"strings"
	"testing"
)

func TestGetExportedVars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []EnvVar
	}{
		{
			name:  "basic export",
			input: "export FOO=bar",
			expected: []EnvVar{
				{Name: "FOO", Value: "bar"},
			},
		},
		{
			name:  "multiple exports",
			input: "export FOO=bar\nexport BAZ=qux",
			expected: []EnvVar{
				{Name: "FOO", Value: "bar"},
				{Name: "BAZ", Value: "qux"},
			},
		},
		{
			name:  "export with double quotes",
			input: `export FOO="bar baz"`,
			expected: []EnvVar{
				{Name: "FOO", Value: "bar baz"},
			},
		},
		{
			name:  "export with single quotes",
			input: `export FOO='bar baz'`,
			expected: []EnvVar{
				{Name: "FOO", Value: "bar baz"},
			},
		},
		{
			name:  "case insensitive EXPORT",
			input: "EXPORT FOO=bar",
			expected: []EnvVar{
				{Name: "FOO", Value: "bar"},
			},
		},
		{
			name:  "mixed case export",
			input: "Export FOO=bar",
			expected: []EnvVar{
				{Name: "FOO", Value: "bar"},
			},
		},
		{
			name:  "with comments and empty lines",
			input: "# This is a comment\nexport FOO=bar\n\n# Another comment\nexport BAZ=qux",
			expected: []EnvVar{
				{Name: "FOO", Value: "bar"},
				{Name: "BAZ", Value: "qux"},
			},
		},
		{
			name:  "with leading/trailing spaces",
			input: "  export FOO=bar  \n\t export BAZ=qux \t",
			expected: []EnvVar{
				{Name: "FOO", Value: "bar"},
				{Name: "BAZ", Value: "qux"},
			},
		},
		{
			name:  "variable names with underscores and numbers",
			input: "export FOO_BAR=test\nexport VAR_123=value",
			expected: []EnvVar{
				{Name: "FOO_BAR", Value: "test"},
				{Name: "VAR_123", Value: "value"},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []EnvVar{},
		},
		{
			name:     "only comments",
			input:    "# Comment 1\n# Comment 2",
			expected: []EnvVar{},
		},
		{
			name:  "non-export lines ignored",
			input: "echo hello\nexport FOO=bar\nls -la",
			expected: []EnvVar{
				{Name: "FOO", Value: "bar"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := GetExportedVars(reader)

			if err != nil {
				t.Errorf("GetExportedVars() error = %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("GetExportedVars() got %d variables, want %d", len(result), len(tt.expected))
				return
			}

			for i, got := range result {
				expected := tt.expected[i]
				if got.Name != expected.Name || got.Value != expected.Value {
					t.Errorf("GetExportedVars()[%d] = {Name: %q, Value: %q}, want {Name: %q, Value: %q}",
						i, got.Name, got.Value, expected.Name, expected.Value)
				}
			}
		})
	}
}
