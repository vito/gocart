package gocart

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lockfile Merger", func() {
	var cartridgeDependencies []Dependency
	var lockDependencies []Dependency

	Context("when there are no differences", func() {
		BeforeEach(func() {
			cartridgeDependencies = []Dependency{
				Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []Dependency{
				Dependency{Path: "a", Version: "x"},
			}
		})

		It("returns the same set of dependencies", func() {
			resolved := MergeDependencies(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]Dependency{
				Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when the lock file has a different version than the cartridge", func() {
		BeforeEach(func() {
			cartridgeDependencies = []Dependency{
				Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []Dependency{
				Dependency{Path: "a", Version: "y"},
			}
		})

		It("uses the version from the lock file", func() {
			resolved := MergeDependencies(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]Dependency{
				Dependency{Path: "a", Version: "y"},
			}))
		})
	})

	Context("when a dependency has been removed from the Cartridge but not the lock", func() {
		BeforeEach(func() {
			cartridgeDependencies = []Dependency{
				Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []Dependency{
				Dependency{Path: "a", Version: "x"},
				Dependency{Path: "b", Version: "x"},
			}
		})

		It("does not include the dependency from the lock", func() {
			resolved := MergeDependencies(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]Dependency{
				Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when a dependency has been added to the Cartridge", func() {
		BeforeEach(func() {
			cartridgeDependencies = []Dependency{
				Dependency{Path: "a", Version: "x"},
				Dependency{Path: "b", Version: "x"},
			}

			lockDependencies = []Dependency{
				Dependency{Path: "a", Version: "x"},
			}
		})

		It("includes the dependency", func() {
			resolved := MergeDependencies(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]Dependency{
				Dependency{Path: "a", Version: "x"},
				Dependency{Path: "b", Version: "x"},
			}))
		})
	})
})
