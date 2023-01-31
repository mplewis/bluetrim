package main

import "os/exec"

// call calls a program with the given arguments and returns the output.
func call(program string, args ...string) (string, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
