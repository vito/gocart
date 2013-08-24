package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/xoebus/gocart/dependencies"
	"github.com/xoebus/gocart/dependency"
)

const CartridgeFile = "Cartridge"
const CartridgeLockFile = "Cartridge.lock"

func main() {
	app := cli.NewApp()
	app.Name = "gocart"
	app.Usage = "a go package manager"

	app.Commands = []cli.Command{
		{
			Name:      "install",
			ShortName: "i",
			Usage:     "install your dependencies",
			Action: func(c *cli.Context) {
				fmt.Println("Installing dependencies...")

				requestedDependencies := loadFile(CartridgeFile)
				lockedDependencies := loadFile(CartridgeLockFile)

				dependencies := dependencies.Merge(requestedDependencies, lockedDependencies)

				err := installDependencies(dependencies)
				if err != nil {
					log.Fatal(err)
				}

				file, err := os.Create("Cartridge.lock")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				err = updateLockFile(file, dependencies)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:      "update",
			ShortName: "u",
			Usage:     "update your dependencies",
			Action: func(c *cli.Context) {
				println("updating...")
			},
		},
	}

	app.Run(os.Args)
}

func loadFile(fileName string) []dependency.Dependency {
	depsFile, err := os.Open(fileName)
	if err != nil {
		return []dependency.Dependency{}
	}

	deps, err := dependencies.Parse(depsFile)
	if err != nil {
		log.Fatal(err)
	}

	return deps
}

func installDependencies(dependencies []dependency.Dependency) error {
	for _, dependency := range dependencies {
		err := dependency.Get()
		if err != nil {
			return err
		}

		gopath, err := InstallationDirectory(os.Getenv("GOPATH"))
		if err != nil {
			return err
		}

		err = dependency.Checkout(gopath)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateLockFile(writer io.Writer, dependencies []dependency.Dependency) error {
	for _, dependency := range dependencies {
		gopath, err := InstallationDirectory(os.Getenv("GOPATH"))
		if err != nil {
			return err
		}

		version, err := dependency.CurrentVersion(gopath)
		if err != nil {
			return err
		}

		writer.Write([]byte(fmt.Sprintf("%s\t%s\n", dependency.Path, version)))
	}

	return nil
}
