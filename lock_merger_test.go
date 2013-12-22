package gocart_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart"
)

var _ = Describe("Lockfile Merger", func() {
	var cartridgeDependencies []gocart.Dependency
	var lockDependencies []gocart.Dependency

	Context("when there are no differences", func() {
		BeforeEach(func() {
			cartridgeDependencies = []gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
			}
		})

		It("returns the same set of dependencies", func() {
			resolved := gocart.MergeDependencies(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when the lock file has a different version than the cartridge", func() {
		BeforeEach(func() {
			cartridgeDependencies = []gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "y"},
			}
		})

		It("uses the version from the lock file", func() {
			resolved := gocart.MergeDependencies(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "y"},
			}))
		})
	})

	Context("when a dependency has been removed from the Cartridge but not the lock", func() {
		BeforeEach(func() {
			cartridgeDependencies = []gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
				gocart.Dependency{Path: "b", Version: "x"},
			}
		})

		It("does not include the dependency from the lock", func() {
			resolved := gocart.MergeDependencies(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when a dependency has been added to the Cartridge", func() {
		BeforeEach(func() {
			cartridgeDependencies = []gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
				gocart.Dependency{Path: "b", Version: "x"},
			}

			lockDependencies = []gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
			}
		})

		It("includes the dependency", func() {
			resolved := gocart.MergeDependencies(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]gocart.Dependency{
				gocart.Dependency{Path: "a", Version: "x"},
				gocart.Dependency{Path: "b", Version: "x"},
			}))
		})
	})
})
