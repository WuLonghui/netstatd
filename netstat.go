package netstatd

import (
	"encoding/json"
	"net/http"

	"github.com/google/gopacket/pcap"
	metrics "github.com/rcrowley/go-metrics"
)

type NetTcpStat struct {
	packets metrics.Meter
	bytes   metrics.Meter
}

type NetHttpStat struct {
	requests metrics.Meter
}

type NetStat struct {
	iface pcap.Interface

	tcp  NetTcpStat
	http NetHttpStat
}

func NewNetStat(iface pcap.Interface) *NetStat {
	return &NetStat{
		iface: iface,

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

func (s *NetStat) MarshalJSON() ([]byte, error) {
	type netStat struct {
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

	j := netStat{}
	j.Tcp.Packets = s.tcp.packets.Count()
	j.Tcp.PacketsPerSec = s.tcp.packets.Rate1()
	j.Tcp.Bytes = s.tcp.bytes.Count()
	j.Tcp.BytesPerSec = s.tcp.bytes.Rate1()

	j.Http.Requests = s.http.requests.Count()
	j.Http.RequestsPerSec = s.http.requests.Rate1()

	return json.Marshal(j)
}
