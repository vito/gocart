package gocart

import (
	"fmt"
	"os/exec"
)

type Dependency struct {
	Path    string
	Version string
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s (%s)", d.Path, d.Version)
}

func (d Dependency) Get() error {
	return exec.Command("go", "get", "-u", "-d", d.Path).Run()
}
