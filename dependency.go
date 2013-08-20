package gocart

import (
	"fmt"
	"os/exec"
	"path"
)

type Dependency struct {
	Path    string
	Version string
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s (%s)", d.Path, d.Version)
}

func (d Dependency) Get() error {
	return exec.Command("go", "get", "-u", "-d", "-v", d.Path).Run()
}

func (d Dependency) Checkout(gopath string) error {
	checkout := exec.Command("git", "checkout", d.Version)
	checkout.Dir = path.Join(gopath, "src", d.Path)
	return checkout.Run()
}
