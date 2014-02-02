package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_fetcher"
	"github.com/vito/gocart/locker"
)

func install(root string, recursive bool, aggregate bool, depth int) {
	if _, err := os.Stat(path.Join(root, CartridgeFile)); err != nil {
		println("no Cartridge file!")
		os.Exit(1)
		return
	}

	requestedDependencies := loadFile(path.Join(root, CartridgeFile))
	lockedDependencies := loadFile(path.Join(root, CartridgeLockFile))

	dependencies := locker.GenerateLock(requestedDependencies, lockedDependencies)

	newLockedDependencies := installDependencies(dependencies, recursive, depth)

	file, err := os.Create(path.Join(root, "Cartridge.lock"))
	if err != nil {
		fatal(err.Error())
	}
	defer file.Close()

	var dependenciesToBeWritten []dependency.Dependency

	if aggregate {
		for _, dep := range FetchedDependencies {
			dependenciesToBeWritten = append(dependenciesToBeWritten, dep)
		}
	} else {
		dependenciesToBeWritten = newLockedDependencies
	}

	err = updateLockFile(file, dependenciesToBeWritten)
	if err != nil {
		fatal(err.Error())
	}

	if depth == 0 {
		fmt.Println(green("OK"))
	}
}

func installDependencies(
	dependencies []dependency.Dependency,
	recursive bool,
	depth int,
) []dependency.Dependency {
	runner := command_runner.New()

	fetcher, err := dependency_fetcher.New(runner)
	if err != nil {
		fatal(err.Error())
	}

	maxWidth := 0

	for _, dep := range dependencies {
		if len(dep.Path) > maxWidth {
			maxWidth = len(dep.Path)
		}
	}

	lockedDependencies := []dependency.Dependency{}

	for _, dep := range dependencies {
		fmt.Println(indent(depth, bold(dep.Path) + padding(maxWidth-len(dep.Path)+2) + cyan(dep.Version)))

		lockedDependency := processDependency(fetcher, dep)

		lockedDependencies = append(lockedDependencies, lockedDependency)

		if recursive {
			dependencyPath := lockedDependency.FullPath(GOPATH)

			if _, err := os.Stat(path.Join(dependencyPath, "Cartridge")); err == nil {
				install(dependencyPath, true, false, depth+1)
			}
		}
	}

	return lockedDependencies
}

func processDependency(
	fetcher *dependency_fetcher.DependencyFetcher,
	dep dependency.Dependency,
) dependency.Dependency {
	currentVersion := findCurrentVersion(dep)

	if currentVersion == dep.Version {
		FetchedDependencies[dep.Path] = dep
		return dep
	}

	checkForConflicts(dep)

	lockedDependency, err := fetcher.Fetch(dep)
	if err != nil {
		fatal(err.Error())
	}

	FetchedDependencies[lockedDependency.Path] = lockedDependency

	return lockedDependency
}

func updateLockFile(writer io.Writer, dependencies []dependency.Dependency) error {
	for _, dependency := range dependencies {
		_, err := writer.Write([]byte(fmt.Sprintf("%s\t%s\n", dependency.Path, dependency.Version)))
		if err != nil {
			return err
		}
	}

	return nil
}

func checkForConflicts(dep dependency.Dependency) {
	_, found := FetchedDependencies[dep.Path]
	if !found {
		return
	}

	status := getDependencyStatus(dep)

	if status != nil {
		if !status.VersionMatches {
			mismatch := VersionMismatch{
				Expected: dep.Version,
				Status:   *status,
			}

			fmt.Println(mismatch.Description())

			os.Exit(1)
		}
	}
}
