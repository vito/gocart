package gocart

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
)

var skippableLine *regexp.Regexp = regexp.MustCompile(`^\s*(#.*)?\s*$`)
var dependencyList *regexp.Regexp = regexp.MustCompile(`^\s*([^\s]+)\s+([^\s]+)\s*$`)

func ParseDependencies(reader io.Reader) ([]Dependency, error) {
	bufferedReader := bufio.NewReader(reader)

	dependencies := []Dependency{}

	for eof := false; !eof; {
		line, err := bufferedReader.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			return nil, err
		}

		if skippableLine.MatchString(line) {
			continue
		}

		dependencyLine := dependencyList.FindStringSubmatch(line)
		if dependencyLine == nil {
			message := fmt.Sprintf("malformed line: %s", line)
			return nil, errors.New(message)
		}

		dependency := Dependency{
			Path:    dependencyLine[1],
			Version: dependencyLine[2],
		}

		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}
