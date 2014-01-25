package main

import (
	"os"

	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/dependency_reader"
)

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
