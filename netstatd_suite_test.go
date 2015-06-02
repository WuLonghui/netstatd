package netstatd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNetstatd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Netstatd Suite")
}
