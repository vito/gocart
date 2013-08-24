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

	BeforeEach(func() {
		dependency = Dependency{
			Path:    "github.com/xoebus/gocart",
			Version: "master",
		}
		runner = &fakes.FakeCommandRunner{}
		fetcher = NewDependencyFetcher(runner)
	})

	Describe("Fetch", func() {
		BeforeEach(func() {
			fetcher.Fetch(dependency)
		})

		It("gets the dependency using go get", func() {
			args := strings.Join(runner.LastCommand.Args, " ")
			Expect(runner.LastCommand.Path).To(MatchRegexp(".*/go"))
			Expect(args).To(ContainSubstring("get -u -d -v " + dependency.Path))
		})

		It("pipes the command's stdout, stdin, and stderr to the user", func() {
			Expect(runner.LastCommand.Stdin).To(Equal(os.Stdin))
			Expect(runner.LastCommand.Stdout).To(Equal(os.Stdout))
			Expect(runner.LastCommand.Stderr).To(Equal(os.Stderr))
		})
	})
})
