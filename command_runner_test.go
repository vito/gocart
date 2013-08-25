package gocart

import (
	"bytes"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shell Command Runner", func() {
	var runner *ShellCommandRunner
	var buffer *bytes.Buffer

	BeforeEach(func() {
		runner = &ShellCommandRunner{}
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
			Expect(err).To(HaveOccured())
		})
	})
})
