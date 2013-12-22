package gocart

import (
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cartridge Reader", func() {
	var reader *Reader
	var inputReader io.Reader

	JustBeforeEach(func() {
		reader = NewReader(inputReader)
	})

	Context("with an empty input", func() {
		BeforeEach(func() {
			inputReader = strings.NewReader("")
		})

		It("returns an empty dependency list", func() {
			deps, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccured())
			Expect(deps).To(BeEmpty())
		})
	})

	Context("with a single dependency line", func() {
		BeforeEach(func() {
			inputReader = strings.NewReader(
				`github.com/xoebus/kingpin	master`,
			)
		})

		It("returns the single dependency", func() {
			deps, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccured())

			Expect(deps).To(HaveLen(1))
			Expect(deps[0]).To(Equal(Dependency{
				Path:    "github.com/xoebus/kingpin",
				Version: "master",
			}))
		})
	})

	Context("with a single dependency line and a line break", func() {
		BeforeEach(func() {
			inputReader = strings.NewReader(
				`github.com/xoebus/kingpin	master
				`,
			)
		})

		It("returns the single dependency", func() {
			deps, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccured())

			Expect(deps).To(HaveLen(1))
			Expect(deps[0]).To(Equal(Dependency{
				Path:    "github.com/xoebus/kingpin",
				Version: "master",
			}))
		})
	})

	Context("with a many dependencies", func() {
		BeforeEach(func() {
			inputReader = strings.NewReader(
				`github.com/xoebus/kingpin master
				github.com/vito/gocart v1.0`,
			)
		})

		It("returns the dependencies in order", func() {
			deps, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccured())

			Expect(deps).To(HaveLen(2))
			Expect(deps[0]).To(Equal(Dependency{
				Path:    "github.com/xoebus/kingpin",
				Version: "master",
			}))
			Expect(deps[1]).To(Equal(Dependency{
				Path:    "github.com/vito/gocart",
				Version: "v1.0",
			}))
		})
	})

	Context("with a many dependencies and a linebreak", func() {
		BeforeEach(func() {
			inputReader = strings.NewReader(
				`github.com/xoebus/kingpin master
				github.com/vito/gocart v1.0
				`,
			)
		})

		It("returns the dependencies in order", func() {
			deps, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccured())

			Expect(deps).To(HaveLen(2))
			Expect(deps[0]).To(Equal(Dependency{
				Path:    "github.com/xoebus/kingpin",
				Version: "master",
			}))
			Expect(deps[1]).To(Equal(Dependency{
				Path:    "github.com/vito/gocart",
				Version: "v1.0",
			}))
		})
	})

	Context("with a single comment line", func() {
		BeforeEach(func() {
			inputReader = strings.NewReader(`# this is a comment`)
		})

		It("ignores the comment and returns an empty dependency list", func() {
			deps, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccured())
			Expect(deps).To(BeEmpty())
		})
	})

	Context("with a combination of dependencies and comments", func() {
		BeforeEach(func() {
			inputReader = strings.NewReader(
				`
				# Foos
				github.com/xoebus/kingpin master

				# Bars
				github.com/vito/gocart v1.0
				`,
			)
		})

		It("returns the dependencies in order", func() {
			deps, err := reader.ReadAll()
			Expect(err).ToNot(HaveOccured())

			Expect(deps).To(HaveLen(2))
			Expect(deps[0]).To(Equal(Dependency{
				Path:    "github.com/xoebus/kingpin",
				Version: "master",
			}))
			Expect(deps[1]).To(Equal(Dependency{
				Path:    "github.com/vito/gocart",
				Version: "v1.0",
			}))
		})
	})

	Context("with a badly formatted dependencies", func() {
		BeforeEach(func() {
			inputReader = strings.NewReader(
				`
				# Foos
				github.com/xoebus/kingpin master

				# Bars
				github.com/vito/gocart v1.0

				# Errors
				github.com/xoebus has too many versions
				`,
			)
		})

		It("returns an errors", func() {
			_, err := reader.ReadAll()
			Expect(err).To(HaveOccured())
		})

		It("returns the dependencies in order up until the error", func() {
			deps, _ := reader.ReadAll()
			Expect(deps).To(HaveLen(2))
		})
	})
})
