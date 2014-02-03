package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/vito/gocart/command_runner"
)

type BzrRepository struct {
	path   string
	runner command_runner.CommandRunner
}

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

func (r *BzrRepository) Checkout(version string) error {
	return r.runner.Run(r.bzrCmd("update", "-r", version))
}

func (r *BzrRepository) CurrentVersion() (string, error) {
	cmd := r.bzrCmd("revno", "--tree")

	out, err := r.cmdOutput(cmd)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(out, "\n"), nil
}

func (r *BzrRepository) Update() error {
	return r.runner.Run(r.bzrCmd("pull"))
}

func (r *BzrRepository) Status() (string, error) {
	out, err := r.cmdOutput(r.bzrCmd("status"))
	if err != nil {
		return out, err
	}

	return strings.Replace(
		out,
		"working tree is out of date, run 'bzr update'\n",
		"",
		1,
	), nil
}

func (r *BzrRepository) Log(from, to string) (string, error) {
	return r.cmdOutput(r.bzrCmd("log", "--line", "-r", fmt.Sprintf("%s..%s", from, to)))
}

func (r *BzrRepository) bzrCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("bzr", args...)
	cmd.Dir = r.path

	return cmd
}

func (r *BzrRepository) cmdOutput(cmd *exec.Cmd) (string, error) {
	buf := new(bytes.Buffer)

	cmd.Stdout = buf
	cmd.Stderr = buf

	err := r.runner.Run(cmd)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
