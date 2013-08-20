package gocart

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

var UnknownDependencyType = errors.New("unknown dependency type")

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

	if findDirectory(repoPath, ".hg") {
		checkout = exec.Command("hg", "update", "-c", d.Version)
	}

	if findDirectory(repoPath, ".git") {
		checkout = exec.Command("git", "checkout", d.Version)
	}

	if findDirectory(repoPath, ".bzr") {
		checkout = exec.Command("bzr", "update", "-r", d.Version)
	}

	if checkout == nil {
		return UnknownDependencyType
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

func findDirectory(root, dir string) bool {
	if root == "/" {
		return false
	}

	_, err := os.Stat(path.Join(root, dir))
	if err == nil {
		return true
	}

	return findDirectory(path.Dir(root), dir)
}
