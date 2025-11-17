package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewTemplateData(t *testing.T) {
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      "/Users/testuser",
		EnvVars: map[string]string{
			"EDITOR": "nvim",
		},
	}

	data := NewTemplateData(user)

	if data.Name != "testuser" {
		t.Errorf("Expected Name=testuser, got %s", data.Name)
	}

	if data.HostName != "testhost" {
		t.Errorf("Expected HostName=testhost, got %s", data.HostName)
	}

	if data.Platform != "darwin" {
		t.Errorf("Expected Platform=darwin, got %s", data.Platform)
	}

	if data.Architecture != "arm64" {
		t.Errorf("Expected Architecture=arm64, got %s", data.Architecture)
	}

	if data.HomeDir != "/Users/testuser" {
		t.Errorf("Expected HomeDir=/Users/testuser, got %s", data.HomeDir)
	}

	if data.EnvVars["EDITOR"] != "nvim" {
		t.Errorf("Expected EDITOR=nvim, got %s", data.EnvVars["EDITOR"])
	}
}

func TestCompileTemplate(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")

	// Write a simple template
	templateContent := `Hello {{.Name}}! Platform: {{.Platform}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create template data
	data := &TemplateData{
		Name:     "alice",
		Platform: "darwin",
	}

	// Compile template
	result, err := CompileTemplate(templatePath, data)
	if err != nil {
		t.Fatalf("CompileTemplate() failed: %v", err)
	}

	expected := "Hello alice! Platform: darwin"
	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, string(result))
	}
}

func TestCompileTemplate_WithEnvVars(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")

	// Write template with EnvVars iteration
	templateContent := `{{range $key, $value := .EnvVars}}{{$key}}={{$value}}
{{end}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create template data
	data := &TemplateData{
		EnvVars: map[string]string{
			"EDITOR":  "nvim",
			"BROWSER": "firefox",
		},
	}

	// Compile template
	result, err := CompileTemplate(templatePath, data)
	if err != nil {
		t.Fatalf("CompileTemplate() failed: %v", err)
	}

	// Check that both variables are present
	resultStr := string(result)
	if !strings.Contains(resultStr, "EDITOR=nvim") {
		t.Error("Expected EDITOR=nvim in result")
	}
	if !strings.Contains(resultStr, "BROWSER=firefox") {
		t.Error("Expected BROWSER=firefox in result")
	}
}

func TestCompileTemplate_NonExistentFile(t *testing.T) {
	data := &TemplateData{Name: "test"}

	_, err := CompileTemplate("/nonexistent/template.tmpl", data)
	if err == nil {
		t.Error("CompileTemplate() should error for non-existent file")
	}
}

func TestCompileTemplate_InvalidTemplate(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "invalid.tmpl")

	// Write invalid template syntax
	templateContent := `{{.Name`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	data := &TemplateData{Name: "test"}

	_, err := CompileTemplate(templatePath, data)
	if err == nil {
		t.Error("CompileTemplate() should error for invalid template syntax")
	}
}

