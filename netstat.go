package netstatd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/gopacket/pcap"
	metrics "github.com/rcrowley/go-metrics"
	. "netstatd/namespace"
)

type NetTcpStat struct {
	packets metrics.Meter
	bytes   metrics.Meter
}

type NetHttpStat struct {
	requests metrics.Meter
}

type NetStat struct {
	namespace *Namespace
	iface     pcap.Interface

	tcp  NetTcpStat
	http NetHttpStat
}

func NewNetStat(namespace *Namespace, iface pcap.Interface) *NetStat {
	return &NetStat{
		namespace: namespace,
		iface:     iface,

		tcp: NetTcpStat{
			packets: metrics.NewMeter(),
			bytes:   metrics.NewMeter(),
		},
		http: NetHttpStat{
			requests: metrics.NewMeter(),
		},
	}
}

func (s *NetStat) MarkTcp(bs []byte) {
	s.tcp.packets.Mark(1)
	s.tcp.bytes.Mark(int64(len(bs)))
}

func (s *NetStat) MarkHttpRequest(req *http.Request) {
	s.http.requests.Mark(1)
}

type NetStatJson struct {
	Namespace *Namespace `json:"namespace"`
	Iface     struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Addresses   []string `json:"addresses"`
	} `json:"interface"`

	Tcp struct {
		Packets       int64   `json:"packets"`
		PacketsPerSec float64 `json:"packets_per_sec"`

		Bytes       int64   `json:"bytes"`
		BytesPerSec float64 `json:"bytes_per_sec"`
	} `json:"tcp"`

	Http struct {
		Requests       int64   `json:"requests"`
		RequestsPerSec float64 `json:"requests_per_sec"`
	} `json:"http"`
}

func (s *NetStat) MarshalJSON() ([]byte, error) {
	j := NetStatJson{}

	j.Namespace = s.namespace

	j.Iface.Name = s.iface.Name
	j.Iface.Description = s.iface.Description
	j.Iface.Addresses = make([]string, 0)
	for _, address := range s.iface.Addresses {
		netMaskSize, _ := address.Netmask.Size()
		j.Iface.Addresses = append(j.Iface.Addresses, fmt.Sprintf("%s/%d", address.IP.String(), netMaskSize))
	}

	j.Tcp.Packets = s.tcp.packets.Count()
	j.Tcp.PacketsPerSec = s.tcp.packets.Rate1()
	j.Tcp.Bytes = s.tcp.bytes.Count()
	j.Tcp.BytesPerSec = s.tcp.bytes.Rate1()

	j.Http.Requests = s.http.requests.Count()
	j.Http.RequestsPerSec = s.http.requests.Rate1()

	return json.Marshal(j)
}
