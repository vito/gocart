package gocart

import (
	"fmt"
	"os"
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
	repoPath := d.fullPath(gopath)

	var checkout *exec.Cmd

	if _, err := os.Stat(path.Join(repoPath, ".git")); err == nil {
		checkout = exec.Command("git", "checkout", d.Version)
	}

	if _, err := os.Stat(path.Join(repoPath, ".bzr")); err == nil {
		checkout = exec.Command("bzr", "update", "-r", d.Version)
	}

	if checkout == nil {
		// TODO
		panic("unknown repo lol")
	}

	checkout.Dir = repoPath

	out, err := checkout.CombinedOutput()
	if err != nil {
		fmt.Printf("output:\n%s\n", out)
		return err
	}

	return nil
}

func (d Dependency) fullPath(gopath string) string {
	return path.Join(gopath, "src", d.Path)
}
