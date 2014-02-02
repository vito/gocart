package dependency_fetcher

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type Repository interface {
	CheckoutCommand(version string) *exec.Cmd
	UpdateCommand() *exec.Cmd
	CurrentVersionCommand() *exec.Cmd
	StatusCommand() *exec.Cmd
	LogCommand(from, to string) *exec.Cmd
}

var UnknownRepositoryType = errors.New("unknown repository type")

func NewRepository(repoPath string) (Repository, error) {
	gitDepth := checkForDir(repoPath, ".git", 0)
	hgDepth := checkForDir(repoPath, ".hg", 0)
	bzrDepth := checkForDir(repoPath, ".bzr", 0)

	if gitDepth < hgDepth && gitDepth < bzrDepth {
		return &GitRepository{}, nil
	}

	if hgDepth < gitDepth && hgDepth < bzrDepth {
		return &HgRepository{}, nil
	}

	if bzrDepth < gitDepth && bzrDepth < hgDepth {
		return &BzrRepository{}, nil
	}

	return nil, UnknownRepositoryType
}

type GitRepository struct{}

func (repo *GitRepository) CheckoutCommand(version string) *exec.Cmd {
	return exec.Command("git", "checkout", version)
}

func (repo *GitRepository) CurrentVersionCommand() *exec.Cmd {
	return exec.Command("git", "rev-parse", "HEAD")
}

func (repo *GitRepository) UpdateCommand() *exec.Cmd {
	return exec.Command("git", "fetch")
}

func (repo *GitRepository) StatusCommand() *exec.Cmd {
	return exec.Command("git", "status", "--porcelain")
}

func (repo *GitRepository) LogCommand(from, to string) *exec.Cmd {
	return exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..%s", from, to))
}

type HgRepository struct{}

func (repo *HgRepository) CheckoutCommand(version string) *exec.Cmd {
	return exec.Command("hg", "update", "-c", version)
}

func (repo *HgRepository) CurrentVersionCommand() *exec.Cmd {
	return exec.Command("hg", "id", "-i")
}

func (repo *HgRepository) UpdateCommand() *exec.Cmd {
	return exec.Command("hg", "pull")
}

func (repo *HgRepository) StatusCommand() *exec.Cmd {
	return exec.Command("hg", "status")
}

func (repo *HgRepository) LogCommand(from, to string) *exec.Cmd {
	return exec.Command("hg", "log", "--template", "{node}\n", "-r", fmt.Sprintf("%s::%s", from, to))
}

type BzrRepository struct{}

func (repo *BzrRepository) CheckoutCommand(version string) *exec.Cmd {
	return exec.Command("bzr", "update", "-r", version)
}

func (repo *BzrRepository) CurrentVersionCommand() *exec.Cmd {
	return exec.Command("bzr", "revno", "--tree")
}

func (repo *BzrRepository) UpdateCommand() *exec.Cmd {
	return exec.Command("bzr", "pull")
}

func (repo *BzrRepository) StatusCommand() *exec.Cmd {
	return exec.Command("bzr", "status")
}

func (repo *BzrRepository) LogCommand(from, to string) *exec.Cmd {
	return exec.Command("bzr", "log", "--line", "-r", fmt.Sprintf("%s..%s", from, to))
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
