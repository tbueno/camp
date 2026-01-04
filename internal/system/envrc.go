package system

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

//go:embed templates/.envrc.tmpl
var envrcTemplateFS embed.FS

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

	// Read embedded template
	templateContent, err := envrcTemplateFS.ReadFile("templates/.envrc.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to read embedded template: %w", err)
	}

	// Create template with custom functions
	funcMap := template.FuncMap{
		"escapeEnvValue": escapeEnvValue,
	}

	tmpl, err := template.New(".envrc").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
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
