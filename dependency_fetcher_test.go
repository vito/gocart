package gocart

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xoebus/gocart/fakes"
)

var _ = Describe("Dependency Fetcher", func() {
	var dependency Dependency
	var fetcher *DependencyFetcher
	var runner *fakes.FakeCommandRunner
	var err error

	BeforeEach(func() {
		dependency = Dependency{
			Path:    "github.com/xoebus/gocart",
			Version: "v1.2",
		}
		runner = &fakes.FakeCommandRunner{}
		fetcher, err = NewDependencyFetcher(runner)
		Expect(err).ToNot(HaveOccured())
	})

	Describe("Fetch", func() {
		BeforeEach(func() {
			fetcher.Fetch(dependency)
		})

		It("gets the dependency using go get", func() {
			args := strings.Join(runner.Commands[0].Args, " ")
			Expect(runner.Commands[0].Path).To(MatchRegexp(".*/go"))
			Expect(args).To(ContainSubstring("get -u -d -v " + dependency.Path))
		})

		It("pipes the command's stdout, stdin, and stderr to the user", func() {
			Expect(runner.Commands[0].Stdin).To(Equal(os.Stdin))
			Expect(runner.Commands[0].Stdout).To(Equal(os.Stdout))
			Expect(runner.Commands[0].Stderr).To(Equal(os.Stderr))
		})

		It("changes the repository version to be the version specified in the dependency", func() {
			gopath, _ := InstallationDirectory(os.Getenv("GOPATH"))

			args := runner.Commands[1].Args
			Expect(runner.Commands[1].Dir).To(Equal(dependency.FullPath(gopath)))
			Expect(args[0]).To(Equal("git"))
			Expect(args[1]).To(Equal("checkout"))
			Expect(args[2]).To(Equal("v1.2"))
		})

		It("pipes the command's stdout, stdin, and stderr to the user for the second command", func() {
			Expect(runner.Commands[1].Stdin).To(Equal(os.Stdin))
			Expect(runner.Commands[1].Stdout).To(Equal(os.Stdout))
			Expect(runner.Commands[1].Stderr).To(Equal(os.Stderr))
		})
	})
})
