package main_test

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vito/cmdtest"
	. "github.com/vito/cmdtest/matchers"

	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/set"
)

var currentDirectory,
	fakeDiverseRepoPath,
	fakeLockedGitRepoPath,
	fakeGitRepoPath, fakeGitRepoWithRevisionPath,
	fakeLockedGitRepoWithNewDepPath,
	fakeLockedGitRepoWithRemovedDepPath,
	fakeHgRepoPath, fakeHgRepoWithRevisionPath,
	fakeBzrRepoPath, fakeBzrRepoWithRevisionPath,
	fakeUnlockedRepoWithRecursiveDependencies,
	fakeUnlockedRepoWithRecursiveConflictingDependencies,
	fakeUnlockedRepoWithTestDependencies string

var _ = BeforeEach(func() {
	var err error

	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory = path.Dir(currentFile)

	destTmpDir, err := ioutil.TempDir(os.TempDir(), "wtf")
	Ω(err).ShouldNot(HaveOccurred())

	err = walkTheDinosaur(currentDirectory, destTmpDir)
	Ω(err).ShouldNot(HaveOccurred())

	gocartDir := path.Join(destTmpDir, "gocart")

	fakeDiverseRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_diverse_repo"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeLockedGitRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo_locked"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeLockedGitRepoWithNewDepPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo_locked_with_new_dep"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeLockedGitRepoWithRemovedDepPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo_locked_with_removed_dep"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeGitRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeGitRepoWithRevisionPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo_with_revision"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeBzrRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_bzr_repo"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeBzrRepoWithRevisionPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_bzr_repo_with_revision"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeHgRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_hg_repo"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeHgRepoWithRevisionPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_hg_repo_with_revision"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeUnlockedRepoWithRecursiveDependencies, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_recursive_repo"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeUnlockedRepoWithRecursiveConflictingDependencies, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_recursive_repo_with_conflicting_dependencies"),
	)
	Ω(err).ShouldNot(HaveOccurred())

	fakeUnlockedRepoWithTestDependencies, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_repo_with_test_tag"),
	)
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = Describe("install", func() {
	gocartPath, err := cmdtest.Build("github.com/vito/gocart")
	if err != nil {
		panic(err)
	}

	// TODO: move to cmdtest
	err = os.Chmod(gocartPath, 0755)
	if err != nil {
		panic(err)
	}

	var installCmd *exec.Cmd
	var gopath string

	teeToStdout := func(w io.Writer) io.Writer {
		return io.MultiWriter(w, os.Stdout)
	}

	installing := func() *cmdtest.Session {
		sess, err := cmdtest.StartWrapped(installCmd, teeToStdout, teeToStdout)
		Expect(err).ToNot(HaveOccurred())

		return sess
	}

	install := func() {
		sess := installing()

		Expect(sess).To(Say("OK"))
		Expect(sess).To(ExitWith(0))
	}

	BeforeEach(func() {
		installCmd = exec.Command(gocartPath, "install")

		var err error

		gopath, err = ioutil.TempDir(os.TempDir(), "fake_repo_GOPATH")
		Expect(err).ToNot(HaveOccurred())

		installCmd.Env = []string{
			"GOPATH=" + gopath,
			"GOROOT=" + os.Getenv("GOROOT"),
			"PATH=" + os.Getenv("PATH"),
			"PYTHONPATH=" + os.Getenv("PYTHONPATH"), // bzr
		}
	})

	Context("without a Cartridge.lock", func() {
		It("downloads git dependencies", func() {
			installCmd.Dir = fakeGitRepoPath

			dependencyPath := path.Join(gopath, "src", "github.com", "vito", "gocart")

			Expect(listing(dependencyPath)).ToNot(ExitWith(0))

			install()

			Expect(listing(dependencyPath)).To(ExitWith(0))
		})

		It("checks out git repos to the given ref", func() {
			installCmd.Dir = fakeGitRepoWithRevisionPath

			dependencyPath := path.Join(gopath, "src", "github.com", "vito", "gocart")

			install()

			Expect(gitRevision(dependencyPath, "HEAD")).To(Say("7c9d1a95d4b7979bc4180d4cb4aebfc036f276de"))
		})

		It("downloads bzr dependencies", func() {
			installCmd.Dir = fakeBzrRepoPath

			dependencyPath := path.Join(gopath, "src", "launchpad.net", "gocheck")

			Expect(listing(dependencyPath)).ToNot(ExitWith(0))

			install()

			Expect(listing(dependencyPath)).To(ExitWith(0))
		})

		It("checks out bzr repos to the given ref", func() {
			installCmd.Dir = fakeBzrRepoWithRevisionPath

			dependencyPath := path.Join(gopath, "src", "launchpad.net", "gocheck")

			install()

			Expect(bzrRevision(dependencyPath)).To(Say("1"))
		})

		It("downloads hg dependencies", func() {
			installCmd.Dir = fakeHgRepoPath

			dependencyPath := path.Join(gopath, "src", "code.google.com", "p", "go.crypto", "ssh")

			Expect(listing(dependencyPath)).ToNot(ExitWith(0))

			install()

			Expect(listing(dependencyPath)).To(ExitWith(0))
		})

		It("checks out hg repos to the given ref", func() {
			installCmd.Dir = fakeHgRepoWithRevisionPath

			dependencyPath := path.Join(gopath, "src", "code.google.com", "p", "go.crypto", "ssh")

			install()

			Expect(hgRevision(dependencyPath)).To(Say("1e7a3e301825"))
		})

		It("generates a Cartridge.lock file", func() {
			installCmd.Dir = fakeDiverseRepoPath

			install()

			set, err := set.LoadFrom(installCmd.Dir)
			Expect(err).ToNot(HaveOccurred())

			dependency0Version := currentGitRevision(path.Join(gopath, "src", "github.com", "vito", "gocart"))
			dependency1Version := currentHgRevision(path.Join(gopath, "src", "code.google.com", "p", "go.crypto", "ssh"))

			Expect(set.Dependencies).To(HaveLen(2))
			Expect(set.Dependencies).To(Equal([]dependency.Dependency{
				{
					Path:    "github.com/vito/gocart",
					Version: dependency0Version,
				},
				{
					Path:    "code.google.com/p/go.crypto/ssh",
					Version: dependency1Version,
				},
			}))
		})
	})

	Context("with a Cartridge.lock", func() {
		It("installs the locked-down versions", func() {
			installCmd.Dir = fakeLockedGitRepoPath

			dependencyPath := path.Join(gopath, "src", "github.com", "vito", "gocart")

			install()

			Expect(gitRevision(dependencyPath, "HEAD")).To(Say("7c9d1a95d4b7979bc4180d4cb4aebfc036f276de"))
		})

		Context("when there are new dependencies in the Cartridge", func() {
			It("locks them down", func() {
				installCmd.Dir = fakeLockedGitRepoWithNewDepPath

				dependencyPath := path.Join(gopath, "src", "github.com", "onsi", "ginkgo")

				install()

				Expect(gitRevision(dependencyPath, "HEAD")).To(Say("334e06b31ec28f58e7f2df287d2bcf68f59af2b3"))

				set, err := set.LoadFrom(installCmd.Dir)
				Ω(err).ShouldNot(HaveOccurred())

				Expect(set.Dependencies).To(HaveLen(2))
				Expect(set.Dependencies).To(ContainElement(
					dependency.Dependency{
						Path: "github.com/onsi/ginkgo",
						// cfd6b07da4e69326bbd6b7057bbb4693cb78577b~1
						Version: "334e06b31ec28f58e7f2df287d2bcf68f59af2b3",
					},
				))
			})
		})

		Context("when there are dependencies removed from the Cartridge", func() {
			It("removes them from the lock", func() {
				installCmd.Dir = fakeLockedGitRepoWithRemovedDepPath

				install()

				set, err := set.LoadFrom(installCmd.Dir)
				Expect(err).ToNot(HaveOccurred())

				Expect(set.Dependencies).To(Equal([]dependency.Dependency{
					{
						Path:    "github.com/vito/gocart",
						Version: "7c9d1a95d4b7979bc4180d4cb4aebfc036f276de",
					},
				}))
			})
		})
	})

	Context("with -r", func() {
		BeforeEach(func() {
			installCmd.Args = append([]string{installCmd.Args[0], "-r"}, installCmd.Args[1:]...)
		})

		It("recursively installs dependencies", func() {
			installCmd.Dir = fakeUnlockedRepoWithRecursiveDependencies

			sess := installing()
			Expect(sess).To(Say("github.com/vito/gocart"))
			Expect(sess).To(Say("39ada75afb9b654b4621822e707258812bff34ac"))
			Expect(sess).To(Say("github.com/vito/cmdtest"))
			Expect(sess).To(Say("4b86f8c2259c55e86e4b971cd7dc5dfb03e41b80"))
			Expect(sess).ToNot(Say("github.com/onsi/(ginkgo|gomega)"))
			Expect(sess).To(Say("OK"))
			Expect(sess).To(ExitWith(0))

			Expect(gitRevision(
				path.Join(gopath, "src", "github.com", "vito", "cmdtest"),
				"HEAD",
			)).To(Say("4b86f8c2259c55e86e4b971cd7dc5dfb03e41b80"))
		})

		It("checks for conflicting dependencies (different SHAs)", func() {
			installCmd.Dir = fakeUnlockedRepoWithRecursiveConflictingDependencies

			sess := installing()
			Expect(sess).To(SayError("version conflict"))
			Expect(sess).ToNot(ExitWith(0))
		})
	})

	Context("with -x", func() {
		BeforeEach(func() {
			installCmd.Args = append(
				[]string{installCmd.Args[0], "-x", "test"},
				installCmd.Args[1:]...,
			)
		})

		It("excludes dependencies matching any of the given tags", func() {
			installCmd.Dir = fakeUnlockedRepoWithTestDependencies

			sess := installing()
			Expect(sess).To(Say("github.com/vito/gocart"))
			Expect(sess).To(Say("origin/master"))
			Expect(sess).ToNot(Say("github.com/onsi/(ginkgo|gomega)"))
			Expect(sess).To(ExitWith(0))
		})
	})
})

