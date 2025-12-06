package main

import "camp/cmd"

var (
	// Version is set via ldflags during build
	Version = "dev"
	// Commit is set via ldflags during build
	Commit = "none"
)

func main() {
	cmd.Execute()
}
