package system

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// EnvrcTemplateData holds data for rendering .envrc template
type EnvrcTemplateData struct {
	HasDevbox bool
	EnvVars   map[string]string
}

// GenerateEnvrc generates a .envrc file for the project using a template
// Includes camp env vars and devbox integration
func GenerateEnvrc(projectPath string, config *CampConfig) (string, error) {
	// Check if devbox.json exists
	devboxPath := filepath.Join(projectPath, "devbox.json")
	hasDevbox := false
	if _, err := os.Stat(devboxPath); err == nil {
		hasDevbox = true
	}

	// Prepare template data
	var envVars map[string]string
	if config != nil && config.Env != nil {
		envVars = config.Env
	} else {
		envVars = make(map[string]string)
	}

	data := EnvrcTemplateData{
		HasDevbox: hasDevbox,
		EnvVars:   envVars,
	}

	// Register escapeEnvValue as a template function
	funcMap := make(map[string]interface{})
	funcMap["escapeEnvValue"] = escapeEnvValue

	// Find template file (try absolute path first, then relative)
	templatePath := "templates/files/.envrc.tmpl"

	// Render template
	result, err := CompileTemplateGeneric(templatePath, data, funcMap)
	if err != nil {
		return "", fmt.Errorf("failed to render .envrc template: %w", err)
	}

	return string(result), nil
}

// escapeEnvValue escapes special characters in environment variable values
// Handles quotes, backslashes, dollar signs, and other special bash characters
func escapeEnvValue(value string) string {
	// Replace backslashes first to avoid double-escaping
	value = strings.ReplaceAll(value, "\\", "\\\\")
	// Escape double quotes
	value = strings.ReplaceAll(value, "\"", "\\\"")
	// Escape dollar signs to prevent variable expansion
	value = strings.ReplaceAll(value, "$", "\\$")
	// Escape backticks
	value = strings.ReplaceAll(value, "`", "\\`")
	return value
}
