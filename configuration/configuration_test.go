package configuration_test

import (
	//"fmt"
	//"os"

	. "github.com/crowley-io/pack/configuration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {
	Describe("define parameters for docker client and modules", func() {

		var (
			c   *Configuration
			err error
		)

		Describe("a new configuration", func() {
			Context("without any specific parameters", func() {
				BeforeEach(func() {
					c = New()
					err = c.Validate()
				})
				It("should not be nil", func() {
					Expect(c).NotTo(BeNil())
				})
				It("should return an error with validate", func() {
					Expect(err).To(HaveOccurred())
				})
			})
			Context("with a remote configuration", func() {
				Context("which is valid", func() {
					Context("the parser", func() {
						BeforeEach(func() {
							c = New()
							err = c.Configure("localhost:5000/app:1.2")
						})
						It("should succeed", func() {
							Expect(err).To(Succeed())
						})
						It("should configure compose parameters", func() {
							Expect(c.Compose.Name).To(Equal("app:1.2"))
						})
						It("should configure publish parameters", func() {
							Expect(c.Publish.Hostname).To(Equal("localhost:5000"))
						})
					})
					Context("the validator", func() {
						BeforeEach(func() {
							c = New()
							_ = c.Configure("localhost:5000/app:1.2")
							err = c.Validate()
						})
						It("should return an error with validate", func() {
							Expect(err).To(HaveOccurred())
						})
					})
				})
				Context("which is invalid", func() {
					BeforeEach(func() {
						c = New()
						err = c.Configure("ftp://localhost:5000/app:1.2")
					})
					It("should return an error", func() {
						Expect(err).To(HaveOccurred())
					})
				})
			})
			Context("with setter function", func() {
				Context("when cache is enabled", func() {
					BeforeEach(func() {
						c = New()
						c.EnableCache()
					})
					It("should enable cache", func() {
						Expect(c.Compose.NoCache).To(BeFalse())
					})
				})
				Context("when cache is disabled", func() {
					BeforeEach(func() {
						c = New()
						c.DisableCache()
					})
					It("should disable cache", func() {
						Expect(c.Compose.NoCache).To(BeTrue())
					})
				})
				Context("when pull is enabled", func() {
					BeforeEach(func() {
						c = New()
						c.EnablePull()
					})
					It("should enable pull", func() {
						Expect(c.Compose.Pull).To(BeTrue())
					})
				})
				Context("when pull is disabled", func() {
					BeforeEach(func() {
						c = New()
						c.DisablePull()
					})
					It("should disable pull", func() {
						Expect(c.Compose.Pull).To(BeFalse())
					})
				})
				Context("when install is enabled", func() {
					BeforeEach(func() {
						c = New()
						c.EnableInstall()
					})
					It("should enable install", func() {
						Expect(c.Install.Disable).To(BeFalse())
					})
				})
				Context("when install is disabled", func() {
					BeforeEach(func() {
						c = New()
						c.DisableInstall()
					})
					It("should disable install", func() {
						Expect(c.Install.Disable).To(BeTrue())
					})
				})
			})
		})
	})
})
