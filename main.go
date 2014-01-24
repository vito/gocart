package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_fetcher"
	"github.com/vito/gocart/dependency_reader"
	"github.com/vito/gocart/gopath"
	"github.com/vito/gocart/locker"
)

const GocartVersion = "0.1.0"

const CartridgeFile = "Cartridge"
const CartridgeLockFile = "Cartridge.lock"

var GOPATH string

var FetchedDependencies = make(map[string]dependency.Dependency)
var TrickleDownDependencies = make(map[string]dependency.Dependency)

var recursive = flag.Bool(
	"r",
	false,
	"recursively fetch dependencies and check for conflicts",
)

var showHelp = flag.Bool(
	"h",
	false,
	"show command usage (this text)",
)

var showVersion = flag.Bool(
	"v",
	false,
	"show gocart version",
)

var trickleDown = flag.Bool(
	"t",
	false,
	"trickle down dependencies into nested Cartridges",
)

var aggregate = flag.Bool(
	"a",
	false,
	"collect recursive dependencies into a single .lock file",
)

func main() {
	flag.Parse()

	args := flag.Args()

	command := ""

	gopath, err := gopath.InstallationDirectory(os.Getenv("GOPATH"))
	if err != nil {
		fatal("GOPATH is not set.")
	}

	GOPATH = gopath

	if len(args) == 0 {
		command = "install"
	} else {
		command = args[0]
	}

	if *showHelp || command == "help" {
		help()
		return
	}

	if *showVersion || command == "version" {
		version()
		return
	}

	if command == "install" {
		install(".", *recursive, *aggregate, *trickleDown, 0)
		return
	}

	if command == "check" {
		check(".")
		return
	}

	unknownCommand()
}

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
		fmt.Println("\x1b[32mOK\x1b[0m")
	}
}

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

func help() {
	fmt.Println(`gocart: a go package manager

Usage:

  'gocart':
    Install dependencies described by Cartridge.lock or Cartridge, and
    update Cartridge.lock with locked-down dependency versions.

    The following flags are handled:

      -r: (recurse) if each dependency has a Cartridge, recursively run
          gocart for it as well

      -t: (trickle down) with -r, override lower-level dependencies with
          root dependencies

      -a: (aggregate) with -r, collect all recursive dependencies into the
          top level lockfile

  'gocart check':
    Check if any of the dependencies are in a modified/dirty state.

Place your dependencies in a file called Cartridge with this format:

[import path]	[vcs ref]

'gocart install' will 'go get' each import path and switch it to the given ref.
After getting them, it will take the "hard" ref (e.g. the sha in git), and save
it in Cartridge.lock. The Cartridge.lock has the same format as Cartridge and
has the same semantics; it will later be used by 'gocart install' if it exists.

To update an individual dependency, simply remove its line from Cartridge.lock
and run 'gocart install'. To update all dependencies, remove Cartridge.lock.
`)
}

func version() {
	fmt.Println(GocartVersion)
}

func unknownCommand() {
	fmt.Println("unknown command:", os.Args[1])
	fmt.Println()
	help()
	os.Exit(1)
}

func loadFile(fileName string) []dependency.Dependency {
	cartridge, err := os.Open(fileName)
	if err != nil {
		return []dependency.Dependency{}
	}

	reader := dependency_reader.New(cartridge)

	dependencies, err := reader.ReadAll()
	if err != nil {
		fatal(err.Error())
	}

	return dependencies
}

func checkForDirtyState(dep dependency.Dependency) bool {
	repoPath := dep.FullPath(GOPATH)

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

	return len(output) != 0
}

func installDependencies(dependencies []dependency.Dependency, recursive bool, trickleDown bool, depth int) []dependency.Dependency {
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

func processDependency(fetcher *dependency_fetcher.DependencyFetcher, dep dependency.Dependency) dependency.Dependency {
	dep = trickledDown(dep)

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

func bold(str string) string {
	return "\x1b[1m" + str + "\x1b[0m"
}

func red(str string) string {
	return "\x1b[31m" + str + "\x1b[0m"
}

func cyan(str string) string {
	return "\x1b[36m" + str + "\x1b[0m"
}

func padding(size int) string {
	return strings.Repeat(" ", size)
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, red(message))
	os.Exit(1)
}
