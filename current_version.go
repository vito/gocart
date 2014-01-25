package main

import (
	"strings"

	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_fetcher"
)

func findCurrentVersion(dep dependency.Dependency) string {
	repoPath := dep.FullPath(GOPATH)

	repo, err := dependency_fetcher.NewRepository(repoPath)
	if err != nil {
		return ""
	}

	current := repo.CurrentVersionCommand()
	current.Dir = repoPath

	currentVersion, err := current.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.Trim(string(currentVersion), "\n ")
}
