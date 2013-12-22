package dependency_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDependency(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependency Suite")
}
