package gocart

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parsing cartridges", func() {
	var (
		cartridge    *bytes.Buffer
		dependencies []Dependency
		err          error
	)

	BeforeEach(func() {
		cartridge = bytes.NewBuffer([]byte(`
foo.com/bar master

# i'm a pretty comment
fizz.buzz/foo last:1`))

		dependencies, err = ParseDependencies(cartridge)
	})

	It("parses without error", func() {
		Expect(err).NotTo(HaveOccured())
	})

	It("has the correct number of dependencies", func() {
		Expect(dependencies).To(HaveLen(2))

		Expect(dependencies[0]).To(Equal(Dependency{
			Path:    "foo.com/bar",
			Version: "master",
		}))

		Expect(dependencies[1]).To(Equal(Dependency{
			Path:    "fizz.buzz/foo",
			Version: "last:1",
		}))
	})
})
