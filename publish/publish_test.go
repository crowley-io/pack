package publish_test

import (
	"errors"
	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"

	. "github.com/crowley-io/pack/publish"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Publish", func() {

	var (
		c   *configuration.Configuration
		ls  = docker.NewLogStream()
		d   *docker.Mock
		err error
	)

	Describe("push an image into a registy", func() {
		Context("when the configuration is valid", func() {
			BeforeEach(func() {

				c = getDefaultConfiguration()
				d = docker.NewMock()

			})
			JustBeforeEach(func() {
				err = Publish(d, ls, c)
			})
			It("should succeed", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should push an image on the registry", func() {
				Expect(d.TagCalled).To(BeTrue())
				Expect(d.TagOptions.Name).To(Equal("debian:8.2"))
				Expect(d.TagOptions.Tag).To(Equal("8.2"))
				Expect(d.TagOptions.Repository).To(Equal("localhost:5000/debian"))
				Expect(d.PushCalled).To(BeTrue())
				Expect(d.PushOptions.Name).To(Equal("debian:8.2"))
				Expect(d.PushOptions.Tag).To(Equal("8.2"))
				Expect(d.PushOptions.Repository).To(Equal("localhost:5000/debian"))
				Expect(d.PushOptions.Registry).To(Equal("localhost:5000"))
			})
			It("should remove the registry tag", func() {
				Expect(d.RemoveImageCalled).To(BeTrue())
				Expect(d.RemoveImageName).To(Equal("localhost:5000/debian:8.2"))
			})
		})
		Context("when the configuration is invalid", func() {
			Context("because its empty", func() {
				BeforeEach(func() {

					c = &configuration.Configuration{}
					d = docker.NewMock()

				})
				JustBeforeEach(func() {
					err = Publish(d, ls, c)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
				It("should not use the docker client", func() {
					Expect(d.TagCalled).To(BeFalse())
					Expect(d.PushCalled).To(BeFalse())
					Expect(d.RemoveImageCalled).To(BeFalse())
				})
			})
			Context("because the image identifier is invalid", func() {
				BeforeEach(func() {

					c = getDefaultConfiguration()
					d = docker.NewMock()

					c.Publish.Hostname = "ftp://localhost:5000"

				})
				JustBeforeEach(func() {
					err = Publish(d, ls, c)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("is not a valid repository/tag"))
				})
				It("should not use the docker client", func() {
					Expect(d.TagCalled).To(BeFalse())
					Expect(d.PushCalled).To(BeFalse())
					Expect(d.RemoveImageCalled).To(BeFalse())
				})
			})
		})
		Context("when the docker client have issue(s)", func() {
			Context("because it cannot tag the image", func() {

				var (
					e = errors.New("error: invalid tag name")
				)

				BeforeEach(func() {

					c = getDefaultConfiguration()
					d = docker.NewMock()

					d.TagHandler = func(option docker.TagOptions) error {
						return e
					}

				})
				JustBeforeEach(func() {
					err = Publish(d, ls, c)
				})
				It("should return a tag error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(e))
				})
				It("should try to tag the image to the repository", func() {
					Expect(d.TagCalled).To(BeTrue())
				})
				It("should not push the docker image", func() {
					Expect(d.PushCalled).To(BeFalse())
				})
				It("should not remove the registry tag", func() {
					Expect(d.RemoveImageCalled).To(BeFalse())
				})
			})
			Context("because it cannot push the image", func() {

				var (
					e = errors.New("error: cannot connect to remote registry")
				)

				BeforeEach(func() {

					c = getDefaultConfiguration()
					d = docker.NewMock()

					d.PushHandler = func(option docker.PushOptions, stream docker.LogStream) error {
						return e
					}

				})
				JustBeforeEach(func() {
					err = Publish(d, ls, c)
				})
				It("should return a push error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(e))
				})
				It("should try to push an image on the registry", func() {
					Expect(d.TagCalled).To(BeTrue())
					Expect(d.PushCalled).To(BeTrue())
				})
				It("should try to remove the registry tag", func() {
					Expect(d.RemoveImageCalled).To(BeTrue())
				})
			})
			Context("because it cannot push the image and remove the registry tag", func() {

				var (
					e1 = errors.New("error: cannot connect to remote registry")
					e2 = errors.New("error: socket timeout")
				)

				BeforeEach(func() {

					c = getDefaultConfiguration()
					d = docker.NewMock()

					d.PushHandler = func(option docker.PushOptions, stream docker.LogStream) error {
						return e1
					}

					d.RemoveImageHandler = func(name string) error {
						return e2
					}

				})
				JustBeforeEach(func() {
					err = Publish(d, ls, c)
				})
				It("should return a push error and not a remove error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(e1))
				})
				It("should try to push an image on the registry", func() {
					Expect(d.TagCalled).To(BeTrue())
					Expect(d.PushCalled).To(BeTrue())
				})
				It("should try to remove the registry tag", func() {
					Expect(d.RemoveImageCalled).To(BeTrue())
				})
			})
			Context("because it cannot remove the registry tag", func() {

				var (
					e = errors.New("error: socket timeout")
				)

				BeforeEach(func() {

					c = getDefaultConfiguration()
					d = docker.NewMock()

					d.RemoveImageHandler = func(name string) error {
						return e
					}

				})
				JustBeforeEach(func() {
					err = Publish(d, ls, c)
				})
				It("should return a remove error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(e))
				})
				It("should push an image on the registry", func() {
					Expect(d.TagCalled).To(BeTrue())
					Expect(d.PushCalled).To(BeTrue())
				})
				It("should try to remove the registry tag", func() {
					Expect(d.RemoveImageCalled).To(BeTrue())
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
			Name: "debian:8.2",
		},
		Publish: configuration.Publish{
			Hostname: "localhost:5000",
		},
	}
}
