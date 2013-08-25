package fakes

import (
	"os/exec"
)

type FakeCommandRunner struct {
	Commands    []*exec.Cmd
	LastCommand *exec.Cmd
}

func (runner *FakeCommandRunner) Run(cmd *exec.Cmd) error {
	runner.LastCommand = cmd
	runner.Commands = append(runner.Commands, cmd)
	return nil
}
