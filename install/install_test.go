package install_test

import (
	"errors"
	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"

	. "github.com/crowley-io/pack/install"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Install", func() {

	var (
		c   *configuration.Configuration
		ls  = docker.NewLogStream()
		d   *docker.Mock
		err error
	)

	Describe("create an artefact from a docker container", func() {
		Context("when the configuration is valid", func() {
			BeforeEach(func() {

				c = getDefaultConfiguration()
				d = docker.NewMock()

			})
			JustBeforeEach(func() {
				err = Install(d, ls, c)
			})
			It("should succeed", func() {
				Expect(err).To(Succeed())
			})
			It("should run a container", func() {
				Expect(d.RunCalled).To(BeTrue())
				Expect(d.RunOptions.Command).To(Equal("make"))
				Expect(d.RunOptions.Image).To(Equal("debian"))
				Expect(d.RunOptions.Env).To(ConsistOf(
					"CROWLEY_PACK_DIRECTORY=/root",
					"CROWLEY_PACK_OUTPUT=/testdata/app.bin",
					ContainSubstring("CROWLEY_PACK_USER="),
					ContainSubstring("CROWLEY_PACK_GROUP="),
				))
				Expect(d.RunOptions.Links).To(BeEmpty())
				Expect(d.RunOptions.Volumes).To(ConsistOf(
					ContainSubstring("src/github.com/crowley-io/pack/install:/root:rw"),
				))
			})
		})
		Context("when install in configuration is disabled", func() {
			BeforeEach(func() {

				c = getDefaultConfiguration()
				c.Install.Disable = true
				d = docker.NewMock()

			})
			JustBeforeEach(func() {
				err = Install(d, ls, c)
			})
			It("should succeed", func() {
				Expect(err).To(Succeed())
			})
			It("should not run a container", func() {
				Expect(d.RunCalled).To(BeFalse())
			})
		})
		Context("when the configuration is invalid", func() {
			Context("because its empty", func() {
				BeforeEach(func() {

					c = &configuration.Configuration{}
					d = docker.NewMock()

				})
				JustBeforeEach(func() {
					err = Install(d, ls, c)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
				It("should not run a container", func() {
					Expect(d.RunCalled).To(BeFalse())
				})
			})
			Context("because volumes configuration has a syntax error", func() {
				BeforeEach(func() {

					c = getDefaultConfiguration()
					c.Install.Volumes = []string{"/home/user/.npm:/root/.npm:rw:3"}
					d = docker.NewMock()

				})
				JustBeforeEach(func() {
					err = Install(d, ls, c)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
				It("should not run a container", func() {
					Expect(d.RunCalled).To(BeFalse())
				})
			})
		})
		Context("when docker have issue(s)", func() {
			Context("because we cannot run a container", func() {

				var (
					e = errors.New("error: socket timeout")
				)

				BeforeEach(func() {

					c = getDefaultConfiguration()
					d = docker.NewMock()
					d.RunHandler = func(option docker.RunOptions, stream docker.LogStream) (int, error) {
						return 0, e
					}

				})
				JustBeforeEach(func() {
					err = Install(d, ls, c)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(e))
				})
				It("should try to run a container", func() {
					Expect(d.RunCalled).To(BeTrue())
				})
			})
			Context("because the container exit with an error", func() {

				var (
					e = 255
				)

				BeforeEach(func() {

					c = getDefaultConfiguration()
					d = docker.NewMock()
					d.RunHandler = func(option docker.RunOptions, stream docker.LogStream) (int, error) {
						return e, nil
					}

				})
				JustBeforeEach(func() {
					err = Install(d, ls, c)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("exit status %d", e))
				})
				It("should try to run a container", func() {
					Expect(d.RunCalled).To(BeTrue())
				})
			})
			Context("because the container doesn't create an output", func() {

				BeforeEach(func() {

					c = getDefaultConfiguration()
					c.Install.Output = "file.txt"
					d = docker.NewMock()

				})
				JustBeforeEach(func() {
					err = Install(d, ls, c)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("file not found"))
				})
				It("should try to run a container", func() {
					Expect(d.RunCalled).To(BeTrue())
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
			Output:  "../testdata/app.bin",
		},
		Compose: configuration.Compose{
			Name: "debian:8.2",
		},
		Publish: configuration.Publish{
			Hostname: "localhost:5000",
		},
	}
}
