package repository_test

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/command_runner/fake_command_runner"
	. "github.com/vito/gocart/command_runner/fake_command_runner/matchers"
	. "github.com/vito/gocart/repository"
)

var _ = Describe("GitRepository", func() {
	var repoPath string

	var gitRepo *GitRepository
	var runner *fake_command_runner.FakeCommandRunner

	BeforeEach(func() {
		runner = fake_command_runner.New()

		tmpdir, err := ioutil.TempDir(os.TempDir(), "git_repo")
		Expect(err).ToNot(HaveOccurred())

		repoPath = tmpdir

		os.Mkdir(path.Join(repoPath, ".git"), 0600)

		repo, err := New(repoPath, runner)
		Expect(err).ToNot(HaveOccurred())

		gitRepo = repo.(*GitRepository)
	})

	AfterEach(func() {
		os.RemoveAll(repoPath)
	})

	Describe("Checkout", func() {
		It("runs git checkout", func() {
			err := gitRepo.Checkout("some-ref")
			Expect(err).ToNot(HaveOccurred())

			Expect(runner).To(HaveExecutedSerially(
				fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"checkout", "some-ref"},
					Dir:  repoPath,
				},
			))
		})

		Context("when git checkout fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"checkout", "some-ref"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				err := gitRepo.Checkout("some-ref")
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("CurrentVersion", func() {
		It("runs git rev-parse HEAD and returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"rev-parse", "HEAD"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					return nil
				},
			)

			ver, err := gitRepo.CurrentVersion()
			Expect(err).ToNot(HaveOccurred())

			Expect(ver).To(Equal("abc"))
		})

		Context("when git rev-parse fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"rev-parse", "HEAD"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := gitRepo.CurrentVersion()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("Update", func() {
		It("runs git fetch", func() {
			err := gitRepo.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(runner).To(HaveExecutedSerially(
				fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"fetch"},
					Dir:  repoPath,
				},
			))
		})

		Context("when git fetch fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"fetch"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				err := gitRepo.Update()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("Status", func() {
		It("runs git status --porcelain and returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"status", "--porcelain"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					cmd.Stderr.Write([]byte("def\n"))
					return nil
				},
			)

			status, err := gitRepo.Status()
			Expect(err).ToNot(HaveOccurred())

			Expect(status).To(Equal("abc\ndef\n"))
		})

		Context("when git status fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"status", "--porcelain"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := gitRepo.Status()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("Log", func() {
		It("runs git log --oneline OLD..NEW and returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("git").Path,
					Args: []string{"log", "--oneline", "OLD..NEW"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					cmd.Stderr.Write([]byte("def\n"))
					return nil
				},
			)

			log, err := gitRepo.Log("OLD", "NEW")
			Expect(err).ToNot(HaveOccurred())

			Expect(log).To(Equal("abc\ndef\n"))
		})

		Context("when git log fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("git").Path,
						Args: []string{"log", "--oneline", "OLD..NEW"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := gitRepo.Log("OLD", "NEW")
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})
})
