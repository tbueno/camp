# Camp - Golang CLI Application

## Project Overview
Camp is a command line application built with Go and the Cobra CLI framework. It serves as your helpful command line companion.

## Tech Stack
- **Language**: Go 1.24.4
- **CLI Framework**: Cobra (github.com/spf13/cobra v1.10.0)
- **Testing**: Go's built-in testing package

## Project Structure
```
camp/
├── go.mod                 # Go module definition
├── main.go               # Application entry point
└── cmd/
    ├── root.go           # Root command implementation
    └── root_test.go      # Unit tests for root command
```

## Development Commands
- **Build**: `go build`
- **Run**: `go run main.go`
- **Test**: `go test ./...`

## Code Style & Conventions
- Follow standard Go conventions. Run `go fmt` after executing commands that alter code.
- Use Cobra's best practices for CLI commands
- All commands should use `cmd.OutOrStdout()` for output to support testing
- Write comprehensive unit tests for all commands
- Keep command logic in the `cmd/` package
- Main entry point should be minimal, just calling `cmd.Execute()`
- When new features are added or changed, make sure to update the documentation in `doc/` folder.

## Testing Guidelines
- Test command names, descriptions, and output
- Use buffer capture for testing command output
- Test both success and error scenarios
- Maintain high test coverage
- Always run the tests when code is changed!

## Current Commands

The source code commands can be executes with the command `go run main.go <command>`. Running the command `go run main.go` will display the help message for all the commands.
