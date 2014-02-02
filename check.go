package main

import (
	"fmt"
	"os"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/repository"
	"github.com/vito/gocart/set"
)

type VersionMismatch struct {
	Expected string
	Status   DependencyStatus
}

func (self VersionMismatch) Error() string {
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

func (self DirtyState) Error() string {
	return fmt.Sprintf("dirty state:\n%s", indent(1, self.Output))
}

func check(root string) {
	cartridge, err := set.LoadFrom(root)
	if err != nil {
		fatal(err)
	}

	dirty := checkDependencies(cartridge, 0)

	if dirty {
		os.Exit(1)
	}
}

func checkDependencies(deps *set.Set, depth int) bool {
	dirty := false

	for _, dep := range deps.Dependencies {
		err := checkForDirtyState(dep)
		if err != nil {
			dirty = true

			fmt.Println(indent(depth, bold(dep.Path)))
			fmt.Println(indent(depth+1, err.Error()))
		} else {
			fmt.Println(indent(depth, bold(dep.Path)), green("OK"))
		}

		nextDeps, err := set.LoadFrom(dep.FullPath(GOPATH))
		if err == set.NoCartridgeError {
			continue
		} else if err != nil {
			fatal(err)
		}

		if checkDependencies(nextDeps, depth+1) {
			dirty = true
		}
	}

	return dirty
}

func checkForDirtyState(dep dependency.Dependency) error {
	repoPath := dep.FullPath(GOPATH)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil
	}

	repo, err := repository.New(repoPath, command_runner.New(false))
	if err != nil {
		fatal(err)
	}

	statusOut, err := repo.Status()
	if err != nil {
		fatal(err)
	}

	if len(statusOut) != 0 {
		return DirtyState{
			Output: statusOut,
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
