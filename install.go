package main

import (
	"fmt"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/fetcher"
	"github.com/vito/gocart/set"
)

func install(root string, recursive bool) {
	cartridge, err := set.LoadFrom(root)
	if err != nil {
		fatal(err)
	}

	err = installDependencies(cartridge, recursive, 0)
	if err != nil {
		fatal(err)
	}

	err = cartridge.SaveTo(root)
	if err != nil {
		fatal(err)
	}

	fmt.Println(green("OK"))
}

func installDependencies(deps *set.Set, recursive bool, depth int) error {
	runner := command_runner.New(false)

	fetcher, err := fetcher.New(runner)
	if err != nil {
		return err
	}

	maxWidth := 0

	for _, dep := range deps.Dependencies {
		if len(dep.Path) > maxWidth {
			maxWidth = len(dep.Path)
		}
	}

	for _, dep := range deps.Dependencies {
		fmt.Println(indent(depth, bold(dep.Path)+padding(maxWidth-len(dep.Path)+2)+cyan(dep.Version)))

		lockedDependency, err := processDependency(fetcher, dep)
		if err != nil {
			return err
		}

		FetchedDependencies[lockedDependency.Path] = lockedDependency

		deps.Replace(lockedDependency)

		if recursive {
			nextDeps, err := set.LoadFrom(lockedDependency.FullPath(GOPATH))
			if err == set.NoCartridgeError {
				continue
			} else if err != nil {
				return err
			}

			err = installDependencies(nextDeps, true, depth+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func processDependency(
	fetcher *fetcher.Fetcher,
	dep dependency.Dependency,
) (dependency.Dependency, error) {
	currentVersion := findCurrentVersion(dep)

	if currentVersion == dep.Version {
		return dep, nil
	}

	if err := checkForConflicts(dep); err != nil {
		return dependency.Dependency{}, err
	}

	lockedDependency, err := fetcher.Fetch(dep)
	if err != nil {
		return dependency.Dependency{}, err
	}

	return lockedDependency, nil
}

func checkForConflicts(dep dependency.Dependency) error {
	_, found := FetchedDependencies[dep.Path]
	if !found {
		return nil
	}

	status := getDependencyStatus(dep)
	if status == nil {
		return nil
	}

	if !status.VersionMatches {
		return VersionMismatch{
			Expected: dep.Version,
			Status:   *status,
		}
	}

	return nil
}
