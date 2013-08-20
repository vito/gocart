package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/codegangsta/cli"

	"github.com/xoebus/gocart"
)

var skippableLine *regexp.Regexp = regexp.MustCompile(`^\s*(#.*)?\s*$`)
var dependencyList *regexp.Regexp = regexp.MustCompile(`^\s*([^\s]+)\s+([^\s]+)\s*$`)

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

				file, err := os.Open("Cartridge")
				if err != nil {
					log.Fatal(err)
				}

				bufferedReader := bufio.NewReader(file)
				for {
					line, err := bufferedReader.ReadString('\n')

					if err == io.EOF {
						break
					} else if err != nil {
						log.Fatal(err)
					}

					if skippableLine.MatchString(line) {
						continue
					}

					dependency := dependencyList.FindStringSubmatch(line)
					if dependency == nil {
						message := fmt.Sprintf("malformed line: %s", line)
						log.Fatalln(message)
					}

					if dependencyList.MatchString(line) {
						dep := &gocart.Dependency{
							Path:    dependency[1],
							Version: dependency[2],
						}

						err := dep.Get()
						if err != nil {
							log.Fatal(err)
						}

						gopath, err := gocart.InstallationDirectory(os.Getenv("GOPATH"))
						if err != nil {
							log.Fatal(err)
						}

						err = dep.Checkout(gopath)
						if err != nil {
							log.Fatal(err)
						}
					}
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
