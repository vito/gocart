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

var _ = Describe("BzrRepository", func() {
	var repoPath string

	var bzrRepo *BzrRepository
	var runner *fake_command_runner.FakeCommandRunner

	BeforeEach(func() {
		runner = fake_command_runner.New()

		tmpdir, err := ioutil.TempDir(os.TempDir(), "bzr_repo")
		Expect(err).ToNot(HaveOccurred())

		repoPath = tmpdir

		os.Mkdir(path.Join(repoPath, ".bzr"), 0600)

		repo, err := New(repoPath, runner)
		Expect(err).ToNot(HaveOccurred())

		bzrRepo = repo.(*BzrRepository)
	})

	AfterEach(func() {
		os.RemoveAll(repoPath)
	})

	Describe("Checkout", func() {
		It("runs bzr update -r", func() {
			err := bzrRepo.Checkout("some-ref")
			Expect(err).ToNot(HaveOccurred())

			Expect(runner).To(HaveExecutedSerially(
				fake_command_runner.CommandSpec{
					Path: exec.Command("bzr").Path,
					Args: []string{"update", "-r", "some-ref"},
					Dir:  repoPath,
				},
			))
		})

		Context("when bzr update -r fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("bzr").Path,
						Args: []string{"update", "-r", "some-ref"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				err := bzrRepo.Checkout("some-ref")
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("CurrentVersion", func() {
		It("runs bzr revno --tree and returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("bzr").Path,
					Args: []string{"revno", "--tree"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					return nil
				},
			)

			ver, err := bzrRepo.CurrentVersion()
			Expect(err).ToNot(HaveOccurred())

			Expect(ver).To(Equal("abc"))
		})

		Context("when bzr revno --tree fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("bzr").Path,
						Args: []string{"revno", "--tree"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := bzrRepo.CurrentVersion()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("Update", func() {
		It("runs bzr pull", func() {
			err := bzrRepo.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(runner).To(HaveExecutedSerially(
				fake_command_runner.CommandSpec{
					Path: exec.Command("bzr").Path,
					Args: []string{"pull"},
					Dir:  repoPath,
				},
			))
		})

		Context("when bzr pull fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("bzr").Path,
						Args: []string{"pull"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				err := bzrRepo.Update()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("Status", func() {
		It("runs bzr status and returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("bzr").Path,
					Args: []string{"status"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					cmd.Stderr.Write([]byte("def\n"))
					return nil
				},
			)

			status, err := bzrRepo.Status()
			Expect(err).ToNot(HaveOccurred())

			Expect(status).To(Equal("abc\ndef\n"))
		})

		Context("when bzr status fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("bzr").Path,
						Args: []string{"status"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := bzrRepo.Status()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})

		Context("when bzr reports that its working tree is out of date", func() {
			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("bzr").Path,
						Args: []string{"status"},
					}, func(cmd *exec.Cmd) error {
						cmd.Stdout.Write([]byte("working tree is out of date, run 'bzr update'\n"))
						cmd.Stderr.Write([]byte("abc\n"))
						return nil
					},
				)
			})

			It("strips it from the output", func() {
				status, err := bzrRepo.Status()
				Expect(err).ToNot(HaveOccurred())

				Expect(status).To(Equal("abc\n"))
			})
		})
	})

	Describe("Log", func() {
		It("runs bzr log --line -r OLD..NEW and returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("bzr").Path,
					Args: []string{"log", "--line", "-r", "OLD..NEW"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					cmd.Stderr.Write([]byte("def\n"))
					return nil
				},
			)

			log, err := bzrRepo.Log("OLD", "NEW")
			Expect(err).ToNot(HaveOccurred())

			Expect(log).To(Equal("abc\ndef\n"))
		})

		Context("when bzr log fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("bzr").Path,
						Args: []string{"log", "--line", "-r", "OLD..NEW"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := bzrRepo.Log("OLD", "NEW")
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})
})
