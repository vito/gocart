package gocart

import (
	"os/exec"
)

type CommandRunner interface {
	Run(cmd *exec.Cmd) error
}

type ShellCommandRunner struct{}

func (runner *ShellCommandRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}
