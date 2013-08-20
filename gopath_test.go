package gocart

import (
	"testing"

	"github.com/remogatto/prettytest"
)

type GoPathSuite struct {
	prettytest.Suite
}

func TestRunnerGoPath(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		new(prettytest.TDDFormatter),
		new(GoPathSuite),
	)
}

func (s *GoPathSuite) TestEmptyGoPath() {
	_, err := InstallationDirectory("")
	s.Equal(err, GoPathNotSet)
}

func (s *GoPathSuite) TestGoPathWithOneElement() {
	path, err := InstallationDirectory("/it/is/a/real/path/honest")
	s.Nil(err)
	s.Equal(path, "/it/is/a/real/path/honest")
}

func (s *GoPathSuite) TestGoPathWithManyElements() {
	path, err := InstallationDirectory("/this/is/a/real/path/too:/it/is/a/real/path/honest")
	s.Nil(err)
	s.Equal(path, "/this/is/a/real/path/too")
}
