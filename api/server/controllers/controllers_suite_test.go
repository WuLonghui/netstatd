package controllers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "netstatd"
	"netstatd/api/server/controllers"
	_ "netstatd/api/server/routers"
	. "netstatd/test_helper"

	"testing"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

var _ = BeforeSuite(func() {
	InitDockerClient()
	netstatd := NewNetstatd()
	controllers.Init(netstatd)
})
