package set_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/vito/gocart/dependency"
	. "github.com/vito/gocart/set"
)

var _ = Describe("Set", func() {
	set := &Set{
		[]dependency.Dependency{
			{
				Path:    "github.com/vito/gocart",
				Version: "origin/master",
			},
			{
				Path:    "github.com/onsi/ginkgo",
				Version: "origin/blaster",
			},
		},
	}

	setCartridge := `github.com/vito/gocart	origin/master
github.com/onsi/ginkgo	origin/blaster
`

	Describe("MarshalText", func() {
		It("formats the dependencies as a Cartridge file", func() {
			out, err := set.MarshalText()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(string(out)).Should(Equal(setCartridge))
		})
	})

	Describe("UnmarshalText", func() {
		It("parses a Cartridge's dependencies", func() {
			newSet := &Set{}

			err := newSet.UnmarshalText([]byte(setCartridge))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(newSet).Should(Equal(set))
		})

		It("skips blank lines, linewise comments, and trailing comments", func() {
			newSet := &Set{}

			err := newSet.UnmarshalText([]byte(`github.com/vito/gocart	origin/master

# foo bar
github.com/onsi/ginkgo	origin/blaster # fizz buzz

github.com/onsi/gomega	origin/faster # what the heck come after buzz
`))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(newSet).Should(Equal(&Set{
				[]dependency.Dependency{
					{
						Path:    "github.com/vito/gocart",
						Version: "origin/master",
					},
					{
						Path:    "github.com/onsi/ginkgo",
						Version: "origin/blaster",
					},
					{
						Path:    "github.com/onsi/gomega",
						Version: "origin/faster",
					},
				},
			}))
		})

		It("parses an asterisk version as bleeding-edge", func() {
			newSet := &Set{}

			err := newSet.UnmarshalText([]byte("github.com/vito/gocart *"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(newSet.Dependencies).Should(Equal([]dependency.Dependency{
				{
					Path:         "github.com/vito/gocart",
					BleedingEdge: true,
				},
			}))
		})

		It("parses the third field as comma-separated tags", func() {
			newSet := &Set{}

			err := newSet.UnmarshalText([]byte(
				"github.com/vito/gocart origin/master test,development",
			))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(newSet.Dependencies).Should(Equal([]dependency.Dependency{
				{
					Path:    "github.com/vito/gocart",
					Version: "origin/master",
					Tags:    []string{"test", "development"},
				},
			}))
		})

		It("fails if a dependency is missing its version", func() {
			newSet := &Set{}

			err := newSet.UnmarshalText([]byte(`github.com/vito/gocart`))
			Ω(err).Should(Equal(MissingVersionError{"github.com/vito/gocart"}))
		})

		It("fails if there is a duplicate dependency", func() {
			newSet := &Set{}

			err := newSet.UnmarshalText([]byte(`github.com/vito/gocart	origin/master
github.com/onsi/ginkgo	origin/blaster
github.com/onsi/ginkgo	origin/faster
`))
			Ω(err).Should(Equal(DuplicateDependencyError{
				dependency.Dependency{
					Path:    "github.com/onsi/ginkgo",
					Version: "origin/blaster",
				},
				dependency.Dependency{
					Path:    "github.com/onsi/ginkgo",
					Version: "origin/faster",
				},
			}))

			newSet = &Set{}

			err = newSet.UnmarshalText([]byte(`github.com/vito/gocart	origin/master
github.com/onsi/ginkgo/foo	origin/blaster
github.com/onsi/ginkgo	origin/faster
`))
			Ω(err).Should(Equal(DuplicateDependencyError{
				dependency.Dependency{
					Path:    "github.com/onsi/ginkgo/foo",
					Version: "origin/blaster",
				},
				dependency.Dependency{
					Path:    "github.com/onsi/ginkgo",
					Version: "origin/faster",
				},
			}))

			newSet = &Set{}

			err = newSet.UnmarshalText([]byte(`github.com/vito/gocart	origin/master
github.com/onsi/ginkgo	origin/blaster
github.com/onsi/ginkgo/foo	origin/faster
`))
			Ω(err).Should(Equal(DuplicateDependencyError{
				dependency.Dependency{
					Path:    "github.com/onsi/ginkgo",
					Version: "origin/blaster",
				},
				dependency.Dependency{
					Path:    "github.com/onsi/ginkgo/foo",
					Version: "origin/faster",
				},
			}))
		})

		It("does not fail if two dependencies look similar but are actually, like, not the same", func() {
			newSet := &Set{}

			err := newSet.UnmarshalText([]byte(`github.com/vito/gocart	origin/master
github.com/onsi/ginkgobiloba	origin/blaster
github.com/onsi/ginkgo	origin/faster
`))
			Ω(err).ShouldNot(HaveOccurred())

			newSet = &Set{}

			err = newSet.UnmarshalText([]byte(`github.com/vito/gocart	origin/master
github.com/onsi/ginkgo	origin/blaster
github.com/onsi/ginkgobiloba	origin/faster
`))
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("WriteTo", func() {
		It("formats the dependencies as a Cartridge file", func() {
			buf := new(bytes.Buffer)

			n, err := set.WriteTo(buf)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(len(setCartridge))))

			Ω(buf.String()).Should(Equal(setCartridge))
		})
	})

	Describe("SaveTo", func() {
		var projectDir string

		BeforeEach(func() {
			tmpdir, err := ioutil.TempDir(os.TempDir(), "gocart-project")
			Ω(err).ShouldNot(HaveOccurred())

			projectDir = tmpdir
		})

		AfterEach(func() {
			os.RemoveAll(projectDir)
		})

		It("formats the dependencies to Cartridge.lock", func() {
			err := set.SaveTo(projectDir)
			Ω(err).ShouldNot(HaveOccurred())

			bytes, err := ioutil.ReadFile(filepath.Join(projectDir, "Cartridge.lock"))
			Ω(err).ShouldNot(HaveOccurred())

			newSet := &Set{}

			err = newSet.UnmarshalText(bytes)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(newSet).Should(Equal(set))
		})
	})

	Describe("LoadFrom", func() {
		var projectDir string
		var cartridgeFilePath string
		var cartridgeLockFilePath string

		BeforeEach(func() {
			tmpdir, err := ioutil.TempDir(os.TempDir(), "gocart-project")
			Ω(err).ShouldNot(HaveOccurred())

			projectDir = tmpdir

			cartridgeFilePath = filepath.Join(projectDir, CartridgeFile)
			cartridgeLockFilePath = filepath.Join(projectDir, CartridgeLockFile)
		})

		AfterEach(func() {
			os.RemoveAll(projectDir)
		})

		Context("without a Cartridge", func() {
			It("returns nil and a NoCartridgeError", func() {
				set, err := LoadFrom(projectDir)
				Ω(err).Should(Equal(NoCartridgeError))
				Ω(set).Should(BeNil())
			})
		})

		Context("with a malformed Cartridge", func() {
			BeforeEach(func() {
				file, err := os.Create(cartridgeFilePath)
				Ω(err).ShouldNot(HaveOccurred())

				defer file.Close()

				file.Write([]byte("butts\n"))
			})

			It("returns the error", func() {
				set, err := LoadFrom(projectDir)
				Ω(err).Should(HaveOccurred())
				Ω(set).Should(BeNil())
			})
		})

		Context("with a Cartridge", func() {
			BeforeEach(func() {
				file, err := os.Create(cartridgeFilePath)
				Ω(err).ShouldNot(HaveOccurred())

				defer file.Close()

				_, err = set.WriteTo(file)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("returns a set containing its dependencies", func() {
				set, err := LoadFrom(projectDir)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(set).Should(Equal(set))
			})

			Context("and a Cartridge.lock", func() {
				lock := &Set{
					[]dependency.Dependency{
						{
							Path:    "github.com/vito/gocart",
							Version: "some-sha",
						},
					},
				}

				BeforeEach(func() {
					file, err := os.Create(cartridgeLockFilePath)
					Ω(err).ShouldNot(HaveOccurred())

					defer file.Close()

					_, err = lock.WriteTo(file)
					Ω(err).ShouldNot(HaveOccurred())
				})

				It("locks the dependencies down and keeps new dependencies", func() {
					set, err := LoadFrom(projectDir)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(set).Should(Equal(&Set{
						[]dependency.Dependency{
							{
								Path:    "github.com/vito/gocart",
								Version: "some-sha",
							},
							{
								Path:    "github.com/onsi/ginkgo",
								Version: "origin/blaster",
							},
						},
					}))
				})

				Context("with bleeding-edge dependencies in Cartridge", func() {
					BeforeEach(func() {
						file, err := os.Create(cartridgeFilePath)
						Ω(err).ShouldNot(HaveOccurred())

						defer file.Close()

						_, err = file.Write([]byte(`github.com/vito/gocart *
github.com/onsi/ginkgo origin/blaster
`))
						Ω(err).ShouldNot(HaveOccurred())
					})

					It("sets BleedingEdge on the locked dependency", func() {
						set, err := LoadFrom(projectDir)
						Ω(err).ShouldNot(HaveOccurred())

						Ω(set).Should(Equal(&Set{
							[]dependency.Dependency{
								{
									Path:         "github.com/vito/gocart",
									Version:      "some-sha",
									BleedingEdge: true,
								},
								{
									Path:    "github.com/onsi/ginkgo",
									Version: "origin/blaster",
								},
							},
						}))
					})
				})
			})
		})
	})

	Describe("Replace", func() {
		It("replaces an existing dependency", func() {
			set.Replace(dependency.Dependency{
				Path:    "github.com/vito/gocart",
				Version: "some-sha",
			})

			Ω(set.Dependencies).Should(Equal([]dependency.Dependency{
				{
					Path:    "github.com/vito/gocart",
					Version: "some-sha",
				},
				{
					Path:    "github.com/onsi/ginkgo",
					Version: "origin/blaster",
				},
			}))
		})
	})
})
