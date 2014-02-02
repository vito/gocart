package fetcher_test

import (
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/command_runner/fake_command_runner"
	dependency_package "github.com/vito/gocart/dependency"
	. "github.com/vito/gocart/fetcher"
	"github.com/vito/gocart/gopath"
)

var _ = Describe("Fetcher", func() {
	var dependency dependency_package.Dependency
	var fetcher *Fetcher
	var runner *fake_command_runner.FakeCommandRunner

	BeforeEach(func() {
		var err error

		dependency = dependency_package.Dependency{
			Path:    "github.com/vito/gocart",
			Version: "v1.2",
		}

		runner = fake_command_runner.New()

		fetcher, err = New(runner)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Fetch", func() {
		It("gets the dependency using go get", func() {
			_, err := fetcher.Fetch(dependency)
			Expect(err).ToNot(HaveOccurred())

			args := strings.Join(runner.ExecutedCommands()[0].Args, " ")
			Expect(runner.ExecutedCommands()[0].Path).To(MatchRegexp(".*/go"))
			Expect(args).To(ContainSubstring("get -d -v " + dependency.Path))
		})

		It("updates the repository", func() {
			_, err := fetcher.Fetch(dependency)
			Expect(err).ToNot(HaveOccurred())

			args := strings.Join(runner.ExecutedCommands()[1].Args, " ")
			Expect(runner.ExecutedCommands()[1].Path).To(MatchRegexp(".*/git"))
			Expect(args).To(ContainSubstring("fetch"))
		})

		It("changes the repository version to be the version specified in the dependency", func() {
			_, err := fetcher.Fetch(dependency)
			Expect(err).ToNot(HaveOccurred())

			gopath, _ := gopath.InstallationDirectory(os.Getenv("GOPATH"))

			args := runner.ExecutedCommands()[2].Args
			Expect(runner.ExecutedCommands()[2].Dir).To(Equal(dependency.FullPath(gopath)))
			Expect(args[0]).To(Equal("git"))
			Expect(args[1]).To(Equal("checkout"))
			Expect(args[2]).To(Equal("v1.2"))
		})

		It("returns the fetched dependency", func() {
			gitPath, err := exec.LookPath("git")
			if err != nil {
				gitPath = "git"
			}

			runner.WhenRunning(fake_command_runner.CommandSpec{
				Path: gitPath,
				Args: []string{"rev-parse", "HEAD"},
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte("some-sha\n"))
				return nil
			})

			dep, err := fetcher.Fetch(dependency)
			Expect(err).ToNot(HaveOccurred())

			Expect(dep.Version).To(Equal("some-sha"))
		})
	})
})
