package namespace_test

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/google/gopacket/pcap"
	. "netstatd/namespace"
	. "netstatd/test_helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Namespace", func() {

	Context("when namespace doesn't exist", func() {
		It("cann't go to the net namespace", func() {
			namespace := NewNamespace(1234, "host")
			Expect(namespace.Exist()).Should(BeFalse())
			err := namespace.Set()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Current Namespace", func() {
		It("should be ok", func() {
			namespace := NewNamespace(CURRENT_NAMESPACE_PID, "host")
			Expect(namespace.Exist()).Should(BeTrue())
			err := namespace.Set()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("Docker Namespace", func() {
		var container *docker.Container

		BeforeEach(func() {
			container = DockerRunContainer()
		})

		It("can goto the docker's  namespace", func() {
			namespace := NewNamespace(container.State.Pid, "docker")
			Expect(namespace.Exist()).Should(BeTrue())
			err := namespace.Set()
			Expect(err).ShouldNot(HaveOccurred())

			devs, err := pcap.FindAllDevs()
			var eth0 pcap.Interface
			for _, dev := range devs {
				if dev.Name == "eth0" {
					eth0 = dev
				}
			}
			Expect(container.NetworkSettings.IPAddress).Should(Equal(eth0.Addresses[0].IP.String()))
		})
	})
})
