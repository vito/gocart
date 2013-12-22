package gopath_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/gopath"
)

var _ = Describe("GOPATH parsing", func() {
	Context("when the GOPATH is empty", func() {
		It("raises an error", func() {
			_, err := gopath.InstallationDirectory("")
			Expect(err).To(Equal(gopath.GoPathNotSet))
		})
	})

	Context("when the GOPATH has a single element", func() {
		It("returns that single element", func() {
			path, err := gopath.InstallationDirectory("/it/is/a/real/path/honest")
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal("/it/is/a/real/path/honest"))
		})
	})

	Context("when the GOPATH has many elements", func() {
		It("returns the first element", func() {
			path, err := gopath.InstallationDirectory("/this/is/a/real/path/too:/it/is/a/real/path/honest")
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal("/this/is/a/real/path/too"))
		})
	})
})
