package netstatd

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

type Netstatd struct {
	NS map[int]*NetstatdInNS
}

type NetstatdInNS struct {
	N         *Namespace
	Direction string
	NetStats  map[string]*NetStat
}

func NewNetstatd() *Netstatd {
	ns := make(map[int]*NetstatdInNS)
	ns[CURRENT_NS_PID] = &NetstatdInNS{
		N:        NewNamespace(CURRENT_NS_PID),
		NetStats: make(map[string]*NetStat),
	}

	return &Netstatd{
		NS: ns,
	}
}

func (d Netstatd) Run() error {
	return d.NS[CURRENT_NS_PID].stat()
}

func (d Netstatd) AddNameSpaceStat(pid int) error {
	n := NewNamespace(pid)
	if !n.Exist() {
		return fmt.Errorf("namespace not found")
	}

	d.NS[pid] = &NetstatdInNS{
		N:        n,
		NetStats: make(map[string]*NetStat),
	}

	return d.NS[pid].stat()
}

func (d NetstatdInNS) stat() error {
	log.Printf("starting stat in namespace(%v)", d.N)

	err := d.N.Set()
	if err != nil {
		log.Printf("error setting namespace, %v", err)
		return err
	}

	ifs, err := d.findDevs("nflog", "nfqueue", "any", "lo")
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

	return nil
}

func (d NetstatdInNS) findDevs(exceptions ...string) ([]pcap.Interface, error) {
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
	neStat := NewNetStat(iface)
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
