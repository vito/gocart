package locker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/gocart/dependency"
	"github.com/vito/gocart/locker"
)

var _ = Describe("Generating a set of locked down dependencies", func() {
	var cartridgeDependencies []dependency.Dependency
	var lockDependencies []dependency.Dependency

	Context("when there are no differences", func() {
		BeforeEach(func() {
			cartridgeDependencies = []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}
		})

		It("returns the same set of dependencies", func() {
			resolved := locker.GenerateLock(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when the lock file has a different version than the cartridge", func() {
		BeforeEach(func() {
			cartridgeDependencies = []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "y"},
			}
		})

		It("uses the version from the lock file", func() {
			resolved := locker.GenerateLock(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "y"},
			}))
		})
	})

	Context("when a dependency has been removed from the Cartridge but not the lock", func() {
		BeforeEach(func() {
			cartridgeDependencies = []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies = []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
				dependency.Dependency{Path: "b", Version: "x"},
			}
		})

		It("does not include the dependency from the lock", func() {
			resolved := locker.GenerateLock(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when a dependency has been added to the Cartridge", func() {
		BeforeEach(func() {
			cartridgeDependencies = []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
				dependency.Dependency{Path: "b", Version: "x"},
			}

			lockDependencies = []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}
		})

		It("includes the dependency", func() {
			resolved := locker.GenerateLock(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
				dependency.Dependency{Path: "b", Version: "x"},
			}))
		})
	})
})
