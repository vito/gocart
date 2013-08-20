package gocart

import (
	"testing"

	"github.com/remogatto/prettytest"
)

type DependencySuite struct {
	prettytest.Suite
}

func TestRunnerDependency(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		new(prettytest.TDDFormatter),
		new(DependencySuite),
	)
}

func (s *DependencySuite) TestString() {
	dependency := Dependency{
		Path:    "github.com/xoebus/kingpin",
		Version: "master",
	}
	s.Equal(dependency.String(), "github.com/xoebus/kingpin\tmaster")
}

func (s *DependencySuite) TestFullPath() {
	dependency := Dependency{
		Path:    "github.com/xoebus/kingpin",
		Version: "master",
	}
	s.Equal(dependency.fullPath("/tmp"), "/tmp/src/github.com/xoebus/kingpin")
}
