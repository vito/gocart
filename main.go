package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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

	newLockedDependencies, err := installDependencies(dependencies, recursive, trickleDown, depth)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(path.Join(root, "Cartridge.lock"))
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	if depth == 0 {
		fmt.Println("\x1b[32mOK\x1b[0m")
	}
}

func help() {
	fmt.Println(`gocart: a go package manager

Usage:

    'gocart install' or 'gocart':
        Install dependencies described by Cartridge.lock or Cartridge, and
        update Cartridge.lock with locked-down dependency versions.

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
		log.Fatal(err)
	}

	return dependencies
}

func installDependencies(dependencies []dependency.Dependency, recursive bool, trickleDown bool, depth int) ([]dependency.Dependency, error) {
	runner := command_runner.New()

	lockedDependencies := []dependency.Dependency{}

	fetcher, err := dependency_fetcher.New(runner)
	if err != nil {
		return []dependency.Dependency{}, err
	}

	maxWidth := 0

	for _, dep := range dependencies {
		if len(dep.Path) > maxWidth {
			maxWidth = len(dep.Path)
		}
	}

	gopath, err := gopath.InstallationDirectory(os.Getenv("GOPATH"))
	if err != nil {
		return nil, err
	}

	nextLevel := []dependency.Dependency{}

	for _, dep := range dependencies {
		fmt.Println("\x1b[1m" + dep.Path + "\x1b[0m" + strings.Repeat(" ", maxWidth-len(dep.Path)+2) + "\x1b[36m" + dep.Version + "\x1b[0m")

		trickled, found := TrickleDownDependencies[dep.Path]
		if found {
			dep = trickled
		}

		lockedDependency, err := fetcher.Fetch(dep)
		if err != nil {
			return []dependency.Dependency{}, err
		}

		currentVersion, found := FetchedDependencies[lockedDependency.Path]
		if found && currentVersion.Version != lockedDependency.Version {
			log.Fatalln("version conflict for", currentVersion.Path)
		}

		FetchedDependencies[lockedDependency.Path] = lockedDependency

		lockedDependencies = append(lockedDependencies, lockedDependency)

		if depth == 0 && trickleDown {
			TrickleDownDependencies[lockedDependency.Path] = lockedDependency
		}

		if !recursive {
			continue
		}

		nextLevel = append(nextLevel, lockedDependency)
	}

	for _, dependency := range nextLevel {
		dependencyPath := dependency.FullPath(gopath)

		if _, err := os.Stat(path.Join(dependencyPath, "Cartridge")); err == nil {
			fmt.Println("\nfetching dependencies for", dependency.Path)
			install(dependencyPath, recursive, false, trickleDown, depth+1)
		}
	}

	return lockedDependencies, nil
}

func updateLockFile(writer io.Writer, dependencies []dependency.Dependency) error {
	for _, dependency := range dependencies {
		writer.Write([]byte(fmt.Sprintf("%s\t%s\n", dependency.Path, dependency.Version)))
	}

	return nil
}
