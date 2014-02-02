package main

import (
	"fmt"
	"os"
	"path"

	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_fetcher"
	"github.com/vito/gocart/locker"
)

type Reason interface {
	Description() string
}

type VersionMismatch struct {
	Expected string
	Status   DependencyStatus
}

func (self VersionMismatch) Description() string {
	want := indent(1, "want "+red(self.Expected))
	have := indent(1, "have "+green(self.Status.CurrentVersion))
	log := indent(1, self.Status.DeltaLog)

	if self.Status.Delta > 0 {
		have = have + " (" + bold(fmt.Sprintf("%d", self.Status.Delta)) + " ahead)"
		have = have + "\n" + indent(1, "extra commits:\n"+log)
	} else if self.Status.Delta == 0 {
		have = have + " (? behind)"
	} else {
		have = have + " (" + bold(fmt.Sprintf("%d", -self.Status.Delta)) + " behind)"
		have = have + "\n" + indent(1, "missing commits:\n"+log)
	}

	return fmt.Sprintf(
		"version mismatch:\n%s\n%s\n",
		want,
		have,
	)
}

type DirtyState struct {
	Output string
}

func (self DirtyState) Description() string {
	return fmt.Sprintf("dirty state:\n%s", indent(1, self.Output))
}

func check(root string) {
	dirty := checkDependencies(root, 0)

	if dirty {
		os.Exit(1)
	}
}

func checkDependencies(root string, depth int) bool {
	requestedDependencies := loadFile(path.Join(root, CartridgeFile))
	lockedDependencies := loadFile(path.Join(root, CartridgeLockFile))

	dependencies := locker.GenerateLock(requestedDependencies, lockedDependencies)

	dirty := false

	for _, dep := range dependencies {
		reason := checkForDirtyState(dep)
		if reason != nil {
			dirty = true

			fmt.Println(indent(depth, bold(dep.Path)))
			fmt.Println(indent(depth+1, reason.Description()))
		} else {
			fmt.Println(indent(depth, bold(dep.Path)), green("OK"))
		}

		if checkDependencies(dep.FullPath(GOPATH), depth+1) {
			dirty = true
		}
	}

	return dirty
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

	currentStatus := getDependencyStatus(dep)
	if currentStatus != nil {
		if !currentStatus.VersionMatches {
			return VersionMismatch{
				Expected: dep.Version,
				Status:   *currentStatus,
			}
		}
	}

	return nil
}
