package system

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// GetExportedVars parses input and returns environment variables
func GetExportedVars(reader io.Reader) ([]EnvVar, error) {
	var envVars []EnvVar
	scanner := bufio.NewScanner(reader)

	// searches for export commands in the .envrc file
	exportRegex := regexp.MustCompile(`(?i)^\s*export\s+([A-Za-z_][A-Za-z0-9_]*)=(.*)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := exportRegex.FindStringSubmatch(line)
		if matches != nil {
			name := matches[1]
			value := strings.Trim(matches[2], `"'`) // Remove surrounding quotes

			envVars = append(envVars, EnvVar{
				Name:  name,
				Value: value,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVars, nil
}
