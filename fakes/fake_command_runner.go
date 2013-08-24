package fakes

import (
	"os/exec"
)

type FakeCommandRunner struct {
	LastCommand *exec.Cmd
}

func (runner *FakeCommandRunner) Run(cmd *exec.Cmd) error {
	runner.LastCommand = cmd
	return nil
}
