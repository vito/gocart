package main

import (
	"fmt"
	"os"
	"path"

	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_fetcher"
	"github.com/vito/gocart/locker"
)

func check(root string) {
	dirtyDependencies := findDirtyDependencies(root)

	if len(dirtyDependencies) > 0 {
		for dep := range dirtyDependencies {
			fmt.Println(bold("dirty dependency:"), dep)
		}

		os.Exit(1)
	}
}

func findDirtyDependencies(root string) map[string]bool {
	requestedDependencies := loadFile(path.Join(root, CartridgeFile))
	lockedDependencies := loadFile(path.Join(root, CartridgeLockFile))

	dependencies := locker.GenerateLock(requestedDependencies, lockedDependencies)

	dirtyDependencies := make(map[string]bool)

	for _, dep := range dependencies {
		if checkForDirtyState(dep) {
			dirtyDependencies[dep.FullPath(GOPATH)] = true
		}

		for dep, _ := range findDirtyDependencies(dep.FullPath(GOPATH)) {
			dirtyDependencies[dep] = true
		}
	}

	return dirtyDependencies
}

func checkForDirtyState(dep dependency.Dependency) bool {
	repoPath := dep.FullPath(GOPATH)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return false
	}

	currentVersion := findCurrentVersion(dep)

	if currentVersion != dep.Version {
		fmt.Println("mismatch:", dep.Path, "should be", dep.Version, "is", currentVersion)
		return true
	}

	repo, err := dependency_fetcher.NewRepository(repoPath)
	if err != nil {
		fatal(err.Error())
	}

	status := repo.StatusCommand()
	status.Dir = repoPath

	output, err := status.CombinedOutput()
	if err != nil {
		fatal(err.Error())
	}

	// Bazaar is bizarre
	if string(output) == "working tree is out of date, run 'bzr update'\n" {
		return false
	}

	if len(output) != 0 {
		fmt.Println("dirty state:", dep.Path)
		return true
	}

	return false
}
