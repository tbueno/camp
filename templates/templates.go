// Package templates embeds the Nix templates camp writes to the user's environment,
// so the binary does not depend on the directory it runs from.
package templates

import "embed"

// FS holds the bootstrap templates under initial/ and the rebuild templates under files/.
//
//go:embed initial files
var FS embed.FS
