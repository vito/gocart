package command_runner

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type CommandRunner interface {
	Run(*exec.Cmd) error
}

type RealCommandRunner struct {
	debug bool
}

type CommandFailedError struct {
	OriginalError error

	Command *exec.Cmd
	Output  []byte
}

func (e CommandFailedError) Error() string {
	return fmt.Sprintf(
		"command failed: %s\nerror: %s\noutput:\n%s",
		prettyCommand(e.Command),
		e.OriginalError.Error(),
		e.Output,
	)
}

func New(debug bool) *RealCommandRunner {
	return &RealCommandRunner{debug}
}

func (r *RealCommandRunner) Run(cmd *exec.Cmd) error {
	if r.debug {
		log.Printf("\x1b[40;36mexecuting: %s\x1b[0m\n", prettyCommand(cmd))
		r.tee(cmd, os.Stderr)
	}

	output := new(bytes.Buffer)

	r.tee(cmd, output)

	err := cmd.Run()

	if r.debug {
		if err != nil {
			log.Printf("\x1b[40;31mcommand failed (%s): %s\x1b[0m\n", prettyCommand(cmd), err)
		} else {
			log.Printf("\x1b[40;32mcommand succeeded (%s)\x1b[0m\n", prettyCommand(cmd))
		}
	}

	if err != nil {
		return CommandFailedError{
			OriginalError: err,

			Command: cmd,
			Output:  output.Bytes(),
		}
	}

	return nil
}
func (r *RealCommandRunner) tee(cmd *exec.Cmd, dst io.Writer) {
	if cmd.Stderr == nil {
		cmd.Stderr = dst
	} else if cmd.Stderr != nil {
		cmd.Stderr = io.MultiWriter(cmd.Stderr, dst)
	}

	if cmd.Stdout == nil {
		cmd.Stdout = dst
	} else if cmd.Stdout != nil {
		cmd.Stdout = io.MultiWriter(cmd.Stdout, dst)
	}
}

func prettyCommand(cmd *exec.Cmd) string {
	return fmt.Sprintf("%v %s %v", cmd.Env, cmd.Path, cmd.Args)
}
