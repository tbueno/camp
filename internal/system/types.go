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
