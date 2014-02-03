package main

import (
	"strings"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/repository"
)

type DependencyStatus struct {
	VersionMatches bool
	CurrentVersion string

	Delta    int
	DeltaLog string
}

func findCurrentVersion(dep dependency.Dependency) string {
	repoPath := dep.FullPath(GOPATH)

	repo, err := repository.New(repoPath, command_runner.New(false))
	if err != nil {
		return ""
	}

	currentVersion, err := repo.CurrentVersion()
	if err != nil {
		return ""
	}

	return currentVersion
}

func getDependencyStatus(dep dependency.Dependency) *DependencyStatus {
	repoPath := dep.FullPath(GOPATH)

	repo, err := repository.New(repoPath, command_runner.New(false))
	if err != nil {
		return nil
	}

	status := &DependencyStatus{}

	status.CurrentVersion = findCurrentVersion(dep)
	status.VersionMatches = status.CurrentVersion == dep.Version

	if status.VersionMatches {
		return status
	}

	newer := true

	logOutput, err := repo.Log(dep.Version, status.CurrentVersion)

	// git or hg with both refs will show empty if newer..older
	if len(logOutput) == 0 {
		newer = false
	}

	// either bazaar, or dep.Version is not fetched
	if err != nil {
		newer = false
	}

	if newer {
		status.DeltaLog = logOutput
		status.Delta = len(strings.Split(status.DeltaLog, "\n")) - 1
	} else {
		logOutput, _ := repo.Log(status.CurrentVersion, dep.Version)

		status.DeltaLog = logOutput
		status.Delta = -(len(strings.Split(status.DeltaLog, "\n")) - 1)
	}

	return status
}
