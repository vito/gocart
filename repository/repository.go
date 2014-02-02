package repository

import (
	"errors"
	"os"
	"path"

	"github.com/vito/gocart/command_runner"
)

type Repository interface {
	Checkout(version string) error
	Update() error
	CurrentVersion() (string, error)
	Status() (string, error)
	Log(from, to string) (string, error)
}

var UnknownRepositoryType = errors.New("unknown repository type")

func New(path string, runner command_runner.CommandRunner) (Repository, error) {
	gitDepth := checkForDir(path, ".git", 0)
	hgDepth := checkForDir(path, ".hg", 0)
	bzrDepth := checkForDir(path, ".bzr", 0)

	if gitDepth < hgDepth && gitDepth < bzrDepth {
		return &GitRepository{path, runner}, nil
	}

	if hgDepth < gitDepth && hgDepth < bzrDepth {
		return &HgRepository{path, runner}, nil
	}

	if bzrDepth < gitDepth && bzrDepth < hgDepth {
		return &BzrRepository{path, runner}, nil
	}

	return nil, UnknownRepositoryType
}

func checkForDir(root, dir string, depth int) int {
	if root == "/" {
		return depth
	}

	_, err := os.Stat(path.Join(root, dir))
	if err == nil {
		return depth
	}

	return checkForDir(path.Dir(root), dir, depth+1)
}
