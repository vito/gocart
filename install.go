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

func install(root string, recursive bool, aggregate bool, trickleDown bool, depth int) {
	if _, err := os.Stat(path.Join(root, CartridgeFile)); err != nil {
		println("no Cartridge file!")
		os.Exit(1)
		return
	}

	requestedDependencies := loadFile(path.Join(root, CartridgeFile))
	lockedDependencies := loadFile(path.Join(root, CartridgeLockFile))

	dependencies := locker.GenerateLock(requestedDependencies, lockedDependencies)

	newLockedDependencies := installDependencies(dependencies, recursive, trickleDown, depth)

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
	trickleDown bool,
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
		fmt.Println(bold(dep.Path) + padding(maxWidth-len(dep.Path)+2) + cyan(dep.Version))

		lockedDependency := processDependency(fetcher, dep)

		if depth == 0 && trickleDown {
			TrickleDownDependencies[lockedDependency.Path] = lockedDependency
		}

		lockedDependencies = append(lockedDependencies, lockedDependency)
	}

	if recursive {
		installNextLevel(lockedDependencies, depth+1)
	}

	return lockedDependencies
}

func processDependency(
	fetcher *dependency_fetcher.DependencyFetcher,
	dep dependency.Dependency,
) dependency.Dependency {
	dep = trickledDown(dep)

	if findCurrentVersion(dep) == dep.Version {
		return dep
	}

	lockedDependency, err := fetcher.Fetch(dep)
	if err != nil {
		fatal(err.Error())
	}

	checkForConflicts(lockedDependency)

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

func installNextLevel(deps []dependency.Dependency, newDepth int) {
	for _, dependency := range deps {
		dependencyPath := dependency.FullPath(GOPATH)

		if _, err := os.Stat(path.Join(dependencyPath, "Cartridge")); err == nil {
			fmt.Println("\nfetching dependencies for", dependency.Path)
			install(dependencyPath, true, false, false, newDepth)
		}
	}
}

func trickledDown(dep dependency.Dependency) dependency.Dependency {
	trickled, found := TrickleDownDependencies[dep.Path]
	if found {
		return trickled
	}

	return dep
}

func checkForConflicts(lockedDependency dependency.Dependency) {
	currentVersion, found := FetchedDependencies[lockedDependency.Path]
	if found && currentVersion.Version != lockedDependency.Version {
		fatal("version conflict for " + currentVersion.Path)
	}
}
