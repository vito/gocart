package gocart

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	. "launchpad.net/gocheck"

	"github.com/remogatto/prettytest"
)

type InstallSuite struct {
	prettytest.Suite

	Install *exec.Cmd
	GOPATH  string
}

func TestRunner(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		new(prettytest.TDDFormatter),
		new(InstallSuite),
	)
}

func (s *InstallSuite) BeforeEach() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory := path.Dir(currentFile)

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

	fakeRepoPath, err := filepath.Abs(path.Join(currentDirectory, "fake_install_repo"))
	s.Nil(err)

	gopath, err := ioutil.TempDir(os.TempDir(), "fake_install_repo_GOPATH")
	s.Nil(err)

	s.Install.Env = []string{"GOPATH=" + gopath}

	s.Install.Dir = fakeRepoPath

	s.GOPATH = gopath
}

func (s *InstallSuite) TestInstallWithoutLockFile() {
	out, err := s.Install.CombinedOutput()
	if err != nil {
		println(string(out))
	}

	s.Nil(err)

	s.Check(string(out), Matches, "Installing dependencies...\n")

	s.Path(path.Join(s.GOPATH, "src", "github.com", "xoebus", "gocart"))
}

func (s *InstallSuite) TestInstallWithLockFile() {
}
