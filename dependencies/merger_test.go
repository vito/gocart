package dependencies

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/xoebus/gocart/dependency"
)

var _ = Describe("Lockfile Merger", func() {
	Context("when there are no differences", func() {
		It("returns the same set of dependencies", func() {
			cartridgeDependencies := []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies := []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}

			resolved := Merge(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when the lock file has a different version than the cartridge", func() {
		It("uses the version from the lock file", func() {
			cartridgeDependencies := []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies := []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "y"},
			}

			resolved := Merge(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "y"},
			}))
		})
	})

	Context("when a dependency has been removed from the Cartridge but not the lock", func() {
		It("does not include the dependency from the lock", func() {
			cartridgeDependencies := []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}

			lockDependencies := []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
				dependency.Dependency{Path: "b", Version: "x"},
			}

			resolved := Merge(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}))
		})
	})

	Context("when a dependency has been added to the Cartridge", func() {
		It("includes the dependency", func() {
			cartridgeDependencies := []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
				dependency.Dependency{Path: "b", Version: "x"},
			}

			lockDependencies := []dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
			}

			resolved := Merge(cartridgeDependencies, lockDependencies)
			Expect(resolved).To(Equal([]dependency.Dependency{
				dependency.Dependency{Path: "a", Version: "x"},
				dependency.Dependency{Path: "b", Version: "x"},
			}))
		})
	})
})
