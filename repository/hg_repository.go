package repository

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/vito/gocart/command_runner"
)

type HgRepository struct {
	path   string
	runner command_runner.CommandRunner
}

func (r *HgRepository) Checkout(version string) error {
	return r.runner.Run(r.hgCmd("update", "-c", version))
}

func (r *HgRepository) CurrentVersion() (string, error) {
	cmd := r.hgCmd("id", "-i")

	out, err := r.cmdOutput(cmd)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(out, "\n"), nil
}

func (r *HgRepository) Update() error {
	return r.runner.Run(r.hgCmd("pull"))
}

func (r *HgRepository) Status() (string, error) {
	return r.cmdOutput(r.hgCmd("status"))
}

func (r *HgRepository) Log(from, to string) (string, error) {
	return r.cmdOutput(r.hgCmd(
		"log",
		"--template", "{rev}:{node}: {desc|firstline}\n",
		"-r", from+"::"+to,
	))
}

func (r *HgRepository) hgCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("hg", args...)
	cmd.Dir = r.path

	return cmd
}

func (r *HgRepository) cmdOutput(cmd *exec.Cmd) (string, error) {
	buf := new(bytes.Buffer)

	cmd.Stdout = buf
	cmd.Stderr = buf

	err := r.runner.Run(cmd)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
