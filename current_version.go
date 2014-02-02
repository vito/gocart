package main

import (
	"strings"

	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_fetcher"
)

type DependencyStatus struct {
	VersionMatches bool
	CurrentVersion string

	Delta    int
	DeltaLog string
}

func findCurrentVersion(dep dependency.Dependency) string {
	repoPath := dep.FullPath(GOPATH)

	repo, err := dependency_fetcher.NewRepository(repoPath)
	if err != nil {
		return ""
	}

	current := repo.CurrentVersionCommand()
	current.Dir = repoPath

	currentVersion, err := current.Output()
	if err != nil {
		return ""
	}

	return strings.Trim(string(currentVersion), "\n ")
}

func getDependencyStatus(dep dependency.Dependency) *DependencyStatus {
	repoPath := dep.FullPath(GOPATH)

	repo, err := dependency_fetcher.NewRepository(repoPath)
	if err != nil {
		return nil
	}

	status := &DependencyStatus{}

	status.CurrentVersion = findCurrentVersion(dep)
	status.VersionMatches = status.CurrentVersion == dep.Version

	if status.VersionMatches {
		return status
	}

	logCmd := repo.LogCommand(dep.Version, status.CurrentVersion)
	logCmd.Dir = repoPath

	newer := true

	output, err := logCmd.Output()

	// git or hg with both refs will show empty if newer..older
	if len(output) == 0 {
		newer = false
	}

	// either bazaar, or dep.Version is not fetched
	if err != nil {
		newer = false
	}

	if newer {
		status.DeltaLog = string(output)
		status.Delta = len(strings.Split(status.DeltaLog, "\n")) - 1
	} else {
		logCmd := repo.LogCommand(status.CurrentVersion, dep.Version)
		logCmd.Dir = repoPath

		output, _ := logCmd.Output()

		status.DeltaLog = string(output)
		status.Delta = -(len(strings.Split(status.DeltaLog, "\n")) - 1)
	}

	return status
}
