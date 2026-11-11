// Package utilstest provides test helpers for code that runs external commands through utils.Exec.
package utilstest

import (
	"io"
	"slices"
	"strings"
	"testing"

	"camp/internal/utils"
)

// ExecRecorder records the commands sent to utils.Exec instead of running them.
type ExecRecorder struct {
	Calls [][]string
	// Err is returned by every recorded command
	Err error
}

// StubExec replaces utils.Exec with a recorder for the duration of the test.
func StubExec(t *testing.T) *ExecRecorder {
	t.Helper()

	recorder := &ExecRecorder{}
	original := utils.Exec
	utils.Exec = func(_, _ io.Writer, comm string, args ...string) error {
		recorder.Calls = append(recorder.Calls, append([]string{comm}, args...))
		return recorder.Err
	}
	t.Cleanup(func() { utils.Exec = original })

	return recorder
}

// Called reports whether a recorded command line contains all the given parts, in order.
func (r *ExecRecorder) Called(parts ...string) bool {
	for _, call := range r.Calls {
		if containsInOrder(call, parts) {
			return true
		}
	}
	return false
}

// String renders the recorded commands, one per line, for test failure messages.
func (r *ExecRecorder) String() string {
	lines := make([]string, len(r.Calls))
	for i, call := range r.Calls {
		lines[i] = strings.Join(call, " ")
	}
	return strings.Join(lines, "\n")
}

func containsInOrder(call, parts []string) bool {
	start := 0
	for _, part := range parts {
		idx := slices.Index(call[start:], part)
		if idx < 0 {
			return false
		}
		start += idx + 1
	}
	return true
}
