package utils

import (
	"os"
	"os/exec"
)

// RunCommand runs a shell command and returns an error if the command fails.
// The command is expected to be in the format of "command arg1 arg2 ..."
// This function is particularly useful for running commands that can't be replaced by Go code.
func RunCommand(comm string, args ...string) error {
	cmd := exec.Command(comm, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
