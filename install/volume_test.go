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
		It("should parse without explicit mode", func() {

			c, e := getSimpleVolumeConfiguration("/home/user/.npm:/root/.npm")
			v, err := GetVolumes(c)

			Expect(err).To(Succeed())
			Expect(v).To(ConsistOf(e))

		})
		It("should parse with a given mode", func() {

			c, e := getSimpleVolumeConfiguration("/home/user/.npm:/root/.npm:rw")
			v, err := GetVolumes(c)

			Expect(err).To(Succeed())
			Expect(v).To(ConsistOf(e))

		})
		It("should parse without an external path", func() {

			c, e := getSimpleVolumeConfiguration("/root/.npm")
			v, err := GetVolumes(c)

			Expect(err).To(Succeed())
			Expect(v).To(ConsistOf(e))

		})
		It("should return an error with a syntax error", func() {

			c, _ := getSimpleVolumeConfiguration("/home/user/.npm:/root/.npm:rw:3")
			v, err := GetVolumes(c)

			Expect(err).To(HaveOccurred())
			Expect(v).To(BeEmpty())

		})
		It("should return an error with an internal relative path", func() {

			c, _ := getSimpleVolumeConfiguration("/home/user/.npm:./foo:rw")
			v, err := GetVolumes(c)

			Expect(err).To(HaveOccurred())
			Expect(v).To(BeEmpty())

		})
		Context("with the resolver", func() {
			It("should resolve home for an external relative path", func() {

				c, e := getSimpleVolumeConfigurationWithResolve(
					"~/.npm:/root/.npm",
					"%s/.npm:/root/.npm",
					home,
				)
				v, err := GetVolumes(c)

				Expect(err).To(Succeed())
				Expect(v).To(ConsistOf(e))

			})
			It("should resolve home for an external relative path with a trailing slash", func() {

				c, e := getSimpleVolumeConfigurationWithResolve(
					"~/:/root/",
					"%s:/root",
					home,
				)
				v, err := GetVolumes(c)

				Expect(err).To(Succeed())
				Expect(v).To(ConsistOf(e))

			})
			It("should resolve a folder in the working directory", func() {

				c, e := getSimpleVolumeConfigurationWithResolve(
					"./bin/:/usr/local/bin/:ro",
					"%s/bin:/usr/local/bin:ro",
					pwd,
				)
				v, err := GetVolumes(c)

				Expect(err).To(Succeed())
				Expect(v).To(ConsistOf(e))

			})
			It("should resolve the working directory", func() {

				c, e := getSimpleVolumeConfigurationWithResolve("./:/var/www",
					"%s:/var/www", pwd)

				v, err := GetVolumes(c)

				Expect(err).To(Succeed())
				Expect(v).To(ConsistOf(e))

			})
			It("should resolve a folder in working directory without any root prefix", func() {

				c, e := getSimpleVolumeConfigurationWithResolve("bin/app:/var/www:rw",
					"%s/bin/app:/var/www:rw", pwd)

				v, err := GetVolumes(c)

				Expect(err).To(Succeed())
				Expect(v).To(ConsistOf(e))

			})
			It("should resolve a folder with a tilde in working directory", func() {

				c, e := getSimpleVolumeConfigurationWithResolve("~bin/app:/var/www:rw",
					"%s/~bin/app:/var/www:rw", pwd)

				v, err := GetVolumes(c)

				Expect(err).To(Succeed())
				Expect(v).To(ConsistOf(e))

			})
			It("should resolve a complex relative path for a folder in working directory", func() {

				c, e := getSimpleVolumeConfigurationWithResolve("../bin/../app:/var/",
					"%s/app:/var", path.Dir(pwd))

				v, err := GetVolumes(c)

				Expect(err).To(Succeed())
				Expect(v).To(ConsistOf(e))

			})
			It("should resolve a multiple level parent directory", func() {

				c, e := getSimpleVolumeConfigurationWithResolve("../../:/var/",
					"%s:/var", path.Dir(path.Dir(pwd)))

				v, err := GetVolumes(c)

				Expect(err).To(Succeed())
				Expect(v).To(ConsistOf(e))

			})
		})
		Context("with an invalid configuration", func() {
			It("should return an error for invalid path", func() {

				c, _ := getSimpleVolumeConfiguration("/media/foo/bar:/var/lib/foo/bar")
				c.Install.Path = "~/foo:/usr/local/app"
				v, err := GetVolumes(c)

				Expect(err).To(HaveOccurred())
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
