package command_runner_test

import (
	"bytes"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/command_runner"
)

var _ = Describe("Shell Command Runner", func() {
	var runner *command_runner.ShellCommandRunner
	var buffer *bytes.Buffer

	BeforeEach(func() {
		runner = command_runner.New()
		buffer = &bytes.Buffer{}
	})

	Describe("running commands", func() {
		It("runs commands", func() {
			cmd := exec.Command("echo", "hello")
			cmd.Stdout = buffer
			runner.Run(cmd)
			Expect(buffer.String()).To(Equal("hello\n"))
		})

		It("returns the errors from the command", func() {
			cmd := exec.Command("adsfasdf")
			err := runner.Run(cmd)
			Expect(err).To(HaveOccurred())
		})

		It("returns the command's stdout and stderr in the error", func() {
			cmd := exec.Command("bash", "-c", "echo out; echo err 1>&2; exit 1")
			err := runner.Run(cmd)
			Expect(err.Error()).To(ContainSubstring("exit status 1"))
			Expect(err.Error()).To(ContainSubstring("bash -c"))
			Expect(err.Error()).To(ContainSubstring("\nout\nerr"))
		})
	})
})
