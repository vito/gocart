package gocart

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/remogatto/prettytest"
)

type InstallSuite struct {
	prettytest.Suite

	mainExecutable *os.File

	Install *exec.Cmd
	GOPATH  string
}

func TestRunnerInstall(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		new(prettytest.TDDFormatter),
		new(InstallSuite),
	)
}

var currentDirectory,
	fakeDiverseRepoPath,
	fakeLockedGitRepoPath,
	fakeGitRepoPath, fakeGitRepoWithRevisionPath,
	fakeLockedGitRepoWithNewDepPath,
	fakeLockedGitRepoWithRemovedDepPath,
	fakeHgRepoPath, fakeHgRepoWithRevisionPath,
	fakeBzrRepoPath, fakeBzrRepoWithRevisionPath string

func walkTheDinosaur(src, dest string) error {
	cp := exec.Command("cp", "-r", src+"/", dest)
	return cp.Run()
}

func init() {
	var err error

	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory = path.Dir(currentFile)

	destTmpDir, err := ioutil.TempDir(os.TempDir(), "wtf")
	if err != nil {
		panic(err)
	}

	walkTheDinosaur(currentDirectory, destTmpDir)

	fakeDiverseRepoPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_diverse_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeLockedGitRepoPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_git_repo_locked"),
	)
	if err != nil {
		panic(err)
	}

	fakeLockedGitRepoWithNewDepPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_git_repo_locked_with_new_dep"),
	)
	if err != nil {
		panic(err)
	}

	fakeLockedGitRepoWithRemovedDepPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_git_repo_locked_with_removed_dep"),
	)
	if err != nil {
		panic(err)
	}

	fakeGitRepoPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_git_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeGitRepoWithRevisionPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_git_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}

	fakeBzrRepoPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_bzr_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeBzrRepoWithRevisionPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_bzr_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}

	fakeHgRepoPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_hg_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeHgRepoWithRevisionPath, err = filepath.Abs(
		path.Join(destTmpDir, "fixtures", "fake_hg_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}
}

func (s *InstallSuite) BeforeAll() {
	mainPath, err := filepath.Abs(path.Join(currentDirectory, "gocart", "main.go"))
	s.Nil(err)

	s.mainExecutable, err = ioutil.TempFile(os.TempDir(), "gocart_test_main")
	s.Nil(err)

	install := exec.Command("go", "build", "-o", s.mainExecutable.Name(), mainPath)
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	install.Stdin = os.Stdin

	err = install.Run()
	s.Nil(err)
}

func (s *InstallSuite) BeforeEach() {
	s.Install = exec.Command(s.mainExecutable.Name(), "install")

	gopath, err := ioutil.TempDir(os.TempDir(), "fake_repo_GOPATH")
	s.Nil(err)

	s.Install.Env = []string{
		"GOPATH=" + gopath,
		"GOROOT=" + os.Getenv("GOROOT"),
		"PATH=" + os.Getenv("PATH"),
	}

	s.Install.Stdout = os.Stdout
	s.Install.Stderr = os.Stderr
	s.Install.Stdin = os.Stdin

	s.GOPATH = gopath
}

func (s *InstallSuite) TestInstallWithoutLockFileDownloadsGitDependencies() {
	s.Install.Dir = fakeGitRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart")

	s.Not(s.Path(dependencyPath))

	err := s.Install.Run()
	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutGitRevision() {
	s.Install.Dir = fakeGitRepoWithRevisionPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart")

	err := s.Install.Run()
	s.Nil(err)

	s.Equal(
		s.gitRevision(dependencyPath, "HEAD"),
		"7c9d1a95d4b7979bc4180d4cb4aebfc036f276de",
	)
}

func (s *InstallSuite) TestInstallWithoutLockFileDownloadsBzrDependencies() {
	s.Install.Dir = fakeBzrRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "launchpad.net", "gocheck")

	s.Not(s.Path(dependencyPath))

	err := s.Install.Run()
	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutBzrRevision() {
	s.Install.Dir = fakeBzrRepoWithRevisionPath

	dependencyPath := path.Join(s.GOPATH, "src", "launchpad.net", "gocheck")

	err := s.Install.Run()
	s.Nil(err)

	s.Equal(s.bzrRevision(dependencyPath), "1")
}

func (s *InstallSuite) TestInstallWithoutLockFileDownloadsHgDependencies() {
	s.Install.Dir = fakeHgRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "code.google.com", "p", "go.crypto", "ssh")

	s.Not(s.Path(dependencyPath))

	err := s.Install.Run()
	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutHgRevision() {
	s.Install.Dir = fakeHgRepoWithRevisionPath

	dependencyPath := path.Join(s.GOPATH, "src", "code.google.com", "p", "go.crypto", "ssh")

	err := s.Install.Run()
	s.Nil(err)

	s.Equal(s.hgRevision(dependencyPath), "1e7a3e301825")
}

func (s *InstallSuite) TestInstallWithoutLockFileGeneratesLockFile() {
	s.Install.Dir = fakeDiverseRepoPath

	err := s.Install.Run()
	s.Nil(err)

	lockFilePath := path.Join(fakeDiverseRepoPath, "Cartridge.lock")
	s.Path(lockFilePath)

	lockFile, err := os.Open(lockFilePath)
	s.Nil(err)

	dependencies, err := ParseDependencies(lockFile)
	s.Nil(err)

	s.Equal(len(dependencies), 2)

	if len(dependencies) != 2 {
		return
	}

	dependency0Version, err := dependencies[0].CurrentVersion(s.GOPATH)
	s.Nil(err)

	dependency1Version, err := dependencies[1].CurrentVersion(s.GOPATH)
	s.Nil(err)

	s.Equal(dependencies[0], Dependency{
		Path:    "github.com/xoebus/gocart",
		Version: dependency0Version,
	})

	s.Equal(dependencies[1], Dependency{
		Path:    "code.google.com/p/go.crypto/ssh",
		Version: dependency1Version,
	})
}

func (s *InstallSuite) TestInstallWithLockFileInstallsLockedVersions() {
	s.Install.Dir = fakeLockedGitRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart")

	err := s.Install.Run()
	s.Nil(err)

	s.Equal(
		s.gitRevision(dependencyPath, "HEAD"),
		"7c9d1a95d4b7979bc4180d4cb4aebfc036f276de",
	)
}

func (s *InstallSuite) TestInstallWithLockFileWithNewDependencies() {
	s.Install.Dir = fakeLockedGitRepoWithNewDepPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "onsi", "ginkgo")

	err := s.Install.Run()
	s.Nil(err)

	s.Equal(
		s.gitRevision(dependencyPath, "HEAD"),
		"cfd6b07da4e69326bbd6b7057bbb4693cb78577b",
	)

	lockFilePath := path.Join(s.Install.Dir, "Cartridge.lock")

	lockFile, err := os.Open(lockFilePath)
	s.Nil(err)

	dependencies, err := ParseDependencies(lockFile)
	s.Nil(err)

	s.Equal(2, len(dependencies))

}

func (s *InstallSuite) TestInstallWithLockFileWithRemovedDependencies() {
	s.Install.Dir = fakeLockedGitRepoWithRemovedDepPath

	err := s.Install.Run()
	s.Nil(err)

	lockFilePath := path.Join(s.Install.Dir, "Cartridge.lock")

	lockFile, err := os.Open(lockFilePath)
	s.Nil(err)

	dependencies, err := ParseDependencies(lockFile)
	s.Nil(err)

	s.Equal(1, len(dependencies))
	s.Equal(Dependency{Path: "github.com/xoebus/gocart", Version: "7c9d1a95d4b7979bc4180d4cb4aebfc036f276de"}, dependencies[0])
}

func (s *InstallSuite) gitRevision(path, rev string) string {
	git := exec.Command("git", "rev-parse", rev)
	git.Dir = path

	out, err := git.CombinedOutput()
	if err != nil {
		s.Error(err)
	}

	return strings.Trim(string(out), "\n")
}

func (s *InstallSuite) bzrRevision(path string) string {
	bzr := exec.Command("bzr", "revno", "--tree")
	bzr.Dir = path

	out, err := bzr.CombinedOutput()
	if err != nil {
		s.Error(err)
	}

	return strings.Trim(string(out), "\n")
}

func (s *InstallSuite) hgRevision(path string) string {
	hg := exec.Command("hg", "id", "-i")
	hg.Dir = path

	out, err := hg.CombinedOutput()
	if err != nil {
		s.Error(err)
	}

	return strings.Trim(string(out), "\n")
}
