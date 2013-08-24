package gocart

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lockfile Merger", func() {
	var merger *LockfileMerger

	BeforeEach(func() {
		merger = &LockfileMerger{}
	})

	Context("when there are no differences", func() {
		It("returns the same set of dependencies", func() {
			cartridgeDependencies := []Dependency{
				Dependency{Path: "a", Version: "x"},
			}

			lockDependencies := []Dependency{
				Dependency{Path: "a", Version: "x"},
			}

			resolved := merger.Merge(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]Dependency{
				Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when the lock file has a different version than the cartridge", func() {
		It("uses the version from the lock file", func() {
			cartridgeDependencies := []Dependency{
				Dependency{Path: "a", Version: "x"},
			}

			lockDependencies := []Dependency{
				Dependency{Path: "a", Version: "y"},
			}

			resolved := merger.Merge(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]Dependency{
				Dependency{Path: "a", Version: "y"},
			}))
		})
	})

	Context("when a dependency has been removed from the Cartridge but not the lock", func() {
		It("does not include the dependency from the lock", func() {
			cartridgeDependencies := []Dependency{
				Dependency{Path: "a", Version: "x"},
			}

			lockDependencies := []Dependency{
				Dependency{Path: "a", Version: "x"},
				Dependency{Path: "b", Version: "x"},
			}

			resolved := merger.Merge(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]Dependency{
				Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when a dependency has been added to the Cartridge", func() {
		It("includes the dependency", func() {
			cartridgeDependencies := []Dependency{
				Dependency{Path: "a", Version: "x"},
				Dependency{Path: "b", Version: "x"},
			}

			lockDependencies := []Dependency{
				Dependency{Path: "a", Version: "x"},
			}

			resolved := merger.Merge(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]Dependency{
				Dependency{Path: "a", Version: "x"},
				Dependency{Path: "b", Version: "x"},
			}))
		})
	})
})
