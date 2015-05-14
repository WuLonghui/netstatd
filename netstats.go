package netstatd

import (
	"encoding/json"

	metrics "github.com/rcrowley/go-metrics"
)

type NetStats struct {
	bytes   metrics.Meter
	packets metrics.Meter
}

func NewNetStats() *NetStats {
	return &NetStats{
		bytes:   metrics.NewMeter(),
		packets: metrics.NewMeter(),
	}
}

func (s *NetStats) Mark(bytes []byte) {
	s.bytes.Mark(int64(len(bytes)))
	s.packets.Mark(1)
}

func (s *NetStats) MarshalJSON() ([]byte, error) {
	type netStats struct {
		Bytes       int64   `json:"bytes"`
		BytesPerSec float64 `json:"bytes_per_second"`

		Packets       int64   `json:"packets"`
		PacketsPerSec float64 `json:"packets_per_second"`
	}

	j := netStats{}
	j.Bytes = s.bytes.Count()
	j.BytesPerSec = s.bytes.Rate1()
	j.Packets = s.packets.Count()
	j.PacketsPerSec = s.packets.Rate1()
	return json.Marshal(j)
}
