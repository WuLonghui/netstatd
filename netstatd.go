package netstatd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"

	. "netstatd/namespace"
	. "netstatd/namespace/discovery"
)

type Netstatd struct {
	NS map[int]*NetstatdInNS
}

func NewNetstatd() *Netstatd {
	ns := make(map[int]*NetstatdInNS)
	return &Netstatd{
		NS: ns,
	}
}

func (d Netstatd) Run(statTarget string) error {
	if strings.Contains(statTarget, "host") {
		namespace := NewNamespace(CURRENT_NAMESPACE_PID, "host")
		err := d.AddNameSpace(namespace)
		if err != nil {
			return err
		}
	}

	if strings.Contains(statTarget, "docker") {
		dockerDiscovery := NewDockerDiscovery()
		go func() {
			ticker := time.NewTicker(time.Second * 1)
			defer ticker.Stop()
			for {
				<-ticker.C
				namespaces := dockerDiscovery.ListAllNamespaces()
				for _, namespace := range namespaces {
					err := d.AddNameSpace(namespace)
					if err != nil {
						log.Printf("error adding namespace, %v", err)
					}
				}
			}
		}()
	}

	return nil
}

func (d Netstatd) Stop() {
	for _, n := range d.NS {
		n.Stop()
	}
}

func (d Netstatd) AddNameSpace(n *Namespace) error {
	if !n.Exist() {
		return fmt.Errorf("namespace not found")
	}

	if _, ok := d.NS[n.Pid]; ok {
		return nil
	}

	d.NS[n.Pid] = &NetstatdInNS{
		N:        n,
		NetStats: make(map[string]*NetStat),
		running:  false,
	}

	return d.NS[n.Pid].Run()
}

func (d Netstatd) GetNetStats(namespace *Namespace) map[string]*NetStat {
	n, ok := d.NS[namespace.Pid]
	if !ok {
		return make(map[string]*NetStat)
	}
	return n.NetStats
}

func (d Netstatd) GetAllNetStats() []*NetStat {
	netStats := make([]*NetStat, 0)
	for _, n := range d.NS {
		for _, netStat := range n.NetStats {
			netStats = append(netStats, netStat)
		}
	}

	return netStats
}

type NetstatdInNS struct {
	N         *Namespace
	Direction string
	NetStats  map[string]*NetStat

	running bool
}

func (d NetstatdInNS) Run() error {
	if d.running {
		return nil
	}

	log.Printf("starting run in namespace(%v)", d.N)

	err := d.N.Set()
	if err != nil {
		log.Printf("error setting namespace, %v", err)
		return err
	}

	ifs, err := d.FindDevs("nflog", "nfqueue", "any", "lo")
	if err != nil {
		log.Printf("error finding devs, %v", err)
		return err
	}

	for _, iface := range ifs {
		err = d.capture(iface)
		if err != nil {
			log.Printf("error capturing on interface, %v", err)
			return err
		}
	}

	d.running = true
	return nil
}

func (d NetstatdInNS) Stop() {

}

func (d NetstatdInNS) FindDevs(exceptions ...string) ([]pcap.Interface, error) {
	ifs := make([]pcap.Interface, 0)
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return ifs, err
	}

	for _, dev := range devs {
		skip := false
		for _, exception := range exceptions {
			if dev.Name == exception {
				skip = true
				break
			}
		}
		if !skip {
			ifs = append(ifs, dev)
		}
	}
	return ifs, nil
}

func (d NetstatdInNS) capture(iface pcap.Interface) error {
	log.Printf("starting capture on interface %s", iface.Name)

	snaplen := 65536
	filter := "tcp and port not 22"

	handle, err := pcap.OpenLive(iface.Name, int32(snaplen), true, pcap.BlockForever)
	if err != nil {
		log.Printf("error opening pcap handle, %v", err)
		return err
	}

	err = handle.SetDirection(pcap.DirectionIn)
	if err != nil {
		log.Printf("error setting direction, %v", err)
		return err
	}

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Printf("error setting BPF filter, %v", err)
		return err
	}

	// Set up assembly
	neStat := NewNetStat(d.N, iface)
	d.NetStats[iface.Name] = neStat
	streamFactory := &statsStreamFactory{
		netStat: neStat,
	}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	go func() {
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
	}()

	return nil
}
