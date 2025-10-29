package util

// CommandExecutor is a function that executes a command and returns the output, error and exit code.
type CommandExecutor func(name string, arg ...string) ([]byte, error, int)
