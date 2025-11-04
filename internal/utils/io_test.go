package utils

import "testing"

func TestReplaceInFile(t *testing.T) {
	t.Run("replace in file", func(t *testing.T) {
		content := []byte("something something to-be-replaced")
		old := "to-be-replaced"
		new := "new-content"
		result := ReplaceInContent(content, old, new)
		expected := "something something new-content"
		if string(result) != expected {
			t.Errorf("expected %q, got %q", expected, string(result))
		}
	})
}

func TestReplaceInContent(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		old      string
		new      string
		expected string
	}{
		{
			name:     "replace occurrence",
			content:  []byte("something something to-be-replaced"),
			old:      "to-be-replaced",
			new:      "new-content",
			expected: "something something new-content",
		},
		{
			name:     "no occurrences",
			content:  []byte("something something"),
			old:      "to-be-replaced",
			new:      "new-content",
			expected: "something something",
		},
		{
			name:     "empty content",
			content:  []byte(""),
			old:      "to-be-replaced",
			new:      "new-content",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceInContent(tt.content, tt.old, tt.new)
			if string(result) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}
