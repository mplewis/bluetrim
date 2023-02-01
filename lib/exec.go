package lib

import (
	"os/exec"
)

// call calls a program with the given arguments and returns the output and exit code.
func call(program string, args ...string) (string, int, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.CombinedOutput()
	return string(out), cmd.ProcessState.ExitCode(), err
}
