package dependency_reader_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDependency_reader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependency_reader Suite")
}
