package gocart_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart"
)

var _ = Describe("Dependency", func() {
	var dependency gocart.Dependency

	BeforeEach(func() {
		dependency = gocart.Dependency{
			Path:    "github.com/xoebus/kingpin",
			Version: "master",
		}
	})

	Describe("Stringer interface", func() {
		It("returns the string as it would appear in a Cartridge", func() {
			Expect(dependency.String()).To(Equal("github.com/xoebus/kingpin\tmaster"))
		})
	})

	Describe("the full path of the dependency", func() {
		It("prepends the passed in root path", func() {
			Expect(dependency.FullPath("/tmp")).To(Equal("/tmp/src/github.com/xoebus/kingpin"))
		})
	})
})
