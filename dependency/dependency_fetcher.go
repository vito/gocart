package dependency

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return fetcher.commandRunner.Run(cmd)
}
