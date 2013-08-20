package gocart

import (
	"fmt"
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
	fakeGitRepoPath, fakeGitRepoWithRevisionPath,
	fakeBzrRepoPath, fakeBzrRepoWithRevisionPath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory = path.Dir(currentFile)

	var err error

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
}

func (s *InstallSuite) BeforeEach() {
	mainPath, err := filepath.Abs(path.Join(currentDirectory, "gocart/main.go"))
	s.Nil(err)

	mainExecutable, err := ioutil.TempFile(os.TempDir(), "gocart_test_main")
	s.Nil(err)

	install := exec.Command("go", "build", "-o", mainExecutable.Name(), mainPath)
	out, err := install.CombinedOutput()
	if err != nil {
		println(string(out))
	}

	s.Nil(err)

	s.Install = exec.Command(mainExecutable.Name(), "install")

	gopath, err := ioutil.TempDir(os.TempDir(), "fake_git_repo_GOPATH")
	s.Nil(err)

	s.Install.Env = []string{
		"GOPATH=" + gopath,
		"GOROOT=" + os.Getenv("GOROOT"),
		"PATH=" + os.Getenv("PATH"),
	}

	s.GOPATH = gopath
}

func (s *InstallSuite) TestInstallWithoutLockFileDownloadsGitDependencies() {
	s.Install.Dir = fakeGitRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart")

	s.Not(s.Path(dependencyPath))

	out, err := s.Install.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}

	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutGitRevision() {
	s.Install.Dir = fakeGitRepoWithRevisionPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart")

	out, err := s.Install.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}

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

	out, err := s.Install.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}

	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutBzrRevision() {
	s.Install.Dir = fakeBzrRepoWithRevisionPath

	dependencyPath := path.Join(s.GOPATH, "src", "launchpad.net", "gocheck")

	out, err := s.Install.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}

	s.Nil(err)

	s.Equal(s.bzrRevision(dependencyPath), "1")
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
