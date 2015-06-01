package discovery_test

import (
	. "netstatd/namespace/discovery"
	. "netstatd/test_helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DockerDiscovery", func() {
	var dockerDiscovery *DockerDiscovery

	BeforeEach(func() {
		dockerDiscovery = NewDockerDiscovery()
		DockerCleanAllContainers()
	})

	AfterEach(func() {
		DockerCleanAllContainers()
	})

	It("gets container's namespace", func() {
		container := DockerRunContainer()
		_, err := dockerDiscovery.GetNamespace(container.ID[0:11])
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("list all namespaces", func() {
		for i := 0; i < 2; i++ {
			DockerRunContainer()
		}

		namespaces := dockerDiscovery.ListAllNamespaces()
		Expect(len(namespaces)).Should(Equal(2))
	})
})
