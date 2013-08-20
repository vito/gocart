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

var currentDirectory, fakeRepoPath, fakeRepoWithRevisionPath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory = path.Dir(currentFile)

	var err error

	fakeRepoPath, err = filepath.Abs(path.Join(currentDirectory, "fake_install_repo"))
	if err != nil {
		panic(err)
	}

	fakeRepoWithRevisionPath, err = filepath.Abs(path.Join(currentDirectory, "fake_install_repo_with_revision"))
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

	gopath, err := ioutil.TempDir(os.TempDir(), "fake_install_repo_GOPATH")
	s.Nil(err)

	s.Install.Env = []string{
		"GOPATH=" + gopath,
		"GOROOT=" + os.Getenv("GOROOT"),
		"PATH=" + os.Getenv("PATH"),
	}

	s.GOPATH = gopath
}

func (s *InstallSuite) TestInstallWithoutLockFileDownloadsDependencies() {
	s.Install.Dir = fakeRepoPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart")

	s.Not(s.Path(dependencyPath))

	out, err := s.Install.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}

	s.Nil(err)

	s.Path(dependencyPath)
}

func (s *InstallSuite) TestInstallWithoutLockFileChecksOutRevision() {
	s.Install.Dir = fakeRepoWithRevisionPath

	dependencyPath := path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart")

	out, err := s.Install.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}

	s.Nil(err)

	s.Equal(
		s.gitRevision(dependencyPath, "HEAD"),
		s.gitRevision(dependencyPath, "7c9d1a95d4b7979bc4180d4cb4aebfc036f276de"),
	)
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
