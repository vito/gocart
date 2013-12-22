package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_fetcher"
	"github.com/vito/gocart/dependency_reader"
	"github.com/vito/gocart/locker"
)

const GocartVersion = "0.1.0"

const CartridgeFile = "Cartridge"
const CartridgeLockFile = "Cartridge.lock"

func main() {
	command := ""

	if len(os.Args) == 1 {
		command = "install"
	} else {
		command = os.Args[1]
	}

	switch command {
	case "install":
		install()
	case "help", "--help", "-h":
		help()
	case "version", "--version", "-v":
		version()
	default:
		unknownCommand()
	}
}

func install() {
	requestedDependencies := loadFile(CartridgeFile)
	lockedDependencies := loadFile(CartridgeLockFile)

	dependencies := locker.GenerateLock(requestedDependencies, lockedDependencies)

	newLockedDependencies, err := installDependencies(dependencies)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("Cartridge.lock")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = updateLockFile(file, newLockedDependencies)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OK")
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

func installDependencies(dependencies []dependency.Dependency) ([]dependency.Dependency, error) {
	runner := command_runner.New()

	lockedDependencies := []dependency.Dependency{}

	fetcher, err := dependency_fetcher.New(runner)
	if err != nil {
		return []dependency.Dependency{}, err
	}

	for _, dep := range dependencies {
		fmt.Println("fetching", dep.Path)
		fmt.Println("  at", dep.Version)

		lockedDependency, err := fetcher.Fetch(dep)
		if err != nil {
			return []dependency.Dependency{}, err
		}

		lockedDependencies = append(lockedDependencies, lockedDependency)
	}

	return lockedDependencies, nil
}

func updateLockFile(writer io.Writer, dependencies []dependency.Dependency) error {
	for _, dependency := range dependencies {
		writer.Write([]byte(fmt.Sprintf("%s\t%s\n", dependency.Path, dependency.Version)))
	}

	return nil
}
