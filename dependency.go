package gocart

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var UnknownDependencyType = errors.New("unknown dependency type")

type Dependency struct {
	Path    string
	Version string
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s\t%s", d.Path, d.Version)
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

	checkout.Stdout = os.Stdout
	checkout.Stderr = os.Stderr
	checkout.Stdin = os.Stdin

	err := checkout.Run()
	if err != nil {
		return err
	}

	return nil
}

func (d Dependency) CurrentVersion(gopath string) (string, error) {
	repoPath := d.fullPath(gopath)

	var version *exec.Cmd

	if findDirectory(repoPath, ".hg") {
		version = exec.Command("hg", "id", "-i")
	}

	if findDirectory(repoPath, ".git") {
		version = exec.Command("git", "rev-parse", "HEAD")
	}

	if findDirectory(repoPath, ".bzr") {
		version = exec.Command("bzr", "revno", "--tree")
	}

	if version == nil {
		return "", UnknownDependencyType
	}

	version.Dir = repoPath

	out, err := version.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(out), "\n"), nil
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
