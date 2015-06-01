package discovery_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "netstatd/test_helper"

	"testing"
)

func TestDiscovery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Discovery Suite")
}

var _ = BeforeSuite(func() {
	InitDockerClient()
})
