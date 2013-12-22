package gocart

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
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

func (d Dependency) CurrentVersion(gopath string) (string, error) {
	repoPath := d.FullPath(gopath)

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

	outBuf := new(bytes.Buffer)

	version.Stdout = outBuf
	version.Stderr = ioutil.Discard

	err := version.Run()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(outBuf.Bytes()), "\n"), nil
}

func (d Dependency) FullPath(gopath string) string {
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
