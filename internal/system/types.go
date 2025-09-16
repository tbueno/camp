package system

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
// This struct is prepared for future expansion
type User struct {
	Name string
	// Future fields: ID, Groups, Shell, etc.
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
