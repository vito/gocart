package gocart

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

type Reader struct {
	reader io.Reader
}

var skippableLine *regexp.Regexp = regexp.MustCompile(`^\s*(#.*)?\s*$`)
var dependencyPattern *regexp.Regexp = regexp.MustCompile(`^\s*([^\s]+)\s+([^\s]+)\s*$`)

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		reader: reader,
	}
}

func (reader *Reader) ReadAll() ([]Dependency, error) {
	bufferedReader := bufio.NewReader(reader.reader)

	dependencies := []Dependency{}

	for {
		line, eof, err := readLine(bufferedReader)

		if err != nil {
			return nil, err
		}

		skippable := skippableLine.MatchString(line)

		if skippable {
			if eof {
				break
			}

			continue
		}

		dependency, present, err := parseLine(line)

		if err != nil {
			return dependencies, err
		}

		if !present && eof {
			break
		}

		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

func readLine(reader *bufio.Reader) (line string, eof bool, err error) {
	line, err = reader.ReadString('\n')

	if err == io.EOF {
		return line, true, nil
	} else if err != nil {
		return "", false, nil
	}

	return line, false, nil
}

func parseLine(line string) (dep Dependency, present bool, err error) {
	dependencyLine := dependencyPattern.FindStringSubmatch(line)
	if dependencyLine == nil {
		return Dependency{}, false, errors.New("malformed line")
	}

	return Dependency{
		Path:    dependencyLine[1],
		Version: dependencyLine[2],
	}, true, nil
}
