package install_test

import (
	"github.com/crowley-io/macchiato"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"os"
	"os/user"
	"testing"
)

func TestInstall(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	macchiato.RunSpecs(t, "Install Suite")
}

var (
	uid  string
	gid  string
	home string
	pwd  string
)

var _ = ginkgo.BeforeSuite(func() {

	u, err := user.Current()
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	p, err := os.Getwd()
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	uid = u.Uid
	gid = u.Gid
	home = u.HomeDir
	pwd = p

})
