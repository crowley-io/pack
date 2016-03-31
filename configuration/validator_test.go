package configuration_test

import (
	//"fmt"
	//"os"

	. "github.com/crowley-io/pack/configuration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	Describe("validate configuration attributes", func() {

		var (
			c   *Configuration
			err error
		)

		JustBeforeEach(func() {
			err = c.Validate()
		})
		Context("when it has no flaw", func() {
			Context("with a lambda configuration", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
				})
				It("should succeed", func() {
					Expect(err).To(Succeed())
				})
			})
			Context("with install disabled", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					c.Install.Disable = true
					c.Install.Image = ""
					c.Install.Output = ""
					c.Install.Path = ""
				})
				It("should succeed", func() {
					Expect(err).To(Succeed())
				})
				It("should be disabled", func() {
					Expect(c.Install.Disable).To(BeTrue())
				})
				It("could have empty attributes for install", func() {
					Expect(c.Install.Image).To(BeEmpty())
					Expect(c.Install.Output).To(BeEmpty())
					Expect(c.Install.Path).To(BeEmpty())
				})
			})
		})
		Context("when it has errors", func() {
			Context("because it's empty", func() {
				BeforeEach(func() {
					c = New()
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
			Context("because it's nil", func() {
				BeforeEach(func() {
					c = nil
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(ErrConfigurationEmpty))
				})
			})
			Context("because no output is defined", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					c.Install.Output = ""
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(ErrOutputRequired))
				})
			})
			Context("because no image is defined", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					c.Install.Image = ""
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(ErrImageRequired))
				})
			})
			Context("because no path is defined", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					c.Install.Path = ""
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(ErrPathRequired))
				})
			})
			Context("because no name is defined", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					c.Compose.Name = ""
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(ErrNameRequired))
				})
			})
			Context("because no hostname is defined", func() {
				BeforeEach(func() {
					c = getDefaultConfiguration()
					c.Publish.Hostname = ""
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(ErrHostnameRequired))
				})
			})
		})
	})
})

func getDefaultConfiguration() *Configuration {

	c := New()

	c.Install = Install{
		Image:  "debian",
		Path:   "/root",
		Output: "file",
	}

	c.Compose = Compose{
		Name: "app",
	}

	c.Publish = Publish{
		Hostname: "localhost:5000",
	}

	return c
}
