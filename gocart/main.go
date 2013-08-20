package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/xoebus/gocart"
)

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

				cartridge, err := os.Open("Cartridge")
				if err != nil {
					log.Fatal(err)
				}

				dependencies, err := gocart.ParseDependencies(cartridge)
				if err != nil {
					log.Fatal(err)
				}

				err = installDependencies(dependencies)
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

func installDependencies(dependencies []gocart.Dependency) error {
	for _, dependency := range dependencies {
		err := dependency.Get()
		if err != nil {
			return err
		}

		gopath, err := gocart.InstallationDirectory(os.Getenv("GOPATH"))
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

func updateLockFile(writer io.Writer, dependencies []gocart.Dependency) error {
	for _, dependency := range dependencies {
		gopath, err := gocart.InstallationDirectory(os.Getenv("GOPATH"))
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
