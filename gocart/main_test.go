package main_test

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vito/cmdtest"
	. "github.com/vito/cmdtest/matchers"

	"github.com/vito/gocart"
)

var currentDirectory,
	fakeDiverseRepoPath,
	fakeLockedGitRepoPath,
	fakeGitRepoPath, fakeGitRepoWithRevisionPath,
	fakeLockedGitRepoWithNewDepPath,
	fakeLockedGitRepoWithRemovedDepPath,
	fakeHgRepoPath, fakeHgRepoWithRevisionPath,
	fakeBzrRepoPath, fakeBzrRepoWithRevisionPath string

func init() {
	var err error

	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory = path.Dir(path.Dir(currentFile))

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
}

var _ = Describe("install", func() {
	gocartPath, err := cmdtest.Build("github.com/vito/gocart/gocart")
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

			reader := gocart.NewReader(lockFile)

			dependencies, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccurred())

			Expect(dependencies).To(HaveLen(2))

			dependency0Version, err := dependencies[0].CurrentVersion(gopath)
			Expect(err).ToNot(HaveOccurred())

			dependency1Version, err := dependencies[1].CurrentVersion(gopath)
			Expect(err).ToNot(HaveOccurred())

			Expect(dependencies[0]).To(Equal(gocart.Dependency{
				Path:    "github.com/vito/gocart",
				Version: dependency0Version,
			}))

			Expect(dependencies[1]).To(Equal(gocart.Dependency{
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

				reader := gocart.NewReader(lockFile)
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

				reader := gocart.NewReader(lockFile)
				dependencies, err := reader.ReadAll()
				Expect(err).ToNot(HaveOccurred())

				Expect(dependencies).To(Equal([]gocart.Dependency{
					{
						Path:    "github.com/vito/gocart",
						Version: "7c9d1a95d4b7979bc4180d4cb4aebfc036f276de",
					},
				}))
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

func listing(path string) *cmdtest.Session {
	sess, err := cmdtest.Start(exec.Command("ls", path))
	Expect(err).ToNot(HaveOccurred())

	return sess
}

func walkTheDinosaur(src, dest string) error {
	cp := exec.Command("cp", "-r", src, dest)
	return cp.Run()
}
