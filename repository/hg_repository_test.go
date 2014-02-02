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

var _ = Describe("HgRepository", func() {
	var repoPath string

	var hgRepo *HgRepository
	var runner *fake_command_runner.FakeCommandRunner

	BeforeEach(func() {
		runner = fake_command_runner.New()

		tmpdir, err := ioutil.TempDir(os.TempDir(), "hg_repo")
		Expect(err).ToNot(HaveOccurred())

		repoPath = tmpdir

		os.Mkdir(path.Join(repoPath, ".hg"), 0600)

		repo, err := New(repoPath, runner)
		Expect(err).ToNot(HaveOccurred())

		hgRepo = repo.(*HgRepository)
	})

	AfterEach(func() {
		os.RemoveAll(repoPath)
	})

	Describe("Checkout", func() {
		It("runs hg checkout", func() {
			err := hgRepo.Checkout("some-ref")
			Expect(err).ToNot(HaveOccurred())

			Expect(runner).To(HaveExecutedSerially(
				fake_command_runner.CommandSpec{
					Path: exec.Command("hg").Path,
					Args: []string{"update", "-c", "some-ref"},
					Dir:  repoPath,
				},
			))
		})

		Context("when hg checkout fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("hg").Path,
						Args: []string{"update", "-c", "some-ref"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				err := hgRepo.Checkout("some-ref")
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("CurrentVersion", func() {
		It("runs hg id and returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("hg").Path,
					Args: []string{"id", "-i"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					return nil
				},
			)

			ver, err := hgRepo.CurrentVersion()
			Expect(err).ToNot(HaveOccurred())

			Expect(ver).To(Equal("abc"))
		})

		Context("when hg rev-parse fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("hg").Path,
						Args: []string{"id", "-i"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := hgRepo.CurrentVersion()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("Update", func() {
		It("runs hg pull", func() {
			err := hgRepo.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(runner).To(HaveExecutedSerially(
				fake_command_runner.CommandSpec{
					Path: exec.Command("hg").Path,
					Args: []string{"pull"},
					Dir:  repoPath,
				},
			))
		})

		Context("when hg pull fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("hg").Path,
						Args: []string{"pull"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				err := hgRepo.Update()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("Status", func() {
		It("runs hg status and returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("hg").Path,
					Args: []string{"status"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					cmd.Stderr.Write([]byte("def\n"))
					return nil
				},
			)

			status, err := hgRepo.Status()
			Expect(err).ToNot(HaveOccurred())

			Expect(status).To(Equal("abc\ndef\n"))
		})

		Context("when hg status fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("hg").Path,
						Args: []string{"status"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := hgRepo.Status()
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})

	Describe("Log", func() {
		It("runs hg log with a fancy template returns its output", func() {
			runner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: exec.Command("hg").Path,
					Args: []string{"log", "--template", "{rev}:{node}: {desc|firstline}\n", "-r", "OLD::NEW"},
				}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("abc\n"))
					cmd.Stderr.Write([]byte("def\n"))
					return nil
				},
			)

			log, err := hgRepo.Log("OLD", "NEW")
			Expect(err).ToNot(HaveOccurred())

			Expect(log).To(Equal("abc\ndef\n"))
		})

		Context("when hg log fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				runner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: exec.Command("hg").Path,
						Args: []string{"log", "--template", "{rev}:{node}: {desc|firstline}\n", "-r", "OLD::NEW"},
					}, func(*exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := hgRepo.Log("OLD", "NEW")
				Expect(err).To(HaveOccurred())

				Expect(err).To(Equal(disaster))
			})
		})
	})
})