var _ = Describe("check", func() {
	gocartPath, err := cmdtest.Build("github.com/vito/gocart")
	if err != nil {
		panic(err)
	}

	// TODO: move to cmdtest
	err = os.Chmod(gocartPath, 0755)
	if err != nil {
		panic(err)
	}

	var installCmd *exec.Cmd
	var checkCmd *exec.Cmd
	var gopath string

	teeToStdout := func(w io.Writer) io.Writer {
		return io.MultiWriter(w, os.Stdout)
	}

	installing := func() *cmdtest.Session {
		sess, err := cmdtest.StartWrapped(installCmd, teeToStdout, teeToStdout)
		Expect(err).ToNot(HaveOccurred())

		return sess
	}

	checking := func() *cmdtest.Session {
		sess, err := cmdtest.StartWrapped(checkCmd, teeToStdout, teeToStdout)
		Expect(err).ToNot(HaveOccurred())

		return sess
	}

	install := func() {
		sess := installing()

		Expect(sess).To(Say("OK"))
		Expect(sess).To(ExitWith(0))
	}

	BeforeEach(func() {
		installCmd = exec.Command(gocartPath, "install")

		var err error

		gopath, err = ioutil.TempDir(os.TempDir(), "fake_repo_GOPATH")
		Expect(err).ToNot(HaveOccurred())

		installCmd.Env = []string{
			"GOPATH=" + gopath,
			"GOROOT=" + os.Getenv("GOROOT"),
			"PATH=" + os.Getenv("PATH"),
			"PYTHONPATH=" + os.Getenv("PYTHONPATH"), // bzr
		}

		checkCmd = exec.Command(gocartPath, "check")
		checkCmd.Env = installCmd.Env
	})

	itCorrectlyDetectsDirtyDependency := func(repo ...string) {
		Context("when the dependency is in a dirty state", func() {
			var repoPath string
			var repoImportPath string

			BeforeEach(func() {
				repoPath = path.Join(append([]string{gopath, "src"}, repo...)...)
				repoImportPath = path.Join(repo...)

				file, err := os.Create(path.Join(repoPath, "butts"))
				Ω(err).ShouldNot(HaveOccurred())

				file.Close()
			})

			It("reports it as dirty", func() {
				check := checking()
				Expect(check).To(Say(repoImportPath))
				Expect(check).To(ExitWith(1))
			})
		})

		Context("when the repo is so fresh and so clean clean", func() {
			It("exits 0", func() {
				check := checking()
				Expect(check).To(ExitWith(0))
			})
		})
	}

	Context("when the dependencies are not on disk", func() {
		BeforeEach(func() {
			checkCmd.Dir = fakeGitRepoPath

			// don't install
		})

		It("exits 0", func() {
			check := checking()
			Expect(check).To(ExitWith(0))
		})
	})

	Context("when the repo is on a different revision", func() {
		var repoPath string

		BeforeEach(func() {
			installCmd.Dir = fakeGitRepoPath
			checkCmd.Dir = fakeGitRepoPath

			install()

			repoPath = path.Join(gopath, "src", "github.com", "vito", "gocart")

			checkout := exec.Command("git", "checkout", "HEAD~1")
			checkout.Dir = repoPath

			err = checkout.Run()
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("reports it as dirty", func() {
			check := checking()
			Expect(check).To(Say("github.com/vito/gocart"))
			Expect(check).To(Say(".*1.* behind"))
			Expect(check).To(ExitWith(1))
		})
	})

	Context("with recursive dependencies", func() {
		BeforeEach(func() {
			installCmd.Args = append([]string{installCmd.Args[0], "-r"}, installCmd.Args[1:]...)

			installCmd.Dir = fakeUnlockedRepoWithRecursiveDependencies
			checkCmd.Dir = fakeUnlockedRepoWithRecursiveDependencies

			install()
		})

		itCorrectlyDetectsDirtyDependency("github.com", "vito", "cmdtest")
	})

	Context("with a git dependency", func() {
		BeforeEach(func() {
			installCmd.Dir = fakeGitRepoPath
			checkCmd.Dir = fakeGitRepoPath

			install()
		})

		itCorrectlyDetectsDirtyDependency("github.com", "vito", "gocart")
	})

	Context("with a hg dependency", func() {
		BeforeEach(func() {
			installCmd.Dir = fakeHgRepoPath
			checkCmd.Dir = fakeHgRepoPath

			install()
		})

		itCorrectlyDetectsDirtyDependency("code.google.com", "p", "go.crypto", "ssh")
	})

	Context("with a bzr dependency", func() {
		BeforeEach(func() {
			installCmd.Dir = fakeBzrRepoPath
			checkCmd.Dir = fakeBzrRepoPath

			install()
		})

		itCorrectlyDetectsDirtyDependency("launchpad.net", "gocheck")
	})
})

func gitRevision(path, rev string) *cmdtest.Session {
	git := exec.Command("git", "rev-parse", rev)
	git.Dir = path

	sess, err := cmdtest.Start(git)
	Expect(err).ToNot(HaveOccurred())

	return sess
}

func bzrRevision(path string) *cmdtest.Session {
	bzr := exec.Command("bzr", "revno", "--tree")
	bzr.Dir = path

	sess, err := cmdtest.Start(bzr)
	Expect(err).ToNot(HaveOccurred())

	return sess
}

func hgRevision(path string) *cmdtest.Session {
	hg := exec.Command("hg", "id", "-i")
	hg.Dir = path

	sess, err := cmdtest.Start(hg)
	Expect(err).ToNot(HaveOccurred())

	return sess
}

func currentGitRevision(path string) string {
	sess := gitRevision(path, "HEAD")
	Expect(sess).To(ExitWith(0))

	return strings.Trim(string(sess.FullOutput()), "\n")
}

func currentBzrRevision(path string) string {
	sess := bzrRevision(path)
	Expect(sess).To(ExitWith(0))

	return strings.Trim(string(sess.FullOutput()), "\n")
}

func currentHgRevision(path string) string {
	sess := hgRevision(path)
	Expect(sess).To(ExitWith(0))

	return strings.Trim(string(sess.FullOutput()), "\n")
}

func listing(path string) *cmdtest.Session {
	sess, err := cmdtest.Start(exec.Command("ls", path))
	Expect(err).ToNot(HaveOccurred())

	return sess
}

func walkTheDinosaur(src, dest string) error {
	cp := exec.Command("cp", "-a", src, dest)
	cp.Stdout = os.Stdout
	cp.Stderr = os.Stderr
	return cp.Run()
}
