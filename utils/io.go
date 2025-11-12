package utils

import "strings"

// ReplaceInContent replaces all occurrences of old with new in the content
func ReplaceInContent(content []byte, old, new string) []byte {
	c := string(content)
	c = strings.Replace(c, old, new, -1)
	return []byte(c)
}
