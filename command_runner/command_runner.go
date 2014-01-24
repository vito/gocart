package command_runner

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

type CommandRunner interface {
	Run(cmd *exec.Cmd) error
}

type CommandError struct {
	RunError error

	Command *exec.Cmd
	Output  []byte
}

func (e CommandError) Error() string {
	return fmt.Sprintf(
		"command %v failed with %s:\n%s",
		e.Command.Args,
		e.RunError,
		e.Output,
	)
}

type ShellCommandRunner struct{}

func New() *ShellCommandRunner {
	return &ShellCommandRunner{}
}

func (runner *ShellCommandRunner) Run(cmd *exec.Cmd) error {
	output := new(bytes.Buffer)

	if cmd.Stdout != nil && cmd.Stdout != ioutil.Discard {
		cmd.Stdout = io.MultiWriter(output, cmd.Stdout)
	} else {
		cmd.Stdout = output
	}

	if cmd.Stderr != nil && cmd.Stderr != ioutil.Discard {
		cmd.Stderr = io.MultiWriter(output, cmd.Stderr)
	} else {
		cmd.Stderr = output
	}

	err := cmd.Run()
	if err != nil {
		return CommandError{
			Command:  cmd,
			Output:   output.Bytes(),
			RunError: err,
		}
	}

	return nil
}
