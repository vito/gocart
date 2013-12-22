package dependency_fetcher

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/gopath"
)

type DependencyFetcher struct {
	commandRunner command_runner.CommandRunner
	gopath        string
}

func New(runner command_runner.CommandRunner) (*DependencyFetcher, error) {
	gopath, err := gopath.InstallationDirectory(os.Getenv("GOPATH"))
	if err != nil {
		return nil, err
	}

	return &DependencyFetcher{
		commandRunner: runner,
		gopath:        gopath,
	}, nil
}

func (f *DependencyFetcher) Fetch(dependency dependency.Dependency) (dependency.Dependency, error) {
	cmd := exec.Command("go", "get", "-u", "-d", "-v", dependency.Path)

	err := f.commandRunner.Run(cmd)
	if err != nil {
		return dependency, err
	}

	repo, err := NewRepository(dependency.FullPath(f.gopath))
	if err != nil {
		return dependency, err
	}

	cmd = repo.CheckoutCommand(dependency.Version)
	cmd.Dir = dependency.FullPath(f.gopath)

	err = f.commandRunner.Run(cmd)
	if err != nil {
		return dependency, err
	}

	current := repo.CurrentVersionCommand()
	current.Dir = dependency.FullPath(f.gopath)

	outBuf := new(bytes.Buffer)

	current.Stdout = outBuf
	current.Stderr = ioutil.Discard

	err = f.commandRunner.Run(current)
	if err != nil {
		return dependency, err
	}

	dependency.Version = strings.Trim(string(outBuf.Bytes()), "\n")

	return dependency, nil
}
