package configuration_test

import (
	"fmt"
	"os"

	. "github.com/crowley-io/pack/configuration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("parse a file and inflate a new configuration", func() {

		var (
			p   string
			c   *Configuration
			err error
		)

		JustBeforeEach(func() {
			c, err = Parse(p)
		})
		Context("when configuration is valid", func() {

			AssertValidConfiguration := func() {
				It("should succeed", func() {
					Expect(err).To(Succeed())
				})
				It("should inflate a valid configuration", func() {
					Expect(c.Install.Output).To(Equal("app.tar.gz"))
					Expect(c.Install.Image).To(Equal("packer-go"))
					Expect(c.Install.Path).To(Equal("/usr/local/go"))
					Expect(c.Compose.Name).To(Equal("shen"))
					Expect(c.Publish.Hostname).To(Equal("localhost:5000"))
				})
			}

			Context("with a yaml file", func() {
				BeforeEach(func() {
					p = getConfigurationFilepath("packer.yml")
				})
				AssertValidConfiguration()
			})
			Context("with a json file", func() {
				BeforeEach(func() {
					p = getConfigurationFilepath("packer.json")
				})
				AssertValidConfiguration()
			})
		})
		Context("when configuration has an error", func() {
			Context("with its syntax", func() {
				BeforeEach(func() {
					p = getConfigurationFilepath("packer-seq.yml")
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
			Context("with its attributes", func() {
				BeforeEach(func() {
					p = getConfigurationFilepath("packer-nooutput.yml")
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(ErrOutputRequired))
				})
			})
		})
		Context("file does not exist", func() {
			BeforeEach(func() {
				p = getConfigurationFilepath("nofile.yml")
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})
	})
})

func getConfigurationFilepath(name string) string {
	return fmt.Sprintf("../testdata/%s", name)
}
