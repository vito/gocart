package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_fetcher"
	"github.com/vito/gocart/locker"
)

type Reason interface {
	Description() string
}

type VersionMismatch struct {
	Expected string
	Current  string
}

func (self VersionMismatch) Description() string {
	return fmt.Sprintf("version mismatch:\n  want %s\n  have %s\n", red(self.Expected), green(self.Current))
}

type DirtyState struct {
	Output string
}

func (self DirtyState) Description() string {
	output := ""

	for _, line := range strings.Split(strings.Trim(self.Output, " \n"), "\n") {
		output = output + "  " + line + "\n"
	}

	return fmt.Sprintf("dirty state:\n%s", output)
}

func check(root string) {
	dirtyDependencies := findDirtyDependencies(root)

	if len(dirtyDependencies) > 0 {
		count := 0

		for path, reason := range dirtyDependencies {
			fmt.Println(bold(path))
			fmt.Println(reason.Description())

			count++
		}

		os.Exit(1)
	}
}

func findDirtyDependencies(root string) map[string]Reason {
	requestedDependencies := loadFile(path.Join(root, CartridgeFile))
	lockedDependencies := loadFile(path.Join(root, CartridgeLockFile))

	dependencies := locker.GenerateLock(requestedDependencies, lockedDependencies)

	dirtyDependencies := make(map[string]Reason)

	for _, dep := range dependencies {
		reason := checkForDirtyState(dep)

		if reason != nil {
			dirtyDependencies[dep.FullPath(GOPATH)] = reason
		}

		for dep, reason := range findDirtyDependencies(dep.FullPath(GOPATH)) {
			dirtyDependencies[dep] = reason
		}
	}

	return dirtyDependencies
}

func checkForDirtyState(dep dependency.Dependency) Reason {
	repoPath := dep.FullPath(GOPATH)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil
	}

	repo, err := dependency_fetcher.NewRepository(repoPath)
	if err != nil {
		fatal(err.Error())
	}

	status := repo.StatusCommand()
	status.Dir = repoPath

	output, err := status.Output()
	if err != nil {
		fatal(err.Error())
	}

	// Bazaar is bizarre
	if string(output) == "working tree is out of date, run 'bzr update'\n" {
		return nil
	}

	if len(output) != 0 {
		return DirtyState{
			Output: string(output),
		}
	}

	currentVersion := findCurrentVersion(dep)

	if currentVersion != dep.Version {
		return VersionMismatch{
			Expected: dep.Version,
			Current:  currentVersion,
		}
	}

	return nil
}
