package dependency

import (
	"fmt"
	"path"
)

type Dependency struct {
	Path         string
	Version      string
	Tags         []string
	BleedingEdge bool
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s\t%s", d.Path, d.Version)
}

func (d Dependency) FullPath(gopath string) string {
	return path.Join(gopath, "src", d.Path)
}
