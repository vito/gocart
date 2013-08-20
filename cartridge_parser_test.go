package gocart

import (
	"bytes"
	"testing"

	"github.com/remogatto/prettytest"
)

type CartridgeParserSuite struct {
	prettytest.Suite
}

func TestRunnerCartridgeParser(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		new(prettytest.TDDFormatter),
		new(CartridgeParserSuite),
	)
}

func (s *CartridgeParserSuite) TestParsingDependencies() {
	cartridge := bytes.NewBuffer([]byte(`
foo.com/bar master

# i'm a pretty comment
fizz.buzz/foo last:1`))

	dependencies, err := ParseDependencies(cartridge)
	s.Nil(err)

	s.Equal(len(dependencies), 2)

	if len(dependencies) != 2 {
		return
	}

	s.Equal(dependencies[0], Dependency{
		Path:    "foo.com/bar",
		Version: "master",
	})

	s.Equal(dependencies[1], Dependency{
		Path:    "fizz.buzz/foo",
		Version: "last:1",
	})
}
