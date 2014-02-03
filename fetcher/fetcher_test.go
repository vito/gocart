package fetcher_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/command_runner/fake_command_runner"
	. "github.com/vito/gocart/command_runner/fake_command_runner/matchers"
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
		It("go gets, updates, checkouts, and returns the fetched dependency", func() {
			runner.WhenRunning(fake_command_runner.CommandSpec{
				Path: exec.Command("git").Path,
				Args: []string{"rev-parse", "HEAD"},
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte("some-sha\n"))
				return nil
			})

			dep, err := fetcher.Fetch(dependency)
			Expect(err).ToNot(HaveOccurred())

			gopath, _ := gopath.InstallationDirectory(os.Getenv("GOPATH"))

			Ω(runner).Should(HaveExecutedSerially(
				fake_command_runner.CommandSpec{
					Path: exec.Command("go").Path,
					Args: []string{"get", "-d", "-v", dependency.Path},
				},
				fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"fetch"},
					Dir:  dependency.FullPath(gopath),
				},
				fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"checkout", "v1.2"},
					Dir:  dependency.FullPath(gopath),
				},
			))

			Ω(dep.Version).Should(Equal("some-sha"))
		})

		Context("when the repo exists and is already on the correct version", func() {
			BeforeEach(func() {
				gopath, _ := gopath.InstallationDirectory(os.Getenv("GOPATH"))

				err := os.MkdirAll(dependency.FullPath(gopath), 0755)
				Ω(err).ShouldNot(HaveOccurred())

				runner.WhenRunning(fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"rev-parse", "HEAD"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte(dependency.Version + "\n"))
					return nil
				})
			})

			It("skips updating", func() {
				_, err := fetcher.Fetch(dependency)
				Expect(err).ToNot(HaveOccurred())

				Ω(runner).ShouldNot(HaveExecutedSerially(
					fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"fetch"},
					},
				))

				Ω(runner).ShouldNot(HaveExecutedSerially(
					fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"checkout", "v1.2"},
					},
				))
			})
		})

		Context("when the dependency is bleeding-edge", func() {
			Context("and the repository exists", func() {
				BeforeEach(func() {
					gopath, _ := gopath.InstallationDirectory(os.Getenv("GOPATH"))

					err := os.MkdirAll(dependency.FullPath(gopath), 0755)
					Ω(err).ShouldNot(HaveOccurred())
				})

				It("fast-forwards it, but otherwise leaves it alone", func() {
					dependency.BleedingEdge = true

					runner.WhenRunning(fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"rev-parse", "HEAD"},
					}, func(cmd *exec.Cmd) error {
						cmd.Stdout.Write([]byte("some-sha\n"))
						return nil
					})

					dep, err := fetcher.Fetch(dependency)
					Expect(err).ToNot(HaveOccurred())

					Ω(runner).Should(HaveExecutedSerially(
						fake_command_runner.CommandSpec{
							Path: exec.Command("go").Path,
							Args: []string{"get", "-u", "-d", "-v", dependency.Path},
						},
					))

					Ω(runner).ShouldNot(HaveExecutedSerially(
						fake_command_runner.CommandSpec{
							Path: exec.Command("git").Path,
							Args: []string{"fetch"},
						},
					))

					Ω(runner).ShouldNot(HaveExecutedSerially(
						fake_command_runner.CommandSpec{
							Path: exec.Command("git").Path,
							Args: []string{"checkout", "v1.2"},
						},
					))

					Ω(dep.Version).Should(Equal("some-sha"))
				})

				Context("and it's dirty", func() {
					BeforeEach(func() {
						runner.WhenRunning(fake_command_runner.CommandSpec{
							Path: exec.Command("git").Path,
							Args: []string{"status", "--porcelain"},
						}, func(cmd *exec.Cmd) error {
							cmd.Stdout.Write([]byte("A taking\n"))
							cmd.Stdout.Write([]byte("M care\n"))
							cmd.Stdout.Write([]byte("D of\n"))
							cmd.Stdout.Write([]byte("T business\n"))
							return nil
						})
					})

					It("does not fast-forward it", func() {
						dependency.BleedingEdge = true

						runner.WhenRunning(fake_command_runner.CommandSpec{
							Path: exec.Command("git").Path,
							Args: []string{"rev-parse", "HEAD"},
						}, func(cmd *exec.Cmd) error {
							cmd.Stdout.Write([]byte("some-sha\n"))
							return nil
						})

						dep, err := fetcher.Fetch(dependency)
						Expect(err).ToNot(HaveOccurred())

						Ω(runner).Should(HaveExecutedSerially(
							fake_command_runner.CommandSpec{
								Path: exec.Command("go").Path,
								Args: []string{"get", "-d", "-v", dependency.Path},
							},
						))

						Ω(runner).ShouldNot(HaveExecutedSerially(
							fake_command_runner.CommandSpec{
								Path: exec.Command("git").Path,
								Args: []string{"fetch"},
							},
						))

						Ω(runner).ShouldNot(HaveExecutedSerially(
							fake_command_runner.CommandSpec{
								Path: exec.Command("git").Path,
								Args: []string{"checkout", "v1.2"},
							},
						))

						Ω(dep.Version).Should(Equal("some-sha"))
					})
				})
			})

			Context("and the repository does not exist", func() {
				It("go gets it and locks it down", func() {
					runner.WhenRunning(fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"rev-parse", "HEAD"},
					}, func(cmd *exec.Cmd) error {
						cmd.Stdout.Write([]byte("some-sha\n"))
						return nil
					})

					dep, err := fetcher.Fetch(dependency)
					Expect(err).ToNot(HaveOccurred())

					gopath, _ := gopath.InstallationDirectory(os.Getenv("GOPATH"))

					Ω(runner).Should(HaveExecutedSerially(
						fake_command_runner.CommandSpec{
							Path: exec.Command("go").Path,
							Args: []string{"get", "-d", "-v", dependency.Path},
						},
						fake_command_runner.CommandSpec{
							Path: exec.Command("git").Path,
							Args: []string{"fetch"},
							Dir:  dependency.FullPath(gopath),
						},
						fake_command_runner.CommandSpec{
							Path: exec.Command("git").Path,
							Args: []string{"checkout", "v1.2"},
							Dir:  dependency.FullPath(gopath),
						},
					))

					Ω(dep.Version).Should(Equal("some-sha"))
				})
			})
		})

		Context("when a different version has already been fetched", func() {
			It("returns a VersionConflictError", func() {
				count := 0

				runner.WhenRunning(fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"rev-parse", "HEAD"},
				}, func(cmd *exec.Cmd) error {
					if count == 0 {
						// initial check
						cmd.Stdout.Write([]byte("xxx\n"))
					} else if count == 1 {
						// first version
						cmd.Stdout.Write([]byte("some-sha\n"))
					} else if count == 2 {
						// second check
						cmd.Stdout.Write([]byte("some-sha\n"))
					} else {
						// new version
						cmd.Stdout.Write([]byte("some-other-sha\n"))
					}

					count++

					return nil
				})

				_, err := fetcher.Fetch(dependency)
				Ω(err).ShouldNot(HaveOccurred())

				_, err = fetcher.Fetch(dependency)
				Expect(err).To(Equal(VersionConflictError{
					Path:     dependency.Path,
					VersionA: "some-sha",
					VersionB: "some-other-sha",
				}))
			})
		})
	})
})
