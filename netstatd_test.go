package netstatd_test

import (
	. "netstatd"
	. "netstatd/namespace"

	"github.com/google/gopacket/pcap"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Netstatd", func() {
	var (
		d *Netstatd
	)

	BeforeEach(func() {
		d = NewNetstatd()
	})

	AfterEach(func() {
		d.Stop()
	})

	Describe("AddNameSpace", func() {
		It("adds namespace to gather statistics", func() {
			namespace := NewNamespace(CURRENT_NAMESPACE_PID, "host")
			err := d.AddNameSpace(namespace)
			Expect(err).ShouldNot(HaveOccurred())

			netStats := d.GetNetStats(namespace)
			Expect(len(netStats)).Should(BeNumerically(">", 0))

			allNetStats := d.GetAllNetStats()
			Expect(len(allNetStats)).Should(BeNumerically(">", 0))
		})

		It("returns an error when namespace doesn't exist", func() {
			namespace := NewNamespace(1234, "unknown")
			err := d.AddNameSpace(namespace)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Run", func() {
		It("should be ok", func() {
			err := d.Run("host,docker")
			Expect(err).ShouldNot(HaveOccurred())

			namespace := NewNamespace(CURRENT_NAMESPACE_PID, "host")
			netStats := d.GetNetStats(namespace)
			Expect(len(netStats)).Should(BeNumerically(">", 0))

			allNetStats := d.GetAllNetStats()
			Expect(len(allNetStats)).Should(BeNumerically(">", 0))
		})

	})
})

var _ = Describe("NetstatdInNS", func() {
	var (
		dn *NetstatdInNS
	)

	BeforeEach(func() {
		namespace := NewNamespace(CURRENT_NAMESPACE_PID, "host")
		dn = NewNetstatdInNS(namespace)
	})

	AfterEach(func() {
		dn.Stop()
	})

	Describe("FindDevs", func() {
		It("returns all devs", func() {
			devs, err := dn.FindDevs()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(devs)).Should(BeNumerically(">", 0))
		})
	})

	Describe("Capture", func() {
		It("can stop", func() {
			eth0 := pcap.Interface{
				Name: "eth0",
			}
			stop, err := dn.Capture(eth0)
			Expect(err).ShouldNot(HaveOccurred())
			stop <- true
		})
	})

	Describe("Run", func() {
		It("should be ok", func() {
			err := dn.Run()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dn.Running).Should(BeTrue())

			dn.Stop()
			Expect(dn.Running).Should(BeFalse())
		})
	})
})
