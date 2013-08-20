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
	fakeGitRepoPath, fakeGitRepoWithRevisionPath,
	fakeHgRepoPath, fakeHgRepoWithRevisionPath,
	fakeBzrRepoPath, fakeBzrRepoWithRevisionPath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory = path.Dir(currentFile)

	var err error

	fakeDiverseRepoPath, err = filepath.Abs(
		path.Join(currentDirectory, "fixtures", "fake_diverse_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeGitRepoPath, err = filepath.Abs(
		path.Join(currentDirectory, "fixtures", "fake_git_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeGitRepoWithRevisionPath, err = filepath.Abs(
		path.Join(currentDirectory, "fixtures", "fake_git_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}

	fakeBzrRepoPath, err = filepath.Abs(
		path.Join(currentDirectory, "fixtures", "fake_bzr_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeBzrRepoWithRevisionPath, err = filepath.Abs(
		path.Join(currentDirectory, "fixtures", "fake_bzr_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}

	fakeHgRepoPath, err = filepath.Abs(
		path.Join(currentDirectory, "fixtures", "fake_hg_repo"),
	)
	if err != nil {
		panic(err)
	}

	fakeHgRepoWithRevisionPath, err = filepath.Abs(
		path.Join(currentDirectory, "fixtures", "fake_hg_repo_with_revision"),
	)
	if err != nil {
		panic(err)
	}
}

func (s *InstallSuite) BeforeAll() {
	mainPath, err := filepath.Abs(path.Join(currentDirectory, "gocart/main.go"))
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
	defer s.cleanupLock()

	s.Install.Dir = fakeGitRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart")

	s.Not(s.Path(dependencyPath))

	err := s.Install.Run()
	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutGitRevision() {
	defer s.cleanupLock()

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
	defer s.cleanupLock()

	s.Install.Dir = fakeBzrRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "launchpad.net", "gocheck")

	s.Not(s.Path(dependencyPath))

	err := s.Install.Run()
	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutBzrRevision() {
	defer s.cleanupLock()

	s.Install.Dir = fakeBzrRepoWithRevisionPath

	dependencyPath := path.Join(s.GOPATH, "src", "launchpad.net", "gocheck")

	err := s.Install.Run()
	s.Nil(err)

	s.Equal(s.bzrRevision(dependencyPath), "1")
}

func (s *InstallSuite) TestInstallWithoutLockFileDownloadsHgDependencies() {
	defer s.cleanupLock()

	s.Install.Dir = fakeHgRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "code.google.com", "p", "go.crypto", "ssh")

	s.Not(s.Path(dependencyPath))

	err := s.Install.Run()
	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutHgRevision() {
	defer s.cleanupLock()

	s.Install.Dir = fakeHgRepoWithRevisionPath

	dependencyPath := path.Join(s.GOPATH, "src", "code.google.com", "p", "go.crypto", "ssh")

	err := s.Install.Run()
	s.Nil(err)

	s.Equal(s.hgRevision(dependencyPath), "1e7a3e301825")
}

func (s *InstallSuite) TestInstallWithoutLockFileGeneratesLockFile() {
	defer s.cleanupLock()

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

func (s *InstallSuite) TestInstallWithLockFile() {
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

func (s *InstallSuite) cleanupLock() {
	lockFile := path.Join(s.Install.Dir, "Cartridge.lock")

	if _, err := os.Stat(lockFile); err == nil {
		os.Remove(lockFile)
	}
}
