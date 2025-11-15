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
	}
	// Load config and populate EnvVars if available
	user.Reload()
	return user
}

// Reload refreshes the user's configuration from camp.yml
// This loads environment variables from ~/.camp/camp.yml or ~/.camp/camp.yaml
func (u *User) Reload() error {
	config, err := LoadUserConfig(u.HomeDir)
	if err != nil {
		// If config loading fails, keep existing EnvVars
		return err
	}

	// Update EnvVars from config
	if config.Env != nil {
		u.EnvVars = config.Env
	} else {
		u.EnvVars = make(map[string]string)
	}

	return nil
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
