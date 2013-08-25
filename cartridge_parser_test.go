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

	Context("with a well-formed cartridge", func() {
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

	Context("when the cartridge has a mal-formed line", func() {
		BeforeEach(func() {
			cartridge = bytes.NewBuffer([]byte(`foo.com/bar`))

			_, err = ParseDependencies(cartridge)
		})

		It("returns an error", func() {
			Expect(err.Error()).To(Equal("malformed line: foo.com/bar"))
		})
	})
})
