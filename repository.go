package gocart

import (
	"errors"
	"os"
	"os/exec"
	"path"
)

type Repository interface {
	CheckoutCommand(version string) *exec.Cmd
}

var UnknownRepositoryType = errors.New("unknown repository type")

func NewRepository(repoPath string) (Repository, error) {
	if checkForDir(repoPath, ".git") {
		return &GitRepository{}, nil
	}

	if checkForDir(repoPath, ".hg") {
		return &HgRepository{}, nil
	}

	if checkForDir(repoPath, ".bzr") {
		return &BzrRepository{}, nil
	}

	return nil, UnknownRepositoryType
}

func checkForDir(root, dir string) bool {
	if root == "/" {
		return false
	}

	_, err := os.Stat(path.Join(root, dir))
	if err == nil {
		return true
	}

	return findDirectory(path.Dir(root), dir)
}

type GitRepository struct{}

func (repo *GitRepository) CheckoutCommand(version string) *exec.Cmd {
	return exec.Command("git", "checkout", version)
}

type HgRepository struct{}

func (repo *HgRepository) CheckoutCommand(version string) *exec.Cmd {
	return exec.Command("hg", "update", "-c", version)
}

type BzrRepository struct{}

func (repo *BzrRepository) CheckoutCommand(version string) *exec.Cmd {
	return exec.Command("bzr", "update", "-r", version)
}
