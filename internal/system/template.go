package system

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData holds the data to be injected into templates
type TemplateData struct {
	Name         string            // Username
	HostName     string            // Machine hostname
	Platform     string            // OS (darwin/linux)
	Architecture string            // CPU arch (amd64/arm64)
	HomeDir      string            // User's home directory
	EnvVars      map[string]string // Custom environment variables
	Packages     []string          // Nix packages to install
	Flakes       []Flake           // External Nix flakes to integrate
}

// NewTemplateData creates template data from a User
func NewTemplateData(user *User) *TemplateData {
	return &TemplateData{
		Name:         user.Name,
		HostName:     user.HostName,
		Platform:     user.Platform,
		Architecture: user.Architecture,
		HomeDir:      user.HomeDir,
		EnvVars:      user.EnvVars,
		Packages:     user.Packages,
		Flakes:       user.Flakes,
	}
}

// renderNixValue converts a Go value to its Nix syntax representation
func renderNixValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		// Escape special characters in strings
		escaped := strings.ReplaceAll(v, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		escaped = strings.ReplaceAll(escaped, "\n", "\\n")
		return fmt.Sprintf("\"%s\"", escaped)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case int, int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case []interface{}:
		// Render list elements
		var elements []string
		for _, elem := range v {
			elements = append(elements, renderNixValue(elem))
		}
		return fmt.Sprintf("[ %s ]", strings.Join(elements, " "))
	default:
		// Fallback for unsupported types (should be caught by validation)
		return fmt.Sprintf("\"%v\"", v)
	}
}

// CompileTemplate parses and renders a template file with the given data
func CompileTemplate(templatePath string, data *TemplateData) ([]byte, error) {
	// Read template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template with custom functions
	funcMap := template.FuncMap{
		"renderNixValue": renderNixValue,
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// CompileTemplateGeneric parses and renders a template file with any data type
// This is a more flexible version of CompileTemplate that accepts custom function maps
func CompileTemplateGeneric(templatePath string, data interface{}, customFuncs template.FuncMap) ([]byte, error) {
	// Read template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Merge default functions with custom functions
	funcMap := template.FuncMap{
		"renderNixValue": renderNixValue,
	}
	for name, fn := range customFuncs {
		funcMap[name] = fn
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// RenderFlakeTemplate renders the flake.nix template with user data
// and saves it to ~/.camp/nix/flake.nix
func RenderFlakeTemplate(user *User) error {
	// Reload user config to get latest env vars
	if err := user.Reload(); err != nil {
		return fmt.Errorf("failed to reload user config: %w", err)
	}

	// Create template data
	data := NewTemplateData(user)

	// Template path (from templates/files/)
	templatePath := "templates/files/flake.nix"

	// Compile template
	rendered, err := CompileTemplate(templatePath, data)
	if err != nil {
		return fmt.Errorf("failed to compile flake template: %w", err)
	}

	// Output path
	outputPath := filepath.Join(user.HomeDir, ".camp", "nix", "flake.nix")

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write rendered template
	if err := os.WriteFile(outputPath, rendered, 0644); err != nil {
		return fmt.Errorf("failed to write rendered template: %w", err)
	}

	return nil
}