func TestRenderFlakeTemplate(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create .camp directory
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	// Create config file
	configPath := filepath.Join(campDir, "camp.yml")
	configContent := `env:
  EDITOR: nvim
  BROWSER: firefox
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create user
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
	}

	// NOTE: We can't fully test RenderFlakeTemplate without the actual template file
	// This test would need the templates/files/flake.nix to exist
	// For now, we'll test that it properly errors when template doesn't exist
	err := RenderFlakeTemplate(user)
	if err == nil {
		// Template file might exist in the project, verify output
		outputPath := filepath.Join(tmpHome, ".camp", "nix", "flake.nix")
		if _, statErr := os.Stat(outputPath); statErr != nil {
			t.Error("RenderFlakeTemplate() should create output file")
		}
	} else {
		// Expected error if template doesn't exist
		if !strings.Contains(err.Error(), "failed to compile flake template") &&
			!strings.Contains(err.Error(), "failed to read template file") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}

func TestRenderFlakeTemplate_CreatesDirectory(t *testing.T) {
	// Skip if template file doesn't exist
	if _, err := os.Stat("templates/files/flake.nix"); os.IsNotExist(err) {
		t.Skip("Skipping test: flake.nix template not found")
	}

	// Create temporary home directory
	tmpHome := t.TempDir()

	// Create .camp directory and config
	campDir := filepath.Join(tmpHome, ".camp")
	if err := os.MkdirAll(campDir, 0755); err != nil {
		t.Fatalf("Failed to create .camp directory: %v", err)
	}

	configPath := filepath.Join(campDir, "camp.yml")
	configContent := `env:
  EDITOR: nvim
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create user
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      tmpHome,
		EnvVars:      make(map[string]string),
	}

	// Render template
	if err := RenderFlakeTemplate(user); err != nil {
		t.Fatalf("RenderFlakeTemplate() failed: %v", err)
	}

	// Verify output directory was created
	nixDir := filepath.Join(tmpHome, ".camp", "nix")
	if _, err := os.Stat(nixDir); os.IsNotExist(err) {
		t.Error("RenderFlakeTemplate() should create .camp/nix directory")
	}

	// Verify output file was created
	flakePath := filepath.Join(nixDir, "flake.nix")
	if _, err := os.Stat(flakePath); os.IsNotExist(err) {
		t.Error("RenderFlakeTemplate() should create flake.nix file")
	}

	// Verify content includes user data
	content, err := os.ReadFile(flakePath)
	if err != nil {
		t.Fatalf("Failed to read rendered flake.nix: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "testuser") {
		t.Error("Rendered template should contain username")
	}
	if !strings.Contains(contentStr, "testhost") {
		t.Error("Rendered template should contain hostname")
	}
	if !strings.Contains(contentStr, "EDITOR") {
		t.Error("Rendered template should contain env vars")
	}
}

// Flake template tests

