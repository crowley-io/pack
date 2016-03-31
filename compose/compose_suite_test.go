package compose_test

import (
	"github.com/crowley-io/macchiato"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"os"
	"testing"
)

func TestInstall(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	macchiato.RunSpecs(t, "Compose Suite")
}

var (
	pwd string
)

var _ = ginkgo.BeforeSuite(func() {
	p, err := os.Getwd()
	if err != nil {
		ginkgo.Fail(err.Error())
	}
	pwd = p
})
