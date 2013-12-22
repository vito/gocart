package dependency_fetcher_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/dependency_fetcher"
)

var _ = Describe("A Repository", func() {
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
				repo, err := dependency_fetcher.NewRepository(repoPath)
				Expect(err).ToNot(HaveOccurred())

				_, correctType := repo.(*dependency_fetcher.GitRepository)
				Expect(correctType).To(BeTrue())
			})
		})

		Describe("the checkout command", func() {
			It("uses the correct one", func() {
				repo, err := dependency_fetcher.NewRepository(repoPath)
				Expect(err).ToNot(HaveOccurred())

				command := repo.CheckoutCommand("v1.4")

				Expect(command.Args[0]).To(Equal("git"))
				Expect(command.Args[1]).To(Equal("checkout"))
				Expect(command.Args[2]).To(Equal("v1.4"))
			})
		})
	})

	Describe("a hg repository", func() {
		Describe("type identification", func() {
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
				repo, err := dependency_fetcher.NewRepository(repoPath)
				Expect(err).ToNot(HaveOccurred())

				_, correctType := repo.(*dependency_fetcher.HgRepository)
				Expect(correctType).To(BeTrue())
			})

			Describe("the checkout command", func() {
				It("uses the correct one", func() {
					repo, err := dependency_fetcher.NewRepository(repoPath)
					Expect(err).ToNot(HaveOccurred())

					command := repo.CheckoutCommand("v1.12")

					Expect(command.Args[0]).To(Equal("hg"))
					Expect(command.Args[1]).To(Equal("update"))
					Expect(command.Args[2]).To(Equal("-c"))
					Expect(command.Args[3]).To(Equal("v1.12"))
				})
			})
		})
	})

	Describe("a bzr repository", func() {
		Describe("type identification", func() {
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
				repo, err := dependency_fetcher.NewRepository(repoPath)
				Expect(err).ToNot(HaveOccurred())

				_, correctType := repo.(*dependency_fetcher.BzrRepository)
				Expect(correctType).To(BeTrue())
			})

			Describe("the checkout command", func() {
				It("uses the correct one", func() {
					repo, err := dependency_fetcher.NewRepository(repoPath)
					Expect(err).ToNot(HaveOccurred())

					command := repo.CheckoutCommand("353")

					Expect(command.Args[0]).To(Equal("bzr"))
					Expect(command.Args[1]).To(Equal("update"))
					Expect(command.Args[2]).To(Equal("-r"))
					Expect(command.Args[3]).To(Equal("353"))
				})
			})
		})
	})

	Describe("an unknown repository", func() {
		Describe("type identification", func() {
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
				_, err := dependency_fetcher.NewRepository(repoPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown repository type"))
			})
		})
	})

	Describe("type identification when a subdirectory of the repository", func() {
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
			repo, err := dependency_fetcher.NewRepository(subDir)
			Expect(err).ToNot(HaveOccurred())

			_, correctType := repo.(*dependency_fetcher.HgRepository)
			Expect(correctType).To(BeTrue())
		})
	})
})
