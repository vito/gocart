package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/vito/gocart/command_runner"
)

type GitRepository struct {
	path   string
	runner command_runner.CommandRunner
}

func (r *GitRepository) Checkout(version string) error {
	return r.runner.Run(r.gitCmd("checkout", version))
}

func (r *GitRepository) CurrentVersion() (string, error) {
	cmd := r.gitCmd("rev-parse", "HEAD")

	out, err := r.cmdOutput(cmd)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(out, "\n"), nil
}

func (r *GitRepository) Update() error {
	return r.runner.Run(r.gitCmd("fetch"))
}

func (r *GitRepository) Status() (string, error) {
	return r.cmdOutput(r.gitCmd("status", "--porcelain"))
}

func (r *GitRepository) Log(from, to string) (string, error) {
	return r.cmdOutput(r.gitCmd("log", "--oneline", fmt.Sprintf("%s..%s", from, to)))
}

func (r *GitRepository) gitCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.path

	return cmd
}

func (r *GitRepository) cmdOutput(cmd *exec.Cmd) (string, error) {
	buf := new(bytes.Buffer)

	cmd.Stdout = buf
	cmd.Stderr = buf

	err := r.runner.Run(cmd)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
