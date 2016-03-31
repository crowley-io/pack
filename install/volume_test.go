package install_test

import (
	"fmt"
	"github.com/crowley-io/pack/configuration"
	"path"

	. "github.com/crowley-io/pack/install"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Volumes", func() {
	Describe("with go path library", func() {
		It("should remove double path sepator with Clean", func() {
			Expect(path.Clean("/media//archive.tar.gz")).To(Equal("/media/archive.tar.gz"))
		})
		It("should remove trailing path sepator with Clean", func() {
			Expect(path.Clean("/media/")).To(Equal("/media"))
		})
		It("should return the parent directory with Dir", func() {
			Expect(path.Dir(path.Clean("/media/archive.tar.gz"))).To(Equal("/media"))
		})
	})
	Describe("with internal parser", func() {

		var (
			c   *configuration.Configuration
			e   []string
			v   []string
			err error
		)

		AssertFailedParsing := func() {
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("should have an empty list of volumes", func() {
				Expect(v).To(BeEmpty())
			})
		}

		AssertSuccessfulParsing := func() {
			It("should succeed", func() {
				Expect(err).To(Succeed())
			})
			It("should match the expected list of volumes", func() {
				Expect(v).To(ConsistOf(e))
			})
		}

		JustBeforeEach(func() {
			v, err = GetVolumes(c)
		})
		Context("with an absolute external path", func() {
			Context("without an explicit mode", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfiguration("/home/user/.npm:/root/.npm")
				})
				AssertSuccessfulParsing()
			})
			Context("with an explicit mode", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfiguration("/home/user/.npm:/root/.npm:rw")
				})
				AssertSuccessfulParsing()
			})
			Context("without an external path", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfiguration("/root/.npm")
				})
				AssertSuccessfulParsing()
			})
			Context("without an external path and with a mode", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfiguration("/root/.npm:ro")
				})
				AssertFailedParsing()
				It("should parse mode as an internal path", func() {
					Expect(err.Error()).To(Equal("internal path 'ro' is not an absolute path"))
				})
			})
			Context("with a syntax error", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfiguration("/home/user/.npm:/root/.npm:rw:3")
				})
				AssertFailedParsing()
			})
			Context("with an internal relative path", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfiguration("/home/user/.npm:./foo:rw")
				})
				AssertFailedParsing()
			})
		})
		Context("with a relative external path", func() {
			Context("with home to resolve", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfigurationWithResolve(
						"~/.npm:/root/.npm", "%s/.npm:/root/.npm", home,
					)
				})
				AssertSuccessfulParsing()
			})
			Context("with home (with a trailing slash) to resolve", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfigurationWithResolve(
						"~/:/root/", "%s:/root", home,
					)
				})
				AssertSuccessfulParsing()
			})
			Context("with a folder in the working to resolve", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfigurationWithResolve(
						"./bin/:/usr/local/bin/:ro", "%s/bin:/usr/local/bin:ro", pwd,
					)
				})
				AssertSuccessfulParsing()
			})
			Context("with the working directory to resolve", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfigurationWithResolve(
						"./:/var/www", "%s:/var/www", pwd,
					)
				})
				AssertSuccessfulParsing()
			})
			Context("with a folder in working directory (without any prefix path) to resolve", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfigurationWithResolve(
						"bin/app:/var/www:rw", "%s/bin/app:/var/www:rw", pwd,
					)
				})
				AssertSuccessfulParsing()
			})
			Context("with a folder starting with a tilde in working directory", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfigurationWithResolve(
						"~bin/app:/var/www:rw", "%s/~bin/app:/var/www:rw", pwd,
					)
				})
				AssertSuccessfulParsing()
			})
			Context("with a complex path in working directory", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfigurationWithResolve(
						"../bin/../app:/var/", "%s/app:/var", path.Dir(pwd),
					)
				})
				AssertSuccessfulParsing()
			})
			Context("with a multiple level parent directory", func() {
				BeforeEach(func() {
					c, e = getSimpleVolumeConfigurationWithResolve(
						"../../:/var/", "%s:/var", path.Dir(path.Dir(pwd)),
					)
				})
				AssertSuccessfulParsing()
			})
		})
		Context("with an invalid configuration", func() {
			BeforeEach(func() {
				c, e = getSimpleVolumeConfiguration("/media/foo/bar:/var/lib/foo/bar")
				c.Install.Path = "~/foo:/usr/local/app"
			})
			AssertFailedParsing()
			It("should mention that the format is incorrect", func() {
				Expect(err.Error()).To(ContainSubstring("has incorrect format, should be external:internal[:mode]"))
				Expect(v).To(BeEmpty())
			})
		})
	})
})

func getSimpleVolumeConfiguration(path string) (*configuration.Configuration, []string) {
	return getSimpleVolumeConfigurationWithResolve(path, "%s", path)
}

func getSimpleVolumeConfigurationWithResolve(path, format, value string) (*configuration.Configuration, []string) {

	c := &configuration.Configuration{
		Install: configuration.Install{
			Path:    "/media",
			Volumes: []string{path},
		},
	}

	v := []string{
		fmt.Sprintf("%s:/media:rw", pwd),
		fmt.Sprintf(format, value),
	}

	return c, v
}
