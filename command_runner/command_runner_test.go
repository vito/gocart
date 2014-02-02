package command_runner_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/command_runner"
)

var _ = Describe("Running commands", func() {
	It("runs the command and returns nil", func() {
		runner := command_runner.New(false)

		cmd := exec.Command("ls")
		Expect(cmd.ProcessState).To(BeNil())

		err := runner.Run(cmd)
		Expect(err).ToNot(HaveOccurred())

		Expect(cmd.ProcessState).ToNot(BeNil())
	})

	Context("when the command fails", func() {
		It("returns an error containing its output", func() {
			runner := command_runner.New(false)

			err := runner.Run(exec.Command(
				"/bin/bash",
				"-c", "echo hi out; echo hi err >/dev/stderr; exit 42",
			))

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("hi out\n"))
			Expect(err.Error()).To(ContainSubstring("hi err\n"))
			Expect(err.Error()).To(ContainSubstring("exit status 42"))
		})
	})
})
