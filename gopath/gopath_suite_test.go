package gopath_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGopath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gopath Suite")
}
