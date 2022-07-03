package exec

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/dannykopping/errata/sample/errata"
)

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Execute executes a given command and arguments, and returns the result
func Execute(command string, args []string) (*Result, error) {
	if command == "" {
		return nil, errata.NewMissingCommandErr(nil)
	}

	if _, err := os.Stat(command); err != nil {
		return nil, errata.NewScriptNotFoundErr(err, command)
	}

	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	res := &Result{
		ExitCode: cmd.ProcessState.ExitCode(),
		Stderr:   stderr.String(),
		Stdout:   stdout.String(),
	}

	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			res.ExitCode = e.ExitCode()
		}

		return res, errata.NewScriptExecutionFailedErr(err)
	}

	return res, nil
}
