package set

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vito/gocart/dependency"
)

const CartridgeFile = "Cartridge"
const CartridgeLockFile = "Cartridge.lock"

type Set struct {
	Dependencies []dependency.Dependency
}

var NoCartridgeError = fmt.Errorf("no %s file present", CartridgeFile)

type DuplicateDependencyError struct {
	Original  dependency.Dependency
	Duplicate dependency.Dependency
}

func (e DuplicateDependencyError) Error() string {
	return fmt.Sprintf(
		"duplicate dependencies: '%s' and '%s'",
		e.Original,
		e.Duplicate,
	)
}

type MissingVersionError struct {
	Path string
}

func (e MissingVersionError) Error() string {
	return fmt.Sprintf("missing version for '%s'", e.Path)
}

func LoadFrom(dir string) (*Set, error) {
	cartridgeFilePath := filepath.Join(dir, CartridgeFile)
	cartridgeLockFilePath := filepath.Join(dir, CartridgeLockFile)

	if _, err := os.Stat(cartridgeFilePath); os.IsNotExist(err) {
		return nil, NoCartridgeError
	}

	set := &Set{}

	if err := set.readFrom(cartridgeFilePath); err != nil {
		return nil, err
	}

	if _, err := os.Stat(cartridgeLockFilePath); err == nil {
		lockedSet := &Set{}

		if err := lockedSet.readFrom(cartridgeLockFilePath); err != nil {
			return nil, err
		}

		set.merge(lockedSet)
	}

	return set, nil
}

func (s *Set) SaveTo(dir string) error {
	file, err := os.Create(filepath.Join(dir, CartridgeLockFile))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.WriteTo(file)

	return err
}

func (s *Set) WriteTo(out io.Writer) (int64, error) {
	var written int64

	for _, dep := range s.Dependencies {
		n, err := out.Write([]byte(dep.Path + "\t" + dep.Version + "\n"))

		written += int64(n)

		if err != nil {
			return written, err
		}
	}

	return written, nil
}

func (s *Set) MarshalText() ([]byte, error) {
	buf := new(bytes.Buffer)

	if _, err := s.WriteTo(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Set) UnmarshalText(text []byte) error {
	lines := bufio.NewScanner(bytes.NewReader(text))

	for lines.Scan() {
		words := bufio.NewScanner(bytes.NewReader(lines.Bytes()))
		words.Split(bufio.ScanWords)

		count := 0
		dep := dependency.Dependency{}

		for words.Scan() {
			if strings.HasPrefix(words.Text(), "#") {
				break
			}

			if count == 0 {
				dep.Path = words.Text()
			} else if count == 1 {
				if words.Text() == "*" {
					dep.BleedingEdge = true
				} else {
					dep.Version = words.Text()
				}
			} else if count == 2 {
				dep.Tags = strings.Split(words.Text(), ",")
			}

			count++
		}

		if count == 0 {
			// blank line
			continue
		}

		if count == 1 {
			return MissingVersionError{dep.Path}
		}

		// check for dupes
		for _, existing := range s.Dependencies {
			if strings.HasPrefix(dep.Path+"/", existing.Path+"/") {
				return DuplicateDependencyError{existing, dep}
			}

			if strings.HasPrefix(existing.Path+"/", dep.Path+"/") {
				return DuplicateDependencyError{existing, dep}
			}
		}

		s.Dependencies = append(s.Dependencies, dep)
	}

	return nil
}

func (s *Set) Replace(ldep dependency.Dependency) {
	for i, dep := range s.Dependencies {
		if dep.Path == ldep.Path {
			s.Dependencies[i].Version = ldep.Version
		}
	}
}

func (s *Set) readFrom(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return s.UnmarshalText(content)
}

func (s *Set) merge(lock *Set) {
	for _, ldep := range lock.Dependencies {
		s.Replace(ldep)
	}
}
