package system

import (
	"camp/internal/utils"
	"os"
	"os/user"
	"runtime"
)

// System represents the system's basic information
type System struct {
	OS           string
	Architecture string
}

// Home represents home directory information
// This struct is prepared for future expansion
type Home struct {
	Path string
	// Future fields: Owner, Permissions, etc.
}

// User represents user information
type User struct {
	Name         string
	HomeDir      string
	Platform     string
	Architecture string
	Shell        string
	HostName     string
	EnvVars      map[string]string // Custom environment variables from camp.yml
	Packages     []string          // Nix packages to install from camp.yml
	Flakes       []Flake           // External Nix flakes from camp.yml
}

// NewUser creates a new User instance using only current user's machine information.
func NewUser() *User {
	u, _ := user.Current()
	shell := os.Getenv("SHELL")
	user := &User{
		Name:         u.Username,
		HomeDir:      u.HomeDir,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		Shell:        shell,
		HostName:     utils.HostName(),
		EnvVars:      make(map[string]string),
		Packages:     []string{},
		Flakes:       []Flake{},
	}
	// Load config and populate EnvVars and Flakes if available
	user.Reload()
	return user
}

// Reload refreshes the user's configuration from camp.yml
// This loads environment variables and flakes from ~/.camp/camp.yml or ~/.camp/camp.yaml
func (u *User) Reload() error {
	config, err := LoadUserConfig(u.HomeDir)
	if err != nil {
		// If config loading fails, keep existing EnvVars and Flakes
		return err
	}

	// Update EnvVars from config
	if config.Env != nil {
		u.EnvVars = config.Env
	} else {
		u.EnvVars = make(map[string]string)
	}

	// Update Packages from config
	if config.Packages != nil {
		u.Packages = config.Packages
	} else {
		u.Packages = []string{}
	}

	// Update Flakes from config
	if config.Flakes != nil {
		u.Flakes = config.Flakes
	} else {
		u.Flakes = []Flake{}
	}

	return nil
}

// FlakeOutputType defines the allowed types for a flake's output
type FlakeOutputType string

const (
	// OutputTypeSystem indicates the output should be applied at system level (nix-darwin)
	OutputTypeSystem FlakeOutputType = "system"
	// OutputTypeHome indicates the output should be applied at user level (home-manager)
	OutputTypeHome FlakeOutputType = "home"
)

// FlakeOutput represents a specific output from a flake
type FlakeOutput struct {
	Name string          `yaml:"name"` // Output name (e.g., "packages", "homeManagerModules.default")
	Type FlakeOutputType `yaml:"type"` // Where to apply: "system" or "home"
}

// Flake represents an external Nix flake reference
type Flake struct {
	Name    string                 `yaml:"name"`    // Unique identifier for the flake
	URL     string                 `yaml:"url"`     // Flake URL (github:user/repo, git+ssh://..., path:/..., etc.)
	Follows map[string]string      `yaml:"follows"` // Input dependency overrides (e.g., nixpkgs: "nixpkgs")
	Args    map[string]interface{} `yaml:"args"`    // Custom arguments to pass to flake outputs (types inferred from YAML)
	Outputs []FlakeOutput          `yaml:"outputs"` // Which outputs to import
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name  string
	Value string
}

// Application represents an application to be installed
type Application struct {
	Name           string
	InstallCommand string
}

// BootstrapConfig represents the internal bootstrap configuration
type BootstrapConfig struct {
	Applications []Application
}
