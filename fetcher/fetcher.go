package fetcher

import (
	"fmt"
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

	fetchedDependencies map[string]dependency.Dependency
}

type VersionConflictError struct {
	Path string

	VersionA string
	VersionB string
}

func (e VersionConflictError) Error() string {
	return fmt.Sprintf("version conflict for %s: %s and %s", e.Path, e.VersionA, e.VersionB)
}

func New(runner command_runner.CommandRunner) (*Fetcher, error) {
	gopath, err := gopath.InstallationDirectory(os.Getenv("GOPATH"))
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		runner: runner,
		gopath: gopath,

		fetchedDependencies: make(map[string]dependency.Dependency),
	}, nil
}

func (f *Fetcher) Fetch(dep dependency.Dependency) (dependency.Dependency, error) {
	var goGet *exec.Cmd

	repoPath := dep.FullPath(f.gopath)

	lockDown := true
	updateRepo := false

	if dep.BleedingEdge {
		// update the repo only if bleeding-edge and repo is clean
		if _, err := os.Stat(repoPath); err == nil {
			lockDown = false

			repo, err := repository.New(repoPath, f.runner)
			if err != nil {
				return dependency.Dependency{}, err
			}

			statusOut, err := repo.Status()
			if err != nil {
				return dependency.Dependency{}, err
			}

			if len(statusOut) == 0 {
				updateRepo = true
			}
		}
	}

	if updateRepo {
		goGet = exec.Command("go", "get", "-u", "-d", "-v", dep.Path)
	} else {
		goGet = exec.Command("go", "get", "-d", "-v", dep.Path)
	}

	err := f.runner.Run(goGet)
	if err != nil {
		return dependency.Dependency{}, err
	}

	repo, err := repository.New(repoPath, f.runner)
	if err != nil {
		return dependency.Dependency{}, err
	}

	if lockDown {
		err := f.syncRepo(repo, dep.Version)
		if err != nil {
			return dependency.Dependency{}, err
		}
	}

	currentVersion, err := repo.CurrentVersion()
	if err != nil {
		return dependency.Dependency{}, err
	}

	dep.Version = currentVersion

	fetched, found := f.fetchedDependencies[dep.Path]
	if found {
		if fetched.Version != dep.Version {
			return dependency.Dependency{}, VersionConflictError{
				Path:     dep.Path,
				VersionA: fetched.Version,
				VersionB: dep.Version,
			}
		}
	} else {
		f.fetchedDependencies[dep.Path] = dep
	}

	return dep, nil
}

func (f *Fetcher) syncRepo(repo repository.Repository, version string) error {
	currentVersion, err := repo.CurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion == version {
		// already up-to-date
		return nil
	}

	err = repo.Update()
	if err != nil {
		return err
	}

	return repo.Checkout(version)
}
