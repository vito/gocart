package gocart

import (
	"os"
	"os/exec"
)

type DependencyFetcher struct {
	commandRunner CommandRunner
}

func NewDependencyFetcher(runner CommandRunner) *DependencyFetcher {
	return &DependencyFetcher{
		commandRunner: runner,
	}
}

func (fetcher *DependencyFetcher) Fetch(dependency Dependency) error {
	cmd := exec.Command("go", "get", "-u", "-d", "-v", dependency.Path)
	fetcher.redirectIO(cmd)
	return fetcher.commandRunner.Run(cmd)
}

func (fetcher *DependencyFetcher) redirectIO(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
}
