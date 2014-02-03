package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vito/gocart/gopath"
)

const GocartVersion = "0.1.0"

const CartridgeFile = "Cartridge"
const CartridgeLockFile = "Cartridge.lock"

var GOPATH string

var recursive = flag.Bool(
	"r",
	false,
	"recursively fetch dependencies and check for conflicts",
)

var exclude = flag.String(
	"x",
	"",
	"exclude dependencies matching any of the (comma-separated) tags",
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
		install(".", *recursive, strings.Split(*exclude, ","))
		return
	}

	if command == "check" {
		check(".")
		return
	}

	unknownCommand()
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
