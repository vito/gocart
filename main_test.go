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
	"github.com/vito/gocart/dependency_reader"
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
	fakeUnlockedRepoWithRecursiveConflictingDependencies string

func init() {
	var err error

	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory = path.Dir(currentFile)

	destTmpDir, err := ioutil.TempDir(os.TempDir(), "wtf")
	if err != nil {
		panic(err)
	}

	err = walkTheDinosaur(currentDirectory, destTmpDir)
	if err != nil {
		panic(err)
	}

	gocartDir := path.Join(destTmpDir, "gocart")

	fakeDiverseRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_diverse_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeLockedGitRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo_locked"),
	)
	if err != nil {
		panic(err)
	}

	fakeLockedGitRepoWithNewDepPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo_locked_with_new_dep"),
	)
	if err != nil {
		panic(err)
	}

	fakeLockedGitRepoWithRemovedDepPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo_locked_with_removed_dep"),
	)
	if err != nil {
		panic(err)
	}

	fakeGitRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeGitRepoWithRevisionPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_git_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}

	fakeBzrRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_bzr_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeBzrRepoWithRevisionPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_bzr_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}

	fakeHgRepoPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_hg_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeHgRepoWithRevisionPath, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_hg_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}

	fakeUnlockedRepoWithRecursiveDependencies, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_recursive_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeUnlockedRepoWithRecursiveConflictingDependencies, err = filepath.Abs(
		path.Join(gocartDir, "fixtures", "fake_recursive_repo_with_conflicting_dependencies"),
	)
	if err != nil {
		panic(err)
	}
}

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

			lockFilePath := path.Join(fakeDiverseRepoPath, "Cartridge.lock")

			lockFile, err := os.Open(lockFilePath)
			Expect(err).ToNot(HaveOccurred())

			reader := dependency_reader.New(lockFile)

			dependencies, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccurred())

			Expect(dependencies).To(HaveLen(2))

			dependency0Version := currentGitRevision(path.Join(gopath, "src", "github.com", "vito", "gocart"))
			dependency1Version := currentHgRevision(path.Join(gopath, "src", "code.google.com", "p", "go.crypto", "ssh"))

			Expect(dependencies[0]).To(Equal(dependency.Dependency{
				Path:    "github.com/vito/gocart",
				Version: dependency0Version,
			}))

			Expect(dependencies[1]).To(Equal(dependency.Dependency{
				Path:    "code.google.com/p/go.crypto/ssh",
				Version: dependency1Version,
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
			It("adds them to Cartridge.lock", func() {
				installCmd.Dir = fakeLockedGitRepoWithNewDepPath

				dependencyPath := path.Join(gopath, "src", "github.com", "onsi", "ginkgo")

				install()

				Expect(gitRevision(dependencyPath, "HEAD")).To(Say("cfd6b07da4e69326bbd6b7057bbb4693cb78577b"))

				lockFilePath := path.Join(installCmd.Dir, "Cartridge.lock")

				lockFile, err := os.Open(lockFilePath)
				Expect(err).ToNot(HaveOccurred())

				reader := dependency_reader.New(lockFile)
				dependencies, err := reader.ReadAll()
				Expect(err).ToNot(HaveOccurred())

				Expect(dependencies).To(HaveLen(2))
			})
		})

		Context("when there are dependencies removed from the Cartridge", func() {
			It("removes them from Cartridge.lock", func() {
				installCmd.Dir = fakeLockedGitRepoWithRemovedDepPath

				install()

				lockFilePath := path.Join(installCmd.Dir, "Cartridge.lock")

				lockFile, err := os.Open(lockFilePath)
				Expect(err).ToNot(HaveOccurred())

				reader := dependency_reader.New(lockFile)
				dependencies, err := reader.ReadAll()
				Expect(err).ToNot(HaveOccurred())

				Expect(dependencies).To(Equal([]dependency.Dependency{
					{
						Path:    "github.com/vito/gocart",
						Version: "7c9d1a95d4b7979bc4180d4cb4aebfc036f276de",
					},
				}))
			})
		})
	})

	Context("when we are recursive when we are recursive when we are", func() {
		BeforeEach(func() {
			installCmd.Args = append([]string{installCmd.Args[0], "-r"}, installCmd.Args[1:]...)
		})

		It("recursively installs dependencies", func() {
			installCmd.Dir = fakeUnlockedRepoWithRecursiveDependencies

			sess := installing()
			Expect(sess).To(Say("github.com/vito/gocart"))
			Expect(sess).To(Say("origin/master"))
			Expect(sess).To(Say("github.com/onsi/ginkgo"))
			Expect(sess).To(Say("9019392d862065b9f2c4461623bd0d1abfd5f435"))
			Expect(sess).To(Say("github.com/onsi/gomega"))
			Expect(sess).To(Say("82aceb33958ceb2758ee32204e02e681d483423c"))
			Expect(sess).To(Say("OK"))
			Expect(sess).To(ExitWith(0))

			Expect(gitRevision(
				path.Join(gopath, "src", "github.com", "onsi", "ginkgo"),
				"HEAD",
			)).To(Say("9019392d862065b9f2c4461623bd0d1abfd5f435"))

			Expect(gitRevision(
				path.Join(gopath, "src", "github.com", "onsi", "gomega"),
				"HEAD",
			)).To(Say("82aceb33958ceb2758ee32204e02e681d483423c"))
		})

		It("checks for conflicting dependencies (different SHAs)", func() {
			installCmd.Dir = fakeUnlockedRepoWithRecursiveConflictingDependencies

			sess := installing()
			Expect(sess).To(SayError("conflict"))
			Expect(sess).ToNot(ExitWith(0))
		})

		Context("with -a (for aggregate)", func() {
			BeforeEach(func() {
				installCmd.Args = append([]string{installCmd.Args[0], "-a"}, installCmd.Args[1:]...)
			})

			It("should collect all the dependencies into a giant .lock file", func() {
				installCmd.Dir = fakeUnlockedRepoWithRecursiveDependencies

				install()

				lockFilePath := path.Join(installCmd.Dir, "Cartridge.lock")
				lockFile, err := os.Open(lockFilePath)
				Expect(err).ToNot(HaveOccurred())

				reader := dependency_reader.New(lockFile)
				dependencies, err := reader.ReadAll()
				Expect(err).ToNot(HaveOccurred())

				Expect(dependencies).To(HaveLen(4))
			})
		})

		Context("with -t (for trickledown)", func() {
			BeforeEach(func() {
				installCmd.Args = append([]string{installCmd.Args[0], "-t"}, installCmd.Args[1:]...)
			})

			It("should enforce dependencies that are defined in the top level cartridge", func() {
				installCmd.Dir = fakeUnlockedRepoWithRecursiveConflictingDependencies

				install()

				lockFilePath := path.Join(gopath, "src", "github.com", "vito", "gocart", "Cartridge.lock")
				lockFile, err := os.Open(lockFilePath)
				Expect(err).ToNot(HaveOccurred())

				reader := dependency_reader.New(lockFile)
				dependencies, err := reader.ReadAll()
				Expect(err).ToNot(HaveOccurred())

				Expect(dependencies).To(ContainElement(
					dependency.Dependency{
						Path:    "github.com/onsi/ginkgo",
						Version: "ed2674365250adb1cae3038ee49a2d8d87a8e4c7",
					},
				))

				Expect(dependencies).To(ContainElement(
					dependency.Dependency{
						Path:    "github.com/onsi/gomega",
						Version: "af00a096625b2a4621cbe96590d2478c906acbd1",
					},
				))

				Expect(dependencies).To(ContainElement(
					dependency.Dependency{
						Path:    "github.com/vito/cmdtest",
						Version: "9193198ec4ce39c99cb25b64f94dea5b5b924e68",
					},
				))

				Expect(gitRevision(
					path.Join(gopath, "src", "github.com", "onsi", "ginkgo"),
					"HEAD",
				)).To(Say("ed2674365250adb1cae3038ee49a2d8d87a8e4c7"))

				Expect(gitRevision(
					path.Join(gopath, "src", "github.com", "onsi", "gomega"),
					"HEAD",
				)).To(Say("af00a096625b2a4621cbe96590d2478c906acbd1"))

				Expect(gitRevision(
					path.Join(gopath, "src", "github.com", "vito", "cmdtest"),
					"HEAD",
				)).To(Say("9193198ec4ce39c99cb25b64f94dea5b5b924e68"))
			})
		})
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
