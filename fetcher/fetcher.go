package fetcher

import (
	"os"
	"os/exec"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/gopath"
	"github.com/vito/gocart/repository"
)

type Fetcher struct {
	runner command_runner.CommandRunner
	gopath string
}

func New(runner command_runner.CommandRunner) (*Fetcher, error) {
	gopath, err := gopath.InstallationDirectory(os.Getenv("GOPATH"))
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		runner: runner,
		gopath: gopath,
	}, nil
}

func (f *Fetcher) Fetch(dependency dependency.Dependency) (dependency.Dependency, error) {
	cmd := exec.Command("go", "get", "-d", "-v", dependency.Path)

	err := f.runner.Run(cmd)
	if err != nil {
		return dependency, err
	}

	repo, err := repository.New(dependency.FullPath(f.gopath), f.runner)
	if err != nil {
		return dependency, err
	}

	err = repo.Update()
	if err != nil {
		return dependency, err
	}

	err = repo.Checkout(dependency.Version)
	if err != nil {
		return dependency, err
	}

	currentVersion, err := repo.CurrentVersion()
	if err != nil {
		return dependency, err
	}

	dependency.Version = currentVersion

	return dependency, nil
}
