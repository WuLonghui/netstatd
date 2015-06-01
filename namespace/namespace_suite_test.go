package namespace_test

import (
	. "netstatd/test_helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNamespace(t *testing.T) {
	RegisterFailHandler(Fail)
	InitDockerClient()
	RunSpecs(t, "Namespace Suite")
}
