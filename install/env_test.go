package install_test

import (
	"fmt"
	"github.com/crowley-io/pack/configuration"
	"os"

	. "github.com/crowley-io/pack/install"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Env", func() {

	var (
		p   = "/media/app"
		o   = "libshaped.so"
		gp  = "GOPATH=$PATH:/usr/local/go"
		sk  = "SECRET_KEY='3ztP7$Xqoef=VUdPa'"
		db  = "DB_URI=mongodb://user:password@host:27017/db"
		fb  = "FOO:BAR"
		c   *configuration.Configuration
		env []string
	)

	Describe("create a list of environment variable", func() {
		BeforeEach(func() {
			c = &configuration.Configuration{
				Install: configuration.Install{
					Path:        p,
					Output:      o,
					Environment: []string{sk, db, gp, fb},
				},
			}
		})
		JustBeforeEach(func() {
			env = GetEnv(c)
		})
		It("should add user's UID and GID", func() {
			Expect(env).To(ContainElement(fmt.Sprintf("CROWLEY_PACK_USER=%s", uid)))
			Expect(env).To(ContainElement(fmt.Sprintf("CROWLEY_PACK_GROUP=%s", gid)))
		})
		It("should add directory and output path", func() {
			Expect(env).To(ContainElement(fmt.Sprintf("CROWLEY_PACK_DIRECTORY=%s", p)))
			Expect(env).To(ContainElement(fmt.Sprintf("CROWLEY_PACK_OUTPUT=%s/%s", p, o)))
		})
		It("should resolve variables", func() {
			path := os.Getenv("PATH")
			Expect(env).To(ContainElement(fmt.Sprintf("GOPATH=%s:/usr/local/go", path)))
		})
		It("should not resolve escaped variables", func() {
			Expect(env).To(ContainElement("SECRET_KEY=3ztP7$Xqoef=VUdPa"))
		})
		It("should add a variable", func() {
			Expect(env).To(ContainElement(db))
		})
		It("should not filter invalid variable", func() {
			Expect(env).To(ContainElement(fb))
		})
	})
})
