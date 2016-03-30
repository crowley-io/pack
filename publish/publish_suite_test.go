package publish_test

import (
	"github.com/crowley-io/macchiato"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"testing"
)

func TestPublish(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	macchiato.RunSpecs(t, "Publish Suite")
}
