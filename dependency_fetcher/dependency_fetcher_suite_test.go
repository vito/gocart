package dependency_fetcher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDependency_fetcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependency_fetcher Suite")
}
