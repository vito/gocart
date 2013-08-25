package gocart

import (
	"os"
	"os/exec"
)

type DependencyFetcher struct {
	commandRunner CommandRunner
	gopath        string
}

func NewDependencyFetcher(runner CommandRunner) (*DependencyFetcher, error) {
	gopath, err := InstallationDirectory(os.Getenv("GOPATH"))
	if err != nil {
		return nil, err
	}

	return &DependencyFetcher{
		commandRunner: runner,
		gopath:        gopath,
	}, nil
}

func (fetcher *DependencyFetcher) Fetch(dependency Dependency) error {
	cmd := exec.Command("go", "get", "-u", "-d", "-v", dependency.Path)
	fetcher.redirectIO(cmd)

	err := fetcher.commandRunner.Run(cmd)
	if err != nil {
		return err
	}

	repo, err := NewRepository(dependency.FullPath(fetcher.gopath))
	if err != nil {
		return err
	}

	cmd = repo.CheckoutCommand(dependency.Version)
	cmd.Dir = dependency.FullPath(fetcher.gopath)
	fetcher.redirectIO(cmd)

	err = fetcher.commandRunner.Run(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (fetcher *DependencyFetcher) redirectIO(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
}
