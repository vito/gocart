package gopath

import (
	"errors"
	"path/filepath"
)

var GoPathNotSet = errors.New("The GOPATH environment variable needs to be set.")

func InstallationDirectory(gopath string) (string, error) {
	if gopath == "" {
		return "", GoPathNotSet
	}

	gopath = filepath.SplitList(gopath)[0]

	return gopath, nil
}