func TestNewTemplateData_WithFlakes(t *testing.T) {
	user := &User{
		Name:         "testuser",
		HostName:     "testhost",
		Platform:     "darwin",
		Architecture: "arm64",
		HomeDir:      "/Users/testuser",
		EnvVars:      make(map[string]string),
		Flakes: []Flake{
			{
				Name: "my-flake",
				URL:  "github:user/repo",
				Follows: map[string]string{
					"nixpkgs": "nixpkgs",
				},
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	data := NewTemplateData(user)

	if len(data.Flakes) != 1 {
		t.Fatalf("Expected 1 flake, got %d", len(data.Flakes))
	}

	if data.Flakes[0].Name != "my-flake" {
		t.Errorf("Expected flake name=my-flake, got %s", data.Flakes[0].Name)
	}

	if data.Flakes[0].URL != "github:user/repo" {
		t.Errorf("Expected flake URL=github:user/repo, got %s", data.Flakes[0].URL)
	}
}

func TestCompileTemplate_WithFlakes(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")

	// Write template with flakes iteration
	templateContent := `{{range .Flakes}}flake: {{.Name}} url: {{.URL}}
{{end}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create template data with flakes
	data := &TemplateData{
		Flakes: []Flake{
			{Name: "flake1", URL: "github:user/flake1"},
			{Name: "flake2", URL: "github:user/flake2"},
		},
	}

	// Compile template
	result, err := CompileTemplate(templatePath, data)
	if err != nil {
		t.Fatalf("CompileTemplate() failed: %v", err)
	}

	// Check that both flakes are present
	resultStr := string(result)
	if !strings.Contains(resultStr, "flake: flake1") {
		t.Error("Expected flake1 in result")
	}
	if !strings.Contains(resultStr, "url: github:user/flake1") {
		t.Error("Expected flake1 URL in result")
	}
	if !strings.Contains(resultStr, "flake: flake2") {
		t.Error("Expected flake2 in result")
	}
}

func TestCompileTemplate_WithFlakeOutputs(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")

	// Write template with flake outputs iteration
	templateContent := `{{range $flake := .Flakes}}{{range .Outputs}}{{$flake.Name}}.{{.Name}} type={{.Type}}
{{end}}{{end}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create template data with flakes and outputs
	data := &TemplateData{
		Flakes: []Flake{
			{
				Name: "my-flake",
				URL:  "github:user/repo",
				Outputs: []FlakeOutput{
					{Name: "packages", Type: OutputTypeHome},
					{Name: "homeManagerModules.default", Type: OutputTypeHome},
				},
			},
		},
	}

	// Compile template
	result, err := CompileTemplate(templatePath, data)
	if err != nil {
		t.Fatalf("CompileTemplate() failed: %v", err)
	}

	// Check that outputs are present
	resultStr := string(result)
	if !strings.Contains(resultStr, "my-flake.packages") {
		t.Error("Expected my-flake.packages in result")
	}
	if !strings.Contains(resultStr, "type=home") {
		t.Error("Expected type=home in result")
	}
	if !strings.Contains(resultStr, "my-flake.homeManagerModules.default") {
		t.Error("Expected my-flake.homeManagerModules.default in result")
	}
}

func TestCompileTemplate_WithFlakeFollows(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")

	// Write template with follows iteration
	templateContent := `{{range .Flakes}}{{range $key, $value := .Follows}}{{$key}}={{$value}}
{{end}}{{end}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create template data with follows
	data := &TemplateData{
		Flakes: []Flake{
			{
				Name: "my-flake",
				URL:  "github:user/repo",
				Follows: map[string]string{
					"nixpkgs":      "nixpkgs",
					"home-manager": "home-manager",
				},
			},
		},
	}

	// Compile template
	result, err := CompileTemplate(templatePath, data)
	if err != nil {
		t.Fatalf("CompileTemplate() failed: %v", err)
	}

	// Check that follows are present
	resultStr := string(result)
	if !strings.Contains(resultStr, "nixpkgs=nixpkgs") {
		t.Error("Expected nixpkgs=nixpkgs in result")
	}
	if !strings.Contains(resultStr, "home-manager=home-manager") {
		t.Error("Expected home-manager=home-manager in result")
	}
}

func TestCompileTemplate_WithSystemAndHomeOutputs(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")

	// Write template that filters by output type
	templateContent := `System outputs:
{{range $flake := .Flakes}}{{range .Outputs}}{{if eq .Type "system"}}{{$flake.Name}}.{{.Name}}
{{end}}{{end}}{{end}}
Home outputs:
{{range $flake := .Flakes}}{{range .Outputs}}{{if eq .Type "home"}}{{$flake.Name}}.{{.Name}}
{{end}}{{end}}{{end}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create template data with mixed output types
	data := &TemplateData{
		Flakes: []Flake{
			{
				Name: "flake1",
				URL:  "github:user/flake1",
				Outputs: []FlakeOutput{
					{Name: "darwinModules.team", Type: OutputTypeSystem},
					{Name: "packages", Type: OutputTypeHome},
				},
			},
		},
	}

	// Compile template
	result, err := CompileTemplate(templatePath, data)
	if err != nil {
		t.Fatalf("CompileTemplate() failed: %v", err)
	}

	// Check that outputs are in correct sections
	resultStr := string(result)

	// System outputs section should have darwinModules.team
	systemSection := strings.Split(resultStr, "Home outputs:")[0]
	if !strings.Contains(systemSection, "flake1.darwinModules.team") {
		t.Error("Expected flake1.darwinModules.team in system outputs section")
	}

	// Home outputs section should have packages
	homeSection := strings.Split(resultStr, "Home outputs:")[1]
	if !strings.Contains(homeSection, "flake1.packages") {
		t.Error("Expected flake1.packages in home outputs section")
	}
}
