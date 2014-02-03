package repository_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/command_runner/fake_command_runner"
	"github.com/vito/gocart/repository"
)

var _ = Describe("A Repository", func() {
	var runner command_runner.CommandRunner

	BeforeEach(func() {
		runner = fake_command_runner.New()
	})

	Describe("a git repository", func() {
		var repoPath string
		var err error

		BeforeEach(func() {
			repoPath, err = ioutil.TempDir(os.TempDir(), "git_repo")
			Expect(err).ToNot(HaveOccurred())

			os.Mkdir(path.Join(repoPath, ".git"), 0600)
		})

		AfterEach(func() {
			os.RemoveAll(repoPath)
		})

		Describe("type identification", func() {
			It("returns that it is a GitRepository", func() {
				repo, err := repository.New(repoPath, runner)
				Expect(err).ToNot(HaveOccurred())

				_, correctType := repo.(*repository.GitRepository)
				Expect(correctType).To(BeTrue())
			})
		})
	})

	Describe("a hg repository", func() {
		var repoPath string
		var err error

		BeforeEach(func() {
			repoPath, err = ioutil.TempDir(os.TempDir(), "hg_repo")
			Expect(err).ToNot(HaveOccurred())

			os.Mkdir(path.Join(repoPath, ".hg"), 0600)
		})

		AfterEach(func() {
			os.RemoveAll(repoPath)
		})

		It("returns that it is a HgRepository", func() {
			repo, err := repository.New(repoPath, runner)
			Expect(err).ToNot(HaveOccurred())

			_, correctType := repo.(*repository.HgRepository)
			Expect(correctType).To(BeTrue())
		})
	})

	Describe("a bzr repository", func() {
		var repoPath string
		var err error

		BeforeEach(func() {
			repoPath, err = ioutil.TempDir(os.TempDir(), "bzr_repo")
			Expect(err).ToNot(HaveOccurred())

			os.Mkdir(path.Join(repoPath, ".bzr"), 0600)
		})

		AfterEach(func() {
			os.RemoveAll(repoPath)
		})

		It("returns that it is a BzrRepository", func() {
			repo, err := repository.New(repoPath, runner)
			Expect(err).ToNot(HaveOccurred())

			_, correctType := repo.(*repository.BzrRepository)
			Expect(correctType).To(BeTrue())
		})
	})

	Describe("an unknown repository", func() {
		var repoPath string
		var err error

		BeforeEach(func() {
			repoPath, err = ioutil.TempDir(os.TempDir(), "unknown_repo")
			Expect(err).ToNot(HaveOccurred())

			os.Mkdir(path.Join(repoPath, ".unknown"), 0600)
		})

		AfterEach(func() {
			os.RemoveAll(repoPath)
		})

		It("returns an error", func() {
			_, err := repository.New(repoPath, runner)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown repository type"))
		})
	})

	Context("when initializing from a subdirectory of the repository", func() {
		var repoPath string
		var err error

		var subDir string

		BeforeEach(func() {
			repoPath, err = ioutil.TempDir(os.TempDir(), "hg_repo")
			Expect(err).ToNot(HaveOccurred())

			subDir = path.Join(repoPath, "a", "subdir")

			os.Mkdir(path.Join(repoPath, ".hg"), 0600)
			os.MkdirAll(subDir, 0600)
		})

		AfterEach(func() {
			os.RemoveAll(repoPath)
		})

		It("recurses up the directory tree until it finds a repo it knows", func() {
			repo, err := repository.New(subDir, runner)
			Expect(err).ToNot(HaveOccurred())

			_, correctType := repo.(*repository.HgRepository)
			Expect(correctType).To(BeTrue())
		})
	})
})
