package utils

import (
	"io"
	"os"
	"os/exec"
)

// Exec runs an external command attached to the terminal's stdin, writing its output to
// stdout and stderr. It is a variable so tests can replace it and never run real commands.
var Exec = func(stdout, stderr io.Writer, comm string, args ...string) error {
	cmd := exec.Command(comm, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

// RunCommand runs a shell command and returns an error if the command fails.
// The command is expected to be in the format of "command arg1 arg2 ..."
// This function is particularly useful for running commands that can't be replaced by Go code.
func RunCommand(comm string, args ...string) error {
	return Exec(os.Stdout, os.Stderr, comm, args...)
}

// RunCommands runs a list of shell commands in order and returns the first error encountered
// or nil if all commands succeed.
// The commands are expected to be in the format of []string{"command", "arg1", "arg2", ...}
func RunCommands(cmds [][]string) error {
	for _, cmd := range cmds {
		err := RunCommand(cmd[0], cmd[1:]...)
		if err != nil {
			return err
		}
	}
	return nil
}

// CommandReturn runs a shell command and returns the output as a string.
// Useful for when the output of the command is needed for further processing.
func CommandReturn(comm string, args ...string) (string, error) {
	cmd := exec.Command(comm, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
