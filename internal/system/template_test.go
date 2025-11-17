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

// renderNixValue tests

func TestRenderNixValue_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello", `"hello"`},
		{"empty string", "", `""`},
		{"string with spaces", "hello world", `"hello world"`},
		{"string with quotes", `hello "world"`, `"hello \"world\""`},
		{"string with backslash", `hello\world`, `"hello\\world"`},
		{"string with newline", "hello\nworld", `"hello\nworld"`},
		{"complex escaping", `"test\n"`, `"\"test\\n\""`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderNixValue(tt.input)
			if result != tt.expected {
				t.Errorf("renderNixValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRenderNixValue_Bool(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected string
	}{
		{"true", true, "true"},
		{"false", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderNixValue(tt.input)
			if result != tt.expected {
				t.Errorf("renderNixValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRenderNixValue_Int(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"zero", 0, "0"},
		{"positive int", 42, "42"},
		{"negative int", -10, "-10"},
		{"large int", 1000000, "1000000"},
		{"int64", int64(9223372036854775807), "9223372036854775807"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderNixValue(tt.input)
			if result != tt.expected {
				t.Errorf("renderNixValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRenderNixValue_Float(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"zero", 0.0, "0"},
		{"positive float", 3.14, "3.14"},
		{"negative float", -2.5, "-2.5"},
		{"large float", 1000.5, "1000.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderNixValue(tt.input)
			if result != tt.expected {
				t.Errorf("renderNixValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRenderNixValue_List(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected string
	}{
		{"empty list", []interface{}{}, "[  ]"},
		{"string list", []interface{}{"vim", "git", "tmux"}, `[ "vim" "git" "tmux" ]`},
		{"int list", []interface{}{1, 2, 3}, "[ 1 2 3 ]"},
		{"bool list", []interface{}{true, false, true}, "[ true false true ]"},
		{"mixed list", []interface{}{"hello", 42, true}, `[ "hello" 42 true ]`},
		{"float list", []interface{}{1.5, 2.7, 3.14}, "[ 1.5 2.7 3.14 ]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderNixValue(tt.input)
			if result != tt.expected {
				t.Errorf("renderNixValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCompileTemplate_WithRenderNixValue(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")

	// Write template that uses renderNixValue function
	templateContent := `{{range $key, $value := .Flakes}}{{range $argKey, $argValue := .Args}}{{ $argKey }} = {{ renderNixValue $argValue }};
{{end}}{{end}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create template data with flake args
	data := &TemplateData{
		Flakes: []Flake{
			{
				Name: "test-flake",
				URL:  "github:user/flake",
				Args: map[string]interface{}{
					"email":         "test@example.com",
					"enableFeature": true,
					"fontSize":      14,
					"threshold":     3.14,
					"packages":      []interface{}{"vim", "git"},
				},
			},
		},
	}

	// Compile template
	result, err := CompileTemplate(templatePath, data)
	if err != nil {
		t.Fatalf("CompileTemplate() failed: %v", err)
	}

	resultStr := string(result)

	// Check that all args are rendered correctly
	expectedLines := []string{
		`email = "test@example.com";`,
		`enableFeature = true;`,
		`fontSize = 14;`,
		`threshold = 3.14;`,
		`packages = [ "vim" "git" ];`,
	}

	for _, expected := range expectedLines {
		if !strings.Contains(resultStr, expected) {
			t.Errorf("Expected result to contain %q, got:\n%s", expected, resultStr)
		}
	}
}

func TestCompileTemplate_FlakeWithArgs(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "flake.nix")

	// Write simplified flake template with args
	templateContent := `{
  outputs = { ... }:
  {
    darwinConfigurations.test = {
      modules = [
        # System modules
        {{- range $flake := .Flakes }}
          {{- range .Outputs }}
            {{- if eq .Type "system" }}
        ({{ $flake.Name }}.{{ .Name }} {
          userName = "{{ $.Name }}";
          hostName = "{{ $.HostName }}";
          home = "{{ $.HomeDir }}";
          {{- range $key, $value := $flake.Args }}
          {{ $key }} = {{ renderNixValue $value }};
          {{- end }}
        })
            {{- end }}
          {{- end }}
        {{- end }}
      ];
    };

    homeConfigurations.test = {
      modules = [
        # Home modules
        {{- range $flake := .Flakes }}
          {{- range .Outputs }}
            {{- if eq .Type "home" }}
        ({{ $flake.Name }}.{{ .Name }} {
          userName = "{{ $.Name }}";
          hostName = "{{ $.HostName }}";
          home = "{{ $.HomeDir }}";
          {{- range $key, $value := $flake.Args }}
          {{ $key }} = {{ renderNixValue $value }};
          {{- end }}
        })
            {{- end }}
          {{- end }}
        {{- end }}
      ];
    };
  };
}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create template data with flake args
	data := &TemplateData{
		Name:     "testuser",
		HostName: "testhost",
		HomeDir:  "/home/testuser",
		Flakes: []Flake{
			{
				Name: "my-config",
				URL:  "github:user/config",
				Args: map[string]interface{}{
					"email":         "test@example.com",
					"enableFeature": true,
					"fontSize":      14,
					"packages":      []interface{}{"vim", "git"},
				},
				Outputs: []FlakeOutput{
					{Name: "darwinModules.default", Type: OutputTypeSystem},
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

	resultStr := string(result)

	// Verify system module has args
	if !strings.Contains(resultStr, `(my-config.darwinModules.default {`) {
		t.Error("Expected system module to be called as function")
	}
	if !strings.Contains(resultStr, `userName = "testuser";`) {
		t.Error("Expected userName arg in system module")
	}
	if !strings.Contains(resultStr, `hostName = "testhost";`) {
		t.Error("Expected hostName arg in system module")
	}
	if !strings.Contains(resultStr, `home = "/home/testuser";`) {
		t.Error("Expected home arg in system module")
	}
	if !strings.Contains(resultStr, `email = "test@example.com";`) {
		t.Error("Expected custom email arg in system module")
	}
	if !strings.Contains(resultStr, `enableFeature = true;`) {
		t.Error("Expected custom enableFeature arg in system module")
	}
	if !strings.Contains(resultStr, `fontSize = 14;`) {
		t.Error("Expected custom fontSize arg in system module")
	}
	if !strings.Contains(resultStr, `packages = [ "vim" "git" ];`) {
		t.Error("Expected custom packages arg in system module")
	}

	// Verify home module has args too
	if !strings.Contains(resultStr, `(my-config.homeManagerModules.default {`) {
		t.Error("Expected home module to be called as function")
	}

	// Count occurrences to ensure both modules got the args
	if strings.Count(resultStr, `userName = "testuser";`) != 2 {
		t.Errorf("Expected userName to appear twice (system + home), got %d occurrences", strings.Count(resultStr, `userName = "testuser";`))
	}
}
