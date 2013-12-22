package locker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLock_merger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Locker Suite")
}
