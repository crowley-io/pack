package compose_test

import (
	"errors"
	//"fmt"
	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"

	. "github.com/crowley-io/pack/compose"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compose", func() {

	var (
		c   *configuration.Configuration
		ls  = docker.NewLogStream()
		d   *docker.Mock
		err error
	)

	AssertNoDockerClientUsage := func() {
		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
		})
		It("should not use the docker client", func() {
			Expect(d.ImageIDCalled).To(BeFalse())
			Expect(d.BuildCalled).To(BeFalse())
			Expect(d.RemoveImageCalled).To(BeFalse())
		})
	}

	ExpectDockerBuildImage := func() {
		Expect(d.BuildCalled).To(BeTrue())
		Expect(d.BuildOptions.Name).To(Equal("debian:8.2"))
		Expect(d.BuildOptions.Pull).To(BeTrue())
		Expect(d.BuildOptions.NoCache).To(BeFalse())
		Expect(d.BuildOptions.Directory).To(Equal(pwd))
	}

	AssertSuccessfulCompose := func() {
		It("should succeed", func() {
			Expect(err).To(Succeed())
		})
		It("should build a new image", ExpectDockerBuildImage)
	}

	Describe("create a docker image", func() {
		JustBeforeEach(func() {
			err = Compose(d, ls, c)
		})
		Context("when the configuration is valid", func() {
			Context("and no previous image exists", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					d = docker.NewMock()
					d.ImageIDHandler = getImageIDHandler("", "e4785f57b931")
				})
				AssertSuccessfulCompose()
				It("should not remove any previous image", func() {
					Expect(d.RemoveImageCalled).To(BeFalse())
				})
			})
			Context("and a previous image exists", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					d = docker.NewMock()
					d.ImageIDHandler = getImageIDHandler("2dc8f5484fd8", "e4785f57b931")
				})
				AssertSuccessfulCompose()
				It("should remove the previous image", func() {
					Expect(d.RemoveImageCalled).To(BeTrue())
					Expect(d.RemoveImageName).To(Equal("2dc8f5484fd8"))
				})
			})
			Context("and a previous image exists and match the new one", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					d = docker.NewMock()
					d.ImageIDHandler = getImageIDHandler("e4785f57b931", "e4785f57b931")
				})
				AssertSuccessfulCompose()
				It("should not remove any previous image", func() {
					Expect(d.RemoveImageCalled).To(BeFalse())
				})
			})
		})
		Context("when the configuration is invalid", func() {
			Context("because its empty", func() {
				BeforeEach(func() {
					c = &configuration.Configuration{}
					d = docker.NewMock()
				})
				AssertNoDockerClientUsage()
			})
			Context("because its nil", func() {
				BeforeEach(func() {
					c = nil
					d = docker.NewMock()
				})
				AssertNoDockerClientUsage()
			})
		})
		Context("when the docker client have issue(s)", func() {
			Context("because it cannot build an image", func() {

				var (
					e = errors.New("error: invalid image name")
				)

				BeforeEach(func() {
					c = getDefaultConfiguration()
					d = docker.NewMock()
					d.BuildHandler = func(option docker.BuildOptions, stream docker.LogStream) error {
						return e
					}
					d.ImageIDHandler = getImageIDHandler("2dc8f5484fd8", "2dc8f5484fd8")
				})
				It("should return a build error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(e))
				})
				It("should try to build an image", ExpectDockerBuildImage)
				It("should not remove any previous image", func() {
					Expect(d.RemoveImageCalled).To(BeFalse())
				})
			})
			Context("because it cannot build an image, but a new image exists somehow", func() {

				var (
					e = errors.New("error: socket timeout")
				)

				BeforeEach(func() {
					c = getDefaultConfiguration()
					d = docker.NewMock()
					d.BuildHandler = func(option docker.BuildOptions, stream docker.LogStream) error {
						return e
					}
					d.ImageIDHandler = getImageIDHandler("2dc8f5484fd8", "e4785f57b931")
				})
				It("should return a build error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(e))
				})
				It("should try to build an image", ExpectDockerBuildImage)
				It("should remove the previous image", func() {
					Expect(d.RemoveImageCalled).To(BeTrue())
					Expect(d.RemoveImageName).To(Equal("2dc8f5484fd8"))
				})
			})
			Context("because it cannot remove the previous image", func() {

				var (
					e = errors.New("error: socket timeout")
				)

				BeforeEach(func() {
					c = getDefaultConfiguration()
					d = docker.NewMock()
					d.RemoveImageHandler = func(name string) error {
						return e
					}
					d.ImageIDHandler = getImageIDHandler("2dc8f5484fd8", "e4785f57b931")
				})
				It("should return a remove error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(e))
				})
				It("should build an image", ExpectDockerBuildImage)
				It("should remove the previous image", func() {
					Expect(d.RemoveImageCalled).To(BeTrue())
					Expect(d.RemoveImageName).To(Equal("2dc8f5484fd8"))
				})
			})
		})
	})
})

func getDefaultConfiguration() *configuration.Configuration {
	return &configuration.Configuration{
		Install: configuration.Install{
			Command: "make",
			Path:    "/root",
			Image:   "debian",
			Output:  "libapp.so",
		},
		Compose: configuration.Compose{
			Name:    "debian:8.2",
			Pull:    true,
			NoCache: false,
		},
		Publish: configuration.Publish{
			Hostname: "localhost:5000",
		},
	}
}

func getImageIDHandler(old, new string) func(name string) string {
	i := 0
	return func(name string) string {
		if i != 0 {
			return new
		}
		i++
		return old
	}
}
