package util

import "os/exec"

// CommandExecutor is a function that executes a command and returns the output, error and exit code.
type CommandExecutor func(name string, arg ...string) ([]byte, error, int)

// ExecCommandExecutor executes a command with exec.Command and returns the output, error and exit code.
var ExecCommandExecutor CommandExecutor = func(name string, arg ...string) ([]byte, error, int) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.Output()
	exitCode := cmd.ProcessState.ExitCode()
	return output, err, exitCode
}
