package netstatd

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

type Netstatd struct {
	*NetStats
}

func NewNetstatd() *Netstatd {
	return &Netstatd{
		NetStats: NewNetStats(),
	}
}

func (d Netstatd) FindDevs(exceptions ...string) ([]string, error) {
	allDevs, err := pcap.FindAllDevs()
	if err != nil {
		return []string{}, err
	}

	devs := make([]string, 0)
	for _, dev := range allDevs {
		skip := false
		for _, exception := range exceptions {
			if dev.Name == exception {
				skip = true
				break
			}
		}
		if !skip {
			devs = append(devs, dev.Name)
		}
	}
	return devs, nil
}

func (d Netstatd) Run() {
	log.Print("netstartd: starting run")

	devs, err := d.FindDevs("nflog", "nfqueue", "any", "lo")
	if err != nil {
		log.Printf("netstartd: fail to find devs, %v", err)
	}

	for _, dev := range devs {
		go d.runStatByDev(dev)
	}
}

func (d Netstatd) runStatByDev(dev string) {
	log.Printf("netstartd: starting stat on interface %s", dev)

	snaplen := 65536
	filter := "tcp and port not 22"

	handle, err := pcap.OpenLive(dev, int32(snaplen), true, pcap.BlockForever)
	if err != nil {
		log.Printf("netstartd: error opening pcap handle, %v", err)
		return
	}

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Printf("netstartd: error setting BPF filter, %v", err)
		return
	}

	// Set up assembly
	streamFactory := &statsStreamFactory{
		NetStats: d.NetStats,
	}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	for {
		select {
		case packet := <-packets:
			if packet == nil {
				continue
			}

			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
				continue
			}

			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		}
	}
}
