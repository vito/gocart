package command_runner

import (
	"os/exec"
)

type CommandRunner interface {
	Run(cmd *exec.Cmd) error
}

type ShellCommandRunner struct{}

func New() *ShellCommandRunner {
	return &ShellCommandRunner{}
}

func (runner *ShellCommandRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}
